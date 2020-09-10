package main

import (
	"fmt"
)

// McAPI offers methods to talk to the MailChimp API
type McAPI struct{}

// NewMcAPI creates a new McAPI
func NewMcAPI() *McAPI {
	return &McAPI{}
}

func (api *McAPI) AddSubscriber(email string) error {
	fmt.Println(email)

	return nil
}
