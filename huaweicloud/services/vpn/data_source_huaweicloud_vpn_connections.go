package vpn

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpn/v5/connections"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/vpn-connection
func DataSourceVpnConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceConnectionRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
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
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tunnel_local_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tunnel_peer_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_nqa": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ikepolicy": {
							Type:     schema.TypeList,
							Elem:     connectionIkePolicySchema(),
							Computed: true,
						},
						"ipsecpolicy": {
							Type:     schema.TypeList,
							Elem:     connectionIpsecPolicySchema(),
							Computed: true,
						},
						"policy_rules": {
							Type:     schema.TypeList,
							Elem:     connectionsPolicyRuleSchema(),
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
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_monitor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func connectionIkePolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ike_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phase1_negotiation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dh_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"local_id_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_id_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dpd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     connectionIkePolicyDPDSchema(),
			},
		},
	}
	return &sc
}

func connectionIkePolicyDPDSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func connectionIpsecPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pfs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifetime_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transform_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encapsulation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func connectionsPolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.VpnV5Client(region)
	if err != nil {
		return diag.Errorf("error creating vpn v5 client: %s", err)
	}
	opts := connections.ListOpts{}
	allConnections, err := connections.List(client, opts)
	if err != nil {
		return diag.Errorf("unable to list connections: %s ", err)
	}

	log.Printf("[DEBUG] retrieved VPN connections: %#v", allConnections)
	filter := map[string]interface{}{
		"ID":     d.Get("connection_id"),
		"Name":   d.Get("name").(string),
		"Status": d.Get("status").(string),
		"Style":  strings.ToUpper(d.Get("vpn_type").(string)),
		"VgwId":  d.Get("gateway_id").(string),
		"VgwIp":  d.Get("gateway_ip").(string),
	}

	filterConnections, err := utils.FilterSliceWithField(allConnections, filter)
	if err != nil {
		return diag.Errorf("filter connections failed: %s", err)
	}

	var cnt []map[string]interface{}
	for _, item := range filterConnections {
		connection := item.(connections.Connections)
		cnt = append(cnt, flattenSourceConnection(connection))
	}

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuidStr)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("connections", cnt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSourceConnection(cnt connections.Connections) map[string]interface{} {
	resourceConnections := map[string]interface{}{
		"id":                    cnt.ID,
		"name":                  cnt.Name,
		"status":                cnt.Status,
		"gateway_id":            cnt.VgwId,
		"gateway_ip":            cnt.VgwIp,
		"vpn_type":              strings.ToLower(cnt.Style),
		"customer_gateway_id":   cnt.CgwId,
		"peer_subnets":          cnt.PeerSubnets,
		"tunnel_local_address":  cnt.TunnelLocalAddress,
		"tunnel_peer_address":   cnt.TunnelPeerAddress,
		"enable_nqa":            cnt.EnableNqa,
		"created_at":            cnt.CreatedAt,
		"updated_at":            cnt.UpdatedAt,
		"enterprise_project_id": cnt.EnterpriseProjectID,
		"connection_monitor_id": cnt.ConnectionMonitorID,
		"ha_role":               cnt.HaRole,
		"policy_rules":          flattenGetConnectionPolicyRule(cnt),
		"ikepolicy":             flattenGetConnectionIkePolicy(cnt),
		"ipsecpolicy":           flattenGetConnectionIpsecPolicy(cnt),
	}
	return resourceConnections
}

func flattenGetConnectionIkePolicy(cnt connections.Connections) []interface{} {
	ikePly := cnt.Ikepolicy
	rst := []interface{}{
		map[string]interface{}{
			"authentication_algorithm": ikePly.AuthenticationAlgorithm,
			"encryption_algorithm":     ikePly.EncryptionAlgorithm,
			"ike_version":              ikePly.IkeVersion,
			"lifetime_seconds":         ikePly.LifetimeSeconds,
			"local_id_type":            ikePly.LocalIdType,
			"local_id":                 ikePly.LocalId,
			"peer_id_type":             ikePly.PeerIdType,
			"peer_id":                  ikePly.PeerId,
			"phase1_negotiation_mode":  ikePly.Phase1NegotiationMode,
			"authentication_method":    ikePly.AuthenticationMethod,
			"dh_group":                 ikePly.DhGroup,
			"dpd":                      flattenGetConnectionDPD(cnt),
		},
	}
	return rst
}

func flattenGetConnectionDPD(cnt connections.Connections) []interface{} {
	dpdConncetion := cnt.Ikepolicy.Dpd
	rst := []interface{}{
		map[string]interface{}{
			"timeout":  dpdConncetion.Timeout,
			"interval": dpdConncetion.Interval,
			"msg":      dpdConncetion.Msg,
		},
	}
	return rst
}

func flattenGetConnectionIpsecPolicy(cnt connections.Connections) []interface{} {
	ipsecPly := cnt.Ipsecpolicy
	rst := []interface{}{
		map[string]interface{}{
			"authentication_algorithm": ipsecPly.AuthenticationAlgorithm,
			"encryption_algorithm":     ipsecPly.EncryptionAlgorithm,
			"pfs":                      ipsecPly.Pfs,
			"lifetime_seconds":         ipsecPly.LifetimeSeconds,
			"transform_protocol":       ipsecPly.TransformProtocol,
			"encapsulation_mode":       ipsecPly.EncapsulationMode,
		},
	}
	return rst
}

func flattenGetConnectionPolicyRule(cnt connections.Connections) []map[string]interface{} {
	plyRules := cnt.PolicyRules
	rst := make([]map[string]interface{}, len(plyRules))
	for i, v := range plyRules {
		rst[i] = map[string]interface{}{
			"rule_index":  v.RuleIndex,
			"destination": v.Destination,
			"source":      v.Source,
		}
	}
	return rst
}
