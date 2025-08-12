package sfsturbo

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

var updatableChargeModeNonUpdatableParams = []string{
	"share_id",
	"period_num",
	"period_type",
	"is_auto_renew",
}

// @API SFSTurbo POST /v2/{project_id}/sfs-turbo/shares/{share_id}/change-charge-mode
func ResourceSFSTurboChangeChargeMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSTurboChangeChargeModeCreate,
		ReadContext:   resourceSFSTurboChangeChargeModeRead,
		UpdateContext: resourceSFSTurboChangeChargeModeUpdate,
		DeleteContext: resourceSFSTurboChangeChargeModeDelete,

		CustomizeDiff: config.FlexibleForceNew(updatableChargeModeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_auto_renew": {
				Type:     schema.TypeInt,
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

func buildSFSTurboChangeChargeModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bss_param": map[string]interface{}{
			"period_num":    d.Get("period_num"),
			"period_type":   d.Get("period_type"),
			"is_auto_renew": d.Get("is_auto_renew"),
			"is_auto_pay":   1,
		},
	}

	return bodyParams
}

func resourceSFSTurboChangeChargeModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
		httpUrl = "v2/{project_id}/sfs-turbo/shares/{share_id}/change-charge-mode"
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSFSTurboChangeChargeModeBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error changing SFS Turbo charge mode: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The api only create one order in fact, so just wait for first order completed.
	orderId := utils.PathSearch("order_ids|[0]", respBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error changing SFS Turbo charge mode: Order ID is not found in API response")
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuid)

	return nil
}

func resourceSFSTurboChangeChargeModeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceSFSTurboChangeChargeModeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceSFSTurboChangeChargeModeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to change SFS Turbo charge mode. Deleting this 
resource will not change the current charge mode result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
