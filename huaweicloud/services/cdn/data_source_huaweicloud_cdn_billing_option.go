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

// @API CDN GET /v1.0/cdn/charge/charge-modes
func DataSourceBillingOption() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBillingOptionRead,
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
				Description: `Specifies the service area.`,
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

func buildDataSourceBillingOptionQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?product_type=%v", d.Get("product_type"))
	if v, ok := d.GetOk("service_area"); ok {
		queryParams = fmt.Sprintf("%s&service_area=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBillingOptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg  = meta.(*config.Config)
		mErr *multierror.Error
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getPath := client.Endpoint + "v1.0/cdn/charge/charge-modes"
	getPath += buildDataSourceBillingOptionQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CDN billing option: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	result := utils.PathSearch("result[0]", getRespBody, nil)
	if result == nil {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("product_type", utils.PathSearch("product_type", result, nil)),
		d.Set("service_area", utils.PathSearch("service_area", result, nil)),
		d.Set("status", utils.PathSearch("status", result, nil)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", result, nil)),
		d.Set("created_at", flattenCreatedAt(result)),
		d.Set("effective_time", flattenEffectiveTime(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
