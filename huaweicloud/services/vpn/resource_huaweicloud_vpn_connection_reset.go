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

var connectionResetNonUpdatableParams = []string{"connection_id"}

// @API VPN POST /v5/{project_id}/vpn-connection/{vpn_connection_id}/reset
// @API VPN GET /v5/{project_id}/vpn-connection/{vpn_connection_id}
func ResourceConnectionReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionResetCreate,
		ReadContext:   resourceConnectionResetRead,
		UpdateContext: resourceConnectionResetUpdate,
		DeleteContext: resourceConnectionResetDelete,

		CustomizeDiff: config.FlexibleForceNew(connectionResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the connection ID of a VPN gateway.`,
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

func resourceConnectionResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	connectionId := d.Get("connection_id").(string)

	createHttpUrl := "v5/{project_id}/vpn-connection/{vpn_connection_id}/reset"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpn_connection_id}", connectionId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating VPN gateway connection reset: %s", err)
	}

	d.SetId(connectionId)

	err = createConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for reseting VPN connection (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func resourceConnectionResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConnectionResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConnectionResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPN gateway connection reset resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
