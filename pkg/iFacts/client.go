package iFacts

import (
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return f(resp)
}
