package defaultprops

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrNoPointer is thrown when one of the structs is not a pointer
	ErrNoPointer = errors.New("the struct sent is not a pointer")
	// ErrDifferentType is thrown when of the structs have different type
	ErrDifferentType = errors.New("the structs are not equal")
)

// Config is the config for the package
// It allows to handle special cases, such as:
// - boolean substitution
type Config struct {
	SetFalseBools  bool
	SetEmptyString bool
}

// SubstituteNonConfig replaces the values of the origin into destination
// In order to do that:
// - The two values must be pointers
// - The two values must have the same type
// NOTE: if the orig value is a zero value, it's not substituted
func SubstituteNonConfig(orig, dest interface{}) error {
	if reflect.ValueOf(orig).Kind() != reflect.Ptr ||
		reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return ErrNoPointer
	}

	origElem := reflect.ValueOf(orig).Elem()
	destElem := reflect.ValueOf(dest).Elem()

	err := substitute(origElem, destElem, &Config{})

	return err
}

func mergeMaps(orig, dest reflect.Value, config *Config) error {
	if orig.Kind() != reflect.Map || dest.Kind() != reflect.Map {
		return fmt.Errorf(
			"types: %s and %s. %w",
			orig.Kind(),
			dest.Kind(),
			ErrDifferentType,
		)
	}
	for _, key := range orig.MapKeys() {
		origVal := orig.MapIndex(key)
		dest.SetMapIndex(key, origVal)
	}
	return nil
}

func substituteStruct(orig, dest reflect.Value, config *Config) error {
	for i := 0; i < orig.NumField(); i++ {
		valOrig := orig.Field(i)
		valDest := dest.Field(i)
		err := substitute(valOrig, valDest, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func substitute(orig, dest reflect.Value, config *Config) error {
	if orig.Kind() != dest.Kind() {
		return fmt.Errorf(
			"types: %s and %s. %w",
			orig.Kind(),
			dest.Kind(),
			ErrDifferentType,
		)
	}
	switch orig.Kind() {
	// simple types
	case reflect.String:
		if !orig.IsZero() {
			dest.SetString(orig.String())
		}
	case reflect.Float32, reflect.Float64:
		if !orig.IsZero() {
			dest.SetFloat(orig.Float())
		}
	case reflect.Bool:
		if !orig.IsZero() {
			dest.SetBool(orig.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !orig.IsZero() {
			dest.SetInt(orig.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !orig.IsZero() {
			dest.SetUint(orig.Uint())
		}
	case reflect.Chan:
		if !orig.IsZero() {
			dest.Set(orig)
		}
	case reflect.Slice:
		if !orig.IsZero() && !orig.IsNil() && orig.Len() > 0 {
			dest.Set(orig)
		}
	case reflect.Array:
		if !orig.IsZero() && orig.Len() > 0 {
			dest.Set(orig)
		}
	case reflect.Map:
		if !orig.IsZero() && !orig.IsNil() && orig.Len() > 0 {
			err := mergeMaps(orig, dest, config)
			if err != nil {
				return err
			}
		}
	case reflect.Func:
	case reflect.Struct:
		err := substituteStruct(orig, dest, config)
		if err != nil {
			return err
		}
	case reflect.Ptr:
		if !orig.IsNil() && dest.IsNil() {
			dest.Set(orig)
		} else if !orig.IsNil() && !dest.IsNil() {
			err := substitute(orig.Elem(), dest.Elem(), config)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%v", orig.Kind())
	}
	return nil
}
