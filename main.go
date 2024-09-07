package main

import (
	"log"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"

	"testapi/controllers"
	"testapi/repositories"
	"testapi/usecases"
)

type Config struct {
	CacheDir string `yaml:"cache_dir"`
	HTMLDir  string `yaml:"html_dir"`
}

func loadConfig() (*Config, error) {
	// Open config file
	configFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	// Read config data
	fileInfo, err := os.Stat("config.yaml")
	if err != nil {
		return nil, err
	}

	// Read file buffer
	data := make([]byte, fileInfo.Size())
	_, err = configFile.Read(data)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	r := gin.Default()

	// disable proxies
	r.SetTrustedProxies(nil)

	// Load HTML templates from the configured directory
	r.LoadHTMLGlob(config.HTMLDir + "/*")

	// Initialize the repository and use case
	cacheRepo := repositories.NewFileCacheRepository(config.CacheDir)
	processUseCase := usecases.NewProcessAddressesUseCase(cacheRepo)
	processController := controllers.NewProcessJsonController(processUseCase)

	// Route for the main page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main.html", nil)
	})

	// Handle JSON processing
	r.POST("/process", processController.Process)

	// Handle 404 error
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// Handle internal server errors (example)
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{
				"error": c.Errors.String(),
			})
		}
	})

	// Start the server
	r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
