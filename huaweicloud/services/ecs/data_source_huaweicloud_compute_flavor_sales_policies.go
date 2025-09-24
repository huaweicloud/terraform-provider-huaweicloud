package ecs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS GET /v1/{project_id}/cloudservers/flavor-sell-policies
func DataSourceEcsComputeFlavorSalesPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEcsComputeFlavorSalesPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sell_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sell_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"longest_spot_duration_hours_gt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"largest_spot_duration_count_gt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"longest_spot_duration_hours": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"largest_spot_duration_count": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interruption_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sell_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sell_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sell_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interruption_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"longest_spot_duration_hours": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"largest_spot_duration_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEcsComputeFlavorSalesPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/{project_id}/cloudservers/flavor-sell-policies"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var marker string
	limit := 1000
	res := make([]map[string]interface{}, 0)
	for {
		queryParams, err := buildGetFlavorSalesPoliciesQueryParams(d, limit, marker)
		if err != nil {
			return diag.Errorf("error building flavor sales policies query params: %s", err)
		}
		getPath := getBasePath + queryParams
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving ECS flavor sales policies: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		flavorSalesPolicies, nextMarker := flattenFlavorSalesPolicies(getRespBody)
		if len(flavorSalesPolicies) == 0 {
			break
		}

		res = append(res, flavorSalesPolicies...)
		marker = nextMarker
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("sell_policies", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
func buildGetFlavorSalesPoliciesQueryParams(d *schema.ResourceData, limit int, marker string) (string, error) {
	res := fmt.Sprintf("?limit=%d", limit)
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	if v, ok := d.GetOk("flavor_id"); ok {
		res = fmt.Sprintf("%s&flavor_id=%v", res, v)
	}
	if v, ok := d.GetOk("sell_status"); ok {
		res = fmt.Sprintf("%s&sell_status=%v", res, v)
	}
	if v, ok := d.GetOk("sell_mode"); ok {
		res = fmt.Sprintf("%s&sell_mode=%v", res, v)
	}
	if v, ok := d.GetOk("availability_zone_id"); ok {
		res = fmt.Sprintf("%s&availability_zone_id=%v", res, v)
	}
	if v, ok := d.GetOk("longest_spot_duration_hours_gt"); ok {
		longestSpotDurationHoursGt, err := strconv.Atoi(v.(string))
		if err != nil {
			return "", fmt.Errorf("invalid parameter longest_spot_duration_hours_gt: %s", err)
		}
		res = fmt.Sprintf("%s&longest_spot_duration_hours_gt=%v", res, longestSpotDurationHoursGt)
	}
	if v, ok := d.GetOk("largest_spot_duration_count_gt"); ok {
		largestSpotDurationCountGt, err := strconv.Atoi(v.(string))
		if err != nil {
			return "", fmt.Errorf("invalid parameter largest_spot_duration_count_gt: %s", err)
		}
		res = fmt.Sprintf("%s&largest_spot_duration_count_gt=%v", res, largestSpotDurationCountGt)
	}
	if v, ok := d.GetOk("longest_spot_duration_hours"); ok {
		longestSpotDurationHours, err := strconv.Atoi(v.(string))
		if err != nil {
			return "", fmt.Errorf("invalid parameter longest_spot_duration_hours: %s", err)
		}
		res = fmt.Sprintf("%s&longest_spot_duration_hours=%v", res, longestSpotDurationHours)
	}
	if v, ok := d.GetOk("largest_spot_duration_count"); ok {
		largestSpotDurationCount, err := strconv.Atoi(v.(string))
		if err != nil {
			return "", fmt.Errorf("invalid parameter largest_spot_duration_count: %s", err)
		}
		res = fmt.Sprintf("%s&largest_spot_duration_count=%v", res, largestSpotDurationCount)
	}
	if v, ok := d.GetOk("interruption_policy"); ok {
		res = fmt.Sprintf("%s&interruption_policy=%v", res, v)
	}
	return res, nil
}

func flattenFlavorSalesPolicies(resp interface{}) ([]map[string]interface{}, string) {
	curJson := utils.PathSearch("sell_policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"flavor_id":            utils.PathSearch("flavor_id", v, nil),
			"sell_status":          utils.PathSearch("sell_status", v, nil),
			"availability_zone_id": utils.PathSearch("availability_zone_id", v, nil),
			"sell_mode":            utils.PathSearch("sell_mode", v, nil),
			"spot_options":         flattenFlavorSalesPolicySpotOptions(v),
		})
	}
	nextMarker := utils.PathSearch("sell_policies[-1].id", resp, float64(0)).(float64)
	return result, strconv.Itoa(int(nextMarker))
}

func flattenFlavorSalesPolicySpotOptions(resp interface{}) []interface{} {
	curJson := utils.PathSearch("spot_options", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"longest_spot_duration_hours": utils.PathSearch("longest_spot_duration_hours", curJson, nil),
			"largest_spot_duration_count": utils.PathSearch("largest_spot_duration_count", curJson, nil),
			"interruption_policy":         utils.PathSearch("interruption_policy", curJson, nil),
		},
	}
	return rst
}
