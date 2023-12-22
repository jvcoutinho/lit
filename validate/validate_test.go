package validate_test

type validatableFields struct {
	String  string
	Bool    bool `query:""`
	Int     int  `json:"int"`
	Pointer *int `uri:"pointer"`
}

func pointerOf[T any](value T) *T {
	return &value
}
