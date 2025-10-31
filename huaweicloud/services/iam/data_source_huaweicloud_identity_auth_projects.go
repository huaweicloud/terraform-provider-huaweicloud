package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityAuthProjects
// @API IAM GET /v3/auth/projects
func DataSourceIdentityAuthProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityAuthProjectsRead,

		Schema: map[string]*schema.Schema{
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

func DataSourceIdentityAuthProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	authProjectsPath := iamClient.Endpoint + "v3/auth/projects"
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", authProjectsPath, &options)
	if err != nil {
		return diag.Errorf("error listAuthProjects: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	projectsBody := utils.PathSearch("projects", respBody, make([]interface{}, 0)).([]interface{})
	projects := flattenProjects(projectsBody)
	if err = d.Set("projects", projects); err != nil {
		return diag.Errorf("error setting projects fields: %s", err)
	}
	return nil
}

func flattenProjects(projectsModel []interface{}) []map[string]interface{} {
	if projectsModel == nil {
		return nil
	}
	projects := make([]map[string]interface{}, 0, len(projectsModel))
	for _, project := range projectsModel {
		projects = append(projects, map[string]interface{}{
			"id":      utils.PathSearch("id", project, nil),
			"name":    utils.PathSearch("name", project, nil),
			"enabled": utils.PathSearch("enabled", project, nil),
		})
	}
	return projects
}
