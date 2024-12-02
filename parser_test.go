package anypath

import (
	"reflect"
	"testing"
)

func TestParsePath(t *testing.T) {
	type test struct {
		Path        string
		Await       []any
		IsError     bool
		Description string
	}

	tests := []test{
		test{
			Path:        "",
			Await:       []any{},
			IsError:     false,
			Description: "Empty path",
		},
		test{
			Path:        "a",
			Await:       []any{"a"},
			IsError:     false,
			Description: "One token",
		},
		test{
			Path:        "abcde",
			Await:       []any{"abcde"},
			IsError:     false,
			Description: "One token",
		},
		test{
			Path:        ".aaa",
			Await:       []any{"aaa"},
			IsError:     false,
			Description: "Point started token",
		},
		test{
			Path:        "abc.cde",
			Await:       []any{"abc", "cde"},
			IsError:     false,
			Description: "token after token",
		},
		test{
			Path:        "abc.cde.def",
			Await:       []any{"abc", "cde", "def"},
			IsError:     false,
			Description: "token after token",
		},
		test{
			Path:        "abc.cde.",
			Await:       []any{},
			IsError:     true,
			Description: "token after token, dot at EOF",
		},
		test{
			Path:        ".",
			Await:       []any{},
			IsError:     true,
			Description: "Empty path (one dot)",
		},
		test{
			Path:        "...",
			Await:       []any{},
			IsError:     true,
			Description: "Error path (some dots)",
		},
		test{
			Path:        "27",
			Await:       []any{"27"},
			IsError:     false,
			Description: "Digits in path",
		},
		test{
			Path:        "[27]",
			Await:       []any{int64(27)},
			IsError:     false,
			Description: "Index in path",
		},
		test{
			Path:        "[-2]",
			Await:       []any{int64(-2)},
			IsError:     false,
			Description: "Negative Index in path",
		},
		test{
			Path:        "[abc]",
			Await:       []any{},
			IsError:     true,
			Description: "Error index",
		},
		test{
			Path:        "[",
			Await:       []any{},
			IsError:     true,
			Description: "Error index brackets",
		},
		test{
			Path:        "[]",
			Await:       []any{},
			IsError:     true,
			Description: "Error index brackets",
		},
		test{
			Path:        "]",
			Await:       []any{},
			IsError:     true,
			Description: "Error index brackets",
		},
		test{
			Path:        "[123",
			Await:       []any{},
			IsError:     true,
			Description: "Error index brackets",
		},
		test{
			Path:        ".abc[123].def[345].cde",
			Await:       []any{"abc", int64(123), "def", int64(345), "cde"},
			IsError:     false,
			Description: "Composite path",
		},
		test{
			Path:        "[123457878728728728728728782728728]",
			Await:       []any{},
			IsError:     true,
			Description: "Too long integer index",
		},
		test{
			Path:        "abc[.cde",
			Await:       []any{},
			IsError:     true,
			Description: "error at the middle of string",
		},
		test{
			Path:        "abc[]cde",
			Await:       []any{},
			IsError:     true,
			Description: "error at the middle of string",
		},
	}

	for _, tc := range tests {
		t.Logf("Checking '%s' (%s)", tc.Path, tc.Description)

		res, err := parsePath(tc.Path)

		equals := true
		if len(res) == len(tc.Await) {
			for i, v := range res {
				if !reflect.DeepEqual(v, tc.Await[i]) {
					equals = false
					break
				}
			}
		} else {
			equals = false
		}

		if !equals {
			t.Fatalf("!Expected: <%+v>, got: <%+v> (err: %s)", tc.Await, res, err)
		}
		if tc.IsError {
			if err == nil {
				t.Fatalf("%s: Expected error, but got 'ok'", tc.Path)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: Unexpected error: %s", tc.Path, err)
			}
		}
	}
}
