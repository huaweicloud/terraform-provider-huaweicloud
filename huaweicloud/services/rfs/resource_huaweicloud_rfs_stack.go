package rfs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/{project_id}/stacks
// @API RFS GET /v1/{project_id}/stacks
// @API RFS POST /v1/{project_id}/stacks/{stack_name}/deployments
// @API RFS GET /v1/{project_id}/stacks/{stack_name}/events
// @API RFS POST /v1/{project_id}/stacks/{stack_name}/deletion
func ResourceStack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStackCreate,
		ReadContext:   resourceStackRead,
		UpdateContext: resourceStackUpdate,
		DeleteContext: resourceStackDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The region where the RFS resource stack is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the resource stack.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the resource stack.",
			},
			"agency": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: utils.SchemaDesc(
								"The name of the provider corresponding to the IAM agency.",
								utils.SchemaDescInput{
									Required: true,
								},
							),
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Description: utils.SchemaDesc(
								"The name of IAM agency authorized to IAC account for resources modification.",
								utils.SchemaDescInput{
									Required: true,
								},
							),
						},
					},
				},
				Description: "The configuration of the agencies authorized to IAC.",
			},
			"template_body": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"template_uri"},
				Description:   "The HCL/JSON template content for deployment resources.",
			},
			"vars_body": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vars_uri"},
				RequiredWith:  []string{"template_body"},
				Description:   "The variable content for deployment resources.",
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The OBS address where the HCL/JSON template archive (**.zip** file, which contains all " +
					"resource **.tf.json** script files to be deployed) or **.tf.json** file is located, which " +
					"describes the target status of the deployment resources.",
			},
			"vars_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"template_uri"},
				Description: "The OBS address where the variable (**.tfvars**) file corresponding to the HCL/JSON " +
					"template located, which describes the target status of the deployment resources.",
			},
			"enable_auto_rollback": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable automatic rollback.",
			},
			"enable_deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable delete protection.",
			},
			"retain_all_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to reserve resources when deleting the resource stack.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the resource stack.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func buildStackAgencies(agencies []interface{}) []interface{} {
	if len(agencies) < 1 || agencies[0] == nil {
		return nil
	}

	result := make([]interface{}, 0, len(agencies))
	for _, agency := range agencies {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"provider_name": utils.ValueIgnoreEmpty(utils.PathSearch("provider_name", agency, nil)),
			"agency_name":   utils.ValueIgnoreEmpty(utils.PathSearch("name", agency, nil)),
		}))
	}

	return result
}

func resourceStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/stacks"
		requestId, _ = uuid.GenerateUUID()
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createtOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateStackBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
			"Client-Request-Id": requestId,
		},
	}
	requestResp, err := client.Request("POST", createPath, &createtOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	stackId := utils.PathSearch("stack_id", respBody, "").(string)
	if stackId == "" {
		return diag.Errorf("unable to find the stack ID from the API response")
	}
	d.SetId(stackId)

	// When creates a stack using the template_body parameter, it is automatically deployed without calling the deployment API.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      stackStatusRefreshFunc(client, stackId, []string{"CREATION_COMPLETE", "DEPLOYMENT_COMPLETE"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStackRead(ctx, d, meta)
}

func buildCreateStackBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"stack_name":                 d.Get("name"),
		"agencies":                   utils.ValueIgnoreEmpty(buildStackAgencies(d.Get("agency").([]interface{}))),
		"description":                utils.ValueIgnoreEmpty(d.Get("description")),
		"template_body":              utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":               utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_body":                  utils.ValueIgnoreEmpty(d.Get("vars_body")),
		"vars_uri":                   utils.ValueIgnoreEmpty(d.Get("vars_uri")),
		"enable_auto_rollback":       d.Get("enable_auto_rollback"),
		"enable_deletion_protection": d.Get("enable_deletion_protection"),
	}
}

func deployStack(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, stackId string,
	timeout time.Duration) error {
	var (
		// Stack name is a required parameter.
		stackName    = d.Get("name").(string)
		httpUrl      = "v1/{project_id}/stacks/{stack_name}/deployments"
		requestId, _ = uuid.GenerateUUID()
	)
	deployPath := client.Endpoint + httpUrl
	deployPath = strings.ReplaceAll(deployPath, "{project_id}", client.ProjectID)
	deployPath = strings.ReplaceAll(deployPath, "{stack_name}", stackName)

	deployOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"template_body": utils.ValueIgnoreEmpty(d.Get("template_body")),
			"template_uri":  utils.ValueIgnoreEmpty(d.Get("template_uri")),
			"vars_body":     utils.ValueIgnoreEmpty(d.Get("vars_body")),
			"vars_uri":      utils.ValueIgnoreEmpty(d.Get("vars_uri")),
			"stack_id":      stackId,
		},
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
			"Client-Request-Id": requestId,
		},
	}
	requestResp, err := client.Request("POST", deployPath, &deployOpt)
	if err != nil {
		return err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      stackStatusRefreshFunc(client, stackId, []string{"DEPLOYMENT_COMPLETE"}),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if status := utils.PathSearch("status", resp, "").(string); status == "DEPLOYMENT_FAILED" {
			deploymentId := utils.PathSearch("deployment_id", respBody, "").(string)
			if deploymentId == "" {
				return fmt.Errorf("unable to find the deployment ID from the API response after request send")
			}
			return queryAllFailedEvents(client, stackId, stackName, deploymentId)
		}
	}
	return err
}

// QueryStackById is a method to query stack details using its ID.
func QueryStackById(client *golangsdk.ServiceClient, stackId string) (interface{}, error) {
	var (
		httpUrl      = "v1/{project_id}/stacks"
		requestId, _ = uuid.GenerateUUID()
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
			"Client-Request-Id": requestId,
		},
	}
	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	stackDetail := utils.PathSearch(fmt.Sprintf("stacks[?stack_id=='%s']|[0]", stackId), respBody, nil)
	if stackDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return stackDetail, nil
}

func stackStatusRefreshFunc(client *golangsdk.ServiceClient, stackId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QueryStackById(client, stackId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "Resource Not Found", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		failedErrors := []string{
			"DEPLOYMENT_FAILED",
			"ROLLBACK_FAILED",
			"DELETION_FAILED",
		}
		status := utils.PathSearch("status", resp, "").(string)
		if utils.StrSliceContains(failedErrors, status) {
			return resp, "ERROR", fmt.Errorf("unexpected status '%s'", status)
		}
		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func queryAllFailedEvents(client *golangsdk.ServiceClient, stackId, stackName, deploymentId string) error {
	var (
		mErr         *multierror.Error
		httpUrl      = "v1/{project_id}/stacks/{stack_name}/events"
		requestId, _ = uuid.GenerateUUID()
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{stack_name}", stackName)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"stack_id":      stackId,
			"deployment_id": deploymentId,
		},
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
			"Client-Request-Id": requestId,
		},
	}
	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	events := utils.PathSearch("stack_events", respBody, make([]interface{}, 0)).([]interface{})
	for _, event := range events {
		if utils.PathSearch("event_type", event, "").(string) == "ERROR" {
			mErr = multierror.Append(mErr, errors.New(utils.PathSearch("event_message", event, "").(string)))
		}
	}
	return mErr.ErrorOrNil()
}

func resourceStackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		stackId = d.Id()
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	resp, err := QueryStackById(client, stackId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "RFS resource stack")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("stack_name", resp, nil)),
		d.Set("description", utils.PathSearch("description", resp, nil)),
		d.Set("status", utils.PathSearch("status", resp, nil)),
		d.Set("created_at", utils.PathSearch("create_time", resp, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", resp, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving stack (%s) fields: %s", stackId, mErr)
	}
	return nil
}

func resourceStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	if d.HasChanges("template_body", "vars_body", "template_uri", "vars_uri") {
		if err = deployStack(ctx, client, d, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceStackRead(ctx, d, meta)
}

func resourceStackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/stacks/{stack_name}/deletion"
		stackName    = d.Get("name").(string)
		stackId      = d.Id()
		requestId, _ = uuid.GenerateUUID()
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{stack_name}", stackName)

	deletetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"stack_id":             stackId,
			"retain_all_resources": d.Get("retain_all_resources"),
		},
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
			"Client-Request-Id": requestId,
		},
	}
	_, err = client.Request("POST", deletePath, &deletetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting stack (%s)", stackId))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      stackStatusRefreshFunc(client, stackId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
