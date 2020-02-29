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
// NOTE: not used yet
type Config struct {
	SetFalseBools      bool // allows false values to be set
	SkipIfNonZeroValue bool // if the destination have non-zero value, do not replace
	ReplaceMaps        bool // if true, the maps are not merged but replaced
}

// SubstituteNonConfig replaces the values of the origin into destination
// In order to do that:
// - the two values must be pointers
// - they should have the same type
func SubstituteNonConfig(orig, dest interface{}) error {
	return Substitute(orig, dest, Config{})
}

// Substitute replaces the values of the origin into the destination
// with the given configuration. In order to work:
// - the two values must be pointers
// - they should have the same type
func Substitute(orig, dest interface{}, config Config) error {
	if reflect.ValueOf(orig).Kind() != reflect.Ptr ||
		reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return ErrNoPointer
	}

	origElem := reflect.ValueOf(orig).Elem()
	destElem := reflect.ValueOf(dest).Elem()

	err := substitute(origElem, destElem, &config)

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

func validSimpleTypes(orig, dest reflect.Value, config *Config) bool {
	skipIfDestNonZero := !(config.SkipIfNonZeroValue && !dest.IsZero())
	isNotZero := !orig.IsZero()
	if orig.Kind() == reflect.Bool {
		return isNotZero || config.SetFalseBools
	}
	return isNotZero && skipIfDestNonZero
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
		if validSimpleTypes(orig, dest, config) {
			dest.SetString(orig.String())
		}
	case reflect.Float32, reflect.Float64:
		if validSimpleTypes(orig, dest, config) {
			dest.SetFloat(orig.Float())
		}
	case reflect.Bool:
		if validSimpleTypes(orig, dest, config) {
			dest.SetBool(orig.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validSimpleTypes(orig, dest, config) {
			dest.SetInt(orig.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validSimpleTypes(orig, dest, config) {
			dest.SetUint(orig.Uint())
		}
	case reflect.Complex64, reflect.Complex128:
		if validSimpleTypes(orig, dest, config) {
			dest.SetComplex(orig.Complex())
		}
	case reflect.Chan:
		if !orig.IsZero() {
			dest.Set(orig)
		}
	case reflect.Slice:
		if !orig.IsZero() && orig.Len() > 0 {
			dest.Set(orig)
		}
	case reflect.Array:
		if !orig.IsZero() && orig.Len() > 0 {
			dest.Set(orig)
		}
	case reflect.Map:
		if config.ReplaceMaps {
			dest.Set(orig)
		} else if !orig.IsZero() && orig.Len() > 0 {
			err := mergeMaps(orig, dest, config)
			if err != nil {
				return err
			}
		}
	case reflect.Func:
		// skip for now
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
