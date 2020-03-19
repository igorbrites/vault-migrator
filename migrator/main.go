package migrator

import (
	"fmt"
	"github.com/igorbrites/vault-migrator/vault"
)

type Migrator struct {
	Origin      vault.Vault
	Destination vault.Vault
	Overwrite   bool
}

func (m Migrator) Start() {
	originPath := m.Origin.Path
	destinationPath := m.Destination.Path

	if m.Origin.IsKVV2 {
		originPath += "data/"
	}

	if m.Destination.IsKVV2 {
		destinationPath += "data/"
	}

	fmt.Printf("Starting migration of %q\n\n", originPath)
	m.recursiveMigration(originPath, destinationPath)

	if m.Origin.IsKVV2 && m.Destination.IsKVV2 {
		originPath = m.Origin.Path+"metadata/"
		destinationPath = m.Destination.Path+"metadata/"

		fmt.Printf("Starting migration of %q\n\n", originPath)
		m.recursiveMigration(originPath, destinationPath)
	}
}

func (m Migrator) recursiveMigration(originPath string, destinationPath string) {
	if originPath[len(originPath)-1:] != "/" {
		m.copyKey(originPath, destinationPath)
		return
	}

	fmt.Printf("Listing %q\n", originPath)
	s, err := m.Origin.Client.Logical().List(originPath)

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
		return
	}

	for i := range r {
		if path, ok := r[i].(string); ok {
			m.recursiveMigration(originPath+path, destinationPath+path)
		}
	}
}

func (m Migrator) copyKey(originPath string, destinationPath string) {
	fmt.Printf("Copying key %q to %q\n", originPath, destinationPath)

	to, err := m.Destination.Read(destinationPath)

	if (to != nil || len(to) > 0) && !m.Overwrite {
		fmt.Println("The destination path exists and overwrite is disabled. Skipping...")
		return
	}

	from, err := m.Origin.Read(originPath)

	if err != nil {
		fmt.Printf("Error reading key %q on origin, err=%v\n", originPath, err)
		return
	}

	err = m.Destination.Write(destinationPath, from)

	if err != nil {
		fmt.Printf("Error writing key %q, err=%v", destinationPath, err)
		return
	}

	fmt.Printf("Key %q copied to %q successfully\n", originPath, destinationPath)
}
