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

// CodeErr represents error code type
type ErrCode int

// Here are some common error codes
// you can define your own error code in your package
const (
	NONE ErrCode = iota
	ERROR
	UNFOUND
)

// Err is an error that has a message and a code
type Err struct {
	msg  string
	code ErrCode
}

func (err *Err) Error() string {
	return err.msg
}

func (err *Err) Code() ErrCode {
	return err.code
}

// NewWithCode returns an error with the supplied message and error code.
func NewWithCode(errCode ErrCode, format string, args ...interface{}) error {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	return &Err{
		msg:  msg,
		code: errCode,
	}
}

// New returns an error with the supplied message, but no error code.
func New(format string, args ...interface{}) error {
	return NewWithCode(ERROR, format, args...)
}

// AppendWithCode is used to append error information and pass the information to the caller on a higher level.
// Unlike Append, you can add error code.
// The new error code will overwrite the old one
func AppendWithCode(e error, errCode ErrCode, format string, args ...interface{}) error {
	var msg string
	if len(args) == 0 {
		msg = format
	} else {
		msg = fmt.Sprintf(format, args...)
	}
	msg += fmt.Sprintf("|%v", e)

	return &Err{
		msg:  msg,
		code: errCode,
	}
}

// Append is used to append error information and pass the information to the caller on a higher level.
// If it is an err type, the code in it will be inherited.
func Append(e error, format string, args ...interface{}) error {
	var errCode ErrCode = ERROR
	if err, ok := e.(*Err); ok {
		errCode = err.Code()
	}
	return AppendWithCode(e, errCode, format, args...)
}

// Is judges whether an error is a specific type of error by comparing the error code,
// and then returns the bool value.
func Is(e error, code ErrCode) bool {
	if err, ok := e.(*Err); ok {
		if err.Code() == code {
			return true
		}
	}
	return false
}

// Code returns an error code. If it is not err type, it returns None.
func Code(e error) ErrCode {
	if err, ok := e.(*Err); ok {
		return err.Code()
	}
	return NONE
}
