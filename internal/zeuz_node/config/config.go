package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// ConvertVersionToInt converts a SEMVER version string into an integer that can
// be compared. Returns 0 in case the version string is invalid - does not
// contain any integers.
func ConvertVersionToInt(s string) (major int, minorPatch int) {
	firstDot := 0
	for i, c := range s {
		if c == '.' {
			firstDot = i
			break
		}
	}

	for i := 0; i < firstDot; i++ {
		if s[i] >= '0' && s[i] <= '9' {
			major = major*10 + int(s[i]-'0')
		}
	}

	for i := firstDot + 1; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			minorPatch = minorPatch*10 + int(s[i]-'0')
		}
	}
	return major, minorPatch
}

// CompareVersion returns the following values depending on different
// conditions:
//
// 0,  if both versions are equal or either/both of them are invalid
//
// greater than 0, if otherVersion is greater than c.CurrentVersion
//
// less than 0, if otherVersion is less than c.CurrentVersion
func (c *Config) CompareVersion(otherVersion string) (result int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("invalid version string: %v", r)
		}
	}()
	curMajor, cur := ConvertVersionToInt(c.CurrentVersion)
	otherMajor, other := ConvertVersionToInt(otherVersion)
	if curMajor < otherMajor {
		return 0
	}
	if cur == 0 || other == 0 {
		return 0
	}
	return other - cur
}

var FirstRunVersion = "0.0.0"

var DefaultConfig = fmt.Sprintf(`
{
	"current_version": "%s"
}
`, FirstRunVersion)
