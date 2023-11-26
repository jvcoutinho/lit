package bind

import (
	"log"

	"github.com/jvcoutinho/lit/validate"
)

func validateFields[T any](target *T) error {
	validatable, ok := any(target).(validate.Validatable)
	if !ok { // we want pointer receiver to implement the interface
		return nil
	}

	_, ok = any(*target).(validate.Validatable)
	if ok { // but we don't want value receiver to implement the interface
		log.Printf("%T: the receiver of Validate() should be a pointer in order to use "+
			"validation from bind functions", *target)
		return nil
	}

	return validate.Fields(target, validatable.Validate()...)
}
