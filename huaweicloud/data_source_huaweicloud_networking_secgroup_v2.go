package huaweicloud

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingSecGroupV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"security_group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_max": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_min": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkingSecGroupV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	listOpts := groups.ListOpts{
		ID:       d.Get("secgroup_id").(string),
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allSecGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve security groups: %s", err)
	}

	if len(allSecGroups) < 1 {
		return fmtp.Errorf("No Security Group found with name: %s", d.Get("name"))
	}

	if len(allSecGroups) > 1 {
		return fmtp.Errorf("More than one Security Group found with name: %s", d.Get("name"))
	}

	secGroup := allSecGroups[0]

	logp.Printf("[DEBUG] Retrieved Security Group %s: %+v", secGroup.ID, secGroup)
	d.SetId(secGroup.ID)

	d.Set("name", secGroup.Name)
	d.Set("description", secGroup.Description)
	d.Set("region", GetRegion(d, config))
	security_group_rules := make([]map[string]string, 0, len(secGroup.Rules))
	for _, v := range secGroup.Rules {
		logp.Printf("[DEBUG] Retrieved Security Group %s", v)
		rule := make(map[string]string)
		rule["id"] = v.ID
		rule["security_group_id"] = v.SecGroupID
		rule["direction"] = v.Direction
		rule["protocol"] = v.Protocol
		rule["description"] = v.Description
		rule["ethertype"] = v.EtherType
		rule["port_range_max"] = fmt.Sprintf("%d", v.PortRangeMax)
		rule["port_range_min"] = fmt.Sprintf("%d", v.PortRangeMin)
		rule["remote_group_id"] = v.RemoteGroupID
		rule["remote_ip_prefix"] = v.RemoteIPPrefix
		security_group_rules = append(security_group_rules, rule)
	}
	d.Set("security_group_rules", security_group_rules)

	return nil
}
