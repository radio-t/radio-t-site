package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
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

func TestShellExecutor_exec(t *testing.T) {
	c := ShellExecutor{}

	t.Run("argv command", func(t *testing.T) {
		assert.NoError(t, c.exec(exec.Command("ls", "-la"), "ls"))
	})

	t.Run("shell script", func(t *testing.T) {
		assert.NoError(t, c.exec(exec.Command("sh", "-c", "ls -la && pwd"), "ls && pwd"))
	})

	t.Run("missing command", func(t *testing.T) {
		assert.Error(t, c.exec(exec.Command("lxxxxxxs", "-la"), "lxxxxxxs"))
	})
}

func TestShellExecutor_Run(t *testing.T) {
	c := ShellExecutor{}

	t.Run("argument with spaces and metacharacters reaches the command intact", func(t *testing.T) {
		dir := t.TempDir()
		name := filepath.Join(dir, "a file $(id -un) `date`.txt")
		c.Run("touch", name)
		assert.FileExists(t, name, "the name must not be split or expanded on the way to touch")

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, entries, 1, "exactly one file, no extra ones from a split argument")
	})

	t.Run("dry run executes nothing", func(t *testing.T) {
		dir := t.TempDir()
		dry := ShellExecutor{Dry: true}
		dry.Run("touch", filepath.Join(dir, "created.txt"))

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestShellExecutor_RunShell(t *testing.T) {
	t.Run("shell syntax still applies", func(t *testing.T) {
		dir := t.TempDir()
		c := ShellExecutor{}
		c.RunShell("cd " + dir + " && touch one.txt && touch two.txt")
		assert.FileExists(t, filepath.Join(dir, "one.txt"))
		assert.FileExists(t, filepath.Join(dir, "two.txt"))
	})

	t.Run("dry run executes nothing", func(t *testing.T) {
		dir := t.TempDir()
		dry := ShellExecutor{Dry: true}
		dry.RunShell("touch " + filepath.Join(dir, "created.txt"))

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}
