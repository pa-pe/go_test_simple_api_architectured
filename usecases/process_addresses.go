package usecases

import (
	"time"

	"testapi/models"
	"testapi/repositories"
	"testapi/utils"
)

type ProcessAddressesUseCase struct {
	CacheRepo repositories.CacheRepository
}

func NewProcessAddressesUseCase(cacheRepo repositories.CacheRepository) *ProcessAddressesUseCase {
	return &ProcessAddressesUseCase{
		CacheRepo: cacheRepo,
	}
}

func (uc *ProcessAddressesUseCase) Execute(request models.RequestData) (models.ResponseData, error) {
	// startTime is needed to calculate the processing time
	startTime := time.Now()

	// Remove duplicates from incoming request
	uniqueAddresses, duplicatesRemoved := utils.RemoveDuplicateAddresses(request.Addresses)

	// Asynchronously update the cache
	//	go uc.CacheRepo.UpdateCache(request.Name, request.Last, uniqueAddresses)
	go func() {
		err := uc.CacheRepo.UpdateCache(request.Name, request.Last, uniqueAddresses)
		if err != nil {
			// Do nothing because all error logging in realized in UpdateCache function
		}
	}()

	// Calculate processing time
	processingTime := time.Since(startTime).String()

	response := models.ResponseData{
		Name:      request.Name,
		Last:      request.Last,
		Addresses: uniqueAddresses,
		ProcessingInfo: models.ProcessingInfo{
			TimeTaken:         processingTime,
			DuplicatesRemoved: duplicatesRemoved,
		},
	}

	return response, nil
}
