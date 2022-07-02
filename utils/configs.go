package utils

// Config Структура конфигурации TCP чата
type Config struct {
	BindAddr       string `json:"Bind_addr"`
	MaxConnections int `json:"max_connections"`
}

// NewConfig конструктор для нашей структуры 
func NewConfig() *Config {
	return&Config{}
}