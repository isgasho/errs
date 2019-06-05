// Package errs provides a simple error package with stack traces.
package errs

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

// Namer is implemented by all errors returned in this package. It returns a
// name for the tag of the error, and a boolean indicating if the name is
// valid.
type Namer interface{ Name() (string, bool) }

// Errorf does the same thing as fmt.Errorf(...) except it captures a stack
// trace on creation.
func Errorf(format string, args ...interface{}) error {
	return Tag("").wrap(fmt.Errorf(format, args...))
}

// Wrap returns an error not contained in any class. It just associates a stack
// trace with the error. Wrap returns nil if err is nil.
func Wrap(err error) error {
	return Tag("").wrap(err)
}

// Tagged is a shorthand for Tag(tag).Wrap(err).
func Tagged(tag string, err error) error {
	return Tag(tag).wrap(err)
}

// Tags returns all the tags that have wrapped the error.
func Tags(err error) (tags []Tag) {
	for {
		e, ok := err.(*errorT)
		if !ok {
			return tags
		}
		if e.tag != "" {
			tags = append(tags, e.tag)
		}
		err = errors.Unwrap(err)
	}
}

//
// error tags
//

// Tag represents some extra information about an error.
type Tag string

// New constructs an error with the format string that will be contained by
// this class. This is the same as calling Wrap(fmt.Errorf(...)).
func (t Tag) Errorf(format string, args ...interface{}) error {
	return t.wrap(fmt.Errorf(format, args...))
}

// Wrap returns a new error based on the passed in error that is contained in
// this class. Wrap returns nil if err is nil.
func (t Tag) Wrap(err error) error {
	return t.wrap(err)
}

// Error returns the class string as the error text. It allows the use of
// errors.Is, or as just an easy way to have a string constant error.
func (t Tag) Error() string { return string(t) }

// create constructs the error, or just adds the class to the error, keeping
// track of the stack if it needs to construct it.
func (t Tag) wrap(err error) error {
	if err == nil {
		return nil
	}

	var pcs []uintptr
	if err, ok := err.(*errorT); ok {
		if t == "" || err.tag == t {
			return err
		}
		pcs = err.pcs
	}

	e := &errorT{
		tag: t,
		err: err,
		pcs: pcs,
	}

	if e.pcs == nil {
		e.pcs = make([]uintptr, 64)
		n := runtime.Callers(3, e.pcs)
		e.pcs = e.pcs[:n:n]
	}

	return e
}

//
// errors
//

// errorT is the type of errors returned from this package.
type errorT struct {
	tag Tag
	err error
	pcs []uintptr
}

var ( // ensure *errorT implements the helper interfaces.
	_ Namer = (*errorT)(nil)
	_ error = (*errorT)(nil)
)

// errorT implements the error interface.
func (e *errorT) Error() string {
	return fmt.Sprintf("%v", e)
}

// Format handles the formatting of the error. Using a "+" on the format string
// specifier will also write the stack trace.
func (e *errorT) Format(f fmt.State, c rune) {
	sep := ""
	if e.tag != "" {
		fmt.Fprintf(f, "%s", string(e.tag))
		sep = ": "
	}
	if text := e.err.Error(); len(text) > 0 {
		fmt.Fprintf(f, "%s%v", sep, text)
	}
	if f.Flag(int('+')) {
		summarizeStack(f, e.pcs)
	}
}

// Cause implements the interface wrapping errors are expected to implement
// to allow getting at underlying causes.
func (e *errorT) Cause() error {
	return e.err
}

// Unwrap implements the draft design for error inspection. Since this is
// on an unexported type, it should not be hard to maintain going forward
// given that it also is the exact same semantics as Cause.
func (e *errorT) Unwrap() error {
	return e.err
}

// Name returns the name for the error, which is the first wrapping tag.
func (e *errorT) Name() (string, bool) {
	return string(e.tag), e.tag != ""
}

// Is is for go1.13 errors so that the Is function reports true if the error is
// part of the class.
func (e *errorT) Is(target error) bool {
	tag, ok := target.(Tag)
	return ok && e.tag == tag
}

// summarizeStack writes stack line entries to the writer.
func summarizeStack(w io.Writer, pcs []uintptr) {
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if !more {
			return
		}
		fmt.Fprintf(w, "\n\t%s:%d", frame.Function, frame.Line)
	}
}
