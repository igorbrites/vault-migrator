package vault

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type Vault struct {
	Client *api.Client
	IsKVV2 bool
	Path   string
}

func New(address string, token string) (*Vault, error) {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)

	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	health, err := client.Sys().Health()

	if err != nil {
		return nil, err
	}

	if !health.Initialized || health.Sealed {
		return nil, errors.New("Vault unable to handle requests")
	}

	return &Vault{
		Client: client,
	}, nil
}

func (v *Vault) KVIsV2(isKVV2 bool) {
	v.IsKVV2 = isKVV2
}

func (v *Vault) SetPath(path string) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}

	v.Path = path
}

func (v *Vault) Read(path string) (map[string]string, error) {
	out := make(map[string]string)

	s, err := v.Client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if s == nil || s.Data == nil {
		return nil, fmt.Errorf("No data to read at path %q", path)
	}

	for k, v := range s.Data {
		switch t := v.(type) {
		case string:
			out[k] = base64.StdEncoding.EncodeToString([]byte(t))
		case json.Number:
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
			return nil, fmt.Errorf("Error reading value at %q, key=%q, type=%T", path, k, v)
		}
	}

	return out, nil
}

func (v *Vault) Write(path string, data map[string]string) error {
	body := make(map[string]interface{})

	for k, v := range data {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
		body[k] = string(b)
	}

	var err error

	if v.IsKVV2 {
		d := make(map[string]interface{})
		d["data"] = body
		_, err = v.Client.Logical().Write(path, d)
	} else {
		_, err = v.Client.Logical().Write(path, body)
	}

	return err
}
