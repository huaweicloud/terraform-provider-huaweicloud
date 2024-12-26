package codeartsdeploy

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy GET /v1/projects/{project_id}/applications/groups
func DataSourceCodeartsDeployApplicationGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsDeployApplicationGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the application group list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application group name.`,
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parent application group ID.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group path.`,
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the group sorting field.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the group creator.`,
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the user who last updates the group.`,
						},
						"application_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total number of applications in the group.`,
						},
						"children": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the child group name list.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsDeployApplicationGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	groups, err := getDeployApplicationGroup(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenCodeartsDeployApplicationGroups(groups.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCodeartsDeployApplicationGroups(currentList []interface{}) []map[string]interface{} {
	groups := currentList
	rst := make([]map[string]interface{}, 0)

	for _, group := range groups {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", group, nil),
			"name":              utils.PathSearch("name", group, nil),
			"parent_id":         utils.PathSearch("parent_id", group, nil),
			"path":              utils.PathSearch("path", group, nil),
			"ordinal":           utils.PathSearch("ordinal", group, nil),
			"application_count": utils.PathSearch("application_count", group, nil),
			"children":          utils.PathSearch("children[*].name", group, nil),
			"updated_by":        utils.PathSearch("last_update_user_id", group, nil),
			"created_by":        utils.PathSearch("create_user_id", group, nil),
		})
	}

	children := utils.PathSearch("[*].children[]", currentList, make([]interface{}, 0)).([]interface{})
	if len(children) != 0 {
		rst = append(rst, flattenCodeartsDeployApplicationGroups(children)...)
	}

	return rst
}
