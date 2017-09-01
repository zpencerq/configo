package configo

import (
	"reflect"

	"github.com/BurntSushi/toml"
)

type TomlConfigo struct {
	file string
}

func NewTomlConfigo(file string) *TomlConfigo {
	return &TomlConfigo{file: file}
}

func (tc *TomlConfigo) Load(v interface{}) error {
	return FromTOML(tc.file, v)
}

// FromTOML decodes the contents of the file `f` in TOML format into a pointer `v`.
//
// A field's value will be determined based on the following order:
//
// 1. If the field exists in the file, its value will be used. The `toml` tag may be used to map TOML keys to fields that don't match the key name exactly.
// 2. If `v` already contains a value for the field, it will be used.
func FromTOML(f string, v interface{}) error {
	rv := reflect.ValueOf(v).Elem()

	// Unmarshalling TOML onto a non-zero struct is inconsistent.
	// One time the value might be the pre-existing value, another time
	// it might be from the TOML. Instead we unmarshal onto a new struct
	// then walk the struct copying non-zero values.

	nv := reflect.New(rv.Type())
	ni := nv.Interface()
	_, err := toml.DecodeFile(f, ni)
	if err != nil {
		return err
	}

	return setToml(&rv, nv.Elem())
}

func setToml(dst *reflect.Value, src reflect.Value) error {
	var err error

	// TODO: Don't assume src and dst are the same

	switch dst.Kind() {
	case reflect.Ptr:
		dnv := dst.Elem()
		snv := src.Elem()

		err = setToml(&dnv, snv)
		if err != nil {
			return err
		}

		dst = &dnv
	case reflect.Struct:
		for i := 0; i < dst.NumField(); i++ {
			dval := dst.Field(i)
			sval := src.Field(i)

			kind := dval.Kind()

			if kind == reflect.Struct {
				err = setToml(&dval, sval)
				if err != nil {
					return err
				}

				continue
			}

			if kind == reflect.Ptr && dval.Elem().Kind() == reflect.Struct {
				err = setToml(&dval, sval)
				if err != nil {
					return err
				}

				continue
			}

			if !isZero(sval) {
				dval.Set(sval)
			}
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
		if !isZero(src) {
			dst.Set(src)
		}
	default:
		// TODO: Do something with unknown value types
	}

	return nil
}
