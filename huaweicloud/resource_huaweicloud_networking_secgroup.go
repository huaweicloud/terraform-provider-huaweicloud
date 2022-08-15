package huaweicloud

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	v1rules "github.com/chnsz/golangsdk/openstack/networking/v1/security/rules"
	v1groups "github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	v2groups "github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	v3groups "github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	v3rules "github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

var securityGroupRuleSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
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
			"port_range_min": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port_range_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_address_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
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
		},
	},
}

func ResourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupCreate,
		ReadContext:   resourceNetworkingSecGroupRead,
		UpdateContext: resourceNetworkingSecGroupUpdate,
		DeleteContext: resourceNetworkingSecGroupDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"delete_default_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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

func resourceNetworkingSecGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	v3Client, err := config.NetworkingV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	// Only name and enterprise project ID are supported.
	createOpts := v3groups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create HuaweiCloud Security Group: %#v", createOpts)
	securityGroup, err := v3groups.Create(v3Client, createOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return resourceNetworkingSecGroupCreateV1(ctx, d, meta)
		}
		return fmtp.DiagErrorf("Error creating Security Group: %s", err)
	}

	d.SetId(securityGroup.ID)

	if val, ok := d.GetOk("description"); ok {
		desc := val.(string)
		updateOpts := v3groups.UpdateOpts{
			Description: &desc,
		}
		_, err = v3groups.Update(v3Client, d.Id(), updateOpts)
		if err != nil {
			return fmtp.DiagErrorf("Error updating the security group (%s) description: %s", d.Id(), err)
		}
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		for _, rule := range securityGroup.SecurityGroupRules {
			if err := v3rules.Delete(v3Client, rule.ID).ExtractErr(); err != nil {
				return fmtp.DiagErrorf("There was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupCreateV1(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	// The v3 API does not exist or has not been published in this region, retry creation using v1 client.
	v1Client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}

	// Only name and enterprise project ID are supported.
	createOpts := v1groups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}
	securityGroup, err := v1groups.Create(v1Client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Security Group: %s", err)
	}
	d.SetId(securityGroup.ID)

	if _, ok := d.GetOk("description"); ok {
		// The v1 API does not support creating and updating methods for description parameters.
		err = resourceNetworkingSecGroupUpdateV2(d, config, region)
		if err != nil {
			return fmtp.DiagErrorf("Error updating description of Security group (%s): %s", d.Id(), err)
		}
	}

	// Delete the default security group rules if it has been requested.
	if d.Get("delete_default_rules").(bool) {
		for _, rule := range securityGroup.SecurityGroupRules {
			if err := v1rules.Delete(v1Client, rule.ID).ExtractErr(); err != nil {
				return fmtp.DiagErrorf("There was a problem deleting a default security group rule: %s", err)
			}
		}
	}
	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	v1Client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}
	v3Client, err := config.NetworkingV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	v1Resp, err := v1groups.Get(v1Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "HuaweiCloud Security group")
	}

	logp.Printf("[DEBUG] Retrieved Security Group (%s) by v1 client: %v", d.Id(), v1Resp)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", v1Resp.Name),
		d.Set("description", v1Resp.Description),
		d.Set("enterprise_project_id", v1Resp.EnterpriseProjectId),
		d.Set("rules", flattenSecurityGroupRulesV1(v1Resp)),
	)

	// If the v3 API is supported, setting related parameters.
	v3Resp, err := v3groups.Get(v3Client, d.Id())
	if err == nil {
		// If the v3 API method has no error, parse its rules list and timestamp attributes and setup.
		logp.Printf("[DEBUG] Retrieved Security Group (%s) by v3 client: %v", d.Id(), v3Resp)
		rules, err := flattenSecurityGroupRulesV3(v3Resp.SecurityGroupRules)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr,
			d.Set("rules", rules), // Override the configuration of the rules list.
			d.Set("created_at", v3Resp.CreatedAt),
			d.Set("updated_at", v3Resp.UpdatedAt),
		)
	}

	// If the query process returns an error, either because the specified region does not exist or the v3 API is
	// not released, or other reasons, skip the setting.
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityGroupRulesV1(secGroup *v1groups.SecurityGroup) []map[string]interface{} {
	sgRules := make([]map[string]interface{}, len(secGroup.SecurityGroupRules))
	for i, rule := range secGroup.SecurityGroupRules {
		sgRules[i] = map[string]interface{}{
			"id":               rule.ID,
			"direction":        rule.Direction,
			"protocol":         rule.Protocol,
			"ethertype":        rule.Ethertype,
			"remote_ip_prefix": rule.RemoteIpPrefix,
			"remote_group_id":  rule.RemoteGroupId,
			"description":      rule.Description,
			"port_range_min":   rule.PortRangeMin,
			"port_range_max":   rule.PortRangeMax,
		}
	}

	return sgRules
}

func flattenSecurityGroupRulesV3(rules []v3rules.SecurityGroupRule) ([]map[string]interface{}, error) {
	sgRules := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		ruleInfo := map[string]interface{}{
			"id":                      rule.ID,
			"direction":               rule.Direction,
			"protocol":                rule.Protocol,
			"ethertype":               rule.Ethertype,
			"remote_ip_prefix":        rule.RemoteIpPrefix,
			"remote_group_id":         rule.RemoteGroupId,
			"remote_address_group_id": rule.RemoteAddressGroupId,
			"description":             rule.Description,
			"action":                  rule.Action,
			"priority":                rule.Priority,
		}
		if rule.MultiPort != "" {
			ruleInfo["ports"] = rule.MultiPort
			if !strings.Contains(rule.MultiPort, ",") {
				re := regexp.MustCompile("^(\\d+)(?:\\-(\\d+))?$")
				rangeSet := re.FindStringSubmatch(rule.MultiPort)
				if len(rangeSet) < 3 {
					logp.Printf("[DEBUG] Regular result for port range (%v) not as expected, should be 3, but %d.",
						rule.MultiPort, len(rule.MultiPort))
					return sgRules, fmtp.Errorf("The parameter format of the 'ports' is invalid.")
				}
				minVal, _ := strconv.Atoi(rangeSet[1])
				ruleInfo["port_range_min"] = minVal
				if rangeSet[2] == "" {
					ruleInfo["port_range_max"] = minVal
				} else {
					maxVal, _ := strconv.Atoi(rangeSet[2])
					ruleInfo["port_range_max"] = maxVal
				}
			}
		}

		sgRules[i] = ruleInfo
	}

	return sgRules, nil
}

func resourceNetworkingSecGroupUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.NetworkingV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	updateOpts := v3groups.UpdateOpts{
		Name:        name,
		Description: &description,
	}

	logp.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
	_, err = v3groups.Update(client, d.Id(), updateOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// The v1 API does not support creating and updating description parameters.
			err = resourceNetworkingSecGroupUpdateV2(d, config, region)
			if err != nil {
				return fmtp.DiagErrorf("Error updating description of security group (%s): %s", d.Id(), err)
			}
		} else {
			return fmtp.DiagErrorf("Error updating security group (%s): %s", d.Id(), err)
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupUpdateV2(d *schema.ResourceData, config *config.Config, region string) error {
	v2Client, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking v2 client: %s", err)
	}

	desc := d.Get("description").(string)
	updateOpts := v2groups.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: &desc,
	}

	_, err = v2groups.Update(v2Client, d.Id(), updateOpts).Extract()

	return err
}

func resourceNetworkingSecGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting security group (%s): %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupDelete(client *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Security Group %s.", secGroupId)

		r, err := v1groups.Get(client, secGroupId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = v1groups.Delete(client, secGroupId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Security Group %s still active", secGroupId)
		return r, "ACTIVE", nil
	}
}
