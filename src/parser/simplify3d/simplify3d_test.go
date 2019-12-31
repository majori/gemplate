package simplify3d

import (
	"reflect"
	"testing"
)

func TestParseStaticSettings(t *testing.T) {
	sample := `
		; Settings Summary
		;   string,value
		;   stringlist,value1,value2
		;   int,1
		;   intlist,1,1
		; 	float,1.1
		;   floatlist,0.5,0.5
		;   empty,
		;   gcode,M400 ; wait for moves to finish,M104 S0 T0 ; turn off back extruder,M104 S0 T1 ...
	`

	static := *parseSettings(&sample)

	testValue := func(key string, kind reflect.Kind) {
		if reflect.TypeOf(static[key]).Kind() != kind {
			t.Errorf("%s value parsing failed, was type of %s instead", kind.String(), reflect.TypeOf(static["string"]))
		}
	}

	testList := func(key string, kind reflect.Kind) {
		if list, ok := static[key].([]interface{}); ok {
			for _, item := range list {
				if reflect.TypeOf(item).Kind() != kind {
					t.Errorf("%s list value parsing failed, item was type of %s instead", kind.String(), reflect.TypeOf(item))
				}
			}
		} else {
			t.Errorf("List parsing failed, got type of %s instead", reflect.TypeOf(static[key]))
		}
	}

	testValue("string", reflect.String)
	testList("stringlist", reflect.String)

	testValue("int", reflect.Int)
	testList("intlist", reflect.Int)

	testValue("float", reflect.Float64)
	testList("floatlist", reflect.Float64)

	testValue("gcode", reflect.String)

	if static["empty"] != nil {
		t.Errorf("Key \"empty\" should not exist")
	}
}
