
terraform {	
	required_providers {	
		sleuth = {		
			version = "1.0-dev"
			source	= "sleuth.io/core/sleuth"	
		}	
	}
}
provider "sleuth" {
	baseurl = "http://dev.sleuth.io"
    api_key = "9208af1845ec3517ee80f8113097f57d9da1a6b2"
	org_slug = "myorg"
}

resource "sleuth_project" "myproject" {
	name = "My Project fixed 2"
}

