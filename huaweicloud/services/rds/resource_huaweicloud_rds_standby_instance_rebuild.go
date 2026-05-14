package rds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var rebuildSlaveNonUpdatableParams = []string{
	"instance_id",
}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/rebuild
// @API RDS GET /v3/{project_id}/instances/{instance_id}/rebuild
// @API RDS GET /v3/{project_id}/instances
func ResourceStandbyInstanceRebuild() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStandbyInstanceRebuildCreate,
		ReadContext:   resourceStandbyInstanceRebuildRead,
		UpdateContext: resourceStandbyInstanceRebuildUpdate,
		DeleteContext: resourceStandbyInstanceRebuildDelete,

		CustomizeDiff: config.FlexibleForceNew(rebuildSlaveNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"workflow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_rebuild_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_rebuild_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStandbyInstanceRebuildCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/rebuild"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS standby instance rebuild:(%s): %s", instanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	workflowID := utils.PathSearch("workflow_id", createRespBody, "").(string)
	if workflowID == "" {
		return diag.Errorf("error creating RDS standby instance rebuild: workflow_id not found in the response")
	}
	err = checkRDSInstanceJobFinish(client, workflowID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)

	return resourceStandbyInstanceRebuildRead(ctx, d, meta)
}

func resourceStandbyInstanceRebuildRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/rebuild"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Id()
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS standby instance rebuild")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	workflowId := utils.PathSearch("workflow_id", getRespBody, "").(string)
	if workflowId == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving RDS standby instance rebuild")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", getRespBody, nil)),
		d.Set("workflow_id", utils.PathSearch("workflow_id", getRespBody, nil)),
		d.Set("last_rebuild_time", utils.PathSearch("last_rebuild_time", getRespBody, nil)),
		d.Set("next_rebuild_time", utils.PathSearch("next_rebuild_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceStandbyInstanceRebuildUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceStandbyInstanceRebuildDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS standby instance rebuild resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
