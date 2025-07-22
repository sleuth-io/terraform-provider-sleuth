package gqlclient

import (
	"context"
	"errors"

	"github.com/shurcooL/graphql"
)

// DeleteChangeSource - Deletes a change source
func (c *Client) DeleteChangeSource(ctx context.Context, projectSlug *string, slug *string) error {

	var m struct {
		DeleteChangeSource struct {
			Success graphql.Boolean
		} `graphql:"deleteChangeSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": DeleteChangeSourceMutationInput{ProjectSlug: *projectSlug, Slug: *slug},
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return err
	}

	if !m.DeleteChangeSource.Success {
		return errors.New("Missing change source")
	} else {
		return nil
	}
}
