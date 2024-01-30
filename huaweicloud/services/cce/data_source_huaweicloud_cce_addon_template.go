package cce

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"
	"github.com/chnsz/golangsdk/openstack/cce/v3/templates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type Template struct {
	UID             string
	Name            string
	Version         string
	Description     string
	Spec            string
	Stable          bool
	SupportVersions []addons.SupportVersions
}

func (t Template) IsEmpty() bool {
	return reflect.DeepEqual(t, Template{})
}

// @API CCE GET /api/v3/addontemplates
func DataSourceAddonTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAddonTemplateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_version": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"bare_metal": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAddonTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CceAddonV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE client : %s", err)
	}
	// Get all addon templates by List function
	cluster_id := d.Get("cluster_id").(string)
	templateList, err := templates.List(client, cluster_id).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve template list: %s", err)
	}

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	template, err := getTemplateByNameAndVersion(templateList, name, version)
	if err != nil {
		return diag.Errorf("unable to find specifies template by name (%s) and version (%s): %s", name, version, err)
	}

	d.SetId(template.UID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("spec", template.Spec),
		d.Set("description", template.Description),
		d.Set("stable", template.Stable),
		setTemplateSupportVersionState(d, template.SupportVersions),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting template fields: %s", err)
	}
	return nil
}

// getTemplateByNameAndVersion is method to using add-on name, version and the cluster id to find a unique
// add-on template in list.
func getTemplateByNameAndVersion(templateList []templates.Template, specName, specVersion string) (Template, error) {
	var result Template
	// For each add-on, they have one or more version templates.
	for _, temp := range templateList {
		if temp.Metadata.Name != specName {
			continue
		}
		// If the specified additional template is found, the loop is interrupted.
		for _, ver := range temp.Spec.Versions {
			if ver.Version == specVersion {
				// Save the description and
				result.UID = temp.Metadata.UID
				result.Name = specName
				result.Version = specVersion
				result.Description = temp.Spec.Description
				// Return a json string to the user, which contains the contents of the basic and custom fields.
				specBytes, err := json.Marshal(ver.Input)
				if err != nil {
					return result, fmt.Errorf("error converting input struct")
				}
				result.Spec = string(specBytes)
				result.Stable = ver.Stable
				result.SupportVersions = ver.SupportVersions
				break
			}
		}
		if !result.IsEmpty() {
			break
		}
	}
	if result.IsEmpty() {
		return result, fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return result, nil
}

func setTemplateSupportVersionState(d *schema.ResourceData, supportList []addons.SupportVersions) error {
	serportVersionMap := map[string]*schema.Set{}
	for _, supports := range supportList {
		v := schema.Set{F: schema.HashString}
		for _, ver := range supports.ClusterVersion {
			v.Add(ver)
		}
		if supports.ClusterType == "VirtualMachine" {
			serportVersionMap["virtual_machine"] = &v
		}
		if supports.ClusterType == "BareMetal" {
			serportVersionMap["bare_metal"] = &v
		}
	}
	if err := d.Set("support_version", []map[string]*schema.Set{serportVersionMap}); err != nil {
		return err
	}
	return nil
}
