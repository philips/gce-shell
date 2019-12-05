package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v0.beta"
)

var (
	project  = flag.String("p", "ifup-shell", "project id")
	zone     = flag.String("z", "us-west1-b", "zone id")
	instance = flag.String("i", "instance-1", "instance id")

	secret = flag.String("s", "pho3KieGh6", "secret")
)

func main() {
	flag.Usage = usage
	flag.Parse()
	ctx := context.Background()

	args := flag.Args()

	if len(args) == 0 {
		usage()
	}

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	var output []byte

	switch args[0] {
	default:
		usage()

	case "start":
		resp, err := computeService.Instances.Start(*project, *zone, *instance).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}
		output, _ = json.MarshalIndent(resp, "", "    ")

	case "stop":
		resp, err := computeService.Instances.Stop(*project, *zone, *instance).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}
		output, _ = json.MarshalIndent(resp, "", "    ")

	case "server":
		server()

	case "deploy":
		deploy()

	}

	fmt.Printf("%s\n", output)
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage: gce-shell -z <zone> -p <project> -i <instance>

Commands:
	start
	stop
	server
	deploy
`)

	os.Exit(2)
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func deploy() {
	fmt.Printf(`gcloud run deploy gce-shell \
		--image quay.io/philips/gce-shell
		--set-env-vars 
`)

	fmt.Printf(`gcloud projects add-iam-policy-binding PROJECT-ID \
--member=serviceAccount:service-PROJECT-NUMBER@gcp-sa-pubsub.iam.gserviceaccount.com \
--role=roles/iam.serviceAccountTokenCreator`)

	fmt.Printf(`gcloud run services add-iam-policy-binding pubsub-tutorial \
--member=serviceAccount:cloud-run-pubsub-invoker@PROJECT-ID.iam.gserviceaccount.com \
--role=roles/run.invoker`)
}
