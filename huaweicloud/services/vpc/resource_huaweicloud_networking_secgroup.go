package vpc

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
	"github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	"github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
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
			"description": {
				Type:     schema.TypeString,
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
				Computed: true,
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

func resourceNetworkingSecGroupCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	// only name and enterprise_project_id are supported
	createOpts := groups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create HuaweiCloud Security Group: %#v", createOpts)
	securityGroup, err := groups.Create(client, createOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Security Group: %s", err)
	}

	d.SetId(securityGroup.ID)

	description := d.Get("description").(string)
	if description != "" {
		updateOpts := groups.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: &description,
		}
		_, err = groups.Update(client, d.Id(), updateOpts)
		if err != nil {
			return fmtp.DiagErrorf("Error updating description of security group %s: %s", d.Id(), err)
		}
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		for _, rule := range securityGroup.SecurityGroupRules {
			if err := rules.Delete(client, rule.ID).ExtractErr(); err != nil {
				return fmtp.DiagErrorf(
					"There was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	logp.Printf("[DEBUG] Retrieve information about security group: %s", d.Id())
	securityGroup, err := groups.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "HuaweiCloud Security group")
	}

	logp.Printf("[DEBUG] Retrieved Security Group %s: %+v", d.Id(), securityGroup)

	secGroupRule, err := flattenSecurityGroupRules(securityGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", securityGroup.Name),
		d.Set("description", securityGroup.Description),
		d.Set("enterprise_project_id", securityGroup.EnterpriseProjectId),
		d.Set("rules", secGroupRule),
		d.Set("created_at", securityGroup.CreatedAt),
		d.Set("updated_at", securityGroup.UpdatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func flattenSecurityGroupRules(secGroup *groups.SecurityGroup) ([]map[string]interface{}, error) {
	sgRules := make([]map[string]interface{}, len(secGroup.SecurityGroupRules))
	for i, rule := range secGroup.SecurityGroupRules {
		ruleInfo := map[string]interface{}{
			"id":               rule.ID,
			"direction":        rule.Direction,
			"protocol":         rule.Protocol,
			"ethertype":        rule.Ethertype,
			"remote_ip_prefix": rule.RemoteIpPrefix,
			"remote_group_id":  rule.RemoteGroupId,
			"description":      rule.Description,
		}
		if rule.MultiPort != "" {
			ruleInfo["ports"] = rule.MultiPort
			if !strings.Contains(rule.MultiPort, ",") {
				re := regexp.MustCompile("^(\\d+)(?:\\-(\\d+))?$")
				rangeSet := re.FindStringSubmatch(rule.MultiPort)
				if len(rangeSet) < 3 {
					return sgRules, fmtp.Errorf("Regular result for port range (%v) not as expected, should be 3, "+
						"but %d.", rule.MultiPort, len(rule.MultiPort))
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
	client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		description := d.Get("description").(string)
		name := d.Get("name").(string)
		updateOpts := groups.UpdateOpts{
			Name:        name,
			Description: &description,
		}

		logp.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
		_, err = groups.Update(client, d.Id(), updateOpts)
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloud SecGroup: %s", err)
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
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
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Security Group: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupDelete(client *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Security Group %s.", secGroupId)

		r, err := groups.Get(client, secGroupId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = groups.Delete(client, secGroupId).ExtractErr()
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
