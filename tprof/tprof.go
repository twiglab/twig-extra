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
