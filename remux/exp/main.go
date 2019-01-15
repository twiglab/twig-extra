package main

import (
	"github.com/twiglab/twig"
	"github.com/twiglab/twig-extra/remux"
)

func main() {
	web := twig.TODO()
	web.WithMuxer(remux.New(remux.Pattern))

	twig.Config(web).
		Get("/a/:name/b/:other", twig.HelloTwig).
		Done()

	twig.Start(web)

	twig.Signal(twig.Graceful(web, 15))
}
