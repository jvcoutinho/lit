package validate_test

type validatableFields struct {
	String  string
	Bool    bool `validate:""`
	Int     int  `validate:"int"`
	Pointer *int `validate:"pointer"`
}

func pointerOf[T any](value T) *T {
	return &value
}
