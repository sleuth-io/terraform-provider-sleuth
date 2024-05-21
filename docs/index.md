# Sleuth Provider
The Sleuth provider is used to interact with [Sleuth](https://sleuth.io) resources.

The provider allows you to manage your Sleuth organization's projects, environments, change sources, and impact sources easily.
It needs to be configured with the proper credentials before it can be used.

## Versions explained:
Since there are different versions it make sense to explain the differences between them.

Terraform provider version 0.4.* is meant to be used by older Terraform versions (0.12.* and on). It is feature complete and will only receive bug fixes.
Terraform provider version 0.5.* is a first step towards rewrite of a provider. It has the same features as 0.4.* but there might be a few missing pieces.
Terraform provider version 0.6.* is a final rewrite of a provider. It has the same features as 0.4.* and will receive new features and bug fixes. It is written in Terraform plugin framework and only works with newer versions of Terraform (1.* and on).

If you are using older terraform version, you should use 0.4 version, if you are using newer terraform version, you should use 0.6 version. Version 0.5 is just a transition version and should not be used regularly.


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

You can provide the API key via `SLEUTH_API_KEY` environment variable. Note that you cannot use a *personal* API key, doing so will result in a *403 Forbidden* error. The key must be an *Organization* API key.

### How to generate an Organization API key

1. Navigate to "Organization Settings" in Sleuth dashboard
1. Under the "Details" tab, there is an "Api Key" field that contains the desired key


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

## Resources deletion caveats

Due to the way Sleuth API works, there are some caveats when deleting resources. When a project resource is created, a default environment is created as well (called `Production`).
If you want to delete the default environment, you will have to do it manually in the Sleuth UI. If you only delete environments in Terraform, the project will be left with a default environment (`Production`) even though the environment will not be present in the state.
This can cause issues when you will try to create `code_change_source` and other resources.

## Schema

### Required

- **api_key** (String) The Sleuth organization's Api key
