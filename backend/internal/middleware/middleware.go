package middleware

import (
	"net/http"
	"strings"

	"github.com/franchise-system/backend/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Fallback: check ?token= query param (used by iframe/embed for PDF viewing)
			if tokenParam := c.Query("token"); tokenParam != "" {
				authHeader = "Bearer " + tokenParam
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Header Authorization diperlukan"})
				c.Abort()
				return
			}
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Format otorisasi tidak valid"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Token tidak valid atau sudah kedaluwarsa"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Token tidak valid"})
			c.Abort()
			return
		}
		userID, _ := uuid.Parse(claims["user_id"].(string))
		c.Set("user_id", userID)
		c.Set("user_email", claims["email"].(string))
		c.Set("user_role", claims["role"].(string))
		if panel, ok := claims["panel"].(string); ok {
			c.Set("panel", panel)
		}
		c.Next()
	}
}

// AdminOnly ensures the request comes from an admin-panel issued JWT
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Akses ditolak"})
			c.Abort()
			return
		}
		if !entity.IsAdminRole(role.(string)) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Akses khusus admin"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoleAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Role tidak ditemukan"})
			c.Abort()
			return
		}
		userRole := role.(string)
		for _, r := range allowedRoles {
			if r == userRole {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Anda tidak memiliki akses untuk fitur ini"})
		c.Abort()
	}
}

func CORS(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := ""
		for _, o := range origins {
			if o == origin || o == "*" {
				allowed = origin
				break
			}
		}
		if allowed == "" && len(origins) > 0 {
			allowed = origins[0]
		}
		c.Header("Access-Control-Allow-Origin", allowed)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
