package utils

import (
	"bufio"
	"fmt"
	"os"
)

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
		return Ask(q + ": ")
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
