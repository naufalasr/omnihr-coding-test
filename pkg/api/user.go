package api

import (
	"errors"
	"fmt"
	"net/http"
	"omnihr-coding-test/pkg/auth"
	"omnihr-coding-test/pkg/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler
func LoginHandler(c *gin.Context) {
	appCtx, exists := c.MustGet("appCtx").(*AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	var incomingUser models.User
	var dbUser models.User

	// Get JSON body
	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Fetch the user from the database
	if err := appCtx.DB.Where("username = ? AND company_id = ?", incomingUser.Username, incomingUser.CompanyID).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(auth.GenerateTokenParam(dbUser.Username, dbUser.CompanyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RegisterHandler
func RegisterHandler(c *gin.Context) {
	appCtx, exists := c.MustGet("appCtx").(*AppContext)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	var user models.LoginUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create new user
	newUser := models.User{Username: user.Username, Password: hashedPassword, CompanyID: user.CompanyID}

	// Save the user to the database
	if err := appCtx.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
