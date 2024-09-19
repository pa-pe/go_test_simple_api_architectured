package controllers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"testapi/controllers"
	"testapi/models"
	"testapi/usecases/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProcessJsonControllerTestSuite struct {
	suite.Suite
	mockUseCase *mocks.IProcessAddressesUseCase
	controller  *controllers.ProcessJsonController
	router      *gin.Engine
}

func (suite *ProcessJsonControllerTestSuite) SetupTest() {
	// Создаем mock usecase
	suite.mockUseCase = new(mocks.IProcessAddressesUseCase)

	// Создаем контроллер с mock-юзкейсом
	suite.controller = controllers.NewProcessJsonController(suite.mockUseCase)

	// Настраиваем роутинг Gin
	suite.router = gin.Default()
	suite.router.POST("/process", suite.controller.Process)
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_ValidRequest() {
	// Задаем ожидаемое поведение mock usecase
	requestData := models.RequestData{
		Name: "John",
		Last: "Doe",
		Addresses: []models.Address{
			{Country: "USA", City: "New York"},
		},
	}
	responseData := models.ResponseData{
		Name:      "John",
		Last:      "Doe",
		Addresses: requestData.Addresses,
		ProcessingInfo: models.ProcessingInfo{
			TimeTaken:         "1ms",
			DuplicatesRemoved: 0,
		},
	}

	suite.mockUseCase.On("Execute", requestData).Return(responseData, nil)

	// Создаем приемник HTTP-ответа
	w := httptest.NewRecorder()

	// Создаем HTTP-запрос с валидным JSON
	req, _ := http.NewRequest("POST", "/process", createJsonBody(`{
		"name": "John",
		"last": "Doe",
		"addresses": [{"country": "USA", "city": "New York"}]
	}`))

	// Выполняем запрос
	suite.router.ServeHTTP(w, req)

	// Проверяем ответ
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Контрольная проверка что мок настройки ожиданий были выполним, в данном случае это вызов метода "Execute"
	suite.mockUseCase.AssertExpectations(suite.T())
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_InvalidJsonRequest() {
	// Создаем приемник HTTP-ответа
	w := httptest.NewRecorder()

	// Создаем HTTP-запрос с невалидным JSON
	req, _ := http.NewRequest("POST", "/process", createJsonBody(`invalid json`))

	// Выполняем запрос
	suite.router.ServeHTTP(w, req)

	// Проверяем, что возвращается ошибка
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	// Проверяем что далее код не выполнялся
	suite.mockUseCase.AssertNotCalled(suite.T(), "Execute")
}

func (suite *ProcessJsonControllerTestSuite) TestProcess_UseCaseError() {
	// Задаем поведение mock usecase при ошибке
	requestData := models.RequestData{
		Name: "John",
		Last: "Doe",
		Addresses: []models.Address{
			{Country: "USA", City: "New York"},
		},
	}
	suite.mockUseCase.On("Execute", requestData).Return(models.ResponseData{}, assert.AnError)

	// Создаем HTTP-запрос с валидным JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/process", createJsonBody(`{
		"name": "John",
		"last": "Doe",
		"addresses": [{"country": "USA", "city": "New York"}]
	}`))

	// Выполняем запрос
	suite.router.ServeHTTP(w, req)

	// Проверяем, что возвращается ошибка
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Контрольная проверка что мок настройки ожиданий были выполним, в данном случае это вызов метода "Execute"
	suite.mockUseCase.AssertExpectations(suite.T())
}

func createJsonBody(jsonStr string) io.Reader {
	return strings.NewReader(jsonStr)
}

func TestProcessJsonControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessJsonControllerTestSuite))
}
