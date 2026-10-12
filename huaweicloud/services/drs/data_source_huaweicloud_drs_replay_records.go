package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/replay-record
func DataSourceReplayRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReplayRecordsRead,

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
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sql_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execute_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"finished_sql": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"abnormal_sql": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slow_sql": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_sql": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildReplayRecordsQueryParams(d *schema.ResourceData, offset int) string {
	return fmt.Sprintf("?limit=2000&start_time=%d&end_time=%d&offset=%d",
		d.Get("start_time").(int), d.Get("end_time").(int), offset)
}

func dataSourceReplayRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/replay-record"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))

	var (
		totalCount interface{}
		sqlRecords = make([]interface{}, 0)
	)
	offset := 0
	for {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		requestResp, err := client.Request("GET", listPath+buildReplayRecordsQueryParams(d, offset), &opt)
		if err != nil {
			return diag.Errorf("error retrieving DRS replay records: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if offset == 0 {
			totalCount = utils.PathSearch("total_count", listRespBody, 0)
		}
		items := utils.PathSearch("sql_records", listRespBody, make([]interface{}, 0)).([]interface{})
		sqlRecords = append(sqlRecords, items...)

		offset += 2000
		if len(items) < 2000 {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_count", totalCount),
		d.Set("sql_records", flattenReplayRecords(sqlRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplayRecords(sqlRecords []interface{}) []interface{} {
	if len(sqlRecords) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(sqlRecords))
	for _, item := range sqlRecords {
		result = append(result, map[string]interface{}{
			"execute_time": utils.PathSearch("execute_time", item, nil),
			"finished_sql": utils.PathSearch("finished_sql", item, nil),
			"abnormal_sql": utils.PathSearch("abnormal_sql", item, nil),
			"slow_sql":     utils.PathSearch("slow_sql", item, nil),
			"total_sql":    utils.PathSearch("total_sql", item, nil),
		})
	}
	return result
}
