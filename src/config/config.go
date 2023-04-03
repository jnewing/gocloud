package config

type Configs struct {
	Cloudflare Cloudflare
	Zones      []Zone
}

type Cloudflare struct {
	API   string `yaml:"api"`
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}

type Zone struct {
	ID     string   `yaml:"id"`
	Update []string `yaml:"update"`
	IP     string   `yaml:"ip"`
}
