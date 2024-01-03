package proc

import (
	"net/http"
	"net/http/httptest"
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

	assert.Equal(t, "https://radio-t.com/p/2023/03/18/podcast-850/", res.Items[1].URL)
}

func TestRSSGenerator_getMp3Size(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/testfile.mp3" && r.Method == "HEAD" {
			w.Header().Set("Content-Length", "1234")
			return
		}
		if r.URL.Path == "/other-file.mp3" && r.Method == "GET" {
			w.Write([]byte("Hello world"))
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
		return
	}))
	defer testServer.Close()

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
