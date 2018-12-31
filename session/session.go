package session

import (
	"github.com/gorilla/sessions"
	"github.com/twiglab/twig"
	"github.com/twiglab/twig/middleware"
)

type Config struct {
	Skipper middleware.Skipper
	Store   sessions.Store
}

const Key = "_gorilla_session_store_"

var DefaultConfig = Config{
	Skipper: middleware.DefaultSkipper,
}

func Middleware(store sessions.Session) twig.MiddlewareFunc {
	c := DefaultConfig
	c.Store = store
	return MiddlewareWithConfig(c)
}

func MiddlewareWithConfig(config Config) twig.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	if config.Store == nil {
		panic("Store is nil")
	}

	return func(next twig.HandlerFunc) twig.HandlerFunc {
		return func(c twig.Ctx) error {
			if config.Skipper(c) {
				return next(c)
			}
			c.Set(Key, config.Store)
			return next(c)
		}
	}
}
