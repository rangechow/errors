# errors [![Build Status](https://travis-ci.org/rangechow/errors.svg?branch=main)](https://travis-ci.org/rangechow/errors) [![Go Reference](https://pkg.go.dev/badge/github.com/rangechow/errors.svg)](https://pkg.go.dev/github.com/rangechow/errors) [![Go Report Card](https://goreportcard.com/badge/github.com/rangechow/errors)](https://goreportcard.com/report/github.com/rangechow/errors)
Package errors provides error with error code.

在golang中，我们一般使用的error是基于字符串的。

	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}

如果是常量字符串，那么会比较方便的处理。

    var ErrNotFound = errors.New("not found")
    if err == ErrNotFound {
        // something wasn't found
    }

如果字符串中带有变量，处理起来就比较棘手了。

	err := errors.New("%s not found", filename)

一般需要定义一个error的类型，通过转换来判断是否是该错误。

	type NotFoundError struct {
		Name string
	}
 
	func (e *NotFoundError) Error() string { return e.Name + ": not found" }

	if e, ok := err.(*NotFoundError); ok {
		// e.Name wasn't found
	}

在C/C++中，我们就没有这种烦恼。因为我们使用的是错误码。通过操作码的比较，可以直接判断是哪种类型的错误。
golang的error设计是非常好的。为了弥补这一点不足，我设计了一个带有错误码的error包。
如果你需要在调用函数中处理不同的错误返回，你可以使用带错误码的error。


    err := errors.NewWithCode(errors.UNFOUND,"not found")

	if err != nil {
		if errors.Is(err, errors.UNFOUND) {
			// something wasn't found
		} else {
			fmt.Print(err)
		}
	}

