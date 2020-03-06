package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/igorbrites/vault-migrator/migrator"
	"github.com/igorbrites/vault-migrator/vault"
)

var (
	m migrator.Migrator

	originAddr   = flag.String("origin-addr", "", "The Vault Address of the backend to be migrated")
	originPath   = flag.String("origin-path", "secret/", "The path to be migrated (no need to pass \"data/\" when using KV-V2)")
	originIsKvV2 = flag.Bool("origin-is-kvv2", false, "Whether the origin backend is in KV-V2 format")

	destinationAddr   = flag.String("destination-addr", "", "The Vault Address of the backend that will receive the migration")
	destinationIsKvV2 = flag.Bool("destination-is-kvv2", false, "Whether the destination backend is in KV-V2 format")

	originToken      = os.Getenv("ORIGIN_VAULT_TOKEN")
	destinationToken = os.Getenv("DESTINATION_VAULT_TOKEN")
)

func main() {
	flag.Usage = func() {
		help := `Usage: %s [args]

First you must define the environment variables bellow:
  ORIGIN_VAULT_TOKEN
        The token with permittion to read the path to be migrated
  DESTINATION_VAULT_TOKEN
        The token with permittion to write in the migrated path

Args:
`

		fmt.Fprintf(flag.CommandLine.Output(), help, os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startup()
}

func startup() {
	m = migrator.Migrator{
		Origin:      *buildClient(*originAddr, originToken, *originPath, *originIsKvV2),
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

func buildClient(addr string, token string, path string, isKVV2 bool) *vault.Vault {
	c, err := vault.New(addr, token)

	if err != nil {
		fmt.Printf("Unable to generate client, err=%v\n", err)
		return nil
	}

	c.KVIsV2(isKVV2)
	c.SetPath(path)

	return c
}
