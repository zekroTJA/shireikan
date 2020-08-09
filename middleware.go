package shireikan

type MiddlewareLayer int

const (
	LayerBeforeCommand MiddlewareLayer = 1 << iota
	LayerAfterCommand
)

// Middleware specifies a command middleware.
type Middleware interface {

	// Handle is called right before the execution of
	// the command handler and is getting passed the
	// Command instance and the Context.
	//
	// When the returned bool is false, the follwoing
	// command handler will not be executed.
	//
	// An error should only be returned when the
	// execution of the middleware handler failed
	// unexpectedly.
	Handle(cmd Command, ctx Context) (error, bool)

	GetLayer() MiddlewareLayer
}
