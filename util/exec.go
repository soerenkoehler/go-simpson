package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Execute ... TODO
func Execute(cmdline []string, env ...string) error {
	fmt.Println(cmdline, env)

	proc := exec.Command(cmdline[0], cmdline[1:]...)
	proc.Env = append(os.Environ(), env...)
	pipeOut, _ := proc.StdoutPipe()
	pipeErr, _ := proc.StderrPipe()

	output := make(chan string)
	defer close(output)

	go copyOutput(pipeOut, output)
	go copyOutput(pipeErr, output)
	go printOutput(output)

	err := proc.Run()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return err
}

func copyOutput(src io.Reader, dest chan<- string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		dest <- scanner.Text()
	}
}

func printOutput(src <-chan string) {
	for line := range src {
		fmt.Println(line)
	}
}
