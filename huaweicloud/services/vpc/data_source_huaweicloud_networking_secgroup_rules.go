package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v3rules "github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPC GET /v3/{project_id}/vpc/security-group-rules
func DataSourceNetworkingSecGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"ethertype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
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
						"remote_address_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkingSecGroupRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	listOpts := v3rules.ListOpts{
		SecurityGroupId: d.Get("security_group_id").(string),
		ID:              d.Get("rule_id").(string),
		Protocol:        d.Get("protocol").(string),
		Description:     d.Get("description").(string),
		RemoteGroupId:   d.Get("remote_group_id").(string),
		Direction:       d.Get("direction").(string),
		Action:          d.Get("action").(string),
	}
	resp, err := v3rules.List(client, listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Security Group rules")
	}
	rules := make([]map[string]interface{}, len(resp))
	for i, rule := range resp {
		ruleInfo := map[string]interface{}{
			"id":                      rule.ID,
			"description":             rule.Description,
			"security_group_id":       rule.SecurityGroupId,
			"direction":               rule.Direction,
			"protocol":                rule.Protocol,
			"ethertype":               rule.Ethertype,
			"ports":                   rule.MultiPort,
			"action":                  rule.Action,
			"priority":                rule.Priority,
			"remote_group_id":         rule.RemoteGroupId,
			"remote_ip_prefix":        rule.RemoteIpPrefix,
			"remote_address_group_id": rule.RemoteAddressGroupId,
			"created_at":              rule.CreateAt,
			"updated_at":              rule.UpdateAt,
		}
		rules[i] = ruleInfo
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", rules),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving Security Group rules data source fields: %s", mErr)
	}
	return nil
}
