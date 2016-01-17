package luar

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

type Metatable struct {
	*lua.LTable
}

func MT(L *lua.LState, value interface{}) *Metatable {
	val := reflect.ValueOf(value)
	switch val.Type().Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct:
		mt := getMetatable(L, val)
		return &Metatable{mt}
	default:
		return nil
	}
}

func (m *Metatable) Remove(name string) {
	if tbl, ok := m.RawGetString("methods").(*lua.LTable); ok {
		tbl.RawSetString(name, lua.LNil)
	}
	if tbl, ok := m.RawGetString("ptr_methods").(*lua.LTable); ok {
		tbl.RawSetString(name, lua.LNil)
	}
	if tbl, ok := m.RawGetString("fields").(*lua.LTable); ok {
		tbl.RawSetString(name, lua.LNil)
	}
}

func (m *Metatable) method(name string) lua.LValue {
	methods := m.RawGetString("methods").(*lua.LTable)
	if fn := methods.RawGetString(name); fn != lua.LNil {
		return fn
	}
	return nil
}

func (m *Metatable) ptrMethod(name string) lua.LValue {
	methods := m.RawGetString("ptr_methods").(*lua.LTable)
	if fn := methods.RawGetString(name); fn != lua.LNil {
		return fn
	}
	return nil
}

func (m *Metatable) fieldIndex(name string) []int {
	fields := m.RawGetString("fields").(*lua.LTable)
	if index := fields.RawGetString(name); index != lua.LNil {
		return index.(*lua.LUserData).Value.([]int)
	}
	return nil
}
