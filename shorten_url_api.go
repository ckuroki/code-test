package main

import (
	"github.com/ckuroki/code-test/store"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

var log = logrus.New()

// ShortURL url to be shorten
type ShortURL struct {
	URL  string `json:"url"`
	Code string `json:"code"`
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
	var short ShortURL
	// Connection to the kv store
	dbFile := viper.Get("db_file").(string)
	db, err := store.Open(dbFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	c.ShouldBindJSON(&short) // Bind JSON body

	// url not present
	if len(short.URL) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	// check if key is valid
	err = isValidCode(short.Code)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}

	// check if key is in use
	_, err = db.Get(short.Code)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{})
		return
	}

	// Valid key store in db
	err = db.Put(short.Code, short.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	// Return JSON result
	c.JSON(http.StatusOK, gin.H{"code": short.Code})
}

//GetCode redirects a short url to the stored url
func GetCode(c *gin.Context) {
	// Connection to the kv store
	dbFile := viper.Get("db_file").(string)
	db, err := store.Open(dbFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Get code from url parameter
	code := c.Params.ByName("code")

	// check if key exists
	url, err := db.Get(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	// Return redirect
	c.Redirect(http.StatusFound, url)
}

//GetStats get statistics about a short url
func GetStats(c *gin.Context) {
	// Connection to the kv store
	dbFile := viper.Get("db_file").(string)
	db, err := store.Open(dbFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
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

// isValidCode validates a shortcode
func isValidCode(code string) (err error) {
	return
}
