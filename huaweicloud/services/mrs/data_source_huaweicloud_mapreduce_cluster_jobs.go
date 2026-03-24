package mrs

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

// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/job-executions
func DataSourceClusterJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterJobsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster jobs are located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the job.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the job.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name of the job submitter.`,
			},
			"job_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the job.`,
			},
			"job_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The execution status of the job.`,
			},
			"job_result": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The execution result of the job.`,
			},
			"queue": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource queue name of the job.`,
			},
			"submitted_time_begin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The begin time of the submitted jobs, in RFC3339 format.`,
			},
			"submitted_time_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the submitted jobs, in RFC3339 format.`,
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of jobs that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the job.`,
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user name of the job submitter.`,
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the job.`,
						},
						"job_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution result of the job.`,
						},
						"job_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution status of the job.`,
						},
						"job_progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The execution progress of the job.`,
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the job.`,
						},
						"started_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the job, in RFC3339 format.`,
						},
						"submitted_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The submit time of the job, in RFC3339 format.`,
						},
						"finished_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The finish time of the job, in RFC3339 format.`,
						},
						"elapsed_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The elapsed time of the job, in milliseconds.`,
						},
						"arguments": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The runtime arguments of the job.`,
						},
						"launcher_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The launcher ID of the job.`,
						},
						"properties": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The properties of the job.`,
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application ID of the job.`,
						},
						"tracking_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tracking URL of the job logs.`,
						},
						"queue": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource queue name of the job.`,
						},
					},
				},
			},
		},
	}
}

func buildClusterJobsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("job_id"); ok {
		res = fmt.Sprintf("%s&job_id=%v", res, v)
	}

	if v, ok := d.GetOk("job_name"); ok {
		res = fmt.Sprintf("%s&job_name=%v", res, v)
	}

	if v, ok := d.GetOk("user"); ok {
		res = fmt.Sprintf("%s&user=%v", res, v)
	}

	if v, ok := d.GetOk("job_type"); ok {
		res = fmt.Sprintf("%s&job_type=%v", res, v)
	}

	if v, ok := d.GetOk("job_state"); ok {
		res = fmt.Sprintf("%s&job_state=%v", res, v)
	}

	if v, ok := d.GetOk("job_result"); ok {
		res = fmt.Sprintf("%s&job_result=%v", res, v)
	}

	if v, ok := d.GetOk("queue"); ok {
		res = fmt.Sprintf("%s&queue=%v", res, v)
	}

	if v, ok := d.GetOk("submitted_time_begin"); ok {
		res = fmt.Sprintf("%s&submitted_time_begin=%v", res, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}

	if v, ok := d.GetOk("submitted_time_end"); ok {
		res = fmt.Sprintf("%s&submitted_time_end=%v", res, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}

	return res
}

func listClusterJobs(client *golangsdk.ServiceClient, d *schema.ResourceData, clusterId string) ([]interface{}, error) {
	var (
		httpURL = "v2/{project_id}/clusters/{cluster_id}/job-executions"
		result  = make([]interface{}, 0)
		limit   = 100
		// The offset indicates the page number, starts from 1.
		offset = 1
	)

	listPath := client.Endpoint + httpURL
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	listPath = fmt.Sprintf("%s?limit=%d%s", listPath, limit, buildClusterJobsQueryParams(d))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{202},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		jobs := utils.PathSearch("job_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, jobs...)
		if len(jobs) < limit {
			break
		}

		offset++
	}

	return result, nil
}

func dataSourceClusterJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	jobs, err := listClusterJobs(client, d, clusterId)
	if err != nil {
		return diag.Errorf("error retrieving jobs under the specified cluster (%s): %s", clusterId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenClusterJobs(jobs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterJobs(jobs []interface{}) []map[string]interface{} {
	if len(jobs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobs))
	for _, v := range jobs {
		result = append(result, map[string]interface{}{
			"job_id":       utils.PathSearch("job_id", v, nil),
			"user":         utils.PathSearch("user", v, nil),
			"job_name":     utils.PathSearch("job_name", v, nil),
			"job_result":   utils.PathSearch("job_result", v, nil),
			"job_state":    utils.PathSearch("job_state", v, nil),
			"job_progress": utils.PathSearch("job_progress", v, nil),
			"job_type":     utils.PathSearch("job_type", v, nil),
			"started_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("started_time",
				v, float64(0)).(float64))/1000, false),
			"submitted_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("submitted_time",
				v, float64(0)).(float64))/1000, false),
			"finished_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("finished_time",
				v, float64(0)).(float64))/1000, false),
			"elapsed_time": utils.PathSearch("elapsed_time", v, nil),
			"arguments":    utils.PathSearch("arguments", v, nil),
			"launcher_id":  utils.PathSearch("launcher_id", v, nil),
			"properties":   utils.JsonToString(utils.PathSearch("properties", v, nil)),
			"app_id":       utils.PathSearch("app_id", v, nil),
			"tracking_url": utils.PathSearch("tracking_url", v, nil),
			"queue":        utils.PathSearch("queue", v, nil),
		})
	}

	return result
}
