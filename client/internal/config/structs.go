package config

var Current Config

type Config struct {
	Server	ServerConfig	`yaml:"server"`
	Games	[]GameConfig	`yaml:"games"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
}

type GameConfig struct {
	ID			string `yaml:"id"`
	SavePath	string `yaml:"save_path"`
}
