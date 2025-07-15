package dcs

import (
	"context"
	"fmt"
	"sort"
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

const (
	floatBitSize = 64
)

// DataSourceDcsFlavorsV2 the function is used for data source 'huaweicloud_dcs_flavors'.
// @API DCS GET /v2/{project_id}/flavors
func DataSourceDcsFlavorsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsFlavorsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Redis",
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_architecture": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_versions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"available_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"charging_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDcsFlavorsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/flavors"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getFlavorsQueryParams := buildGetDcsFlavorsDetailQueryParams(d)
	getPath += getFlavorsQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS flavors: %s", err)
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

	flavorLists := flattenDcsFlavors(d, getRespBody)
	sort.Slice(flavorLists, func(i, j int) bool {
		a := flavorLists[i]
		b := flavorLists[j]
		v1 := a["ip_count"].(float64)
		v2 := b["ip_count"].(float64)

		return v1 <= v2
	})

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavors", flavorLists),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDcsFlavors(d *schema.ResourceData, resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	var rawCapacity string
	if c, ok := d.GetOk("capacity"); ok {
		rawCapacity = strconv.FormatFloat(c.(float64), 'f', -1, floatBitSize)
	}
	curJson := utils.PathSearch("flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		availableZones := utils.PathSearch("flavors_available_zones", v, make([]interface{}, 0)).([]interface{})
		if len(availableZones) == 0 {
			continue
		}
		azCodes := utils.PathSearch("[0].az_codes", availableZones, make([]interface{}, 0)).([]interface{})
		if len(azCodes) == 0 {
			continue
		}
		for _, availableZone := range availableZones {
			capacity := utils.PathSearch("capacity", availableZone, "").(string)
			if rawCapacity != "" && rawCapacity != capacity {
				continue
			}
			capacityFloat, _ := strconv.ParseFloat(capacity, floatBitSize)
			rst = append(rst, map[string]interface{}{
				"name":             utils.PathSearch("spec_code", v, nil),
				"cache_mode":       utils.PathSearch("cache_mode", v, nil),
				"engine":           utils.PathSearch("engine", v, nil),
				"engine_versions":  utils.PathSearch("engine_version", v, nil),
				"cpu_architecture": utils.PathSearch("cpu_type", v, nil),
				"capacity":         capacityFloat,
				"charging_modes":   utils.PathSearch("billing_mode", v, nil),
				"ip_count":         utils.PathSearch("tenant_ip_count", v, nil),
			})
		}
	}
	return rst
}

func buildGetDcsFlavorsDetailQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&spec_code=%v", res, v)
	}
	if v, ok := d.GetOk("cache_mode"); ok {
		res = fmt.Sprintf("%s&cache_mode=%v", res, v)
	}
	if v, ok := d.GetOk("engine"); ok {
		res = fmt.Sprintf("%s&engine=%v", res, v)
	}
	if v, ok := d.GetOk("engine_version"); ok {
		res = fmt.Sprintf("%s&engine_version=%v", res, v)
	}
	if v, ok := d.GetOk("cpu_architecture"); ok {
		res = fmt.Sprintf("%s&cpu_type=%v", res, v)
	}
	if v, ok := d.GetOk("capacity"); ok {
		res = fmt.Sprintf("%s&capacity=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
