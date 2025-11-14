package aom

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

var uniAgentBatchUpgradeNonUpdatableParams = []string{"version", "agent_list"}

// @API AOM POST /v1/{project_id}/uniagent-console/upgrade/batch-upgrade
func ResourceUniAgentBatchUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUniAgentBatchUpgradeCreate,
		ReadContext:   resourceUniAgentBatchUpgradeRead,
		UpdateContext: resourceUniAgentBatchUpgradeUpdate,
		DeleteContext: resourceUniAgentBatchUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(uniAgentBatchUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the target machines to be operated are located.`,
			},

			// Required parameters.
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version number of UniAgent to be upgraded.`,
			},
			"agent_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        uniAgentBatchUpgradeAgentListSchema(),
				Description: `The list of host information for upgrading UniAgent.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func uniAgentBatchUpgradeAgentListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The unique agent ID.`,
			},
			"inner_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The host IP address.`,
			},
		},
	}
}

func buildUniAgentBatchUpgradeAgentList(agentList []interface{}) []interface{} {
	if len(agentList) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(agentList))
	for _, agent := range agentList {
		result = append(result, map[string]interface{}{
			"agent_id": utils.PathSearch("agent_id", agent, nil),
			"inner_ip": utils.PathSearch("inner_ip", agent, nil),
		})
	}
	return result
}

func buildUniAgentBatchUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"version":    d.Get("version"),
		"agent_list": buildUniAgentBatchUpgradeAgentList(d.Get("agent_list").([]interface{})),
	}
}

func resourceUniAgentBatchUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/uniagent-console/upgrade/batch-upgrade"
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUniAgentBatchUpgradeBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating UniAgent batch upgrade task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	state := utils.PathSearch("state", respBody, false).(bool)
	if !state {
		return diag.Errorf("error dispatching UniAgent batch upgrade task")
	}

	return resourceUniAgentBatchUpgradeRead(ctx, d, meta)
}

func resourceUniAgentBatchUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUniAgentBatchUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUniAgentBatchUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch upgrading AOM UniAgents. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
