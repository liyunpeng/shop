package config

type Config struct {
	Server ServerConf   `yaml:"server"`
	Mysql  MysqlConf `yaml:"mysql"`
	Redis  RedisConf `yaml:"redis"`
}

type ServerConf struct {
	Port int `yaml:"port"`
	List []string `yaml:"list,flow"`
}

type RedisConf struct {
	Enable bool `yaml:"enable"`
}

type MysqlConf struct {
	MaxIdle int `yaml:"maxidle"`
	MaxOpen int `yaml:"maxopen"`
	User string `yaml:"user"`
	Host string `yaml:"host"`
	Password string `yaml:"password"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

