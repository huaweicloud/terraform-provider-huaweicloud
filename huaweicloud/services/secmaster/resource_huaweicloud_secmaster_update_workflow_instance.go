package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances/{instance_id}
func ResourceUpdateWorkflowInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpdateWorkflowInstanceCreate,
		UpdateContext: resourceUpdateWorkflowInstanceUpdate,
		ReadContext:   resourceUpdateWorkflowInstanceRead,
		DeleteContext: resourceUpdateWorkflowInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"workspace_id",
			"instance_id",
			"command_type",
			"task_id",
			"input_dataobject",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"command_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"input_dataobject": {
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

func buildUpdateWorkflowInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"command_type":     d.Get("command_type"),
		"task_id":          utils.ValueIgnoreEmpty(d.Get("task_id")),
		"input_dataobject": utils.ValueIgnoreEmpty(d.Get("input_dataobject")),
	}
}

func resourceUpdateWorkflowInstanceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances/{instance_id}"
		workspaceId = d.Get("workspace_id").(string)
		instanceId  = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildUpdateWorkflowInstanceBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster workflow instance: %s", err)
	}

	d.SetId(instanceId)

	return nil
}

func resourceUpdateWorkflowInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateWorkflowInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateWorkflowInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource using to update workflow instance. Deleting this resource will not
	 change the status of the currently SecMaster workflow instance resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
