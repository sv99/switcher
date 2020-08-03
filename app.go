package switcher

import (
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
	"github.com/gofiber/recover"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

const (
	VERSION  = "0.3"
	ConfFile = "switcher.conf"
)

type App struct {
	App    *fiber.App
	Logger *zerolog.Logger
	Config *Config
	WorkDir   string
}

func NewApp(workDir string, zeroLogger *zerolog.Logger) *App {
	conf := DefaultConfig()
	err := conf.TOML(filepath.Join(workDir, ConfFile))
	if err != nil {
		log.Panicf("Config load problems: %v", err)
	}

	app := fiber.New()

	srv := App{
		Config:  &conf,
		App:     app,
		WorkDir: workDir,
		Logger: zeroLogger,
	}

	// set global log level
	if conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// set LogLevel from config
	//app.Logger().SetLevel(conf.LogLevel)


	// dump config
	srv.Logger.Debug().Msgf("Config: %s", workDir)
	srv.Logger.Debug().Msgf("ServerAddr: %s", srv.Config.ServerAddr)
	srv.Logger.Debug().Msgf("API version: ", VERSION)

	app.Use(NewLoggerMiddleware(zeroLogger))

	app.Use(recover.New(recover.Config{
		Handler: func(c *fiber.Ctx, err error) {
			c.Status(500)
			_ = c.JSON(fiber.Map{"Message": err.Error()})
		},
	}))
	app.Use(helmet.New())

	// init ssh key
	buffer, err := ioutil.ReadFile(conf.Key)
	if err != nil {
		srv.Logger.Fatal().Msgf("Read ssh key problems: %v", err)
	}

	conf.SshSigner, err = ssh.ParsePrivateKey(buffer)
	if err != nil {
		srv.Logger.Fatal().Msgf("Parse ssh key problems: %v", err)
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
		srv.Logger.Fatal().Msgf("run problems: %v", err)
	}
	conf.MikrotikVersion = strings.TrimSpace(ver)
	srv.Logger.Info().Msg(conf.MikrotikVersion)

	// handlers
	app.Get("/", func(ctx *fiber.Ctx) {
		//_ = ctx.SendFile("./index.html")
		ctx.Type("html", "utf-8")
		ctx.Send(_indexHtml)
	})
	app.Use(middleware.Favicon("./static/favicon.ico"))
	app.Static("/static", "./static")

	// API
	v1 := app.Group("/api/v1")
	{
		v1.Get("/version", func(ctx *fiber.Ctx) {
			_ = ctx.JSON(fiber.Map{"version": VERSION})
		})

		v1.Get("/mikrotik", func(ctx *fiber.Ctx) {
			_ = ctx.JSON(fiber.Map{"version": conf.MikrotikVersion})
		})

		v1.Get("/provider", func(ctx *fiber.Ctx) {
			provider, err := mikrotikRunScript(conf,
				"/system script run get_provider")
			if err != nil {
				srv.Error(ctx, fiber.StatusInternalServerError,
					"run get_provider err: " + err.Error())
				return
			}
			_ = ctx.JSON(fiber.Map{"provider": strings.TrimSpace(provider)})
		})

		v1.Post("/switch", func(ctx *fiber.Ctx) {
			provider, err := mikrotikRunScript(conf,
				"/system script run switch_provider")
			if err != nil {
				srv.Error(ctx, fiber.StatusInternalServerError,
					"run switch_provider err: " + err.Error())
				return
			}
			_ = ctx.JSON(fiber.Map{"provider": strings.TrimSpace(provider)})
		})

		// Dune test
		v1.Get("/dune/names", func(ctx *fiber.Ctx) {
			_ = ctx.JSON(fiber.Map{"names": conf.Dunes.Name})
		})

		v1.Get("/dune/:id/status", func(ctx *fiber.Ctx) {
			id, _ := strconv.Atoi(ctx.Params("id"))
			status := "unknown"
			if id >= 0 && id < len(conf.Dunes.IP) {
				ip := conf.Dunes.IP[id]
				status, err = getDuneStatus(ip)
				if err != nil {
					srv.Logger.Error().Msgf("getDuneStatus error: %v", err)
				}
			}
			_ = ctx.JSON(fiber.Map{"status": status})
		})

		v1.Get("/dune/:id/on", func(ctx *fiber.Ctx) {
			id, _ := strconv.Atoi(ctx.Params("id"))
			status := false
			if id >= 0 && id < len(conf.Dunes.IP) {
				ip := conf.Dunes.IP[id]
				status, _ = getDuneOn(ip)
			}
			_ = ctx.JSON(fiber.Map{"status": status})
		})

		v1.Get("/dune/:id/off", func(ctx *fiber.Ctx) {
			id, _ := strconv.Atoi(ctx.Params("id"))
			status := false
			if id >= 0 && id < len(conf.Dunes.IP) {
				ip := conf.Dunes.IP[id]
				status, _ = getDuneOff(ip)
			}
			_ = ctx.JSON(fiber.Map{"status": status})
		})
	}
	app.Use(srv.NotFound)
	return &srv
}

func (srv *App) Serve() {
	// Create tls certificate
	srv.Logger.Info().Msg("Start server")
	err := srv.App.Listen(srv.Config.ServerAddr)
	if err != nil {
		srv.Logger.Fatal().Msgf("Listen error: %v", err)
	}
	srv.Logger.Info().Msg("Server stopped")
}

func (srv *App) Error(c *fiber.Ctx, status int, message string) {
	srv.Logger.Info().Msg(message)
	c.SendStatus(status)
	_ = c.JSON(fiber.Map{"status": status, "error": message})
}

func (srv *App) NotFound(c *fiber.Ctx) {
	srv.Error(
		c,
		fiber.StatusNotFound,
		"Sorry, but the page you were looking for could not be found.",
	)
}
