// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/services/specifications
func DataSourceServiceFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceServiceFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"infer_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Inference mode.`,
			},
			"is_personal_cluster": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether this flavors is supported by dedicated resource pool.`,
			},
			"is_open": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether this flavor is open or not.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Flavor status.`,
			},
			"is_free": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the flavor is free of charge.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        serviceFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of flavors.`,
			},
		},
	}
}

func serviceFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"is_open": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this flavor is open or not.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Flavor status.`,
			},
			"billing_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `ID of the billing specifications.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model type, which can be empty or **auto**.`,
			},
			"is_free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the flavor is free of charge.`,
			},
			"over_quota": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the quota exceeds the upper limit.`,
			},
			"extend_params": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Billing item.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Description of the flavor.`,
			},
		},
	}
	return &sc
}

func resourceServiceFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of ModelArts service deployment flavors
	var (
		listFlavorsHttpUrl = "v1/{project_id}/services/specifications"
		listFlavorsProduct = "modelarts"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsqueryParams := buildListServiceFlavorsQueryParams(d)
	listFlavorsPath += listFlavorsqueryParams

	listFlavorsResp, err := pagination.ListAllItems(
		listFlavorsClient,
		"offset",
		listFlavorsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving service flavors")
	}

	listFlavorsRespJson, err := json.Marshal(listFlavorsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listFlavorsRespBody interface{}
	err = json.Unmarshal(listFlavorsRespJson, &listFlavorsRespBody)
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
		d.Set("flavors", filterServiceFlavors(flattenListServiceFlavorsFlavors(listFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListServiceFlavorsFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("specifications", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("specification", v, nil),
			"is_open":       utils.PathSearch("is_open", v, nil),
			"status":        utils.PathSearch("spec_status", v, nil),
			"billing_spec":  utils.PathSearch("billing_spec", v, nil),
			"source_type":   utils.PathSearch("source_type", v, nil),
			"is_free":       utils.PathSearch("is_free", v, nil),
			"over_quota":    utils.PathSearch("over_quota", v, nil),
			"extend_params": utils.PathSearch("extend_params", v, nil),
			"description":   utils.PathSearch("display_en", v, nil),
		})
	}
	return rst
}

func filterServiceFlavors(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("is_open"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("is_open", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("is_free"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("is_free", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListServiceFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("infer_type"); ok {
		res = fmt.Sprintf("%s&infer_type=%v", res, v)
	}

	if v, ok := d.GetOk("is_personal_cluster"); ok {
		res = fmt.Sprintf("%s&is_personal_cluster=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
