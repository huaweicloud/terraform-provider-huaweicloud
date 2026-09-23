package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/export-jobs
func DataSourceExportJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExportJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbs_instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"completed_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inner_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_type": {
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
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"async_job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildExportJobsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"job_type":              d.Get("job_type"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"status":                utils.ValueIgnoreEmpty(d.Get("status")),
		"dbs_instance_ids":      utils.ValueIgnoreEmpty(d.Get("dbs_instance_ids")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"create_at":             utils.ValueIgnoreEmpty(d.Get("create_at")),
		"completed_at":          utils.ValueIgnoreEmpty(d.Get("completed_at")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"engine_type":           utils.ValueIgnoreEmpty(d.Get("engine_type")),
		"net_type":              utils.ValueIgnoreEmpty(d.Get("net_type")),
		"billing_tag":           utils.ValueIgnoreEmpty(d.Get("billing_tag")),
		"billing_mode":          utils.ValueIgnoreEmpty(d.Get("billing_mode")),
		"public_ip":             utils.ValueIgnoreEmpty(d.Get("public_ip")),
		"instance_ip":           utils.ValueIgnoreEmpty(d.Get("instance_ip")),
		"inner_ip":              utils.ValueIgnoreEmpty(d.Get("inner_ip")),
		"spec_type":             utils.ValueIgnoreEmpty(d.Get("spec_type")),
		"direction":             utils.ValueIgnoreEmpty(d.Get("direction")),
		"task_type":             utils.ValueIgnoreEmpty(d.Get("task_type")),
		"sort_key":              utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"sort_dir":              utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"tag":                   utils.ValueIgnoreEmpty(d.Get("tag")),
	}
}

func dataSourceExportJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/export-jobs"
	)

	client, err := cfg.DrsV5Client(region)
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
		JSONBody: utils.RemoveNil(buildExportJobsBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error exporting DRS jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("async_job_id", utils.PathSearch("id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
