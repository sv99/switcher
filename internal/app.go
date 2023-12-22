package internal

import (
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

type App struct {
	App     *fiber.App
	Config  *Config
	WorkDir string
}

func NewApp(conf *Config) *App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		StrictRouting:         true,
		AppName:               "switcher",
	})

	srv := App{
		Config: conf,
		App:    app,
	}

	app.Use(NewLoggerMiddleware(&log.Logger))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
			log.Error().Msgf("panic: %v\n%s", e, debug.Stack())
		},
	}))
	app.Use(helmet.New())

	// init ssh key
	buffer, err := os.ReadFile(filepath.Join(conf.WorkDir, conf.Key))
	if err != nil {
		log.Fatal().Msgf("Read ssh key problems: %v", err)
	}

	conf.SshSigner, err = ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal().Msgf("Parse ssh key problems: %v", err)
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

	// handlers
	app.Get("/", func(ctx *fiber.Ctx) error {
		//_ = ctx.SendFile("./index.html")
		ctx.Type("html", "utf-8")
		return ctx.Send(_indexHtml)
	})
	app.Use(favicon.New(favicon.Config{
		File: filepath.Join(conf.WorkDir, "static/favicon.ico"),
	}))
	app.Static("/static", filepath.Join(conf.WorkDir, "static"))

	// API
	v1 := app.Group("/api/v1")
	{
		v1.Get("/version", func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{"version": Version})
		})

		v1.Get("/mikrotik", func(ctx *fiber.Ctx) error {
			ver, err := MikrotikGetVersion(conf)
			if err != nil {
				log.Fatal().Msgf("run problems: %v", err)
			}
			ver = strings.TrimSpace(ver)
			log.Debug().Msg(ver)

			return ctx.JSON(fiber.Map{"version": ver})
		})

		v1.Get("/provider", func(ctx *fiber.Ctx) error {
			provider, err := MikrotikGetProvider(conf)
			if err != nil {
				return srv.Error(ctx, fiber.StatusInternalServerError,
					"run get_provider err: "+err.Error())
			}
			return ctx.JSON(fiber.Map{"provider": strings.TrimSpace(provider)})
		})

		v1.Post("/switch", func(ctx *fiber.Ctx) error {
			provider, err := MikrotikSwitchProvider(conf)
			if err != nil {
				return srv.Error(ctx, fiber.StatusInternalServerError,
					"run switch_provider err: "+err.Error())
			}
			return ctx.JSON(fiber.Map{"provider": strings.TrimSpace(provider)})
		})

		// Dune test
		v1.Get("/dune/names", func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{"names": conf.Dunes.Name})
		})

		v1.Get("/dune/:id/status", func(ctx *fiber.Ctx) error {
			status := GetDuneStatusById(conf, ctx.Params("id"))
			return ctx.JSON(fiber.Map{"status": status})
		})

		v1.Get("/dune/:id/on", func(ctx *fiber.Ctx) error {
			status := DuneOnById(conf, ctx.Params("id"))
			return ctx.JSON(fiber.Map{"status": status})
		})

		v1.Get("/dune/:id/off", func(ctx *fiber.Ctx) error {
			status := DuneOffById(conf, ctx.Params("id"))
			return ctx.JSON(fiber.Map{"status": status})
		})
	}
	app.Use(srv.NotFound)
	return &srv
}

func (srv *App) Serve() error {
	// Create tls certificate
	log.Info().Msgf("Start server http://%s", srv.Config.ServerAddr)

	if err := srv.App.Listen(srv.Config.ServerAddr); err != nil {
		log.Fatal().Msgf("Listen error: %v", err)
		return err
	}
	log.Info().Msg("Server stopped")
	return nil
}

func (srv *App) Shutdown() error {
	if err := srv.App.Shutdown(); err != nil {
		log.Fatal().Msgf("Shutdown error %v", err)
		return err
	}
	return nil
}

func (srv *App) Error(c *fiber.Ctx, status int, message string) error {
	c.Status(status)
	log.Info().Msg(message)
	return c.JSON(fiber.Map{"status": status, "error": message})
}

func (srv *App) NotFound(c *fiber.Ctx) error {
	return srv.Error(
		c,
		fiber.StatusNotFound,
		"Sorry, but the page you were looking for could not be found.",
	)
}
