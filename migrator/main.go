package migrator

import (
	"github.com/igorbrites/vault-migrator/vault"
)

type Migrator struct {
	Origin *vault.Vault,
	Destination *vault.Vault,
}
