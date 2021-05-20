package partial

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTag = "testingTag"
)

type sithLord struct {
	name     string `testingTag:"name"`
	age      int    `testingTag:"age"`
	darkSide bool   `testingTag:"dark_side"`
}

func TestFieldHasTag(t *testing.T) {
	ok, tag := fieldHasTag("darkSide", testTag, reflect.TypeOf(sithLord{
		name:     "Count Dooku",
		age:      83,
		darkSide: true,
	}))
	assert.True(t, ok, "should not be false")
	assert.Equal(t, "dark_side", tag, "tag value should match")
}

func TestGetFieldsWithTag(t *testing.T) {
	fields, err := getFieldsWithTag(testTag, reflect.TypeOf(sithLord{
		name:     "Asajj Ventress",
		age:      32,
		darkSide: true,
	}))
	assert.NoError(t, err, "should not error")
	assert.Len(t, fields, 3, "should have 3 fields")
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
}
