# Terraform Provider for Sleuth

![](https://github.com/sleuth-io/terraform-provider-sleuth/actions/workflows/test.yml/badge.svg)
[![Latest release](https://img.shields.io/github/v/release/sleuth-io/terraform-provider-sleuth)](https://github.com/sleuth-io/terraform-provider-sleuth/releases)


This repository is a Terraform provider for [Sleuth](https://sleuth.io), allowing a team to manage Sleuth configuration via Terraform instead of having to click around in the UI.

* [Documentation](https://registry.terraform.io/providers/sleuth-io/sleuth/latest/docs)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Known Limitations
- No support for Azure as repository provider
- No support for New Relic as impact source provider
- No support for custom impact sources
- Limited support for multiple integrations of same type

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Quick Starts

- [Using the provider](docs/index.md)
- [Provider development](docs/contributing)
