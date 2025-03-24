package auth

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Techeer-Hogwarts/slack-bot/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/slack-go/slack"
)

var db *sql.DB

type JWTClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

var (
	jwtSecret      = config.GetEnvVarAsString("JWT_SECRET", "ak")
	signingKey     = config.GetEnvVarAsString("SLACK_SIGNING_SECRET", "ak")
	wildcardAPIKey = config.GetEnvVarAsString("CICD_API_KEY", "ak")
)

func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key"})
			c.Abort()
			return
		}

		// wildcard key
		if apiKey == wildcardAPIKey {
			c.Set("valid_api_key", true)
			c.Next()
			return
		}

		var isActive bool
		err := db.QueryRow("SELECT active FROM users WHERE api_key = $1", strings.TrimSpace(apiKey)).Scan(&isActive)

		if err != nil || !isActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or inactive API key"})
			c.Abort()
			return
		}

		c.Set("valid_api_key", true)
		c.Next()
	}
}

func VerifySlackRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := slack.NewSecretsVerifier(c.Request.Header, signingKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Slack signature"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if _, err := s.Write(body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify request body"})
			c.Abort()
			return
		}

		if err := s.Ensure(); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Slack signature verification failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		access_token, err := c.Cookie("access_token")
		if err != nil {
			c.Set("valid_jwt", false)
			c.Next()
			return
		}

		claims, err := validateToken(access_token)
		if err != nil {
			c.Set("valid_jwt", false)
			c.Next()
			return
		}

		c.Set("user_id", claims.ExpiresAt)
		c.Set("valid_jwt", true)
		c.Next()
	}
}

func validateToken(access_token string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(access_token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func GenerateJWT(userID int, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(6 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}
