package codeartspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/pipelines-modify-historys
func DataSourceCodeArtsPipelineModifyHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineModifyHistoriesRead,

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
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"histories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the history list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"modify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the modify type.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"creator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator name.`,
						},
						"creator_nick_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator nick name.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineModifyHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/pipelines-modify-historys?limit=10"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error getting modify histories: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error getting modify histories: %s", err)
		}

		histories := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(histories) == 0 {
			break
		}

		for _, history := range histories {
			rst = append(rst, map[string]interface{}{
				"create_time":       utils.PathSearch("create_time", history, nil),
				"creator_name":      utils.PathSearch("creator_name", history, nil),
				"creator_nick_name": utils.PathSearch("creator_nick_name", history, nil),
				"modify_type":       utils.PathSearch("modify_type", history, nil),
			})
		}

		offset += 10
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("histories", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
