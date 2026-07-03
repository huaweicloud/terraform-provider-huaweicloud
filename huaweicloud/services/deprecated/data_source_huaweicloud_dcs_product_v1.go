package deprecated

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/products"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceDcsProductV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDcsProductV1Read,
		DeprecationMessage: "this is deprecated." +
			"This data source is used for the \"product_id\" of the \"huaweicloud_dcs_instance\" resource. " +
			"Now \"product_id\" has been deprecated and this data source is no longer used.",

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

			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cache_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsProductV1Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dcsV1Client, err := cfg.DcsV1Client(region)
	if err != nil {
		return fmt.Errorf("error getting DCS client: %s", err)
	}

	v, err := products.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	specCode := d.Get("spec_code").(string)
	log.Printf("[DEBUG] query DCS products with %s", specCode)

	var filteredPd *products.Product
	for _, pd := range v.Products {
		if specCode != "" && pd.SpecCode != specCode {
			continue
		}
		filteredPd = &pd
		break
	}

	if filteredPd == nil {
		return errors.New("your query returned no results, please change your filters and try again")
	}

	log.Printf("[DEBUG] get DCS product: %+v", filteredPd)
	d.SetId(filteredPd.ProductID)
	d.Set("spec_code", filteredPd.SpecCode)
	d.Set("engine", filteredPd.Engine)
	d.Set("engine_version", filteredPd.EngineVersion)
	d.Set("cache_mode", filteredPd.CacheMode)

	return nil
}
