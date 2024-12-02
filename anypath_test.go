package anypath

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRefs(t *testing.T) {
	x := map[string]any{
		"foo": []int{1, 2, 3},
	}

	ap := Anypath{Raw: x}
	ary1any, err := ap.Extract("foo")
	if err != nil {
		t.Fatalf("%s", err)
	}
	mp := ap.Raw.(map[string]any)
	ary1 := (*ary1any).([]int)
	ary2 := mp["foo"]

	if !reflect.DeepEqual(ary1, ary2) {
		t.Fatal("Arrays aren't the same")
	}
	ary1[1] = 27
	if !reflect.DeepEqual(ary1, ary2) {
		t.Fatal("Arrays aren't the same")
	}
}

func TestExtract(t *testing.T) {

	type TstType struct {
		A int
		b string
		C map[string]any
	}

	x := map[string]any{
		"Test":  "Passed",
		"Hello": "World",
		"Foo":   "Bar",
		"Nested": map[string]any{
			"Help": "Me!",
			"Null": nil,
		},
		"counters": []any{
			1,
			"ary 2",
			"ary c",
		},

		"s": TstType{
			A: 10,
			b: "Hello, world",
			C: map[string]any{
				"foo":  "bar",
				"some": "other",
			},
		},
		"easyAry": []int{1, 2, 3},
	}

	type oneTest struct {
		Path        string
		Success     bool
		Expect      any
		CheckExpect bool
	}

	tests := []oneTest{
		oneTest{Path: "Test", Success: true, Expect: "Passed"},
		oneTest{Path: "", Success: true, Expect: x},
		oneTest{Path: ".Test", Success: true, Expect: "Passed"},
		oneTest{Path: ".Test.abc", Success: false},
		oneTest{Path: "Nested.Help", Success: true, Expect: "Me!"},
		oneTest{Path: "Nested.Null", Success: true, Expect: nil, CheckExpect: true},
		oneTest{Path: "counters[0]", Success: true, Expect: 1},
		oneTest{Path: "counters[1]", Success: true, Expect: "ary 2"},
		oneTest{Path: "counters[2]", Success: true, Expect: "ary c"},
		oneTest{Path: "counters[-1]", Success: true},
		oneTest{Path: "counters[-2]", Success: true},
		oneTest{Path: "counters[-3]", Success: true},
		oneTest{Path: "Nested.Foo", Success: false},
		oneTest{Path: "Nested.Null[4]", Success: false},
		oneTest{Path: "Nested.Foo[123]", Success: false},
		oneTest{Path: "Nested[11]", Success: false},
		oneTest{Path: "counters.abc", Success: false},
		oneTest{Path: "s", Success: true},
		oneTest{Path: "s.A", Success: true},
		oneTest{Path: "s.b", Success: false},
		oneTest{Path: "s.C", Success: true},
		oneTest{Path: "s.C.xxx", Success: false},
		oneTest{Path: "s.C.foo", Success: true},
		oneTest{Path: "Test", Success: true, Expect: "Passed"},
		oneTest{Path: "easyAry", Success: true, Expect: []int{1, 2, 3}},
		oneTest{
			Path:    "s.C",
			Success: true,
			Expect: map[string]any{
				"foo":  "bar",
				"some": "other",
			},
		},
	}

	ap := Anypath{Raw: x}

	for _, tc := range tests {
		t.Run(
			fmt.Sprintf("check_path:%s", tc.Path),
			func(t *testing.T) {
				res, err := ap.Extract(tc.Path)
				if tc.Success && err != nil {
					t.Fatalf("%s", err)
				}

				if tc.Expect != nil || tc.CheckExpect {
					if !reflect.DeepEqual(*res, tc.Expect) {
						t.Fatalf("expected: %+v\ngot: %+v", tc.Expect, *res)
					}
				}
			})

	}
}

func TestExists(t *testing.T) {
	type oneTest struct {
		Path   string
		Result bool
	}

	x := map[string]any{
		"foo": "bar",
		"nested": map[int64]any{
			10:  11,
			-11: 27,
		},
		"nested2": map[uint]any{
			10: 11,
			17: 28,
		},
	}

	tests := []oneTest{
		oneTest{Path: "", Result: true},
		oneTest{Path: "foo", Result: true},
		oneTest{Path: "bar", Result: false},

		// oneTest{Path: "nested[10]", Result: true},
		// oneTest{Path: "nested[-11]", Result: true},
		// oneTest{Path: "nested[110]", Result: false},
		// oneTest{Path: "nested[10].abc", Result: false},
		// oneTest{Path: "nested[10].foo.bar[11]", Result: false},

		// oneTest{Path: "nested2[10]", Result: true},
		// oneTest{Path: "nested2[126]", Result: false},
		// oneTest{Path: "nested2[-12345]", Result: false},
	}

	a := Anypath{Raw: x}

	for _, tc := range tests {
		t.Run(
			fmt.Sprintf("Checking path '%s' (expect %t)", tc.Path, tc.Result),
			func(t *testing.T) {

				if a.Exists(tc.Path) != tc.Result {
					t.Fatalf("Unexpected result with path '%s', expected %t, got %t",
						tc.Path,
						tc.Result,
						!tc.Result)
				}
			})
	}
}
