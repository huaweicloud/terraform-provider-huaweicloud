package secmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsWorkflowAction = []string{
	"workspace_id",
	"workflow_id",
	"command_type",
	"action_type",
	"action_instance_id",
	"playbook_context",
	"simulation_context",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/instances
func ResourceWorkflowAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowActionCreate,
		UpdateContext: resourceWorkflowActionUpdate,
		ReadContext:   resourceWorkflowActionRead,
		DeleteContext: resourceWorkflowActionDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsWorkflowAction),

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
			"workflow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"command_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"playbook_context": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"simulation_context": {
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

func bulidPlaybookContextBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	v, ok := d.GetOk("playbook_context")
	if !ok {
		return nil, nil
	}

	var playbookContext map[string]interface{}
	err := json.Unmarshal([]byte(v.(string)), &playbookContext)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling playbook_context json: %s", err)
	}

	return playbookContext, err
}

func bulidSimulationContextBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	v, ok := d.GetOk("simulation_context")
	if !ok {
		return nil, nil
	}

	var simulationContext map[string]interface{}
	err := json.Unmarshal([]byte(v.(string)), &simulationContext)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling playbook_context json: %s", err)
	}

	return simulationContext, err
}

func resourceWorkflowActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createWorkflowActionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/instances"
		createWorkflowActionProduct = "secmaster"
	)
	createWorkflowActionClient, err := cfg.NewServiceClient(createWorkflowActionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createWorkflowActionPath := createWorkflowActionClient.Endpoint + createWorkflowActionHttpUrl
	createWorkflowActionPath = strings.ReplaceAll(createWorkflowActionPath, "{project_id}", createWorkflowActionClient.ProjectID)
	createWorkflowActionPath = strings.ReplaceAll(createWorkflowActionPath, "{workspace_id}", d.Get("workspace_id").(string))
	createWorkflowActionPath = strings.ReplaceAll(createWorkflowActionPath, "{workflow_id}", d.Get("workflow_id").(string))

	createWorkflowActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	playbookContext, err := bulidPlaybookContextBodyParams(d)
	if err != nil {
		return nil
	}

	simulationContext, err := bulidSimulationContextBodyParams(d)
	if err != nil {
		return nil
	}

	createOpts := map[string]interface{}{
		"command_type":       d.Get("command_type"),
		"action_type":        d.Get("action_type"),
		"action_id":          d.Get("workflow_id"),
		"action_instance_id": utils.ValueIgnoreEmpty(d.Get("action_instance_id")),
		"playbook_context":   playbookContext,
		"simulation_context": simulationContext,
	}

	createWorkflowActionOpt.JSONBody = utils.RemoveNil(createOpts)

	resp, err := createWorkflowActionClient.Request("POST", createWorkflowActionPath, &createWorkflowActionOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster workflow action: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return nil
}

func resourceWorkflowActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for SecMaster workflow action resource. Deleting this resource will not
	 change the status of the currently SecMaster workflow action resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
