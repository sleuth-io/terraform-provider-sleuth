package main

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	frameworkProvider "github.com/sleuth-io/terraform-provider-sleuth/internal/sleuth"
)

var (
	// Version can be updated by goreleaser on release
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/sleuth-io/sleuth",
		Debug:   debug,
	}

	sleuthProvider := func() provider.Provider {
		return frameworkProvider.New(version)
	}

	err := providerserver.Serve(context.Background(), sleuthProvider, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
