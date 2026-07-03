package deprecated

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/availablezones"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceDcsAZV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsAZV1Read,
		DeprecationMessage: "this is deprecated. " +
			"This data source is used for the \"available_zones\" of the \"huaweicloud_dcs_instance\" resource. " +
			"Now `available_zones` has been deprecated and this data source is no longer used.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
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

func dataSourceDcsAZV1Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dcsV1Client, err := cfg.DcsV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating DCS client: %s", err)
	}

	v, err := availablezones.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] fetching DCS available zones : %+v", v)
	var filteredAZs []availablezones.AvailableZone
	if v.RegionID == region {
		AZs := v.AvailableZones
		for _, newAZ := range AZs {
			if newAZ.ResourceAvailability != "true" {
				continue
			}

			code := d.Get("code").(string)
			if code != "" && newAZ.Code != code {
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
		return errors.New("not found any available zones")
	}

	az := filteredAZs[0]
	log.Printf("[DEBUG] filter DCS available zone: %+v", az)

	d.SetId(az.ID)
	d.Set("code", az.Code)
	d.Set("name", az.Name)
	d.Set("port", az.Port)

	return nil
}
