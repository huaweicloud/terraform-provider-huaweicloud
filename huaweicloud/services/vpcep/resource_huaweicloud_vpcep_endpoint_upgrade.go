package vpcep

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

var vpcepEndpointUpgradeNonUpdatableParams = []string{
	"vpc_endpoint_id", "action",
}

// @API VPCEP POST /v2/{project_id}/vpc-endpoints/{vpc_endpoint_id}/upgrade
func ResourceVPCEndpointUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEndpointUpgradeCreate,
		ReadContext:   resourceVPCEndpointUpgradeRead,
		UpdateContext: resourceVPCEndpointUpgradeUpdate,
		DeleteContext: resourceVPCEndpointUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(vpcepEndpointUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceVPCEndpointUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v2/{project_id}/vpc-endpoints/{vpc_endpoint_id}/upgrade"
		vpcEndpointId = d.Get("vpc_endpoint_id").(string)
	)

	client, err := cfg.NewServiceClient("vpcep", region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpc_endpoint_id}", vpcEndpointId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUpgradeRequestBody(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error upgrading VPC endpoint: %s", err)
	}

	d.SetId(vpcEndpointId)

	return nil
}

func buildUpgradeRequestBody(d *schema.ResourceData) map[string]interface{} {
	requestBody := map[string]interface{}{
		"action": d.Get("action").(string),
	}
	return requestBody
}

func resourceVPCEndpointUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPC endpoint upgrade resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
