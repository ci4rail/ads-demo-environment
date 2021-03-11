package main

import (
	"log"
	"os"

	"github.com/ci4rail/ads-demo-environment/eventhub2db/cmd"
	"github.com/ci4rail/ads-demo-environment/eventhub2db/internal/eventhub"
)

func main() {
	versionArgFound := false
	for _, v := range os.Args {
		if v == "version" || v == "help" || v == "--help" || v == "-h" {
			versionArgFound = true
		}
	}
	if !versionArgFound {
		_, ok := os.LookupEnv(eventhub.EnvEventHubConnectionsString)

		if !ok {
			log.Fatalf("Error: environment variable %s missing", eventhub.EnvEventHubConnectionsString)
		}
	}
	cmd.Execute()
}
