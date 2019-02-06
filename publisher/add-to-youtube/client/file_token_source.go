package client

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type fileTokenSource struct {
	locker      sync.Locker
	pathToToken string
}

func newFileTokenSource(pathToToken string) *fileTokenSource {
	return &fileTokenSource{locker: &sync.Mutex{}, pathToToken: pathToToken}
}

// Token returns token from a file.
func (s *fileTokenSource) Token() (*oauth2.Token, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, err := os.Stat(s.pathToToken); err == os.ErrNotExist {
		return nil, errors.Wrap(err, "Need authorize an user")
	}

	t, err := readToken(s.pathToToken)
	if err != nil {
		return nil, err
	}

	return t, nil
}
