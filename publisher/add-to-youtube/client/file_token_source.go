package client

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type fileTokenSource struct {
	locker    sync.Locker
	tokenPath string
}

func newFileTokenSource(tokenPath string) *fileTokenSource {
	return &fileTokenSource{locker: &sync.Mutex{}, tokenPath: tokenPath}
}

// Token returns token from a file.
func (s *fileTokenSource) Token() (*oauth2.Token, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, err := os.Stat(s.tokenPath); err == os.ErrNotExist {
		return nil, errors.Wrap(err, "Required user authorization")
	}

	t, err := readToken(s.tokenPath)
	if err != nil {
		return nil, err
	}

	return t, nil
}
