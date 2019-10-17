/*
Copyright © 2019 Ivan Goncharov <i.morph@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"time"

	"github.com/imorph/blank/lib/buildinfo"
	"github.com/imorph/blank/lib/signl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start to serve some requests",
	Long:  `This is main server entry-point`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()

		// configure logging
		logger, _ := initZap(viper.GetString("logLevel"))
		defer func() {
			err := logger.Sync()
			if err != nil {
				log.Println("error syncing logger", err)
			}
		}()

		stdLog := zap.RedirectStdLog(logger)
		defer stdLog()

		logger.Info("Application started",
			zap.Duration("startup_duration", time.Since(startTime)),
			zap.String("listen_address", viper.GetString("address")),
			zap.String("log_level", viper.GetString("logLevel")),
			zap.String("app_name", buildinfo.GetAppName()),
			zap.String("version", buildinfo.GetVersion()),
		)
		// log.Println("Started App in:", time.Since(startTime))
		// log.Println("Listen on:", viper.GetString("address"))
		// log.Println("Log Level:", viper.GetString("logLevel"))

		signal := signl.WaitForSigterm()
		log.Println("recieved signal:", signal)
		log.Println("exiting")
	},
}

var address, logLevel string

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.
	serverCmd.PersistentFlags().StringVar(&address, "address", "127.0.0.1:9999", "An adress server will listen on")
	serverCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "Logger verbosity level")
	err := viper.BindPFlag("logLevel", serverCmd.PersistentFlags().Lookup("logLevel"))
	if err != nil {
		log.Println("unable to bind to flag", err)
	}
	err = viper.BindPFlag("address", serverCmd.PersistentFlags().Lookup("address"))
	if err != nil {
		log.Println("unable to bind to flag", err)
	}

}

func initZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		log.Fatal("Dont know this log level:", logLevel, "known levels are: debug, info, warn, error, fatal, panic")
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}
