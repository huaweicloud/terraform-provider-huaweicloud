package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceIECVpcSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECVpcSubnetIdsRead,

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
			"site_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIECVpcSubnetIdsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	vpcID := d.Get("vpc_id").(string)
	listOpts := subnets.ListOpts{
		VpcID:  vpcID,
		SiteID: d.Get("site_id").(string),
	}

	allSubnets, err := subnets.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve subnets: %s", err)
	}

	total := len(allSubnets.Subnets)
	if total == 0 {
		return fmtp.Errorf("no matching subnet found for vpc with id %s", vpcID)
	}

	iecSubnets := make([]map[string]interface{}, total)
	allIDs := make([]string, total)
	for i, item := range allSubnets.Subnets {
		val := map[string]interface{}{
			"id":         item.ID,
			"name":       item.Name,
			"cidr":       item.Cidr,
			"gateway_ip": item.GatewayIP,
			"site_id":    item.SiteID,
			"site_info":  item.SiteInfo,
			"dns_list":   item.DNSList,
			"status":     item.Status,
		}
		iecSubnets[i] = val
		allIDs[i] = item.ID
	}
	if err := d.Set("subnets", iecSubnets); err != nil {
		return fmtp.Errorf("Error saving IEC subnets: %s", err)
	}

	// set id
	d.SetId(hashcode.Strings(allIDs))
	d.Set("region", GetRegion(d, config))

	return nil
}
