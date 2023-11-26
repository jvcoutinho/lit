package bind_test

import (
	"time"

	"github.com/jvcoutinho/lit/validate"

	"github.com/jvcoutinho/lit"
)

type unbindableField struct {
	Request lit.Request `uri:"field" query:"field" header:"field"`
}

type ignorableFields struct {
	unexported string `uri:"unexported" query:"unexported" header:"unexported"`
	Missing    string `uri:"missing" query:"missing" header:"missing"`
	Untagged   string
}

type bindableFields struct {
	String     string     `uri:"string" query:"string" header:"string" json:"string" form:"string"`
	Pointer    *int       `uri:"pointer" query:"pointer" header:"pointer" json:"pointer" form:"pointer"`
	Uint       uint       `uri:"uint" query:"uint" header:"uint" json:"uint" form:"uint"`
	Uint8      uint8      `uri:"uint8" query:"uint8" header:"uint8" json:"uint8" form:"uint8"`
	Uint16     uint16     `uri:"uint16" query:"uint16" header:"uint16" json:"uint16" form:"uint16"`
	Uint32     uint32     `uri:"uint32" query:"uint32" header:"uint32" json:"uint32" form:"uint32"`
	Uint64     uint64     `uri:"uint64" query:"uint64" header:"uint64" json:"uint64" form:"uint64"`
	Int        int        `uri:"int" query:"int" header:"int" json:"int" form:"int"`
	Int8       int8       `uri:"int8" query:"int8" header:"int8" json:"int8" form:"int8"`
	Int16      int16      `uri:"int16" query:"int16" header:"int16" json:"int16" form:"int16"`
	Int32      int32      `uri:"int32" query:"int32" header:"int32" json:"int32" form:"int32"`
	Int64      int64      `uri:"int64" query:"int64" header:"int64" json:"int64" form:"int64"`
	Float32    float32    `uri:"float32" query:"float32" header:"float32" json:"float32" form:"float32"`
	Float64    float64    `uri:"float64" query:"float64" header:"float64" json:"float64" form:"float64"`
	Complex64  complex64  `uri:"complex64" query:"complex64" header:"complex64" json:"complex64" form:"complex64"`
	Complex128 complex128 `uri:"complex128" query:"complex128" header:"complex128" json:"complex128" form:"complex128"`
	Bool       bool       `uri:"bool" query:"bool" header:"bool" json:"bool" form:"bool"`
	Time       time.Time  `uri:"time" query:"time" header:"time" json:"time" form:"time"`
	Slice      []int      `uri:"slice" query:"slice" header:"slice" json:"slice" form:"slice"`
	Array      [2]int     `uri:"array" query:"array" header:"array" json:"array" form:"array"`
}

type nonPointerReceiverValidatableFields struct {
	String string `uri:"string" query:"string" header:"string" json:"string" form:"string"`
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
	String string `uri:"string" query:"string" header:"string" json:"string" form:"string"`
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
