package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// DataSourceIdentityProjects
// @API IAM GET /v3/projects
// @API IAM GET /v3/projects/{project_id}
func DataSourceIdentityProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the IAM project name to query",
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Specifies the IAM project id to query. This parameter conflicts with `name`.",
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
	cfg := meta.(*config.Config)
	client, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	var result []map[string]interface{}
	var ids []string
	if projectId := d.Get("project_id").(string); projectId != "" {
		project, err := projects.Get(client, projectId).Extract()
		if err != nil {
			return diag.Errorf("error retrieving IAM project: %v", err)
		}
		result, ids = flattenProjectList([]projects.Project{*project})
	} else {
		listOpts := projects.ListOpts{
			Name: d.Get("name").(string),
		}
		pages, err := projects.List(client, listOpts).AllPages()
		if err != nil {
			return diag.Errorf("error retrieving IAM project list: %v", err)
		}
		projectList, err := projects.ExtractProjects(pages)
		if err != nil {
			return diag.Errorf("error fetching IAM project objects: %v", err)
		}
		result, ids = flattenProjectList(projectList)
	}
	d.SetId(hashcode.Strings(ids))
	return diag.FromErr(d.Set("projects", result))
}
