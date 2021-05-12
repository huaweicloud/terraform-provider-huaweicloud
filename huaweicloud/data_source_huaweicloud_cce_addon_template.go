package huaweicloud

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/templates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceCCEAddonTemplateV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEAddonTemplateV3Read,

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
		},
	}
}

func dataSourceCCEAddonTemplateV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.CceAddonV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client : %s", err)
	}
	// Get all addon templates by List function
	cluster_id := d.Get("cluster_id").(string)
	templateList, err := templates.List(client, cluster_id).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve template list: %s", err)
	}

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	spec, description, err := getTemplateByNameAndVersion(templateList, name, version)
	if err != nil {
		return fmt.Errorf("Unable to find specifies template by name (%s) and version (%s): %s", name, version, err)
	}

	d.SetId(hashcode.Strings([]string{name, version}))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("spec", spec),
		d.Set("description", description),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getTemplateByNameAndVersion(templateList []templates.Template, name, version string) (string, string, error) {
	versions := make([]templates.Versions, 1)
	var description string
LOOP:
	for _, template := range templateList {
		if template.Metadata.Name != name {
			continue
		}
		for _, ver := range template.Spec.Versions {
			if version == ver.Version {
				description = template.Spec.Description
				versions[0] = ver
				break LOOP
			}
		}
	}
	if len(versions) < 1 {
		return "", "", fmt.Errorf("Your query returned no results, please change your search criteria and try again")
	}

	specStruct := versions[0].Input
	// Return a json string to the user, which contains the contents of the basic and custom fields (In add-on values).
	specBytes, err := json.Marshal(specStruct)
	if err != nil {
		return "", "", fmt.Errorf("Error converting input struct")
	}
	spec := string(specBytes)
	return spec, description, nil
}
