package iam

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

func DataSourceIdentityProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenProjectList(projectList []projects.Project) ([]map[string]interface{}, []string) {
	if len(projectList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(projectList))
	ids := make([]string, len(projectList))
	for i, val := range projectList {
		ids[i] = val.ID
		result[i] = map[string]interface{}{
			"id":      val.ID,
			"name":    val.Name,
			"enabled": val.Enabled,
		}
	}
	return result, ids
}

func dataSourceIdentityProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	listOpts := projects.ListOpts{
		Name: d.Get("name").(string),
	}
	pages, err := projects.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("error retrieving project list: %v", err)
	}
	projectList, err := projects.ExtractProjects(pages)
	if err != nil {
		return diag.Errorf("error fetching project objects: %v", err)
	}

	result, ids := flattenProjectList(projectList)

	d.SetId(hashcode.Strings(ids))
	return diag.FromErr(d.Set("projects", result))
}
