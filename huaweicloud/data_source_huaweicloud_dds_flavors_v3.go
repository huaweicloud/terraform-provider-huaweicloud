package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
				ForceNew: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
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
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
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
		Region: GetRegion(d, config),
	}

	if v, ok := d.GetOk("engine_name"); ok {
		listOpts.EngineName = v.(string)
	}

	pages, err := flavors.List(ddsClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	allFlavors, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract flavors: %s", err)
	}

	if len(allFlavors) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var refinedFlavors []flavors.Flavor
	var flavor flavors.Flavor
	if v, ok := d.GetOk("spec_code"); ok {
		for _, flavor = range allFlavors {
			if flavor.SpecCode == v.(string) {
				refinedFlavors = append(refinedFlavors, flavor)
			}
		}
		if len(refinedFlavors) < 1 {
			return fmt.Errorf("Your query returned no results. " +
				"Please change your search criteria and try again.")
		}
		flavor = refinedFlavors[0]
	} else {
		flavor = allFlavors[0]
	}

	log.Printf("[DEBUG] Retrieved DDS Flavor: %+v ", flavor)
	d.SetId(flavor.SpecCode)

	d.Set("engine_name", flavor.EngineName)
	d.Set("spec_code", flavor.SpecCode)
	d.Set("type", flavor.Type)
	d.Set("vcpus", flavor.Vcpus)
	d.Set("ram", flavor.Ram)
	d.Set("region", GetRegion(d, config))

	return nil
}
