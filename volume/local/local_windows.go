// Package local provides the default implementation for volumes. It
// is used to mount data volume containers and directories local to
// the host server.
package local

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/volume"
)

// scopedPath verifies that the path where the volume is located
// is under Docker's root and the valid local paths.
func (r *Root) scopedPath(realPath string) bool {
	if strings.HasPrefix(realPath, filepath.Join(r.scope, volumesPathName)) && realPath != filepath.Join(r.scope, volumesPathName) {
		return true
	}
	return false
}

// validateVolumeName does platform specific validation of the name of the
// volume being created.
func validateVolumeName(name string) error {
	nameExp := regexp.MustCompile(`^` + volume.RXName + `$`)
	if !nameExp.MatchString(name) {
		return fmt.Errorf("Volume name '%s' is invalid. Must match regex %s", name, volume.RXName)
	}
	return nil
}
