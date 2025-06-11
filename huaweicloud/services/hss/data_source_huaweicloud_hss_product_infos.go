package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/product/productdata/offering-infos
func DataSourceProductInfos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductInfosRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"site_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the site information.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        productSchema(),
				Description: "The product information list.",
			},
		},
	}
}

func periodSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"period_vals": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value string of the required duration.",
			},
			"period_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Required duration unit.",
			},
		},
	}
}

func productSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The billing modes.",
			},
			"is_auto_renew": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable automatic renewal.",
			},
			"version_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"periods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     periodSchema(),
						},
					},
				},
				Description: "The edition information.",
			},
		},
	}
}

func dataSourceProductInfosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/product/productdata/offering-infos"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildProductInfosQueryParams(d, cfg)

	listProductInfosOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &listProductInfosOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS product information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenProductInfos(utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildProductInfosQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsID := cfg.GetEnterpriseProjectID(d)

	if v, ok := d.GetOk("site_code"); ok {
		res = fmt.Sprintf("%s&site_code=%v", res, v)
	}
	if epsID != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsID)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenProductInfos(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		versionInfo := utils.PathSearch("version_info", v, make(map[string]interface{})).(map[string]interface{})
		versionList := flattenVersionInfoList(versionInfo)

		rst = append(rst, map[string]interface{}{
			"charging_mode": utils.PathSearch("charging_mode", v, nil),
			"is_auto_renew": utils.PathSearch("is_auto_renew", v, nil),
			"version_info":  versionList,
		})
	}
	return rst
}

func flattenVersionInfoList(raw map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(raw))
	for ver, periods := range raw {
		result = append(result, map[string]interface{}{
			"version": ver,
			"periods": flattenPeriodList(periods),
		})
	}
	return result
}

func flattenPeriodList(raw interface{}) []interface{} {
	periods, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(periods))
	for _, p := range periods {
		result = append(result, map[string]interface{}{
			"period_vals": utils.PathSearch("period_vals", p, nil),
			"period_unit": utils.PathSearch("period_unit", p, nil),
		})
	}
	return result
}
