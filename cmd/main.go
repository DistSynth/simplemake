package main

import (
	"os"

	"github.com/distsynth/simplemake/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

func main() {
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	app := &cli.App{
		Name:      "simplemake",
		Usage:     "simple replacement for make",
		Args:      true,
		ArgsUsage: "[tasks...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
			&cli.StringFlag{Name: "makefile", Aliases: []string{"f"}},
			&cli.BoolFlag{Name: "force", Value: false},
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all tasks",
				Action:  internal.ListTasks,
			},
			{
				Name:   "validate",
				Usage:  "Validate makefile tasks",
				Action: internal.ValidateMakeFile,
			},
		},
		Before: func(cCtx *cli.Context) error {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if cCtx.Bool("verbose") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			return nil
		},
		Action: internal.RunTasks,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msgf("Error occured: %s", err.Error())
	}
}
