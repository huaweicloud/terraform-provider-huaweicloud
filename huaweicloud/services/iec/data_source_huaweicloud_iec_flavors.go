package iec

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/flavors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/cloudservers/flavors
func DataSourceFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFlavorsRead,

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

func dataSourceFlavorsRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)

	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
	}

	listOpts := flavors.ListOpts{
		Name:     d.Get("name").(string),
		SiteIDS:  d.Get("site_ids").(string),
		Area:     d.Get("area").(string),
		Province: d.Get("province").(string),
		City:     d.Get("city").(string),
		Operator: d.Get("operator").(string),
	}

	log.Printf("[DEBUG] fetching IEC flavors by filter: %#v", listOpts)
	allFlavors, err := flavors.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmt.Errorf("unable to extract IEC flavors: %s", err)
	}
	total := len(allFlavors.Flavors)
	if total < 1 {
		return fmt.Errorf("your query returned no results of iec_flavors, " +
			"please change your search criteria and try again")
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
	mErr := multierror.Append(d.Set("flavors", iecFlavors))
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error saving IEC flavors: %s", err)
	}

	d.SetId(allFlavors.Flavors[0].ID)
	return nil
}
