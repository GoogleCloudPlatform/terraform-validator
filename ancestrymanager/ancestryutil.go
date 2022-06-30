package ancestrymanager

import (
	"fmt"
	"strings"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"

	"go.uber.org/zap"
)

// assetParent derives a resource's parent from its ancestors.
func assetParent(cai *resources.Asset, ancestors []string) (string, error) {
	if cai == nil {
		return "", fmt.Errorf("asset not provided")
	}
	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Folder":
		if len(ancestors) < 2 {
			return "", fmt.Errorf("unexpected value for ancestors: %s", ancestors)
		}
		parent := ancestors[1]
		if strings.HasPrefix(parent, "folders/") || strings.HasPrefix(parent, "organizations/") {
			return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
		}
	case "cloudresourcemanager.googleapis.com/Organization":
		return "", nil
	case "cloudresourcemanager.googleapis.com/Project":
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected value for ancestors: %s", ancestors)
		}
		if strings.HasPrefix(ancestors[0], "projects/") {
			if len(ancestors) > 1 {
				return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[1]), nil
			}
		}
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	default:
		if len(ancestors) < 1 {
			return "", fmt.Errorf("unexpected value for ancestors: %s", ancestors)
		}
		return fmt.Sprintf("//cloudresourcemanager.googleapis.com/%s", ancestors[0]), nil
	}
	return "", fmt.Errorf("unexpected value for ancestors: %v", ancestors)
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

// getProjectFromResource reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFromResource(d resources.TerraformResourceData, config *resources.Config, cai resources.Asset, errorLogger *zap.Logger) (string, error) {

	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Project",
		"cloudbilling.googleapis.com/ProjectBillingInfo":
		res, ok := d.GetOk("number")
		if ok {
			return res.(string), nil
		}
		// Fall back to project_id if number is not available.
		res, ok = d.GetOk("project_id")
		if ok {
			return res.(string), nil
		} else {
			errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from resource", cai.Name))
		}
	case "storage.googleapis.com/Bucket":
		if cai.Resource != nil {
			res, ok := cai.Resource.Data["project"]
			if ok {
				return res.(string), nil
			}
		}
		errorLogger.Warn(fmt.Sprintf("Failed to retrieve project_id for %s from cai resource", cai.Name))
	}

	return getProjectFromSchema("project", d, config)
}

func getProjectFromSchema(projectSchemaField string, d resources.TerraformResourceData, config *resources.Config) (string, error) {
	res, ok := d.GetOk(projectSchemaField)
	if ok && projectSchemaField != "" {
		return res.(string), nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("required field '%s' is not set, you may use --project=my-project to provide a default project to resolve the issue", projectSchemaField)
}

// getOrganizationFromResource reads org_id field from terraform data.
func getOrganizationFromResource(tfData resources.TerraformResourceData) (string, bool) {
	orgID, ok := tfData.GetOk("org_id")
	if ok {
		return orgID.(string), ok
	}
	return "", false
}

// getFolderFromResource reads folder_id or folder field from terraform data.
func getFolderFromResource(tfData resources.TerraformResourceData) (string, bool) {
	folderID, ok := tfData.GetOk("folder_id")
	if ok {
		return folderID.(string), ok
	}
	folderID, ok = tfData.GetOk("folder")
	if ok {
		return folderID.(string), ok
	}
	return "", false
}

// isGoogleApiErrorWithCode cheks if the error code is of given type or not.
func isGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}
