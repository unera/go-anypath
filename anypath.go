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

// Exists returns true if path is present in the object
func (a *Anypath) Exists(path string) bool {
	_, err := a.Extract(path)
	return err == nil
}

func (a *Anypath) extract(path []any, tracePath string, o *any) (*any, error) {
	if len(path) == 0 {
		return o, nil
	}

	path0 := path[0]
	var (
		fullPath string
		index    int64
		newO     any
		uindex   uint64
		aindex   any
	)
	switch path0.(type) {
	case int64:
		fullPath = fmt.Sprintf("%s[%d]", tracePath, path0)
	case string:
		fullPath = fmt.Sprintf("%s.%s", tracePath, path0)
	}

	e := func(msg ...string) (*any, error) {
		if len(msg) == 0 {
			return nil, fmt.Errorf("%s: Can not find path", fullPath)
		}
		return nil, fmt.Errorf("%s: %s", fullPath, msg[0])
	}

	if *o == nil {
		return e("Can not dereference Nil")
	}

	t, v := reflect.TypeOf(*o), reflect.ValueOf(*o)

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		switch path0.(type) {
		case int64:
			index = path0.(int64)
			if int(index) >= v.Len() {
				goto wrongPath
			}
			if index < 0 {
				index = int64(v.Len()) + index
				if index < 0 {
					goto wrongPath
				}
			}
			newO = v.Index(int(index)).Interface()
			return a.extract(path[1:], fullPath, &newO)
		wrongPath:
			return e()

		default:
			return e()
		}

	case reflect.Map:
		switch t.Key().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if _, ok := path0.(int); !ok {
				return e()
			}
			extractedO := v.MapIndex(reflect.ValueOf(path0))
			if !extractedO.IsValid() {
				return e()
			}
			newO := extractedO.Interface()
			return a.extract(path[1:], fullPath, &newO)

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			index, ok := path0.(int64)
			if !ok || index < 0 {
				return e()
			}
			uindex = uint64(index)
			aindex = uint(uindex)
			goto mapIndexUint
		mapIndexUint:
			extractedO := v.MapIndex(reflect.ValueOf(aindex))
			if !extractedO.IsValid() {
				return e()
			}
			newO := extractedO.Interface()
			return a.extract(path[1:], fullPath, &newO)

		case reflect.String:
			if _, ok := path0.(string); !ok {
				return e()
			}
			extractedO := v.MapIndex(reflect.ValueOf(path0))
			if !extractedO.IsValid() {
				return e()
			}
			newO := extractedO.Interface()
			return a.extract(path[1:], fullPath, &newO)
		}

	case reflect.Struct:
		name, ok := path0.(string)
		if !ok {
			return e()
		}
		sfield, ok := t.FieldByName(name)
		if !ok {
			return e()
		}
		if !sfield.IsExported() {
			return e()
		}
		field := v.FieldByName(name)
		if !field.IsValid() {
			return e()
		}
		newO := field.Interface()
		return a.extract(path[1:], fullPath, &newO)

	}
	return e()
}
