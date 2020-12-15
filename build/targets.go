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
	archiveType *util.ArchiveType
}

// Desc returns a string representation of the TargetSpec.
func (target *TargetSpec) Desc() string {
	return fmt.Sprintf("%s-%s", target.os, target.arch)
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

// GetTargets ... TODO
func GetTargets(filterList string) []*TargetSpec {
	result := []*TargetSpec{}
	for _, filter := range strings.Split(filterList, ",") {
		if target := findTarget(filter); target != nil {
			result = append(result, target)
		}
	}
	return result
}

func findTarget(filter string) *TargetSpec {
	for _, target := range AllTargets {
		if target.Desc() == filter {
			return target
		}
	}
	return nil
}

// Predefined TargetSpecs
var (
	TargetWinAmd64   = &TargetSpec{"windows", "amd64", util.ZIP}
	TargetLinuxAmd64 = &TargetSpec{"linux", "amd64", util.TGZ}
	TargetLinuxArm   = &TargetSpec{"linux", "arm", util.TGZ}
	TargetLinuxArm64 = &TargetSpec{"linux", "arm64", util.TGZ}

	AllTargets = []*TargetSpec{
		TargetWinAmd64,
		TargetLinuxAmd64,
		TargetLinuxArm,
		TargetLinuxArm64,
	}
)
