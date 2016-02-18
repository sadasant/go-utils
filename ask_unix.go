// +build darwin freebsd linux netbsd openbsd

package utils

import (
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
