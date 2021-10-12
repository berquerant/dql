package dql

import (
	"fmt"
	"strings"
)

var (
	Commit    string
	Version   string
	GoVersion string
)

func Meta() string {
	return strings.Join([]string{
		fmt.Sprintf("Version:\t%s", Version),
		fmt.Sprintf("Go version:\t%s", GoVersion),
		fmt.Sprintf("Git commit:\t%s", Commit),
	}, "\n")
}
