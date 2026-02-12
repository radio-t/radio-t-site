package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/radio-t-site/publisher/app/cmd/mocks"
)

func TestDeploy_Do(t *testing.T) {
	ex := &mocks.ExecutorMock{
		RunFunc: func(cmd string, params ...string) {},
	}

	d := Deploy{Executor: ex}
	d.Do()

	require.Equal(t, 2, len(ex.RunCalls()))
	assert.Equal(t, `git pull && git add . && git diff --staged --exit-code --quiet || git commit -m auto && git push`, ex.RunCalls()[0].Cmd)
	assert.Equal(t, 0, len(ex.RunCalls()[0].Params))

	assert.Equal(t, `ssh umputun@master.radio-t.com`, ex.RunCalls()[1].Cmd)
	assert.Equal(t, 1, len(ex.RunCalls()[1].Params))
	assert.Equal(t, `"cd /srv/site.hugo && git pull && docker compose run --rm hugo"`, ex.RunCalls()[1].Params[0])
}
