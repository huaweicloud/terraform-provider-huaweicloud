package dli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI GET /v2.0/{project_id}/batches
func DataSourceSparkJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSparkJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the spark jobs are located.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The DLI queue name of the spark job to be queried.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The queue name of the spark job to be queried.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the spark job to be queried.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the spark job to be queried.`,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The state of the spark job to be queried.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The owner of the spark job to be queried.`,
			},

			// Attributes
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of spark jobs that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the spark job.`,
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The backend app ID of the spark job.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the spark job.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the spark job.`,
						},
						"queue": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The queue of the spark job.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster name of the spark job.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The state of the spark job.`,
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the spark job.`,
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The running duration of the spark job, in milliseconds.`,
						},
						"sc_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The compute resource type of the spark job.`,
						},
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The custom image of the spark job.`,
						},
						"log": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The last 10 log records of the spark job.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"req_body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request body details of the spark job.`,
						},
						"proxy_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The proxy user of the spark job.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the spark job, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the spark job, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildSparkJobsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("cluster_name"); ok {
		res += fmt.Sprintf("&cluster_name=%v", v)
	}
	if v, ok := d.GetOk("queue_name"); ok {
		res += fmt.Sprintf("&queue_name=%v", v)
	}
	if v, ok := d.GetOk("job_name"); ok {
		res += fmt.Sprintf("&job_name=%v", v)
	}
	if v, ok := d.GetOk("job_id"); ok {
		res += fmt.Sprintf("&job-id=%v", v)
	}
	if v, ok := d.GetOk("state"); ok {
		res += fmt.Sprintf("&state=%v", v)
	}
	if v, ok := d.GetOk("owner"); ok {
		res += fmt.Sprintf("&owner=%v", v)
	}

	return res
}

func querySparkJobs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl  = "v2.0/{project_id}/batches?size={size}"
		size     = 100
		from     = 0
		result   = make([]interface{}, 0)
		respBody interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{size}", strconv.Itoa(size))
	listPath += buildSparkJobsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithFrom := fmt.Sprintf("%s&from=%d", listPath, from)
		requestResp, err := client.Request("GET", listPathWithFrom, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err = utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		sessions := utils.PathSearch("sessions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, sessions...)
		if len(sessions) < size {
			break
		}

		from += size
	}

	return result, nil
}

func flattenSparkJobs(jobs []interface{}) []map[string]interface{} {
	if len(jobs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", job, nil),
			"app_id":       utils.PathSearch("appId", job, nil),
			"name":         utils.PathSearch("name", job, nil),
			"owner":        utils.PathSearch("owner", job, nil),
			"queue":        utils.PathSearch("queue", job, nil),
			"cluster_name": utils.PathSearch("cluster_name", job, nil),
			"state":        utils.PathSearch("state", job, nil),
			"kind":         utils.PathSearch("kind", job, nil),
			"duration":     int(utils.PathSearch("duration", job, float64(0)).(float64)),
			"sc_type":      utils.PathSearch("sc_type", job, nil),
			"image":        utils.PathSearch("image", job, nil),
			"log":          utils.PathSearch("log", job, make([]interface{}, 0)),
			"req_body":     utils.PathSearch("req_body", job, nil),
			"proxy_user":   utils.PathSearch("proxyUser", job, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", job, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("update_time", job, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func dataSourceSparkJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	jobList, err := querySparkJobs(client, d)
	if err != nil {
		return diag.Errorf("error querying spark jobs: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("jobs", flattenSparkJobs(jobList)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
