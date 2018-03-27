package main

import (
	"io/ioutil"
	"net"
	"flag"
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"golang.org/x/crypto/ssh"
	"strings"
)

const VERSION = "0.1"

type Config struct {
	LogLevel        string
	ServerAddr      string
	Key             string
	MikrotikAddr    string
	MikrotikUser    string
	MikrotikVersion string
	SshClientConfig *ssh.ClientConfig
	SshSigner       ssh.Signer
}

func mikrotikPing(conf Config) (bool, error) {
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return false, err
	}
	defer client.Close()
	return true, nil
}

func mikrotikRunScript(conf Config, command string) (string, error) {
	res := ""
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", conf.MikrotikAddr, conf.SshClientConfig)
	if err != nil {
		return res, err
	}
	defer client.Close()
	// Create a session
	session, err := client.NewSession()
	if err != nil {
		return res, err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)
	return stdoutBuf.String(), nil
}

func sendJsonError(app *iris.Application, ctx iris.Context, status int, message string) {
	app.Logger().Error(message)
	ctx.StatusCode(status)
	ctx.JSON(iris.Map{"status": status, "error": message})
}

func main() {
	isDev := flag.Bool("dev", false, "development mode")
	flag.Parse()

	app := iris.New()
	app.Configure(iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:     !*isDev,
		DisableVersionChecker: !*isDev,
		Charset:               "UTF-8",
	}))
	app.Use(recover.New())
	app.Use(logger.New())

	//read configuration
	conf := Config{
		LogLevel:     "info",
		ServerAddr:   ":8080",
		Key:          "mikrotik.dsa",
		MikrotikUser: "switcherUser",
		MikrotikAddr: "192.168.1.202:22",
	}
	if _, err := toml.DecodeFile("./switcher.conf", &conf); err != nil {
		app.Logger().Warn("Config problems: " + err.Error())
	}
	// set LogLevel from config
	app.Logger().SetLevel(conf.LogLevel)
	app.Logger().Info("API version: ", VERSION)

	// init ssh key
	buffer, err := ioutil.ReadFile(conf.Key)
	if err != nil {
		app.Logger().Fatal("Read ssh key problems: " + err.Error())
	}

	conf.SshSigner, err = ssh.ParsePrivateKey(buffer)
	if err != nil {
		app.Logger().Fatal("Parse ssh key problems: " + err.Error())
	}

	conf.SshClientConfig = &ssh.ClientConfig{
		User: conf.MikrotikUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(conf.SshSigner),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	ver, err := mikrotikRunScript(conf, "/system script run get_version")
	if err != nil {
		app.Logger().Fatal("run problems: " + err.Error())
	}
	conf.MikrotikVersion = strings.TrimSpace(ver)
	app.Logger().Info(conf.MikrotikVersion)

	app.Favicon("./static/favicon.ico")

	app.StaticWeb("/static", "./static")

	// not parse index.html - conflict vith VueJS
	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("index.html", true)
	})

	// api
	v1 := app.Party("/api/v1")
	{
		v1.Get("/version", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"version": VERSION})
		})

		v1.Get("/mikrotik", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"version": conf.MikrotikVersion})
		})

		v1.Get("/provider", func(ctx iris.Context) {
			provider, err := mikrotikRunScript(conf,
				"/system script run get_provider")
			if err != nil {
				sendJsonError(app, ctx, iris.StatusInternalServerError,
					"run get_provider err: "+err.Error())
				return
			}
			ctx.JSON(iris.Map{"provider": strings.TrimSpace(provider)})
		})

		v1.Post("/switch", func(ctx iris.Context) {
			provider, err := mikrotikRunScript(conf,
				"/system script run switch_provider")
			if err != nil {
				sendJsonError(app, ctx, iris.StatusInternalServerError,
					"run switch_provider err: "+err.Error())
				return
			}
			ctx.JSON(iris.Map{"provider": strings.TrimSpace(provider)})
		})
	}
	app.Run(iris.Addr(conf.ServerAddr))
}
