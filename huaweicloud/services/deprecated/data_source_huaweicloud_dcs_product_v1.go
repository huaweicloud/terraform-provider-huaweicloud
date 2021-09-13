package deprecated

import (
	"github.com/chnsz/golangsdk/openstack/dcs/v1/products"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	dcsV1Client, err := config.DcsV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error get dcs product client: %s", err)
	}

	v, err := products.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}

	specCode := d.Get("spec_code").(string)
	logp.Printf("[DEBUG] query DCS products with %s", specCode)

	var filteredPd *products.Product
	for _, pd := range v.Products {
		if specCode != "" && pd.SpecCode != specCode {
			continue
		}
		filteredPd = &pd
		break
	}

	if filteredPd == nil {
		return fmtp.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	logp.Printf("[DEBUG] get DCS product: %+v", filteredPd)
	d.SetId(filteredPd.ProductID)
	d.Set("spec_code", filteredPd.SpecCode)
	d.Set("engine", filteredPd.Engine)
	d.Set("engine_version", filteredPd.EngineVersion)
	d.Set("cache_mode", filteredPd.CacheMode)

	return nil
}
