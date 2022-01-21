package version

// To be set by Go build tools.
var buildVersion string = "dev"

// BuildVersion returns the build version of Terraform Validator.
func BuildVersion() string {
	return buildVersion
}
