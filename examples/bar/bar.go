package bar

import (
	"github.com/mattmoor/hellogen/examples/foo"
)

// +hello:function=a long value
func Bar(f foo.Foo) (string, error) {
	return "", nil
}
