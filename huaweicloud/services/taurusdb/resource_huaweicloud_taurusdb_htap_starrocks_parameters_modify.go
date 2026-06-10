package taurusdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var starrocksParametersModifyNoneUpdatableParams = []string{
	"instance_id", "starrocks_instance_id", "node_type", "parameter_values",
}

// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/configurations
// @API TaurusDB GET /v3/{project_id}/jobs
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/{starrocks_instance_id}
// @API TaurusDB PUT /v3/{project_id}/instances/{starrocks_instance_id}/starrocks/restart
func ResourceTaurusDBHtapStarrocksParametersModify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksParametersModifyCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksParametersModifyRead,
		UpdateContext: resourceTaurusDBHtapStarrocksParametersModifyUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksParametersModifyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(starrocksParametersModifyNoneUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"taurusdb_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"starrocks_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameter_values": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceTaurusDBHtapStarrocksParametersModifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		taurusdbInstanceId  = d.Get("taurusdb_instance_id").(string)
		starrocksInstanceId = d.Get("starrocks_instance_id").(string)
		timeout             = d.Timeout(schema.TimeoutCreate)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	// Modify StarRocks parameters
	restartRequired, err := modifyStarrocksParameters(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// If restart is required, restart the StarRocks instance
	if restartRequired {
		err := restartStarrocksInstance(ctx, client, taurusdbInstanceId, starrocksInstanceId, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func modifyStarrocksParameters(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) (bool, error) {
	var (
		instanceId          = d.Get("taurusdb_instance_id").(string)
		starrocksInstanceId = d.Get("starrocks_instance_id").(string)
		timeout             = d.Timeout(schema.TimeoutCreate)
	)
	modifyPath := client.Endpoint + "v3/{project_id}/instances/{instance_id}/starrocks/configurations"
	modifyPath = strings.ReplaceAll(modifyPath, "{project_id}", client.ProjectID)
	modifyPath = strings.ReplaceAll(modifyPath, "{instance_id}", starrocksInstanceId)

	modifyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	modifyOpt.JSONBody = buildModifyStarrocksParametersBodyParams(d)

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", modifyPath, &modifyOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, instanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return false, fmt.Errorf("error modifying TaurusDB Htap StarRocks instance(%s) parameters: %s",
			starrocksInstanceId, err)
	}
	modifyRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return false, err
	}

	jobId := utils.PathSearch("job_id", modifyRespBody, "").(string)
	if jobId == "" {
		return false, fmt.Errorf("error modifying TaurusDB Htap StarRocks instance(%s) parameters, job_id is not found in the response",
			starrocksInstanceId)
	}

	// Wait for the modify job to complete
	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return false, fmt.Errorf("error waiting for modifying TaurusDB Htap StarRocks parameters job (%s) to complete: %s", jobId, err)
	}

	restartRequired := utils.PathSearch("restart_required", modifyRespBody, false).(bool)

	return restartRequired, nil
}

func buildModifyStarrocksParametersBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_type":        d.Get("node_type"),
		"parameter_values": utils.ExpandToStringMap(d.Get("parameter_values").(map[string]interface{})),
	}
	return bodyParams
}

func resourceTaurusDBHtapStarrocksParametersModifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksParametersModifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksParametersModifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameters modify resource is not supported. The resource is only removed from the state," +
		" the StarRocks instance parameters remain in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
