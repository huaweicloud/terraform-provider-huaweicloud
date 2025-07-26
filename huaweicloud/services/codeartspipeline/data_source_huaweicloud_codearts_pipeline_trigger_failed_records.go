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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/trigger-failed-record
func DataSourceCodeArtsPipelineTriggerFailedRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineTriggerFailedRecordsRead,

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
				Description: `Indicates the records list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger type.`,
						},
						"trigger_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the trigger time.`,
						},
						"executor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the executor ID.`,
						},
						"executor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the executor name.`,
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the cause of trigger failure.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineTriggerFailedRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/trigger-failed-record?limit=10"
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
			return diag.Errorf("error getting trigger failed records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error getting trigger failed records: %s", err)
		}

		records := utils.PathSearch("triggerFailedRecordVOS", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		for _, record := range records {
			rst = append(rst, map[string]interface{}{
				"trigger_type":  utils.PathSearch("trigger_type", record, nil),
				"trigger_time":  utils.PathSearch("trigger_time", record, nil),
				"executor_id":   utils.PathSearch("executor_id", record, nil),
				"executor_name": utils.PathSearch("executor_name", record, nil),
				"reason":        utils.PathSearch("reason", record, nil),
			})
		}

		offset += len(records)
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
