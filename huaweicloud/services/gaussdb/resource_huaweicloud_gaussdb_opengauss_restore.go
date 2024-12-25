package gaussdb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/recovery
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceOpenGaussRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussRestoreCreate,
		ReadContext:   resourceOpenGaussRestoreRead,
		DeleteContext: resourceOpenGaussRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"restore_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"backup_id"},
			},
		},
	}
}

func resourceOpenGaussRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/recovery"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	targetInstanceId := d.Get("target_instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateRestoreBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("target_instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error restoring GaussDB OpenGauss instance (%s): %s", targetInstanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error restoring GaussDB OpenGauss instance(%s), job_id is not found in the response",
			targetInstanceId)
	}

	d.SetId(jobId.(string))

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 120, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for restoring GaussDB OpenGauss instance(%s) to complete: %s",
			targetInstanceId, err)
	}

	return nil
}

func buildCreateRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source": map[string]interface{}{
			"instance_id":  d.Get("source_instance_id"),
			"type":         d.Get("type"),
			"backup_id":    utils.ValueIgnoreEmpty(d.Get("backup_id")),
			"restore_time": utils.ValueIgnoreEmpty(d.Get("restore_time")),
		},
		"target": map[string]interface{}{
			"instance_id": d.Get("target_instance_id"),
		},
	}
	return bodyParams
}

func resourceOpenGaussRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration record is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
