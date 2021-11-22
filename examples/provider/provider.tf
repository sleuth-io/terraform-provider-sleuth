terraform {
  required_providers {
    sleuth = {
      source  = "sleuth-io/sleuth"
      version = "~> x.x.x"
    }
  }
}

provider "sleuth" {
  api_key = "this-api-key-is-your-sleuth-organization-api-key"
}