package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/bogem/id3v2/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radio-t/radio-t-site/publisher/app/cmd/mocks"
)

func TestUpload_Do(t *testing.T) {
	tempDir := "/tmp/publisher_test"
	defer os.RemoveAll(tempDir)

	dest := tempDir + "/rt_podcast123"
	err := os.MkdirAll(dest, 0o755)
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

	ex := &mocks.ExecutorMock{
		RunFunc: func(cmd string, params ...string) {},
	}

	d := Upload{
		Executor:     ex,
		LocationMp3:  tempDir,
		LocationPost: "testdata",
	}

	err = d.Do(123)
	require.NoError(t, err)

	require.Equal(t, 2, len(ex.RunCalls()))
	assert.Equal(t, "spot", ex.RunCalls()[0].Cmd)
	assert.Equal(t, []string{"-e mp3:/tmp/publisher_test/rt_podcast123/rt_podcast123.mp3", "--task=\"deploy to master", "-v", "/tmp/publisher_test/rt_podcast123/rt_podcast123.mp3"}, ex.RunCalls()[0].Params)

	assert.Equal(t, "spot", ex.RunCalls()[1].Cmd)
	assert.Equal(t, 4, len(ex.RunCalls()[1].Params))
	assert.Equal(t, []string{"-e mp3:/tmp/publisher_test/rt_podcast123/rt_podcast123.mp3", "--task=\"deploy to nodes\"", "-v", "/tmp/publisher_test/rt_podcast123/rt_podcast123.mp3"}, ex.RunCalls()[1].Params)
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

	t.Run("without chapters", func(t *testing.T) {
		u := Upload{LocationMp3: tempDir}
		err = u.setMp3Tags(123, nil)
		require.NoError(t, err)

		tag, err := id3v2.Open(dst.Name(), id3v2.Options{Parse: true})
		require.NoError(t, err)
		assert.Equal(t, "Радио-Т 123", tag.Title())
		assert.Equal(t, "Umputun, Bobuk, Gray, Ksenks, Alek.sys", tag.Artist())
		assert.Equal(t, "Радио-Т", tag.Album())
		assert.Equal(t, fmt.Sprintf("%d", time.Now().Year()), tag.Year())
		assert.Equal(t, "Podcast", tag.Genre())
	})

	t.Run("with chapters", func(t *testing.T) {
		u := Upload{LocationMp3: tempDir}
		err = u.setMp3Tags(123, []chapter{
			{"Chapter One", "http://example.com/one", time.Second},
			{"Chapter Two", "http://example.com/two", time.Second * 5},
		})
		require.NoError(t, err)

		tag, err := id3v2.Open(dst.Name(), id3v2.Options{Parse: true})
		require.NoError(t, err)
		assert.Equal(t, "Радио-Т 123", tag.Title())
		assert.Equal(t, "Umputun, Bobuk, Gray, Ksenks, Alek.sys", tag.Artist())
		assert.Equal(t, "Радио-Т", tag.Album())
		assert.Equal(t, fmt.Sprintf("%d", time.Now().Year()), tag.Year())
		assert.Equal(t, "Podcast", tag.Genre())

		chapterFrames := tag.GetFrames(tag.CommonID("CHAP"))
		require.Len(t, chapterFrames, 2)
		assert.Equal(t, "1", chapterFrames[0].(id3v2.ChapterFrame).ElementID)
		assert.Equal(t, "Chapter One", chapterFrames[0].(id3v2.ChapterFrame).Title.Text)
		assert.Equal(t, time.Second, chapterFrames[0].(id3v2.ChapterFrame).StartTime)
		assert.Equal(t, 5*time.Second, chapterFrames[0].(id3v2.ChapterFrame).EndTime)
		t.Logf("%+v", chapterFrames[0].(id3v2.ChapterFrame))
	})

}

func TestUpload_parseChapters(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    []chapter
		expectError bool
	}{
		{
			name: "Valid chapters",
			content: `
- [Chapter One](http://example.com/one) - *00:01:00*.
- [Chapter Two](http://example.com/two) - *00:02:30*.
`,
			expected: []chapter{
				{"Chapter One", "http://example.com/one", time.Minute},
				{"Chapter Two", "http://example.com/two", time.Minute*2 + time.Second*30},
			},
			expectError: false,
		},
		{
			name: "Invalid timestamp format",
			content: `
- [Chapter One](http://example.com/one) - *00:100*.
`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty content",
			content:     "",
			expected:    []chapter{},
			expectError: false,
		},
	}

	u := &Upload{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := u.parseChapters(tc.content)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestUpload_parseChaptersWithRealData(t *testing.T) {
	realDataContent := `
+++
title = "Радио-Т 686"
date = "2020-01-25T18:02:42"
categories = ["podcast"]
image = "https://radio-t.com/images/radio-t/rt686.jpg"
filename = "rt_podcast686"
+++

![](https://radio-t.com/images/radio-t/rt686.jpg)

- [Первому Macintosh 36 лет](https://www.macrumors.com/2020/01/24/macintosh-36th-anniversary/) - *00:04:18*.
- [JetBrains придумает новую IntelliJ](https://devclass.com/2020/01/21/jetbrains-reimagines-intellij-as-text-editor-machine-learning/) - *00:11:10*.
- [Мы учим не тому](https://www.bloomberg.com/tosv2.html?vid=&uuid=59d32d10-31cd-11ea-a482-59e1177b04c0&url=L29waW5pb24vYXJ0aWNsZXMvMjAyMC0wMS0wNy9jb2RpbmctaXMtY29sbGFib3JhdGl2ZS1hbmQtc3RlbS1lZHVjYXRpb24tc2hvdWxkLWJlLXRvbw==) - *00:28:16*.
- [Google всех опять разозлил](https://www.theverge.com/tech/2020/1/24/21079696/google-serp-design-change-altavisa-ads-trust) - *00:45:31*.
- [Google пообещал исправиться](https://www.businessinsider.com/google-walks-back-search-results-design-changes-following-criticism-2020-1) - *00:52:02*.
- [Sonos заразит всю систему](https://www.extremetech.com/electronics/305263-sonos-frantic-flailing-illustrates-the-stupidity-of-smart-tech) - *00:56:36*.
- [Как начать работать когда не хочется](https://www.deprocrastination.co/blog/3-tricks-to-start-working-despite-not-feeling-like-it) - *01:16:03*.
- [Apple отказалась от полного шифрования](https://www.reuters.com/article/us-apple-fbi-icloud-exclusive-idUSKBN1ZK1CT) - *01:41:06*.
- [Темы слушателей](https://radio-t.com/p/2020/01/21/prep-686/) - *01:51:56*.

*Спонсор этого выпуска [DigitalOcean](https://do.co/radiot)*

[аудио](https://cdn.radio-t.com/rt_podcast686.mp3) • [лог чата](https://chat.radio-t.com/logs/radio-t-686.html)
<audio src="https://cdn.radio-t.com/rt_podcast686.mp3" preload="none"></audio>
`

	expectedChapters := []chapter{
		{"Первому Macintosh 36 лет", "https://www.macrumors.com/2020/01/24/macintosh-36th-anniversary/", 4*time.Minute + 18*time.Second},
		{"JetBrains придумает новую IntelliJ", "https://devclass.com/2020/01/21/jetbrains-reimagines-intellij-as-text-editor-machine-learning/", 11*time.Minute + 10*time.Second},
		{"Мы учим не тому", "https://www.bloomberg.com/tosv2.html?vid=&uuid=59d32d10-31cd-11ea-a482-59e1177b04c0&url=L29waW5pb24vYXJ0aWNsZXMvMjAyMC0wMS0wNy9jb2RpbmctaXMtY29sbGFib3JhdGl2ZS1hbmQtc3RlbS1lZHVjYXRpb24tc2hvdWxkLWJlLXRvbw==", 28*time.Minute + 16*time.Second},
		{"Google всех опять разозлил", "https://www.theverge.com/tech/2020/1/24/21079696/google-serp-design-change-altavisa-ads-trust", 45*time.Minute + 31*time.Second},
		{"Google пообещал исправиться", "https://www.businessinsider.com/google-walks-back-search-results-design-changes-following-criticism-2020-1", 52*time.Minute + 2*time.Second},
		{"Sonos заразит всю систему", "https://www.extremetech.com/electronics/305263-sonos-frantic-flailing-illustrates-the-stupidity-of-smart-tech", 56*time.Minute + 36*time.Second},
		{"Как начать работать когда не хочется", "https://www.deprocrastination.co/blog/3-tricks-to-start-working-despite-not-feeling-like-it", 76*time.Minute + 3*time.Second},
		{"Apple отказалась от полного шифрования", "https://www.reuters.com/article/us-apple-fbi-icloud-exclusive-idUSKBN1ZK1CT", 101*time.Minute + 6*time.Second},
		{"Темы слушателей", "https://radio-t.com/p/2020/01/21/prep-686/", 111*time.Minute + 56*time.Second},
	}

	u := &Upload{}

	result, err := u.parseChapters(realDataContent)
	assert.NoError(t, err)
	assert.Equal(t, expectedChapters, result)
}