package iFacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL     string
	CookieToken string
}

func (c *Client) Request(method, endPoint string, body io.Reader, f func(*http.Response) error) error {

	url := fmt.Sprintf("%s%s", c.BaseURL, endPoint)
	fmt.Printf("Forward request to iFacts: %s\n", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", c.CookieToken)
	req.Header.Set("Content-type", "application/json-patch+json")
	req.Header.Set("Accept", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// it's cleaner this way, since we should wait for f to return and then close the body
	err = f(resp)

	return err
}

func (c *Client) SearchByName(name string, f func(*http.Response) error) error {

	type bodyT struct {
		AssetName string `json:"assetName"`
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(bodyT{
		AssetName: name,
	}); err != nil {
		return err
	}

	return c.Request(http.MethodPost, "/api/v1/assets/search", &body, f)
}
