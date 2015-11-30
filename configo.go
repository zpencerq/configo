// Package configo provides functions to handle configuration items.
//
// Example:
//  type Config struct {
//    SomeThing SomeThing
//    Mysql  MysqlConfig `toml:"mysql"`
//    Port Int `default:"8080"`
//  }
//
//  type MysqlConfig struct {
//    Dsn string
//  }
//
//  config := Config{
//    Mysql: MysqlConfig{
//        Dsn: "/mydb",
//    },
//  }
//
//  err := configo.UnmarshalFile("/path/to/config.toml", &config)
//  // error check
//
//  fmt.Printf("Listening on port %d", *config.Port)
package configo

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
)

// UnmarshalFile decodes the contents of the file `f` in TOML format into a pointer `v`. If `v` contains data, that data will be used as "defaults".
//
// A field's value will be determined based on the following order:
//
// 1. If an "env" tag exists for a field and an environment variable matching the tag's value exists, the environment variable's value will be used, subject to type casting.
// 2. If the field exists in the file, its value will be used. The `toml` tag may be used to map TOML keys to fields that don't match the key name exactly.
// 3. If `v` already contains a value for the field, it will be used.
// 4. If a "default" tag exists for a field, its value will be used, subject to type casting.
// 5. The field will be initialized to its zero value (i.e., "" for string, 0 for int, etc).
func UnmarshalFile(f string, v interface{}) error {
	var err error

	rv := reflect.ValueOf(v).Elem()

	err = setDefaults(&rv)
	if err != nil {
		return err
	}

	// Unmarshalling TOML onto a non-zero struct is inconsistent.
	// One time the value might be the pre-existing value, another time
	// it might be from the TOML. Instead we unmarshal onto a new struct
	// then walk the struct copying non-zero values.

	nv := reflect.New(rv.Type())
	ni := nv.Interface()
	_, err = toml.DecodeFile(f, ni)
	if err != nil {
		return err
	}

	err = setToml(&rv, nv.Elem())
	if err != nil {
		return err
	}

	err = setEnv(&rv)
	if err != nil {
		return err
	}

	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	case reflect.Ptr:
		if !v.Elem().IsValid() {
			return true
		}

		if v.Elem().Interface() == nil {
			return true
		}
	}

	return false
}

func set(v *reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		v.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s == "" {
			s = "0"
		}

		n, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}

		v.SetInt(n)
		return nil
	case reflect.Ptr:
		if !v.Elem().IsValid() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		switch v.Interface().(type) {
		case *bool:
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}

			v.Elem().SetBool(b)
			return nil
		case *int, *int8, *int16, *int32, *int64:
			n, err := strconv.ParseInt(s, 10, v.Elem().Type().Bits())
			if err != nil {
				return err
			}

			v.Elem().SetInt(n)
			return nil
		case *string:
			v.Elem().SetString(s)
			return nil
		default:
			// TODO: Do something with unknown pointer types
		}
	case reflect.String:
		v.SetString(s)
		return nil
	default:
		// TODO: Do something with unknown reflect types
	}

	return nil
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
