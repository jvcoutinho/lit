package bind_test

import (
	"github.com/jvcoutinho/lit"
	"net/http"
	"net/http/httptest"
	"time"
)

type bindableTypes struct {
	String     string     `uri:"string" query:"string" header:"string"`
	Uint       uint       `uri:"uint" query:"uint" header:"uint"`
	Uint8      uint8      `uri:"uint8" query:"uint8" header:"uint8"`
	Uint16     uint16     `uri:"uint16" query:"uint16" header:"uint16"`
	Uint32     uint32     `uri:"uint32" query:"uint32" header:"uint32"`
	Uint64     uint64     `uri:"uint64" query:"uint64" header:"uint64"`
	Int        int        `uri:"int" query:"int" header:"int"`
	Int8       int8       `uri:"int8" query:"int8" header:"int8"`
	Int16      int16      `uri:"int16" query:"int16" header:"int16"`
	Int32      int32      `uri:"int32" query:"int32" header:"int32"`
	Int64      int64      `uri:"int64" query:"int64" header:"int64"`
	Float32    float32    `uri:"float32" query:"float32" header:"float32"`
	Float64    float64    `uri:"float64" query:"float64" header:"float64"`
	Complex64  complex64  `uri:"complex64" query:"complex64" header:"complex64"`
	Complex128 complex128 `uri:"complex128" query:"complex128" header:"complex128"`
	Bool       bool       `uri:"bool" query:"bool" header:"bool"`
	Time       time.Time  `uri:"time" query:"time" header:"time"`
	Slice      []int      `uri:"slice" query:"slice" header:"slice"`
	Array      [2]int     `uri:"array" query:"array" header:"array"`
}

var request = lit.NewRequest(
	httptest.NewRequest(http.MethodPost, "/users/user_1/books/book_1", nil),
	map[string]string{"user_id": "123", "book_id": "book_1"},
)

//func testBinding[T any, V bindableSingleValueTypes | bindableMultipleValueTypes](
//	t *testing.T,
//	r *lit.Request,
//	bind func(*lit.Request) (T, error),
//) {
//	t.Helper()
//
//	tests := []struct{
//		description string
//		request string
//	}
//}
