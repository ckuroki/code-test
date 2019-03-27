package main

import (
	"fmt"
	"github.com/appleboy/gofight"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Unit Tests
func TestsShortenUrl(t *testing.T) {
	var err error

	// Read configuration
	err = readconfig()
	if err != nil {
		log.Error("configuration file error: \n")
		return
	}

	r := gofight.New()

	// mock data
	goodCode := "12345"
	badReq := "{\"whatever\":\"1234\"}"
	good := fmt.Sprintf("{\"url\":\"http://example.com\",\"code\":\"%s\"}", goodCode)
	badCode := "{\"url\":\"http://example.com\",\"code\":\"example\"}"
	postURL := "/urls"
	getURL := fmt.Sprintf("/%s", goodCode)
	getNotURL := "/z54321"

	t.Log("Testing POST /urls")

	t.Info("Should return bad request")
	r.POST(postURL).
		SetDebug(true).
		SetBody(string(badReq)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		})

	t.Info("Should return ok")
	r.POST(postURL).
		SetDebug(true).
		SetBody(string(good)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			// Checking valid result content
			data := []byte(res.Body.String())
			value, _ := jsonparser.GetString(data, "code")
			assert.Equal(t, "a12345", value)
			// Now check http code
			assert.Equal(t, http.StatusOk, res.Code)
		})

	t.Info("Should return conflict. (repeating last api call)")
	r.POST(postURL).
		SetDebug(true).
		SetBody(string(good)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.Conflict, res.Code)
		})

	t.Info("Should return Unprocessable Entity")
	r.POST(postURL).
		SetDebug(true).
		SetBody(string(badCode)).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.UnprocessableEntity, res.Code)
		})

	t.Log("Testing GET /:code")
	t.Info("Should return Found")
	r.GET(getURL).
		SetDebug(true).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusFound, res.Code)
		})

	t.Info("Should return NotFound")
	r.GET(getNotURL).
		SetDebug(true).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNotFound, res.Code)
		})

	t.Log("Testing GET /:code/stats")
	r.GET(getStatsUrl).
		SetDebug(true).
		Run(restEngine(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			// Checking valid result content
			data := []byte(res.Body.String())
			// Check usage count
			count, err := jsonparser.GetInt(data, "usage_count")
			log.Info("Checking that usage_count exists ")
			assert.Nil(t, err)
			log.Info("Checking usage_count value")
			assert.Greater(t, count, 0)
			// Check created_at
			log.Info("Checking that created_at exists ")
			_, err := jsonparser.GetString(data, "created_at")
			assert.Nil(t, err)
			log.Info("Checking that created_at is in ISO8601 format ")
			_, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
			assert.Nil(t, err)
			// Check last_usage
			log.Info("Checking that last_usage exists ")
			_, err := jsonparser.GetString(data, "created_at")
			assert.Nil(t, err)
			log.Info("Checking that last_usage is in ISO8601 format ")
			_, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
			assert.Nil(t, err)
		})
}
