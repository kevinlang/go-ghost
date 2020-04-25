package ghost

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"time"
)

// Pagination has all pagination related info for the response.
type Pagination struct {
	Page  *int
	Limit *int
	Pages *int
	Total *int
	Next  *int
	Prev  *int
}

// Meta encompasses meta data from the response we get back from the API.
type Meta struct {
	Pagination *Pagination
}

func (m Meta) String() string {
	return Stringify(m)
}

// QueryParams are query params that can be used for get and list requests.
type QueryParams struct {
	// TODO
}

// ListParams are params that can be used for list requests.
type ListParams struct {
	QueryParams
	Filter string `url:"filter,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Page   int    `url:"page,omitempty"`
	Order  string `url:"order,omitempty"`
}

func (lp ListParams) String() string {
	return Stringify(lp)
}

// String returns a pointer to the string.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the bool.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to the int.
func Int(i int) *int {
	return &i
}

// Time creates a timestamp from the RFC3339 string and returns a pointer,
// ignoring any errors that occur during construction.
func Time(s string) *time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}

// Stringify attempts to create a reasonable string representation of types in
// the Ghost library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
