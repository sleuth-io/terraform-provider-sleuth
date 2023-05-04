projects = {
  sleuth = {
    project = {
      name = "sleuth"
    }
    environments = {
      prod = {
        name        = "Production"
        description = "Production Environment"
        color       = "#337ab7"
      }
      staging = {
        name        = "Staging"
        description = "Staging Environment"
        color       = "#e8b017"
      }
    }
    repositories = {
      sleuth = {
        name     = "application"
        owner    = "sleuth-io"
        provider = "GITHUB"
        url      = "https://github.com/sleuth-io/sleuth"
        environment_mappings = {
          prod = {
            environment_name = "Production"
            branch           = "master"
          }
          staging = {
            environment_name = "Staging"
            branch           = "master"
          }
        }
      }
      documentation = {
        name     = "documentation"
        owner    = "sleuth-io"
        provider = "GITHUB"
        url      = "https://github.com/sleuth-io/sleuth-gitbook-docs"
        environment_mappings = {
          prod = {
            environment_name = "Production"
            branch           = "master"
          }
        }
      }
      terraform = {
        name     = "terraform"
        owner    = "sleuth-io"
        provider = "GITHUB"
        url      = "https://github.com/sleuth-io/sleuth-terraform"
        environment_mappings = {
          prod = {
            environment_name = "Production"
            branch           = "master"
          }
        }
      }
    }
  },
  sleuth_clients = {
    project = {
      name = "sleuth clients"
    }
    environments = {
      prod = {
        name        = "Production"
        description = "Production Environment"
        color       = "#059041"
      }
    }
    repositories = {
      sleuth-client = {
        name     = "sleuth-client"
        owner    = "sleuth-io"
        provider = "GITHUB"
        url      = "https://github.com/sleuth-io/sleuth-client"
        environment_mappings = {
          prod = {
            environment_name = "Production"
            branch           = "master"
          }
        }
      }
      terraform-provider-sleuth = {
        name     = "documentation"
        owner    = "sleuth-io"
        provider = "GITHUB"
        url      = "https://github.com/sleuth-io/terraform-provider-sleuth"
        environment_mappings = {
          prod = {
            environment_name = "Production"
            branch           = "master"
          }
        }
      }
    }
  }
}
