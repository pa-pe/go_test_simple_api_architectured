package usecases_test

import (
	"sync"
	"testing"
	"time"

	"testapi/models"
	"testapi/repositories/mocks"
	"testapi/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessAddressesUseCase_Execute(t *testing.T) {
	mockCacheRepo := new(mocks.CacheRepository)

	// test data
	request := models.RequestData{
		Name: "John",
		Last: "Doe",
		Addresses: []models.Address{
			{Country: "USA", City: "Baltimore"},
			{Country: "USA", City: "Baltimore"}, // Duplicate
		},
	}

	// Create a UseCase with a mock repository
	uc := usecases.NewProcessAddressesUseCase(mockCacheRepo)

	var wg sync.WaitGroup
	wg.Add(1)

	mockCacheRepo.On("UpdateCache", "John", "Doe", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		wg.Done()
	})

	response, err := uc.Execute(request)

	// Waiting for the goroutine to complete
	wg.Wait()

	assert.NoError(t, err)

	// Check that duplicates have been removed
	assert.Equal(t, 1, len(response.Addresses))

	// We check that the processed information is correct
	assert.Equal(t, "John", response.Name)
	assert.Equal(t, "Doe", response.Last)

	// Check that 1 duplicate was removed
	assert.Equal(t, 1, response.ProcessingInfo.DuplicatesRemoved)

	// Check that the processing time is not empty
	assert.NotEmpty(t, response.ProcessingInfo.TimeTaken)

	// Check that the UpdateCache method was called
	mockCacheRepo.AssertCalled(t, "UpdateCache", request.Name, request.Last, mock.Anything)

	// We check that the processing time does not exceed the specified interval
	processingDuration, _ := time.ParseDuration(response.ProcessingInfo.TimeTaken)
	assert.Less(t, processingDuration.Milliseconds(), int64(100))
}
