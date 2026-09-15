package cmd

import (
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
	assert.Equal(t, `git pull && git add . && git diff --staged --exit-code --quiet || git commit -m auto && git push`,
		ex.RunShellCalls()[0].Script)

	require.Equal(t, 1, len(ex.RunCalls()))
	assert.Equal(t, "ssh", ex.RunCalls()[0].Cmd)
	assert.Equal(t, []string{"umputun@master.radio-t.com",
		"cd /srv/radio-t/site.hugo && git pull && docker compose run --rm hugo"}, ex.RunCalls()[0].Params)
}
