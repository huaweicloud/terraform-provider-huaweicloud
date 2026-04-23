package rfs

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var applyExecutionPlanNonUpdatableParams = []string{
	"stack_name",
	"execution_plan_name",
	"execution_plan_id",
	"stack_id",
}

// @API RFS POST /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}
// @API RFS GET /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/metadata
func ResourceApplyExecutionPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplyExecutionPlanCreate,
		ReadContext:   resourceApplyExecutionPlanRead,
		UpdateContext: resourceApplyExecutionPlanUpdate,
		DeleteContext: resourceApplyExecutionPlanDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(applyExecutionPlanNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_plan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildApplyExecutionPlanBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"execution_plan_id": utils.ValueIgnoreEmpty(d.Get("execution_plan_id")),
		"stack_id":          utils.ValueIgnoreEmpty(d.Get("stack_id")),
	}
}

func readExecutionPlanDetail(client *golangsdk.ServiceClient, reqUUID, stackName, planName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", planName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForExecutionPlanApplied(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, reqUUID string) error {
	stackName := d.Get("stack_name").(string)
	executionPlanName := d.Get("execution_plan_name").(string)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := readExecutionPlanDetail(client, reqUUID, stackName, executionPlanName)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("unable to find status in API response")
			}

			if status == "APPLIED" {
				return respBody, "COMPLETED", nil
			}

			// The documentation does not provide the status values ​​for application failures;
			// To be on the safe side, all other statuses will be treated as PENDING.
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceApplyExecutionPlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		stackName         = d.Get("stack_name").(string)
		executionPlanName = d.Get("execution_plan_name").(string)
		httpUrl           = "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", executionPlanName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApplyExecutionPlanBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error triggering RFS stack execution plan: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentId := utils.PathSearch("deployment_id", respBody, "").(string)
	if deploymentId == "" {
		return diag.Errorf("error triggering RFS stack execution plan: Deployment ID is not found in API response")
	}

	d.SetId(deploymentId)

	if err := waitingForExecutionPlanApplied(ctx, client, d, d.Timeout(schema.TimeoutCreate), requestId); err != nil {
		return diag.Errorf("error waiting for RFS execution plan to be applied: %s", err)
	}

	return nil
}

func resourceApplyExecutionPlanRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceApplyExecutionPlanUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceApplyExecutionPlanDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to trigger stack execution plan. Deleting this resource
    will not cancel the execution plan operation, but will only remove resource information from
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
