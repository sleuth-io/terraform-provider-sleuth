package gqlclient

import (
	"context"

	// 	"encoding/json"
	"fmt"
	"github.com/shurcooL/graphql"
	"io/ioutil"
	"net/http"
	// 	"strings"
	"time"
)

// Client -
type Client struct {
	Baseurl    string
	HTTPClient *http.Client
	GQLClient  *graphql.Client
	ApiKey     string
	OrgSlug    string
}

type AuthenticatedTransport struct {
	T      http.RoundTripper
	ApiKey string
}

func (transport *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "sleuth-terraform")
	req.Header.Add("Authorization", fmt.Sprintf("apikey %s", transport.ApiKey))
	return transport.T.RoundTrip(req)
}

// NewClient -
func NewClient(baseurl, apiKey *string) (*Client, error) {
	httpClient := http.Client{Timeout: 10 * time.Second,
		Transport: &AuthenticatedTransport{http.DefaultTransport, *apiKey}}
	c := Client{
		GQLClient:  graphql.NewClient(*baseurl+"/graphql", &httpClient),
		HTTPClient: &httpClient,
		Baseurl:    *baseurl,
		ApiKey:     *apiKey,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) doQuery(query interface{}, variables map[string]interface{}) error {

	err := c.GQLClient.Query(context.Background(), query, variables)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) doMutate(query interface{}, variables map[string]interface{}) error {

	err := c.GQLClient.Mutate(context.Background(), query, variables)
	if err != nil {
		return err
	}
	return nil
}
