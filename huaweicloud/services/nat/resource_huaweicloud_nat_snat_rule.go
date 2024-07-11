package nat

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/nat/v2/snats"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type SourceType int

const (
	SourceTypeVpc SourceType = 0
	SourceTypeDc  SourceType = 1
)

// @API NAT POST /v2/{project_id}/snat_rules
// @API NAT GET /v2/{project_id}/snat_rules/{snat_rule_id}
// @API NAT PUT /v2/{project_id}/snat_rules/{snat_rule_id}
// @API NAT DELETE /v2/{project_id}/nat_gateways/{nat_gateway_id}/snat_rules/{snat_rule_id}
// @API EIP GET /v1/{project_id}/publicips/{publicip_id}
func ResourcePublicSnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicSnatRuleCreate,
		ReadContext:   resourcePublicSnatRuleRead,
		UpdateContext: resourcePublicSnatRuleUpdate,
		DeleteContext: resourcePublicSnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SNAT rule is located.",
			},
			"floating_ip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"floating_ip_id", "global_eip_id"},
				DiffSuppressFunc: utils.SuppressSnatFiplistDiffs,
				Description:      "The IDs of floating IPs connected by SNAT rule.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IDs (separated by commas) of global EIPs connected by SNAT rule.",
			},
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "schema: Required; The ID of the gateway to which the SNAT rule belongs.",
			},
			"source_type": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					int(SourceTypeVpc),
					int(SourceTypeDc),
				}),
				Optional:    true,
				ForceNew:    true,
				Description: "The resource type of the SNAT rule.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"cidr", "network_id"},
				Description:  "The network IDs of subnet connected by SNAT rule (VPC side).",
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id", "network_id"},
				Description:  "The CIDR block connected by SNAT rule (DC side).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the SNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The floating IP addresses (separated by commas) connected by SNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global EIP addresses (separated by commas) connected by SNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the SNAT rule.",
			},

			// deprecated
			"network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "schema: Deprecated; Use 'subnet_id' instead.",
			},
		},
	}
}

func buildPublicSnatRuleCreateOpts(d *schema.ResourceData) (snats.CreateOpts, error) {
	result := snats.CreateOpts{
		GatewayId:    d.Get("nat_gateway_id").(string),
		FloatingIpId: d.Get("floating_ip_id").(string),
		GlobalEipId:  d.Get("global_eip_id").(string),
		Cidr:         d.Get("cidr").(string),
		Description:  d.Get("description").(string),
	}
	var subnetId string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	} else {
		subnetId = d.Get("network_id").(string)
	}
	result.NetworkId = subnetId

	sourceType := d.Get("source_type").(int)
	if sourceType == 1 && subnetId != "" {
		return result, fmt.Errorf("in the DC (Direct Connect) scenario (source_type is 1), only the parameter 'cidr' " +
			"is valid, and the parameter 'subnet_id' must be empty")
	}
	result.SourceType = sourceType
	return result, nil
}

func publicSnatRuleStateRefreshFunc(client *golangsdk.ServiceClient, ruleId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := snats.Get(client, ruleId)
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

func resourcePublicSnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatGatewayClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	opts, err := buildPublicSnatRuleCreateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] The create options of the public SNAT rule is: %#v", opts)
	resp, err := snats.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating public SNAT rule: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicSnatRuleStateRefreshFunc(client, d.Id(), []string{"ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePublicSnatRuleRead(ctx, d, meta)
}

func resourcePublicSnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	resp, err := snats.Get(natClient, d.Id())
	if err != nil {
		// If the SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving SNAT rule")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nat_gateway_id", resp.GatewayId),
		d.Set("floating_ip_id", resp.FloatingIpId),
		d.Set("floating_ip_address", resp.FloatingIpAddress),
		d.Set("global_eip_id", resp.GlobalEipId),
		d.Set("global_eip_address", resp.GlobalEipAddress),
		d.Set("source_type", resp.SourceType),
		d.Set("subnet_id", resp.NetworkId),
		d.Set("cidr", resp.Cidr),
		d.Set("status", resp.Status),
		d.Set("description", resp.Description),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving public SNAT rule fields: %s", err)
	}
	return nil
}

func resourcePublicSnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	ruleId := d.Id()
	opts := snats.UpdateOpts{
		GatewayId:   d.Get("nat_gateway_id").(string),
		Description: utils.String(d.Get("description").(string)),
	}
	if d.HasChange("floating_ip_id") {
		eipClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v1 client: %s", err)
		}

		eipIds := d.Get("floating_ip_id").(string)
		eipList := strings.Split(eipIds, ",")
		eipAddrs := make([]string, len(eipList))

		// get EIP address from ID
		for i, eipId := range eipList {
			eIP, err := eips.Get(eipClient, eipId).Extract()
			if err != nil {
				return diag.Errorf("error fetching EIP (%s): %s", eipId, err)
			}
			eipAddrs[i] = eIP.PublicAddress
		}

		opts.FloatingIpAddress = strings.Join(eipAddrs, ",")
	}

	if d.HasChange("global_eip_id") {
		opts.GlobalEipId = d.Get("global_eip_id").(string)
	}

	log.Printf("[DEBUG] The update options of the public SNAT rule is: %#v", opts)
	_, err = snats.Update(natClient, ruleId, opts)
	if err != nil {
		return diag.Errorf("error updating public SNAT rule: %s", err)
	}

	log.Printf("[DEBUG] waiting for public SNAT rule (%s) to become available", ruleId)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicSnatRuleStateRefreshFunc(natClient, ruleId, []string{"ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePublicSnatRuleRead(ctx, d, meta)
}

func resourcePublicSnatRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatGatewayClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	ruleId := d.Id()
	gatewayId := d.Get("nat_gateway_id").(string)
	err = snats.Delete(client, gatewayId, ruleId)
	if err != nil {
		// If the SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting SNAT rule")
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicSnatRuleStateRefreshFunc(client, ruleId, nil),
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
