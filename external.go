package qp

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
)

// in expands slice values in args, returning the modified query string
// and a new arg list that can be executed by a database. The `query` should
// use the `?` bindVar. The return value uses the `?` bindVar.
func in(query string, args ...interface{}) (string, []interface{}, error) {
	// argMeta stores reflect.Value and length for slices and
	// the value itself for non-slice arguments
	type argMeta struct {
		v      reflect.Value
		i      interface{}
		length int
	}

	var (
		flatArgsCount int
		anySlices     bool
		meta          = make([]argMeta, len(args))
	)

	for i, arg := range args {
		if a, ok := arg.(driver.Valuer); ok {
			arg, _ = a.Value()
		}
		v := reflect.ValueOf(arg)
		t := deref(v.Type())

		// []byte is a driver.Value type, so it should not be expanded
		if t.Kind() == reflect.Slice && t != reflect.TypeOf([]byte{}) {
			meta[i] = argMeta{v: v, length: v.Len()}
			anySlices = true
			flatArgsCount += meta[i].length

			if meta[i].length == 0 {
				return "", nil, errors.New("empty slice passed to 'in' query")
			}
		} else {
			meta[i] = argMeta{i: arg}
			flatArgsCount++
		}
	}

	// Don't do any parsing if there aren't any slices; note that this means
	// some errors that we might have caught below will not be returned.
	if !anySlices {
		return query, args, nil
	}

	newArgs := make([]interface{}, 0, flatArgsCount)
	buf := make([]byte, 0, len(query)+len(", ?")*flatArgsCount)

	var arg, offset int

	for i := strings.IndexByte(query[offset:], '?'); i != -1; i = strings.IndexByte(query[offset:], '?') {
		if arg >= len(meta) {
			return "", nil, errors.New("number of bindVars exceeds arguments")
		}

		argMeta := meta[arg]
		arg++

		// Not a slice, continue.
		if argMeta.length == 0 {
			offset = offset + i + 1
			newArgs = append(newArgs, argMeta.i)
			continue
		}

		// Write everything up to and including our ? character
		buf = append(buf, query[:offset+i+1]...)

		for si := 1; si < argMeta.length; si++ {
			buf = append(buf, ", ?"...)
		}

		newArgs = appendReflectSlice(newArgs, argMeta.v, argMeta.length)

		// Slice the query and reset the offset.
		query = query[offset+i+1:]
		offset = 0
	}

	buf = append(buf, query...)

	if arg < len(meta) {
		return "", nil, errors.New("number of bindVars less than number of arguments")
	}

	return string(buf), newArgs, nil
}

// appendReflectSlice appends elements of a reflect.Value slice to args.
func appendReflectSlice(args []interface{}, v reflect.Value, vlen int) []interface{} {
	switch val := v.Interface().(type) {
	case []interface{}:
		args = append(args, val...)
	case []int:
		for i := range val {
			args = append(args, val[i])
		}
	case []string:
		for i := range val {
			args = append(args, val[i])
		}
	default:
		for si := 0; si < vlen; si++ {
			args = append(args, v.Index(si).Interface())
		}
	}
	return args
}

// deref is Indirect for reflect.Types.
func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}
