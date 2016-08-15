package flatmap

import (
	"fmt"
	"reflect"
)

// Config is holds the FlattenMap configuration
type Config struct {
	AddLengthForArrays bool
}

// FlattenWithConfig transform complex map to flatten map
func FlattenWithConfig(data map[string]interface{}, config Config) (flatmap map[string]interface{}, err error) {
	flatmap = make(map[string]interface{})
	for k, raw := range data {
		err = flatten(flatmap, k, config, reflect.ValueOf(raw))
		if err != nil {
			return nil, err
		}
	}
	return
}

// Flatten transform complex map to flatten map
func Flatten(data map[string]interface{}) (flatmap map[string]interface{}, err error) {
	return FlattenWithConfig(data, Config{false})
}

func flatten(result map[string]interface{}, prefix string, config Config, v reflect.Value) (err error) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		result[prefix] = v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result[prefix] = v.Int()
	case reflect.Float64, reflect.Float32:
		result[prefix] = v.Float()
	case reflect.Map:
		err = flattenMap(result, prefix, config, v)
		if err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		err = flattenSliceArray(result, prefix, config, v)
		if err != nil {
			return err
		}
	case reflect.Struct:
		err = flattenStruct(result, prefix, config, v)
		if err != nil {
			return err
		}
	case reflect.String:
		result[prefix] = v.String()
	case reflect.Invalid:
		result[prefix] = interface{}(nil)
	default:
		return fmt.Errorf("Unknown: %s", v)
	}
	return nil
}

func flattenMap(result map[string]interface{}, prefix string, config Config, v reflect.Value) (err error) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}
		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", prefix, k))
		}
		err = flatten(result, fmt.Sprintf("%s.%s", prefix, k.String()), config, v.MapIndex(k))
		if err != nil {
			return err
		}
	}
	return nil
}

func flattenSliceArray(result map[string]interface{}, prefix string, config Config, v reflect.Value) (err error) {
	for i := 0; i < v.Len(); i++ {
		err = flatten(result, fmt.Sprintf("%s[%d]", prefix, i), config, v.Index(i))
		if err != nil {
			return err
		}
	}

	if config.AddLengthForArrays {
		err := flatten(result, fmt.Sprintf("%s.length", prefix), config, reflect.ValueOf(v.Len()))

		if err != nil {
			return err
		}
	}

	return nil
}

func flattenStruct(result map[string]interface{}, prefix string, config Config, v reflect.Value) (err error) {
	prefix = prefix + "."
	ty := v.Type()
	for i := 0; i < ty.NumField(); i++ {
		err = flatten(result, fmt.Sprintf("%s%s", prefix, ty.Field(i).Name), config, v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}
