package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/bogem/id3v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/publisher/cmd/mocks"
)

func TestUpload_Do(t *testing.T) {
	ex := &mocks.ExecutorMock{
		RunFunc: func(cmd string, params ...string) {},
	}

	d := Upload{
		Executor: ex,
		Location: "/tmp",
	}

	err := d.Do(123)
	require.NoError(t, err)

	require.Equal(t, 2, len(ex.RunCalls()))
	assert.Equal(t, "spot", ex.RunCalls()[0].Cmd)
	assert.Equal(t, []string{"-e mp3:/tmp/rt_podcast123/rt_podcast123.mp3", "--task=\"deploy to master", "-v", "/tmp/rt_podcast123/rt_podcast123.mp3"}, ex.RunCalls()[0].Params)

	assert.Equal(t, "spot", ex.RunCalls()[1].Cmd)
	assert.Equal(t, 4, len(ex.RunCalls()[1].Params))
	assert.Equal(t, []string{"-e mp3:/tmp/rt_podcast123/rt_podcast123.mp3", "--task=\"deploy to nodes\"", "-v", "/tmp/rt_podcast123/rt_podcast123.mp3"}, ex.RunCalls()[1].Params)
}

func TestUpload_setMp3Tags(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tags")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dest := tempDir + "/rt_podcast123"
	err = os.MkdirAll(dest, 0o755)
	require.NoError(t, err)

	// copy test file to dest
	src, err := os.Open("testdata/test.mp3")
	require.NoError(t, err)
	defer src.Close()
	dst, err := os.Create(dest + "/rt_podcast123.mp3")
	require.NoError(t, err)
	defer dst.Close()
	_, err = io.Copy(dst, src)
	require.NoError(t, err)

	u := Upload{Location: tempDir}
	err = u.setMp3Tags(123)
	require.NoError(t, err)

	tag, err := id3v2.Open(dst.Name(), id3v2.Options{Parse: true})
	require.NoError(t, err)
	assert.Equal(t, "Радио-Т 123", tag.Title())
	assert.Equal(t, "Umputun, Bobuk, Gray, Ksenks, Alek.sys", tag.Artist())
	assert.Equal(t, "Радио-Т", tag.Album())
	assert.Equal(t, fmt.Sprintf("%d", time.Now().Year()), tag.Year())
	assert.Equal(t, "Podcast", tag.Genre())
}
