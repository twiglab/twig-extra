package tprof

import (
	"net/http/pprof"

	"github.com/twiglab/twig"
)

var TprofIndex = twig.WrapHttpHandlerFunc(pprof.Index)
var TprofCmdLine = twig.WrapHttpHandlerFunc(pprof.Cmdline)
var TprofProfile = twig.WrapHttpHandlerFunc(pprof.Profile)
var TprofSymbol = twig.WrapHttpHandlerFunc(pprof.Symbol)
var TprofTrace = twig.WrapHttpHandlerFunc(pprof.Trace)

type Prefix string

func (t Prefix) String() string {
	return string(t)
}

func (t Prefix) Url(postfix string) string {
	return t.String() + postfix
}

func (t Prefix) Mount(mux twig.Register) {
	twig.Cfg(mux).
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
