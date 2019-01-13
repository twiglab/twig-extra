package remux

import (
	"net/http"

	"github.com/twiglab/twig"
)

type node struct {
	handler twig.HandlerFunc
	regex   UrlRegex
	path    string
}

func newNode(h twig.HandlerFunc, regex string) *node {
	return &node{
		handler: h,
		regex:   Pattern(regex),
		path:    regex,
	}
}

type ctg map[string][]*node

func newCtg() ctg {
	return make(ctg)
}

func (n ctg) add(method, path string, h twig.HandlerFunc) *node {
	node := newNode(h, path)
	n[method] = append(n[method], node)
	return node
}

func (n ctg) get(method string) []*node {
	return n[method]
}

type RegexMux struct {
	root   ctg
	m      []twig.MiddlewareFunc
	routes map[string]twig.Route
	t      *twig.Twig
}

func New() *RegexMux {
	return &RegexMux{
		root:   newCtg(),
		routes: make(map[string]twig.Route),
	}
}

func (r *RegexMux) Attach(t *twig.Twig) {
	r.t = t
}

func (r *RegexMux) Use(m ...twig.MiddlewareFunc) {
	r.m = append(r.m, m...)
}

func (r *RegexMux) Lookup(method, path string, req *http.Request) twig.Ctx {

	nodes := r.root.get(method)

	for _, node := range nodes {
		if param, err := node.regex.Match(path); err == nil { // ok
			c.SetPath(node.path)
			c.SetHandler(twig.Enhance(node.handler, r.m))
			c.SetRoutes(r.routes)

			twigParam(param, c)

			return
		}
	}
}

func (r *RegexMux) AddHandler(method string, path string, h twig.HandlerFunc, m ...twig.MiddlewareFunc) twig.Route {
	handler := twig.Enhance(h, m)
	r.root.add(method, path, handler)
	rd := &twig.NamedRoute{
		M: method,
		P: path,
		N: twig.HandlerName(h),
	}
	r.routes[rd.ID()] = rd
	return rd
}

func twigParam(p map[string]string, mc twig.MCtx) {
	l := len(p)
	if l > twig.MaxParam {
		panic("len(param) > twig.MaxParam")
	}

	names := make([]string, l, l)
	values := mc.ParamValues()

	i := 0
	for k, v := range p {
		names[i] = k
		values[i] = v
		i++
	}

	mc.SetParamNames(names)
	return
}
