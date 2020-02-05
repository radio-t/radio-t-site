package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploy_Do(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/active/last/12", r.URL.Path)
		assert.Equal(t, "Basic YWRtaW46cGFzc3dk", r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	ex := &MockExecutor{}
	ex.On("Run", "git pull && git commit -am episode %d && git push", 123)
	ex.On("Run", `ssh umputun@master.radio-t.com "cd /srv/site.hugo && git pull && docker-compose run --rm hugo"`)
	ex.On("Run", `ssh umputun@master.radio-t.com "docker exec -i super-bot /srv/telegram-rt-bot --super=umputun --super=bobuk --super=ksenks --super=grayru --dbg --export-num=%d --export-path=/srv/html"`, 123)

	d := Deploy{
		NewsPasswd:   "passwd",
		NewsAPI:      ts.URL,
		NewsDuration: time.Hour * 12,
		Client:       http.Client{Timeout: 10 * time.Millisecond},
		Executor:     ex,
	}

	require.NoError(t, d.Do(123))
	ex.AssertNumberOfCalls(t, "Run", 3)
}

func TestDeploy_archiveNews(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/active/last/12", r.URL.Path)
		assert.Equal(t, "Basic YWRtaW46cGFzc3dk", r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	d := Deploy{
		NewsPasswd:   "passwd",
		NewsAPI:      ts.URL,
		NewsDuration: time.Hour * 12,
		Client:       http.Client{Timeout: 10 * time.Millisecond},
	}

	assert.NoError(t, d.archiveNews())
}
