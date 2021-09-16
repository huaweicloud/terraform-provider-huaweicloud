package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceVpcSubnetIdsV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcSubnetIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVpcSubnetIdsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := subnets.ListOpts{
		VPC_ID: d.Get("vpc_id").(string),
	}

	refinedSubnets, err := subnets.List(subnetClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve subnets: %s", err)
	}

	if len(refinedSubnets) == 0 {
		return fmtp.Errorf("no matching subnet found for vpc with id %s", d.Get("vpc_id").(string))
	}

	Subnets := make([]string, 0)

	for _, subnet := range refinedSubnets {
		Subnets = append(Subnets, subnet.ID)
	}

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", Subnets)

	d.Set("region", config.GetRegion(d))

	return nil
}
