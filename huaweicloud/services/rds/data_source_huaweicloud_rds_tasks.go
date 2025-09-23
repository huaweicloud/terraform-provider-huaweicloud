package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/tasklist/detail
func DataSourceRdsTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     tasksJobSchema(),
			},
		},
	}
}

func tasksJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ended": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobsInstanceSchema(),
			},
		},
	}
}

func jobsInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	url := "v3/{project_id}/instances/{instance_id}/tasklist/detail"
	getUrl := client.Endpoint + url
	getUrl = strings.ReplaceAll(getUrl, "{project_id}", client.ProjectID)
	getUrl = strings.ReplaceAll(getUrl, "{instance_id}", d.Get("instance_id").(string))
	getUrl += buildTasksListQueryParam(d)

	getOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getResp, err := client.Request("GET", getUrl, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS task list: %s", err)
	}

	body, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenTasksJobs(body)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTasksListQueryParam(d *schema.ResourceData) string {
	params := fmt.Sprintf("?start_time=%v", d.Get("start_time"))
	if v, ok := d.GetOk("end_time"); ok {
		params = fmt.Sprintf("%s&end_time=%v", params, v)
	}
	return params
}

func flattenTasksJobs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rawJobs := utils.PathSearch("jobs", resp, nil)
	jobs, ok := rawJobs.([]interface{})
	if !ok {
		return nil
	}

	out := make([]interface{}, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, map[string]interface{}{
			"id":          utils.PathSearch("id", j, nil),
			"name":        utils.PathSearch("name", j, nil),
			"status":      utils.PathSearch("status", j, nil),
			"created":     utils.PathSearch("created", j, nil),
			"ended":       utils.PathSearch("ended", j, nil),
			"process":     utils.PathSearch("process", j, nil),
			"task_detail": utils.PathSearch("task_detail", j, nil),
			"fail_reason": utils.PathSearch("fail_reason", j, nil),
			"entities":    utils.PathSearch("entities", j, nil),
			"instance":    flattenTasksInstance(j),
		})
	}
	return out
}

func flattenTasksInstance(job interface{}) []interface{} {
	raw := utils.PathSearch("instance", job, nil)
	if raw == nil {
		return []interface{}{}
	}

	inst, ok := raw.(map[string]interface{})
	if !ok || len(inst) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", inst, nil),
			"name": utils.PathSearch("name", inst, nil),
		},
	}
}
