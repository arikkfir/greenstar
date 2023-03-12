package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

const sessionContextKey = "$$$$____session____$$$$"

type Session struct {
	Claims      jwt.RegisteredClaims `json:"claims"`
	Token       *oauth2.Token        `json:"token"`
	Tenant      string               `json:"tenant"`
	Permissions []string             `json:"permissions"`
}

func GetSession(ctx context.Context) *Session {
	v := ctx.Value(sessionContextKey)
	if v == nil {
		panic("session not found in context")
	} else if session, ok := v.(*Session); ok {
		return session
	} else {
		panic(fmt.Sprintf("unexpected session object encountered: %+v", v))
	}
}

func setSession(c *gin.Context, session *Session) {
	c.Set(sessionContextKey, session)
}
