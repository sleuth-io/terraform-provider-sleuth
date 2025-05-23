package gqlclient

// For team members query
// TeamMembersQueryResult is used for unmarshalling the team members GraphQL query
// (add this if not already present)
type TeamMembersQueryResult struct {
	Organization struct {
		Team struct {
			Members struct {
				Objects []struct {
					Email string `graphql:"email"`
				} `graphql:"objects"`
			} `graphql:"members(page: $page, pageSize: $pageSize)"`
		} `graphql:"team(teamSlug: $teamSlug)"`
	} `graphql:"organization(orgSlug: $orgSlug)"`
}
