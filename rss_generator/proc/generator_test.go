package proc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSSGenerator_MakeFeed(t *testing.T) {
	g := RSSGenerator{
		BaseArchiveURL: "https://archive.radio-t.com/media",
		BaseURL:        "https://radio-t.com",
		BaseCdnURL:     "https://cdn.rucast.net",
	}
	s := HugoSite{Location: "testdata/hugo", BaseURL: "https://radio-t.com"}

	posts, err := s.ReadPosts()
	require.NoError(t, err)
	res, err := g.MakeFeed(FeedConfig{
		Name:            "name1",
		Title:           "title1",
		Image:           "image1",
		Count:           10,
		Size:            false,
		FeedSubtitle:    "sub",
		FeedDescription: "desc",
	}, posts)
	require.NoError(t, err)
	t.Logf("%+v", res)
	require.Equal(t, 2, len(res.Items))
	assert.Equal(t, "https://radio-t.com/p/2023/12/23/podcast-889/", res.Items[0].URL)
	assert.Equal(t, "Радио-Т 889", res.Items[0].Title)
	assert.Equal(t, "https://radio-t.com/images/radio-t/rt889.jpg", res.Items[0].Image)
	assert.Equal(t, "https://cdn.rucast.net/rt_podcast889.mp3", res.Items[0].EnclosureURL)
	assert.Equal(t, 0, res.Items[0].FileSize)
	assert.Contains(t, res.Items[0].Description, "Apple останавливает продажи часов")
	assert.Contains(t, res.Items[0].Summary, "Apple останавливает продажи часов")
	assert.Contains(t, res.Items[0].ItunesSubtitle, "Apple останавливает продажи часов - 00:08:46. Весь код это технический долг - 00:27:21. Нет, это не так - 00:34:02.")

	assert.Equal(t, "https://radio-t.com/p/2023/03/18/podcast-850/", res.Items[1].URL)
}

func TestRSSGenerator_getMp3Size(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/testfile.mp3" && r.Method == "HEAD" {
			w.Header().Set("Content-Length", "1234")
			return
		}
		if r.URL.Path == "/other-file.mp3" && r.Method == "GET" {
			_, _ = w.Write([]byte("Hello world"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()

	client := &http.Client{Timeout: time.Second}

	g := RSSGenerator{
		Client:         client,
		BaseArchiveURL: testServer.URL,
	}

	t.Run("HEAD request", func(t *testing.T) {
		size, err := g.getMp3Size("testfile.mp3")
		require.NoError(t, err)
		assert.Equal(t, 1234, size)
	})

	t.Run("GET request", func(t *testing.T) {
		size, err := g.getMp3Size("other-file.mp3")
		require.NoError(t, err)
		assert.Equal(t, 11, size)
	})

}

func TestRSSGenerator_createItemData(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1234")
	}))
	defer testServer.Close()

	t.Run("long description", func(t *testing.T) {
		client := &http.Client{Timeout: time.Second}
		g := RSSGenerator{
			Client:         client,
			BaseArchiveURL: testServer.URL,
			BaseURL:        "https://example.com",
			BaseCdnURL:     "https://cdn.com",
		}

		data := strings.Repeat("1234567890", 50)
		res, err := g.createItemData(FeedConfig{
			Name:            "name1",
			Title:           "title1",
			Image:           "image1",
			Count:           10,
			Size:            false,
			FeedSubtitle:    "sub",
			FeedDescription: "desc",
		}, Post{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			URL:       "http://example.com/post1/",
			Config: map[string]interface{}{
				"title":      "title1",
				"filename":   "rt_podcast850",
				"categories": []interface{}{"podcast"},
				"image":      "https://radio-t.com/images/radio-t/rt850.jpg",
			},
			Data: data,
		})
		require.NoError(t, err)
		t.Logf("%+v", res)
		assert.Contains(t, res.ItunesSubtitle, "...")
		assert.Len(t, res.ItunesSubtitle, 240+3)
	})

	t.Run("no mp3 size", func(t *testing.T) {
		client := &http.Client{Timeout: time.Second}

		g := RSSGenerator{
			Client:         client,
			BaseArchiveURL: testServer.URL,
			BaseURL:        "https://example.com",
			BaseCdnURL:     "https://cdn.com",
		}

		res, err := g.createItemData(FeedConfig{
			Name:            "name1",
			Title:           "title1",
			Image:           "image1",
			Count:           10,
			Size:            false,
			FeedSubtitle:    "sub",
			FeedDescription: "desc",
		}, Post{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			URL:       "http://example.com/post1/",
			Config: map[string]interface{}{
				"title":      "title1",
				"filename":   "rt_podcast850",
				"categories": []interface{}{"podcast"},
				"image":      "https://radio-t.com/images/radio-t/rt850.jpg",
			},
			Data: "data1",
		})
		require.NoError(t, err)
		t.Logf("%+v", res)
		assert.Equal(t, "http://example.com/post1/", res.URL)
		assert.Equal(t, "Wed, 01 Jan 2020 00:00:00 UTC", res.Date)
		assert.Equal(t, "title1", res.Title)
		assert.Equal(t, "<p>data1</p>", res.Description)
		assert.Equal(t, "https://radio-t.com/images/radio-t/rt850.jpg", res.Image)
		assert.Equal(t, "https://cdn.com/rt_podcast850.mp3", res.EnclosureURL)
		assert.Equal(t, 0, res.FileSize)
	})

	t.Run("with mp3 size", func(t *testing.T) {
		client := &http.Client{Timeout: time.Second}

		g := RSSGenerator{
			Client:         client,
			BaseArchiveURL: testServer.URL,
			BaseURL:        "https://example.com",
			BaseCdnURL:     "https://cdn.com",
		}

		res, err := g.createItemData(FeedConfig{
			Name:            "name1",
			Title:           "title1",
			Image:           "image1",
			Count:           10,
			Size:            true,
			FeedSubtitle:    "sub",
			FeedDescription: "desc",
		}, Post{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			URL:       "http://example.com/post1/",
			Config: map[string]interface{}{
				"title":      "title1",
				"filename":   "rt_podcast850",
				"categories": []interface{}{"podcast"},
				"image":      "https://radio-t.com/images/radio-t/rt850.jpg",
			},
			Data: "data1",
		})
		require.NoError(t, err)
		t.Logf("%+v", res)
		assert.Equal(t, "http://example.com/post1/", res.URL)
		assert.Equal(t, "Wed, 01 Jan 2020 00:00:00 UTC", res.Date)
		assert.Equal(t, "title1", res.Title)
		assert.Equal(t, "<p>data1</p>", res.Description)
		assert.Equal(t, "https://radio-t.com/images/radio-t/rt850.jpg", res.Image)
		assert.Equal(t, "https://cdn.com/rt_podcast850.mp3", res.EnclosureURL)
		assert.Equal(t, 1234, res.FileSize)
	})
}

func TestRSSGenerator_htmlToPlainText(t *testing.T) {
	g := RSSGenerator{}

	t.Run("converts HTML content to plain text", func(t *testing.T) {
		htmlContent := "<p>Hello, World!</p>"
		expected := "Hello, World!"

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("handles empty HTML content", func(t *testing.T) {
		htmlContent := ""
		expected := ""

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("handles HTML content with multiple text nodes", func(t *testing.T) {
		htmlContent := "<p>Hello,</p><p>World!</p>"
		expected := "Hello,World!"

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("handles HTML content with multiple lines", func(t *testing.T) {
		htmlContent := "<p>Hello,</p><p>World!</p>\n<p>Another line</p>"
		expected := "Hello,World! Another line"

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("ignores HTML tags", func(t *testing.T) {
		htmlContent := "<p><strong>Hello,</strong> <em>World!</em></p>"
		expected := "Hello, World!"

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("real-life example", func(t *testing.T) {
		htmlContent := `
<p><img src="https://radio-t.com/images/radio-t/rt894.jpg" alt="" /></p>

<p><em>Темы</em><ul>
<li>Вступление - <em>00:00:00</em>.</li>
<li><a href="https://9to5mac.com/2024/01/25/third-party-default-browsers-engines/">Apple откроется в народ, но только в Европе</a> - <em>00:01:02</em>.</li>
<li><a href="https://trunk.io/blog/git-commit-messages-are-useless">Git commit messages</a> - <em>00:14:12</em>.</li>
<li><a href="https://www.docker.com/blog/introducing-docker-build-cloud/">Docker Build Cloud</a> - <em>00:49:23</em>.</li>
<li><a href="https://archive.ph/1waXO">Удаленная работа не испортила ничего</a> - <em>01:15:11</em>.</li>
<li><a href="https://www.hottakes.space/p/remote-work-won-dont-let-anyone-gaslight">Но, при этом, победила</a> - <em>01:16:15</em>.</li>
<li><a href="https://seykafu.medium.com/a-realistic-day-in-the-life-of-an-ai-product-manager-354d5b86318b">Тяжелая жизнь Product Manager</a> - <em>01:27:45</em>.</li>
<li><a href="https://radio-t.com/p/2024/01/23/prep-894/">Темы слушателей</a> - <em>01:51:47</em>.</li>
</ul></p>

<p><a href="https://cdn.radio-t.com/rt_podcast894.mp3">аудио</a> • <a href="https://chat.radio-t.com/logs/radio-t-894.html">лог чата</a>
<audio src="https://cdn.radio-t.com/rt_podcast894.mp3" preload="none"></audio></p>
`
		expected := "Темы Вступление - 00:00:00. Apple откроется в народ, но только в Европе - 00:01:02. Git commit messages - 00:14:12. Docker Build Cloud - 00:49:23. Удаленная работа не испортила ничего - 01:15:11. Но, при этом, победила - 01:16:15. Тяжелая жизнь Product Manager - 01:27:45. Темы слушателей - 01:51:47."

		result, err := g.htmlToPlainText(htmlContent)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
