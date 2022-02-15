package config

import "testing"

func TestConvertVersionStringToInt(t *testing.T) {
	cases := []struct {
		versionStr string
		want       int
	}{
		{"", 0},
		{"invalid version string", 0},
		{"v0.0.0", 0},
		{"0.0.0", 0},
		{"v0.0.1", 1},
		{"v1.0.0", 0},
		{"1.2.3", 23},
		{"v1.2.3", 23},
	}

	for _, c := range cases {
		_, got := ConvertVersionToInt(c.versionStr)
		if got != c.want {
			t.Fatalf("got: %d, want: %d", got, c.want)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		current string
		latest  string
		want    int
	}{
		{"", "", 0},
		{"invalid version string", "another invalid", 0},
		{"1.0.0", "", 0},
		{"1.0.0", "invalid", 0},
		{"v0.0.0", "v0.0.0", 0},
		{"0.0.1", "v0.0.2", 1},
		{"0.0.2", "0.0.1", -1},
		{"1.0.2", "2.0.0", 0},
		{"1.0.2", "2.0.3", 0},
		{"1.0.2", "1.3.5", 33},
		{"1.2", "1", 0},
	}

	for _, c := range cases {
		conf := &Config{
			CurrentVersion: c.current,
		}
		got := conf.CompareVersion(c.latest)
		if got != c.want {
			t.Fatalf("got: %d, want: %d for the case: %+v", got, c.want, c)
		}
	}
}
