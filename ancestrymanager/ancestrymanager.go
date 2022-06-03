// Package ancestrymanager provides an interface to query the ancestry information for a resource.
package ancestrymanager

import (
	"fmt"
	"strings"

	"google.golang.org/api/cloudresourcemanager/v3"

	resources "github.com/GoogleCloudPlatform/terraform-validator/converters/google/resources"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

// AncestryManager is the interface that fetch ancestors for a resource.
type AncestryManager interface {
	// Ancestors returns a list of ancestors.
	Ancestors(config *resources.Config, tfData resources.TerraformResourceData, cai *resources.Asset) ([]string, string, error)
}

type manager struct {
	// The logger.
	errorLogger *zap.Logger
	// GCP resource manager service. If this field is nil, online lookups will
	// be disabled.
	resourceManager *cloudresourcemanager.Service
	// Cache to prevent multiple network calls for looking up the same
	// resource's ancestry. The map key is the resource itself, in the format of
	// "<type>/<id>", ancestors are sorted from closest to furthest.
	ancestorCache map[string][]string
}

// New returns AncestryManager that can be used to fetch ancestry information.
// Entries takes `projects/<number>` or `folders/<id>` as key and ancestry path
// as value to the offline cache. If the key is not prefix with `projects/` or
// `folders/`, it will be considered as a project. If no resourceManager is
// provided, API requests for ancestry will be disabled.
func New(resourceManager *cloudresourcemanager.Service, entries map[string]string, errorLogger *zap.Logger) (AncestryManager, error) {
	return newManager(resourceManager, entries, errorLogger)
}

func newManager(resourceManager *cloudresourcemanager.Service, entries map[string]string, errorLogger *zap.Logger) (*manager, error) {
	am := &manager{
		ancestorCache:   map[string][]string{},
		errorLogger:     errorLogger,
		resourceManager: resourceManager,
	}
	err := am.initAncestryCache(entries)
	if err != nil {
		return nil, err
	}
	return am, nil
}

func (m *manager) initAncestryCache(entries map[string]string) error {
	for item, ancestry := range entries {
		if item != "" && ancestry != "" {
			ancestors, err := parseAncestryPath(ancestry)
			if err != nil {
				continue
			}
			key, err := parseAncestryKey(item)
			if err != nil {
				return err
			}
			// ancestry path should include the item itself
			if ancestors[0] != key {
				ancestors = append([]string{key}, ancestors...)
			}
			m.store(key, ancestors)
		}
	}
	return nil
}

func parseAncestryKey(val string) (string, error) {
	key := normalizeAncestry(val)
	ix := strings.LastIndex(key, "/")
	if ix == -1 {
		// If not containing /, then treat it as a project.
		return fmt.Sprintf("projects/%s", key), nil
	} else {
		k := key[:ix]
		if k == "projects" || k == "folders" || k == "organizations" {
			return key, nil
		}
		return "", fmt.Errorf("key with can only start with projects/, folders/, or organizations/")
	}
}

// Ancestors uses the resource manager API to get ancestors for resource.
// It implements a cache because many resources share the same ancestors.
func (m *manager) Ancestors(config *resources.Config, tfData resources.TerraformResourceData, cai *resources.Asset) ([]string, string, error) {
	results, err := m.fetchAncestors(config, tfData, cai)
	if err != nil {
		return nil, "", err
	}

	parent, err := assetParent(cai, results)
	if err != nil {
		return nil, "", err
	}
	return results, parent, nil
}

// fetchAncestors uses the resource manager API to get ancestors for resource.
// It implements a cache because many resources share the same ancestors.
func (m *manager) fetchAncestors(config *resources.Config, tfData resources.TerraformResourceData, cai *resources.Asset) ([]string, error) {
	if cai == nil {
		return nil, fmt.Errorf("CAI asset is nil")
	}
	m.errorLogger.Info(fmt.Sprintf("Retrieving ancestry from resource (type=%s)", cai.Type))
	key := ""
	orgKey := ""
	folderKey := ""

	orgID, orgOK := getOrganizationFromResource(tfData)
	if orgOK {
		orgKey = orgID
		if !strings.HasPrefix(orgKey, "organizations/") {
			orgKey = fmt.Sprintf("organizations/%s", orgKey)
		}
	}
	folderID, folderOK := getFolderFromResource(tfData)
	if folderOK {
		folderKey = folderID
		if !strings.HasPrefix(folderKey, "folders/") {
			folderKey = fmt.Sprintf("folders/%s", folderKey)
		}
	}
	switch cai.Type {
	case "cloudresourcemanager.googleapis.com/Folder":
		if !folderOK {
			return nil, fmt.Errorf("folder id not found in terraform data")
		}
		key = folderKey
	case "cloudresourcemanager.googleapis.com/Organization":
		if !orgOK {
			return nil, fmt.Errorf("organization id not found in terraform data")
		}
		key = orgKey
	case "iam.googleapis.com/Role":
		// google_organization_iam_custom_role or google_project_iam_custom_role
		if orgOK {
			key = orgKey
		} else {
			project, err := getProjectFromResource(tfData, config, *cai, m.errorLogger)
			if err != nil {
				return nil, err
			}
			key = fmt.Sprintf("projects/%s", project)
		}
	case "cloudresourcemanager.googleapis.com/Project", "cloudbilling.googleapis.com/ProjectBillingInfo":
		projectID, err := getProjectFromResource(tfData, config, *cai, m.errorLogger)
		if err != nil {
			return nil, err
		}

		// Changing project_id forces a new project to be created.
		// Changing folder_id or org_id forces project to be migrated.
		// Hence ancestors should not be fetched from API in those scenarios.
		if orgOK {
			ancestors := []string{
				fmt.Sprintf("projects/%s", projectID),
				fmt.Sprintf("organizations/%s", orgID),
			}
			return ancestors, nil
		}
		if folderOK {
			key = folderKey
			ret, err := m.getAncestorsWithCache(key)
			if err != nil {
				return nil, err
			}
			ancestors := append([]string{fmt.Sprintf("projects/%s", projectID)}, ret...)
			return ancestors, nil
		}
		key = fmt.Sprintf("projects/%s", projectID)
		ancestors, err := m.getAncestorsWithCache(key)
		if err != nil {
			return nil, err
		}
		return ancestors, nil
	default:
		project, err := getProjectFromResource(tfData, config, *cai, m.errorLogger)
		if err != nil {
			return nil, err
		}
		key = fmt.Sprintf("projects/%s", project)
	}
	return m.getAncestorsWithCache(key)
}

func (m *manager) getAncestorsWithCache(key string) ([]string, error) {
	var ancestors []string
	cur := key
	for cur != "" {
		if cachedAncestors, ok := m.ancestorCache[cur]; ok {
			ancestors = append(ancestors, cachedAncestors...)
			break
		}
		if strings.HasPrefix(cur, "organizations/") {
			ancestors = append(ancestors, cur)
			break
		}
		if m.resourceManager == nil {
			return nil, fmt.Errorf("resourceManager required to fetch ancestry for %s from the API", cur)
		}
		project, err := m.resourceManager.Projects.Get(cur).Do()
		if err != nil {
			return nil, errors.Wrap(err, "for more info- https://cloud.google.com/docs/terraform/policy-validation/troubleshooting#ProjectCallerForbidden")
		}
		ancestors = append(ancestors, project.Name)
		cur = project.Parent
	}
	m.store(key, ancestors)
	return ancestors, nil
}

func (m *manager) store(key string, ancestors []string) {
	if key == "" || len(ancestors) == 0 {
		return
	}
	if _, ok := m.ancestorCache[key]; !ok {
		m.ancestorCache[key] = ancestors
	}
	// cache ancestors along the ancestry path
	for i, ancestor := range ancestors {
		if _, ok := m.ancestorCache[ancestor]; !ok {
			m.ancestorCache[ancestor] = ancestors[i:]
		}
	}
}

func parseAncestryPath(path string) ([]string, error) {
	normStr := normalizeAncestry(path)
	splits := strings.Split(normStr, "/")
	if len(splits)%2 != 0 {
		return nil, fmt.Errorf("unexpected format of ancestry path %s", path)
	}
	var ancestors []string
	allowedPrefixes := map[string]bool{
		"projects":      true,
		"folders":       true,
		"organizations": true,
	}
	for i := 0; i < len(splits); i = i + 2 {
		if _, ok := allowedPrefixes[splits[i]]; !ok {
			return nil, fmt.Errorf("invalid ancestry path %s with %s", path, splits[i])
		}
		ancestors = append(ancestors, fmt.Sprintf("%s/%s", splits[i], splits[i+1]))
	}
	// reverse the sequence
	i, j := 0, len(ancestors)-1
	for i < j {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
		i++
		j--
	}
	return ancestors, nil
}

func normalizeAncestry(val string) string {
	for _, r := range []struct {
		old string
		new string
	}{
		{"organization/", "organizations/"},
		{"folder/", "folders/"},
		{"project/", "projects/"},
	} {
		val = strings.ReplaceAll(val, r.old, r.new)
	}
	return val
}

type NoOpAncestryManager struct{}

func (*NoOpAncestryManager) Ancestors(config *resources.Config, tfData resources.TerraformResourceData, cai *resources.Asset) ([]string, string, error) {
	return nil, "", nil
}
