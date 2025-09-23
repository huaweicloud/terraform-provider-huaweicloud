package sfsturbo

import (
	"context"
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

// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/share-types
func DataSourceSfsTurboShareTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSfsTurboShareTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scenario": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attribution": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"step": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"bandwidth": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"step": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"density": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"base": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"iops": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"single_channel_4k_latency": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"support_period": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"available_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_media": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"features": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceSfsTurboShareTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/share-types"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	allShareTypes := make([]interface{}, 0)
	offset := 0

	for {
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPathWithParams := getPath + "?limit=1000" + "&offset=" + strconv.Itoa(offset)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPathWithParams, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving SFS turbo share types: %s", err)
		}

		respBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		shareTypes := utils.PathSearch("share_types", respBody, []interface{}{}).([]interface{})
		if len(shareTypes) == 0 {
			break
		}
		allShareTypes = append(allShareTypes, shareTypes...)
		offset += len(shareTypes)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("share_types", flattenShareTypes(allShareTypes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenShareTypes(shareTypes []interface{}) []map[string]interface{} {
	if len(shareTypes) == 0 {
		return nil
	}

	results := make([]map[string]interface{}, 0, len(shareTypes))
	for _, shareType := range shareTypes {
		result := map[string]interface{}{
			"share_type":      utils.PathSearch("share_type", shareType, nil),
			"scenario":        utils.PathSearch("scenario", shareType, nil),
			"attribution":     flattenAttribution(utils.PathSearch("attribution", shareType, nil)),
			"support_period":  utils.PathSearch("support_period", shareType, nil),
			"available_zones": flattenAvailableZones(utils.PathSearch("available_zones", shareType, []interface{}{}).([]interface{})),
			"spec_code":       utils.PathSearch("spec_code", shareType, nil),
			"storage_media":   utils.PathSearch("storage_media", shareType, nil),
			"features":        utils.ExpandToStringList(utils.PathSearch("features", shareType, []interface{}{}).([]interface{})),
		}
		results = append(results, result)
	}

	return results
}

func flattenAttribution(attribution interface{}) []map[string]interface{} {
	if attribution == nil {
		return nil
	}

	result := map[string]interface{}{
		"capacity":                  flattenCapacity(utils.PathSearch("capacity", attribution, nil)),
		"bandwidth":                 flattenBandwidth(utils.PathSearch("bandwidth", attribution, nil)),
		"iops":                      flattenIops(utils.PathSearch("iops", attribution, nil)),
		"single_channel_4k_latency": flattenSingleChannel4kLatency(utils.PathSearch("single_channel_4k_latency", attribution, nil)),
	}

	return []map[string]interface{}{result}
}

func flattenCapacity(capacity interface{}) []map[string]interface{} {
	if capacity == nil {
		return nil
	}

	result := map[string]interface{}{
		"max":  utils.PathSearch("max", capacity, nil),
		"min":  utils.PathSearch("min", capacity, nil),
		"step": utils.PathSearch("step", capacity, nil),
	}

	return []map[string]interface{}{result}
}

func flattenBandwidth(bandwidth interface{}) []map[string]interface{} {
	if bandwidth == nil {
		return nil
	}

	result := map[string]interface{}{
		"max":     utils.PathSearch("max", bandwidth, nil),
		"min":     utils.PathSearch("min", bandwidth, nil),
		"step":    utils.PathSearch("step", bandwidth, nil),
		"density": utils.PathSearch("density", bandwidth, nil),
		"base":    utils.PathSearch("base", bandwidth, nil),
	}

	return []map[string]interface{}{result}
}

func flattenIops(iops interface{}) []map[string]interface{} {
	if iops == nil {
		return nil
	}

	result := map[string]interface{}{
		"max": utils.PathSearch("max", iops, nil),
		"min": utils.PathSearch("min", iops, nil),
	}

	return []map[string]interface{}{result}
}

func flattenSingleChannel4kLatency(latency interface{}) []map[string]interface{} {
	if latency == nil {
		return nil
	}

	result := map[string]interface{}{
		"max": utils.PathSearch("max", latency, nil),
		"min": utils.PathSearch("min", latency, nil),
	}

	return []map[string]interface{}{result}
}

func flattenAvailableZones(availableZones []interface{}) []map[string]interface{} {
	if len(availableZones) == 0 {
		return nil
	}

	results := make([]map[string]interface{}, 0, len(availableZones))
	for _, zone := range availableZones {
		result := map[string]interface{}{
			"available_zone": utils.PathSearch("available_zone", zone, nil),
			"status":         utils.PathSearch("status", zone, nil),
		}
		results = append(results, result)
	}

	return results
}
