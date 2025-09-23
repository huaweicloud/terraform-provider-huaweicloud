package iec

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/security/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC POST /v1/security-group-rules
// @API IEC DELETE /v1/security-group-rules/{security_group_rule_id}
// @API IEC GET /v1/security-group-rules/{security_group_rule_id}
func ResourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupRuleCreate,
		ReadContext:   resourceSecurityGroupRuleRead,
		DeleteContext: resourceSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"egress", "ingress",
				}, true),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"icmp", "tcp", "udp", "gre",
				}, true),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port_range_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"port_range_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ethertype": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
				Default:      "IPv4",
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"remote_group_id"},
			},
			"remote_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"remote_ip_prefix"},
			},
		},
	}
}

func resourceSecurityGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	if d.Get("protocol").(string) != "icmp" && d.Get("port_range_min").(int) > d.Get("port_range_max").(int) {
		return diag.Errorf("the value of `port_range_min` can not be greater than the value of `port_range_max`")
	}

	sgRule := ieccommon.ReqSecurityGroupRuleEntity{
		Direction:       d.Get("direction").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		Description:     d.Get("description").(string),
		EtherType:       d.Get("ethertype").(string),
		Protocol:        d.Get("protocol").(string),
		RemoteIPPrefix:  d.Get("remote_ip_prefix").(string),
		RemoteGroupID:   d.Get("remote_group_id").(string),
	}

	if d.Get("protocol").(string) != "icmp" {
		sgRule.PortRangeMin = d.Get("port_range_min")
		sgRule.PortRangeMax = d.Get("port_range_max")
	}

	createOpts := rules.CreateOpts{
		SecurityGroupRule: &sgRule,
	}

	rule, err := rules.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC security group rule: %s", err)
	}

	d.SetId(rule.SecurityGroupRule.ID)
	return resourceSecurityGroupRuleRead(ctx, d, meta)
}

func resourceSecurityGroupRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	rule, err := rules.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IEC security group rule")
	}

	mErr := multierror.Append(
		d.Set("description", rule.SecurityGroupRule.Description),
		d.Set("direction", rule.SecurityGroupRule.Direction),
		d.Set("ethertype", rule.SecurityGroupRule.EtherType),
		d.Set("protocol", rule.SecurityGroupRule.Protocol),
		d.Set("security_group_id", rule.SecurityGroupRule.SecurityGroupID),
		d.Set("remote_ip_prefix", rule.SecurityGroupRule.RemoteIPPrefix),
		d.Set("remote_group_id", rule.SecurityGroupRule.RemoteGroupID),
	)

	if ret, err := strconv.Atoi(rule.SecurityGroupRule.PortRangeMin.(string)); err == nil {
		mErr = multierror.Append(d.Set("port_range_min", ret))
	}
	if ret, err := strconv.Atoi(rule.SecurityGroupRule.PortRangeMax.(string)); err == nil {
		mErr = multierror.Append(d.Set("port_range_max", ret))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields: %s", err)
	}

	return nil
}

func resourceSecurityGroupRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecurityGroupRuleDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      8 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting IEC security group rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecurityGroupRuleDelete(client *golangsdk.ServiceClient, ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rule, err := rules.Get(client, ruleID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted IEC security group rule %s", ruleID)
				return rule, "DELETED", nil
			}
			return err, "ACTIVE", err
		}

		err = rules.Delete(client, ruleID).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted IEC security group rule %s", ruleID)
				return rule, "DELETED", nil
			}
			return rule, "ACTIVE", err
		}
		log.Printf("[DEBUG] IEC security group rule %s still active.\n", ruleID)
		return rule, "ACTIVE", nil
	}
}
