package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/availablezones"
)

func dataSourceDmsAZV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDmsAZV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDmsAZV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.dmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	v, err := availablezones.Get(dmsV1Client).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dms az : %+v", v)
	var filteredAZs []availablezones.AvailableZone
	if v.RegionID == GetRegion(d, config) {
		AZs := v.AvailableZones
		for _, newAZ := range AZs {
			if newAZ.ResourceAvailability != "true" {
				continue
			}

			name := d.Get("name").(string)
			if name != "" && newAZ.Name != name {
				continue
			}

			port := d.Get("port").(string)
			if port != "" && newAZ.Port != port {
				continue
			}
			filteredAZs = append(filteredAZs, newAZ)
		}
	}

	if len(filteredAZs) < 1 {
		return fmt.Errorf("Not found any available zones")
	}

	az := filteredAZs[0]
	log.Printf("[DEBUG] Dms az : %+v", az)

	d.SetId(az.ID)
	d.Set("code", az.Code)
	d.Set("name", az.Name)
	d.Set("port", az.Port)

	return nil
}
