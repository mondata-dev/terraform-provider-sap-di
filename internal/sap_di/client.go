package sap_di

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"encoding/base64"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		HostURL: *host,
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	return &c, nil
}

func basicAuth(username, password string) string {
  auth := username + ":" + password
  return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	// Note: this will have problems if there are redirects
	// see https://stackoverflow.com/a/31309385
	req.Header.Set("Authorization", "Basic " + basicAuth(c.Auth.Username, c.Auth.Password))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
