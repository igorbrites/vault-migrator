package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/igorbrites/vault-migrator/client"
)

var origin *client.Vault
var destination *client.Vault
var originIsKvV2 bool
var destinationIsKvV2 bool

func main() {
	Startup()
	RecursiveMigration(origin.Path, destination.Path)
}

func Startup() {
	originAddrFlag := flag.String("origin-addr", "", "The Vault Address of the backend to be migrated")
	originIsKvV2Flag := flag.Bool("origin-is-kvv2", false, "Whether the origin backend is in KV-V2 format")
	originPathFlag := flag.String("origin-path", "secret/", "The path to be migrated (no need to pass \"data/\" when using KV-V2)")

	destinationAddrFlag := flag.String("destination-addr", "", "The Vault Address of the backend that will receive the migration")
	destinationIsKvV2Flag := flag.Bool("destination-is-kvv2", false, "Whether the destination backend is in KV-V2 format")

	originToken := os.Getenv("ORIGIN_VAULT_TOKEN")
	destinationToken := os.Getenv("DESTINATION_VAULT_TOKEN")

	flag.Parse()
	origin = BuildClient(*originAddrFlag, os.Getenv("ORIGIN_VAULT_TOKEN"), *originPathFlag, *originIsKvV2Flag)
	destination = BuildClient(*destinationAddrFlag, os.Getenv("DESTINATION_VAULT_TOKEN"), *destinationPathFlag, *destinationIsKvV2Flag)
}

func BuildClient(addr string, token string, path string, isKVV2 bool) *client.Vault {
	c, err := client.New(addr, token)

	if err != nil {
		fmt.Printf("Unable to generate client, err=%v\n", err)
		return nil
	}

	c.KVIsV2(isKVV2)
	c.SetPath(isKVV2)

	return c
}

func RecursiveMigration(originPath string, destinationPath string) {
	if originPath[len(originPath)-1:] != "/" {
		Migrate(originPath, destinationPath)
		return
	}

	fmt.Printf("Listing %q\n", originPath)
	s, err := origin.Logical().List(originPath)

	if s == nil {
		fmt.Printf("Unable to read path %q, err=response was empty\n", originPath)
		return
	}
	if err != nil {
		fmt.Printf("Unable to read path %q, err=%v\n", originPath, err)
		return
	}

	r, ok := s.Data["keys"].([]interface{})
	if !ok {
		fmt.Println("Error listing path")
	}

	for i := range r {
		if path, ok := r[i].(string); ok {
			RecursiveMigration(originPath+path, destinationPath+path)
		}
	}
}

func Migrate(originPath string, destinationPath string) {
	fmt.Printf("Migrating key %q to %q\n", originPath, destinationPath)

	from := Read(originPath)
	fmt.Println(from)
	// to := Write(destinationPath, from)
}
