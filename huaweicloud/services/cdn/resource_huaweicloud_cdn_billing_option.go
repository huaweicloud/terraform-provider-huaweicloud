// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cdnv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN PUT /v1.0/cdn/charge/charge-modes
// @API CDN GET /v1.0/cdn/charge/charge-modes
func ResourceBillingOption() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBillingOptionCreate,
		UpdateContext: resourceBillingOptionUpdate,
		ReadContext:   resourceBillingOptionRead,
		DeleteContext: resourceBillingOptionDelete,

		Schema: map[string]*schema.Schema{
			"charge_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the billing option.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the product mode.`,
			},
			"service_area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the service area.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"effective_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The effective time of the option.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status.`,
			},
			"current_charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing option of the account.`,
			},
		},
	}
}

func resourceBillingOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	if err := changeBillingOption(hcCdnClient, d); err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return resourceBillingOptionRead(ctx, d, meta)
}

func resourceBillingOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	if err := changeBillingOption(hcCdnClient, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceBillingOptionRead(ctx, d, meta)
}

func changeBillingOption(hcCdnClient *cdnv2.CdnClient, d *schema.ResourceData) error {
	req := model.SetChargeModesRequest{
		Body: &model.SetChargeModesBody{
			ChargeMode:  d.Get("charge_mode").(string),
			ProductType: d.Get("product_type").(string),
			ServiceArea: d.Get("service_area").(string),
		},
	}
	_, err := hcCdnClient.SetChargeModes(&req)
	if err != nil {
		return fmt.Errorf("error modifying CDN billing options: %s", err)
	}
	return nil
}

func resourceBillingOptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := model.ShowChargeModesRequest{
		ProductType: d.Get("product_type").(string),
	}

	resp, err := hcCdnClient.ShowChargeModes(&request)
	if err != nil {
		// The billing model is always valuable and there is no need to pay attention to scenarios where resources do not exist.
		return diag.Errorf("error retrieving CDN billing option: %s", err)
	}

	if resp == nil || resp.Result == nil || len(*resp.Result) == 0 {
		return diag.Errorf("error retrieving CDN billing option: Result is not found in API response")
	}

	var mErr *multierror.Error
	resultArray := *resp.Result
	resultMap := resultArray[0]

	mErr = multierror.Append(
		mErr,
		d.Set("product_type", resultMap["product_type"]),
		d.Set("service_area", resultMap["service_area"]),
		d.Set("created_at", flattenTimeStamp(resultMap["create_time"])),
		d.Set("effective_time", flattenTimeStamp(resultMap["effective_time"])),
		d.Set("status", resultMap["status"]),
		d.Set("current_charge_mode", resultMap["charge_mode"]),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTimeStamp(timestamp interface{}) string {
	if timestamp == nil {
		return ""
	}
	timeInt64, err := timestamp.(json.Number).Int64()
	if err != nil {
		return ""
	}
	return utils.FormatTimeStampRFC3339(timeInt64/1000, false)
}

func resourceBillingOptionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
