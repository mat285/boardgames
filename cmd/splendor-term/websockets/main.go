package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mat285/boardgames/games/splendor/pkg/client/terminal"
	client "github.com/mat285/boardgames/pkg/client/http/v1alpha1"
)

func main() {

	ctx := context.Background()

	scheme := "http"
	host := "localhost"
	port := uint16(8080)

	if len(os.Args) > 1 {
		addr, err := Parse(os.Args[1])
		if err == nil {
			scheme = addr.Protocol
			host = addr.Host
			port = addr.Port
		}
	}

	cli := client.New(
		client.OptConfig(client.Config{
			Scheme: scheme,
			Host:   host,
			Port:   port,
		}),
		client.OptUsername("michael"),
	)

	term := terminal.NewTerminal(ctx, "michael", cli)
	if err := term.Start(); err != nil {
		log.Fatal(err)
	}
}

type Address struct {
	Protocol string
	Host     string
	Port     uint16
}

func Parse(addr string) (Address, error) {
	addr = strings.TrimSuffix(addr, "/")
	scheme := "http"
	if strings.HasPrefix(addr, "https") {
		scheme = "https"
	}
	addr = strings.TrimPrefix(addr, scheme+"://")
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return Address{}, fmt.Errorf("invalid address")
	}
	host := parts[0]
	port, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return Address{}, err
	}
	return Address{Protocol: scheme, Host: host, Port: uint16(port)}, nil
}
