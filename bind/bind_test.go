package bind_test

import (
	"time"

	"github.com/jvcoutinho/lit"
	"github.com/jvcoutinho/lit/validate"
)

type unbindableField struct {
	Request lit.Request `header:"field" query:"field" uri:"field"`
}

type ignorableFields struct {
	unexported string `header:"unexported" query:"unexported" uri:"unexported"`
	Missing    string `header:"missing"    query:"missing"    uri:"missing"`
	Untagged   string
}

type bindableFields struct {
	String     string     `form:"string"     header:"string"     json:"string"     query:"string"     uri:"string"`
	Pointer    *int       `form:"pointer"    header:"pointer"    json:"pointer"    query:"pointer"    uri:"pointer"`
	Uint       uint       `form:"uint"       header:"uint"       json:"uint"       query:"uint"       uri:"uint"`
	Uint8      uint8      `form:"uint8"      header:"uint8"      json:"uint8"      query:"uint8"      uri:"uint8"`
	Uint16     uint16     `form:"uint16"     header:"uint16"     json:"uint16"     query:"uint16"     uri:"uint16"`
	Uint32     uint32     `form:"uint32"     header:"uint32"     json:"uint32"     query:"uint32"     uri:"uint32"`
	Uint64     uint64     `form:"uint64"     header:"uint64"     json:"uint64"     query:"uint64"     uri:"uint64"`
	Int        int        `form:"int"        header:"int"        json:"int"        query:"int"        uri:"int"`
	Int8       int8       `form:"int8"       header:"int8"       json:"int8"       query:"int8"       uri:"int8"`
	Int16      int16      `form:"int16"      header:"int16"      json:"int16"      query:"int16"      uri:"int16"`
	Int32      int32      `form:"int32"      header:"int32"      json:"int32"      query:"int32"      uri:"int32"`
	Int64      int64      `form:"int64"      header:"int64"      json:"int64"      query:"int64"      uri:"int64"`
	Float32    float32    `form:"float32"    header:"float32"    json:"float32"    query:"float32"    uri:"float32"`
	Float64    float64    `form:"float64"    header:"float64"    json:"float64"    query:"float64"    uri:"float64"`
	Complex64  complex64  `form:"complex64"  header:"complex64"  json:"complex64"  query:"complex64"  uri:"complex64"`
	Complex128 complex128 `form:"complex128" header:"complex128" json:"complex128" query:"complex128" uri:"complex128"`
	Bool       bool       `form:"bool"       header:"bool"       json:"bool"       query:"bool"       uri:"bool"`
	Time       time.Time  `form:"time"       header:"time"       json:"time"       query:"time"       uri:"time"`
	Slice      []int      `form:"slice"      header:"slice"      json:"slice"      query:"slice"      uri:"slice"`
	Array      [2]int     `form:"array"      header:"array"      json:"array"      query:"array"      uri:"array"`
}

type nonPointerReceiverValidatableFields struct {
	String string `form:"string" header:"string" json:"string" query:"string" uri:"string"`
}

func (f nonPointerReceiverValidatableFields) Validate() []validate.Field {
	return []validate.Field{
		{
			Valid:   len(f.String) > 6,
			Message: "{0} should have a length greater than 6",
			Fields:  []any{&f.String},
		},
	}
}

type pointerReceiverValidatableFields struct {
	String string `form:"string" header:"string" json:"string" query:"string" uri:"string"`
}

func (f *pointerReceiverValidatableFields) Validate() []validate.Field {
	return []validate.Field{
		{
			Valid:   len(f.String) > 6,
			Message: "{0} should have a length greater than 6",
			Fields:  []any{&f.String},
		},
	}
}

func pointerOf[T any](value T) *T {
	return &value
}
