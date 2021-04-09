/*
 *	Package errors provides error with error code.
 *
 *	In golang, the error we generally use is string-based.
 *
 *		err := errors.New("emit macho dwarf: elf header corrupted")
 *		if err != nil {
 *			fmt.Print(err)
 *		}
 *
 *	If it is a constant string, it will be more convenient to deal with.
 *
 *		var ErrNotFound = errors.New("not found")
 *		if err == ErrNotFound {
 *			// something wasn't found
 *		}
 *
 *	there are variables in the string, it will be tricky to deal with.
 *	Generally, it is necessary to define an error type, and determine whether it is the error through conversion.
 *
 *		type NotFoundError struct {
 *			Name string
 *		}
 *
 *		func (e *NotFoundError) Error() string { return e.Name + ": not found" }
 *
 *		if e, ok := err.(*NotFoundError); ok {
 *			// e.Name wasn't found
 *		}
 *
 *	In C/C++, we don't have this trouble. Because we are using error codes.
 *	By comparing the opcodes, you can directly determine which type of error it is.
 *	The error design of golang is very good.
 *	In order to make up for this shortcoming, I designed an error package with an error code.
 *
 *	If you need to handle different error returns in the calling function, you can use error with error code.
 *
 *		err := errors.NewWithCode(errors.UNFOUND,"not found")
 *
 *		if err != nil {
 *			if errors.Is(err, errors.UNFOUND) {
 *				// something wasn't found
 *			} else {
 *				fmt.Print(err)
 *			}
 *		}
 *
 *  See the documentation for more details.
 */

package errors

import (
	"fmt"
)

// error code type
type CodeErr int

// Here are some common error codes
// you can define your own error code in your package
const (
	NONE CodeErr = iota
	ERROR
	UNFOUND
)

// err is an error that has a message and a code
type err struct {
	msg  string
	code CodeErr
}

func (e *err) Error() string {
	return e.msg
}

func (e *err) Code() CodeErr {
	return e.code
}

// New returns an error with the supplied message, but no error code.
func New(format string, args ...interface{}) error {
	return &err{
		msg:  fmt.Sprintf(format, args...),
		code: ERROR,
	}
}

// NewWithCode returns an error with the supplied message and error code.
func NewWithCode(errCode CodeErr, format string, args ...interface{}) error {
	return &err{
		msg:  fmt.Sprintf(format, args...),
		code: errCode,
	}
}

// Append is used to append error information and pass the information to the caller on a higher level.
// If it is an err type, the code in it will be inherited.
func Append(format string, args ...interface{}) error {
	var errCode CodeErr = ERROR
	if len(args) > 0 {
		if err, ok := args[len(args)-1].(err); ok {
			format = format + "|%v"
			errCode = err.Code()
		} else if _, ok := args[len(args)-1].(error); ok {
			format = format + "|%v"
		}
	}
	return &err{
		msg:  fmt.Sprintf(format, args...),
		code: errCode,
	}
}

// AppendWithCode is used to append error information and pass the information to the caller on a higher level.
// Unlike Append, you can add error code.
func AppendWithCode(errCode CodeErr, format string, args ...interface{}) error {
	var code CodeErr = errCode
	if len(args) > 0 {
		if err, ok := args[len(args)-1].(err); ok {
			format = format + "|%v"
			code = err.Code()
		} else if _, ok := args[len(args)-1].(error); ok {
			format = format + "|%v"
		}
	}
	return &err{
		msg:  fmt.Sprintf(format, args...),
		code: code,
	}
}

// Is judges whether an error is a specific type of error by comparing the error code,
// and then returns the bool value.
func Is(e error, code CodeErr) bool {
	if err, ok := e.(*err); ok {
		if err.Code() == code {
			return true
		}
	}
	return false
}

func Code(e error) CodeErr {
	if err, ok := e.(*err); ok {
		return err.Code()
	}
	return NONE
}
