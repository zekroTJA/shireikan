package shireikan

// MiddlewareLayer defines the layer when the middleware
// handler shall be executed during command parsing and
// handling.
type MiddlewareLayer int

const (
	// Execute right before command handler is executed.
	LayerBeforeCommand MiddlewareLayer = 1 << iota
	// Execute right after command handler was executed successfully.
	LayerAfterCommand
)

// Middleware specifies a command middleware.
type Middleware interface {

	// Handle is called right before the execution of
	// the command handler and is getting passed the
	// Command instance and the Context.
	//
	// When the returned bool is false, the following
	// command handler will not be executed.
	//
	// An error should only be returned when the
	// execution of the middleware handler failed
	// unexpectedly.
	Handle(cmd Command, ctx Context) (bool, error)

	// GetLayer returns the layer(s) when the middleware
	// shall be executed.
	//
	// This value is defined as a bitmask value, so you
	// can combine multiple layers to execute the
	// middleware at multiple points during command
	// handling.
	GetLayer() MiddlewareLayer
}
