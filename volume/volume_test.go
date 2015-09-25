package volume

import (
	"runtime"
	"strings"
	"testing"
)

func TestParseMountSpec(t *testing.T) {
	var (
		valid   []string
		invalid map[string]string
	)

	if runtime.GOOS == "windows" {
		valid = []string{
			`d:\`,
			`d:`,
			`d:\path`,
			`d:\path with space`,
			`d:\pathandmode:rw`,
			// TODO Windows post TP4 - readonly support `d:\pathandmode:ro`,
			`c:\:d:\`,
			`c:\windows\:d:`,
			`c:\windows:d:\s p a c e`,
			`c:\windows:d:\s p a c e:RW`,
			`c:\program files:d:\s p a c e i n h o s t d i r`,
			`0123456789name:d:`,
			`MiXeDcAsEnAmE:d:`,
			`name:D:`,
			`name:D::rW`,
			`name:D::RW`,
			// TODO Windows post TP4 - readonly support `name:D::RO`,
			`c:/:d:/forward/slashes/are/good/too`,
			// TODO Windows post TP4 - readonly support `c:/:d:/including with/spaces:ro`,
		}
		invalid = map[string]string{
			``:                                 "Invalid volume specification: ",
			`.`:                                "Invalid volume specification: ",
			`..\`:                              "Invalid volume specification: ",
			`c:\:..\`:                          "Invalid volume specification: ",
			`c:\:d:\:xyzzy`:                    "Invalid volume specification: ",
			`c:`:                               "cannot be c:",
			`c:\`:                              `cannot be c:\`,
			`c:\notexist:d:`:                   `The system cannot find the file specified`,
			`c:\windows\system32\ntdll.dll:d:`: `Source 'c:\windows\system32\ntdll.dll' is not a directory`,
			`_name:d:`:                         `Invalid volume specification`,
			`-name:d:`:                         `Invalid volume specification`,
			`.name:d:`:                         `Invalid volume specification`,
			`name!:d:`:                         `Invalid volume specification`,
			`name@:d:`:                         `Invalid volume specification`,
			`name#:d:`:                         `Invalid volume specification`,
			`name$:d:`:                         `Invalid volume specification`,
			`name%:d:`:                         `Invalid volume specification`,
			`name^:d:`:                         `Invalid volume specification`,
			`name&:d:`:                         `Invalid volume specification`,
			`name*:d:`:                         `Invalid volume specification`,
			`name(:d:`:                         `Invalid volume specification`,
			`name):d:`:                         `Invalid volume specification`,
			`name+:d:`:                         `Invalid volume specification`,
			`name=:d:`:                         `Invalid volume specification`,
			`name{:d:`:                         `Invalid volume specification`,
			`name}:d:`:                         `Invalid volume specification`,
			`name[:d:`:                         `Invalid volume specification`,
			`name]:d:`:                         `Invalid volume specification`,
			`name\:d:`:                         `Invalid volume specification`,
			`name|:d:`:                         `Invalid volume specification`,
			`name;d:`:                          `Invalid volume specification`,
			`name':d:`:                         `Invalid volume specification`,
			`name<:d:`:                         `Invalid volume specification`,
			`name>:d:`:                         `Invalid volume specification`,
			`name?:d:`:                         `Invalid volume specification`,
			`name/:d:`:                         `Invalid volume specification`,
		}

	} else {
		valid = []string{
			"/home",
			"/home:/home",
			"/home:/something/else",
			"/with space",
			"/home:/with space",
			"relative:/absolute-path",
			"hostPath:/containerPath:ro",
			"/hostPath:/containerPath:rw",
			"/rw:/ro",
			"/path:rw",
			"/path:ro",
			"/rw:rw",
		}
		invalid = map[string]string{
			"":                "Invalid volume specification",
			"./":              "Invalid volume destination",
			"../":             "Invalid volume destination",
			"/:../":           "Invalid volume destination",
			"/:path":          "Invalid volume destination",
			":":               "Invalid volume specification",
			"/tmp:":           "Invalid volume destination",
			":test":           "Invalid volume specification",
			":/test":          "Invalid volume specification",
			"tmp:":            "Invalid volume destination",
			":test:":          "Invalid volume specification",
			"::":              "Invalid volume specification",
			":::":             "Invalid volume specification",
			"/tmp:::":         "Invalid volume specification",
			":/tmp::":         "Invalid volume specification",
			"path:ro":         "mount path must be absolute",
			"/path:/path:sw":  "invalid mode for",
			"/path:/path:rwz": "invalid mode for",
		}
	}

	for _, path := range valid {
		if _, err := ParseMountSpec(path, "local"); err != nil {
			t.Fatalf("ParseMountSpec(`%q`) should succeed: error %q", path, err)
		}
	}

	var i = 1
	for path, expectedError := range invalid {
		if _, err := ParseMountSpec(path, "local"); err == nil {
			t.Fatalf("ParseMountSpec(`%q`) should have failed validation test %d", path, i)
		} else {
			if !strings.Contains(err.Error(), expectedError) {
				t.Fatalf("ParseMountSpec(`%q`) error should contain %q, got %q test %d", path, expectedError, err.Error(), i)
			}
		}
		i++
	}
}
