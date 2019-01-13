package main

import (
	"github.com/twiglab/twig"
	"mikewang.info/remux"
)

func main() {
	web := twig.TODO()
	web.WithMuxer(remux.New())

	twig.Config(web).
		Get("/a/:name/b/:other/", twig.HelloTwig).
		Done()

	twig.Start(web)

	twig.Signal(twig.Quit())
}
