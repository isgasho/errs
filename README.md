# errs/v2

[![GoDoc](https://godoc.org/github.com/zeebo/errs/v2?status.svg)](https://godoc.org/github.com/zeebo/errs)
[![Sourcegraph](https://sourcegraph.com/github.com/zeebo/errs/v2/-/badge.svg)](https://sourcegraph.com/github.com/zeebo/errs?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeebo/errs/v2)](https://goreportcard.com/report/github.com/zeebo/errs)

errs is a package for making errors friendly and easy.

### Creating Errors

The easiest way to use it, is to use the package level [Errorf][Errorf] function. It's much like  `fmt.Errorf`, but better. For example:

```go
func checkThing() error {
	return errs.Errorf("what's up with %q?", "zeebo")
}
```

Why is it better? Errors come with a stack trace that is only printed when a `"+"` character is used in the format string. This should retain the benefits of being able to diagnose where and why errors happen, without all of the noise of printing a stack trace in every situation. For example:

```go
func doSomeRealWork() {
	err := checkThing()
	if err != nil {
		fmt.Printf("%+v\n", err) // contains stack trace if it's a errs error.
		fmt.Printf("%v\n", err)  // does not contain a stack trace
		return
	}
}
```

### Error Tags

You can create a [Tag][Tag] for errors and check if any error has been associated with that tag. The tag is prefixed to all of the error strings it creates, and tags are just strings: two tags with the same contents are the same tag. For example:

```go
const Unauthorized = errs.Tag("unauthorized")

func checkUser(username, password string) error {
	if username != "zeebo" {
		return Unauthorized.Errorf("who is %q?", username)
	}
	if password != "hunter2" {
		return Unauthorized.Errorf("that's not a good password, jerkmo!")
	}
	return nil
}

func handleRequest() {
	if err := checkUser("zeebo", "hunter3"); errors.Is(err, Unauthorized) {
		fmt.Println(err)
		fmt.Println(errors.Is(err, Tag("unauthorized"))))
	}

	// output:
	// unauthorized: that's not a good password, jerkmo!
	// true
}
```

Tags can also [Wrap][TagWrap] other errors, and errors may be wrapped multiple times. For example:

```go
const (
	Package      = errs.Tag("mypackage")
	Unauthorized = errs.Tag("unauthorized")
)

func deep3() error {
	return fmt.Errorf("ouch")
}

func deep2() error {
	return Unauthorized.Wrap(deep3())
}

func deep1() error {
	return Package.Wrap(deep2())
}

func deep() {
	fmt.Println(deep1())

	// output:
	// mypackage: unauthorized: ouch
}
```

In the above example, both `errors.Is(deep1(), Package)` and `errors.Is(deep1()), Unauthorized)` would return `true`, and the stack trace would only be recorded once at the `deep2` call.

In addition, when an error has been wrapped, wrapping it again with the same class will not do anything. For example:

```go
func doubleWrap() {
	fmt.Println(Package.Wrap(Package.Errorf("foo")))

	// output:
	// mypackage: foo
}
```

This is to make it an easier decision if you should wrap or not (you should).

### Utilities

[Tags][Tags] is a helper function to get a slice of tags that an error has. The latest wrap is first in the slice. For example:

```go
func getTags() {
	tags := errs.Tags(deep1())
	fmt.Println(tags[0] == Package)
	fmt.Println(tags[1] == Unauthorized)

	// output:
	// true
	// true
}
```

If you don't have a tag available but don't really want to make an exported one but do want to have the error tagged for monitoring purposes, you can create a one of tag with the [Tagged][Tagged] helper:

```go
func oneOff() error {
	fh, err := fh.Open("somewhere")
	if err != nil {
		return errs.Tagged("open", err)
	}
	return errs.Tagged("close", fh.Close())
}
```

### Groups

[Groups][Group] allow one to collect a set of errors. For example:

```go
func tonsOfErrors() error {
	var group errs.Group
	for _, work := range someWork {
		group.Add(maybeErrors(work))
	}
	return group.Err()
}
```

Some things to note:

- The [Add][GroupAdd] method only adds to the group if the passed in error is non-nil.
- The [Err][GroupErr] method returns an error only if non-nil errors have been added, and aditionally returns just the error if only one error was added. Thus, we always have that if you only call `group.Add(err)`, then `group.Err() == err`.

The returned error will format itself similarly:

```go
func groupFormat() {
	var group errs.Group
	group.Add(errs.Errorf("first"))
	group.Add(errs.Errorf("second"))
	err := group.Err()

	fmt.Printf("%v\n", err)
	fmt.Println()
	fmt.Printf("%+v\n", err)

	// output:
	// first; second
	//
	// group:
	// --- first
	//     ... stack trace
	// --- second
	//     ... stack trace
}
```

### Contributing

errs is released under an MIT License. If you want to contribute, be sure to add yourself to the list in AUTHORS.

[Errorf]: https://godoc.org/github.com/zeebo/errs#Errorf
[Wrap]: https://godoc.org/github.com/zeebo/errs#Wrap
[Tag]: https://godoc.org/github.com/zeebo/errs#Tag
[TagWrap]: https://godoc.org/github.com/zeebo/errs#Tag.Wrap
[Tags]: https://godoc.org/github.com/zeebo/errs#Tags
[Group]: https://godoc.org/github.com/zeebo/errs#Group
[GroupAdd]: https://godoc.org/github.com/zeebo/errs#Group.Add
[GroupErr]: https://godoc.org/github.com/zeebo/errs#Group.Err
