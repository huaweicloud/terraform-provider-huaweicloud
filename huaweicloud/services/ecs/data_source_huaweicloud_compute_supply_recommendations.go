package ecs

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS POST /v1/{domain_id}/recommendations/ecs-supply
func DataSourceEcsSupplyRecommendations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEcsSupplyRecommendationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_constraint": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     supplyRecommendationFlavorConstraintSchema(),
			},
			"flavor_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"locations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     supplyRecommendationLocationsSchema(),
			},
			"option": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     supplyRecommendationOptionSchema(),
			},
			"strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"supply_recommendations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     supplyRecommendationSupplyRecommendationsSchema(),
			},
		},
	}
}

func supplyRecommendationFlavorConstraintSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"architecture_type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavor_requirements": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     supplyRecommendationFlavorConstraintFlavorRequirementsSchema(),
			},
		},
	}
}

func supplyRecommendationFlavorConstraintFlavorRequirementsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vcpu_count": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     supplyRecommendationFlavorConstraintFlavorRequirementsVcpuCountSchema(),
			},
			"memory_mb": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     supplyRecommendationFlavorConstraintFlavorRequirementsMemoryMbSchema(),
			},
			"cpu_manufacturers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"memory_gb_per_vcpu": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     supplyRecommendationFlavorConstraintFlavorRequirementsMemoryGbPerVcpuSchema(),
			},
			"instance_generations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func supplyRecommendationFlavorConstraintFlavorRequirementsVcpuCountSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func supplyRecommendationFlavorConstraintFlavorRequirementsMemoryMbSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func supplyRecommendationFlavorConstraintFlavorRequirementsMemoryGbPerVcpuSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"min": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func supplyRecommendationLocationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func supplyRecommendationOptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"result_granularity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_spot": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
		},
	}
}

func supplyRecommendationSupplyRecommendationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"score": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceEcsSupplyRecommendationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{domain_id}/recommendations/ecs-supply"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetSupplyRecommendationsBodyParams(d))

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ECS supply recommendations: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("supply_recommendations", flattenSupplyRecommendationsResponseBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSupplyRecommendationsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_constraint": buildGetSupplyRecommendationsFlavorConstraint(d.Get("flavor_constraint")),
		"flavor_ids":        utils.ValueIgnoreEmpty(d.Get("flavor_ids")),
		"locations":         buildGetSupplyRecommendationsLocations(d.Get("locations")),
		"option":            buildGetSupplyRecommendationsOption(d.Get("option")),
		"strategy":          utils.ValueIgnoreEmpty(d.Get("strategy")),
	}
	return bodyParams
}

func buildGetSupplyRecommendationsFlavorConstraint(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]interface{}{
		"architecture_type":   utils.ValueIgnoreEmpty(raw["architecture_type"]),
		"flavor_requirements": buildGetSupplyRecommendationsFlavorConstraintFlavorRequirements(raw["flavor_requirements"]),
	}
	return params
}

func buildGetSupplyRecommendationsFlavorConstraintFlavorRequirements(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"vcpu_count":           buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsVcpuCount(raw["vcpu_count"]),
				"memory_mb":            buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsMemoryMb(raw["memory_mb"]),
				"cpu_manufacturers":    utils.ValueIgnoreEmpty(raw["cpu_manufacturers"]),
				"memory_gb_per_vcpu":   buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsMemoryGbPerVcpu(raw["memory_gb_per_vcpu"]),
				"instance_generations": utils.ValueIgnoreEmpty(raw["instance_generations"]),
			})
		}
	}
	return rst
}

func buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsVcpuCount(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]interface{}{
		"max": utils.ValueIgnoreEmpty(raw["max"]),
		"min": utils.ValueIgnoreEmpty(raw["min"]),
	}
	return params
}

func buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsMemoryMb(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]interface{}{
		"max": utils.ValueIgnoreEmpty(raw["max"]),
		"min": utils.ValueIgnoreEmpty(raw["min"]),
	}
	return params
}

func buildGetSupplyRecommendationsFlavorConstraintFlavorRequirementsMemoryGbPerVcpu(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]interface{}{
		"max": utils.ValueIgnoreEmpty(raw["max"]),
		"min": utils.ValueIgnoreEmpty(raw["min"]),
	}
	return params
}

func buildGetSupplyRecommendationsLocations(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"region_id":            raw["region_id"],
				"availability_zone_id": utils.ValueIgnoreEmpty(raw["availability_zone_id"]),
			})
		}
	}
	return rst
}

func buildGetSupplyRecommendationsOption(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	params := map[string]interface{}{
		"result_granularity": raw["result_granularity"],
	}
	if v, enableSpotOk := raw["enable_spot"]; enableSpotOk {
		enableSpot, _ := strconv.ParseBool(v.(string))
		params["enable_spot"] = enableSpot
	}
	return params
}

func flattenSupplyRecommendationsResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("supply_recommendations", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id":            utils.PathSearch("flavor_id", v, nil),
			"region_id":            utils.PathSearch("region_id", v, nil),
			"availability_zone_id": utils.PathSearch("availability_zone_id", v, nil),
			"score":                utils.PathSearch("score", v, nil),
		})
	}
	return rst
}
