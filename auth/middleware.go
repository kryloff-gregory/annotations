package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"main/controller"
	"main/user"
)

var jwtKey = []byte("SECRET_KEY")

func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

type Middleware struct {
	userProvider *user.Provider
}

func NewMiddleware(userProvider *user.Provider) *Middleware {
	return &Middleware{userProvider: userProvider}
}

func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		requiredToken := c.Request.Header["Authorization"]

		if len(requiredToken) == 0 {
			responseWithError(c, 403, "Please login to your account")
		}

		userName, _ := getUsernameFromToken(requiredToken[0])
		usr := m.userProvider.GerUserByName(userName)

		if usr == nil {
			responseWithError(c, 404, "User account not found")
			return
		}

		c.Set("User", usr)
		c.Next()
	}
}

func getUsernameFromToken(tkStr string) (string, error) {
	claims := &controller.Claims{}

	tkn, err := jwt.ParseWithClaims(tkStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}
		return "", err
	}

	if !tkn.Valid {
		return "", err
	}

	return claims.Username, nil
}
