# Vault Migrator

Simply migrate a [Vault](https://www.vaultproject.io/) path from a backend to another!

## Command line

```
Usage: vault-migrator [args]

First you must define the environment variables bellow:
  ORIGIN_VAULT_TOKEN
        The token with permittion to read the path to be migrated
  DESTINATION_VAULT_TOKEN
        The token with permittion to write in the migrated path

Args:
  -destination-addr string
        The Vault Address of the backend that will receive the migration
  -destination-is-kvv2
        Whether the destination backend is in KV-V2 format
  -origin-addr string
        The Vault Address of the backend to be migrated
  -origin-is-kvv2
        Whether the origin backend is in KV-V2 format
  -origin-path string
        The path to be migrated (no need to pass "data/" when using KV-V2) (default "secret/")
```

## Example

To migrate an generic old backend to a new one using KV-V2:

```
$ vault-migrator -origin-addr=https://my-old-vault.example.com -origin-path=secret/ -destination-addr=https://my-new-vault.example.com -destination-is-kvv2
Starting migration of secret/

Listing "secret/"
Listing "secret/my-secret-path/"
Listing "secret/my-secret-path/the-subpath/"
Copying key "secret/my-secret-path/the-subpath/123" to "secret/data/my-secret-path/the-subpath/123"
Key "secret/my-secret-path/the-subpath/123" copied to "secret/data/my-secret-path/the-subpath/123" successfully
Copying key "secret/my-secret-path/the-subpath/1234" to "secret/data/my-secret-path/the-subpath/1234"
Key "secret/my-secret-path/the-subpath/1234" copied to "secret/data/my-secret-path/the-subpath/1234" successfully
Copying key "secret/my-secret-path/the-subpath/12345" to "secret/data/my-secret-path/the-subpath/12345"
Key "secret/my-secret-path/the-subpath/12345" copied to "secret/data/my-secret-path/the-subpath/12345" successfully
Copying key "secret/my-secret-path/the-subpath/14265" to "secret/data/my-secret-path/the-subpath/14265"
Key "secret/my-secret-path/the-subpath/14265" copied to "secret/data/my-secret-path/the-subpath/14265" successfully
Listing "secret/my-secret-path/another-subpath/"
Copying key "secret/my-secret-path/another-subpath/123" to "secret/data/my-secret-path/another-subpath/123"
Key "secret/my-secret-path/another-subpath/123" copied to "secret/data/my-secret-path/another-subpath/123" successfully
Copying key "secret/my-secret-path/another-subpath/1234" to "secret/data/my-secret-path/another-subpath/1234"
```
