package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsAlertRuleSimulation = []string{
	"workspace_id",
	"pipeline_id",
	"query_type",
	"from_time",
	"to_time",
	"event_grouping",
	"triggers",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/alert-rules/simulation
func ResourceAlertRuleSimulation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertRuleSimulationCreate,
		UpdateContext: resourceAlertRuleSimulationUpdate,
		ReadContext:   resourceAlertRuleSimulationRead,
		DeleteContext: resourceAlertRuleSimulationDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsAlertRuleSimulation),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the workspace to which the playbook belongs.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID of the alert rule.`,
			},
			"query_rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the query rule of the alert rule.`,
			},
			"query_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the query type of the alert rule.`,
			},
			"from_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time of the alert rule simulation.`,
			},
			"to_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time of the alert rule simulation.`,
			},
			"event_grouping": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Specifies whether to put events in a group.`,
			},
			"triggers": {
				Type:        schema.TypeList,
				MinItems:    1,
				Elem:        alertRuleSimulationTriggerSchema(),
				Required:    true,
				Description: `Specifies the triggers of the alert rule.`,
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

func alertRuleSimulationTriggerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the expression.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operator.`,
			},
			"accumulated_times": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the accumulated times.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the trigger mode.`,
			},
			"severity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the severity of the trigger.`,
			},
		},
	}
	return &sc
}

func resourceAlertRuleSimulationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAlertRuleSimulationHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/alert-rules/simulation"
		createAlertRuleSimulationProduct = "secmaster"
	)
	createAlertRuleSimulationClient, err := cfg.NewServiceClient(createAlertRuleSimulationProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createAlertRuleSimulationPath := createAlertRuleSimulationClient.Endpoint + createAlertRuleSimulationHttpUrl
	createAlertRuleSimulationPath = strings.ReplaceAll(createAlertRuleSimulationPath, "{project_id}", createAlertRuleSimulationClient.ProjectID)
	createAlertRuleSimulationPath = strings.ReplaceAll(createAlertRuleSimulationPath, "{workspace_id}", d.Get("workspace_id").(string))

	createAlertRuleSimulationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	bodyParams, err := buildCreateAlertRuleSimulationBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createAlertRuleSimulationOpt.JSONBody = utils.RemoveNil(bodyParams)

	_, err = createAlertRuleSimulationClient.Request("POST", createAlertRuleSimulationPath, &createAlertRuleSimulationOpt)
	if err != nil {
		return diag.Errorf("error creating alert rule simulation: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func buildCreateAlertRuleSimulationBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"pipe_id":        d.Get("pipeline_id"),
		"query":          d.Get("query_rule"),
		"query_type":     d.Get("query_type"),
		"event_grouping": utils.ValueIgnoreEmpty(d.Get("event_grouping")),
		"triggers":       d.Get("triggers"),
	}

	if v, ok := d.GetOk("from_time"); ok {
		startTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return nil, err
		}
		params["from"] = startTime * 1000
	}
	if v, ok := d.GetOk("to_time"); ok {
		startTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return nil, err
		}
		params["to"] = startTime * 1000
	}

	return params, nil
}

func resourceAlertRuleSimulationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlertRuleSimulationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlertRuleSimulationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for alert rule simulation resource. Deleting this resource will not change
		the status of the currently alert rule simulation resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
