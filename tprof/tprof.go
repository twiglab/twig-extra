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

func (t Prefix) Mount(mux twig.Register) {
	mux.AddHandler(twig.GET, t.String()+"/", TprofIndex)
	mux.AddHandler(twig.GET, t.String()+"/*", TprofIndex)

	mux.AddHandler(twig.GET, t.String()+"/cmdline", TprofCmdLine)
	mux.AddHandler(twig.GET, t.String()+"/profile", TprofProfile)
	mux.AddHandler(twig.GET, t.String()+"/symbol", TprofSymbol)
	mux.AddHandler(twig.GET, t.String()+"/trace", TprofTrace)
}
