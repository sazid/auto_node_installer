package config

type Config struct {
	Dirs Dirs
}

type Dirs struct {
	HomeDir                 string
	ZeuzRootDir             string
	ZeuzNodeDir             string
	ZeuzLogDir              string
	ZeuzPayloadDir          string
	DefaultPythonInstallDir string
}
