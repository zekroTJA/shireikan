package shireikan

import "strconv"

// Argument extends string to provide general
// transformation functionality.
type Argument string

// AsString returns the argument as string.
func (a Argument) AsString() string {
	return string(a)
}

// AsInt tries to parse the given argument
// as integer. If this fails, an error is
// returned.
func (a Argument) AsInt() (int, error) {
	return strconv.Atoi(a.AsString())
}

// AsFloat64 tries to parse the given argument
// as float64. If this fails, an error is
// returned.
func (a Argument) AsFloat64() (float64, error) {
	return strconv.ParseFloat(a.AsString(), 64)
}

// AsBool tries to parse the given argument
// as bool. If this fails, an error is
// returned.
//
// As described in the strconv.ParseBool docs,
// the following values are accepted:
// "It accepts 1, t, T, TRUE, true, True, 0, f, F,
// FALSE, false, False. Any other value returns
// an error."
func (a Argument) AsBool() (bool, error) {
	return strconv.ParseBool(a.AsString())
}

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
