package validate

/* ------------------------------- Imports --------------------------- */

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	libValidator "github.com/go-playground/validator/v10"
)

/* ---------------------- Constants/Types/Variables ------------------ */

const (
	sliceKind  = 23
	stringText = "string"
	sliceText  = "slice"
)

// libValidate is a single instance of Validate, it caches struct info
var libValidate = libValidator.New() // *validator.Validate

var (
	// the following two default values can be overwriten by SetMaxErrors()
	locMaxFieldErrors  = 50 // maximum number of field errors
	locMaxRecordErrors = 10 // maximum number of record errors

	fieldErrorsCount = 0
)

/* -------------------------- Methods/Functions ---------------------- */

/*
Struct validates a model.

	for more info about the validator in use, check out the following page:
	https://godoc.org/gopkg.in/go-playground/validator.v10
*/
func Struct(mdl interface{}) (err error) {
	structValue := reflect.ValueOf(mdl)
	if structValue.Kind() != reflect.Ptr {
		oName := fmt.Sprintf("%T", mdl)
		return fmt.Errorf("model '%s' should be a pointer", oName)
	}

	if structValue.Elem().Kind() == reflect.Ptr {
		oName := fmt.Sprintf("%T", mdl)
		return fmt.Errorf("'%s' is a pointer to a pointer, need a pointer to a struct", oName)
	}

	// init
	fieldErrorsCount = 0

	// is it a slice?
	if structValue.Elem().Kind() == sliceKind {
		err = checkAllStructs(mdl)
	} else {
		err = checkMyStruct(mdl)
	}

	return err
}

/*
SetMaxErrors sets the maximum number of errors to be returned
*/
func SetMaxErrors(maxFieldErrors, maxRecordErrors int) {
	locMaxFieldErrors = maxFieldErrors
	locMaxRecordErrors = maxRecordErrors
}

/*
checkAllStructs checks all structs in a slice
*/
func checkAllStructs(mdl interface{}) error {
	mdlValue := reflect.ValueOf(mdl)
	mdlValue = mdlValue.Elem()

	var (
		errMessages []string
		err         error
	)

	for i := 0; i < mdlValue.Len(); i++ {
		val := mdlValue.Index(i).Interface()
		if err = checkMyStruct(val); err != nil {
			errMessages = append(errMessages, fmt.Errorf("error validating item %d: %w", i, err).Error())
			if len(errMessages) == locMaxRecordErrors {
				errMessages = append(errMessages, "too many records have validation errors")
				break
			}

			if fieldErrorsCount == locMaxFieldErrors {
				break
			}
		}
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "; "))
	}

	return nil
}

/*
checkMyStruct checks the struct and can catch a panic
*/
func checkMyStruct(mdl interface{}) (err error) {
	defer func() {
		// this error will be returned if a panic happens
		if errInt := recover(); errInt != nil {
			err = fmt.Errorf("%v", errInt)
		}
	}()

	// register function to get tag name from json tags.
	libValidate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	_ = libValidate.RegisterValidation("isDateValid", isDateValid) // register custom validator

	verrs := libValidate.Struct(mdl)
	err = analyzeErrors(verrs)

	return err
}

/*
isDateValid custom validator for date
*/
func isDateValid(fl libValidator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

/*
analyzeErrors makes a combined error from many field errors
*/
func analyzeErrors(verrs error) error {
	if verrs == nil {
		return nil
	}

	// in case of a strange validation error
	if err, ok := verrs.(*libValidator.InvalidValidationError); ok {
		return fmt.Errorf("validation failed with InvalidValidationError: %w", err)
	}

	var (
		errMessages []string
		err         error
	)

	for _, verr := range verrs.(libValidator.ValidationErrors) {
		err = analyseError(verr.Namespace(), verr)
		errMessages = append(errMessages, err.Error())
		fieldErrorsCount++

		if fieldErrorsCount == locMaxFieldErrors {
			errMessages = append(errMessages, "too many fields have validation errors")
			break
		}
	}

	if errMessages != nil {
		err = errors.New(strings.Join(errMessages, ","))
	}

	return err
}

/*
analyseError makes an error out of a validation error
*/
func analyseError(field string, e libValidator.FieldError) error { // nolint: gocyclo // is not that complex
	structFieldType := e.Type()
	fType := structFieldType.Name()
	fieldType := e.Kind().String()

	switch e.Tag() {
	case "email":
		value := e.Value().(string)
		return fmt.Errorf("'%s=%s' (%s) should be a valid emailaddress", field, value, fType)
	case "gt":
		value := fmt.Sprintf("%v", e.Value())
		if fieldType == sliceText {
			return fmt.Errorf("'%s=%s' (%s) should have a more elements than '%s'", field, value, fieldType, e.Param())
		}

		return fmt.Errorf("'%s=%s' (%s) should have a value greater than '%s'", field, value, fType, e.Param())
	case "gte":
		value := fmt.Sprintf("%v", e.Value())
		return fmt.Errorf("'%s=%s' (%s) should have a minimal value of '%s'", field, value, fType, e.Param())
	case "isDateValid":
		value := fmt.Sprintf("%v", e.Value())
		return fmt.Errorf("'%s=%s' should have a date (YYYY-MM-DD) value", field, value)
	case "lt":
		value := fmt.Sprintf("%v", e.Value())
		if fieldType == sliceText {
			return fmt.Errorf("'%s=%s' (%s) should have a less elements than '%s'", field, value, fieldType, e.Param())
		}

		return fmt.Errorf("'%s=%s' (%s) should have a value less than '%s'", field, value, fType, e.Param())
	case "lte":
		value := fmt.Sprintf("%v", e.Value())
		return fmt.Errorf("'%s=%s' (%s) should have a maximal value of '%s'", field, value, fType, e.Param())
	case "len":
		value := fmt.Sprintf("%v", e.Value())
		return fmt.Errorf("'%s=%s' should have a length of '%s'", field, value, e.Param())
	case "max":
		value := fmt.Sprintf("%v", e.Value())
		if fType == stringText {
			return fmt.Errorf("'%s=%s' should have a maximal length of '%s'", field, value, e.Param())
		}

		return fmt.Errorf("'%s=%s' (%s) should have a maximal value of '%s'", field, value, fType, e.Param())
	case "min":
		value := fmt.Sprintf("%v", e.Value())
		if fType == stringText {
			return fmt.Errorf("'%s=%s' should have a minimal length of '%s'", field, value, e.Param())
		}

		return fmt.Errorf("'%s=%s' (%s) should have a minimal value of '%s'", field, value, fType, e.Param())
	case "numeric":
		return fmt.Errorf("'%s=%s' should have a numeric value", field, e.Value())
	case "oneof":
		return fmt.Errorf("'%s' should be one of these allowed values: '%s'", field, e.Param())
	case "required":
		return fmt.Errorf("'%s' is required", field)
	}

	return fmt.Errorf(`'%s' has an unhandled validation error '%s'`, field, e.Tag())
}
