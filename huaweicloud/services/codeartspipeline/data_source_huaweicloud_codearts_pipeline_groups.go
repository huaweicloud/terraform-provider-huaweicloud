package codeartspipeline

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy GET /v5/{project_id}/api/pipeline-group/tree
func DataSourceCodeArtsPipelineGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineGroupsRead,

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
				Description: `Specifies the CodeArts project ID.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the pipeline groups list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group name.`,
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parent group ID.`,
						},
						"path_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group path.`,
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the group sorting field.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the group creator.`,
						},
						"updater": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the user who last updates the group.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the update time.`,
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

func dataSourceCodeArtsPipelineGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts pipeline client: %s", err)
	}

	groups, err := GetPipelineGroups(client, d.Get("project_id").(string))
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
		d.Set("groups", flattenCodeartsPipelineGroups(groups.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCodeartsPipelineGroups(currentList []interface{}) []map[string]interface{} {
	groups := currentList
	rst := make([]map[string]interface{}, 0)

	for _, group := range groups {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", group, nil),
			"name":        utils.PathSearch("name", group, nil),
			"parent_id":   utils.PathSearch("parent_id", group, nil),
			"path_id":     utils.PathSearch("path_id", group, nil),
			"ordinal":     utils.PathSearch("ordinal", group, nil),
			"children":    utils.PathSearch("children[*].name", group, nil),
			"updater":     utils.PathSearch("updater", group, nil),
			"creator":     utils.PathSearch("creator", group, nil),
			"create_time": utils.PathSearch("create_time", group, nil),
			"update_time": utils.PathSearch("update_time", group, nil),
		})
	}

	children := utils.PathSearch("[*].children[]", currentList, make([]interface{}, 0)).([]interface{})
	if len(children) != 0 {
		rst = append(rst, flattenCodeartsPipelineGroups(children)...)
	}

	return rst
}
