# Plan: Add support for project_name in code change source resource build mappings

## Goal
Add support for a new optional parameter `project_name` (in addition to the existing `project_key`) for build mappings in the Terraform code change source resource.
- `project_name` is only used in GraphQL mutations (not returned in queries).
- Both `project_name` and `project_key` are optional.
- **Note:** If both are provided, `project_key` is preferred and takes precedence over `project_name`.

---

## Steps

### 1. Schema Changes
- Update the Terraform resource schema for `build_mappings` to add an optional `project_name` attribute.
- Ensure the attribute is documented as "used only for creation/update, not returned by API".
- **Document that if both are provided, `project_key` takes precedence.**

### 2. Model Changes
- Update the Go struct(s) representing the build mapping to include a `ProjectName` field.
- Ensure this field is optional and handled similarly to `ProjectKey`.
- **Logic should ensure that if both are set, `project_key` is used in preference to `project_name`.**

### 3. GraphQL Mutation Input
- Update the logic that builds the GraphQL mutation input for build mappings to include `project_name` if set.
- Ensure that both `project_key` and `project_name` are included in the mutation input if provided.
- **If both are set, only send `project_key` in the mutation input.**

### 4. Read/State Handling
- Since `project_name` is not returned by the API, ensure that:
  - The value is not expected in the state after a read.
  - The value is not set in the state from the API response.
  - The value is only used for create/update operations.

### 5. Documentation
- Update resource and attribute documentation to clarify the usage and limitations of `project_name` and `project_key`.
- **Explicitly state that `project_key` takes precedence if both are provided.**

### 6. Tests
- Add/Update tests to cover:
  - Supplying only `project_name`
  - Supplying only `project_key`
  - Supplying both (ensure `project_key` is used)
  - Supplying neither

---

## Summary Table

| Area           | Change                                                                 |
|----------------|------------------------------------------------------------------------|
| Schema         | Add `project_name` to build_mappings, optional, doc note, precedence   |
| Model          | Add `ProjectName` field to Go struct(s), precedence logic              |
| Mutation Input | Pass `project_name` in mutation if set, but prefer `project_key`       |
| Read/State     | Do not expect or set `project_name` in state from API                  |
| Docs           | Clarify usage and precedence of `project_name`/`project_key`           |
| Tests          | Add/Update tests for all combinations, check precedence                | 

---

## Files to Modify

1. **Resource Schema and Model**
   - `internal/sleuth/code_change_source_resource.go`
     - Update the Terraform resource schema for `build_mappings` to add the `project_name` attribute.
     - Update the Go struct for build mappings to include a `ProjectName` field.
     - Update logic for mutation input and precedence handling.

2. **GraphQL Client Types (if applicable)**
   - `internal/gqlclient/types.go` (or similar, depending on your codebase)
     - Update the GraphQL mutation input struct for build mappings to support `project_name`.
     - Ensure the mutation input logic respects the precedence of `project_key` over `project_name`.

3. **Tests**
   - `internal/sleuth/code_change_source_resource_test.go` (or wherever resource tests are located)
     - Add or update tests to cover all combinations of `project_key` and `project_name`.

4. **Documentation**
   - `docs/resources/code_change_source.md` (or wherever the Terraform resource documentation is maintained)
     - Update documentation to describe the new `project_name` attribute and precedence rules. 