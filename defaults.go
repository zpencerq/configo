package configo

import (
	"fmt"
	"reflect"
)

type DefaultsConfigo struct{}

func NewDefaultsConfigo() *DefaultsConfigo { return &DefaultsConfigo{} }

func (dc *DefaultsConfigo) Load(v interface{}) error {
	return FromDefaults(v)
}

// FromDefaults sets pointer `v` based on default values of `v`.
//
// A field's value will be determined based on the following order:
//
// 1. If `v` already contains a value for the field, it will be used.
// 2. If a "default" tag exists for a field, its value will be used, subject to type casting.
// 3. The field will be initialized to its zero value (i.e., "" for string, 0 for int, etc).
func FromDefaults(v interface{}) error {
	rv := reflect.ValueOf(v).Elem()

	return setDefaults(&rv)
}

func setDefaults(v *reflect.Value) error {
	var err error

	// TODO: Properly initialize struct pointers
	// e.g. given "type Foo struct { Bar *OtherStruct }" set Bar's fields
	// to their zero-value.

	switch v.Kind() {
	case reflect.Ptr:
		if isZero(*v) {
			v.Set(reflect.New(v.Type().Elem()))
		}

		nv := v.Elem()

		err = setDefaults(&nv)
		if err != nil {
			return err
		}

		v = &nv
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			val := v.Field(i)
			typ := v.Type().Field(i)

			kind := val.Kind()

			if kind == reflect.Struct {
				err = setDefaults(&val)
				if err != nil {
					return err
				}

				continue
			}

			if kind == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
				err = setDefaults(&val)
				if err != nil {
					return err
				}

				continue
			}

			if isZero(val) {
				switch val.Kind() {
				case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
					val.Set(reflect.ValueOf(val.Interface()))
				case reflect.Ptr:
					val.Set(reflect.New(val.Type().Elem()))
				default:
					// TODO: Do something with unknown types
				}

				tag := typ.Tag.Get("default")
				if tag != "" {
					err = set(&val, tag)
					if err != nil {
						err = fmt.Errorf("default %s: %s", typ.Name, err)
						return err
					}
				}
			}
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
		if isZero(*v) {
			v.Set(reflect.ValueOf(v.Interface()))
		}
	default:
		// TODO: Do something with unknown value types
	}

	return nil
}
