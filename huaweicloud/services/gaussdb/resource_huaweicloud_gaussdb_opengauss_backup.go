package gaussdb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/backups
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3.1/{project_id}/backups
// @API GaussDB DELETE /v3/{project_id}/backups/{backup_id}
func ResourceGaussDBOpenGaussBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussBackupCreate,
		ReadContext:   resourceOpenGaussBackupRead,
		DeleteContext: resourceOpenGaussBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpenGaussBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussBackupBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss backup: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	backupId := utils.PathSearch("backup.id", createRespBody, nil)
	if backupId == nil {
		return diag.Errorf("error creating GaussDB OpenGauss backup: backup ID is not found in API response")
	}
	d.SetId(backupId.(string))

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating GaussDB OpenGauss backup: job_id is not found in API response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 30, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenGaussBackupRead(ctx, d, meta)
}

func buildCreateOpenGaussBackupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_id": d.Get("instance_id"),
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceOpenGaussBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3.1/{project_id}/backups?backup_id={backup_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{backup_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB OpenGauss backup")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	backup := utils.PathSearch("backups|[0]", getRespBody, nil)
	if backup == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", backup, nil)),
		d.Set("name", utils.PathSearch("name", backup, nil)),
		d.Set("description", utils.PathSearch("description", backup, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", backup, nil)),
		d.Set("end_time", utils.PathSearch("end_time", backup, nil)),
		d.Set("size", utils.PathSearch("size", backup, nil)),
		d.Set("type", utils.PathSearch("type", backup, nil)),
		d.Set("datastore", flattenGetGaussDBBackupResponseBodyDatastore(backup)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussDBBackupResponseBodyDatastore(backup interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("datastore.type", backup, nil),
			"version": utils.PathSearch("datastore.version", backup, nil),
		},
	}
	return rst
}

func resourceOpenGaussBackupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/{backup_id}"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{backup_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200751"),
			"error deleting GaussDB OpenGauss backup")
	}

	return nil
}
