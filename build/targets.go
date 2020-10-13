package build

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// TargetSpec describes a build target architecture.
type TargetSpec struct {
	os   string
	arch string
}

// Desc returns a string representation of the TargetSpec.
func (target *TargetSpec) Desc() string {
	return fmt.Sprintf("%s-%s", target.os, target.arch)
}

// Mkdir creates a subdirectory in 'parent' based on the TargetSpec and returns
// the absolute path.
func (target *TargetSpec) Mkdir(parent string) string {
	targetDir := path.Join(parent, target.Desc())
	os.MkdirAll(targetDir, 0777)
	result, error := filepath.Abs(targetDir)
	if error != nil {
		panic(error)
	}
	return result
}

// Env returns a list of environment variables for the Go compiler based on the
// TargetSpec.
func (target *TargetSpec) Env() []string {
	environment := []string{
		fmt.Sprintf("GOOS=%s", target.os),
		fmt.Sprintf("GOARCH=%s", target.arch)}
	if target.arch == "arm" {
		environment = append(environment, "GOARM=7")
	}
	return environment
}

// Predefined TargetSpecs
var (
	TargetWinAmd64   = &TargetSpec{"windows", "amd64"}
	TargetLinuxAmd64 = &TargetSpec{"linux", "amd64"}
	TargetLinuxArm   = &TargetSpec{"linux", "arm"}
	TargetLinuxArm64 = &TargetSpec{"linux", "arm64"}

	AllTargets = []*TargetSpec{
		TargetWinAmd64,
		TargetLinuxAmd64,
		TargetLinuxArm,
		TargetLinuxArm64,
	}
)
