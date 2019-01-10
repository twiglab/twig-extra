package tprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/twiglab/twig"
)

var TprofIndex = twig.WrapHttpHandler(http.HandlerFunc(pprof.Index))
var TprofCmdLine = twig.WrapHttpHandler(http.HandlerFunc(pprof.Cmdline))
var TprofProfile = twig.WrapHttpHandler(http.HandlerFunc(pprof.Profile))
var TprofSymbol = twig.WrapHttpHandler(http.HandlerFunc(pprof.Symbol))
var TprofTrace = twig.WrapHttpHandler(http.HandlerFunc(pprof.Trace))

type Prefix string

func (t Prefix) String() string {
	return string(t)
}

func (t Prefix) Url(postfix string) string {
	return t.String() + postfix
}

func (t Prefix) Mount(mux twig.Register) {
	twig.Config(mux).
		Get(t.Url("/"), TprofIndex).
		Get(t.Url("/*"), TprofIndex).
		Get(t.Url("/cmdline"), TprofCmdLine).
		Get(t.Url("/profile"), TprofProfile).
		Get(t.Url("/symbol"), TprofSymbol).
		Get(t.Url("/trace"), TprofTrace).
		Done()
}

var prof = Prefix("/debug/pprof")

func Mount(r twig.Register) {
	prof.Mount(r)
}
