package huaweicloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/flavors"
)

func DataSourceEcsFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEcsFlavorsRead,

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
			"performance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"generation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceEcsFlavorsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ecsClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud ECS client: %s", err)
	}

	listOpts := &flavors.ListOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	pages, err := flavors.List(ecsClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s ", err)
	}

	cpu := d.Get("cpu_core_count").(int)
	mem := int64(d.Get("memory_size").(int)) * 1024
	pType := d.Get("performance_type").(string)
	gen := d.Get("generation").(string)
	az := d.Get("availability_zone").(string)

	var ids []string
	for _, flavor := range allFlavors {
		vCpu, _ := strconv.Atoi(flavor.Vcpus)
		if cpu > 0 && vCpu != cpu {
			continue
		}

		if mem > 0 && flavor.Ram != mem {
			continue
		}

		if pType != "" && flavor.OsExtraSpecs.PerformanceType != pType {
			continue
		}

		if gen != "" && flavor.OsExtraSpecs.Generation != gen {
			continue
		}

		if az != "" {
			status := flavor.OsExtraSpecs.OperationStatus
			azStatusRaw := flavor.OsExtraSpecs.OperationAz
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
	}

	if len(ids) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("region", GetRegion(d, config))

	return nil
}
