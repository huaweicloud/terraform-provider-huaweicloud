package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/products"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceDcsProductV1() *schema.Resource {
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
	config := meta.(*config.Config)
	dcsV1Client, err := config.DcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error get dcs product client: %s", err)
	}

	v, err := products.Get(dcsV1Client).Extract()
	if err != nil {
		return err
	}
	logp.Printf("[DEBUG] Dcs get products : %+v", v)
	var FilteredPd []products.Product
	for _, pd := range v.Products {
		spec_code := d.Get("spec_code").(string)
		if spec_code != "" && pd.SpecCode != spec_code {
			continue
		}
		FilteredPd = append(FilteredPd, pd)
	}

	if len(FilteredPd) < 1 {
		return fmtp.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	pd := FilteredPd[0]
	d.SetId(pd.ProductID)
	d.Set("spec_code", pd.SpecCode)
	logp.Printf("[DEBUG] Dcs product : %+v", pd)

	return nil
}
