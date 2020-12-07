package huaweicloud

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/availabilityzones"
)

func dataSourceComputeAvailabilityZonesV2() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceComputeAvailabilityZonesV2Read,
		DeprecationMessage: "use huaweicloud_availability_zones data source instead",
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:         schema.TypeString,
				Default:      "available",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "unavailable"}, true),
			},
		},
	}
}

func dataSourceComputeAvailabilityZonesV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	allPages, err := availabilityzones.List(computeClient).AllPages()
	if err != nil {
		return fmt.Errorf("Error retrieving huaweicloud_compute_availability_zones_v2: %s", err)
	}
	zoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting huaweicloud_compute_availability_zones_v2 from response: %s", err)
	}

	stateBool := d.Get("state").(string) == "available"
	zones := make([]string, 0, len(zoneInfo))
	for _, z := range zoneInfo {
		if z.ZoneState.Available == stateBool {
			zones = append(zones, z.ZoneName)
		}
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(zones)

	d.SetId(hashcode.Strings(zones))
	d.Set("names", zones)

	return nil
}
