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

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-tag/list
func DataSourceCodeArtsPipelineTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the tag list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tag ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tag name.`,
						},
						"color": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tag color.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineTag(client, d.Get("project_id").(string))
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
		d.Set("tags", flattenDataSourcePipelineTagsStages(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourcePipelineTagsStages(resp interface{}) []interface{} {
	if tags, ok := resp.([]interface{}); ok && len(tags) > 0 {
		result := make([]interface{}, 0, len(tags))
		for _, v := range tags {
			tag := v.(map[string]interface{})
			m := map[string]interface{}{
				"id":    utils.PathSearch("tag_id", tag, nil),
				"name":  utils.PathSearch("name", tag, nil),
				"color": utils.PathSearch("color", tag, nil),
			}
			result = append(result, m)
		}

		return result
	}

	return nil
}
