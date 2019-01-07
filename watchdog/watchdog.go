package doorgod

import (
	"github.com/casbin/casbin"
	"github.com/twiglab/twig"
	"github.com/twiglab/twig/middleware"
)

type CheckFunc func(twig.Ctx, *casbin.Enforcer) bool

type Config struct {
	Skipper  middleware.Skipper
	Enforcer *casbin.Enforcer

	Check CheckFunc
}

var DefaultConfig = Config{
	Skipper: middleware.DefaultSkipper,
	Check:   check,
}

func NewWithConfig(config Config) twig.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	if config.Check == nil {
		config.Check = DefaultConfig.Check
	}

	return func(next twig.HandlerFunc) twig.HandlerFunc {
		return func(c twig.Ctx) error {
			if config.Skipper(c) || config.Check(c, config.Enforcer) {
				return next(c)
			}

			return twig.ErrForbidden
		}
	}
}

func New(ce *casbin.Enforcer, check CheckFunc) twig.MiddlewareFunc {
	c := DefaultConfig
	c.Enforcer = ce
	c.Check = check
	return NewWithConfig(c)
}

func check(c twig.Ctx, ce *casbin.Enforcer) bool {
	username, _, ok := c.Req().BasicAuth()
	if !ok {
		return false
	}

	method := c.Req().Method
	path := c.Req().URL.Path

	return ce.Enforce(username, path, method)
}
