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
	"github.com/stretchr/testify/suite"
)

type ProcessAddressesTestSuite struct {
	suite.Suite
	mockCacheRepo *mocks.CacheRepository
	useCase       *usecases.ProcessAddressesUseCase
	wg            sync.WaitGroup
	request       models.RequestData
}

func (suite *ProcessAddressesTestSuite) SetupTest() {
	suite.mockCacheRepo = new(mocks.CacheRepository)
	suite.useCase = usecases.NewProcessAddressesUseCase(suite.mockCacheRepo)

	// Sample request data
	suite.request = models.RequestData{
		Name: "John",
		Last: "Doe",
		Addresses: []models.Address{
			{Country: "USA", City: "Baltimore"},
			{Country: "USA", City: "Baltimore"}, // Duplicate
		},
	}

}

func (suite *ProcessAddressesTestSuite) TestExecute() {
	// Adding WaitGroup for goroutine
	suite.wg.Add(1)

	// Mock UpdateCache behavior
	suite.mockCacheRepo.On("UpdateCache", "John", "Doe", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		suite.wg.Done()
	})

	// Execute the use case
	response, err := suite.useCase.Execute(suite.request)

	// Wait for the goroutine to complete
	suite.wg.Wait()

	// Assert no errors
	assert.NoError(suite.T(), err)

	// Assert duplicates removed
	assert.Equal(suite.T(), 1, len(response.Addresses))

	// Assert response data is correct
	assert.Equal(suite.T(), "John", response.Name)
	assert.Equal(suite.T(), "Doe", response.Last)

	// Assert 1 duplicate was removed
	assert.Equal(suite.T(), 1, response.ProcessingInfo.DuplicatesRemoved)

	// Assert processing time is not empty
	assert.NotEmpty(suite.T(), response.ProcessingInfo.TimeTaken)

	// Assert processing time is within an acceptable range
	processingDuration, _ := time.ParseDuration(response.ProcessingInfo.TimeTaken)
	assert.Less(suite.T(), processingDuration.Milliseconds(), int64(100))

	// Assert UpdateCache method was called
	suite.mockCacheRepo.AssertCalled(suite.T(), "UpdateCache", suite.request.Name, suite.request.Last, mock.Anything)
}

// Execute the suite
func TestProcessAddressesTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessAddressesTestSuite))
}
