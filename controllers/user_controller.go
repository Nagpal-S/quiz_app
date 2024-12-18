package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"quizapp/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// private function to generate OTP
func GenerateOTP() string {
	// Create a new random number generator with a unique seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random 4-digit number
	otp := rng.Intn(9000) + 1000 // Ensures the range is [1000, 9999]

	// Return the OTP as a string
	return fmt.Sprintf("%04d", otp)
}

// GetUserByPhone generateing OTP for the user
//
//	@Summary		generateing OTP for the user
//	@Description	generateing OTP for the user
//	@Schemes
//	@Tags		User AUth
//	@Accept		json
//	@Produce	json
//	@Param		phone	path		string	true	"user phone number"
//	@Success	200		{object}	OTPResponse
//	@Router		/users/check-user/{phone} [get]
func (uc *UserController) GetUserByPhone(c *gin.Context) {
	phone := c.Param("phone")
	var user models.User

	// Check if the phone number exists in the database
	err := uc.DB.Where("phone = ?", phone).First(&user).Error
	if err == nil {
		// If user exists, update the OTP
		user.Otp = GenerateOTP() // Generate a new OTP
		if err := uc.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update OTP", "status": "0"})
			return
		}

		details := map[string]interface{}{
			"otp":    user.Otp,
			"userId": user.ID,
		}

		// Send the updated OTP in the response
		c.JSON(http.StatusOK, gin.H{
			"message": "OTP updated",
			"status":  "1",
			"details": details,
		})
		return
	}

	// If user does not exist, create a new user
	if err.Error() == "record not found" {
		newUser := models.User{
			Phone: phone,
			Otp:   GenerateOTP(), // Generate OTP for the new user
		}

		// Save the new user to the database
		if err := uc.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		details := map[string]interface{}{
			"userId": newUser.ID,  // Add userId in the response
			"otp":    newUser.Otp, // Include the OTP as well
		}

		// Send the OTP and user ID in the response
		c.JSON(http.StatusCreated, gin.H{
			"message": "New user created",
			"status":  "1",
			"details": details,
		})
		return
	}

	// Handle any other unexpected errors
	c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
}

type OTPResponse struct {
	Details struct {
		OTP    string `json:"otp" example:"8162"`
		UserID int    `json:"userId" example:"3"`
	} `json:"details"`
	Message string `json:"message" example:"New user created"`
	Status  string `json:"status" example:"1"`
}

// VerifyOTP This API will verify user OTP with userId
//
//	@Summary		This API will verify user OTP with userId
//	@Description	This API will verify user OTP with userId
//	@Schemes
//	@Tags		User AUth
//	@Accept		json
//	@Produce	json
//	@Param		userId	path		string	true	"user app Id"
//	@Param		otp		path		string	true	"user otp"
//	@Success	200		{object}	verifyOtpResponse
//	@Router		/users/verify-otp/{userId}/{otp} [get]
func (uc *UserController) VerifyOTP(c *gin.Context) {
	// Get userId and otp from the request params
	userId := c.Param("userId")
	otp := c.Param("otp")

	var user models.User
	var wallet models.UserWallet

	// Find user by userId
	err := uc.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		// If user is not found
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "status": "0"})
		return
	}

	// Check if the OTP matches
	if user.Otp != otp {
		// If OTP does not match
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP", "status": "0"})
		return
	}

	wallet.UserId = uint64(user.ID)
	wallet.Created = user.Created
	if err := uc.DB.Create(&wallet).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "Db error while creating user wallet",
		})
		return

	}

	// Check if the user is registered or not (Register == 0 or 1)
	if user.Register == "0" {
		// User is not registered, update the status
		user.Register = "1" // Update register to "1" (registered)

		if err := uc.DB.Model(&user).Select("Register").Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update registration status",
				"status":  "0",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"status":  "1",
			"details": user, // Send user information in the response
		})
		return
	} else if user.Register == "1" {
		// User is already registered
		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in successfully",
			"status":  "1",
			"details": user, // Send user information in the response
		})
		return
	}

	// Handle unexpected cases
	c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error occurred", "status": "0"})
}

type verifyOtpResponse struct {
	Details struct {
		OTP      string `json:"otp" example:"8162"`
		UserID   int    `json:"ID" example:"3"`
		Name     string `json:"name" example:"Shivam"`
		Email    string `json:"email" example:"shivam@gmail.com"`
		Phone    string `json:"phone" example:"9144"`
		Register string `json:"register" example:"1"`
		Created  string `json:"created" example:"2024-12-10T07:04:37Z"`
	} `json:"details"`
	Message string `json:"message" example:"User logged in successfully"`
	Status  string `json:"status" example:"1"`
}

// CheckUser This API will provide user info bu id
//
//	@Summary		This API will provide user info bu id
//	@Description	This API will provide user info bu id
//	@Schemes
//	@Tags		User AUth
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"this is user id"
//	@Success	200	{object}	userInfo
//	@Router		/users/get-user/{id} [get]
func (uc *UserController) CheckUser(c *gin.Context) {
	// Get the userId from the URL parameter
	userId := c.Param("id")
	var user models.User

	// Find the user by userId
	if err := uc.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		// If no user found, return "user not found" message
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "status": "0"})
		return
	}

	// Check if user.register is 1
	if user.Register == "1" {
		// If register == 1, return user info
		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in successfully",
			"status":  "1",
			"details": user,
		})
		return
	}

	// If register is not 1, return a different message
	c.JSON(http.StatusForbidden, gin.H{
		"message": "User not registered",
		"status":  "0",
	})

}

type userInfo struct {
	Details struct {
		OTP      string `json:"otp" example:"8162"`
		UserID   int    `json:"ID" example:"3"`
		Name     string `json:"name" example:"Shivam"`
		Email    string `json:"email" example:"shivam@gmail.com"`
		Phone    string `json:"phone" example:"9144"`
		Register string `json:"register" example:"1"`
		Created  string `json:"created" example:"2024-12-10T07:04:37Z"`
	} `json:"details"`
	Message string `json:"message" example:"User logged in successfully"`
	Status  string `json:"status" example:"1"`
}

type profileInfo struct {
	ID     string `json:"userId" example:"1"`
	Name   string `json:"name" example:"Shivam Nagpal"`
	Image  string `json:"image" example:"url-of-the-image"`
	Gender string `json:"gender" example:"Male"`
	Phone  string `json:"phone" example:"0987656"`
	Email  string `json:"email" example:"sn@gmail.com"`
}

// EditUserProfile This API edit user profile
//
//	@Summary		This API edit user profile
//	@Description	This API edit user profile
//	@Schemes
//	@Tags		User AUth
//	@Accept		json
//	@Produce	json
//	@Param		id	body		profileInfo	true	"this is user info json"
//	@Success	200	{object}	editProfileResponse
//	@Router		/users/edit-user-profile [post]
func (uc *UserController) EditUserProfile(c *gin.Context) {

	var requestBody profileInfo

	var user models.User

	if err := c.BindJSON(&requestBody); err != nil {

		c.JSON(400, gin.H{"status": "0", "message": "Bad Request. Invalid parameters."})
		return

	}

	if err := uc.DB.Where("id = ?", requestBody.ID).First(&user).Error; err != nil {

		c.JSON(404, gin.H{
			"status":  "0",
			"message": "Invalid userId",
		})
		return

	}

	user.Name = requestBody.Name
	user.Email = requestBody.Email
	user.Phone = requestBody.Phone
	user.Image = requestBody.Image
	user.Gender = requestBody.Gender

	if err := uc.DB.Save(&user).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "Error while updating user info.",
		})
		return

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "User info updated successfully.",
	})

}

type editProfileResponse struct {
	Status  string `json:"status" example:"1"`
	Message string `json:"message" example:"User info updated successfully."`
}
