package rds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var primaryStandbySwitchNonUpdatableParams = []string{"instance_id", "force"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/failover
// @API RDS GET /v3/{project_id}/jobs
func ResourceRdsPrimaryStandbySwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPrimaryStandbySwitchCreate,
		ReadContext:   resourceRdsPrimaryStandbySwitchRead,
		UpdateContext: resourceRdsPrimaryStandbySwitchUpdate,
		DeleteContext: resourceRdsPrimaryStandbySwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(primaryStandbySwitchNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of instance.`,
			},
			"force": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to perform a forcible primary/standby switchover.`,
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

func resourceRdsPrimaryStandbySwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/failover"
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

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreatePrimaryStandbySwitchBodyParams(d))

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
		return diag.Errorf("error switching primary/standby RDS instance (%s): %s", instanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("workflowId", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("unable to find the workflowId from the response: %s", err)
	}

	d.SetId(instanceId)

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for switching primary/standby RDS instance(%s) to complete: %s",
			instanceId, err)
	}

	return nil
}

func buildCreatePrimaryStandbySwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"force": d.Get("force"),
	}
	return bodyParams
}

func resourceRdsPrimaryStandbySwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsPrimaryStandbySwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsPrimaryStandbySwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting primary standby switch is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before switch."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
