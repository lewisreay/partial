package partial

import (
	"errors"
	"fmt"
	"reflect"
)

// structField defines the structure of a struct field. name = field name, tag = tag value
type structField struct {
	name string
	tag  string
}

// Partials interface allows for custom types.
type Partials interface {
	Value(i interface{}) (interface{}, error)
}

// fieldHasTag will check if a field has the tag.
func fieldHasTag(name, tag string, t reflect.Type) (bool, string) {
	field, found := t.FieldByName(name)
	if !found {
		return false, ""
	}
	v, found := field.Tag.Lookup(tag)
	if !found {
		return false, ""
	}
	return true, v
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
		found, tagValue := fieldHasTag(field.Name, tag, t)
		if found {
			fields = append(fields, structField{
				name: field.Name,
				tag:  tagValue,
			})
		}
	}
	return fields, nil
}

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

		v := reflect.ValueOf(i).FieldByName(field.name)

		// We don't want fields that don't have a value.
		if v.IsZero() {
			break
		}

		switch v.Kind() {
		case reflect.Bool:
			values[field.tag] = v.Bool()
		case reflect.Int:
			values[field.tag] = v.Int()
		case reflect.Int8:
			values[field.tag] = v.Int()
		case reflect.Int16:
			values[field.tag] = v.Int()
		case reflect.Int32:
			values[field.tag] = v.Int()
		case reflect.Int64:
			values[field.tag] = v.Int()
		case reflect.Uint:
			values[field.tag] = v.Uint()
		case reflect.Uint8:
			values[field.tag] = v.Uint()
		case reflect.Uint16:
			values[field.tag] = v.Uint()
		case reflect.Uint32:
			values[field.tag] = v.Uint()
		case reflect.Uint64:
			values[field.tag] = v.Uint()
		case reflect.Float32:
			values[field.tag] = v.Float()
		case reflect.Float64:
			values[field.tag] = v.Float()
		case reflect.Complex64:
			values[field.tag] = v.Complex()
		case reflect.Complex128:
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
			value, err := p.Value(v.Interface())
			if err != nil {
				return nil, err
			}
			values[field.tag] = value
		}
	}

	return values, nil
}
