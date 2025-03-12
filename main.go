package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gekoke/terraform-provider-zone/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version string = "@providerVersion@"

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	var debug bool
	var printVersion bool
	flag.BoolVar(&debug, "debug", false, "run the provider with debugging support")
	flag.BoolVar(&printVersion, "version", false, "print the version and exit")
	flag.Parse()

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/gekoke/zone",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
