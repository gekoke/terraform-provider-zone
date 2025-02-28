package main

import (
	"context"
	"flag"
	"log"

	"github.com/gekoke/terraform-provider-zone/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version string = "0.0.0-alpha.0"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with debugging support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/gekoke/zone",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
