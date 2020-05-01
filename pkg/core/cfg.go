package core

import (
	"encoding/json"
	"os"
)

// Config represents the json config file
type Config struct {
	ConnectionStrings struct {
		MySQL string
	}
	Bot struct {
		Commands struct {
			Prefix string
		}
		OAuthToken string
		Nickname   string
	}
}

// LoadConfig loads the config file at given location
func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(cfg)
	return cfg, err
}
