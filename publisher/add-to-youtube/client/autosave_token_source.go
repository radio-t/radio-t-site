package client

import (
	"sync"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"
)

// autoSaveTokenSource represents a token source, that save token when it changed in TokenSource.
type autoSaveTokenSource struct {
	oauth2.TokenSource
	current   *oauth2.Token
	previous  *oauth2.Token
	locker    sync.Locker
	tokenPath string
}

// newAutoSaveTokenSource returns autoSaveTokenSource.
func newAutoSaveTokenSource(tokenPath string, t *oauth2.Token, ts oauth2.TokenSource) *autoSaveTokenSource {
	return &autoSaveTokenSource{
		locker:      &sync.Mutex{},
		tokenPath:   tokenPath,
		previous:    t,
		TokenSource: ts,
	}
}

// Token returns a token and saves it if changed.
func (s *autoSaveTokenSource) Token() (*oauth2.Token, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	var err error
	s.current, err = s.TokenSource.Token()
	if err != nil {
		return nil, err
	}
	if s.previous == s.current {
		return s.current, nil
	}
	if err := saveToken(s.tokenPath, s.current); err != nil {
		return nil, errors.Wrap(err, "Error saving token")
	}
	s.previous = s.current

	return s.current, nil
}
