package client

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

type Vault struct {
	Client *api.Client,
	IsKVV2 bool,
	Path string
}

func New(address string, token string) *Vault, error {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)

	if err != nil {
		return nil, err
	}

	origin.SetToken(originToken)

	health, err := client.Sys().Health()

	if err != nil {
		return nil, err
	}

	if !health.Initialized || health.Sealed {
		return nil, error {
			func Error() {
				return "Vault unable to handle requests."
			}
		}
	}

	return &Vault{
		Client: client
	}, nil
}

func (v *Vault) KVIsV2(isKVV2 bool) {
	v.IsKVV2 = isKVV2
}

func (v *Vault) SetPath(path string) {
	if v.IsKVV2 {
		path += "data/"
	}

	v.Path = path
}

// Read accepts a vault path to read the data out of. It will return a map
// of base64 encoded values.
func (v *Vault) Read(path string) map[string]interface{} {
	out := make(map[string]interface{})

	s, err := v.Client.Logical().Read(path)
	if err != nil {
		fmt.Printf("Error reading secrets, err=%v", err)
		return nil
	}

	// Encode all k,v pairs
	if s == nil || s.Data == nil {
		fmt.Printf("No data to read at path, %s\n", path)
		return out
	}
	for k, v := range s.Data {
		switch t := v.(type) {
		case string:
			out[k] = base64.StdEncoding.EncodeToString([]byte(t))
		case map[string]interface{}:
			if k == "data" {
				for x, y := range t {
					if z, ok := y.(string); ok {
						out[x] = base64.StdEncoding.EncodeToString([]byte(z))
					}
				}
			}
		default:
			fmt.Printf("error reading value at %s, key=%s, type=%T\n", path, k, v)
		}
	}

	return out
}

// Write takes in a vault path and base64 encoded data to be written at that path.
func (v *Vault) Write(path string, data map[string]string) error {
	body := make(map[string]interface{})

	// Decode the base64 values
	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		body[k] = string(b)
	}

	var err error

	if destinationIsKvV2 {
		d := make(map[string]interface{})
		d["data"] = body
		_, err = v.Client.Logical().Write(path, d)
	} else {
		_, err = v.Client.Logical().Write(path, body)
	}

	return err
}
