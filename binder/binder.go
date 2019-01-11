package binder

import "github.com/twiglab/twig"

// Binder 数据绑定接口
// Binder 作为一个插件集成到Twig中,请实现Plugin接口
type Binder interface {
	Bind(interface{}, twig.Ctx) error
}

// GetBinder 获取绑定接口
func GetBinder(id string, c twig.Ctx) (binder Binder, ok bool) {
	var plugin twig.Plugin
	plugin, ok = twig.GetPlugin(id, c)
	if !ok {
		return
	}
	binder, ok = plugin.(Binder)
	return
}
