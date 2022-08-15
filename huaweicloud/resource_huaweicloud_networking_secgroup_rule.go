package huaweicloud

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	v1Rules "github.com/chnsz/golangsdk/openstack/networking/v1/security/rules"
	v3Rules "github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// Some parameters are only support creation in ver.3 API.
var advancedParams = []string{"ports", "remote_address_group_id", "action", "priority"}

func ResourceNetworkingSecGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupRuleCreate,
		ReadContext:   resourceNetworkingSecGroupRuleRead,
		DeleteContext: resourceNetworkingSecGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress", "egress",
				}, true),
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IPv4", "IPv6",
				}, true),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The port range parameters conflict with advanced parameters.
			"port_range_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"protocol"},
			},
			"port_range_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"port_range_min"},
			},
			"ports": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"port_range_min", "port_range_max"},
				RequiredWith:  []string{"protocol"},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"tcp", "udp", "icmp", "icmpv6"}, false),
					validation.StringMatch(regexp.MustCompile("^([0-1]?[0-9]?[0-9]|2[0-4][0-9]|25[0-5])$"),
						"The valid protocol is range from 0 to 255.",
					),
				),
			},
			"remote_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"remote_address_group_id", "remote_ip_prefix"},
			},
			"remote_address_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"port_range_min", "port_range_max"},
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateCIDR,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny",
				}, false),
				ConflictsWith: []string{"port_range_min", "port_range_max"},
			},
			"priority": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IntBetween(1, 100),
				ConflictsWith: []string{"port_range_min", "port_range_max"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func doesAdvanceddParamUsed(d *schema.ResourceData, params []string) bool {
	for _, pk := range params {
		if _, ok := d.GetOk(pk); ok {
			return true
		}
	}
	return false
}

func resourceNetworkingSecGroupRuleCreateV1(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	v1Client, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}

	opt := v1Rules.CreateOpts{
		Description:     d.Get("description").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
		RemoteGroupId:   d.Get("remote_group_id").(string),
		RemoteIpPrefix:  d.Get("remote_ip_prefix").(string),
		Protocol:        d.Get("protocol").(string),
		Ethertype:       d.Get("ethertype").(string),
		Direction:       d.Get("direction").(string),
		PortRangeMin:    d.Get("port_range_min").(int),
		PortRangeMax:    d.Get("port_range_max").(int),
	}

	logp.Printf("[DEBUG] The createOpts of the Security Group rule is: %#v", opt)
	resp, err := v1Rules.Create(v1Client, opt)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Security Group rule: %s", err)
	}
	d.SetId(resp.ID)

	return resourceNetworkingSecGroupRuleRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleCreateV3(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	v3Client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	opt := v3Rules.CreateOpts{
		Description:          d.Get("description").(string),
		SecurityGroupId:      d.Get("security_group_id").(string),
		RemoteGroupId:        d.Get("remote_group_id").(string),
		RemoteAddressGroupId: d.Get("remote_address_group_id").(string),
		RemoteIpPrefix:       d.Get("remote_ip_prefix").(string),
		Protocol:             d.Get("protocol").(string),
		Ethertype:            d.Get("ethertype").(string),
		Direction:            d.Get("direction").(string),
		MultiPort:            d.Get("ports").(string),
		Action:               d.Get("action").(string),
		Priority:             d.Get("priority").(int),
	}

	logp.Printf("[DEBUG] The createOpts of the Security Group rule is: %#v", opt)
	resp, err := v3Rules.Create(v3Client, opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return fmtp.DiagErrorf("The current region does not support creating security group rules through the "+
				"ver.3 API: %#v", err)
		}
		return fmtp.DiagErrorf("Error creating Security Group rule: %s", err)
	}
	d.SetId(resp.ID)

	return resourceNetworkingSecGroupRuleRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	if doesAdvanceddParamUsed(d, advancedParams) {
		return resourceNetworkingSecGroupRuleCreateV3(ctx, d, meta)
	}
	return resourceNetworkingSecGroupRuleCreateV1(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := GetRegion(d, config)

	v1Client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}
	v3Client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	resp, err := v1Rules.Get(v1Client, d.Id())
	if err != nil {
		logp.Printf("[DEBUG] Unable to find the specified Security group rule (%s).", d.Id())
		return common.CheckDeletedDiag(d, err, "HuaweiCloud Security Group Rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("direction", resp.Direction),
		d.Set("description", resp.Description),
		d.Set("ethertype", resp.Ethertype),
		d.Set("protocol", resp.Protocol),
		d.Set("remote_group_id", resp.RemoteGroupId),
		d.Set("remote_ip_prefix", resp.RemoteIpPrefix),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("port_range_min", resp.PortRangeMin),
		d.Set("port_range_max", resp.PortRangeMax),
	)

	rule, err := v3Rules.Get(v3Client, d.Id())
	if err == nil {
		// If the v3 API method has no error, parse its ports attribute and setup.
		logp.Printf("[DEBUG] Retrieved Security Group rule (%s): %+v", d.Id(), rule)
		mErr = multierror.Append(mErr,
			d.Set("ports", rule.MultiPort),
			d.Set("action", rule.Action),
			d.Set("priority", rule.Priority),
			d.Set("remote_address_group_id", rule.RemoteAddressGroupId),
		)
	}

	// If the query process returns an error, either because the specified region does not exist or the v3 API is
	// not released, or other reasons, skip the setting.
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkingSecGroupRuleDelete(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	logp.Printf("[DEBUG] Destroy security group rule: %s", d.Id())

	config := meta.(*config.Config)
	client, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupRuleDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      8 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Security Group Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupRuleDelete(client *golangsdk.ServiceClient, ruleId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Security Group Rule %s.", ruleId)

		r, err := v1Rules.Get(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group Rule %s", ruleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = v1Rules.Delete(client, ruleId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group Rule %s", ruleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Security Group Rule %s still active.", ruleId)
		return r, "ACTIVE", nil
	}
}
