package config

type Config struct {
	Engine      string      `yaml:"engine" validate:"required"`
	Debug       bool        `yaml:"debug"`
	GrpcConfig  GrpcConfig  `yaml:"grpc"`
	MysqlConfig MysqlConfig `yaml:"mysql"`
	StoreType   string      `yaml:"store_type" validate:"required"`
	LogConfig   LogConfig   `yaml:"log"`
}

type LogConfig struct {
	Path         string `yaml:"path" validate:"required"`
	Level        string `yaml:"level" validate:"required"`
	MaxAge       int64  `yaml:"max_age" validate:"required"`
	RotationTime int64  `yaml:"rotation_time" validate:"required"`
}

type GrpcConfig struct {
	Port int `yaml:"server_port"`
}

type MysqlConfig struct {
	Ip       string `yaml:"ip" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Database string `yaml:"database" validate:"required"`
	Debug    bool   `yaml:"debug" `
}
