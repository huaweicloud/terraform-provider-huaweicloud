package vpn

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

var gatewayUpgradeNonUpdatableParams = []string{"vgw_id", "action"}

// @API VPN POST /v5/{project_id}/vpn-gateways/{vgw_id}/upgrade
func ResourceGatewayUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatewayUpgradeCreate,
		ReadContext:   resourceGatewayUpgradeRead,
		UpdateContext: resourceGatewayUpgradeUpdate,
		DeleteContext: resourceGatewayUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(gatewayUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID of a VPN gateway.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies an upgrade operation.`,
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

func resourceGatewayUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	gatewayId := d.Get("vgw_id").(string)

	createHttpUrl := "v5/{project_id}/vpn-gateways/{vgw_id}/upgrade"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vgw_id}", gatewayId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": d.Get("action"),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating VPN gateway upgrade: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("job_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN gateway upgrade: job ID is not found in API response")
	}
	d.SetId(id)

	return nil
}

func resourceGatewayUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGatewayUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGatewayUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPN gateway upgrade resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
