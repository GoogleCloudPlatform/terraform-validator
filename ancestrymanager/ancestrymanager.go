// Package ancestrymanager provides an interface to query the ancestry information for a project.
package ancestrymanager

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

// AncestryManager is the interface that wraps the GetAncestry method.
type AncestryManager interface {
	// GetAncestry takes a project name and return a ancestry path
	GetAncestry(project string) (string, error)
}

type ancestryManager struct {
	// Talk to GCP resource manager. This field would be nil in offline mode.
	resourceManager *cloudresourcemanager.Service
	// Cache to prevent multiple network calls for looking up the same project's ancestry
	// map[project]ancestryPath
	ancestryCache map[string]string
}

// GetAncestry uses the resource manager API to get ancestry paths for
// projects. It implements a cache because many resources share the same
// project.
func (am *ancestryManager) GetAncestry(project string) (string, error) {
	if path, ok := am.ancestryCache[project]; ok {
		return path, nil
	}
	if am.resourceManager == nil {
		return "", fmt.Errorf("cannot fetch ancestry in offline mode for project %s", project)
	}
	ancestry, err := am.resourceManager.Projects.GetAncestry(project, &cloudresourcemanager.GetAncestryRequest{}).Do()
	if err != nil {
		return "", err
	}
	path := ancestryPath(ancestry.Ancestor)
	am.store(project, path)
	return path, nil
}

func (am *ancestryManager) store(project, ancestry string) {
	if project != "" && ancestry != "" {
		am.ancestryCache[project] = ancestry
	}
}

// New returns AncestryManager that can be used to fetch ancestry information for a project.
func New(ctx context.Context, project, ancestry string, offline bool, opts ...option.ClientOption) (AncestryManager, error) {
	am := &ancestryManager{
		ancestryCache: map[string]string{},
	}
	am.store(project, ancestry)
	if !offline {
		rm, err := cloudresourcemanager.NewService(ctx, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "constructing resource manager client")
		}
		am.resourceManager = rm
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
