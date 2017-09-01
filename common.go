package configo

import (
	"reflect"
	"strconv"
)

type Configo interface {
	Load(interface{}) error
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
