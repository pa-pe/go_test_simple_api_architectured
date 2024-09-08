package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	LogDir   string `yaml:"log_dir"`
	LogFile  string `yaml:"log_file"`
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

	// setup LogDir
	if _, err := os.Stat(config.LogDir); os.IsNotExist(err) {
		err := os.MkdirAll(config.LogDir, 0755)
		if err != nil {
			log.Fatalf("Error creating log directory %s: %v", config.LogDir, err)
		}
	}

	// setup LogFile
	logPath := filepath.Join(config.LogDir, config.LogFile)
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("Error opening log file %s: %v", logPath, err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Channel for intercepting system signals
	// interception of signals is needed to log the termination of the application
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for a system signal in a goroutine
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v. Shutting down.", sig)
		os.Exit(0)
	}()

	// The main logic of the application
	runApp(config)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runApp(config *Config) error {
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

	log.Println("Application started.")
	defer log.Println("Application stopped.")

	// Start the server
	return r.Run(":8080")
}
