package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var IamBigQueryDatasetSchema = map[string]*schema.Schema{
	"dataset_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareBigQueryDatasetName,
	},
}

type BigQueryDatasetIamUpdater struct {
	resourceId string
	d          TerraformResourceData
	Config     *Config
}

func NewBigQueryDatasetIamUpdater(d TerraformResourceData, config *Config) (ResourceIamUpdater, error) {
	return &BigQueryDatasetIamUpdater{
		resourceId: d.Get("dataset_id").(string),
		d:          d,
		Config:     config,
	}, nil
}

func BigQueryDatasetIdParseFunc(d *schema.ResourceData, _ *Config) error {
	if err := d.Set("dataset_id", d.Id()); err != nil {
		return fmt.Errorf("Error setting dataset_id: %s", err)
	}
	return nil
}

func (u *BigQueryDatasetIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	dataset_idId := GetResourceNameFromSelfLink(u.resourceId)

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	p, err := u.Config.NewResourceManagerClient(userAgent).Projects.GetIamPolicy(dataset_idId,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: iamPolicyVersion,
			},
		}).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return p, nil
}

func (u *BigQueryDatasetIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	dataset_idId := GetResourceNameFromSelfLink(u.resourceId)

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewResourceManagerClient(userAgent).Projects.SetIamPolicy(dataset_idId,
		&cloudresourcemanager.SetIamPolicyRequest{
			Policy:     policy,
			UpdateMask: "bindings,etag,auditConfigs",
		}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BigQueryDatasetIamUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *BigQueryDatasetIamUpdater) GetMutexKey() string {
	return getBigQueryDatasetIamPolicyMutexKey(u.resourceId)
}

func (u *BigQueryDatasetIamUpdater) DescribeResource() string {
	return fmt.Sprintf("dataset_id %q", u.resourceId)
}

func compareBigQueryDatasetName(_, old, new string, _ *schema.ResourceData) bool {
	// We can either get "dataset_ids/dataset_id-id" or "dataset_id-id", so strip any prefixes
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

func getBigQueryDatasetIamPolicyMutexKey(pid string) string {
	return fmt.Sprintf("iam-dataset_id-%s", pid)
}
