package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/igorbrites/vault-migrator/client"
	"github.com/igorbrites/vault-migrator/migrator"
)

var (
	m migrator.Migrator
	origin client.Vault
	destination client.Vault

	originAddr = flag.String("origin-addr", "", "The Vault Address of the backend to be migrated")
	originIsKvV2 = flag.Bool("origin-is-kvv2", false, "Whether the origin backend is in KV-V2 format")
	originPath = flag.String("origin-path", "secret/", "The path to be migrated (no need to pass \"data/\" when using KV-V2)")

	destinationAddr = flag.String("destination-addr", "", "The Vault Address of the backend that will receive the migration")
	destinationIsKvV2 = flag.Bool("destination-is-kvv2", false, "Whether the destination backend is in KV-V2 format")

	originToken = os.Getenv("ORIGIN_VAULT_TOKEN")
	destinationToken = os.Getenv("DESTINATION_VAULT_TOKEN")
)

func main() {
	if err := validate(); err != nil {
		fmt.Println(err)
		return
	}

	startup()
}

func startup() {
	m = migrator.Migrator{
		Origin: *buildClient(*originAddr, originToken, *originPath, *originIsKvV2),
		Destination: *buildClient(*destinationAddr, destinationToken, *originPath, *destinationIsKvV2),
	}

	m.Start()
}

func validate() error {
	flag.Parse()

	if *originAddr == "" {
		return errors.New("Missing -origin-addr")
	}

	if *originPath == "" {
		return errors.New("Missing -origin-path")
	}

	if *destinationAddr == "" {
		return errors.New("Missing -destination-addr")
	}

	if originToken == "" {
		return errors.New("Missing ORIGIN_VAULT_TOKEN environment variable")
	}

	if destinationToken == "" {
		return errors.New("Missing DESTINATION_VAULT_TOKEN environment variable")
	}

	return nil
}

func buildClient(addr string, token string, path string, isKVV2 bool) *client.Vault {
	c, err := client.New(addr, token)

	if err != nil {
		fmt.Printf("Unable to generate client, err=%v\n", err)
		return nil
	}

	c.KVIsV2(isKVV2)
	c.SetPath(path)

	return c
}
