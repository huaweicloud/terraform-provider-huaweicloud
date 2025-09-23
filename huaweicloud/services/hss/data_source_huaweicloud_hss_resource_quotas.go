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

// @API HSS GET /v5/{project_id}/billing/quotas
func DataSourceResourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the HSS version.",
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the billing mode.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource{Schema: quotaSchema},
				Description: "The quota information list.",
			},
		},
	}
}

var resourceSchema = map[string]*schema.Schema{
	"resource_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The resource ID.",
	},
	"current_time": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The current time.",
	},
	"shared_quota": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Whether the quota is shared.",
	},
}

var quotaSchema = map[string]*schema.Schema{
	"version": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The HSS version.",
	},
	"total_num": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The total quota number.",
	},
	"used_num": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The used quota number.",
	},
	"available_num": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The available quota number.",
	},
	"available_resources_list": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Resource{Schema: resourceSchema},
		Description: "The list of available resources.",
	},
}

func dataSourceResourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	requestPath := client.Endpoint + "v5/{project_id}/billing/quotas"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildResourceQuotasQueryParams(d, cfg)

	listQuotasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &listQuotasOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS resource quotas: %s", err)
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

	quotaList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenQuotaList(quotaList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourceQuotasQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("version"); ok {
		res = fmt.Sprintf("%s&version=%v", res, v)
	}
	if v, ok := d.GetOk("charging_mode"); ok {
		res = fmt.Sprintf("%s&charging_mode=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenQuotaList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"version":                  utils.PathSearch("version", v, nil),
			"total_num":                utils.PathSearch("total_num", v, nil),
			"used_num":                 utils.PathSearch("used_num", v, nil),
			"available_num":            utils.PathSearch("available_num", v, nil),
			"available_resources_list": flattenResourceList(utils.PathSearch("available_resources_list", v, nil)),
		})
	}
	return rst
}

func flattenResourceList(resources interface{}) []interface{} {
	if resources == nil {
		return nil
	}

	resourceList, ok := resources.([]interface{})
	if !ok {
		return nil
	}
	if len(resourceList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resourceList))
	for _, v := range resourceList {
		rst = append(rst, map[string]interface{}{
			"resource_id":  utils.PathSearch("resource_id", v, nil),
			"current_time": utils.PathSearch("current_time", v, nil),
			"shared_quota": utils.PathSearch("shared_quota", v, nil),
		})
	}
	return rst
}
