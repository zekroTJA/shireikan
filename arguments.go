package shireikan

// Argument extends string to provide general
// transformation functionality.
type Argument string

// ArgumentList wraps a string list to get
// arguments in that list as Argument object.
type ArgumentList []string

// Get returns the Argument at the given Index.
// If there is no argument at that index, an
// empty string is returned.
func (al ArgumentList) Get(i int) Argument {
	if i < 0 || i >= len(al) {
		return Argument("")
	}

	return Argument(al[i])
}
