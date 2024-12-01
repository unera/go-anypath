package anypath

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

// Anypath is a common data container
type Anypath struct {
	Raw any // Raw unmarshaled data
}

// LoadYamlString construct container from yaml
func LoadYamlString(y string) (*Anypath, error) {
	var data any = nil

	if err := yaml.Unmarshal([]byte(y), &data); err != nil {
		return nil, err
	}
	return &Anypath{Raw: data}, nil
}

// Extract raw data by path
func (a *Anypath) Extract(path string) (*any, error) {

	pathElements, err := parsePath(path)
	if err != nil {
		return nil, err
	}
	return a.extract(pathElements, "", &a.Raw)
}

func (a *Anypath) extract(path []any, tracePath string, o *any) (*any, error) {
	if len(path) == 0 {
		return o, nil
	}

	path0 := path[0]
	var fullPath string
	switch path0.(type) {
	case int:
		fullPath = fmt.Sprintf("%s[%d]", tracePath, path0)
	case string:
		fullPath = fmt.Sprintf("%s.%s", tracePath, path0)
	}
	if *o == nil {
		return nil, fmt.Errorf("%s: Can not dereference Nil", fullPath)
	}
	// Invalid Kind = iota
	// Bool
	// Int
	// Int8
	// Int16
	// Int32
	// Int64
	// Uint
	// Uint8
	// Uint16
	// Uint32
	// Uint64
	// Uintptr
	// Float32
	// Float64
	// Complex64
	// Complex128
	// Chan
	// Func
	// Interface
	// Map
	// Pointer
	// String
	// Struct
	// UnsafePointer

	t, v := reflect.TypeOf(*o), reflect.ValueOf(*o)

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		switch path0.(type) {
		case int:
			var (
				index int
				newO  any
			)
			index = path0.(int)
			if index >= v.Len() {
				goto wrongPath
			}
			if index < 0 {
				index = v.Len() + index
				if index < 0 {
					goto wrongPath
				}
			}
			newO = v.Index(index).Interface()
			return a.extract(path[1:], fullPath, &newO)
		wrongPath:
			return nil, fmt.Errorf("%s: Can not find path", fullPath)

		default:
			return nil, fmt.Errorf("%s: Can not find path", fullPath)
		}
	case reflect.Map:
		if _, ok := path0.(string); !ok {
			return nil, fmt.Errorf("%s: Can not find path 6", fullPath)
		}
		extractedO := v.MapIndex(reflect.ValueOf(path0))
		if !extractedO.IsValid() {
			return nil, fmt.Errorf("%s: Can not find path 3", fullPath)
		}
		newO := extractedO.Interface()
		return a.extract(path[1:], fullPath, &newO)

	case reflect.Struct:
		name, ok := path0.(string)
		if !ok {
			return nil, fmt.Errorf("%s: Can not find path 10", fullPath)
		}
		sfield, ok := t.FieldByName(name)
		if !ok {
			return nil, fmt.Errorf("%s: Can not find path 11", fullPath)
		}
		if !sfield.IsExported() {
			return nil, fmt.Errorf("%s: Path not exported", fullPath)
		}
		field := v.FieldByName(name)
		if !field.IsValid() {
			return nil, fmt.Errorf("%s: Can not find path 13", fullPath)
		}
		newO := field.Interface()
		return a.extract(path[1:], fullPath, &newO)

	}
	return nil, fmt.Errorf("%s: Can not find path", fullPath)
}
