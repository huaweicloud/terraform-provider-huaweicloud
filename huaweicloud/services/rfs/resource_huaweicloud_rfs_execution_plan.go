package rfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/{project_id}/stacks/{stack_name}/execution-plans
// @API RFS DELETE /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}
func ResourceExecutionPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExecutionPlanCreate,
		ReadContext:   resourceExecutionPlanRead,
		DeleteContext: resourceExecutionPlanDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the resource stack to which the execution plan belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the execution plan.`,
			},
			"stack_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the resource stack to which the execution plan belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The description of the execution plan.`,
			},
			"template_body": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The HCL/JSON template content for deployment resources.",
			},
			"vars_body": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The variable content for deployment resources.",
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The OBS address where the HCL/JSON template archive (**.zip** file, which contains all " +
					"resource **.tf.json** script files to be deployed) or **.tf.json** file is located, which " +
					"describes the target status of the deployment resources.",
			},
			"vars_uri": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The OBS address where the variable (**.tfvars**) file corresponding to the HCL/JSON " +
					"template located, which describes the target status of the deployment resources.",
			},
		},
	}
}

func resourceExecutionPlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/stacks/{stack_name}/execution-plans"
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
	createPath = strings.ReplaceAll(createPath, "{stack_name}", d.Get("stack_name").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
		},
		JSONBody: utils.RemoveNil(buildExecutionPlanCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating execution plan: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("execution_plan_id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the execution plan ID from the API response: %#v", respBody)
	}
	d.SetId(resourceId)

	return resourceExecutionPlanRead(ctx, d, meta)
}

func buildExecutionPlanCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"execution_plan_name": d.Get("name"),
		"stack_id":            utils.ValueIgnoreEmpty(d.Get("stack_id")),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"template_body":       utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":        utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_body":           utils.ValueIgnoreEmpty(d.Get("vars_body")),
		"vars_uri":            utils.ValueIgnoreEmpty(d.Get("vars_uri")),
	}
}

func resourceExecutionPlanRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceExecutionPlanDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}"
		executionPlanId = d.Id()
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{stack_name}", d.Get("stack_name").(string))
	deletePath = strings.ReplaceAll(deletePath, "{execution_plan_name}", d.Get("name").(string))

	deletetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
			"Content-Type":      "application/json",
			"X-Language":        "en-us",
		},
		JSONBody: map[string]interface{}{
			"execution_plan_id": executionPlanId,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deletetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting execution plan (%s)", executionPlanId))
	}

	return nil
}
