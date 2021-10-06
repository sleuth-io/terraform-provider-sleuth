package gqlclient

import (
	"errors"
	"github.com/shurcooL/graphql"
)

// DeleteImpactSource - Deletes a impact source
func (c *Client) DeleteImpactSource(projectSlug *string, slug *string) error {

	var m struct {
		DeleteImpactSource struct {
			Success graphql.Boolean
		} `graphql:"deleteImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": DeleteImpactSourceMutationInput{ProjectSlug: *projectSlug, Slug: *slug},
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return err
	}

	if !m.DeleteImpactSource.Success {
		return errors.New("Missing impact source")
	} else {
		return nil
	}
}
