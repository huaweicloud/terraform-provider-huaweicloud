package huaweicloud

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVirtualPrivateCloudVpcIdsV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualPrivateCloudIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceVirtualPrivateCloudIdsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := vpcs.ListOpts{}
	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	listVpcs := make([]string, 0)

	for _, vpc := range refinedVpcs {
		listVpcs = append(listVpcs, vpc.ID)
	}
	d.SetId(listVpcs[0])
	d.Set("ids", listVpcs)

	return nil
}
