package switcher

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh"
)

type Dunes struct {
	Name []string
	IP   []string
}

type Config struct {
	LogLevel        string
	ServerAddr      string
	Debug           bool
	Key             string
	MikrotikAddr    string
	MikrotikUser    string
	MikrotikVersion string
	SshClientConfig *ssh.ClientConfig
	SshSigner       ssh.Signer
	Dunes           Dunes
}

func DefaultConfig() Config {
	return Config{
		LogLevel:     "info",
		ServerAddr:   ":8080",
		Debug: 		  false,
		Key:          "mikrotik.dsa",
		MikrotikUser: "switcherUser",
		MikrotikAddr: "192.168.1.202:22",
		Dunes:        Dunes{},
	}
}

func (conf *Config) TOML(path string) error {
	_, err := toml.DecodeFile(path, &conf)
	return err
}
