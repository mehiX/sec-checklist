package iFacts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mehix/sec-checklist/pkg/domain"
)

type Client interface {
	SearchByName(name string, f func(*http.Response) error) error
	GetClassifications(id string) ([]domain.Classification, error)
}

type client struct {
	BaseURL      string
	CookieToken  string
	clientID     string
	clientSecret string

	m                  *sync.Mutex // protect the token
	token              string
	tokenLastRequested time.Time
	tokenLifespan      time.Duration
}

func NewClient(baseURL, clientID, secret string) Client {
	fmt.Printf("iFacts client for %s [%s]\n", baseURL, clientID)
	c := &client{
		BaseURL:       baseURL,
		clientID:      clientID,
		clientSecret:  secret,
		m:             new(sync.Mutex),
		tokenLifespan: time.Hour,
	}

	go func() {
		if _, err := c.requestToken(); err != nil {
			log.Printf("Could not fetch iFacts access token: %v\n", err)
		} else {
			fmt.Println("Got iFacts access token")
		}
	}()

	return c
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

func (c *client) GetClassifications(id string) ([]domain.Classification, error) {

	var classifications []domain.Classification

	readClassifications := func(resp *http.Response) error {
		respData := struct {
			Classifications []domain.Classification `json:"Classifications"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return err
		}

		classifications = respData.Classifications

		return nil
	}

	if err := c.request(http.MethodGet, fmt.Sprintf("/api/v1/assets/getclassifications/%s", id), nil, readClassifications); err != nil {
		return nil, err
	}

	return classifications, nil
}
