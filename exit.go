package utils

import (
	"fmt"
	"os"
)

func Errorf(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", values...)
}

func Exit(code int, values ...interface{}) {
	fmt.Fprintln(os.Stderr, values...)
	os.Exit(code)
}

func Exitf(code int, format string, values ...interface{}) {
	Errorf(format, values...)
	os.Exit(code)
}

func ExitOn(err error) {
	if err != nil {
		Exit(1, err.Error())
	}
}

func ExitfOn(str string, values ...interface{}) {
	last := len(values) - 1
	err, ok := values[last].(error)
	if !ok {
		return
	}
	values = values[:last]
	if err != nil {
		values = append(values, err.Error())
		Exitf(1, str, values...)
	}
}
