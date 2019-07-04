package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/k0kubun/pp"

	"github.com/cnosuke/gonzo/config"
	"github.com/cnosuke/gonzo/handler"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	// Version and Revision are replaced when building.
	// To set specific version, edit Makefile.
	Version  = "0.0.1"
	Revision = "xxx"

	Name  = "Gonzo"
	Usage = "Gyazo like image uploader"
)

var logger *zap.SugaredLogger

func main() {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapLogger, err := zapConfig.Build()
	if err != nil {
		fmt.Printf("Building logger error: %v", err)
		os.Exit(1)
	}

	defer zapLogger.Sync()
	logger = zapLogger.Sugar()

	undo := zap.ReplaceGlobals(zapLogger)
	defer undo()

	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s (%s)", Version, Revision)
	app.Name = Name
	app.Usage = Usage

	var (
		binding    string
		configPath string
	)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "binding, b",
			Usage:       "Server binding address",
			Value:       "127.0.0.1:8888",
			Destination: &binding,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Path to config file",
			Value:       "",
			Destination: &configPath,
		},
	}

	app.Action = func(c *cli.Context) error {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "revision", Revision)

		conf, err := config.NewConfig(configPath)

		if err != nil {
			pp.Println(err)
			return err
		}

		h := handler.NewHandler(ctx, conf)

		zap.S().Infow("Init imagine",
			"revision", Revision,
			"binding", binding,
			"conf", conf,
		)

		if err := http.ListenAndServe(binding, h.Routing()); err != nil {
			zap.S().Error("Failed to bind")
			return err
		}

		return nil
	}

	app.Run(os.Args)
}
