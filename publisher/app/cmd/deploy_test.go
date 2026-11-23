package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/radio-t-site/publisher/app/cmd/mocks"
)

func TestDeploy_Do(t *testing.T) {
	ex := &mocks.ShellExecutorIfaceMock{
		RunFunc:      func(_ string, _ ...string) {},
		RunShellFunc: func(_ string) {},
	}

	d := Deploy{ShellExecutorIface: ex}
	d.Do()

	require.Equal(t, 1, len(ex.RunShellCalls()))
	assert.Equal(t, `git pull && git add . && { git diff --staged --exit-code --quiet || git commit -m auto; } && git push`,
		ex.RunShellCalls()[0].Script)

	require.Equal(t, 1, len(ex.RunCalls()))
	assert.Equal(t, "ssh", ex.RunCalls()[0].Cmd)
	assert.Equal(t, []string{"umputun@master.radio-t.com",
		"cd /srv/radio-t/site.hugo && git pull && docker compose run --rm hugo"}, ex.RunCalls()[0].Params)
}

func TestDeploy_DoCommitScript(t *testing.T) {
	tbl := []struct {
		name           string
		prep           func(t *testing.T, work string)
		wantErr        bool
		wantAutoCommit bool
		wantPushed     bool
	}{
		{name: "failed pull with a staged change makes no commit", prep: func(t *testing.T, work string) { // regression: #537
			gitRun(t, work, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
			writeFile(t, work, "f.txt", "changed")
			gitRun(t, work, "add", "f.txt")
		}, wantErr: true},
		{name: "changes are committed and pushed", prep: func(t *testing.T, work string) {
			writeFile(t, work, "f.txt", "changed")
			writeFile(t, work, "new.txt", "new")
		}, wantAutoCommit: true, wantPushed: true},
		{name: "nothing to commit pushes the local commit without an auto commit", prep: func(t *testing.T, work string) {
			writeFile(t, work, "f.txt", "local")
			gitRun(t, work, "commit", "-am", "local")
		}, wantPushed: true},
	}

	for _, tt := range tbl {
		t.Run(tt.name, func(t *testing.T) {
			origin, work := gitSandbox(t)
			tt.prep(t, work)
			localBefore, originBefore := gitRun(t, work, "rev-parse", "HEAD"), gitRun(t, origin, "rev-parse", "master")

			var runErr error
			ex := &mocks.ShellExecutorIfaceMock{
				RunFunc: func(_ string, _ ...string) {},
				RunShellFunc: func(script string) {
					c := exec.Command("sh", "-c", script)
					c.Dir = work
					out, err := c.CombinedOutput()
					t.Logf("%s", out)
					runErr = err
				},
			}
			d := Deploy{ShellExecutorIface: ex}
			d.Do()

			if tt.wantErr {
				assert.Error(t, runErr)
			} else {
				assert.NoError(t, runErr)
			}
			localAfter, originAfter := gitRun(t, work, "rev-parse", "HEAD"), gitRun(t, origin, "rev-parse", "master")
			if tt.wantAutoCommit {
				assert.Equal(t, localBefore, gitRun(t, work, "rev-parse", "HEAD~1"))
				assert.Equal(t, "auto", gitRun(t, work, "log", "-1", "--format=%s"))
			} else {
				assert.Equal(t, localBefore, localAfter, "no commit expected")
			}
			if tt.wantPushed {
				assert.Equal(t, localAfter, originAfter)
			} else {
				assert.Equal(t, originBefore, originAfter, "origin must not change")
			}
		})
	}
}

func gitSandbox(t *testing.T) (origin, work string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	gitRun(t, home, "config", "--global", "user.name", "test")
	gitRun(t, home, "config", "--global", "user.email", "test@example.com")
	gitRun(t, home, "config", "--global", "init.defaultBranch", "master")

	root := t.TempDir()
	origin, work = filepath.Join(root, "origin.git"), filepath.Join(root, "work")
	gitRun(t, root, "init", "--bare", origin)
	gitRun(t, root, "clone", origin, work)
	writeFile(t, work, "f.txt", "initial")
	gitRun(t, work, "add", "f.txt")
	gitRun(t, work, "commit", "-m", "initial")
	gitRun(t, work, "push", "origin", "master")
	return origin, work
}

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	require.NoError(t, err, "git %s: %s", strings.Join(args, " "), out)
	return strings.TrimSpace(string(out))
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}
