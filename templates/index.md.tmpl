# Sleuth Provider
The Sleuth provider is used to interact with [Sleuth](https://sleuth.io) resources.

The provider allows you to manage your Sleuth organization's projects, environments, change sources, and impact sources easily.
It needs to be configured with the proper credentials before it can be used.

## Example Usage

Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    sleuth = {
      source = "sleuth-io/sleuth"
      version = "~> 0.2.0"
    }
  }
}

provider "sleuth" {
  api_key = "API_KEY_FROM_SLEUTH_ORGANIZATION_SETTINGS"
}
```

## Authentication

The Sleuth provider authenticates to Sleuth using your Sleuth organization API key. Find this by clicking on your
organization name at the top left, selecting "Organization settings", and looking under "Api key".


## Slugs

Sleuth resources often refer to a `slug` as a way to identify another resource. This is particularly important when
importing state from an existing Sleuth instance. Unfortunately, slugs are more of an internal Sleuth identifier and as
such, aren't exposed in the web UI.

There are slugs you may encounter and where you can find them in Sleuth:

* `project_slug` -- In the URL when looking at the project metrics or status, usually `https://app.sleuth.io/ORG_SLUG/PROJECT_SLUG`
* `environment_slug` -- In the URL when switching between environments, usually `?env_slug=ENVIRONMENT_SLUG`
* `change_source_slug` -- In the URL when looking at a specific code deployment, usually `https://app.sleuth.io/ORG_SLUG/CHANGE_SOURCE_SLUG`
* `impact_source_slug` -- In the URL when looking at a specific impact source, usually `https://app.sleuth.io/ORG_SLUG/PROJECT_SLUG/impacts/IMPACT_SOURCE_SLUG`

When importing resources, all but the project slug are prefixed by the project slug. For example, to import an environment named `production`, you'd run:

```
terraform import sleuth_environment.production PROJECT_SLUG/ENVIRONMENT_SLUG
```

## Schema

### Required

- **api_key** (String) The Sleuth organization's Api key
