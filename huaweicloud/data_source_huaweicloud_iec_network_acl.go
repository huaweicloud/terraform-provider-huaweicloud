package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIECNetworkACL() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECNetworkACLRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"id"},
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed but always be empty due to the API response
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIECNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := firewalls.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	logp.Printf("[DEBUG] query firewall using given filter: %+v", listOpts)
	allFWs, err := firewalls.List(iecClient, listOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve firewall: %s", err)
	}

	total := len(allFWs.Firewalls)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if total > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	fwGroup := allFWs.Firewalls[0]
	logp.Printf("[DEBUG] Retrieved IEC firewall %s: %+v", fwGroup.ID, fwGroup)

	d.SetId(fwGroup.ID)
	d.Set("name", fwGroup.Name)
	d.Set("status", fwGroup.Status)
	d.Set("description", fwGroup.Description)

	// currently, the following attributes are empty due to the API response
	networkList := make([]map[string]interface{}, 0, len(fwGroup.Subnets))
	for _, val := range fwGroup.Subnets {
		subnet := make(map[string]interface{})
		subnet["vpc_id"] = val.VpcID
		subnet["subnet_id"] = val.ID
		networkList = append(networkList, subnet)
	}
	d.Set("networks", networkList)
	d.Set("inbound_rules", getFirewallRuleIDs(fwGroup.IngressFWPolicy))
	d.Set("outbound_rules", getFirewallRuleIDs(fwGroup.EgressFWPolicy))

	return nil
}
