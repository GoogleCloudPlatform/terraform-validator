package ancestrymanager

import (
	"fmt"
	"strings"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
)

func AssetParent(cai *resources.Asset, ancestors []string) (string, error) {
	if err := validateAncestors(ancestors); err != nil {
		return "", err
	}
	if cai == nil {
		return "", fmt.Errorf("asset not provided")
	}
	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Folder":
		if len(ancestors) < 2 {
			return "", fmt.Errorf("unexpected ancestors %s", ancestors)
		}
		parent := ancestors[1]
		if strings.HasPrefix(parent, "folders/") || strings.HasPrefix(parent, "organizations/") {
			return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
		}
	case "cloudresourcemanager.googleapis.com/Organization":
		return "", nil
	case "cloudresourcemanager.googleapis.com/Project":
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected ancestors %s", ancestors)
		}
		if len(ancestors) > 1 {
			return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
		}
		// project creation/update
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	default:
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected ancestors %s", ancestors)
		}
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	}
	return "", fmt.Errorf("unexpected ancestors: %v", ancestors)
}

func validateAncestors(ancestors []string) error {
	for _, ancestor := range ancestors {
		if strings.HasPrefix(ancestor, "organizations/") {
			s := strings.TrimPrefix(ancestor, "organizations/")
			if len(s) == 0 {
				return fmt.Errorf("ancestor %v does not have an ID in %v", ancestor, ancestors)
			}
		} else if strings.HasPrefix(ancestor, "folders/") {
			s := strings.TrimPrefix(ancestor, "folders/")
			if len(s) == 0 {
				return fmt.Errorf("ancestor %v does not have an ID in %v", ancestor, ancestors)
			}
		} else if strings.HasPrefix(ancestor, "projects/") {
			s := strings.TrimPrefix(ancestor, "projects/")
			if len(s) == 0 {
				return fmt.Errorf("ancestor %v does not have an ID in %v", ancestor, ancestors)
			}
		} else {
			return fmt.Errorf("unexpected type of ancestor %v in %v", ancestor, ancestors)
		}
	}
	return nil
}

// ConvertToAncestryPath composes a path containing organization/folder/project
// (i.e. "organization/my-org/folder/my-folder/project/my-prj").
func ConvertToAncestryPath(as []string) string {
	var path []string
	for i := len(as) - 1; i >= 0; i-- {
		path = append(path, as[i])
	}
	str := strings.Join(path, "/")
	return sanitizeAncestryPath(str)
}

func sanitizeAncestryPath(s string) string {
	ret := s
	// convert back to match existing ancestry path style.
	for _, r := range []struct {
		old string
		new string
	}{
		{"organizations/", "organization/"},
		{"folders/", "folder/"},
		{"projects/", "project/"},
	} {
		ret = strings.ReplaceAll(ret, r.old, r.new)
	}
	return ret
}
