package bms

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/bms/v1/flavors"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceBmsFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBmsFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_arch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "x86_64",
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_arch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBmsFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	bmsClient, err := config.BmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud BMS client: %s", err)
	}

	az := d.Get("availability_zone").(string)
	listOpts := flavors.ListOpts{
		AvailabilityZone: az,
	}

	allFlavors, err := flavors.List(bmsClient, listOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve BMS flavors: %s ", err)
	}

	var vcpus string
	if v, ok := d.GetOk("vcpus"); ok {
		vcpus = strconv.Itoa(v.(int))
	}
	mem := d.Get("memory").(int) * 1024

	filter := map[string]interface{}{
		"VCPUs":                vcpus,
		"RAM":                  mem,
		"OsExtraSpecs.CPUArch": d.Get("cpu_arch"),
	}

	filterFlavors, err := utils.FilterSliceWithField(allFlavors, filter)
	if err != nil {
		return fmtp.DiagErrorf("filter BMS flavors failed: %s", err)
	}
	logp.Printf("filter %d bms flavors from %d through options %v", len(filterFlavors), len(allFlavors), filter)

	var ids []string
	var resultFlavors []map[string]interface{}

	for _, item := range filterFlavors {
		flavor := item.(flavors.Flavor)

		// ignore abandon and sellout flavors
		if az != "" {
			status := flavor.OsExtraSpecs.OperationStatus
			azStatusRaw := flavor.OsExtraSpecs.OperationAZ
			azStatusList := strings.Split(azStatusRaw, ",")
			if strings.Contains(azStatusRaw, az) {
				for i := 0; i < len(azStatusList); i++ {
					azStatus := azStatusList[i]
					if azStatus == (az+"(abandon)") || azStatus == (az+"(sellout)") || azStatus == (az+"obt_sellout") {
						continue
					}
				}
			} else if status == "abandon" || strings.Contains(status, "sellout") {
				continue
			}
		}

		ids = append(ids, flavor.ID)
		resultFlavors = append(resultFlavors, flattenBmsFlavor(flavor))
	}

	if len(resultFlavors) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(hashcode.Strings(ids))
	d.Set("region", region)
	d.Set("flavors", resultFlavors)

	return nil
}

func flattenBmsFlavor(flavor flavors.Flavor) map[string]interface{} {
	vcpus, _ := strconv.Atoi(flavor.VCPUs)
	ram := int(flavor.RAM / 1024)

	return map[string]interface{}{
		"id":        flavor.ID,
		"vcpus":     vcpus,
		"memory":    ram,
		"cpu_arch":  flavor.OsExtraSpecs.CPUArch,
		"operation": flavor.OsExtraSpecs.OperationAZ,
	}
}
