package main

import (
	"os"

	upx "github.com/drone-stack/drone-upx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.1"
)

type formatter struct{}

func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func init() {
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	DisableTimestamp: true,
	// 	DisableColors:    true,
	// })
	logrus.SetFormatter(new(formatter))
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		_ = godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "docker upx"
	app.Usage = "docker upx"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.IntFlag{
			Name:   "level",
			Usage:  "level",
			EnvVar: "PLUGIN_LEVEL",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "path",
			EnvVar: "PLUGIN_PATH",
		},
		cli.StringFlag{
			Name:   "include",
			Usage:  "include",
			EnvVar: "PLUGIN_INCLUDE",
		},
		cli.StringFlag{
			Name:   "exclude",
			Usage:  "exclude",
			EnvVar: "PLUGIN_EXCLUDE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := upx.Plugin{
		Level:   c.Int("level"),
		Path:    c.String("path"),
		Include: c.String("include"),
		Exclude: c.String("exclude"),
		Debug:   c.Bool("debug"),
	}

	if plugin.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := plugin.Exec(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	return nil
}
