# errors [![Build Status](https://travis-ci.org/rangechow/errors.svg?branch=main)](https://travis-ci.org/rangechow/errors) [![Go Reference](https://pkg.go.dev/badge/github.com/rangechow/errors.svg)](https://pkg.go.dev/github.com/rangechow/errors) [![Go Report Card](https://goreportcard.com/badge/github.com/rangechow/errors)](https://goreportcard.com/report/github.com/rangechow/errors) [![Sourcegraph](https://sourcegraph.com/github.com/rangechow/errors/-/badge.svg)](https://sourcegraph.com/github.com/rangechow/errors?badge)

## What

Package errors provides error with error code.

## Why & How

在golang中，任何结构只需要实现Error函数都可以作为error对象。Error函数返回一个string。因此我们说，error是基于字符串的。error接口的定义如下：

	type error interface {
		Error() string
	}

获得一个error的方式有很多。通常情况下我们可以使用一个字符串New一个出来或者通过fmt生成。

	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}

如果只是简单判断一个函数是否正确执行，只需要判断error对象是否为nil。但是如果我们需要进一步判断是何种错误，然后做出相应的操作，就稍微复杂一些了。
如果是常量字符串，那么会比较方便的处理。我们只需要生成一个此字符串对应的error对象，然后进行比较即可。

    var ErrNotFound = errors.New("not found")
    if err == ErrNotFound {
        // something wasn't found
    }

如果字符串中带有变量，处理起来就比较棘手了。

	err := errors.New("%s not found", Filename)

go1.13以前，一般需要定义一个新的error类型，通过转换来判断是否是该错误。

	type NotFoundError struct {
		Name string
	}
 
	func (e *NotFoundError) Error() string { return e.Name + ": not found" }

	if e, ok := err.(*NotFoundError); ok {
		// e.Name wasn't found
	}

go1.13开始提供了Wrap方式。通过`fmt.Errorf`函数的`%w`格式可以将error对象包起来。除了打印error中的字符串，还可以通过`Unwrap`函数获得上一层的包裹的error对象。此外go还提供了`Is`、`As`函数直接获取包裹中的error对象。

这个方案并不完美，如果中间过程中带有变量，还是需要定一个错误类型。

在C/C++中，我们通常没有这种烦恼。因为我们使用的是错误码。通过操作码的比较，可以直接判断是哪种类型的错误。
golang的error设计是非常好的，字符串信息可以更加直观的看到错误的问题或者原因。为了弥补判断error上的不足，我设计了一个带有错误码的error包。
如果你需要在调用函数中处理不同的错误返回，你可以使用带错误码的error。


    err := errors.NewWithCode(errors.UNFOUND,"%s not found", Filename)

	if err != nil {
		if errors.Is(err, errors.UNFOUND) {
			// something wasn't found
		} else {
			fmt.Print(err)
		}
	}

如果是中间过程的错误处理，您可以使用`AppendWithCode`函数来追加错误。此时您可以设置新的错误码。新的错误码会覆盖原有的错误码。

	err := errors.New("%s not found", Filename)

	if err != nil {
		return errors.AppendWithCode(errors.UNFOUND, "new err", err)
	}

这里我没有提供打印堆栈的功能，而是在Append中自动将err追加到新错误字符串的最后。基于做好自己本职工作的原则，每次函数调用如果有错误都会打印错误日志。但并不需要每一层都打印大量的堆栈信息。	

当你拿到一个错误，要对不同的错误类型进行处理时，你可以使用`Is`函数来进行判断。

	err := errors.NewWithCode(errors.UNFOUND,"%s not found", Filename)

	if err != nil && errors.Is(err, errors.UNFOUND) {
		// something wasn't found
		return
	}


我还提供了`New`、`Append`等不带错误码的error对象生成函数。欢迎使用！

## FAQ

1.	关于错误码
	
	你可以在自己的文件中定义`CodeErr`类型的错误码。

2.	返回的error在什么情况下需要打Log？

	一般情况下返回的error不需要打印Log，我们可以看到基础的API都没有因error而打印Log。只有当error被处理，即error的发生不影响后续的流程时，需要将error的信息以日志的形式输出。
	中间的过程不打Log，如何判断出错的位置呢？这种场景下，使用堆栈信息确实是一种比较快捷的方案。堆栈信息拥有执行路径上的文件和函数，但缺少了环境信息。
	因此，我选择了另外一种方案，即追加错误的方案。继续基于golang以字符串记录错误信息的设计，你可以在每一层Append错误，附带所需要的环境信息。
