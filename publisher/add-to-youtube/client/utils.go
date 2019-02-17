package client

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
)

// readToken retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func readToken(pathToToken string) (*oauth2.Token, error) {
	log.Info("Trying to read token from fs")

	b, err := ioutil.ReadFile(pathToToken)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading token file")
	}

	var token oauth2.Token
	if err := json.Unmarshal(b, &token); err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling token from json")
	}

	log.Info("Token read")
	return &token, nil
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(pathToToken string, token *oauth2.Token) error {
	log.Info("Trying to save token")
	log.Infof("Saving token to: %s", pathToToken)

	b, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "Error marshaling token to json")
	}

	if err := ioutil.WriteFile(pathToToken, b, 0600); err != nil {
		return errors.Wrap(err, "Error writing token to fs")
	}
	log.Info("Token saved")
	return nil
}
