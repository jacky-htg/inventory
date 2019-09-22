package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
func Decode(r *http.Request, val interface{}) error {
	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		return ErrBadRequest(err, "")
	}

	validate = validator.New()

	err = validate.Struct(val)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		for _, verr := range err.(validator.ValidationErrors) {
			err = errors.New(verr.Field() + " is " + verr.Tag())
			break
		}
	}

	if err != nil {
		return ErrBadRequest(err, err.Error())
	}

	return nil
}
