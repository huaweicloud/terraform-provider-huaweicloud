package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/available-flavors
func DataSourceRdsAvailableFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsAvailableFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone_ids": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec_code_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_category_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_rha_flavor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"optional_flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_ipv6_supported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_status": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tps": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_volume_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_volume_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsAvailableFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/available-flavors"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl + buildAvailableFlavorsQueryParams(d)
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var flavors []interface{}
	offset := 0
	for {
		getGaussDBAccountPath := getBasePath + buildAvailableFlavorsPageQueryParams(offset)
		getResp, err := client.Request("GET", getGaussDBAccountPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS available flavors: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res := flattenGetAvailableFlavorsResponseBody(getRespBody)
		flavors = append(flavors, res...)

		total := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if len(flavors) >= int(total) {
			break
		}
		offset++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("optional_flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAvailableFlavorsPageQueryParams(offset int) string {
	return fmt.Sprintf("&limit=100&offset=%d", offset)
}

func buildAvailableFlavorsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?availability_zone_ids=%s&ha_mode=%s", d.Get("availability_zone_ids").(string),
		d.Get("ha_mode").(string))
	if specCodeLike, ok := d.GetOk("spec_code_like"); ok {
		res = fmt.Sprintf("%s&spec_code_like=%v", res, specCodeLike)
	}
	if flavorCategoryType, ok := d.GetOk("flavor_category_type"); ok {
		res = fmt.Sprintf("%s&flavor_category_type=%v", res, flavorCategoryType)
	}
	if isRhaFlavor, ok := d.GetOk("is_rha_flavor"); ok {
		res = fmt.Sprintf("%s&is_rha_flavor=%v", res, isRhaFlavor)
	}
	return res
}

func flattenGetAvailableFlavorsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("optional_flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"vcpus":             utils.PathSearch("vcpus", v, nil),
			"ram":               utils.PathSearch("ram", v, nil),
			"spec_code":         utils.PathSearch("spec_code", v, nil),
			"is_ipv6_supported": utils.PathSearch("is_ipv6_supported", v, nil),
			"type_code":         utils.PathSearch("type_code", v, nil),
			"az_status":         utils.PathSearch("az_status", v, nil),
			"group_type":        utils.PathSearch("group_type", v, nil),
			"max_connection":    utils.PathSearch("max_connection", v, nil),
			"tps":               utils.PathSearch("tps", v, nil),
			"qps":               utils.PathSearch("qps", v, nil),
			"min_volume_size":   utils.PathSearch("min_volume_size", v, nil),
			"max_volume_size":   utils.PathSearch("max_volume_size", v, nil),
		})
	}
	return res
}
