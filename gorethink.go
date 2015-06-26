package gorethink

import (
	"reflect"

	"github.com/lamielle/gorethink/encoding"
)

func init() {
	// Set encoding package
	encoding.IgnoreType(reflect.TypeOf(Term{}))
}
