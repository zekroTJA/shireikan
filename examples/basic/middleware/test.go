package middleware

import "github.com/zekroTJA/shireikan"

// Test is a test command middleware.
type Test struct {
}

// Handle is the Middlewares handler.
func (m *Test) Handle(cmd shireikan.Command, ctx shireikan.Context, layer shireikan.MiddlewareLayer) (bool, error) {
	ctx.SetObject("test", "this is a test object")

	return true, nil
}

// GetLayer returns the execution layer.
func (m *Test) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerBeforeCommand
}
