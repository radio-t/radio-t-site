package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/radio-t-site/publisher/app/cmd/mocks"
)

func TestDeploy_Do(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/active/last/12", r.URL.Path)
		assert.Equal(t, "Basic YWRtaW46cGFzc3dk", r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	ex := &mocks.ExecutorMock{
		RunFunc: func(cmd string, params ...string) {},
	}

	d := Deploy{
		NewsPasswd: "passwd",
		NewsAPI:    ts.URL,
		NewsHrs:    12,
		Client:     http.Client{Timeout: 10 * time.Millisecond},
		Executor:   ex,
	}

	require.NoError(t, d.Do(123))
	require.Equal(t, 3, len(ex.RunCalls()))
	assert.Equal(t, "git pull && git commit -am \"episode 123\" && git push", ex.RunCalls()[0].Cmd)
	assert.Equal(t, 0, len(ex.RunCalls()[0].Params))
	assert.Equal(t, `ssh umputun@master.radio-t.com`, ex.RunCalls()[1].Cmd)
	assert.Equal(t, 1, len(ex.RunCalls()[1].Params))
	assert.Equal(t, `cd /srv/site.hugo && git pull && docker-compose run --rm hugo`, ex.RunCalls()[1].Params[0])

	assert.Equal(t, `ssh umputun@master.radio-t.com`, ex.RunCalls()[2].Cmd)
	assert.Equal(t, []string{"docker exec -i super-bot /srv/telegram-rt-bot " +
		"--super=umputun --super=bobuk --super=ksenks --super=grayodesa --super=aleks_sys " +
		"--dbg --export-num=123 --export-path=/srv/html"}, ex.RunCalls()[2].Params)
}

func TestDeploy_archiveNews(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/active/last/12", r.URL.Path)
		assert.Equal(t, "Basic YWRtaW46cGFzc3dk", r.Header.Get("Authorization"))
	}))
	defer ts.Close()

	d := Deploy{
		NewsPasswd: "passwd",
		NewsAPI:    ts.URL,
		NewsHrs:    12,
		Client:     http.Client{Timeout: 10 * time.Millisecond},
	}

	assert.NoError(t, d.archiveNews())
}
