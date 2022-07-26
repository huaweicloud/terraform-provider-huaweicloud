package vpc

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceVpcSubnetV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcSubnetV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_list": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	subnetClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc client: %s", err)
	}

	listOpts := subnets.ListOpts{
		ID:               d.Get("id").(string),
		Name:             d.Get("name").(string),
		CIDR:             d.Get("cidr").(string),
		Status:           d.Get("status").(string),
		GatewayIP:        d.Get("gateway_ip").(string),
		PRIMARY_DNS:      d.Get("primary_dns").(string),
		SECONDARY_DNS:    d.Get("secondary_dns").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		VPC_ID:           d.Get("vpc_id").(string),
	}

	refinedSubnets, err := subnets.List(subnetClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve subnets: %s", err)
	}

	if len(refinedSubnets) == 0 {
		return fmtp.Errorf("No matching subnet found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedSubnets) > 1 {
		return fmtp.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	Subnets := refinedSubnets[0]

	logp.Printf("[INFO] Retrieved Subnet using given filter %s: %+v", Subnets.ID, Subnets)
	d.SetId(Subnets.ID)

	d.Set("name", Subnets.Name)
	d.Set("description", Subnets.Description)
	d.Set("cidr", Subnets.CIDR)
	d.Set("dns_list", Subnets.DnsList)
	d.Set("status", Subnets.Status)
	d.Set("gateway_ip", Subnets.GatewayIP)
	d.Set("ipv6_enable", Subnets.EnableIPv6)
	d.Set("dhcp_enable", Subnets.EnableDHCP)
	d.Set("primary_dns", Subnets.PRIMARY_DNS)
	d.Set("secondary_dns", Subnets.SECONDARY_DNS)
	d.Set("availability_zone", Subnets.AvailabilityZone)
	d.Set("vpc_id", Subnets.VPC_ID)
	d.Set("subnet_id", Subnets.SubnetId)
	d.Set("ipv6_subnet_id", Subnets.IPv6SubnetId)
	d.Set("ipv6_cidr", Subnets.IPv6CIDR)
	d.Set("ipv6_gateway", Subnets.IPv6Gateway)
	d.Set("region", config.GetRegion(d))

	return nil
}
