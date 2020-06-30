package huaweicloud

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/networkipavailabilities"
)

func dataSourceVpcSubnetIdsV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcSubnetIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceVpcSubnetIdsV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

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

	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	sortedSubnets := make([]SubnetIP, 0)
	for _, subnet := range refinedSubnets {
		net, err := networkipavailabilities.Get(networkingClient, subnet.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error retrieving network ip availabilities: %s", err)
		}
		subnetIPAvail := net.SubnetIPAvailabilities[0]
		newSubnet := SubnetIP{
			ID:  subnet.ID,
			IPs: subnetIPAvail.TotalIPs - subnetIPAvail.UsedIPs,
		}
		sortedSubnets = append(sortedSubnets, newSubnet)
	}

	// Returns the Subnet contains most available IPs out of a slice of subnets.
	sort.Sort(sort.Reverse(subnetSort(sortedSubnets)))
	Subnets := make([]string, 0)
	for _, subnet := range sortedSubnets {
		Subnets = append(Subnets, subnet.ID)
	}

	d.SetId(d.Get("vpc_id").(string))
	d.Set("ids", Subnets)

	d.Set("region", GetRegion(d, config))

	return nil
}

type SubnetIP struct {
	ID  string
	IPs int
}

type subnetSort []SubnetIP

func (a subnetSort) Len() int      { return len(a) }
func (a subnetSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a subnetSort) Less(i, j int) bool {
	return a[i].IPs < a[j].IPs
}
