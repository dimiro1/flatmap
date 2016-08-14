package flatmap

import (
	"reflect"
	"testing"
)

func TestFlattenWithConfig(t *testing.T) {
	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		{
			Input: map[string]interface{}{
				"array": []string{"one", "two"},
			},
			Output: map[string]interface{}{
				"array.0":      "one",
				"array.1":      "two",
				"array.length": 2,
			},
		},
	}

	for _, tc := range cases {
		result, err := FlattenWithConfig(tc.Input, Config{AddLengthForArrays: true})
		if err != nil {
			t.Fatal(err)
		}

		compare(t, tc.Input, result, tc.Output)
	}
}

func compare(t *testing.T, input, result, output map[string]interface{}) {
	for k, v := range output {
		reflectValue := reflect.ValueOf(v)

		switch reflectValue.Kind() {
		case reflect.Bool:
			checkError(t, reflectValue.Bool() == result[k].(bool), input, result, output)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			checkError(t, reflectValue.Int() == result[k].(int64), input, result, output)
		case reflect.Float64, reflect.Float32:
			checkError(t, reflectValue.Float() == result[k].(float64), input, result, output)
		case reflect.String:
			checkError(t, reflectValue.String() == result[k].(string), input, result, output)
		default:
			// Throw error by default
			checkError(t, false, input, result, output)
		}
	}
}

func checkError(t *testing.T, ok bool, input, result, output map[string]interface{}) {
	if !ok {
		t.Fatalf(
			"Input:\n\n%#v\n\nOutput:\n\n%#v\n\nExpected:\n\n%#v\n",
			input,
			result,
			output)
	}
}

func TestFlatten(t *testing.T) {
	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		{
			Input: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			Output: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
		},
		{
			Input: map[string]interface{}{
				"foo": []string{
					"one",
					"two",
				},
			},
			Output: map[string]interface{}{
				"foo.0": "one",
				"foo.1": "two",
			},
		},
		{
			Input: map[string]interface{}{
				"foo": []map[interface{}]interface{}{
					map[interface{}]interface{}{
						"name":    "bar",
						"port":    3000,
						"enabled": true,
					},
				},
			},
			Output: map[string]interface{}{
				"foo.0.name":    "bar",
				"foo.0.port":    3000,
				"foo.0.enabled": true,
			},
		},
		{
			Input: map[string]interface{}{
				"foo": []map[interface{}]interface{}{
					map[interface{}]interface{}{
						"name": "bar",
						"ports": []int{
							1,
							2,
						},
					},
				},
			},
			Output: map[string]interface{}{
				"foo.0.name":    "bar",
				"foo.0.ports.0": 1,
				"foo.0.ports.1": 2,
			},
		},
		{
			Input: map[string]interface{}{
				"foo": struct {
					Name string
					Age  int
				}{
					"astaxie",
					30,
				},
			},
			Output: map[string]interface{}{
				"foo.Name": "astaxie",
				"foo.Age":  30,
			},
		},
	}

	for _, tc := range cases {
		result, err := Flatten(tc.Input)

		if err != nil {
			t.Fatal(err)
		}

		compare(t, tc.Input, result, tc.Output)
	}
}
