package youtube

import (
	"golang.org/x/oauth2"

	yt "google.golang.org/api/youtube/v3"
)

// Config represents a configuratio to YouTube client.
type Config struct {
	OAuth2      *oauth2.Config
	SecretsPath string
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
	return upload(c.Config.OAuth2, audioPath, title, description, category, keywords, privacy, c.Config.SecretsPath)
}

// Authorize authorizes an user in Youtube service.
func (c *Client) Authorize() error {
	return authorize(c.Config.OAuth2, c.Config.SecretsPath)
}
