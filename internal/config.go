package internal

import (
	"path/filepath"
	"reflect"
	"strconv"
	"switcher"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type Dunes struct {
	Name []string
	IP   []string
}

type Config struct {
	WorkDir         string
	LogLevel        string
	ServerAddr      string
	Debug           bool
	Key             string
	MikrotikAddr    string
	MikrotikUser    string
	SshClientConfig *ssh.ClientConfig
	SshSigner       ssh.Signer
	Dunes           Dunes
}

func DefaultConfig() *Config {
	return &Config{
		LogLevel:     "info",
		ServerAddr:   ":8080",
		Debug:        false,
		Key:          "mikrotik.dsa",
		MikrotikUser: "switcherUser",
		MikrotikAddr: "192.168.1.202:22",
		Dunes:        Dunes{},
	}
}

// LoadConfig initialize app config from toml file <package_name>.conf
func LoadConfig(workDir string) *Config {
	packageName := reflect.TypeOf(switcher.Empty{}).PkgPath()
	confPath := filepath.Join(workDir, packageName+".conf")
	conf := DefaultConfig()
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		log.Fatal().Msgf("Config load problems: %v", err)
	}
	conf.WorkDir = workDir
	return conf
}

func (conf *Config) Dump() {
	// TODO: dump config
}

func (conf *Config) GetDuneIP(strId string) string {
	var ip = ""
	id, _ := strconv.Atoi(strId)
	if id >= 0 && id < len(conf.Dunes.IP) {
		ip = conf.Dunes.IP[id]
	}
	return ip
}
