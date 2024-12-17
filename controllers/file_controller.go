package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

// FileController handles file upload operations.
type FileController struct{}

// AWS S3 Configuration
var (
	bucketName = "quizbuck"
	region     = "ap-south-1"
	accessKey  = "AKIAZ24IS5EDEPMJWAAE"
	secretKey  = "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"
)

// UploadToS3 upload file to s3
//
//	@Summary		upload file to s3
//	@Description	upload file to s3
//	@Schemes
//	@Tags		User AUth
//	@Accept		json
//	@Produce	json
//	@Param		file	formData	file	true	"user file to upload"
//	@Success	200		{object}	fileUploadResponse
//	@Router		/upload/file [post]
func (fc *FileController) UploadToS3(c *gin.Context) {
	// Parse the file from the request
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get file from request", "status": "0"})
		return
	}
	defer file.Close()

	// Read the file into memory
	fileBuffer := bytes.NewBuffer(nil)
	if _, err := fileBuffer.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read file", "status": "0"})
		return
	}

	// Initialize AWS S3 client with explicit credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load AWS configuration", "status": "0"})
		return
	}

	client := s3.NewFromConfig(cfg)

	// Define the object key (file name in S3)
	objectKey := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Upload the file to S3
	contentType := fileHeader.Header.Get("Content-Type")
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(fileBuffer.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file to S3", "status": "0"})
		return
	}

	// Generate the S3 file URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)

	details := map[string]interface{}{
		"url": fileURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "1",
		"message": "File uploaded successfully",
		"details": details,
	})
}

type fileUploadResponse struct {
	Details struct {
		Url string `json:"url" example:"https://quizbuck.s3.ap-south-1.amazonaws.com/uploads/1734090491_new.jpg"`
	} `json:"details"`
	Message string `json:"message" example:"File uploaded successfully"`
	Status  string `json:"status" example:"1"`
}
