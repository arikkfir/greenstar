package main

import (
	"encoding/base64"
	auth "github.com/arikkfir/greenstar/backend/auth/pkg"
	"github.com/arikkfir/greenstar/backend/util/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

/*
	consent required when login claims missing refresh token
	consent not required when logged in
	post-login URL taken from URL query when provided
	default post-login URL used when not given in URL query
	offline access is requested
	offline access is requested
*/

var _ = Describe("Login", func() {
	_ = ghttp.GHTTPWithGomega{}

	const defaultPostLoginURL = "https://test.com"
	const oauthStateCookieName = "oauth-state-test"
	const claimsCookieName = "claims-test"

	var router *gin.Engine
	var oauthConfig *oauth2.Config
	var err error
	var req *http.Request
	var res *httptest.ResponseRecorder

	BeforeEach(func() {
		oauthConfig = &oauth2.Config{
			RedirectURL:  "https://test.com/google/callback",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Scopes: []string{
				"scope1",
				"scope2",
			},
			Endpoint: google.Endpoint,
		}
		router = ginutil.NewGin(true)
		router.Use(auth.CreateObtainTokenMiddleware(claimsCookieName, oauthConfig.ClientSecret, false))
		router.GET("/google/login", CreateAuthGoogleLoginHandler(oauthConfig, oauthStateCookieName, defaultPostLoginURL))

		req, err = http.NewRequest("GET", "/google/login", nil)
		Expect(err).ToNot(HaveOccurred())

		res = httptest.NewRecorder()
	})

	When("the user navigates to login URL", func() {
		var location string

		JustBeforeEach(func() {
			router.ServeHTTP(res, req)
		})

		AssertLoginRedirect := func() {
			Expect(res).To(HaveHTTPStatus(http.StatusTemporaryRedirect))

			location = res.Header().Get("Location")
			Expect(location).To(HavePrefix(oauthConfig.Endpoint.AuthURL))

			locationURL, err := url.Parse(location)
			Expect(err).ToNot(HaveOccurred())
			Expect(locationURL.Query().Get("access_type")).To(Equal("offline"))
			Expect(locationURL.Query().Get("client_id")).To(Equal(oauthConfig.ClientID))
			Expect(locationURL.Query().Get("scope")).To(Equal(strings.Join(oauthConfig.Scopes, " ")))
			Expect(locationURL.Query().Get("response_type")).To(Equal("code"))
			Expect(locationURL.Query().Get("redirect_uri")).To(Equal(oauthConfig.RedirectURL))

			state := locationURL.Query().Get("state")
			Expect(state).ToNot(BeEmpty())

			encodedNonce, encodedPostLoginURL, found := strings.Cut(state, "|")
			Expect(found).To(BeTrue())

			nonce, err := base64.URLEncoding.DecodeString(encodedNonce)
			Expect(err).ToNot(HaveOccurred())
			Expect(nonce).To(HaveLen(16))

			postLoginURL, err := base64.URLEncoding.DecodeString(encodedPostLoginURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(postLoginURL)).To(Equal(defaultPostLoginURL))
		}

		AssertConsentRequested := func() {
			locationURL, err := url.Parse(res.Header().Get("Location"))
			Expect(err).ToNot(HaveOccurred())
			Expect(locationURL.Query().Get("prompt")).To(Equal("consent"))
		}

		CreateToken := func(iat time.Time, accessToken, refreshToken string) string {
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Token{
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        uuid.NewString(),
					Issuer:    "greenstar.auth.test",
					Subject:   "google|1234567890",
					Audience:  []string{"greenstar.admin", "greenstar.auth", "greenstar.operations", "greenstar.public"},
					IssuedAt:  jwt.NewNumericDate(iat),
					NotBefore: jwt.NewNumericDate(iat),
					ExpiresAt: jwt.NewNumericDate(iat.Add(time.Hour)),
				},
				Tenant:       "test.com",
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			})
			signedToken, err := jwtToken.SignedString([]byte(oauthConfig.ClientSecret))
			Expect(err).ToNot(HaveOccurred())
			return signedToken
		}

		Context("and the user has no claims cookie", func() {
			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should request user consent", func(ctx SpecContext) {
				AssertConsentRequested()
			})
		})

		Context("and the user has an invalid claims cookie", func() {
			BeforeEach(func() {
				req.AddCookie(&http.Cookie{Name: claimsCookieName, Value: "invalid"})
			})

			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should request user consent", func(ctx SpecContext) {
				AssertConsentRequested()
			})

			It("should clear the claims cookie", func(ctx SpecContext) {
				Expect(res.Header().Get("Set-Cookie")).To(Equal(claimsCookieName + "=; Path=/; Max-Age=0; HttpOnly"))
			})
		})

		Context("and the user has an expired claims cookie", func() {
			BeforeEach(func() {
				req.AddCookie(&http.Cookie{
					Name:  claimsCookieName,
					Value: CreateToken(time.Now().Add(time.Hour*24*-7), "at", "rt"),
				})
			})

			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should request user consent", func(ctx SpecContext) {
				AssertConsentRequested()
			})

			It("should clear the cookie", func(ctx SpecContext) {
				Expect(res.Header().Get("Set-Cookie")).To(Equal(claimsCookieName + "=; Path=/; Max-Age=0; HttpOnly"))
			})
		})

		Context("and the user has a valid claims cookie", func() {
			BeforeEach(func() {
				req.AddCookie(&http.Cookie{
					Name:  claimsCookieName,
					Value: CreateToken(time.Now(), "at", "rt"),
				})
			})

			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should not clear the cookie", func(ctx SpecContext) {
				setCookieHeader := res.Header().Get("Set-Cookie")
				Expect(setCookieHeader).ToNot(ContainSubstring(claimsCookieName))
			})
		})

		Context("and the user has a valid claims cookie sans refresh token", func() {
			BeforeEach(func() {
				req.AddCookie(&http.Cookie{
					Name:  claimsCookieName,
					Value: CreateToken(time.Now(), "at", ""),
				})
			})

			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should request user consent", func(ctx SpecContext) {
				AssertConsentRequested()
			})

			It("should clear the cookie", func(ctx SpecContext) {
				Expect(res.Header().Get("Set-Cookie")).To(Equal(claimsCookieName + "=; Path=/; Max-Age=0; HttpOnly"))
			})
		})

		Context("and the user has a valid claims cookie sans access token", func() {
			BeforeEach(func() {
				req.AddCookie(&http.Cookie{
					Name:  claimsCookieName,
					Value: CreateToken(time.Now(), "", "rt"),
				})
			})

			It("should redirect to OAuth endpoint", func(ctx SpecContext) {
				AssertLoginRedirect()
			})

			It("should request user consent", func(ctx SpecContext) {
				AssertConsentRequested()
			})

			It("should clear the cookie", func(ctx SpecContext) {
				Expect(res.Header().Get("Set-Cookie")).To(Equal(claimsCookieName + "=; Path=/; Max-Age=0; HttpOnly"))
			})
		})
	})
})
