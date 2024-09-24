package controllers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testapi/controllers"
	"testapi/models"
	"testapi/usecases/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProcessJsonControllerTestSuite struct {
	suite.Suite
	mockUseCase      *mocks.IProcessAddressesUseCase
	controller       *controllers.ProcessJsonController
	w                *httptest.ResponseRecorder
	ctx              *gin.Context
	requestValidJson string
	requestValidData *models.RequestData
}

func (suite *ProcessJsonControllerTestSuite) SetupSuite() {
	// Create a mock use case
	suite.mockUseCase = new(mocks.IProcessAddressesUseCase)

	// Create the controller with the mock use case
	suite.controller = controllers.NewProcessJsonController(suite.mockUseCase)

	suite.requestValidJson = `{
		"name": "John",
		"last": "Doe",
		"addresses": [{"country": "USA", "city": "New York"}]
	}`

	if err := json.Unmarshal([]byte(suite.requestValidJson), &suite.requestValidData); err != nil {
		suite.T().Fatalf("Error unmarshalling JSON: %v", err)
	}
}

func (suite *ProcessJsonControllerTestSuite) SetupTest() {
	suite.w = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.w)
}

func (suite *ProcessJsonControllerTestSuite) TearDownTest() {
	suite.mockUseCase.ExpectedCalls = nil // Reset the mock's expected calls
	//suite.ctx.Request = nil          // not works
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_ValidRequest() {
	// Set the mock usecase behavior on error
	responseData := models.ResponseData{
		Name:      suite.requestValidData.Name,
		Last:      suite.requestValidData.Last,
		Addresses: suite.requestValidData.Addresses,
		ProcessingInfo: models.ProcessingInfo{
			TimeTaken:         "1ms",
			DuplicatesRemoved: 0,
		},
	}

	suite.mockUseCase.On("Execute", *suite.requestValidData).Return(responseData, nil)

	suite.ctx.Request = &http.Request{
		Body:   io.NopCloser(strings.NewReader(suite.requestValidJson)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}
	suite.controller.Process(suite.ctx)

	// Check the response
	assert.Equal(suite.T(), http.StatusOK, suite.w.Code)

	// Verifying that the mock expectations were fulfilled ("Execute")
	suite.mockUseCase.AssertExpectations(suite.T())
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_InvalidJsonRequest() {
	// Making request with invalid json
	suite.ctx.Request = &http.Request{
		Body:   io.NopCloser(strings.NewReader(`invalid json`)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	suite.controller.Process(suite.ctx)

	// Checking the response code
	assert.Equal(suite.T(), http.StatusBadRequest, suite.w.Code)

	// Verify that the UseCase was not called
	suite.mockUseCase.AssertNotCalled(suite.T(), "Execute")
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_UseCaseError() {
	// Set the mock usecase behavior on error
	suite.mockUseCase.On("Execute", *suite.requestValidData).Return(models.ResponseData{}, assert.AnError)

	suite.ctx.Request = &http.Request{
		Body:   io.NopCloser(strings.NewReader(suite.requestValidJson)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}
	suite.controller.Process(suite.ctx)

	// Checking the response error code
	assert.Equal(suite.T(), http.StatusInternalServerError, suite.w.Code)

	// Verifying that the mock expectations were fulfilled ("Execute")
	suite.mockUseCase.AssertExpectations(suite.T())
}

func TestProcessJsonControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessJsonControllerTestSuite))
}
