package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v3/{project_id}/jobs/batch-progress
func DataSourceDrsBatchProgresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsBatchProgressesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incre_trans_delay": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incre_trans_delay_millis": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transfer_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remaining_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress_map": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"completed": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remaining_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBatchProgressesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"jobs": utils.ExpandToStringList(d.Get("job_ids").([]interface{})),
	}
}

func dataSourceDrsBatchProgressesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/batch-progress"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildBatchProgressesBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS batch progresses: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("results", flattenBatchProgressesResults(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchProgressesResults(respBody interface{}) []interface{} {
	resultsRaw := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(resultsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resultsRaw))
	for _, item := range resultsRaw {
		progressMap := utils.PathSearch("progress_map", item, make(map[string]interface{})).(map[string]interface{})
		result = append(result, map[string]interface{}{
			"job_id":                   utils.PathSearch("job_id", item, nil),
			"progress":                 utils.PathSearch("progress", item, nil),
			"incre_trans_delay":        utils.PathSearch("incre_trans_delay", item, nil),
			"incre_trans_delay_millis": utils.PathSearch("incre_trans_delay_millis", item, nil),
			"task_mode":                utils.PathSearch("task_mode", item, nil),
			"transfer_status":          utils.PathSearch("transfer_status", item, nil),
			"process_time":             utils.PathSearch("process_time", item, nil),
			"remaining_time":           utils.PathSearch("remaining_time", item, nil),
			"progress_map":             flattenProgressMap(progressMap),
			"error_code":               utils.PathSearch("error_code", item, nil),
			"error_msg":                utils.PathSearch("error_msg", item, nil),
		})
	}
	return result
}

func flattenProgressMap(progressMapRaw map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(progressMapRaw))
	for key, value := range progressMapRaw {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"key":            key,
			"completed":      utils.PathSearch("completed", item, nil),
			"remaining_time": utils.PathSearch("remaining_time", item, nil),
		})
	}
	return result
}
