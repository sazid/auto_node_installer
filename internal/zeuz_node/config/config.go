package config

type Config struct {
	// lastUpdated *string `json:"last_updated,omitempty"`
}

type Paths struct {
	HomeDir                 string
	WorkingDir              string
	ZeuzNodeDir             string
	ZeuzLogDir              string
	ZeuzPayloadDir          string
	DefaultPythonInstallDir string

	PythonPath string
}
