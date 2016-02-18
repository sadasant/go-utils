// +build windows

package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func AskPass(q string) string {
	var pass string

	cloudwalk_dir := GetPathFromHome("/.cloudwalk")
	pass_file := cloudwalk_dir + "/getpass.sh"
	pass_script := `#!/bin/bash
function finish {
	reset
}
trap finish EXIT
unset password
read -e -s password
echo $password
`

	os.Mkdir(cloudwalk_dir, 0755)
	ioutil.WriteFile(pass_file, []byte(pass_script), 0755)

	os.Chmod(cloudwalk_dir, 0755)
	os.Chmod(pass_file, 0755)

	cmd := exec.Command("bash", pass_file)

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stdin = os.Stdin

	fmt.Printf("%s", q)

	// Clean OS Exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		if pass == "" {
			println("\nDone. Press Enter.")
		}
	}()

	if err := cmd.Start(); err != nil {
		Exitf(2, "Failed to start the password request. Error: %s\n", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		Exitf(2, "Failed to wait for the password request. Error: %s\n", err.Error())
	}

	pass = b.String()
	pass = pass[:len(pass)-1]

	os.Remove(pass_file)

	close(c)
	println()
	return pass
}
