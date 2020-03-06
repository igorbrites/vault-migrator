package migrator

import (
	"fmt"
	"github.com/igorbrites/vault-migrator/vault"
)

type Migrator struct {
	Origin      vault.Vault
	Destination vault.Vault
}

func (m Migrator) Start() {
	m.recursiveMigration(m.Origin.Path, m.Destination.Path)
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
	}

	for i := range r {
		if path, ok := r[i].(string); ok {
			m.recursiveMigration(originPath+path, destinationPath+path)
		}
	}
}

func (m Migrator) copyKey(originPath string, destinationPath string) {
	fmt.Printf("Copying key %q to %q\n", originPath, destinationPath)

	from, err := m.Origin.Read(originPath)

	if err != nil {
		fmt.Printf("Error reading key %q, err=%v", originPath, err)
		return
	}

	err = m.Destination.Write(destinationPath, from)

	if err != nil {
		fmt.Printf("Error writing key %q, err=%v", destinationPath, err)
		return
	}

	fmt.Printf("Key %q copied to %q successfully\n", originPath, destinationPath)
}
