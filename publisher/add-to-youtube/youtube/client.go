package youtube

import (
	yt "google.golang.org/api/youtube/v3"
)

// Config represents a configuratio to YouTube client.
type Config struct {
	OAuth2    []byte
	TokenPath string
	Scopes    []string
}

// Client represents an client to YouTube service.
type Client struct {
	Config *Config
}

// New returns youtube client.
func New(config *Config) (*Client, error) {
	return &Client{config}, nil
}

// Upload uses an audio file to create a video file, then upload it with metadatas to Youtube.
func (c *Client) Upload(audioPath, title, description, category, keywords, privacy string) (*yt.Video, error) {
	return upload(c, audioPath, title, description, category, keywords, privacy)
}

// Authorize authorizes an user in Youtube service.
func (c *Client) Authorize() error {
	return authorize(c)
}
