// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------
package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var PrivatecaCertificateTemplateIamSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"location": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"certificate_template": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: compareSelfLinkOrResourceName,
	},
}

type PrivatecaCertificateTemplateIamUpdater struct {
	project             string
	location            string
	certificateTemplate string
	d                   TerraformResourceData
	Config              *Config
}

func PrivatecaCertificateTemplateIamUpdaterProducer(d TerraformResourceData, config *Config) (ResourceIamUpdater, error) {
	values := make(map[string]string)

	project, _ := getProject(d, config)
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
	}
	values["project"] = project
	location, _ := getLocation(d, config)
	if location != "" {
		if err := d.Set("location", location); err != nil {
			return nil, fmt.Errorf("Error setting location: %s", err)
		}
	}
	values["location"] = location
	if v, ok := d.GetOk("certificate_template"); ok {
		values["certificate_template"] = v.(string)
	}

	// We may have gotten either a long or short name, so attempt to parse long name if possible
	m, err := getImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/certificateTemplates/(?P<certificate_template>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<certificate_template>[^/]+)", "(?P<location>[^/]+)/(?P<certificate_template>[^/]+)"}, d, config, d.Get("certificate_template").(string))
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &PrivatecaCertificateTemplateIamUpdater{
		project:             values["project"],
		location:            values["location"],
		certificateTemplate: values["certificate_template"],
		d:                   d,
		Config:              config,
	}

	if err := d.Set("project", u.project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("location", u.location); err != nil {
		return nil, fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("certificate_template", u.GetResourceId()); err != nil {
		return nil, fmt.Errorf("Error setting certificate_template: %s", err)
	}

	return u, nil
}

func PrivatecaCertificateTemplateIdParseFunc(d *schema.ResourceData, config *Config) error {
	values := make(map[string]string)

	project, _ := getProject(d, config)
	if project != "" {
		values["project"] = project
	}

	location, _ := getLocation(d, config)
	if location != "" {
		values["location"] = location
	}

	m, err := getImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/certificateTemplates/(?P<certificate_template>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<certificate_template>[^/]+)", "(?P<location>[^/]+)/(?P<certificate_template>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &PrivatecaCertificateTemplateIamUpdater{
		project:             values["project"],
		location:            values["location"],
		certificateTemplate: values["certificate_template"],
		d:                   d,
		Config:              config,
	}
	if err := d.Set("certificate_template", u.GetResourceId()); err != nil {
		return fmt.Errorf("Error setting certificate_template: %s", err)
	}
	d.SetId(u.GetResourceId())
	return nil
}

func (u *PrivatecaCertificateTemplateIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url, err := u.qualifyCertificateTemplateUrl("getIamPolicy")
	if err != nil {
		return nil, err
	}

	project, err := getProject(u.d, u.Config)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	url, err = addQueryParams(url, map[string]string{"options.requestedPolicyVersion": fmt.Sprintf("%d", iamPolicyVersion)})
	if err != nil {
		return nil, err
	}

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return nil, err
	}

	policy, err := sendRequest(u.Config, "GET", project, url, userAgent, obj)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	out := &cloudresourcemanager.Policy{}
	err = Convert(policy, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a policy to a resource manager policy: {{err}}", err)
	}

	return out, nil
}

func (u *PrivatecaCertificateTemplateIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	json, err := ConvertToMap(policy)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	obj["policy"] = json

	url, err := u.qualifyCertificateTemplateUrl("setIamPolicy")
	if err != nil {
		return err
	}
	project, err := getProject(u.d, u.Config)
	if err != nil {
		return err
	}

	userAgent, err := generateUserAgentString(u.d, u.Config.userAgent)
	if err != nil {
		return err
	}

	_, err = sendRequestWithTimeout(u.Config, "POST", project, url, userAgent, obj, u.d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *PrivatecaCertificateTemplateIamUpdater) qualifyCertificateTemplateUrl(methodIdentifier string) (string, error) {
	urlTemplate := fmt.Sprintf("{{PrivatecaBasePath}}%s:%s", fmt.Sprintf("projects/%s/locations/%s/certificateTemplates/%s", u.project, u.location, u.certificateTemplate), methodIdentifier)
	url, err := replaceVars(u.d, u.Config, urlTemplate)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *PrivatecaCertificateTemplateIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/locations/%s/certificateTemplates/%s", u.project, u.location, u.certificateTemplate)
}

func (u *PrivatecaCertificateTemplateIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-privateca-certificatetemplate-%s", u.GetResourceId())
}

func (u *PrivatecaCertificateTemplateIamUpdater) DescribeResource() string {
	return fmt.Sprintf("privateca certificatetemplate %q", u.GetResourceId())
}
