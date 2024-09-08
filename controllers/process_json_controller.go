package controllers

import (
	"log"
	"net/http"

	"testapi/models"
	"testapi/usecases"

	"github.com/gin-gonic/gin"
)

type ProcessJsonController struct {
	UseCase *usecases.ProcessAddressesUseCase
}

func NewProcessJsonController(useCase *usecases.ProcessAddressesUseCase) *ProcessJsonController {
	return &ProcessJsonController{
		UseCase: useCase,
	}
}

func (ctrl *ProcessJsonController) Process(c *gin.Context) {
	var requestData models.RequestData

	// Expect JSON data from the request body
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Execute the use case
	responseData, err := ctrl.UseCase.Execute(requestData)
	if err != nil {
		log.Printf("Error executing use case: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Processing failed"})
		return
	}

	// Return the response
	c.JSON(http.StatusOK, responseData)
}
