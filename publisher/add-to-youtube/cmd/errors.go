package cmd

import "fmt"

func errSiteAPIRequest(err error) error {
	return fmt.Errorf("Error making site api request, got: %v", err)
}

func errJSONUnmarshal(err error) error {
	return fmt.Errorf("Error json unmarshaling, got: %v", err)
}
