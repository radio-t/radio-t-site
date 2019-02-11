package client

import (
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// New returns a configured oauth2 http client.
func New(oauth2Config []byte, pathToToken string, skipPrompt bool, scopes ...string) (*http.Client, error) {
	ctx := context.Background()

	config, err := google.ConfigFromJSON(oauth2Config, scopes...)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshal config from json")
	}

	s1 := newFileTokenSource(pathToToken)

	t, err := s1.Token()
	fileNotExist := os.IsNotExist(errors.Cause(err))
	if skipPrompt && fileNotExist {
		return nil, errors.Wrap(err, "Required user authorization")
	}
	if fileNotExist {
		s5 := newPromptTokenSource(pathToToken, config)
		t, err = s5.Token()
	}
	if err != nil {
		return nil, err
	}

	s2 := config.TokenSource(ctx, t)

	s3 := oauth2.ReuseTokenSource(t, s2)

	s4 := newAutoSaveTokenSource(pathToToken, t, s3)

	client := oauth2.NewClient(ctx, s4)

	return client, nil
}
