package iFacts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func (c *client) request(method, endPoint string, body io.Reader, f func(*http.Response) error) error {

	tkn := c.getToken()

	url := fmt.Sprintf("%s%s", c.BaseURL, endPoint)
	fmt.Printf("Forward request to iFacts: %s\n", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	//req.Header.Set("Cookie", c.CookieToken)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if dmp, err := httputil.DumpResponse(resp, true); err == nil {
		fmt.Println("Response:")
		fmt.Println(string(dmp))
	} else {
		log.Println(err)
	}

	// it's cleaner this way, since we should wait for f to return and then close the body
	err = f(resp)

	return err
}

func (c *client) getToken() string {

	c.m.Lock()
	defer c.m.Unlock()

	if c.token != "" && time.Since(c.tokenLastRequested) <= c.tokenLifespan {
		return c.token
	}

	c.tokenLastRequested = time.Now()

	tkn, err := c.requestToken()
	if err != nil {
		c.token = ""
		log.Println(err)
	} else {
		c.token = tkn
	}

	return c.token
}

func (c *client) requestToken() (string, error) {

	fmt.Println("Requesting new token")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	enc := "application/x-www-form-urlencoded"

	resp, err := http.Post(c.BaseURL+"/idp/connect/token", enc, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bdy, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tokenResp := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Error       string `json:"error"`
	}{}
	if err := json.Unmarshal(bdy, &tokenResp); err != nil || tokenResp.Error != "" {
		return "", fmt.Errorf("iFacts auth response: %s", string(bdy))
	}

	if tokenResp.ExpiresIn != 0 {
		c.tokenLifespan = time.Duration(tokenResp.ExpiresIn) * time.Second
	}

	return tokenResp.AccessToken, nil
}
