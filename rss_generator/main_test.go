package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/radio-t/radio-t-site/rss_generator/mocks"
)

func TestCreateItemData(t *testing.T) {
	mockClient := &mocks.HttpClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Length", "1234")
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	tests := []struct {
		name      string
		client    HttpClient
		post      Post
		feed      Feed
		baseURL   string
		expect    ItemData
		expectErr bool
	}{
		{
			name:   "Valid Post",
			client: mockClient,
			post: Post{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				URL:       "post-url",
				Config: map[string]any{
					"title":    "Test Post",
					"filename": "testfile.mp3",
					"image":    "testimage.jpg",
				},
				Data: "Test content",
			},
			feed: Feed{
				Size: true,
			},
			baseURL: "https://example.com",
			expect: ItemData{
				Title:        "Test Post",
				Description:  "Test content",
				URL:          "https://example.com/post-url",
				Date:         "Wed, 01 Jan 2020 00:00:00 MST",
				Summary:      "Test content",
				Image:        "testimage.jpg",
				EnclosureURL: "https://example.com/testfile.mp3",
				FileSize:     "1234",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := createItemData(tc.client, tc.post, tc.feed, tc.baseURL)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, result)
			}
		})
	}
}
func TestGetMp3Size(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "12345") // Set a test Content-Length value
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Replace the URL with the test server's URL in your function call
	client := &http.Client{}

	size, err := getMp3Size(client, testServer.URL+"/test.mp3")
	assert.NoError(t, err)
	assert.Equal(t, "12345", size, "Expected size to match the Content-Length set in the test server")
}
