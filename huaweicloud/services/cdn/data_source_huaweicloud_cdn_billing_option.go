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
			// Required parameters
			"product_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The product mode.`,
			},

			// Optional parameters
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The service area.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The billing option status.`,
			},

			// Attributes
			"charge_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing option.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time, in RFC3339 format.`,
			},
			"effective_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The effective time, in RFC3339 format.`,
			},
		},
	}
}

func buildDataSourceBillingOptionQueryParams(d *schema.ResourceData) string {
	res := "?"

	res = fmt.Sprintf("%sproduct_type=%v", res, d.Get("product_type"))

	if v, ok := d.GetOk("service_area"); ok {
		res = fmt.Sprintf("%s&service_area=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	return res
}

func dataSourceBillingOptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
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
		return diag.Errorf("error retrieving billing option: %s", err)
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

	mErr := multierror.Append(
		d.Set("product_type", utils.PathSearch("product_type", result, nil)),
		d.Set("service_area", utils.PathSearch("service_area", result, nil)),
		d.Set("status", utils.PathSearch("status", result, nil)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", result, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", result, float64(0)).(float64))/1000, false)),
		d.Set("effective_time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("effective_time", result, float64(0)).(float64))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
