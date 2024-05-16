package functional

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Links(t *testing.T) {
	tc := Context(t)
	assert := assert.New(t)

	response := tc.Get("/")
	assert.Equal(http.StatusOK, response.StatusCode)
	assert.Equal(1, response.Find("table").Length())
	assert.Equal(0, response.Find("table tbody tr").Length())

	form := url.Values{}
	form.Add("path", "home")
	form.Add("url", "https://home.example.com")

	response = tc.PostForm("/links", form)
	assert.Equal(http.StatusOK, response.StatusCode)
	assert.Equal(1, response.Find("table").Length())
	assert.Equal(1, response.Find("table tbody tr").Length())
	assert.Equal("home", response.Find("table tbody tr:nth-child(1) th:nth-child(1)").Text())
}
