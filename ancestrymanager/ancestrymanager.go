// Package ancestrymanager provides an interface to query the ancestry information for a project.
package ancestrymanager

import (
	"fmt"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"go.uber.org/zap"
)

// AncestryManager is the interface that wraps the GetAncestry method.
type AncestryManager interface {
	// GetAncestry takes a project name to return an ancestry path
	GetAncestry(project string) (string, error)
	// GetAncestry takes a project name and resource data to return an ancestry path
	GetAncestryWithResource(project string, tfData resources.TerraformResourceData, cai resources.Asset) (string, error)
}

// ClientRetriever is the interface that returns an instance of various clients.
type ClientRetriever interface {
	// NewResourceManagerClient returns an initialized *cloudresourcemanager.Service
	NewResourceManagerClient(userAgent string) *cloudresourcemanager.Service
}

// resourceAncestryManager provides common methods for retrieving ancestry from resources
type resourceAncestryManager struct {
	errorLogger *zap.Logger
}

func (m *resourceAncestryManager) getFolderAncestry(folderID string) ([]*cloudresourcemanager.Ancestor, error) {
	// TODO(morgantep): Incorporate folders.GetAncestry from v2alpha1 API
	m.errorLogger.Info(fmt.Sprintf("Retrieving ancestry for folder: %s", folderID))

	return []*cloudresourcemanager.Ancestor{
		&cloudresourcemanager.Ancestor{
			ResourceId: &cloudresourcemanager.ResourceId{
				Type: "folder",
				Id:   folderID,
			},
		},
		&cloudresourcemanager.Ancestor{
			ResourceId: &cloudresourcemanager.ResourceId{
				Type: "organization",
				Id:   "unknown",
			},
		},
	}, nil
}

func (m *resourceAncestryManager) getAncestryFromResource(tfData resources.TerraformResourceData, cai resources.Asset) ([]*cloudresourcemanager.Ancestor, bool) {
	m.errorLogger.Info(fmt.Sprintf("Retrieving ancestry from resource (type=%s)", cai.Type))

	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Project", "cloudbilling.googleapis.com/ProjectBillingInfo":
		// Prefer project number to project id if available;
		// CAI exports use project number.
		projectID, ok := tfData.GetOk("number")
		if !ok {
			projectID, ok = tfData.GetOk("project_id")
			if !ok || projectID == "" {
				return nil, false
			}
		}

		ancestry := []*cloudresourcemanager.Ancestor{
			&cloudresourcemanager.Ancestor{
				ResourceId: &cloudresourcemanager.ResourceId{
					Type: "project",
					Id:   projectID.(string),
				},
			},
		}

		orgID, ok := tfData.GetOk("org_id")
		if ok && orgID != "" {
			s := strings.Split(orgID.(string), "/")
			return append(ancestry, &cloudresourcemanager.Ancestor{
				ResourceId: &cloudresourcemanager.ResourceId{
					Type: "organization",
					Id:   s[len(s)-1],
				},
			}), true
		}

		folderID, ok := tfData.GetOk("folder_id")
		if ok && folderID != "" {
			folderAncestry, err := m.getFolderAncestry(folderID.(string))
			if err != nil {
				m.errorLogger.Error(fmt.Sprintf("Failed to retrieve folder ancestry: %s", err))
				return nil, false
			}
			return append(ancestry, folderAncestry...), true
		}
		return nil, false
	default:
		m.errorLogger.Info(fmt.Sprintf("Resource of type %s does not include sufficient data for ancestry retrieval", cai.Type))
		return nil, false
	}
}

type manager struct {
	resourceAncestryManager
	// Talk to GCP resource manager. This field would be nil in offline mode.
	resourceManager *cloudresourcemanager.Service
	// Cache to prevent multiple network calls for looking up the same project's ancestry
	// map[project]ancestryPath
	ancestryCache map[string]string
	offline       bool
}

// GetAncestry uses the resource manager API to get ancestry paths for
// projects. It implements a cache because many resources share the same
// project.
func (m *manager) GetAncestry(project string) (string, error) {
	if path, ok := m.ancestryCache[project]; ok {
		return path, nil
	}

	if m.offline {
		return "", fmt.Errorf("cannot fetch ancestry in offline mode")
	}
	ancestry, err := m.resourceManager.Projects.GetAncestry(project, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		return "", err
	}
	path := ancestryPath(ancestry.Ancestor)
	m.store(project, path)
	return path, nil
}

func (m *manager) store(project, ancestry string) {
	if project != "" && ancestry != "" {
		m.ancestryCache[project] = ancestry
	}
}

// GetAncestryWithResource first attempts to get Ancestry from the resource
// If that fails, it falls back to the offline cache.
func (m *manager) GetAncestryWithResource(project string, tfData resources.TerraformResourceData, cai resources.Asset) (string, error) {
	ancestry, ok := m.getAncestryFromResource(tfData, cai)
	if ok {
		path := ancestryPath(ancestry)
		m.errorLogger.Info(fmt.Sprintf("[Offline] Retrieved ancestry for %s: %s", project, path))
		return path, nil
	}

	path, err := m.GetAncestry(project)
	if err != nil {
		return "", err
	}

	return path, nil
}

// New returns AncestryManager that can be used to fetch ancestry information for a project.
// entries takes project as key and ancestry as value
func New(offline bool, retriever ClientRetriever, entries map[string]string, userAgent string, errorLogger *zap.Logger) (AncestryManager, error) {
	am := &manager{
		ancestryCache: map[string]string{},
		resourceAncestryManager: resourceAncestryManager{
			errorLogger: errorLogger,
		},
		offline: offline,
	}
	if !offline {
		rm := retriever.NewResourceManagerClient(userAgent)
		am.resourceManager = rm
	}
	for project, ancestry := range entries {
		if ancestry != "" {
			ancestry = fmt.Sprintf("%s/project/%s", ancestry, project)
		}
		am.store(project, ancestry)
	}
	return am, nil
}

// ancestryPath composes a path containing organization/folder/project
// (i.e. "organization/my-org/folder/my-folder/project/my-prj").
func ancestryPath(as []*cloudresourcemanager.Ancestor) string {
	var path []string
	for i := len(as) - 1; i >= 0; i-- {
		path = append(path, as[i].ResourceId.Type, as[i].ResourceId.Id)
	}
	return strings.Join(path, "/")
}
