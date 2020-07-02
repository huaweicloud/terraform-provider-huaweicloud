package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/flavors"
)

func dataSourceDDSFlavorV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDDSFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DDS-Community", "DDS-Enhanced",
				}, true),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"mongos", "shard", "config", "replica", "single",
				}, true),
			},
			"vcpus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDDSFlavorV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ddsClient, err := config.ddsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Region:     GetRegion(d, config),
		EngineName: d.Get("engine_name").(string),
	}

	pages, err := flavors.List(ddsClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract flavors: %s", err)
	}

	flavorList := make([]map[string]interface{}, 0)
	filterType := d.Get("type").(string)
	filterVcpus := d.Get("vcpus").(string)
	filterMemory := d.Get("memory").(string)

	for _, item := range allFlavors {
		if filterFlavor(item, filterType, filterVcpus, filterMemory) {
			continue
		}

		flavor := map[string]interface{}{
			"spec_code": item.SpecCode,
			"type":      item.Type,
			"vcpus":     item.Vcpus,
			"memory":    item.Ram,
		}
		flavorList = append(flavorList, flavor)
	}

	log.Printf("Extract %d/%d flavors by filters.", len(flavorList), len(allFlavors))
	if len(flavorList) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId("dds flavors")
	d.Set("flavors", flavorList)
	d.Set("region", GetRegion(d, config))

	return nil
}

func filterFlavor(item flavors.Flavor, flavorType, vcpus, memory string) bool {
	if flavorType != "" && flavorType != item.Type {
		return true
	}
	if vcpus != "" && vcpus != item.Vcpus {
		return true
	}
	if memory != "" && memory != item.Ram {
		return true
	}

	return false
}
