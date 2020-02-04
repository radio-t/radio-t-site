package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmd_LastShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/last/1?categories=podcast", r.URL.String())
		w.Write([]byte(`[{"show_num": 683}]`))
	}))
	defer ts.Close()

	res, err := LastShow(http.Client{Timeout: 10 * time.Millisecond}, ts.URL)
	require.NoError(t, err)
	assert.Equal(t, 683, res)
}

func TestCmd_LastShowFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/last/1?categories=podcast", r.URL.String())
		w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	_, err := LastShow(http.Client{Timeout: 10 * time.Millisecond}, "http://127.0.0.2:9999/xyz")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can't get last shows")
}

func TestShellExecutor_Do(t *testing.T) {
	c := ShellExecutor{}
	err := c.Do("ls -la")
	assert.NoError(t, err)

	err = c.Do("ls -la && pwd")
	assert.NoError(t, err)

	err = c.Do("lxxxxxxs -la")
	assert.Error(t, err)
}
