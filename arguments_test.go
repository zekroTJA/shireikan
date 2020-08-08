package shireikan

import "testing"

func TestArgumentAsString(t *testing.T) {
	if "test" != Argument("test").AsString() {
		t.Error("string conversion fialed")
	}
}

func TestArgumentAsInt(t *testing.T) {
	_, err := Argument("non-int").AsInt()
	if err == nil {
		t.Error("non-int conversion shall return error")
	}

	i, err := Argument("123").AsInt()
	if err != nil {
		t.Error("int conversion shall not return error")
	}
	if i != 123 {
		t.Error("recovered int is invalid")
	}
}

func TestArgumentAsFloat64(t *testing.T) {
	_, err := Argument("non-float").AsFloat64()
	if err == nil {
		t.Error("non-float conversion shall return error")
	}

	i, err := Argument("1.23").AsFloat64()
	if err != nil {
		t.Error("float conversion shall not return error")
	}
	if i != 1.23 {
		t.Error("recovered float is invalid")
	}
}

func TestArgumentAsBool(t *testing.T) {
	_, err := Argument("non-bool").AsBool()
	if err == nil {
		t.Error("non-bool conversion shall return error")
	}

	i, err := Argument("TRUE").AsBool()
	if err != nil {
		t.Error("bool conversion shall not return error")
	}
	if i != true {
		t.Error("recovered bool is invalid")
	}
}

func TestArgumentListGet(t *testing.T) {
	arr := []string{"a", "b", "c"}
	list := ArgumentList(arr)

	for i, v := range arr {
		if list.Get(i) != Argument(v) {
			t.Error("recovered value does not match")
		}
	}

	if list.Get(-1) != Argument("") {
		t.Error("getting value of invalid index shall return empty argument")
	}

	if list.Get(len(arr)) != Argument("") {
		t.Error("getting value of invalid index shall return empty argument")
	}
}

func TestArgumentListIndexOf(t *testing.T) {
	list := ArgumentList([]string{"a", "b", "c"})

	if list.IndexOf("b") != 1 {
		t.Error("recovered index is invalid")
	}

	if list.IndexOf("d") != -1 {
		t.Error("index of non-existent value shall be -1")
	}
}

func TestArgumentListContains(t *testing.T) {
	list := ArgumentList([]string{"a", "b", "c"})

	if !list.Contains("c") {
		t.Error("contained value shall be recovered as contained")
	}

	if list.Contains("e") {
		t.Error("non-contained value shall not be recovered as contained")
	}
}

func TestArgumentListSplice(t *testing.T) {
	list := ArgumentList([]string{"a", "b", "c", "d", "e"})
	listCheck := []string{"a", "d", "e"}

	listSpliced := list.Splice(1, 2)

	if len(listSpliced) != len(listCheck) {
		t.Error("spliced list does not match check list")
	}

	for i, v := range listSpliced {
		if v != listCheck[i] {
			t.Error("spliced list has different values to check list")
		}
	}

	listSpliced = list.Splice(len(list), 1)
	if len(listSpliced) != len(list) {
		t.Error("list splice with i > len(list) shall return list")
	}

	listSpliced = list.Splice(len(list)-1, 3)
	if len(listSpliced) != len(list)-1 {
		t.Error("list splice with i+r > len(list) shall return list[:i]")
	}
}
