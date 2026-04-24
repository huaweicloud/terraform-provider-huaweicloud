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

var rfsContinueDeployStackNonUpdatableParams = []string{
	"stack_name", "stack_id"}

// @API RFS POST /v1/{project_id}/stacks/{stack_name}/continuations
// @API RFS GET /v1/{project_id}/stacks/{stack_name}/metadata
func ResourceContinueDeployStack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRfsContinueDeployStackCreate,
		ReadContext:   resourceRfsContinueDeployStackRead,
		UpdateContext: resourceRfsContinueDeployStackUpdate,
		DeleteContext: resourceRfsContinueDeployStackDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(rfsContinueDeployStackNonUpdatableParams),

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

func buildStackContinuationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"stack_id": utils.ValueIgnoreEmpty(d.Get("stack_id")),
	}
	return bodyParams
}

func readStackMetadata(client *golangsdk.ServiceClient, reqUUID, stackName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/stacks/{stack_name}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": reqUUID,
			"Content-Type":      "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForStackDeploymentCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, reqUUID string) error {
	stackName := d.Get("stack_name").(string)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := readStackMetadata(client, reqUUID, stackName)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("unable to find status in API response")
			}

			if status == "DEPLOYMENT_COMPLETE" {
				return respBody, "COMPLETED", nil
			}

			if status == "DEPLOYMENT_FAILED" {
				return respBody, "ERROR", errors.New("stack deployment failed")
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

func resourceRfsContinueDeployStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/stacks/{stack_name}/continuations"
		stackName = d.Get("stack_name").(string)
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{stack_name}", stackName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
		},
		JSONBody: utils.RemoveNil(buildStackContinuationBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error continuing stack deployment (%s): %s", stackName, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentId := utils.PathSearch("deployment_id", respBody, "").(string)
	if deploymentId == "" {
		return diag.Errorf("unable to find the deployment_id from the API response")
	}

	d.SetId(deploymentId)

	if err := waitingForStackDeploymentCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate), requestId); err != nil {
		return diag.Errorf("error waiting for stack deployment to complete: %s", err)
	}

	return nil
}

func resourceRfsContinueDeployStackRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRfsContinueDeployStackUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRfsContinueDeployStackDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to continue deploying the resource stack. 
Deleting this resource will not undo the operation of continuing to deploy the resource stack, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
