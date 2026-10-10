package middleware

import (
	"strings"

	"catatan-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "userID"

func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		raw := strings.TrimPrefix(header, "Bearer ")
		if raw == "" || raw == header {
			c.Abort()
			utils.Fail(c, 401, "token diperlukan")
			return
		}
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) { return secret, nil },
			jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !tok.Valid {
			c.Abort()
			utils.Fail(c, 401, "token tidak valid atau kedaluwarsa")
			return
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		sub, ok2 := claims["sub"].(float64)
		if !ok || !ok2 {
			c.Abort()
			utils.Fail(c, 401, "token tidak valid")
			return
		}
		c.Set(UserIDKey, uint(sub))
		c.Next()
	}
}

func UserID(c *gin.Context) uint {
	return c.GetUint(UserIDKey)
}
