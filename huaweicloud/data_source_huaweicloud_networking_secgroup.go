package huaweicloud

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
			"rules": sgRuleComputedSchema,

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "tenant_id is deprecated",
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

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", secGroup.Name),
		d.Set("description", secGroup.Description),
		d.Set("rules", flattenSecurityGroupDataRules(&secGroup)),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func flattenSecurityGroupDataRules(secGroup *groups.SecGroup) []map[string]interface{} {
	sgRules := make([]map[string]interface{}, len(secGroup.Rules))
	for i, rule := range secGroup.Rules {
		sgRules[i] = map[string]interface{}{
			"id":               rule.ID,
			"direction":        rule.Direction,
			"protocol":         rule.Protocol,
			"ethertype":        rule.EtherType,
			"port_range_max":   rule.PortRangeMax,
			"port_range_min":   rule.PortRangeMin,
			"remote_ip_prefix": rule.RemoteIPPrefix,
			"remote_group_id":  rule.RemoteGroupID,
			"description":      rule.Description,
		}
	}

	return sgRules
}
