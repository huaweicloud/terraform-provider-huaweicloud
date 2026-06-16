package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v3/{project_id}/jobs/{job_id}/table/compare
func DataSourceDrsTableCompare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTableCompareRead,

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
			"compare_jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
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
						"compare_job_tag": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"proportion_value": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func buildTableCompareQueryParams(offset int) string {
	if offset == 0 {
		return ""
	}

	return fmt.Sprintf("?offset=%d", offset)
}

func dataSourceTableCompareRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v3/{project_id}/jobs/{job_id}/table/compare"
		offset  = 0
		result  = make([]interface{}, 0)
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

	for {
		requestPathWithOffset := requestPath + buildTableCompareQueryParams(offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS table compare: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		compareJobs := utils.PathSearch("compare_jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(compareJobs) == 0 {
			break
		}

		result = append(result, compareJobs...)
		offset += len(compareJobs)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("compare_jobs", flattenCompareJobs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCompareJobs(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		optionsMap := utils.PathSearch("options", item, make(map[string]interface{})).(map[string]interface{})
		compareJobTagMap := utils.PathSearch("compare_job_tag", item, make(map[string]interface{})).(map[string]interface{})

		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", item, nil),
			"type":                  utils.PathSearch("type", item, nil),
			"options":               utils.ExpandToStringMap(optionsMap),
			"start_time":            utils.PathSearch("start_time", item, nil),
			"end_time":              utils.PathSearch("end_time", item, nil),
			"status":                utils.PathSearch("status", item, nil),
			"export_status":         utils.PathSearch("export_status", item, nil),
			"report_remain_seconds": utils.PathSearch("report_remain_seconds", item, nil),
			"compare_job_tag":       utils.ExpandToStringMap(compareJobTagMap),
			"proportion_value":      utils.PathSearch("proportion_value", item, nil),
			"database_info":         flattenTableDatabaseInfo(utils.PathSearch("database_info", item, nil)),
		})
	}
	return result
}

func flattenTableDatabaseInfo(respMap interface{}) []interface{} {
	if respMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"service_database":           utils.PathSearch("service_database", respMap, nil),
			"disaster_recovery_database": utils.PathSearch("disaster_recovery_database", respMap, nil),
		},
	}
}
