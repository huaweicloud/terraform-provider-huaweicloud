package huaweicloud

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupRead,

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
			"rules": securityGroupRuleSchema,
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingSecGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking client: %s", err)
	}

	listOpts := groups.ListOpts{
		ID:   d.Get("secgroup_id").(string),
		Name: d.Get("name").(string),
	}

	allSecGroups, err := groups.List(networkingClient, listOpts)
	if err != nil {
		return fmtp.DiagErrorf("Unable to get security groups: %s", err)
	}

	if len(allSecGroups) < 1 {
		return fmtp.DiagErrorf("No Security Group found with name: %s", d.Get("name"))
	}

	if len(allSecGroups) > 1 {
		return fmtp.DiagErrorf("More than one Security Group found with name: %s", d.Get("name"))
	}

	secGroup := allSecGroups[0]

	logp.Printf("[DEBUG] Retrieved Security Group %s: %+v", secGroup.ID, secGroup)
	d.SetId(secGroup.ID)

	secGroupRule, err := flattenSecurityGroupRules(&secGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", secGroup.Name),
		d.Set("description", secGroup.Description),
		d.Set("rules", secGroupRule),
		d.Set("created_at", secGroup.CreatedAt),
		d.Set("updated_at", secGroup.UpdatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
