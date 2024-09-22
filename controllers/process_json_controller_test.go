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

const validJson = `{
    "name": "John",
    "last": "Doe",
    "addresses": [{"country": "USA", "city": "New York"}]
}`

//var requestData models.RequestData
//
//func init() {
//	if err := json.Unmarshal([]byte(validJson), &requestData); err != nil {
//		log.Fatalf("Error unmarshalling JSON: %v", err)
//	}
//}

type ProcessJsonControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.IProcessAddressesUseCase
	controller  *controllers.ProcessJsonController
	w           *httptest.ResponseRecorder
	ctx         *gin.Context
}

func (suite *ProcessJsonControllerTestSuite) SetupTest() {
	// Create a mock use case
	suite.mockUseCase = new(mocks.IProcessAddressesUseCase)

	// Create the controller with the mock use case
	suite.controller = controllers.NewProcessJsonController(suite.mockUseCase)

	suite.w = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.w)
}

func (suite *ProcessJsonControllerTestSuite) TearDownTest() {
	suite.mockUseCase = nil
	suite.controller = nil
	suite.w = nil
	suite.ctx = nil
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_ValidRequest() {
	// Set the mock usecase behavior on error
	var requestData models.RequestData
	if err := json.Unmarshal([]byte(validJson), &requestData); err != nil {
		suite.T().Fatalf("Error unmarshalling JSON: %v", err)
	}
	responseData := models.ResponseData{
		Name:      requestData.Name,
		Last:      requestData.Last,
		Addresses: requestData.Addresses,
		ProcessingInfo: models.ProcessingInfo{
			TimeTaken:         "1ms",
			DuplicatesRemoved: 0,
		},
	}

	suite.mockUseCase.On("Execute", requestData).Return(responseData, nil)

	suite.ctx.Request = &http.Request{
		Body:   io.NopCloser(strings.NewReader(validJson)),
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
	var requestData models.RequestData
	if err := json.Unmarshal([]byte(validJson), &requestData); err != nil {
		suite.T().Fatalf("Error unmarshalling JSON: %v", err)
	}
	suite.mockUseCase.On("Execute", requestData).Return(models.ResponseData{}, assert.AnError)

	suite.ctx.Request = &http.Request{
		Body:   io.NopCloser(strings.NewReader(validJson)),
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
