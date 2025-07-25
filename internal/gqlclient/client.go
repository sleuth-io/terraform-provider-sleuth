package gqlclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"
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
	UA     string
}

func (transport *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", transport.UA)
	req.Header.Add("Authorization", fmt.Sprintf("apikey %s", transport.ApiKey))
	return transport.T.RoundTrip(req)
}

// NewClient -
func NewClient(baseurl, apiKey *string, ua string, timeout time.Duration) (*Client, error) {
	httpClient := http.Client{Timeout: timeout,
		Transport: &AuthenticatedTransport{http.DefaultTransport, *apiKey, ua}}
	c := Client{
		GQLClient:  graphql.NewClient(*baseurl+"/graphql", &httpClient),
		HTTPClient: &httpClient,
		Baseurl:    *baseurl,
		ApiKey:     *apiKey,
	}

	return &c, nil
}

func (c *Client) doQuery(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	err := c.GQLClient.Query(ctx, query, variables)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) doMutate(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	err := c.GQLClient.Mutate(ctx, query, variables)
	if err != nil {
		return err
	}
	return nil
}
