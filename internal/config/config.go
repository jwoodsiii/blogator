package config

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() Config {
	// read blogatorconfig.json from home directory and marshal into Config struct
	return Config{}
}

func (c *Config) SetUser() {
	// write config struct to json after setting current_user_name
}

func getConfigFilePath() (string, error) {
	return "", nil
}

func write(cfg Config) error {
	return nil
}
