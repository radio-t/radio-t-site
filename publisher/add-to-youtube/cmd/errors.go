package cmd

import "github.com/pkg/errors"

func errSiteAPIRequest(err error) error {
	return errors.Wrap(err, "Error making site api request, got: %v")
}

func errJSONUnmarshal(err error) error {
	return errors.Wrap(err, "Error json unmarshaling, got: %v")
}
