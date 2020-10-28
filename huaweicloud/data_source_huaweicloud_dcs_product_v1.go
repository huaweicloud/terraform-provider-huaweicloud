package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/products"
)

func dataSourceDcsProductV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsProductV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsProductV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error get dcs product client: %s", err)
	}

	v, err := products.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Dcs get products : %+v", v)
	var FilteredPd []products.Product
	for _, pd := range v.Products {
		spec_code := d.Get("spec_code").(string)
		if spec_code != "" && pd.SpecCode != spec_code {
			continue
		}
		FilteredPd = append(FilteredPd, pd)
	}

	if len(FilteredPd) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	pd := FilteredPd[0]
	d.SetId(pd.ProductID)
	d.Set("spec_code", pd.SpecCode)
	log.Printf("[DEBUG] Dcs product : %+v", pd)

	return nil
}
