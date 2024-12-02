package taurusdb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL DELETE /v3/{project_id}/scheduled-jobs
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
func ResourceGaussDBScheduledTaskCancel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlScheduledTaskCancelCreate,
		ReadContext:   resourceGaussDBMysqlScheduledTaskCancelRead,
		DeleteContext: resourceGaussDBMysqlScheduledTaskCancelDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_name": {
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
			"job_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBMysqlScheduledTaskCancelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/scheduled-jobs"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	jobId := d.Get("job_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMySQLScheduledTaskCancelBodyParams(jobId))

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error canceling GaussDB MySQL scheduled task(%s): %s", jobId, err)
	}

	d.SetId(jobId)

	return resourceGaussDBMysqlScheduledTaskCancelRead(ctx, d, meta)
}

func buildCreateGaussDBMySQLScheduledTaskCancelBodyParams(jobId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_ids": []string{jobId},
	}
	return bodyParams
}

func resourceGaussDBMysqlScheduledTaskCancelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/scheduled-jobs?job_id={job_id}"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB MySQL scheduled task")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	scheduledTask := utils.PathSearch("schedules[0]", getRespBody, nil)
	if scheduledTask == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", scheduledTask, nil)),
		d.Set("instance_name", utils.PathSearch("instance_name", scheduledTask, nil)),
		d.Set("instance_status", utils.PathSearch("instance_status", scheduledTask, nil)),
		d.Set("project_id", utils.PathSearch("project_id", scheduledTask, nil)),
		d.Set("job_name", utils.PathSearch("job_name", scheduledTask, nil)),
		d.Set("create_time", utils.PathSearch("create_time", scheduledTask, nil)),
		d.Set("start_time", utils.PathSearch("start_time", scheduledTask, nil)),
		d.Set("end_time", utils.PathSearch("end_time", scheduledTask, nil)),
		d.Set("job_status", utils.PathSearch("job_status", scheduledTask, nil)),
		d.Set("datastore_type", utils.PathSearch("datastore_type", scheduledTask, nil)),
		d.Set("job_status", utils.PathSearch("job_status", scheduledTask, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBMysqlScheduledTaskCancelDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB MySQL scheduled task cancel resource is not supported. The GaussDB MySQL scheduled " +
		"task cancel resource is only removed from the state, the GaussDB MySQL scheduled task remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
