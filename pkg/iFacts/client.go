package iFacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

type Client interface {
	SearchByName(name string, f func(*http.Response) error) error
}

type client struct {
	BaseURL      string
	CookieToken  string
	clientID     string
	clientSecret string

	m                  *sync.Mutex // protect the token
	token              string
	tokenLastRequested time.Time
}

var tokenLifespan = time.Hour

func NewClient(baseURL, clientID, secret string) Client {
	fmt.Printf("iFacts client for %s [%s]\n", baseURL, clientID)
	return &client{
		BaseURL:      baseURL,
		clientID:     clientID,
		clientSecret: secret,
		m:            new(sync.Mutex),
	}
}

func (c *client) getToken() string {

	c.m.Lock()
	defer c.m.Unlock()

	if c.token != "" && time.Since(c.tokenLastRequested) <= tokenLifespan {
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

	b := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s", "client_credentials", c.clientID, c.clientSecret)

	resp, err := http.Post(c.BaseURL+"/idp/connect/token", "application/x-www-form-urlencoded", strings.NewReader(b))
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
		return "", fmt.Errorf("iFacts auth response: %s", string(b))
	}

	if tokenResp.ExpiresIn != 0 {
		tokenLifespan = time.Duration(tokenResp.ExpiresIn) * time.Second
	}

	return tokenResp.AccessToken, nil
}

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
	req.Header.Set("Content-type", "application/json-patch+json")
	req.Header.Set("Accept", "text/plain")

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

func (c *client) SearchByName(name string, f func(*http.Response) error) error {

	type bodyT struct {
		AssetName                    string `json:"assetName"`
		IncludeInactiveOrganizations bool   `json:"includeInactiveOrganizations"`
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(bodyT{
		AssetName:                    name,
		IncludeInactiveOrganizations: false,
	}); err != nil {
		return err
	}

	return c.request(http.MethodPost, "/api/v1/assets/search", &body, f)
}

func (c *client) AssetGeneralSection(id string, f func(*http.Response) error) error {

	return c.request(http.MethodGet, fmt.Sprintf("/api/v1/assets/getgeneralsection/%s", id), nil, f)
}
