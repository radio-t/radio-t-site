package youtube

import "fmt"

func errYoutubeClientCreate(err error) error {
	return fmt.Errorf("Error creating YouTube client: %v", err)
}

func errAPICall(err error) error {
	return fmt.Errorf("Error making API call, got: %v", err)
}
