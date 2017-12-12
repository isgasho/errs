package errs

import (
	"fmt"
	"strings"
	"testing"
)

func TestErrs(t *testing.T) {
	assert := func(t *testing.T, v bool, err ...interface{}) {
		if !v {
			t.Fatal(err...)
		}
	}

	var (
		foo = Class("foo")
		bar = Class("bar")
		baz = Class("baz")
	)

	t.Run("Class", func(t *testing.T) {
		t.Run("Has", func(t *testing.T) {
			assert(t, foo.Has(foo.New("t")))
			assert(t, !foo.Has(bar.New("t")))
			assert(t, !foo.Has(baz.New("t")))

			assert(t, !bar.Has(foo.New("t")))
			assert(t, bar.Has(bar.New("t")))
			assert(t, !bar.Has(baz.New("t")))

			assert(t, foo.Has(bar.Wrap(foo.New("t"))))
			assert(t, bar.Has(bar.Wrap(foo.New("t"))))
			assert(t, !baz.Has(bar.Wrap(foo.New("t"))))

			assert(t, foo.Has(foo.Wrap(bar.New("t"))))
			assert(t, bar.Has(foo.Wrap(bar.New("t"))))
			assert(t, !baz.Has(foo.Wrap(bar.New("t"))))
		})

		t.Run("Same Name", func(t *testing.T) {
			c1 := Class("c")
			c2 := Class("c")

			assert(t, c1.Has(c1.New("t")))
			assert(t, !c2.Has(c1.New("t")))

			assert(t, !c1.Has(c2.New("t")))
			assert(t, c2.Has(c2.New("t")))
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Format Contains Classes", func(t *testing.T) {
			assert(t, strings.Contains(foo.New("t").Error(), "foo"))
			assert(t, strings.Contains(bar.New("t").Error(), "bar"))

			assert(t, strings.Contains(bar.Wrap(foo.New("t")).Error(), "foo"))
			assert(t, strings.Contains(bar.Wrap(foo.New("t")).Error(), "bar"))

			assert(t, strings.Contains(foo.Wrap(bar.New("t")).Error(), "foo"))
			assert(t, strings.Contains(foo.Wrap(bar.New("t")).Error(), "bar"))
		})

		t.Run("Format With Stack", func(t *testing.T) {
			err := foo.New("t")

			assert(t,
				!strings.Contains(fmt.Sprintf("%v", err), "\n"),
				"%v format contains newline",
			)
			assert(t,
				strings.Contains(fmt.Sprintf("%+v", err), "\n"),
				"%+v format does not contain newline",
			)
		})

		t.Run("Format Nil", func(t *testing.T) {
			var err *Error
			assert(t, fmt.Sprintf("%v", err) == "<nil>")
		})

		t.Run("Unwrap", func(t *testing.T) {
			err := fmt.Errorf("t")

			assert(t, nil == Unwrap(nil))
			assert(t, err == Unwrap(err))
			assert(t, err == Unwrap(foo.Wrap(err)))
			assert(t, err == Unwrap(bar.Wrap(foo.Wrap(err))))
		})
	})
}