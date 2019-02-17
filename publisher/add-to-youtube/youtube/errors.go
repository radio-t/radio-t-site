package youtube

import "github.com/pkg/errors"

func errOAuth2HTTPClientCreate(err error) error {
	return errors.Wrap(err, "Error creating OAuth2 http client")
}

func errYoutubeClientCreate(err error) error {
	return errors.Wrap(err, "Error creating YouTube client")
}
