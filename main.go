package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

var log = logrus.New()

const (
	mainGroup = "api/v1"
)

// APIError JSONAPI compatible error
type APIError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func main() {
	// Read configuration
	err := readconfig()
	if err != nil {
		log.Error("configuration file error")
		return
	}

	// Setup logging
	logf, err := os.OpenFile("shorten_url.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("log file error:")
		return
	}
	defer logf.Close()

	log.Formatter = new(logrus.JSONFormatter)
	log.Out = logf

	port := viper.Get("port").(string)
	restEngine().Run(port)
}

//restEngine returns a new gin engine
func restEngine() *gin.Engine {

	r := gin.Default()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	r.POST("/urls", PostURL)
	r.GET("/:code", GetCode)
	r.GET("/:code/stats", GetStats)

	return r
}

//PostURL posts a new url
func PostURL(c *gin.Context) {
	// Return JSON result
	c.JSON(http.StatusOK, gin.H{})
}

//GetCode redirects a short url to the stored url
func GetCode(c *gin.Context) {
	// Return JSON result
	c.JSON(http.StatusOK, gin.H{})
}

//GetStats get statistics about a short url
func GetStats(c *gin.Context) {
	// Return JSON result
	c.JSON(http.StatusOK, gin.H{})
}

//readconfig reads and parses a configuration file
func readconfig() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}

	return nil
}
