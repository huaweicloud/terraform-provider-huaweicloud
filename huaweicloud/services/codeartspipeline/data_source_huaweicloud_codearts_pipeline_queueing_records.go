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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/queued-pipeline
func DataSourceCodeArtsPipelineQueueingRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineQueueingRecordsRead,

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
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the record list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the record ID.`,
						},
						"pipeline_run_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the pipeline run ID.`,
						},
						"enqueue_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the enqueuing time.`,
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger type.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status.`,
						},
						"creator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator ID.`,
						},
						"creator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator name.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineQueueingRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/queued-pipeline?limit=10"
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
			return diag.Errorf("error getting ququeing records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error getting ququeing records: %s", err)
		}

		records := utils.PathSearch("queuedPipelines", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		for _, record := range records {
			rst = append(rst, map[string]interface{}{
				"id":              utils.PathSearch("id", record, nil),
				"pipeline_run_id": utils.PathSearch("pipeline_run_id", record, nil),
				"enqueue_time":    utils.PathSearch("enqueue_time", record, nil),
				"trigger_type":    utils.PathSearch("trigger_type", record, nil),
				"status":          utils.PathSearch("status", record, nil),
				"creator_id":      utils.PathSearch("creator_id", record, nil),
				"creator_name":    utils.PathSearch("creator_name", record, nil),
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
		d.Set("records", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
