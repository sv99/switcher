package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"switcher"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type WithLogFileCallback func() error

func WithLogFile(workDir string, withConsole bool, callback WithLogFileCallback) error {
	// init logger
	logFile, err := GetLoggerFile(workDir)
	if err != nil {
		log.Fatal().Err(err).Msg("error open log file")
	}
	defer func(f *os.File) {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    true,
			TimeFormat: time.DateTime,
		})
		_ = f.Close()
	}(logFile)
	InitLogger(logFile, withConsole)

	return callback()
}

func GetLoggerFile(workDir string) (*os.File, error) {
	// log to file for all
	packageName := reflect.TypeOf(switcher.Empty{}).PkgPath()
	logDir := filepath.Join(workDir, "log")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(workDir, "log"), 0755)
		if err != nil {
			return nil, err
		}
	}
	filename := filepath.Join(logDir, packageName+".log")
	// open an output file, this will append to the today's file if server restarted.
	// TODO: logfile rotation
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func InitLogger(logFile *os.File, console bool) {
	fileWriter := zerolog.ConsoleWriter{
		Out:        logFile,
		NoColor:    true,
		TimeFormat: time.DateTime,
	}
	// service not log to file if exists console out!!!
	if console {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: time.DateTime,
		}
		multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(fileWriter).With().Timestamp().Logger()
	}
}

func NewLoggerMiddleware(logger *zerolog.Logger) fiber.Handler {
	// default format "${status} ${method} ${path} - ${ip} ${latency}\n"
	// Middleware function
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// handle request
		err := c.Next()
		// build log
		stop := time.Now()

		logger.Info().Msgf("%d %s %s - %s %s",
			c.Response().StatusCode(), c.Method(), c.Path(), c.IP(),
			stop.Sub(start).String())
		return err
	}
}
