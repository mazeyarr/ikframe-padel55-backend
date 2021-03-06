package auth

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

const OAuthToken = "OAuthToken"
const IdTokenName = "FIREBASE_ID_TOKEN"

// FirebaseAuthMiddleware is middleware for Firebase Authentication
type FirebaseAuthMiddleware struct {
	cli          *auth.Client
	unAuthorized func(c *gin.Context)
}

func New(credFileName string, unAuthorized func(c *gin.Context)) (*FirebaseAuthMiddleware, error) {
	opt := option.WithCredentialsFile(credFileName)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	authCli, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseAuthMiddleware{
		cli:          authCli,
		unAuthorized: unAuthorized,
	}, nil
}

func (fam *FirebaseAuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		idToken, err := fam.cli.VerifyIDToken(context.Background(), token)
		if err != nil {
			if fam.unAuthorized != nil {
				fam.unAuthorized(c)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": http.StatusText(http.StatusUnauthorized),
				})

				c.Abort()
			}
			return
		}

		c.Set(OAuthToken, token)
		c.Set(IdTokenName, idToken)
		c.Next()
	}
}

func ExtractClaims(c *gin.Context) *auth.Token {
	idToken, ok := c.Get(IdTokenName)
	if !ok {
		return new(auth.Token)
	}

	return idToken.(*auth.Token)
}
