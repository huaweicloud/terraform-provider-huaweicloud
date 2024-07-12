package nat

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/nat/v2/dnats"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v2/{project_id}/dnat_rules
// @API NAT GET /v2/{project_id}/dnat_rules/{dnat_rule_id}
// @API NAT PUT /v2/{project_id}/dnat_rules/{dnat_rule_id}
// @API NAT DELETE /v2/{project_id}/nat_gateways/{nat_gateway_id}/dnat_rules/{dnat_rule_id}
func ResourcePublicDnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicDnatRuleCreate,
		ReadContext:   resourcePublicDnatRuleRead,
		UpdateContext: resourcePublicDnatRuleUpdate,
		DeleteContext: resourcePublicDnatRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the DNAT rule is located.",
			},
			"floating_ip_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"floating_ip_id", "global_eip_id"},
				Description:  "The ID of the floating IP address.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the global EIP connected by the DNAT rule.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol type.",
			},
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the NAT gateway to which the DNAT rule belongs.",
			},
			"internal_service_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"internal_service_port_range"},
				RequiredWith: []string{"external_service_port"},
				Description:  "The port used by Floating IP provide services for external systems.",
			},
			"internal_service_port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"external_service_port_range"},
				Description:  "The port used by ECSs or BMSs to provide services for external systems.",
			},
			"external_service_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"external_service_port_range"},
				Description:  "The port range used by Floating IP provide services for external systems.",
			},
			"external_service_port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port range used by ECSs or BMSs to provide services for external systems.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the DNAT rule.",
			},
			"port_id": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"port_id", "private_ip"},
				Optional:     true,
				Computed:     true,
				Description:  "The port ID of network.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private IP address of a user.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the DNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The floating IP address of the DNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global EIP address connected by the DNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the DNAT rule.",
			},
		},
	}
}

func buildPublicDnatRuleCreateOpts(d *schema.ResourceData) dnats.CreateOpts {
	return dnats.CreateOpts{
		GatewayId:                d.Get("nat_gateway_id").(string),
		FloatingIpId:             d.Get("floating_ip_id").(string),
		GlobalEipId:              d.Get("global_eip_id").(string),
		Protocol:                 d.Get("protocol").(string),
		InternalServicePort:      utils.Int(d.Get("internal_service_port").(int)),
		ExternalServicePort:      utils.Int(d.Get("external_service_port").(int)),
		InternalServicePortRange: d.Get("internal_service_port_range").(string),
		EXternalServicePortRange: d.Get("external_service_port_range").(string),
		Description:              d.Get("description").(string),
		PortId:                   d.Get("port_id").(string),
		PrivateIp:                d.Get("private_ip").(string),
	}
}

func publicDnatRuleStateRefreshFunc(client *golangsdk.ServiceClient, ruleId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := dnats.Get(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "COMPLETED", nil
			}
			return resp, "", err
		}

		if utils.StrSliceContains([]string{"INACTIVE", "EIP_FREEZED"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourcePublicDnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	resp, err := dnats.Create(client, buildPublicDnatRuleCreateOpts(d))
	if err != nil {
		return diag.Errorf("error creating DNAT rule: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicDnatRuleStateRefreshFunc(client, d.Id(), []string{"ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePublicDnatRuleRead(ctx, d, meta)
}

func resourcePublicDnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	ruleId := d.Id()
	resp, err := dnats.Get(client, ruleId)
	if err != nil {
		// If the DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving DNAT rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nat_gateway_id", resp.GatewayId),
		d.Set("floating_ip_id", resp.FloatingIpId),
		d.Set("global_eip_id", resp.GlobalEipId),
		d.Set("protocol", resp.Protocol),
		d.Set("internal_service_port", resp.InternalServicePort),
		d.Set("external_service_port", resp.ExternalServicePort),
		d.Set("internal_service_port_range", resp.InternalServicePortRange),
		d.Set("external_service_port_range", resp.EXternalServicePortRange),
		d.Set("description", resp.Description),
		d.Set("port_id", resp.PortId),
		d.Set("private_ip", resp.PrivateIp),
		d.Set("created_at", resp.CreatedAt),
		d.Set("floating_ip_address", resp.FloatingIpAddress),
		d.Set("global_eip_address", resp.GlobalEipAddress),
		d.Set("status", resp.Status),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving DNAT rule fields: %s", err)
	}
	return nil
}

func buildPublicDnatRuleUpdateOpts(d *schema.ResourceData) dnats.UpdateOpts {
	return dnats.UpdateOpts{
		GatewayId:                d.Get("nat_gateway_id").(string),
		FloatingIpId:             d.Get("floating_ip_id").(string),
		GlobalEipId:              d.Get("global_eip_id").(string),
		Protocol:                 d.Get("protocol").(string),
		InternalServicePort:      utils.Int(d.Get("internal_service_port").(int)),
		ExternalServicePort:      utils.Int(d.Get("external_service_port").(int)),
		InternalServicePortRange: d.Get("internal_service_port_range").(string),
		ExternalServicePortRange: d.Get("external_service_port_range").(string),
		Description:              utils.String(d.Get("description").(string)),
		PortId:                   d.Get("port_id").(string),
		PrivateIp:                d.Get("private_ip").(string),
	}
}

func resourcePublicDnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	ruleId := d.Id()
	_, err = dnats.Update(client, ruleId, buildPublicDnatRuleUpdateOpts(d))
	if err != nil {
		return diag.Errorf("error updating DNAT rule (%s): %s", ruleId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicDnatRuleStateRefreshFunc(client, ruleId, []string{"ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePublicDnatRuleRead(ctx, d, meta)
}

func resourcePublicDnatRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	var (
		gatewayId = d.Get("nat_gateway_id").(string)
		ruleId    = d.Id()
	)
	err = dnats.Delete(client, gatewayId, ruleId)
	if err != nil {
		// If the DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting DNAT rule")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicDnatRuleStateRefreshFunc(client, ruleId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
