// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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

func buildChangeBillingOptionBodyParams(d *schema.ResourceData) interface{} {
	return map[string]interface{}{
		"charge_mode":  d.Get("charge_mode").(string),
		"product_type": d.Get("product_type").(string),
		"service_area": d.Get("service_area").(string),
	}
}

func changeBillingOption(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	chargePath := client.Endpoint + "v1.0/cdn/charge/charge-modes"
	chargeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildChangeBillingOptionBodyParams(d),
	}
	_, err := client.Request("PUT", chargePath, &chargeOpt)

	return err
}

func resourceBillingOptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if err := changeBillingOption(client, d); err != nil {
		return diag.Errorf("error modifying CDN billing option in creation operation: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceBillingOptionRead(ctx, d, meta)
}

func buildBillingOptionQueryParams(productType string) string {
	return fmt.Sprintf("?product_type=%v", productType)
}

func GetBillingOptionDetail(client *golangsdk.ServiceClient, productType string) (interface{}, error) {
	getPath := client.Endpoint + "v1.0/cdn/charge/charge-modes"
	getPath += buildBillingOptionQueryParams(productType)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// The billing model is always valuable and there is no need to pay attention to scenarios where resource does
		// not exist.
		return nil, fmt.Errorf("error retrieving CDN billing option: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	result := utils.PathSearch("result[0]", getRespBody, nil)
	if result == nil {
		return nil, fmt.Errorf("error retrieving CDN billing option: result is not found in API response")
	}

	return result, nil
}

func resourceBillingOptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	result, err := GetBillingOptionDetail(client, d.Get("product_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("product_type", utils.PathSearch("product_type", result, nil)),
		d.Set("service_area", utils.PathSearch("service_area", result, nil)),
		d.Set("created_at", flattenCreatedAt(result)),
		d.Set("effective_time", flattenEffectiveTime(result)),
		d.Set("status", utils.PathSearch("status", result, nil)),
		d.Set("current_charge_mode", utils.PathSearch("charge_mode", result, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEffectiveTime(result interface{}) string {
	effectiveTime := utils.PathSearch("effective_time", result, float64(0)).(float64)
	return utils.FormatTimeStampRFC3339(int64(effectiveTime)/1000, false)
}

func resourceBillingOptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if err := changeBillingOption(client, d); err != nil {
		return diag.Errorf("error modifying CDN billing option in update operation: %s", err)
	}

	return resourceBillingOptionRead(ctx, d, meta)
}

func resourceBillingOptionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
