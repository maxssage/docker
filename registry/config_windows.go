package registry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	// DefaultV1Registry is the URI of the default v1 registry
	DefaultV1Registry = "https://registry-win-tp3.docker.io"

	// DefaultV2Registry is the URI of the default (official) v2 registry.
	// This is the windows-specific endpoint.
	//
	// Currently it is a TEMPORARY link that allows Microsoft to continue
	// development of Docker Engine for Windows.
	DefaultV2Registry = "https://registry-win-tp3.docker.io"

	// CertsDir is the directory where certificates are stored
	CertsDir = os.Getenv("programdata") + `\docker\certs.d`
)

// init here checks an override environment variable to allow an alternate
// registry. This is for development purposes only.
func init() {
	if len(os.Getenv("DOCKER_REGISTRY_OVERRIDE")) > 0 {
		DefaultV1Registry = "https://index.docker.io"
		DefaultV2Registry = "https://registry-1.docker.io"
		logrus.Warnf("***DEVELOPMENT OVERRIDE FOR REGISTRY***")
	}

}

// cleanPath is used to ensure that a directory name is valid on the target
// platform. It will be passed in something *similar* to a URL such as
// https:\index.docker.io\v1. Not all platforms support directory names
// which contain those characters (such as : on Windows)
func cleanPath(s string) string {
	return filepath.FromSlash(strings.Replace(s, ":", "", -1))
}
