package server

type Config struct {
	Port    int               `yaml:"port"`
	Host    string            `yaml:"host"`
	Users   map[string]string `yaml:"users"`
	Session struct {
		Secret     string `yaml:"secret"`
		CookieName string `yaml:"cookie_name"`
	} `yaml:"session"`
}
