package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const JWTSecret = "secret-key"

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Sub   string   `json:"sub"`
	Homes []string `json:"homes"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and extracts home information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("user_email", claims.Sub)
		c.Set("user_homes", claims.Homes)

		c.Next()
	}
}

// RequireHomeAccess checks if the user has access to the specified home
func RequireHomeAccess(homeID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userHomes, exists := c.Get("user_homes")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User homes not found in context"})
			c.Abort()
			return
		}

		homes, ok := userHomes.([]string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user homes format"})
			c.Abort()
			return
		}

		// Check if user has access to the requested home
		hasAccess := false
		for _, home := range homes {
			if home == homeID {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this home"})
			c.Abort()
			return
		}

		c.Next()
	}
}
