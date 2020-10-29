package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AccessLoggerConfig struct {
	// Type determines whether zap will be initialized as a file logger or,
	// by default, as a console logger.
	Type string

	// Environment determines whether zap will be initialized using a production
	// or a development logger.
	Environment string

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool
}

func AccessLogger(config *AccessLoggerConfig) fiber.Handler {
	var logger *zap.Logger

	switch strings.ToLower(config.Type) {
		case "file":
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   config.Filename,
				MaxSize:    config.MaxSize,
				MaxAge:     config.MaxAge,
				MaxBackups: config.MaxBackups,
				LocalTime:  config.LocalTime,
				Compress:   config.Compress,
			})
			// Create a zap core object for JSON encoding
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				w,
				zap.InfoLevel,
			)
			// Create a zap logger object
			logger = zap.New(core)
			break
		// Use Access Logger in console based on environment by default
		default:
			var err error
			if strings.ToLower(config.Environment) == "production" {
				logger, err = zap.NewProduction()
			} else {
				logger, err = zap.NewDevelopment()
			}
			if err != nil {
				fmt.Println(err.Error())
			}
			break
	}

	// Flush logger buffers, if any
	defer logger.Sync()

	return func(ctx *fiber.Ctx) error {
		// Handle the request to calculate the number of bytes sent
		err := ctx.Next()

		// Chained error
		if err != nil {
			if chainErr := ctx.App().Config().ErrorHandler(ctx, err); chainErr != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}

		// Send structured information message to the logger
		logger.Info(ctx.IP()+" - "+ctx.Method()+" "+ctx.OriginalURL()+" - "+strconv.Itoa(ctx.Response().StatusCode())+
			" - "+strconv.Itoa(len(ctx.Response().Body())),

			zap.String("ip", ctx.IP()),
			zap.String("hostname", ctx.Hostname()),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.OriginalURL()),
			zap.String("protocol", ctx.Protocol()),
			zap.Int("status", ctx.Response().StatusCode()),

			zap.String("x-forwarded-for", ctx.Get(fiber.HeaderXForwardedFor)),
			zap.String("user-agent", ctx.Get(fiber.HeaderUserAgent)),
			zap.String("referer", ctx.Get(fiber.HeaderReferer)),

			zap.Int("bytes_received", len(ctx.Request().Body())),
			zap.Int("bytes_sent", len(ctx.Response().Body())),
		)

		return err
	}
}
