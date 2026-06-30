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

// @API DRS GET /v5/{project_id}/batch-async-jobs
func DataSourceBatchAsyncJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBatchAsyncJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"async_job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"async_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBatchAsyncJobsQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("async_job_id"); ok {
		queryParams += fmt.Sprintf("&async_job_id=%s", v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams += fmt.Sprintf("&status=%s", v.(string))
	}
	if v, ok := d.GetOk("domain_name"); ok {
		queryParams += fmt.Sprintf("&domain_name=%s", v.(string))
	}
	if v, ok := d.GetOk("user_name"); ok {
		queryParams += fmt.Sprintf("&user_name=%s", v.(string))
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams += fmt.Sprintf("&sort_key=%s", v.(string))
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams += fmt.Sprintf("&sort_dir=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceBatchAsyncJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/batch-async-jobs"
		offset  = 0
		result  = make([]interface{}, 0)
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
	}

	for {
		requestPathWithQuery := requestPath + buildBatchAsyncJobsQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS batch async jobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobsResp := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobsResp) == 0 {
			break
		}

		result = append(result, jobsResp...)
		offset += len(jobsResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenBatchAsyncJobs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchAsyncJobs(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"async_job_id": utils.PathSearch("async_job_id", item, nil),
			"status":       utils.PathSearch("status", item, nil),
			"domain_name":  utils.PathSearch("domain_name", item, nil),
			"user_name":    utils.PathSearch("user_name", item, nil),
			"create_time":  utils.PathSearch("create_time", item, nil),
		})
	}

	return result
}
