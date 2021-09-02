package huaweicloud

import (
	"strconv"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/security/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIECSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remote_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIECSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	allSecGroups, err := groups.List(iecClient, nil).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve security groups: %s", err)
	}

	total := len(allSecGroups.SecurityGroups)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results")
	}

	// filter security groups by name
	var groupItem *groups.RespSecurityGroupEntity
	name := d.Get("name").(string)
	for _, group := range allSecGroups.SecurityGroups {
		if group.Name == name {
			groupItem = &group
			break
		}
	}
	if groupItem == nil {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	logp.Printf("[DEBUG] Retrieved IEC security group %s: %+v", groupItem.ID, groupItem)
	d.SetId(groupItem.ID)
	d.Set("name", groupItem.Name)
	d.Set("description", groupItem.Description)

	secRules := make([]map[string]interface{}, len(groupItem.SecurityGroupRules))
	for index, rule := range groupItem.SecurityGroupRules {
		secRules[index] = map[string]interface{}{
			"id":                rule.ID,
			"security_group_id": rule.SecurityGroupID,
			"description":       rule.Description,
			"direction":         rule.Direction,
			"ethertype":         rule.EtherType,
			"protocol":          rule.Protocol,
			"remote_group_id":   rule.RemoteGroupID,
			"remote_ip_prefix":  rule.RemoteIPPrefix,
		}
		if ret, err := strconv.Atoi(rule.PortRangeMax.(string)); err == nil {
			secRules[index]["port_range_max"] = ret
		}
		if ret, err := strconv.Atoi(rule.PortRangeMin.(string)); err == nil {
			secRules[index]["port_range_min"] = ret
		}
	}
	d.Set("security_group_rules", secRules)

	return nil
}
