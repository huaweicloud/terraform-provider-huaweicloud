package huaweicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/flavors"
)

func dataSourceIecFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIecFlavorsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"area": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"province": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
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
						"name": {
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
					},
				},
			},
		},
	}
}

func dataSourceIecFlavorsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Name:     d.Get("name").(string),
		SiteIDS:  d.Get("site_ids").(string),
		Area:     d.Get("area").(string),
		Province: d.Get("province").(string),
		City:     d.Get("city").(string),
		Operator: d.Get("operator").(string),
	}
	allFlavors, err := flavors.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmt.Errorf("Unable to extract iec flavors: %s", err)
	}
	total := len(allFlavors.Flavors)
	if total < 1 {
		return fmt.Errorf("Your query returned no results of huaweicloud_iec_flavors. " +
			"Please change your search criteria and try again.")
	}

	log.Printf("[INFO] Retrieved [%d] IEC flavors using given filter", total)
	iecFlavors := make([]map[string]interface{}, 0, total)
	for _, item := range allFlavors.Flavors {
		val := map[string]interface{}{
			"id":     item.ID,
			"name":   item.Name,
			"memory": item.Ram,
		}
		if vcpus, err := strconv.Atoi(item.Vcpus); err == nil {
			val["vcpus"] = vcpus
		}
		iecFlavors = append(iecFlavors, val)
	}
	if err := d.Set("flavors", iecFlavors); err != nil {
		return fmt.Errorf("Error saving IEC flavors: %s", err)
	}

	d.SetId(allFlavors.Flavors[0].ID)
	return nil
}
