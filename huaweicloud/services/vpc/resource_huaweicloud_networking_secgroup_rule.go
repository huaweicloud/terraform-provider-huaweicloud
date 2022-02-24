package vpc

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
			"port_range_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Deprecated:   "Use ports instead",
				RequiredWith: []string{"protocol"},
			},
			"port_range_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"port_range_min"},
				Deprecated:   "Use ports instead",
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingSecGroupRuleCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	opt := rules.CreateOpts{
		Description:     d.Get("description").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
		RemoteGroupId:   d.Get("remote_group_id").(string),
		RemoteIpPrefix:  d.Get("remote_ip_prefix").(string),
		Protocol:        d.Get("protocol").(string),
		Ethertype:       d.Get("ethertype").(string),
		Direction:       d.Get("direction").(string),
	}
	ports := buildNetworkingSecGroupRulePorts(d)
	if ports != "" {
		opt.MultiPort = ports
	}

	logp.Printf("[DEBUG] Create HuaweiCloud security group: %#v", opt)

	resp, err := rules.Create(client, opt)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] HuaweiCloud Security Group Rule created: %#v", resp)

	d.SetId(resp.ID)

	return resourceNetworkingSecGroupRuleRead(ctx, d, meta)
}

func buildNetworkingSecGroupRulePorts(d *schema.ResourceData) string {
	if val, ok := d.GetOk("port_range_min"); ok {
		if d.Get("port_range_max").(int) == val.(int) {
			return strconv.Itoa(val.(int))
		}
		return fmt.Sprintf("%d-%d", val.(int), d.Get("port_range_max").(int))
	}
	return d.Get("ports").(string)
}

func resourceNetworkingSecGroupRuleRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	logp.Printf("[DEBUG] Retrieve information about security group rule: %s", d.Id())

	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking client: %s", err)
	}

	resp, err := rules.Get(client, d.Id())

	if err != nil {
		return common.CheckDeletedDiag(d, err, "HuaweiCloud Security Group Rule")
	}

	mErr := multierror.Append(nil,
		d.Set("direction", resp.Direction),
		d.Set("description", resp.Description),
		d.Set("ethertype", resp.Ethertype),
		d.Set("protocol", resp.Protocol),
		d.Set("remote_group_id", resp.RemoteGroupId),
		d.Set("remote_ip_prefix", resp.RemoteIpPrefix),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("region", config.GetRegion(d)),
		setSecurityGroupRulePorts(d, resp.MultiPort),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func setSecurityGroupRulePorts(d *schema.ResourceData, portRange string) error {
	if portRange == "" {
		return nil
	}
	mErr := multierror.Append(nil, d.Set("ports", portRange))

	if !strings.Contains(portRange, ",") {
		re := regexp.MustCompile("^(\\d+)(?:\\-(\\d+))?$")
		rangeSet := re.FindStringSubmatch(portRange)
		if len(rangeSet) < 3 {
			return fmtp.Errorf("Regular result for port range (%v) not as expected, should be 3, but %d.", portRange,
				len(portRange))
		}
		minVal, _ := strconv.Atoi(rangeSet[1]) // The 2nd element type has passed the regular expression check.
		mErr = multierror.Append(mErr, d.Set("port_range_min", minVal))

		// For simple numbers, the last element of string array of the regular result is empty, but for range, it will
		// be a string number.
		if rangeSet[2] == "" {
			mErr = multierror.Append(mErr, d.Set("port_range_max", minVal))
		} else {
			maxVal, _ := strconv.Atoi(rangeSet[2])
			mErr = multierror.Append(mErr, d.Set("port_range_max", maxVal))
		}
	}
	return mErr.ErrorOrNil()
}

func resourceNetworkingSecGroupRuleDelete(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	logp.Printf("[DEBUG] Destroy security group rule: %s", d.Id())

	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking client: %s", err)
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

		r, err := rules.Get(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group Rule %s", ruleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = rules.Delete(client, ruleId).ExtractErr()
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
