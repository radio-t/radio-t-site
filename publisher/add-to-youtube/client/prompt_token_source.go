package client

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type promptTokenSource struct {
	locker    sync.Locker
	config    *oauth2.Config
	tokenPath string
}

func newPromptTokenSource(tokenPath string, config *oauth2.Config) *promptTokenSource {
	return &promptTokenSource{locker: &sync.Mutex{}, config: config, tokenPath: tokenPath}
}

func (s *promptTokenSource) Token() (*oauth2.Token, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	authURL := s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	var code string
	fmt.Printf("Go to the following link in your browser. After completing "+
		"the authorization flow, enter the authorization code on the command "+
		"line: \n\n%v\n\n", authURL)

	fmt.Print("Enter the code here: ")

	if _, err := fmt.Scan(&code); err != nil {
		return nil, errors.Wrap(err, "Unable to read authorization code")
	}

	fmt.Print("\n")

	token, err := s.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve token")
	}

	if err := saveToken(s.tokenPath, token); err != nil {
		return nil, errors.Wrap(err, "Error saving token")
	}

	return token, nil
}
