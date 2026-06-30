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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/health-compare-jobs
func DataSourceDrsHealthCompareJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsHealthCompareJobsRead,

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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compare_jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     compareJobsSchema(),
			},
		},
	}
}

func compareJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
			"compute_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     compareJobsdatabaseInfoSchema(),
			},
		},
	}
}

func compareJobsdatabaseInfoSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func buildHealthCompareJobsQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDrsHealthCompareJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/health-compare-jobs"
		jobId   = d.Get("job_id").(string)
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := requestPath + buildHealthCompareJobsQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS health compare jobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobs := utils.PathSearch("compare_jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) == 0 {
			break
		}

		result = append(result, jobs...)
		offset += len(jobs)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("compare_jobs", flattenHealthCompareJobs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHealthCompareJobs(compareJobs []interface{}) []interface{} {
	if len(compareJobs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(compareJobs))
	for _, v := range compareJobs {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"start_time":   utils.PathSearch("start_time", v, nil),
			"end_time":     utils.PathSearch("end_time", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"compute_type": utils.PathSearch("compute_type", v, nil),
			"database_info": flattenCompareJobsDatabaseInfo(
				utils.PathSearch("database_info", v, nil)),
		})
	}
	return result
}

func flattenCompareJobsDatabaseInfo(databaseInfo interface{}) []interface{} {
	if databaseInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"service_database":           utils.PathSearch("service_database", databaseInfo, nil),
			"disaster_recovery_database": utils.PathSearch("disaster_recovery_database", databaseInfo, nil),
		},
	}
}
