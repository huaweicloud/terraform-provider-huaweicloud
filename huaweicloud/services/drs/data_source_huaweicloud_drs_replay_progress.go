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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/replay-progress
func DataSourceDrsReplayProgress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsReplayProgressRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"progress": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"parse_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replay_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"task_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transfer_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"now_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_export_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_export_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replay_sql_now_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replaySqlNowListSchema(),
			},
		},
	}
}

func replaySqlNowListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"thread_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shard_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_statement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"execute_latency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDrsReplayProgressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/replay-progress"
		jobId   = d.Get("job_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS replay progress: %s", err)
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
		d.Set("progress", utils.PathSearch("progress", respBody, nil)),
		d.Set("parse_count", utils.PathSearch("parse_count", respBody, nil)),
		d.Set("replay_count", utils.PathSearch("replay_count", respBody, nil)),
		d.Set("task_mode", utils.PathSearch("task_mode", respBody, nil)),
		d.Set("process_time", utils.PathSearch("process_time", respBody, nil)),
		d.Set("transfer_status", utils.PathSearch("transfer_status", respBody, nil)),
		d.Set("max_time", utils.PathSearch("max_time", respBody, nil)),
		d.Set("min_time", utils.PathSearch("min_time", respBody, nil)),
		d.Set("now_time", utils.PathSearch("now_time", respBody, nil)),
		d.Set("min_export_time", utils.PathSearch("min_export_time", respBody, nil)),
		d.Set("max_export_time", utils.PathSearch("max_export_time", respBody, nil)),
		d.Set("replay_sql_now_list", flattenReplaySqlNowList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplaySqlNowList(respBody interface{}) []interface{} {
	replaySqlNowListRaw := utils.PathSearch("replay_sql_now_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(replaySqlNowListRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(replaySqlNowListRaw))
	for _, item := range replaySqlNowListRaw {
		result = append(result, map[string]interface{}{
			"thread_id":       utils.PathSearch("thread_id", item, nil),
			"created_at":      utils.PathSearch("created_at", item, nil),
			"modified_at":     utils.PathSearch("modified_at", item, nil),
			"shard_id":        utils.PathSearch("shard_id", item, nil),
			"schema_name":     utils.PathSearch("schema_name", item, nil),
			"sql_statement":   utils.PathSearch("sql_statement", item, nil),
			"latency":         utils.PathSearch("latency", item, nil),
			"execute_latency": utils.PathSearch("execute_latency", item, nil),
			"target_type":     utils.PathSearch("target_type", item, nil),
			"target_name":     utils.PathSearch("target_name", item, nil),
			"status":          utils.PathSearch("status", item, nil),
		})
	}
	return result
}
