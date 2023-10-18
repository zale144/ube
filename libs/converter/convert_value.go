package converter

/* ------------------------------- Imports --------------------------- */

import (
	"fmt"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

/* -------------------------- Methods/Functions ---------------------- */

/*
ConvertBool makes a bool from a string
*/
func ConvertBool(field, value string) (objValue bool, err error) {
	objValue, err = strconv.ParseBool(value)
	if err != nil {
		err = fmt.Errorf("%s: value '%s' is not bool", field, value)
	}

	return objValue, err
}

/*
ConvertInt makes an int from a string
*/
func ConvertInt(field, value string) (objValue int, err error) {
	var tmpValue int64

	// make it a zero value on zero value input
	if value == "" {
		objValue = 0
		return
	}

	tmpValue, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		err = fmt.Errorf("%s: value '%s' is not int", field, value)
		return
	}

	objValue = int(tmpValue)

	return objValue, err
}

/*
ConvertFloat64 makes a float64 from a string
*/
func ConvertFloat64(field, value string) (objValue float64, err error) {
	// make it a zero value on zero value input
	if value == "" {
		objValue = 0
		return
	}

	objValue, err = strconv.ParseFloat(value, 64)
	if err != nil {
		err = fmt.Errorf("%s: value '%s' is not float64", field, value)
		return
	}

	return objValue, err
}

/*
ConvertDateTime makes a time.Time from a string
*/
func ConvertDateTime(field, value string) (objValuePtr *time.Time, err error) {
	// make it a zero value on zero value input
	if value == "" {
		return
	}

	var objValue time.Time

	objValue, err = dateparse.ParseAny(value)
	if err != nil {
		err = fmt.Errorf("%s: value '%s' is not date-time: %w", field, value, err)
		return
	}

	objValuePtr = &objValue

	return objValuePtr, err
}
