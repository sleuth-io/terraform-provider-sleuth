terraform {
  required_providers {
    sleuth = {
      source  = "sleuth-io/sleuth"
      version = "~> 0.2.0"
    }
  }
}

provider "sleuth" {
  api_key = "this-api-key-is-your-sleuth-organization-api-key"
}