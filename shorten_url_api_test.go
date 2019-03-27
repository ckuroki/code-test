package main

import (
	"fmt"
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/ckuroki/code-test/store"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var db *store.DB
var r *gofight.RequestConfig

// mock data
const (
	goodCode  = "a12345"
	badCode   = "{\"url\":\"http://example.com\",\"code\":\"example\"}"
	getNotURL = "/z54321"
	postURL   = "/urls"
	badReq    = "{\"whatever\":\"1234\"}"
)

func init() {
	// Read configuration
	err := readconfig()
	if err != nil {
		log.Error("configuration file error: \n")
		return
	}

	// Connection to the kv store
	dbFile := viper.Get("db_file").(string)
	db, err := store.Open(dbFile)
	if err != nil {
		log.Fatal("DB Open error: " + err.Error())
	}
	defer db.Close()

	// Delete mock code key
	db.Delete(goodCode)

	r = gofight.New()

}

func TestPostBadFormed(t *testing.T) {

	r.POST(postURL).
		SetBody(string(badReq)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		})
}

func TestPost(t *testing.T) {
	good := fmt.Sprintf("{\"url\":\"http://example.com\",\"code\":\"%s\"}", goodCode)
	r.POST(postURL).
		SetBody(string(good)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			// Checking valid result content
			data := []byte(res.Body.String())
			value, _ := jsonparser.GetString(data, "code")
			assert.Equal(t, "a12345", value)
			// Now check http code
			assert.Equal(t, http.StatusOK, res.Code)
		})
}

func TestPostConflict(t *testing.T) {
	good := fmt.Sprintf("{\"url\":\"http://example.com\",\"code\":\"%s\"}", goodCode)
	r.POST(postURL).
		SetBody(string(good)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusConflict, res.Code)
		})
}

func TestPostBadCode(t *testing.T) {
	r.POST(postURL).
		SetBody(string(badCode)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		})
}

func TestGetCode(t *testing.T) {
	getURL := fmt.Sprintf("/%s", goodCode)
	r.GET(getURL).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusFound, res.Code)
		})
}

func TestGetNotFound(t *testing.T) {
	r.GET(getNotURL).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, res.Code)
		})
}

func TestGetStats(t *testing.T) {
	getStatsURL := fmt.Sprintf("/%s/stats", goodCode)
	r.GET(getStatsURL).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			// Checking valid result content
			data := []byte(res.Body.String())
			// Check usage count
			count, err := jsonparser.GetInt(data, "usage_count")
			log.Info("Checking that usage_count exists ")
			assert.Nil(t, err)
			log.Info("Checking usage_count value")
			assert.True(t, count > 0)
			// Check created_at
			log.Info("Checking that created_at exists ")
			_, err = jsonparser.GetString(data, "created_at")
			assert.Nil(t, err)
			log.Info("Checking that created_at is in ISO8601 format ")
			_, err = time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
			assert.Nil(t, err)
			// Check last_usage
			log.Info("Checking that last_usage exists ")
			_, err = jsonparser.GetString(data, "created_at")
			assert.Nil(t, err)
			log.Info("Checking that last_usage is in ISO8601 format ")
			_, err = time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
			assert.Nil(t, err)
		})
}
