package remux

import "github.com/twiglab/twig"

type reCtx struct {
	*twig.VCtx
	path    string
	handler twig.HandlerFunc
	params  twig.UrlParams
}

func (c *reCtx) Path() string {
	return c.path
}

func (c *reCtx) Params() twig.UrlParams {
	return c.params
}

func (c *reCtx) Param(name string) string {
	return c.params[name]
}

func (c *reCtx) Handler() twig.HandlerFunc {
	return c.handler
}

func (c *reCtx) URL(name string, i ...interface{}) string {
	return ""
}

func (c *reCtx) Release() {
	c.SetFact(nil)
}
