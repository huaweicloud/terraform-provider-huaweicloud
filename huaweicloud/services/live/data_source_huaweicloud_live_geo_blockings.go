package live

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE GET /v1/{project_id}/domain/geo-blocking
func DataSourceGeoBlockings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeoBlockingsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ingest domain name to which the recording rules belong.`,
			},
			"apps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the recording rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The recording rule ID.`,
						},
						"area_whitelist": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The default recording configuration rule.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceGeoBlockingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "live"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	respBody, err := ReadGeoBlocking(client, d.Get("domain_name").(string))
	if err != nil {
		return diag.Errorf("error retrieving Live geo blockings: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	apps := utils.PathSearch("apps", respBody, make([]interface{}, 0)).([]interface{})
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("apps", flattenApps(apps)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApps(apps []interface{}) []map[string]interface{} {
	if len(apps) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(apps))
	for i, v := range apps {
		result[i] = map[string]interface{}{
			"app_name":       utils.PathSearch("app", v, nil),
			"area_whitelist": utils.PathSearch("area_whitelist", v, nil),
		}
	}

	return result
}
