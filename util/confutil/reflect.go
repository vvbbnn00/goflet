// Package confutil provides the utility functions for the configuration
package confutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// loadArr loads the array from the string array
func loadArr(field reflect.Value, arr []string) {
	switch field.Type().Elem().Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(arr))
	case reflect.Int:
		intArr := make([]int, len(arr))
		for i, v := range arr {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				// Handle the error if the conversion fails
				fmt.Println("Error converting string to int:", err)
				return
			}
			intArr[i] = intValue
		}
		field.Set(reflect.ValueOf(intArr))
	default:
		panic("unsupported type:" + field.Type().Elem().Kind().String())
	}
}

// SetDefaults sets the default values for the fields of the configuration
func SetDefaults(config interface{}) {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	// Set default values for the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		defaultTag := fieldType.Tag.Get("default")

		if defaultTag == "" {
			continue
		}

		//nolint:exhaustive
		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.SetString(defaultTag)
			}
		case reflect.Int, reflect.Int64:
			if field.Int() == 0 {
				if defaultValue, err := strconv.ParseInt(defaultTag, 10, 64); err == nil {
					field.SetInt(defaultValue)
				}
			}
		case reflect.Bool, reflect.Ptr:
			// Should have a better way to parse bool
			if field.IsNil() {
				if defaultValue, err := strconv.ParseBool(defaultTag); err == nil {
					field.Set(reflect.ValueOf(&defaultValue))
				}
			}
		case reflect.Float64:
			if field.Float() == 0 {
				if defaultValue, err := strconv.ParseFloat(defaultTag, 64); err == nil {
					field.SetFloat(defaultValue)
				}
			}
		case reflect.Array, reflect.Slice:
			// Split the default value by comma
			if field.Len() == 0 {
				s := reflect.ValueOf(defaultTag)
				arr := strings.Split(s.String(), ",")
				loadArr(field, arr)
			}
		default:
			panic("unsupported type:" + field.Kind().String())
		}
	}

	// Set default values for the fields of the nested structs
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			SetDefaults(field.Addr().Interface())
		}
	}
}

// FormatStruct prints the fields of the struct
func FormatStruct(s interface{}) string {
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	var result string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		var fieldValue string

		//nolint:exhaustive
		switch field.Kind() {
		case reflect.Struct:
			fieldValue = FormatStruct(field.Addr().Interface())
		case reflect.Ptr:
			if !field.IsNil() {
				fieldValue = fmt.Sprintf("%v", field.Elem())
			} else {
				fieldValue = "<nil>"
			}
		default:
			fieldValue = fmt.Sprintf("%v", field)
		}

		// If the field value contains newline, add indentation
		if strings.Contains(fieldValue, "\n") {
			fieldValue = strings.ReplaceAll(fieldValue, "\n", "\n  ")
			fieldValue = "{\n  " + fieldValue + "\b\b}"
		}

		result += fmt.Sprintf("%s: %s\n", fieldName, fieldValue)
	}
	return result
}
