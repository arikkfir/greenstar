package auth

import (
	"context"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(_ context.Context) error {
	return nil
}

func CreateJWTValidationMiddleware(auth0Domain string, audiences []string) func(c *gin.Context) {
	issuerURL, err := url.Parse("https://" + auth0Domain + "/")
	if err != nil {
		panic(fmt.Errorf("failed to parse issuer URL: %w", err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audiences,
		validator.WithCustomClaims(func() validator.CustomClaims { return &CustomClaims{} }),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("failed to set up a JWT validator: %w", err))
	}

	return func(c *gin.Context) {
		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to validate JWT: %w", err))
			}),
		)

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validatedClaims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			registeredClaims := validatedClaims.RegisteredClaims
			customClaims := validatedClaims.CustomClaims.(*CustomClaims)
			claims := Claims{
				RegisteredClaims: registeredClaims,
				CustomClaims:     *customClaims,
			}
			c.Set(claimsContextKey, &claims)
			c.Next()
		})
		handler := middleware.CheckJWT(next)
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
