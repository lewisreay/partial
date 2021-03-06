package partial

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTag = "testingTag"
)

type sithLord struct {
	name         string       `testingTag:"name"`
	age          int          `testingTag:"age"`
	darkSide     bool         `testingTag:"dark_side"`
	Apprentices  apprentices  `testingTag:"apprentices"`
	DefeatedJedi defeatedJedi `testingTag:"defeated_jedi"`
}

type apprentices []string

type defeatedJedi []string

func (dj defeatedJedi) Value(i interface{}) (interface{}, error) {
	dj, ok := i.(defeatedJedi)
	if !ok {
		return nil, fmt.Errorf("got %T expected defeatedJedi", i)
	}
	return dj, nil
}

func TestFieldHasTag(t *testing.T) {
	r := reflect.TypeOf(sithLord{
		name:     "Count Dooku",
		age:      83,
		darkSide: true,
	})
	field, ok := r.FieldByName("darkSide")
	assert.True(t, ok, "should not be false")
	tag, ok := field.Tag.Lookup(testTag)
	assert.True(t, ok, "should not be false")
	assert.Equal(t, "dark_side", tag, "tag value should match")
}

func TestGetFieldsWithTag(t *testing.T) {
	fields, err := getFieldsWithTag(testTag, reflect.TypeOf(sithLord{
		name:         "Asajj Ventress",
		age:          32,
		darkSide:     true,
		Apprentices:  []string{"Savage Opress"},
		DefeatedJedi: []string{"None!"},
	}))
	assert.NoError(t, err, "should not error")
	assert.Len(t, fields, 5, "should have 5 fields")
}

func TestGet(t *testing.T) {
	_, err := Get([]string{"bad"}, testTag)
	assert.Error(t, err, "should error")

	sith := sithLord{
		name:     "Darth Vader",
		age:      32,
		darkSide: true,
	}
	vals, err := Get(sith, testTag)
	assert.NoError(t, err, "should not error")

	for k, v := range vals {
		switch k {
		case "name":
			assert.EqualValues(t, sith.name, v, "should match")
		case "age":
			assert.EqualValues(t, sith.age, v, "should match")
		case "dark_side":
			assert.EqualValues(t, sith.darkSide, v, "should match")
		}
	}

	sith.DefeatedJedi = []string{"Sifo Dyas"}
	vals, err = Get(sith, testTag)
	assert.NoError(t, err, "should not error")

	sith.Apprentices = []string{"Asajj Ventress"}
	vals, err = Get(sith, testTag)
	assert.Error(t, err, "should error")
}
