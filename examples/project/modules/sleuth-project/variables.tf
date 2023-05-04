variable "project" {
  description = <<EOF
  Sleuth project resource definition
  Item syntax:
  {
    name = "Sleuth"
    description = "Sleuth application"
  }
  NOTE: some of the arguments are optional
  EOF

  type = object({
    name                         = string
    build_provider               = optional(string, null)
    change_failure_rate_boundary = optional(string, null)
    failure_sensitivity          = optional(number, null)
    impact_sensitivity           = optional(string, null)
    issue_tracker_provider_type  = optional(string, null)
  })
  default = {
    name = ""
  }
}

variable "environments" {
  description = <<EOF
  List of environments.
  Item syntax:
  [
    {
      name = "production"
      description = "Production Environment"
      color = "#00d176"
    },
    {
      name        = "development"
      description = "Dev Environment"
      color       = "#00d176"
    }
  ]
  NOTE: some of the arguments are optional
  EOF

  type = list(object({
    name        = string
    description = optional(string, null)
    color       = optional(string, null)
  }))
  default = []
}

variable "repositories" {
  description = <<EOF
  List of repository objects.
  Item syntax:
[
    {
      name     = "terraform-provider-sleuth"
      owner    = "sleuth-io
      provider = "GITHUB"
      url      = "https://github.com/sleuth-io/terraform-provider-sleuth"
      environment_mappings = {
        prod = {
          environment_name = "production"
          branch           = "master"
        },
        dev = {
          environment_name = "development"
          branch           = "dev"
        }
      }
    },
    {
      name     = "sleuth-client"
      owner    = "sleuth-io"
      provider = "GITHUB"
      url      = "https://github.com/sleuth-io/sleuth-client"
      environment_mappings = {
        prod = {
          environment_name = "production"
          branch           = "master"
        },
        dev = {
          environment_name = "dev"
          branch           = "dev"
        }
      }
    },
    {
      name     = "sleuth-gitbook-docs"
      owner    = "sleuth-io"
      provider = "GITHUB"
      url      = "https://github.com/sleuth-io/sleuth-gitbook-docs"
      environment_mappings = {
        prod = {
          environment_name = "production"
          branch           = "master"
        }
      }
    }
  ]
  NOTE: some of the arguments are optional
  EOF
  type = list(object({
    name                 = string
    owner                = string
    provider             = string
    url                  = string
    deploy_tracking_type = optional(string, "manual")
    collect_impact       = optional(bool, true)
    environment_mappings = map(object({
      branch           = string
      environment_name = string
    }))
  }))
  default = [
  ]
}


