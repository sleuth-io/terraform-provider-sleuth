terraform {
  required_version = ">= 1.4.0"

  required_providers {
    sleuth = {
      source  = "sleuth-io/sleuth"
      version = "~> 0.3.8"
    }
  }
}

provider "sleuth" {
  api_key = "enter_org_api_key"
}
