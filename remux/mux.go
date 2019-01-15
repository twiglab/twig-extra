package remux

import (
	"net/http"

	"github.com/twiglab/twig"
)

type node struct {
	handler twig.HandlerFunc
	regex   *UrlRegex
	path    string
	method  string
}

type table map[string][]*node

func (t table) add(n *node) {
	ns := t[n.method]
	ns = append(ns, n)
	t[n.method] = ns
}

type ctg struct {
	table table
	fn    UrlRegexFunc
}

func newCtg(fn UrlRegexFunc) *ctg {
	return &ctg{
		table: make(table),
		fn:    fn,
	}
}

func (n *ctg) add(method, path string, h twig.HandlerFunc) {
	n.table.add(
		&node{
			handler: h,
			regex:   n.fn(path),
			path:    path,
			method:  method,
		})
}

func (n *ctg) get(method string) []*node {
	return n.table[method]
}

type RegexMux struct {
	root   *ctg
	m      []twig.MiddlewareFunc
	routes map[string]twig.Route
	t      *twig.Twig
}

func New(fn UrlRegexFunc) *RegexMux {
	return &RegexMux{
		root:   newCtg(fn),
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

	c := &reCtx{
		VCtx:    twig.NewVCtx(r.t),
		handler: twig.NotFoundHandler,
	}

	nodes := r.root.get(method)
	for _, node := range nodes {
		if params, ok := node.regex.Match(path); ok {
			c.params = params
			c.path = node.path
			c.handler = node.handler
			c.SetFact(c)
			return c
		}
	}
	return c
}

func (r *RegexMux) AddHandler(method string, path string, h twig.HandlerFunc, m ...twig.MiddlewareFunc) twig.Route {
	handler := twig.Merge(h, m)
	r.root.add(method, path, handler)
	rd := &twig.NamedRoute{
		M: method,
		P: path,
		N: twig.HandlerName(h),
	}
	r.routes[rd.ID()] = rd
	return rd
}
