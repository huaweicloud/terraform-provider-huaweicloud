// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CPH
// ---------------------------------------------------------------

package cph

import (
	"context"
	"fmt"
	"strings"

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

// @API CPH GET /v1/{project_id}/cloud-phone/phone-models
func DataSourcePhoneFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePhoneFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1",
				Description: `The flavor status.`,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cloud phone type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"0", "1",
				}, false),
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The vcpus of the CPH phone.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The ram of the CPH phone in MB.`,
			},
			"server_flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The CPH server flavor.`,
			},
			"image_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The label of image.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        phoneFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of flavor detail.`,
			},
		},
	}
}

func phoneFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the flavor.`,
			},
			"server_flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the CPH server flavor.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The vcpus of the CPH phone.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ram of the CPH phone in MB.`,
			},
			"disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The storage size in GB.`,
			},
			"resolution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resolution of the CPH phone.`,
			},
			"phone_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of cloud phones of the current flavor.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor status.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud phone type.`,
			},
			"extend_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description, which is a string in JSON format and can contain a maximum of 512 bytes.`,
			},
			"image_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The label of image.`,
			},
		},
	}
	return &sc
}

func resourcePhoneFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of CPH phone flavors
	var (
		listFlavorsHttpUrl = "v1/{project_id}/cloud-phone/phone-models"
		listFlavorsProduct = "cph"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsqueryParams := buildListPhoneFlavorsQueryParams(d)
	listFlavorsPath += listFlavorsqueryParams

	listFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listFlavorsResp, err := listFlavorsClient.Request("GET", listFlavorsPath, &listFlavorsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving phone flavors")
	}

	listFlavorsRespBody, err := utils.FlattenResponse(listFlavorsResp)
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
		d.Set("flavors", filterListPhoneModelsFlavors(
			flattenListPhoneModelsFlavors(listFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPhoneModelsFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("phone_models", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id":        utils.PathSearch("phone_model_name", v, nil),
			"server_flavor_id": utils.PathSearch("server_model_name", v, nil),
			"vcpus":            utils.PathSearch("cpu", v, nil),
			"memory":           utils.PathSearch("memory", v, nil),
			"disk":             utils.PathSearch("disk", v, nil),
			"resolution":       utils.PathSearch("resolution", v, nil),
			"phone_capacity":   utils.PathSearch("phone_capacity", v, nil),
			"status":           fmt.Sprint(utils.PathSearch("status", v, "")),
			"type":             fmt.Sprint(utils.PathSearch("product_type", v, "")),
			"image_label":      utils.PathSearch("image_label", v, ""),
			"extend_spec":      utils.PathSearch("extend_spec", v, nil),
		})
	}
	return rst
}

func filterListPhoneModelsFlavors(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("vcpus"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("vcpus", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("memory"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("memory", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("server_flavor_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("server_flavor_id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("image_label"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("image_label", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListPhoneFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
