package hrmux

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/twiglab/twig"
)

func HelloHttpRouter(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{'H', 'e', 'l', 'l', 'o', '!'})
}

func Warp(h httprouter.Handle, p httprouter.Params) twig.HandlerFunc {
	if h == nil {
		return twig.NotFoundHandler
	}

	return func(c twig.Ctx) error {
		h(c.Resp(), c.Req(), p)
		return nil
	}
}

type HttpRouterMux struct {
	*httprouter.Router
	m []twig.MiddlewareFunc
}

func New() *HttpRouterMux {
	r := &HttpRouterMux{
		Router: httprouter.New(),
	}
	return r
}

func (hrx *HttpRouterMux) Lookup(method, path string, r *http.Request, c twig.Ctx) {
	if h, p, _ := hrx.Router.Lookup(method, path); h != nil {
		handle := Warp(h, p)
		c.SetHandler(twig.Enhance(handle, hrx.m))
	}
}

func (hrx *HttpRouterMux) Use(m ...twig.MiddlewareFunc) {
	hrx.m = append(hrx.m, m...)
}

func (hrx *HttpRouterMux) Add(method, path string, h twig.HandlerFunc, m ...twig.MiddlewareFunc) *twig.Route {
	panic("HttpRouterMux is not supports Add func!")
}
