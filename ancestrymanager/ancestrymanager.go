// Package ancestrymanager provides an interface to query the ancestry information for a project.
package ancestrymanager

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v1"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

// AncestryManager is the interface that wraps the GetAncestry method.
type AncestryManager interface {
	// GetAncestry takes a project name to return an ancestry path
	GetAncestry(project string) (string, error)
	// GetAncestry takes a project name and resource data to return an ancestry path
	GetAncestryWithResource(project string, tfData converter.TerraformResourceData, cai converter.Asset) (string, error)
}

// ClientRetriever is the interface that returns an instance of various clients.
type ClientRetriever interface {
	// NewResourceManagerClient returns an initialized *cloudresourcemanager.Service
	NewResourceManagerClient(userAgent string) *cloudresourcemanager.Service
}

// resourceAncestryManager provides common methods for retrieving ancestry from resources
type resourceAncestryManager struct {
}

func (m *resourceAncestryManager) getFolderAncestry(folderID string) ([]*cloudresourcemanager.Ancestor, error) {
	// TODO(morgantep): Incorporate folders.GetAncestry from v2alpha1 API
	log.Printf("[INFO] Retrieve ancestry for folder: %s", folderID)

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

func (m *resourceAncestryManager) getAncestryFromResource(tfData converter.TerraformResourceData, cai converter.Asset) ([]*cloudresourcemanager.Ancestor, bool) {
	log.Printf("[INFO] Retrieving ancestry from resource (type=%s)", cai.Type)

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
				log.Printf("[ERROR] Failed to retrieve folder ancestry: %s", err)
				return nil, false
			}
			return append(ancestry, folderAncestry...), true
		}
		return nil, false
	default:
		log.Printf("[INFO] Resource of type %s does not include sufficient data for ancestry retrieval", cai.Type)
		return nil, false
	}
}

type onlineAncestryManager struct {
	resourceAncestryManager
	// Talk to GCP resource manager. This field would be nil in offline mode.
	resourceManager *cloudresourcemanager.Service
	// Cache to prevent multiple network calls for looking up the same project's ancestry
	// map[project]ancestryPath
	ancestryCache map[string]string
}

// GetAncestry uses the resource manager API to get ancestry paths for
// projects. It implements a cache because many resources share the same
// project.
func (m *onlineAncestryManager) GetAncestry(project string) (string, error) {
	if path, ok := m.ancestryCache[project]; ok {
		return path, nil
	}
	ancestry, err := m.resourceManager.Projects.GetAncestry(project, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		return "", err
	}
	path := ancestryPath(ancestry.Ancestor)
	m.store(project, path)
	return path, nil
}

// GetAncestryWithResource first attempts to get Ancestry from the API
// If that fails, it falls back to inspecting the resource.
func (m *onlineAncestryManager) GetAncestryWithResource(project string, tfData converter.TerraformResourceData, cai converter.Asset) (string, error) {
	path, err := m.GetAncestry(project)
	if path != "" {
		return path, nil
	}

	ancestry, ok := m.getAncestryFromResource(tfData, cai)
	if !ok {
		return "", err
	}
	path = ancestryPath(ancestry)
	log.Printf("[INFO] [Online] Retrieved ancestry for %s: %s", project, path)
	m.store(project, path)
	return path, nil
}

func (m *onlineAncestryManager) store(project, ancestry string) {
	if project != "" && ancestry != "" {
		m.ancestryCache[project] = ancestry
	}
}

type offlineAncestryManager struct {
	resourceAncestryManager
	project  string
	ancestry string
}

// GetAncestry returns the ancestry for the project. It returns an error if
// the project does not equal to the one provided during initialization.
func (m *offlineAncestryManager) GetAncestry(project string) (string, error) {
	if project != m.project {
		return "", fmt.Errorf("cannot fetch ancestry in offline mode")
	}
	return m.ancestry, nil
}

// GetAncestryWithResource first attempts to get Ancestry from the resource
// If that fails, it falls back to the offline cache.
func (m *offlineAncestryManager) GetAncestryWithResource(project string, tfData converter.TerraformResourceData, cai converter.Asset) (string, error) {
	ancestry, ok := m.getAncestryFromResource(tfData, cai)
	if ok {
		path := ancestryPath(ancestry)
		log.Printf("[INFO] [Offline] Retrieved ancestry for %s: %s", project, path)
		return path, nil
	}

	path, err := m.GetAncestry(project)
	if err != nil {
		return "", err
	}

	return path, nil
}

// New returns AncestryManager that can be used to fetch ancestry information for a project.
func New(retriever ClientRetriever, project, ancestry, userAgent string, offline bool) (AncestryManager, error) {
	if ancestry != "" {
		ancestry = fmt.Sprintf("%s/project/%s", ancestry, project)
	}
	if offline {
		return &offlineAncestryManager{project: project, ancestry: ancestry}, nil
	}
	am := &onlineAncestryManager{ancestryCache: map[string]string{}}
	am.store(project, ancestry)
	rm := retriever.NewResourceManagerClient(userAgent)
	am.resourceManager = rm
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
