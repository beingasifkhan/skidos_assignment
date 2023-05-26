package controllers

import (
	"assignment_skidos/models"
	"assignment_skidos/utils"
	"net/http"

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
				token := utils.GenerateToken(v.Username)
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
