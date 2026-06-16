package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/scheduled-jobs
func DataSourceGeminiDBScheduledTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBScheduledTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
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
					},
				},
			},
		},
	}
}

func dataSourceGeminiDBScheduledTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/scheduled-jobs"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	// Build query parameters
	queryParams := buildScheduledTasksQueryParams(d)
	getPath += queryParams

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving scheduled jobs: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	schedules := utils.PathSearch("schedules", respBody, []interface{}{}).([]interface{})

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("schedules", flattenScheduledTasks(schedules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildScheduledTasksQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("job_name"); ok {
		queryParams = fmt.Sprintf("%s&job_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("job_status"); ok {
		queryParams = fmt.Sprintf("%s&job_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenScheduledTasks(schedules []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(schedules))

	for _, schedule := range schedules {
		result = append(result, map[string]interface{}{
			"job_id":          utils.PathSearch("job_id", schedule, ""),
			"job_name":        utils.PathSearch("job_name", schedule, ""),
			"job_status":      utils.PathSearch("job_status", schedule, ""),
			"instance_id":     utils.PathSearch("instance_id", schedule, ""),
			"instance_name":   utils.PathSearch("instance_name", schedule, ""),
			"instance_status": utils.PathSearch("instance_status", schedule, ""),
			"datastore_type":  utils.PathSearch("datastore_type", schedule, ""),
			"create_time":     utils.PathSearch("create_time", schedule, ""),
			"start_time":      utils.PathSearch("start_time", schedule, ""),
			"end_time":        utils.PathSearch("end_time", schedule, ""),
		})
	}

	return result
}
