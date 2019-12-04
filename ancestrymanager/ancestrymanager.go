// Package ancestrymanager provides an interface to query the ancestry information for a project.
package ancestrymanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"

	converter "github.com/GoogleCloudPlatform/terraform-google-conversion/google"
)

// AncestryManager is the interface that wraps the GetAncestry method.
type AncestryManager interface {
	// GetAncestry takes a project name and return a ancestry path
	GetAncestry(project string, options ...GetAncestryOption) (string, error)
}

type getAncestryRequest struct {
	project     string
	hasResource bool
	tfData      converter.TerraformResourceData
	hasAsset    bool
	asset       converter.Asset
}

// GetAncestryOption sets options for an ancestry tequest
type GetAncestryOption func(*getAncestryRequest)

// TerraformResource sets the Terraform resource for a getAncestry request
func TerraformResource(tfData converter.TerraformResourceData) GetAncestryOption {
	return func(req *getAncestryRequest) {
		req.hasResource = true
		req.tfData = tfData
	}
}

// Asset sets the Terraform resource for a getAncestry request
func Asset(cai converter.Asset) GetAncestryOption {
	return func(req *getAncestryRequest) {
		req.hasAsset = true
		req.asset = cai
	}
}

type resourceAncestryManager struct {
}

func (m *resourceAncestryManager) getFolderAncestry(folder_id string) ([]*cloudresourcemanager.Ancestor, error) {
	// TODO(morgantep): Incorporate folders.GetAncestry from v2alpha1 API
	log.Printf("[INFO] Retrieve ancestry for folder: %s", folder_id)

	return []*cloudresourcemanager.Ancestor{
		&cloudresourcemanager.Ancestor{
			ResourceId: &cloudresourcemanager.ResourceId{
				Type: "folder",
				Id:   folder_id,
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

func (m *resourceAncestryManager) getAncestryFromResource(req *getAncestryRequest) ([]*cloudresourcemanager.Ancestor, bool) {
	if !req.hasAsset || !req.hasResource {
		log.Printf("[INFO] %s does not have resource data", req.project)
		return nil, false
	}

	log.Printf("[INFO] Retrieving ancestry for %s from resource (type=%s)", req.project, req.asset.Type)

	switch req.asset.Type {
	case "cloudresourcemanager.googleapis.com/Project", "cloudbilling.googleapis.com/ProjectBillingInfo":
		ancestry := []*cloudresourcemanager.Ancestor{
			&cloudresourcemanager.Ancestor{
				ResourceId: &cloudresourcemanager.ResourceId{
					Type: "project",
					Id:   req.asset.Resource.Data["projectId"].(string),
				},
			},
		}

		org_id, ok := req.tfData.GetOk("org_id")
		if ok && org_id != "" {
			s := strings.Split(org_id.(string), "/")
			return append(ancestry, &cloudresourcemanager.Ancestor{
				ResourceId: &cloudresourcemanager.ResourceId{
					Type: "organization",
					Id:   s[len(s)-1],
				},
			}), true
		}

		folder_id, ok := req.tfData.GetOk("folder_id")
		if ok && folder_id != "" {
			folderAncestry, err := m.getFolderAncestry(folder_id.(string))
			if err != nil {
				log.Printf("[ERROR] Failed to retrieve folder ancestry: %s", err)
				return nil, false
			}
			return append(ancestry, folderAncestry...), true
		}
		return nil, true
	default:
		log.Printf("[INFO] Resource %s of type %s does not include sufficient data for ancestry retrieval", req.project, req.asset.Type)
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
func (m *onlineAncestryManager) GetAncestry(project string, options ...GetAncestryOption) (string, error) {
	req := new(getAncestryRequest)
	req.project = project

	for _, option := range options {
		option(req)
	}

	if path, ok := m.ancestryCache[project]; ok {
		return path, nil
	}
	response, err := m.resourceManager.Projects.GetAncestry(project, &cloudresourcemanager.GetAncestryRequest{}).Do()
	var ancestry []*cloudresourcemanager.Ancestor
	if err != nil {
		rAncestry, ok := m.getAncestryFromResource(req)
		if !ok {
			return "", err
		}
		ancestry = rAncestry
	} else {
		ancestry = response.Ancestor
	}
	path := ancestryPath(ancestry)
	log.Printf("[INFO] Retrieved ancestry for %s: %s", project, path)
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
func (m *offlineAncestryManager) GetAncestry(project string, options ...GetAncestryOption) (string, error) {
	req := new(getAncestryRequest)
	req.project = project

	for _, option := range options {
		option(req)
	}

	if req.project == m.project {
		return m.ancestry, nil
	}

	// attempt to retrieve ancestry from resource
	ancestry, ok := m.getAncestryFromResource(req)
	if ok {
		path := ancestryPath(ancestry)
		log.Printf("[INFO] Retrieved ancestry for %s from resource: %s", project, path)
		return path, nil
	}

	return "", fmt.Errorf("cannot fetch ancestry in offline mode")
}

// New returns AncestryManager that can be used to fetch ancestry information for a project.
func New(ctx context.Context, project, ancestry string, offline bool, opts ...option.ClientOption) (AncestryManager, error) {
	if ancestry != "" {
		ancestry = fmt.Sprintf("%s/project/%s", ancestry, project)
	}
	if offline {
		return &offlineAncestryManager{project: project, ancestry: ancestry}, nil
	}
	am := &onlineAncestryManager{ancestryCache: map[string]string{}}
	am.store(project, ancestry)
	rm, err := cloudresourcemanager.NewService(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "constructing resource manager client")
	}
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
