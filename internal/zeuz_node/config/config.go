package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	CurrentVersion string `json:"current_version"`
}

func NewConfig(f io.Reader) (Config, error) {
	c := Config{}

	err := json.NewDecoder(f).Decode(&c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func (c *Config) WriteToFile(filepath string) {
	// write to config file.
	f, err := os.OpenFile(
		filepath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		log.Fatalf("failed to open `%s` file for writing: %v", filepath, err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		log.Fatalf("failed to write config file: %v", err)
	}
}

var FirstRunVersion = "v0.0.0"

var DefaultConfig = fmt.Sprintf(`
{
	"current_version": "%s"
}
`, FirstRunVersion)
