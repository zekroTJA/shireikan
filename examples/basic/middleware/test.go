package middleware

import "github.com/zekroTJA/shireikan"

type Test struct {
}

func (m *Test) Handle(cmd shireikan.Command, ctx shireikan.Context) (error, bool) {
	ctx.SetObject("test", "this is a test object")

	return nil, true
}

func (m *Test) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}
