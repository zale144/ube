package converter

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

/*
ConvertStruct converts a struct to a new type
*/
func ConvertStruct(typ interface{}, newVal reflect.Value) error {
	if typ == nil {
		return fmt.Errorf("source type is not provided")
	}

	if !newVal.IsValid() {
		return fmt.Errorf("destination type is not provided")
	}

	if newVal.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer structs are not assignable")
	}

	if newVal.IsZero() {
		newVal.Set(reflect.New(reflect.TypeOf(newVal.Interface()).Elem()))
	}

	val := reflect.Indirect(reflect.ValueOf(typ)) // feed

	for i := 0; i < val.NumField(); i++ {
		err := convertField(val, newVal, i)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
convertField converts a field to a new type
*/
func convertField(val, newVal reflect.Value, counter int) error {
	fldType := reflect.Indirect(val).Type().Field(counter)

	newFld := newVal.Elem().FieldByName(fldType.Name)
	if !newFld.IsValid() || !newFld.CanSet() {
		return fmt.Errorf("cannot set value for field '%s'", fldType.Name)
	}

	fld := val.Field(counter)

	isDone := setSimilarKind(fld, newFld)
	if isDone {
		return nil
	}

	if err := convertType(fld, newFld); err != nil {
		return fmt.Errorf("failed to convert field '%s': %w", fldType.Name, err)
	}

	return nil
}

func setSimilarKind(fld, newFld reflect.Value) bool {
	if fld.Kind() != newFld.Kind() {
		return false
	}

	if fld.Kind() == reflect.Array || fld.Kind() == reflect.Slice {
		if fld.Type().Elem().Kind() == newFld.Type().Elem().Kind() &&
			fld.Type().Elem().String() == newFld.Type().Elem().String() {
			newFld.Set(fld)

			return true
		}
	} else if fld.Type() == newFld.Type() {
		newFld.Set(fld)

		return true
	}

	return false
}

func convertType(val, newVal reflect.Value) error { // nolint: gocognit
	if val.IsZero() {
		return nil
	}

	switch val.Kind() {
	case reflect.String:
		err := convertString(val, newVal)
		if err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		err := convertArray(val, newVal)
		if err != nil {
			return err
		}
	case reflect.Map:
		err := convertMap(val, newVal)
		if err != nil {
			return err
		}

	case reflect.Struct:
		// test it
		if err := ConvertStruct(val.Interface(), reflect.Indirect(newVal)); err != nil {
			return fmt.Errorf("failed to convert struct element: %w", err)
		}
	case reflect.Ptr:
		if err := ConvertStruct(val.Interface(), newVal); err != nil {
			return fmt.Errorf("failed to convert struct element: %w", err)
		}
	default:
		return fmt.Errorf("sorry, we don't support converting your type '%s' to type '%s'", val.Type(), newVal.Type())
	}

	return nil
}

/*
convertArray converts an array (or slice) to a new type
*/
func convertArray(val, newVal reflect.Value) error {
	for j := 0; j < val.Len(); j++ {
		element := val.Index(j)
		target := reflect.Indirect(reflect.New(newVal.Type().Elem()))

		if err := convertType(element, target); err != nil {
			return fmt.Errorf("failed to convert array element: %w", err)
		}

		newVal.Set(reflect.Append(newVal, target))
	}

	return nil
}

/*
convertMap converts a map to a new type
*/
func convertMap(val, newVal reflect.Value) error {
	keys := val.MapKeys()
	if len(keys) > 0 {
		key := keys[0]
		aMapType := reflect.MapOf(key.Type(), newVal.Type().Elem())
		newVal.Set(reflect.MakeMapWithSize(aMapType, len(keys)))
	}

	for _, key := range keys {
		element := val.MapIndex(key)
		target := reflect.Indirect(reflect.New(newVal.Type().Elem()))

		if err := convertType(element, target); err != nil {
			return fmt.Errorf("failed to convert array element: %w", err)
		}

		newVal.SetMapIndex(key, target)
	}

	return nil
}

// convertString converts a string to a new type
func convertString(val, newVal reflect.Value) error { // nolint: gocognit
	switch newVal.Type().String() {
	case "int", "int8", "int16", "int32", "int64":
		intVal, err := strconv.ParseInt(val.String(), 10, 64)
		if err != nil {
			return err
		}

		newVal.SetInt(intVal)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		intVal, err := strconv.ParseUint(val.String(), 10, 64)
		if err != nil {
			return err
		}

		newVal.SetUint(intVal)
	case "bool":
		boolVal, err := strconv.ParseBool(val.String())
		if err != nil {
			return err
		}

		newVal.SetBool(boolVal)
	case "float32", "float64":
		fltVal, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return err
		}

		newVal.SetFloat(fltVal)
	case "time.Time":
		t, err := dateparse.ParseAny(val.String())
		if err != nil {
			return err
		}

		newVal.Set(reflect.ValueOf(t))
	case "*time.Time":
		t, err := dateparse.ParseAny(val.String())
		if err != nil {
			return err
		}

		x := new(time.Time)
		reflect.ValueOf(x).Elem().Set(reflect.ValueOf(&t).Elem())
		newVal.Set(reflect.ValueOf(x))
	default:
		return fmt.Errorf("sorry, we don't support converting your string type '%s' to type '%s'", val.Type(), newVal.Type())
	}

	return nil
}
