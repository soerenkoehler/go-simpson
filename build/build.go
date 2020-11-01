package build

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

var buildDate = time.Now().UTC().Format("2006.01.02-15.04.05")

// TestAndBuild performs the standard build process.
func TestAndBuild(packageName string, targets []*TargetSpec) {
	if Test() == nil {
		Build(packageName, targets)
	}
}

// Test runs 'go test' for all packages in the current module.
func Test() error {
	return execute([]string{"go", "test", "./..."})
}

// Build runs 'go build' for the named package and supplied target definitions.
// The resulting binary is stored in target specific subdirectories of the
// directory 'artifacts'.
func Build(packageName string, targets []*TargetSpec) {
	for _, target := range targets {
		buildOneTarget(packageName, target, "artifacts")
	}
}

func buildOneTarget(packageName string,
	target *TargetSpec,
	artifactDir string) error {
	// TODO include git ref info
	version := fmt.Sprintf("%s %s", buildDate, target.Desc())

	return execute(
		[]string{
			"go",
			"build",
			"-a",
			"-ldflags", fmt.Sprintf("-X \"main._Version=%s\"", version),
			"-o",
			target.Mkdir(artifactDir),
			packageName},
		target.Env()...)
}

func execute(cmdline []string, env ...string) error {
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
