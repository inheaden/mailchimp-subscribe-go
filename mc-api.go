package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
)

// McAPI offers methods to talk to the MailChimp API
type McAPI struct {
	config Config
}

// NewMcAPI creates a new McAPI
func NewMcAPI(config Config) *McAPI {
	return &McAPI{config: config}
}

func (api *McAPI) AddSubscriber(email string) error {
	fmt.Println(email)

	client := &http.Client{}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members/%s?skip_merge_validation=true", api.config.APIServer, api.config.ListID, hash)

	log.Debugf("Using URL: %s", url)

	body := fmt.Sprintf(`{ "email_address": "%s", "status_if_new": "%s" }`, email, api.config.StatusIfNew)
	log.Debugf("Using Body: %s", body)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Authorization", fmt.Sprintf(`Basic %s`, api.config.APIKey))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Debugf("response Status: %s", resp.Status)

	if resp.StatusCode != 200 {
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error Body:", string(respBody))

		return errors.New(string(respBody))
	}

	return nil
}
