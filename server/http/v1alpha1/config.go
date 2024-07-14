package v1alpha1

import (
	"context"
	"os"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Web web.Config `json:"web" yaml:"web"`
	TLS TLS        `json:"tls" yaml:"tls"`
}

// Resolve populates configuration fields from a variety of input sources
func (c *Config) Resolve(ctx context.Context, files ...string) error {
	if err := resolveFromFiles(&c, files...); err != nil {
		return err
	}
	return configutil.Resolve(ctx,
		(&c.Web).Resolve,
	)
}

func (c Config) Context(ctx context.Context) (context.Context, error) {
	ctx = logger.WithLogger(ctx, logger.All())
	return ctx, nil
}

type FileResolver struct {
	out   interface{}
	files []string
}

func (fr *FileResolver) Resolve(ctx context.Context) error {
	return resolveFromFiles(fr.out, fr.files...)
}

func readYamlFile(file string, out interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

func resolveFromFiles(out interface{}, files ...string) error {
	for _, file := range files {
		err := readYamlFile(file, out)
		if err != nil {
			return err
		}
	}
	return nil
}
