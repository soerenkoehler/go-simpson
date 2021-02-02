package build

import (
	"fmt"
	"strings"

	"github.com/soerenkoehler/simpson/util"
)

// TargetSpec describes a build target architecture.
type TargetSpec struct {
	os          string
	arch        string
	archiveType util.ArchiveType
}

// Desc returns a string representation of the TargetSpec.
func (target TargetSpec) Desc() string {
	return fmt.Sprintf("%v-%v", target.os, target.arch)
}

// Env returns a list of environment variables for the Go compiler based on the
// TargetSpec.
func (target TargetSpec) Env() []string {
	environment := []string{
		fmt.Sprintf("GOOS=%v", target.os),
		fmt.Sprintf("GOARCH=%v", target.arch)}
	if target.arch == "arm" {
		environment = append(environment, "GOARM=7")
	}
	return environment
}

// GetTargets ... TODO
func GetTargets(filterList string) ([]TargetSpec, []string) {
	result := []TargetSpec{}
	unknown := []string{}
	for _, filter := range strings.Split(filterList, ",") {
		if target, found := findTarget(filter); found {
			result = append(result, target)
		} else {
			unknown = append(unknown, filter)
		}
	}
	return result, unknown
}

func findTarget(filter string) (TargetSpec, bool) {
	for _, target := range AllTargets {
		if target.Desc() == filter {
			return target, true
		}
	}
	return TargetSpec{}, false
}

// Predefined TargetSpecs
var (
	targetWinAmd64   = TargetSpec{"windows", "amd64", util.ZIP}
	targetLinuxAmd64 = TargetSpec{"linux", "amd64", util.TGZ}
	targetLinuxArm   = TargetSpec{"linux", "arm", util.TGZ}
	targetLinuxArm64 = TargetSpec{"linux", "arm64", util.TGZ}

	AllTargets = []TargetSpec{
		targetWinAmd64,
		targetLinuxAmd64,
		targetLinuxArm,
		targetLinuxArm64,
	}
)
