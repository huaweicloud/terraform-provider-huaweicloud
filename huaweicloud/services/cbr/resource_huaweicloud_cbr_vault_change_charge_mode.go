package cbr

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableChargeModeParams = []string{
	"vault_ids",
	"charging_mode",
	"period_type",
	"period_num",
	"is_auto_renew",
}

// @API CBR POST /v3/{project_id}/vaults/change-charge-mode
func ResourceVaultChangeChargeMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultChangeChargeModeCreate,
		ReadContext:   resourceVaultChangeChargeModeRead,
		UpdateContext: resourceVaultChangeChargeModeUpdate,
		DeleteContext: resourceVaultChangeChargeModeDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableChargeModeParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
region will be used.`,
			},
			"vault_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the IDs of the vaults to change charge mode.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// This field `charging_mode` is optional in the openapi document, but it is required in the actual test.
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the charging mode of the vault.`,
			},
			"period_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the period type of the vault.`,
			},
			"period_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of periods to purchase.`,
			},
			"is_auto_renew": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to auto-renew the vault when it expires.`,
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

func buildVaultChangeChargeModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	chargeModeMap := map[string]interface{}{
		"vault_ids":     d.Get("vault_ids"),
		"charging_mode": d.Get("charging_mode"),
		"period_type":   d.Get("period_type"),
		"period_num":    d.Get("period_num"),
		"is_auto_pay":   true,
	}

	if d.Get("is_auto_renew").(bool) {
		chargeModeMap["is_auto_renew"] = true
	}
	return chargeModeMap
}

func resourceVaultChangeChargeModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/vaults/change-charge-mode"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildVaultChangeChargeModeBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error changing vault charge mode: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("orderId", respBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error getting order ID from response: %v", respBody)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(id)
	return resourceVaultChangeChargeModeRead(ctx, d, meta)
}

func resourceVaultChangeChargeModeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time action resource, so we just return without any data
	return nil
}

func resourceVaultChangeChargeModeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time action resource, so we just return without any data
	return nil
}

func resourceVaultChangeChargeModeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to change vault charge mode. Deleting this 
resource will not change the current charge mode result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
