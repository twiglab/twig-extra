package fcgi

import (
	"context"
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/twiglab/twig"
)

type FcgiServnat struct {
	file    string
	ln      net.Listener
	handler http.Handler
}

func NewFcgiServant(filename string) *FcgiServnat {
	return &FcgiServnat{
		file: filename,
	}
}

func (s *FcgiServnat) Start() (err error) {
	if s.ln, err = net.Listen("unix", s.file); err != nil {
		return
	}

	go func() {
		err = fcgi.Serve(s.ln, s.handler)
	}()

	return
}

func (s *FcgiServnat) Shutdown(c context.Context) error {
	return s.ln.Close()
}

func (s *FcgiServnat) Attach(t *twig.Twig) {
	s.handler = t
}
