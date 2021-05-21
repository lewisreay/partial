package partial

import (
	"errors"
	"fmt"
	"reflect"
)

// structField defines the structure of a struct field. name = field name, tag = tag value
type structField struct {
	name  string
	tag   string
	index int
}

// Partials interface allows for custom types.
// If the value type isn't of Go's basic types, implementing the Partials interface is required otherwise an error will be returned.
type Partials interface {
	// Parameter is the reflect.Value as an interface and the returning arg is the value, asserted based on the interface implementation.
	Value(i interface{}) (interface{}, error)
}

// getFieldsWithTag will get all the fields from a struct where the tag is present.
func getFieldsWithTag(tag string, t reflect.Type) ([]structField, error) {
	amt := t.NumField()
	// We don't know how many fields have the requested tag. Size must be zero.
	fields := make([]structField, 0)
	// Add fields that have the tag into a list. Fields that don't have this tag will be ignored.
	for i := 0; i < amt; i++ {
		field := t.Field(i)
		// Only add fields where the requested tag is present.
		v, found := field.Tag.Lookup(tag)
		if found {
			fields = append(fields, structField{
				name:  field.Name,
				tag:   v,
				index: i,
			})
		}
	}
	return fields, nil
}

// Get will return all of the fields with the tag and that don't have a zero value.
func Get(i interface{}, tag string) (map[string]interface{}, error) {
	if reflect.ValueOf(i).Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected type struct got %T", i)
	}

	t := reflect.TypeOf(i)
	if t == nil {
		return nil, errors.New("interface is nil")
	}

	// Get all fields with the matching tag.
	fields, err := getFieldsWithTag(tag, t)
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{}, len(fields))

	// Loop through fields that have the tag and not a zero value, create map with values. K = field name, V = value of field.
	for _, field := range fields {

		_, found := values[field.tag]
		if found {
			return nil, errors.New("cannot have duplicate key")
		}

		v := reflect.ValueOf(i).Field(field.index)

		// We don't want fields that don't have a value.
		if v.IsZero() {
			continue
		}

		switch v.Kind() {
		case reflect.Bool:
			values[field.tag] = v.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values[field.tag] = v.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			values[field.tag] = v.Uint()
		case reflect.Float32, reflect.Float64:
			values[field.tag] = v.Float()
		case reflect.Complex64, reflect.Complex128:
			values[field.tag] = v.Complex()
		case reflect.String:
			values[field.tag] = v.String()
		default:
			// Check if interface can be used.
			if !v.CanInterface() {
				return nil, errors.New("unable to determine kind and cannot interface on value")
			}
			// Check if the value implements the Partials interface.
			p, ok := v.Interface().(Partials)
			if !ok {
				return nil, fmt.Errorf("%v does not implement the Partials interface", v.Type().Name())
			}
			// Run the interface implementation and set the value.
			iv, err := p.Value(v.Interface())
			if err != nil {
				return nil, err
			}
			values[field.tag] = iv
		}
	}

	return values, nil
}
