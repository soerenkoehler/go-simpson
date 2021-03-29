package build

import (
	"fmt"

	"github.com/soerenkoehler/simpson/util"
)

type TargetSpec struct {
	os                  string
	arch                string
	executableExtension string
	archiveType         util.ArchiveType
}

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

func GetTargets(filters []string) ([]TargetSpec, []string) {
	result := []TargetSpec{}
	unknown := []string{}
	for _, filter := range filters {
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
	targetWinAmd64 = TargetSpec{
		os:                  "windows",
		arch:                "amd64",
		executableExtension: "exe",
		archiveType:         util.ZIP}
	targetLinuxAmd64 = TargetSpec{
		os:                  "linux",
		arch:                "amd64",
		executableExtension: "",
		archiveType:         util.TGZ}
	targetLinuxArm = TargetSpec{
		os:                  "linux",
		arch:                "arm",
		executableExtension: "",
		archiveType:         util.TGZ}
	targetLinuxArm64 = TargetSpec{
		os:                  "linux",
		arch:                "arm64",
		executableExtension: "",
		archiveType:         util.TGZ}

	AllTargets = []TargetSpec{
		targetWinAmd64,
		targetLinuxAmd64,
		targetLinuxArm,
		targetLinuxArm64,
	}
)
