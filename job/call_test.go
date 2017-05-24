package job

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBind(t *testing.T) {
	global := NewFuncs(1)
	e := global.Bind("dummy", dummy)
	assert.NoError(t, e, "Error in binding function")
	var fakefunction int
	fakefunction = 10
	e = global.Bind("fake", fakefunction)
	assert.Error(t, e, "Excepting a binding error not a function ")
}

func TestCall(t *testing.T) {
	global := NewFuncs(1)
	e := global.Bind("dummy", dummy)
	assert.NoError(t, e, "Error in binding function")
	e = global.Bind("testsprintf", testsprintf)
	assert.NoError(t, e, "Error in binding function")
	values, err := global.Call("dummy")
	assert.NoError(t, err)
	assert.Equal(t, values[0].Interface().(int), 10)
	values, err = global.Call("testsprintf", "Testing")
	assert.NoError(t, err)
	assert.Equal(t, values[0].Interface().(string), "Hi Will from Testing")
	_, err = global.Call("mumbojumbo", 10)
	assert.Error(t, err, "Function expected to not exist")
	_, err = global.Call("dummy", []interface{}{10, "string"})
	assert.Error(t, err, "Function expected to not exist")
	_, err = global.Call("testsprintf", 10)
	assert.Error(t, err, "Function parameter is not the expecting type")
}
func testsprintf(pass string) string {
	return fmt.Sprintf("Hi Will from %s", pass)
}
func dummy() int {
	return 10
}
