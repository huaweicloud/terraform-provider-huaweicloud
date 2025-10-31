package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityFederationProjects
// @API IAM GET /v3/OS-FEDERATION/projects
func DataSourceIdentityFederationProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityFederationProjectsRead,

		Schema: map[string]*schema.Schema{
			"federation_token": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceIdentityFederationProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	federationProjectsPath := client.ResourceBase + "v3/OS-FEDERATION/projects"
	federationProjectsPath = strings.ReplaceAll(federationProjectsPath, "{region_id}", cfg.GetRegion(d))
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("federation_token").(string),
		},
	}
	response, err := client.Request("GET", federationProjectsPath, &options)
	if err != nil {
		return diag.Errorf("error federationProjects: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate UUID: %s", err)
	}
	d.SetId(id)
	projectsBody := utils.PathSearch("projects", respBody, make([]interface{}, 0)).([]interface{})
	projects := flattenProjects(projectsBody)
	if err = d.Set("projects", projects); err != nil {
		return diag.Errorf("error set projects filed: %s", err)
	}
	return nil
}
