// Package youtube provides functions to authorize at YouTube service and make video from podcast then upload it to YouTube.
package youtube

import (
	yt "google.golang.org/api/youtube/v3"
)

var scopes = []string{yt.YoutubeUploadScope}

// Client represents an client to YouTube service.
type Client struct {
	OAuth2    []byte
	TokenPath string
	Scopes    []string
	CoverPath string
}

// New returns youtube client.
func New(oauth2 []byte, tokenPath, coverPath string) (*Client, error) {
	return &Client{OAuth2: oauth2, TokenPath: tokenPath, CoverPath: coverPath, Scopes: scopes}, nil
}

// Upload uses an audio file to create a video file, then upload it with metadatas to Youtube.
func (c *Client) Upload(audioPath, title, description, category, keywords, privacy string) (*yt.Video, error) {
	return upload(c, audioPath, title, description, category, keywords, privacy)
}

// Authorize authorizes an user in Youtube service.
func (c *Client) Authorize() error {
	return authorize(c)
}
