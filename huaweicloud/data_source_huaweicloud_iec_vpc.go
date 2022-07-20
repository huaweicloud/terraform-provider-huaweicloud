package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceIECVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECVpcRead,

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
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceIECVpcRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	logp.Printf("[DEBUG] query VPCs using given filter: %+v", listOpts)
	allVpcs, err := vpcs.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve vpcs: %s", err)
	}

	total := len(allVpcs.Vpcs)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if total > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	vpcInfo := allVpcs.Vpcs[0]
	logp.Printf("[DEBUG] Retrieved IEC VPC %s: %+v", vpcInfo.ID, vpcInfo)

	d.SetId(vpcInfo.ID)
	d.Set("name", vpcInfo.Name)
	d.Set("cidr", vpcInfo.Cidr)
	d.Set("mode", vpcInfo.Mode)
	d.Set("subnet_num", vpcInfo.SubnetNum)
	d.Set("region", GetRegion(d, config))

	return nil
}
