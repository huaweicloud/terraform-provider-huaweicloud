package huaweicloud

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceVpcSubnetIdsV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcSubnetIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVpcSubnetIdsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))

	listOpts := subnets.ListOpts{
		VPC_ID: d.Get("vpc_id").(string),
	}

	refinedSubnets, err := subnets.List(subnetClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve subnets: %s", err)
	}

	if len(refinedSubnets) == 0 {
		return fmt.Errorf("no matching subnet found for vpc with id %s", d.Get("vpc_id").(string))
	}

	Subnets := make([]string, 0)

	for _, subnet := range refinedSubnets {
		Subnets = append(Subnets, subnet.ID)
	}

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", Subnets)

	d.Set("region", GetRegion(d, config))

	return nil
}
