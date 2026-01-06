package vpcep

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

var vpcepEndpointUpgradeNonUpdatableParams = []string{"endpoint_id"}

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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the endpoint is located.`,
			},

			// Required parameters.
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the endpoint to be upgraded.`,
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

func resourceVPCEndpointUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		endpointId = d.Get("endpoint_id").(string)
		httpUrl    = "v2/{project_id}/vpc-endpoints/{vpc_endpoint_id}/upgrade"
	)

	client, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpc_endpoint_id}", endpointId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error upgrading VPCEP endpoint: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceVPCEndpointUpgradeRead(ctx, d, meta)
}

func resourceVPCEndpointUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVPCEndpointUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for upgrading VPCEP endpoint. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
