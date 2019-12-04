package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v0.beta"
)

var (
	project  = flag.String("p", "ifup-shell", "project id")
	zone     = flag.String("z", "us-west1-b", "zone id")
	instance = flag.String("i", "instance-1", "instance id")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := computeService.Instances.Start(*project, *zone, *instance).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", resp)
}
