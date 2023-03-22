package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/redisutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
	"time"
)

type Authenticator struct {
	Config              *Config
	SecureCookies       bool
	DefaultPostLoginURL string
	OAuth               *oauth2.Config
}

func (a *Authenticator) createTokenVerificationCallback(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	} else {
		return []byte(a.Config.Google.ClientSecret), nil
	}
}

func (a *Authenticator) parseClaimsFromCookie(c *gin.Context, expectedAudience string) (*jwt.RegisteredClaims, error) {
	claimsCookieValue, err := c.Cookie(a.Config.ClaimsCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			log.Ctx(c).Debug().Str("cookieName", a.Config.ClaimsCookieName).Msg("no claims cookie found")
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed getting cookie: %w", err)
		}
	}

	claims := jwt.RegisteredClaims{}
	if _, err := jwt.ParseWithClaims(claimsCookieValue, &claims, a.createTokenVerificationCallback); err != nil {
		return nil, fmt.Errorf("failed parsing JWT token '%s': %w", claimsCookieValue, err)
	} else if !claims.VerifyAudience(expectedAudience, true) {
		return nil, fmt.Errorf("invalid audience: %s", claims.Audience)
	} else if !claims.VerifyExpiresAt(time.Now(), true) {
		return nil, fmt.Errorf("token expired at: %s", claims.ExpiresAt)
	} else if !claims.VerifyIssuedAt(time.Now(), true) {
		return nil, fmt.Errorf("token issued at: %s", claims.IssuedAt)
	} else if !claims.VerifyIssuer("greenstar.auth", true) {
		return nil, fmt.Errorf("invalid issuer: %s", claims.Issuer)
	} else if !claims.VerifyNotBefore(time.Now(), true) {
		return nil, fmt.Errorf("token not valid before: %s", claims.NotBefore)
	} else if strings.HasPrefix(claims.ID, "test:") {
		if err := a.createTestSession(c, &claims); err != nil {
			return nil, fmt.Errorf("failed creating test session: %w", err)
		}
	}
	return &claims, nil
}

func (a *Authenticator) parseOAuthStateFromCookie(c *gin.Context) (*string, error) {
	state, err := c.Cookie(a.Config.StateCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed getting OAuth state cookie: %w", err)
		}
	}
	return &state, nil
}

func (a *Authenticator) createTestSession(ctx context.Context, claims *jwt.RegisteredClaims) error {
	session := Session{}
	session.Claims = *claims
	session.Token = &oauth2.Token{
		AccessToken:  "bogus_access_token",
		TokenType:    "Bearer",
		RefreshToken: "bogus_refresh_token",
		Expiry:       time.Now().Add(15 * time.Minute),
	}
	session.Tenant = "test"
	session.Permissions = []string{
		"greenstar.auth.getUserInfo",
		"greenstar.auth.getUserInfo:mock",
		"greenstar.admin.createTenant",
	}
	session.MockUserInfo = &GoogleAPIUserInfoResponse{
		Email:        "jack@ryan.com",
		FamilyName:   "Ryan",
		GivenName:    "Jack",
		HostedDomain: "test",
		ID:           claims.Subject,
		Link:         "",
		Locale:       "en",
		Name:         "Jack Ryan",
		PictureURL:   "",
	}
	r := redisutil.GetRedis(ctx)
	if sessionBytes, err := json.Marshal(session); err != nil {
		return fmt.Errorf("failed marshalling session before persisting to Redis: %w", err)
	} else if result := r.Do(ctx, r.B().Set().Key("session:"+claims.ID).Value(string(sessionBytes)).ExatTimestamp(claims.ExpiresAt.Unix()).Build()); result.Error() != nil {
		return fmt.Errorf("failed persisting test session to redis: %w", err)
	} else {
		return nil
	}
}

func (a *Authenticator) loadSession(ctx context.Context, id string) (*Session, error) {
	r := redisutil.GetRedis(ctx)

	session := Session{}
	if result := r.Do(ctx, r.B().Get().Key("session:"+id).Build()); result.Error() != nil {
		if rueidis.IsRedisNil(result.Error()) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed getting session '%s' from redis: %w", id, result.Error())
		}
	} else if err := result.DecodeJSON(&session); err != nil {
		return nil, fmt.Errorf("failed decoding session '%s' from redis: %w", id, err)
	} else {
		return &session, nil
	}
}

func (a *Authenticator) verifyClaimsAndSession(claims *jwt.RegisteredClaims, session *Session) error {
	if session.Claims.Issuer != claims.Issuer {
		return fmt.Errorf("incorrect issuer in claims vs. session, expected '%s', got '%s'", session.Claims.Issuer, claims.ID)
	} else if session.Claims.Subject != claims.Subject {
		return fmt.Errorf("incorrect subject in claims vs. session, expected '%s', got '%s'", session.Claims.Subject, claims.Subject)
	} else if strings.Join(session.Claims.Audience, ",") != strings.Join(claims.Audience, ",") {
		return fmt.Errorf("incorrect audience in claims vs. session, expected '%s', got '%s'", session.Claims.Audience, claims.Audience)
	} else if session.Claims.ExpiresAt.UnixNano() != claims.ExpiresAt.UnixNano() {
		return fmt.Errorf("incorrect expiration in claims vs. session, expected '%s', got '%s'", session.Claims.ExpiresAt, claims.ExpiresAt)
	} else if session.Claims.NotBefore.UnixNano() != claims.NotBefore.UnixNano() {
		return fmt.Errorf("incorrect activation time in claims vs. session, expected '%s', got '%s'", session.Claims.NotBefore, claims.NotBefore)
	} else if session.Claims.IssuedAt.UnixNano() != claims.IssuedAt.UnixNano() {
		return fmt.Errorf("incorrect issuance time in claims vs. session, expected '%s', got '%s'", session.Claims.IssuedAt, claims.IssuedAt)
	} else if session.Claims.ID != claims.ID {
		return fmt.Errorf("incorrect ID in claims vs. session, expected '%s', got '%s'", session.Claims.ID, claims.ID)
	} else {
		return nil
	}
}

func (a *Authenticator) loadUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleAPIUserInfoResponse, error) {
	httpClient := a.OAuth.Client(ctx, token)
	defer httpClient.CloseIdleConnections()

	userInfoResponse, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer userInfoResponse.Body.Close()

	if userInfoResponse.StatusCode >= 200 && userInfoResponse.StatusCode <= 299 {
		userInfo := GoogleAPIUserInfoResponse{}
		decoder := json.NewDecoder(userInfoResponse.Body)
		if err := decoder.Decode(&userInfo); err != nil {
			return nil, fmt.Errorf("failed decoding userinfo response: %w", err)
		} else {
			return &userInfo, nil
		}
	} else {
		body, _ := io.ReadAll(userInfoResponse.Body)
		return nil, fmt.Errorf("userinfo request failed with status %d, body:\n%s", userInfoResponse.StatusCode, body)
	}
}
