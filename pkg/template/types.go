package template

import (
	"text/template"

	"github.com/bakito/docs-gen/internal/common"
)

type Option func(c *config)

func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.prefix = prefix
	}
}

func WithSuffix(suffix string) Option {
	return func(c *config) {
		c.suffix = suffix
	}
}

type config struct {
	common.Config
	prefix   string
	suffix   string
	template *template.Template
}
