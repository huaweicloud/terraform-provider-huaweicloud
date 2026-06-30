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

// @API DRS GET /v3/{project_id}/jobs/{job_id}/object/compare
func DataSourceDrsObjectCompare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectCompareRead,

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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"report_remain_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"compare_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"database_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disaster_recovery_database": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceObjectCompareRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v3/{project_id}/jobs/{job_id}/object/compare"
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
		return diag.Errorf("error retrieving DRS object compare: %s", err)
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
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("export_status", utils.PathSearch("export_status", respBody, nil)),
		d.Set("report_remain_seconds", utils.PathSearch("report_remain_seconds", respBody, nil)),
		d.Set("compare_job_id", utils.PathSearch("compare_job_id", respBody, nil)),
		d.Set("error_msg", utils.PathSearch("error_msg", respBody, nil)),
		d.Set("compare_result", flattenCompareResult(respBody)),
		d.Set("database_info", flattenDatabaseInfo(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCompareResult(respBody interface{}) []interface{} {
	compareResultRaw := utils.PathSearch("compare_result", respBody, make([]interface{}, 0)).([]interface{})
	if len(compareResultRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(compareResultRaw))
	for _, item := range compareResultRaw {
		result = append(result, map[string]interface{}{
			"type":         utils.PathSearch("type", item, nil),
			"source_count": utils.PathSearch("source_count", item, nil),
			"target_count": utils.PathSearch("target_count", item, nil),
			"status":       utils.PathSearch("status", item, nil),
		})
	}
	return result
}

func flattenDatabaseInfo(respBody interface{}) []interface{} {
	databaseInfoRaw := utils.PathSearch("database_info", respBody, nil)
	if databaseInfoRaw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"service_database":           utils.PathSearch("service_database", databaseInfoRaw, nil),
			"disaster_recovery_database": utils.PathSearch("disaster_recovery_database", databaseInfoRaw, nil),
		},
	}
}
