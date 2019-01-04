package xsqlt

import (
	"github.com/twiglab/sqlt"
	"github.com/twiglab/twig"
)

type SqltPlugin struct {
	*sqlt.Dbop

	id   string
	name string
}

func (x *SqltPlugin) ID() string {
	return x.id
}

func (x *SqltPlugin) Name() string {
	return x.name
}

func New(id, name string, dbop *sqlt.Dbop) *SqltPlugin {
	return &SqltPlugin{
		Dbop: dbop,
		id:   id,
		name: name,
	}
}

func GetSqltPlugin(id string, c twig.Ctx) *SqltPlugin {
	p := twig.GetPlugin(id, c)
	if plugin, ok := p.(*SqltPlugin); ok {
		return plugin
	}
	return nil
}
