# Plan: Add Terraform Support for Sleuth Teams Resource

## Overview
This plan describes the steps to add a new Terraform resource for managing Sleuth teams and subteams, including membership by user email, using the remote Sleuth GraphQL API as defined in `schema.graphql`. The implementation follows the style and structure of existing resources such as `internal/sleuth/project_resource.go` and `internal/sleuth/project_resource_test.go`.

## Steps

### 1. Define the Team Resource Schema and Model
- **File:** `internal/sleuth/team_resource.go`
- Created a new resource struct and model for `teamResource`.
- Attributes:
  - `id` (computed, string)
  - `name` (required, string)
  - `slug` (computed, string)
  - `parent_slug` (optional, string, for subteams)
  - `members` (optional, list of string, user emails)
- Implemented schema, metadata, and configure methods similar to `projectResource`.

### 2. Implement CRUD Operations for Team Resource
- **File:** `internal/sleuth/team_resource.go`
- Implemented `Create`, `Read`, `Update`, `Delete`, and `ImportState` methods.
- Used the following GraphQL mutations/queries from `schema.graphql`:
  - `createTeam(input: CreateTeamMutationInput!)`
  - `updateTeam(input: UpdateTeamMutationInput!)`
  - `deleteTeam(input: DeleteTeamMutationInput!)`
  - `addTeamMembers(input: AddTeamMembersMutationInput!)`
  - `removeTeamMembers(input: RemoveTeamMembersMutationInput!)`
  - `team(teamSlug: ID!, orgSlug: ID): TeamType!` (for read)
- For member management:
  - On create/update, resolved user emails to user IDs (using `organization.users` query) and called `addTeamMembers`/`removeTeamMembers` as needed.
  - Fixed the GraphQL query structure for resolving users by email during implementation.
- State consistency:
  - Ensured that `parent_slug` and `members` fields in the Terraform state remain consistent, even if the API omits or changes values.
  - After create/update, fetched the actual list of team members from the API to set the state accurately.

### 3. Extend GraphQL Client for Team Operations
- **File:** `internal/gqlclient/client.go` (and related files)
- Added methods for:
  - Creating, updating, deleting teams
  - Adding/removing team members by user ID
  - Querying users by email (using `organization.users`)
  - Querying team details
- Ensured all necessary input/output structs are defined.
- Fixed issues with input struct configuration and GraphQL query parameters during implementation.

### 4. Add Acceptance and Unit Tests
- **File:** `internal/sleuth/team_resource_test.go`
- Added acceptance tests for:
  - Creating a team
  - Creating a subteam (with `parent_slug`)
  - Adding/removing members by email
  - Updating team name and membership
  - Deleting a team
- Used a similar structure to `project_resource_test.go`.
- Added unit tests for helper functions as needed.

### 5. Update Provider Registration
- **File:** `internal/sleuth/provider.go`
- Registered the new `teamResource` in the provider's resource map.

### 6. Update Changelog
- **File:** `CHANGELOG.md`
- Added an entry describing the new team resource, subteam support, and membership management by email.

### 7. (Optional) Add Example Terraform Configurations
- **File:** `examples/resources/terraform_team_example.tf` (new file)
- Provided example usage for teams, subteams, and membership.

---

## Implementation Notes / Changes from Original Plan
- Debug print statements were added throughout the implementation for troubleshooting (e.g., printing GraphQL input/output, resolved user IDs, etc.).
- All debug output was removed from the final code after tests passed, ensuring clean production code.
- Special care was taken to ensure Terraform state consistency for fields like `parent_slug` and `members`, including fetching the actual member list from the API after create/update.
- The GraphQL query for resolving users by email was corrected to match the API's requirements, fixing issues encountered during testing.
- No attempts to get a user's team slugs from GraphQL remain in the codebase, as per later requirements.
- The final implementation closely matches the plan, with additional robustness and cleanup based on test results and troubleshooting experience.

## References
- `internal/sleuth/project_resource.go` (resource implementation example)
- `internal/sleuth/project_resource_test.go` (acceptance test example)
- `schema.graphql` (GraphQL API reference)
- `CHANGELOG.md` (for documenting changes) 