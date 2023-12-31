package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/publisher/cmd/mocks"
)

func TestUpload_Do(t *testing.T) {
	ex := &mocks.ExecutorMock{
		RunFunc: func(cmd string, params ...interface{}) {},
	}

	d := Upload{
		Executor: ex,
		Location: "/tmp",
	}

	d.Do(123)
	require.Equal(t, 7, len(ex.RunCalls()))
	assert.Equal(t, "mp3tags set-tags %d", ex.RunCalls()[0].Cmd)
	assert.Equal(t, []any{123}, ex.RunCalls()[0].Params)

	assert.Equal(t, "scp %s umputun@master.radio-t.com:/srv/master-node/var/media/%s", ex.RunCalls()[1].Cmd)
	assert.Equal(t, 2, len(ex.RunCalls()[1].Params))
	assert.Equal(t, []any{"/tmp/rt_podcast123/rt_podcast123.mp3", "rt_podcast123.mp3"}, ex.RunCalls()[1].Params)

	assert.Equal(t, `ssh umputun@master.radio-t.com "chmod 644 /data/archive/radio-t/media/%s"`, ex.RunCalls()[2].Cmd)
	assert.Equal(t, []any{"rt_podcast123.mp3"}, ex.RunCalls()[2].Params)

	assert.Equal(t, `ssh umputun@master.radio-t.com "find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'"`, ex.RunCalls()[3].Cmd)
	assert.Equal(t, 0, len(ex.RunCalls()[3].Params))

	assert.Equal(t, `ssh umputun@master.radio-t.com "docker exec -i ansible /srv/deploy_radiot.sh %d"`, ex.RunCalls()[4].Cmd)
	assert.Equal(t, []any{123}, ex.RunCalls()[4].Params)

	assert.Equal(t, `scp -P 2222 %s umputun@192.168.1.24:/data/archive.rucast.net/radio-t/media/"`, ex.RunCalls()[5].Cmd)
	assert.Equal(t, []any{"/tmp/rt_podcast123/rt_podcast123.mp3"}, ex.RunCalls()[5].Params)

	assert.Equal(t, `scp %s umputun@master.radio-t.com:/data/archive/radio-t/media/%s`, ex.RunCalls()[6].Cmd)
	assert.Equal(t, []any{"/tmp/rt_podcast123/rt_podcast123.mp3", "rt_podcast123.mp3"}, ex.RunCalls()[6].Params)
}
