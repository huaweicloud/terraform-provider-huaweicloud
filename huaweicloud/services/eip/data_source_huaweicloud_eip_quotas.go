package eip

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v1/{project_id}/quotas
func DataSourceEipQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the quota type to filter. Valid values include publicIp, shareBandwidth, shareBandwidthIP.",
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": quotasSchema(),
					},
				},
			},
		},
	}
}

func quotasSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The type of the quota.",
				},
				"used": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The number of used quotas.",
				},
				"quota": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The total number of quotas.",
				},
				"min": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The minimum value of the quota that can be modified.",
				},
			},
		},
	}
}

func buildListQuotasQueryParams(d *schema.ResourceData) string {
	params := url.Values{}
	if v, ok := d.GetOk("type"); ok {
		params.Add("type", v.(string))
	}

	if len(params) > 0 {
		return "?" + params.Encode()
	}
	return ""
}

func dataSourceEipQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/quotas"
	)

	client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC EIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListQuotasQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error querying EIP quotas: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("quotas", flattenQuotasWrapper(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotasWrapper(respBody interface{}) []interface{} {
	quotasObj := utils.PathSearch("quotas", respBody, nil)
	if quotasObj == nil {
		return nil
	}

	resources := utils.PathSearch("resources", quotasObj, make([]interface{}, 0)).([]interface{})
	return []interface{}{
		map[string]interface{}{
			"resources": flattenEipQuotas(resources),
		},
	}
}

func flattenEipQuotas(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"used":  utils.PathSearch("used", v, 0),
			"quota": utils.PathSearch("quota", v, 0),
			"min":   utils.PathSearch("min", v, 0),
		})
	}
	return rst
}
