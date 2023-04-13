// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CBH
// ---------------------------------------------------------------

package cbh

import (
	"context"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCbhFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCbhFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the project.",
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        FlavorsFlavorSchema(),
				Required:    true,
				Description: `Indicates the list of the product info.`,
			},
		},
	}
}

func FlavorsFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_spec": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource specifications of cloud service types.`,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cloud service region code.`,
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the period type of a yearly/monthly product order.`,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of periods of a yearly/monthly product order.`,
			},
			"subscription_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of subscriptions of a yearly/monthly product order.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the product.`,
			},
		},
	}
	return &sc
}

func resourceCbhFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCbhFlavors: Query the List of CBH flavors
	var (
		getCbhFlavorsHttpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		getCbhFlavorsProduct = "bss"
	)
	getCbhFlavorsClient, err := cfg.NewServiceClient(getCbhFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating BSS Client: %s", err)
	}

	getCbhFlavorsPath := getCbhFlavorsClient.Endpoint + getCbhFlavorsHttpUrl

	getCbhFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCbhFlavorsOpt.JSONBody = utils.RemoveNil(buildGetCbhFlavorsBodyParams(d))
	getCbhFlavorsResp, err := getCbhFlavorsClient.Request("POST", getCbhFlavorsPath, &getCbhFlavorsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CbhFlavors")
	}

	getCbhFlavorsRespBody, err := utils.FlattenResponse(getCbhFlavorsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavors", flattenGetFlavorsResponseBodyFlavor(d, getCbhFlavorsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetFlavorsResponseBodyFlavor(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("official_website_rating_result.product_rating_results",
		resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	flavorMap := generateFlavorMap(curArray)
	flavors := d.Get("flavors").([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for index, flavor := range flavors {
		rst = append(rst, map[string]interface{}{
			"resource_spec":    utils.PathSearch("resource_spec", flavor, nil),
			"region":           utils.PathSearch("region", flavor, nil),
			"period_unit":      utils.PathSearch("period_unit", flavor, nil),
			"period":           utils.PathSearch("period", flavor, nil),
			"subscription_num": utils.PathSearch("subscription_num", flavor, nil),
			"flavor_id":        flavorMap[strconv.Itoa(index)],
		})
	}
	return rst
}

func generateFlavorMap(flavors []interface{}) map[string]string {
	res := make(map[string]string)
	for _, flavor := range flavors {
		id := utils.PathSearch("id", flavor, "")
		flavorId := utils.PathSearch("product_id", flavor, "")
		res[id.(string)] = flavorId.(string)
	}
	return res
}

func buildGetCbhFlavorsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":    utils.ValueIngoreEmpty(d.Get("project_id")),
		"product_infos": buildGetCbhFlavorsFlavorsChildBody(d),
	}
	return bodyParams
}

func buildGetCbhFlavorsFlavorsChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("flavors").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0)
	for index, param := range rawParams {
		perm := make(map[string]interface{})
		periodType := utils.PathSearch("period_unit", param, nil)
		if periodType == "month" {
			periodType = "2"
		} else {
			periodType = "3"
		}
		perm["id"] = strconv.Itoa(index)
		perm["cloud_service_type"] = "hws.service.type.cbh"
		perm["resource_type"] = "hws.resource.type.cbh.ins"
		perm["resource_spec"] = utils.PathSearch("resource_spec", param, nil)
		perm["region"] = utils.PathSearch("region", param, nil)
		perm["period_type"] = periodType
		perm["period_num"] = utils.PathSearch("period", param, nil)
		perm["subscription_num"] = utils.PathSearch("subscription_num", param, nil)
		params = append(params, perm)
	}
	return params
}
