package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func AskPass(q string) string {
	print(q)

	oldState, err := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	ExitfOn("Error on terminal.MakeRaw: %v\n", err)

	reset := CatchSignal(os.Interrupt, func() {
		terminal.Restore(0, oldState)
	})

	pass, err := terminal.ReadPassword(0)
	ExitfOn("Error on terminal.ReadPassword: %v\n", err)

	reset()

	println()
	return string(pass)
}

func Ask(q string) string {
	fmt.Printf("%s", q)
	r := bufio.NewReader(os.Stdin)
	line, _, err := r.ReadLine()
	if err != nil {
		Exit(1, err.Error())
	}
	return string(line)
}

func AskPrevious(q, old string) string {
	if old == "" {
		return Ask(q)
	} else {
		resp := Ask(fmt.Sprintf("%s [%s]", q, old))
		if resp == "" {
			return old
		}
		return resp
	}
}

func AskYN(q string) bool {
	q = fmt.Sprintf("%s (y/n)", q)
	switch Ask(q) {
	case "y", "Y":
		return true
	case "n", "N":
		return false
	default:
		return AskYN(q)
	}
}
