package configo

import (
	"fmt"
	"os"
	"reflect"
)

type EnvConfigo struct{}

func NewEnvConfigo() *EnvConfigo { return &EnvConfigo{} }

func (c *EnvConfigo) Load(v interface{}) error {
	return FromEnv(v)
}

// FromEnv sets pointer `v` based on the environment.
//
// A field's value will be determined based on the following order:
//
// 1. If an "env" tag exists for a field and an environment variable matching the tag's value exists, the environment variable's value will be used, subject to type casting.
// 2. If `v` already contains a value for the field, it will be used.
func FromEnv(v interface{}) error {
	rv := reflect.ValueOf(v).Elem()

	return setEnv(&rv)
}

func setEnv(v *reflect.Value) error {
	var err error

	switch v.Kind() {
	case reflect.Ptr:
		nv := v.Elem()

		err = setEnv(&nv)
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
				err = setEnv(&val)
				if err != nil {
					return err
				}

				continue
			}

			if kind == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
				err = setEnv(&val)
				if err != nil {
					return err
				}

				continue
			}

			tag := typ.Tag.Get("env")
			if tag == "" {
				continue
			}

			getenv := os.Getenv(tag)
			if getenv == "" {
				continue
			}

			err = set(&val, getenv)
			if err != nil {
				err = fmt.Errorf("default %s: %s", typ.Name, err)
				return err
			}
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
		v.Set(reflect.ValueOf(v.Interface()))
	default:
		// TODO: Do something with unknown value types
	}

	return nil
}
