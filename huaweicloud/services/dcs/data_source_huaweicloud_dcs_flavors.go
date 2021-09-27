package dcs

import (
	"context"
	"sort"
	"strconv"

	"github.com/chnsz/golangsdk/openstack/dcs/v2/flavors"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	floatBitSize = 64
)

// DataSourceDcsFlavorsV2 the function is used for data source 'huaweicloud_dcs_flavors'.
func DataSourceDcsFlavorsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsFlavorsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Redis",
				ValidateFunc: validation.StringInSlice([]string{
					"Redis", "Memcached",
				}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"single", "ha", "cluster", "proxy", "ha_rw_split",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"x86_64", "aarch64",
				}, false),
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
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DCS client: %s", err)
	}

	capacity := strconv.FormatFloat(d.Get("capacity").(float64), 'f', -1, floatBitSize)
	// build a list options
	opts := flavors.ListOpts{
		CacheMode:     d.Get("cache_mode").(string),
		Engine:        d.Get("engine").(string),
		EngineVersion: d.Get("engine_version").(string),
		Capacity:      capacity,
		SpecCode:      d.Get("name").(string),
		CPUType:       d.Get("cpu_architecture").(string),
	}
	logp.Printf("[DEBUG] The options of list DCS flavors : %#v", opts)

	list, err := flavors.List(client, opts).Extract()
	logp.Printf("[DEBUG] Get DCS flavors : %#v", list)
	if len(list) == 0 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	ids := make([]string, 0)
	flavors := make([]map[string]interface{}, 0, len(list))
	if len(list) > 0 {
		for _, v := range list {
			if len(v.AvailableZones) == 0 || len(v.AvailableZones[0].AzCodes) == 0 {
				continue
			}
			cap, _ := strconv.ParseFloat(v.Capacity[0], floatBitSize)
			fla := map[string]interface{}{
				"name":             v.SpecCode,
				"cache_mode":       v.CacheMode,
				"engine":           v.Engine,
				"engine_versions":  v.EngineVersion,
				"cpu_architecture": v.CPUType,
				"capacity":         cap,
				"available_zones":  v.AvailableZones[0].AzCodes,
				"charging_modes":   v.BillingMode,
				"ip_count":         v.TenantIPCount,
			}
			flavors = append(flavors, fla)
			ids = append(ids, v.SpecCode)
		}
	}

	if len(flavors) == 0 {
		return fmtp.DiagErrorf("Your query did not return valid data. " +
			"Please change your search criteria and try again.")
	}

	sort.Slice(flavors, func(i, j int) bool {
		a := flavors[i]
		b := flavors[j]
		v1 := a["ip_count"].(int)
		v2 := b["ip_count"].(int)

		return v1 <= v2
	})

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("flavors", flavors),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting DCS flavors attributes: %s", mErr)
	}

	return nil
}
