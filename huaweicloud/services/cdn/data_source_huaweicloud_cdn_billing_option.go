// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/charge/charge-modes
func DataSourceBillingOption() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBillingOptionRead,
		Schema: map[string]*schema.Schema{
			"product_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the product mode.`,
			},
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the service area`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the billing option status.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the billing option.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"effective_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the effective time of the option.`,
			},
		},
	}
}

func dataBillingOptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := model.ShowChargeModesRequest{
		ProductType: d.Get("product_type").(string),
		Status:      utils.StringIgnoreEmpty(d.Get("status").(string)),
		ServiceArea: utils.StringIgnoreEmpty(d.Get("service_area").(string)),
	}

	resp, err := hcCdnClient.ShowChargeModes(&request)
	if err != nil {
		return diag.Errorf("error retrieving CDN billing option: %s", err)
	}

	if resp == nil || resp.Result == nil {
		return diag.Errorf("error retrieving CDN billing option: Result is not found in API response")
	}

	if len(*resp.Result) == 0 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	resultArray := *resp.Result
	resultMap := resultArray[0]
	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("product_type", resultMap["product_type"]),
		d.Set("service_area", resultMap["service_area"]),
		d.Set("status", resultMap["status"]),
		d.Set("charge_mode", resultMap["charge_mode"]),
		d.Set("created_at", flattenTimeStamp(resultMap["create_time"])),
		d.Set("effective_time", flattenTimeStamp(resultMap["effective_time"])),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
