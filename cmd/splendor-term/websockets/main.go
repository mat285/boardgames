package main

import (
	"context"
	"fmt"

	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	client "github.com/mat285/boardgames/pkg/client/http/v1alpha1"
)

func main() {

	ctx := context.Background()

	cli := client.New(
		client.OptConfig(client.Config{
			Scheme: "http",
			Host:   "localhost",
			Port:   8080,
		}),
		client.OptUsername("michael"),
	)

	term := terminal.NewTerminalPlayerHTTP("michael", cli)
	err := term.Start(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
