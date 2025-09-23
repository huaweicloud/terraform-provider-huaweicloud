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

var p2CGatewayConnectionDisconnectNonUpdatableParams = []string{"p2c_vgw_id", "connection_id"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/connections/{connection_id}/disconnect
func ResourceP2CGatewayConnectionDisconnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceP2CGatewayConnectionDisconnectCreate,
		ReadContext:   resourceP2CGatewayConnectionDisconnectRead,
		UpdateContext: resourceP2CGatewayConnectionDisconnectUpdate,
		DeleteContext: resourceP2CGatewayConnectionDisconnectDelete,

		CustomizeDiff: config.FlexibleForceNew(p2CGatewayConnectionDisconnectNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"p2c_vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID of a P2C VPN gateway.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the connection ID of a P2C VPN gateway.`,
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

func resourceP2CGatewayConnectionDisconnectCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	p2cGateWayId := d.Get("p2c_vgw_id").(string)
	connectionId := d.Get("connection_id").(string)

	createHttpUrl := "v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/connections/{connection_id}/disconnect"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{p2c_vgw_id}", p2cGateWayId)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", connectionId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating VPN P2C gateway connection disconnect: %s", err)
	}

	d.SetId(p2cGateWayId + "/" + connectionId)

	return nil
}

func resourceP2CGatewayConnectionDisconnectRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceP2CGatewayConnectionDisconnectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceP2CGatewayConnectionDisconnectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPN P2C gateway connection disconnect resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
