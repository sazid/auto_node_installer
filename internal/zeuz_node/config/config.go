package config

import (
	"encoding/json"
	"io"
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

// convertVersionToInt converts a SEMVER version string into an integer that can
// be compared. Returns 0 in case the version string is invalid - does not
// contain any integers.
func convertVersionToInt(v string) int {
	res := 0
	for _, c := range v {
		if c >= '0' && c <= '9' {
			res = res*10 + int(c-'0')
		}
	}
	return res
}

func (c *Config) CompareVersion(otherVersion string) int {
	cur := convertVersionToInt(c.CurrentVersion)
	other := convertVersionToInt(otherVersion)
	if cur == 0 || other == 0 {
		return 0
	}
	return other - cur
}

var DefaultConfig = `
{
	"current_version": "0.0.0"
}
`
