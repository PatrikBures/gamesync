package config

type Current struct {
	Config Config
	ConfigPath string
	Verbose bool
}

type Config struct {
	Server	ServerConfig	`yaml:"server"`
	Games	[]GameConfig	`yaml:"games"`
}

type ServerConfig struct {
	Host			string `yaml:"host"`
	User			string `yaml:"user"`
	Port			string `yaml:"port"`
	IdentityFile	string `yaml:"identity_file"`
}

type GameConfig struct {
	ID			string `yaml:"id"`
	SavePath	string `yaml:"save_path"`
}
