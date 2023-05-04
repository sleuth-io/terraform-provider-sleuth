terraform {
  required_version = ">= 1.4.0"

  required_providers {
    sleuth = {
      source  = "sleuth-io/sleuth"
      version = "~> 0.3.8"
    }
  }
}

#provider "sleuth" {
#  baseurl = "http://dev.sleuth.io"
#  api_key = "7d61fd5d3b66917ca0064b330e3ec3b319a4d137"
#}
