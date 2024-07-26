package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/jobs
func DataSourceFactoryJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFactoryJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the jobs belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the job name to be queried.`,
			},
			"process_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specified the job type to be queried.`,
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All jobs that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the job.`,
						},
						"process_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the job.`,
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The priority of the job.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the job.`,
						},
						"is_single_task_job": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the job is single task.`,
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The directory tree path of the job.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the job.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the job scheduling, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the job scheduling, in RFC3339 format.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the job.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the job, in RFC3339 format.`,
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user who last updated the job.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the job, in RFC3339 format.`,
						},
						"last_instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest running status of the instance corresponding to the job.`,
						},
						"last_instance_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest end time of the instance corresponding to the job, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildFactoryJobQueryParams(d *schema.ResourceData) string {
	res := ""
	if jobName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&jobName=%v", res, jobName)
	}

	if jobType, ok := d.GetOk("process_type"); ok {
		res = fmt.Sprintf("%s&jobType=%v", res, jobType)
	}

	if res != "" {
		res = "&" + res[1:]
	}
	return res
}

func queryFactoryJobs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/jobs?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	queryParams := buildFactoryJobQueryParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving DataArts Factory jobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		jobs := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) < 1 {
			break
		}
		result = append(result, jobs...)
		offset += len(jobs)
	}
	return result, nil
}

func dataSourceFactoryJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	jobs, err := queryFactoryJobs(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("jobs", flattenJobs(jobs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobs(all []interface{}) []interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(all))
	for _, job := range all {
		result = append(result, map[string]interface{}{
			"name":                   utils.PathSearch("name", job, nil),
			"process_type":           utils.PathSearch("jobType", job, nil),
			"priority":               utils.PathSearch("priority", job, nil),
			"owner":                  utils.PathSearch("owner", job, nil),
			"is_single_task_job":     utils.PathSearch("singleNodeJobFlag", job, false),
			"directory":              utils.PathSearch("path", job, nil),
			"status":                 utils.PathSearch("status", job, nil),
			"start_time":             flattenTimeToRFC3339(utils.PathSearch("startTime", job, float64(0))),
			"end_time":               flattenTimeToRFC3339(utils.PathSearch("endTime", job, float64(0))),
			"created_by":             utils.PathSearch("createUser", job, nil),
			"created_at":             flattenTimeToRFC3339(utils.PathSearch("createTime", job, float64(0))),
			"updated_at":             flattenTimeToRFC3339(utils.PathSearch("lastUpdateTime", job, float64(0))),
			"updated_by":             utils.PathSearch("lastUpdateUser", job, nil),
			"last_instance_status":   utils.PathSearch("lastInstanceStatus", job, nil),
			"last_instance_end_time": flattenTimeToRFC3339(utils.PathSearch("lastInstanceEndTime", job, float64(0))),
		})
	}
	return result
}

// Formats the time according to the local computer's time.
func flattenTimeToRFC3339(timeStamp interface{}) string {
	return utils.FormatTimeStampRFC3339(int64(timeStamp.(float64))/1000, false)
}
