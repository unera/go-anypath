package anypath

import (
	"testing"
)

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
	}

	type oneTest struct {
		Path    string
		Success bool
	}

	tests := []oneTest{
		oneTest{Path: "Test", Success: true},
		oneTest{Path: "", Success: true},
		oneTest{Path: ".Test", Success: true},
		oneTest{Path: ".Test.abc", Success: false},
		oneTest{Path: "Nested.Help", Success: true},
		oneTest{Path: "Nested.Null", Success: true},
		oneTest{Path: "counters[0]", Success: true},
		oneTest{Path: "counters[1]", Success: true},
		oneTest{Path: "counters[2]", Success: true},
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
		oneTest{Path: "s.C.some", Success: true},
	}

	ap := Anypath{Raw: x}

	for _, tc := range tests {
		t.Logf("Checking path `%s`: %v\n", tc.Path, tc.Success)
		_, err := ap.Extract(tc.Path)
		if tc.Success && err != nil {
			t.Fatalf("%s", err)
		}
	}

}
