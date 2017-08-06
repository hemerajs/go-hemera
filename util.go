package hemera

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

func CleanPattern(p interface{}) interface{} {
	s := structs.New(p)
	var pattern = make(map[string]interface{})

	// pattern contains only primitive values
	// and no meta, delegate informations
	for _, f := range s.Fields() {
		fn := f.Name()

		if !strings.HasSuffix(fn, "_") {
			fk := f.Kind()

			switch fk {
			case reflect.Struct:
			case reflect.Map:
			case reflect.Array:
			case reflect.Func:
			case reflect.Chan:
			case reflect.Slice:
			default:
				pattern[f.Name()] = f.Value()
			}
		}
	}

	return pattern
}

// Dissect the cb Handler's signature
func ArgInfo(cb Handler) ([]reflect.Type, int) {
	cbType := reflect.TypeOf(cb)

	if cbType.Kind() != reflect.Func {
		panic("hemera: Handler needs to be a func")
	}

	numArgs := cbType.NumIn()
	argTypes := []reflect.Type{}

	for i := 0; i < numArgs; i++ {
		argTypes = append(argTypes, cbType.In(i))
	}

	return argTypes, numArgs
}
