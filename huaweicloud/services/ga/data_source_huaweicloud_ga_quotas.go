package ga

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/{domain_id}/ga/quotas
func DataSourceGaQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaQuotasRead,

		Schema: map[string]*schema.Schema{
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of quotas.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The quota mark.`,
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The minimum quota threshold.`,
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum quota threshold.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The quota size.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ga"
		httpUrl = "v1/{domain_id}/ga/quotas"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GA quotas: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(nil,
		d.Set("quotas", flattenQuotas(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(quotasResp interface{}) []interface{} {
	if quotasResp == nil {
		return make([]interface{}, 0)
	}

	resources := utils.PathSearch("quotas.resources", quotasResp, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		resource := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", resource, nil),
			"min":   utils.PathSearch("min", resource, nil),
			"max":   utils.PathSearch("max", resource, nil),
			"quota": utils.PathSearch("quota", resource, nil),
		})
	}

	return rst
}
