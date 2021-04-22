package shireikan

// ReadonlyObjectMap provides a thread save
// key-value map to get previously set items
// from.
type ReadonlyObjectMap interface {

	// GetObject returns a value from the
	// object map by its key.
	//
	// Returns 'nil' when no object is stored
	// with the given key.
	GetObject(key string) interface{}

	// SetObject IS DEPRECTAED and will be
	// removed in later versions!
	//
	// SetObject sets a value to the object
	// map linked with the given key.
	SetObject(key string, value interface{})
}

// ObjectMap provides a thread save key-value
// map to set values to and get values from back.
type ObjectMap interface {
	ReadonlyObjectMap

	// SetObject sets a value to the object
	// map linked with the given key.
	SetObject(key string, value interface{})
}
