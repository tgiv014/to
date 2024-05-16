package functional

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tgiv014/to/app"
	"github.com/tgiv014/to/domains/config"
	"github.com/tgiv014/to/domains/identity"
)

type TestContext struct {
	t      *testing.T
	server *httptest.Server
}

type Response struct {
	*goquery.Document
	StatusCode int
}

func (tc *TestContext) Path(p string) string {
	u, err := url.Parse(tc.server.URL)
	assert.NoError(tc.t, err)

	return u.JoinPath(p).String()
}

func (tc *TestContext) Get(p string) Response {
	resp, err := tc.server.Client().Get(tc.Path(p))
	defer resp.Body.Close()
	assert.NoError(tc.t, err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	assert.NoError(tc.t, err)

	return Response{
		doc, resp.StatusCode,
	}
}

func (tc *TestContext) PostForm(p string, data url.Values) Response {
	resp, err := tc.server.Client().PostForm(tc.Path(p), data)
	defer resp.Body.Close()
	assert.NoError(tc.t, err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	assert.NoError(tc.t, err)

	return Response{
		doc, resp.StatusCode,
	}
}

func Context(t *testing.T) *TestContext {
	t.Helper()

	cfg := config.Config{
		Hostname: "test",
		DBPath:   ":memory:",
	}

	app := &app.App{Cfg: cfg, Identifier: MockIdentifier{}}

	err := app.Setup()
	assert.NoError(t, err)

	server := httptest.NewServer(app.Router().Handler())
	t.Cleanup(func() {
		server.Close()
	})

	return &TestContext{
		t,
		server,
	}
}

type MockIdentifier struct {
}

func (m MockIdentifier) Identify(c *gin.Context) (identity.User, error) {
	return identity.User{
		Name:  "Dev User",
		Login: "dev@example.com",
	}, nil
}
