// Package confutil provides the utility functions for the configuration
package confutil

import (
	"reflect"
	"strconv"
	"strings"
)

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
		case reflect.Bool:
			if !field.Bool() {
				if defaultValue, err := strconv.ParseBool(defaultTag); err == nil {
					field.SetBool(defaultValue)
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
				field.Set(reflect.ValueOf(arr))
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
