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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/replay-execute-record
func DataSourceDrsReplayExecuteRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsReplayExecuteRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to query replay execute records.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the start time of the SQL execution record, in second-level timestamp format.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the end time of the SQL execution record, in second-level timestamp format.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of items in the list, regardless of paging.",
			},
			"sql_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of SQL execution records of the replay task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execute_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The execution time of the record, in second-level timestamp format.",
						},
						"finished_sql": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of finished SQL statements.",
						},
						"abnormal_sql": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of abnormal SQL statements.",
						},
						"slow_sql": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of slow SQL statements.",
						},
						"total_sql": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of SQL statements in the execution record.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsReplayExecuteRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/replay-execute-record"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
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
		requestResp, err := client.Request("GET", listPath+buildListReplayExecuteRecordsQueryParams(d, offset), &opt)
		if err != nil {
			return diag.Errorf("error retrieving DRS replay execute records: %s", err)
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

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("total_count", totalCount),
		d.Set("sql_records", flattenReplayExecuteRecords(sqlRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListReplayExecuteRecordsQueryParams(d *schema.ResourceData, offset int) string {
	res := ""
	res = fmt.Sprintf("%s&start_time=%v", res, d.Get("start_time"))
	res = fmt.Sprintf("%s&end_time=%v", res, d.Get("end_time"))
	res = fmt.Sprintf("%s&offset=%d", res, offset)

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenReplayExecuteRecords(sqlRecords []interface{}) []interface{} {
	if len(sqlRecords) == 0 {
		return nil
	}

	res := make([]interface{}, 0, len(sqlRecords))
	for _, v := range sqlRecords {
		res = append(res, map[string]interface{}{
			"execute_time": utils.PathSearch("execute_time", v, nil),
			"finished_sql": utils.PathSearch("finished_sql", v, nil),
			"abnormal_sql": utils.PathSearch("abnormal_sql", v, nil),
			"slow_sql":     utils.PathSearch("slow_sql", v, nil),
			"total_sql":    utils.PathSearch("total_sql", v, nil),
		})
	}
	return res
}
