package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

func CatchSignal(sig os.Signal, f func()) func() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, sig)

	go func() {
		for {
			select {
			case _, ok := <-signalChan:
				if ok {
					f()
					os.Exit(1)
				} else {
					return
				}
			}
		}
	}()

	return func() {
		signal.Reset(sig)
		signal.Stop(signalChan)
	}
}

func getCommandParts(command string) []string {
	parts := []string{""}
	var index int
	var in_block bool
	for _, v := range command {
		if v == '"' || v == '\'' {
			in_block = true
		}
		if !in_block && v == '(' {
			in_block = true
		}
		if in_block && v == ')' {
			in_block = false
		}
		if !in_block && v == ' ' {
			index++
			parts = append(parts, "")
			continue
		}
		parts[index] += string(v)
	}
	return parts
}

func RunCommands(commands string) {
	command_lines := strings.Split(commands, "\n")
	for _, v := range command_lines {
		if v != "" {
			RunCommand(v)
		}
	}
}

func RunCommandf(format string, values ...interface{}) {
	RunCommand(fmt.Sprintf(format, values...))
}

func RunCommand(command string) {

	var last_wait bool // To know if the last printed value was one of the wait dots
	var done bool      // To know when the command finished, to avoid printing more dots
	var exit bool      // To know if we should exit the whole process, in case the output was wrong

	parts := getCommandParts(command)
	// Exception for sudo docker login: change the password for ****
	print_parts := []string{}
	for k, v := range parts {
		if k > 0 && parts[1] == "docker" && parts[k-1] == "-p" {
			print_parts = append(print_parts, "****")
			continue
		}
		print_parts = append(print_parts, v)
	}
	println(Bold("Running:"), strings.Join(print_parts, " "))
	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Stdin = os.Stdin

	// Showing some wait information to make sure the user knows we're working
	go func() {
		time.Sleep(10 * time.Second)
		if done {
			return
		}
		print("Please wait...")
		last_wait = true
	}()

	cmdOutput, err := cmd.StdoutPipe()
	ExitfOn("Failed to StdoutPipe: %s\nError: %s", Red(command), err)

	scanner := bufio.NewScanner(cmdOutput)
	go func() {
		for scanner.Scan() {
			if last_wait {
				println()
				last_wait = false
			}
			text := scanner.Text() + "\r\n"
			if strings.Contains(text, "fatal:") {
				exit = true
				print(Red(text))
				break
			}
			if strings.Contains(text, "cannot") {
				exit = true
				print(Red(text))
				break
			}
			print(Green(text))
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "scanner error:", err)
		}
	}()

	cmdErr, err := cmd.StderrPipe()
	ExitfOn("Failed to StderrPipe: %s\nError: %s", Red(command), err)

	scanner2 := bufio.NewScanner(cmdErr)
	go func() {
		for scanner2.Scan() {
			if last_wait {
				println()
				last_wait = false
			}
			text := scanner2.Text() + "\r\n"
			if strings.Contains(text, "fatal:") {
				exit = true
				print(Red(text))
				break
			}
			if strings.Contains(text, "cannot") {
				exit = true
				print(Red(text))
				break
			}
			print(Green(text))
		}
		if err := scanner2.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "scanner2 error:", err)
		}
	}()

	err = cmd.Start()
	ExitfOn("Failed to start: %s\nError: %s", Red(command), err)

	for {
		err = cmd.Wait()
		ExitfOn("Failed to wait for: %s\nError: %s", Red(command), err)
		break
	}

	done = true

	if exit {
		Exit(1)
	}
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
