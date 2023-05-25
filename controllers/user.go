package controllers

import (
	"Assignment/models"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var users []models.User

func RegisterUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user.Password = string(hashPassword)
	users = append(users, user)

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "registration successful",
	})
}

func Login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	for _, v := range users {
		if v.Username == user.Username {
			err := bcrypt.CompareHashAndPassword([]byte(v.Password), []byte(user.Password))
			if err == nil {
				token := GenerateToken(v.Username)
				c.IndentedJSON(http.StatusOK, gin.H{
					"token": token,
				})
				return
			}
		}
	}
	c.IndentedJSON(http.StatusUnauthorized, gin.H{
		"error": "invalid credentials",
	})
}

func GenerateToken(username string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secretKey := []byte("your_secret_key") // Replace with your own secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return ""
	}

	return tokenString
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// Token verification
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey := []byte("your_secret_key")
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid username",
			})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
