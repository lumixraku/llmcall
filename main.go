package main

import (
	"encoding/base64"
	"net/http"

	"llmcall/gemini"

	"github.com/gin-gonic/gin"
)

// ImageEditRequest represents the request body for image edit API
type ImageEditRequest struct {
	Image  string `json:"image" binding:"required"`  // base64 encoded image
	Prompt string `json:"prompt" binding:"required"` // edit prompt
}

// ImageEditResponse represents the response for image edit API
type ImageEditResponse struct {
	Success bool     `json:"success"`
	Images  []string `json:"images,omitempty"` // base64 encoded result images
	Error   string   `json:"error,omitempty"`
}

func main() {
	r := gin.Default()

	// Serve static files from assets directory
	r.Static("/assets", "./assets")

	// Serve index.html at root
	r.GET("/", func(c *gin.Context) {
		c.File("./assets/index.html")
	})

	// Image edit API endpoint
	r.POST("/api/imageedit/", handleImageEdit)

	// Start server on port 8080
	r.Run(":8080")
}

func handleImageEdit(c *gin.Context) {
	var req ImageEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ImageEditResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Call the image edit function
	results, err := gemini.ImageEdit(req.Image, req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ImageEditResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert results to base64 strings
	var images []string
	for _, result := range results {
		b64 := base64.StdEncoding.EncodeToString(result.Data)
		// Include data URL prefix for easy display
		dataURL := "data:image/" + result.Format + ";base64," + b64
		images = append(images, dataURL)
	}

	c.JSON(http.StatusOK, ImageEditResponse{
		Success: true,
		Images:  images,
	})
}
