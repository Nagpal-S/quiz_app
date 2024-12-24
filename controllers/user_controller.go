package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"quizapp/models"
	"strconv"
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
	if err := uc.DB.Where("id = ? AND register = 1", userId).First(&user).Error; err != nil {
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

// EditUserProfile This API edits the user profile
//
//	@Summary		This API edits the user profile
//	@Description	This API edits the user profile
//	@Schemes
//	@Tags		User Auth
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		id		formData	string	true	"User ID"
//	@Param		name	formData	string	true	"User Name"
//	@Param		email	formData	string	true	"User Email"
//	@Param		phone	formData	string	true	"User Phone"
//	@Param		image	formData	string	true	"User Image"
//	@Param		gender	formData	string	false	"User Gender (Male, Female, Others)"
//	@Success	200		{object}	editProfileResponse
//	@Router		/users/edit-user-profile [post]
func (uc *UserController) EditUserProfile(c *gin.Context) {
	var user models.User

	// Parse form-data values
	id := c.PostForm("id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	image := c.PostForm("image")
	gender := c.PostForm("gender")

	// Validate mandatory fields
	if id == "" || name == "" || email == "" || phone == "" || gender == "" || image == "" {
		c.JSON(400, gin.H{"status": "0", "message": "Bad Request. Missing required parameters."})
		return
	}

	// Fetch the user from the database
	if err := uc.DB.Where("id = ? AND register = 1", id).First(&user).Error; err != nil {
		c.JSON(404, gin.H{
			"status":  "0",
			"message": "Invalid userId",
		})
		return
	}

	// Update user info
	user.Name = name
	user.Email = email
	user.Phone = phone
	user.Image = image
	user.Gender = gender

	// Save the updated user info
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "Error while updating user info.",
		})
		return
	}

	// Respond with success
	c.JSON(200, gin.H{
		"status":  "1",
		"message": "User info updated successfully.",
	})
}

type editProfileResponse struct {
	Status  string `json:"status" example:"1"`
	Message string `json:"message" example:"User info updated successfully."`
}

// var requestBody struct {
// 	UserId          string  `json:"user_id" example:"1"`
// 	Amount          float64 `json:"amount" example:"200"`
// 	TransactionType string  `json:"transaction_type" example:"DEBIT"`
// }

// InitiateUserTransaction This API will make user transactions
//
//	@Summary		This API will make user transactions
//	@Description	This API will make user transactions
//	@Schemes
//	@Tags		User Wallet
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		user_id				formData	string	true	"User ID"
//	@Param		amount				formData	float64	true	"Transaction Amount"
//	@Param		transaction_type	formData	string	true	"Transaction Type (CREDIT/DEBIT)"
//	@Success	200					{object}	transactionResponse
//	@Router		/users/initiate-user-transaction [post]
func (uc *UserController) InitiateUserTransaction(c *gin.Context) {

	// Extract form-data parameters
	userId := c.PostForm("user_id")
	amountStr := c.PostForm("amount")
	transactionType := c.PostForm("transaction_type")

	var wallet models.UserWallet
	var transaction models.UserTransactions
	var user models.User

	// Validate required fields
	if userId == "" || amountStr == "" || transactionType == "" {
		c.JSON(400, gin.H{
			"status":  "0",
			"message": "Bad request. All fields (user_id, amount, transaction_type) are required.",
		})
		return
	}

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "0",
			"message": "Invalid amount. Must be a valid number.",
		})
		return
	}

	// Check valid user
	if err := uc.DB.Where("id = ? AND register = 1", userId).Take(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "Invalid user_id, user not found. " + err.Error(),
		})
		return
	}

	// Get user wallet
	if err := uc.DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"status":  "0",
				"message": "User wallet not found.",
			})
		} else {
			c.JSON(500, gin.H{
				"status":  "0",
				"message": "DB error while fetching user wallet. " + err.Error(),
			})
		}
		return
	}

	// Manage transaction for DEBIT and CREDIT
	if transactionType == "DEBIT" {
		// If transaction is debit
		if wallet.Amount < amount {
			c.JSON(500, gin.H{
				"status":  "0",
				"message": "Amount cannot be debited. INSUFFICIENT BALANCE.",
			})
			return
		}

		wallet.Amount -= amount
		transaction.Title = "Withdrawal"
	} else if transactionType == "CREDIT" {
		// If transaction is credit
		wallet.Amount += amount
		transaction.Title = "Deposit"
	} else {
		c.JSON(400, gin.H{
			"status":  "0",
			"message": "Invalid transaction type. Must be CREDIT or DEBIT.",
		})
		return
	}

	transaction.Amount = amount
	transaction.UserId = user.ID
	transaction.TransactionType = transactionType

	// Create transaction
	if err := uc.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "DB error while creating transaction.",
		})
		return
	}

	// Update wallet
	if err := uc.DB.Model(&wallet).Update("amount", wallet.Amount).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "DB error while updating wallet.",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "1",
		"message": "Transaction successful and user wallet updated.",
	})
}

type transactionResponse struct {
	Status  string `json:"status" example:"1"`
	Message string `json:"message" example:"Transaction successful and user wallet updated."`
}

type TResponse struct {
	Title           string    `json:"title"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	Created         time.Time `json:"created"`
}

// GetUserWalletDetails This API will get user details
// @Summary      Get user wallet details
// @Description  Fetches user wallet and transaction details
// @Tags         User Wallet
// @Accept       json
// @Produce      json
// @Param        user_id  path     int  true  "User ID"
// @Success      200      {object} APIResponse
// @Router       /users/get-user-wallet-details/{user_id} [get]
func (uc *UserController) GetUserWalletDetails(c *gin.Context) {

	userId := c.Param("user_id")

	var wallet models.UserWallet
	var transaction []TResponse
	// get user wallet
	if err := uc.DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "User wallet not found. " + err.Error(),
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while getting user wallet. " + err.Error(),
			})

		}

		return

	}

	// get user transaction
	if err := uc.DB.Raw("SELECT title, transaction_type, amount, created FROM user_transactions WHERE user_id = ?", userId).Scan(&transaction).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while getting user transactions. " + err.Error(),
		})
		return

	}

	response := gin.H{

		"wallet":       wallet.Amount,
		"transactions": transaction,
	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "User wallet details found.",
		"details": response,
	})

}

type Transaction struct {
	Title           string    `json:"title" example:"Deposit"`
	TransactionType string    `json:"transaction_type" example:"CREDIT"`
	Amount          float64   `json:"amount" example:"20"`
	Created         time.Time `json:"created" example:"2024-12-18T23:04:50+05:30"`
}

// WalletDetails defines the structure for wallet details
type WalletDetails struct {
	Transactions []Transaction `json:"transactions"`
	Wallet       float64       `json:"wallet"`
}

// APIResponse defines the overall structure for the response
type APIResponse struct {
	Status  string        `json:"status" example:"1"`
	Message string        `json:"message" example:"User wallet details found."`
	Details WalletDetails `json:"details"`
}
