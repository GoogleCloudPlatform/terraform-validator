// Package version provides the list of supported Terraform versions.
package version

const (
	TF11 string = "0.11"
	TF12 string = "0.12"
)

var supportedMap = map[string]bool{}

func init() {
	for _, v := range supportedList {
		supportedMap[v] = true
	}
}

// Supported checks if the version of Terraform is supported.
func Supported(version string) bool {
	return supportedMap[version]
}

// LeastSupportedVersion returns the minimal supported version.
func LeastSupportedVersion() string {
	return supportedList[0]
}

// AllSupportedVersions returns all versions supported.
func AllSupportedVersions() []string {
	return supportedList
}
