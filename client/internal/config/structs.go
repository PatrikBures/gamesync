package config

type Config struct {
	Server	ServerConfig	`yaml:"server"`
	Pools	[]PoolConfig	`yaml:"pools"`
	Games	[]GameConfig	`yaml:"games"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
}

type PoolConfig struct {
	ID		string `yaml:"id"`
	Path	string `yaml:"path"`
}

type GameConfig struct {
	ID			string `yaml:"id"`
	PoolID		string `yaml:"pool_id"`
	SavePath	string `yaml:"save_path"`
}
