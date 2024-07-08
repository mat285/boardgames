package main

import (
	"context"
	"fmt"
	"os"

	"github.com/blend/go-sdk/graceful"
	"github.com/mat285/boardgames/pkg/server/http/v1alpha1"
	"github.com/spf13/cobra"
)

func run() error {
	ctx := context.Background()
	cfg := v1alpha1.Config{}
	paths := []string{}
	cmd := &cobra.Command{
		Use:           "run-web",
		Short:         "Run web server",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			err := cfg.Resolve(ctx, paths...)
			if err != nil {
				return err
			}

			ctx, err := cfg.Context(ctx)
			if err != nil {
				return err
			}
			app := v1alpha1.New(ctx, cfg)
			return graceful.Shutdown(app)
		},
	}

	cmd.PersistentFlags().StringSliceVar(
		&paths,
		"file",
		paths,
		"Path to a file where '.yml' configuration is stored; can be specified multiple times, last provided has highest precedence when merging",
	)

	return cmd.Execute()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
