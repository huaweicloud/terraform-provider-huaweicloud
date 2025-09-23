package vpc

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	v1rules "github.com/chnsz/golangsdk/openstack/networking/v1/security/rules"
	v1groups "github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	v2groups "github.com/chnsz/golangsdk/openstack/networking/v2/extensions/security/groups"
	v3groups "github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	v3rules "github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
			"port_range_min": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "schema: Deprecated",
			},
			"port_range_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "schema: Deprecated",
			},
		},
	},
}

// @API VPC PUT /v2.0/security-groups/{id}
// @API VPC DELETE /v3/{project_id}/vpc/security-group-rules/{ruleId}
// @API VPC GET /v3/{project_id}/vpc/security-groups/{secgroupId}
// @API VPC PUT /v3/{project_id}/vpc/security-groups/{secgroupId}
// @API VPC POST /v3/{project_id}/vpc/security-groups
// @API VPC DELETE /v1/{project_id}/security-group-rules/{ruleId}
// @API VPC DELETE /v1/{project_id}/security-groups/{securityGroupId}
// @API VPC GET /v1/{project_id}/security-groups/{securityGroupId}
// @API VPC POST /v1/{project_id}/security-groups
// @API VPC POST /v2.0/{project_id}/security-groups/{id}/tags/action
// @API VPC DELETE /v2.0/{project_id}/security-groups/{id}/tags/action
// @API VPC GET /v2.0/{project_id}/security-groups/{id}/tags
func ResourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupCreate,
		ReadContext:   resourceNetworkingSecGroupRead,
		UpdateContext: resourceNetworkingSecGroupUpdate,
		DeleteContext: resourceNetworkingSecGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"tags":  common.TagsSchema(),
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v3Client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	// Only name and enterprise project ID are supported.
	createOpts := v3groups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create Security Group: %#v", createOpts)
	securityGroup, err := v3groups.Create(v3Client, createOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return resourceNetworkingSecGroupCreateV1(ctx, d, meta)
		}
		return diag.Errorf("error creating Security Group: %s", err)
	}

	d.SetId(securityGroup.ID)

	if val, ok := d.GetOk("description"); ok {
		desc := val.(string)
		updateOpts := v3groups.UpdateOpts{
			Description: &desc,
		}
		_, err = v3groups.Update(v3Client, d.Id(), updateOpts)
		if err != nil {
			return diag.Errorf("error updating the security group (%s) description: %s", d.Id(), err)
		}
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		for _, rule := range securityGroup.SecurityGroupRules {
			if err := v3rules.Delete(v3Client, rule.ID).ExtractErr(); err != nil {
				return diag.Errorf("there was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		v2Client, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating networking v2 client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(v2Client, "security-groups", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of security group %q: %s", d.Id(), tagErr)
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupCreateV1(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// The v3 API does not exist or has not been published in this region, retry creation using v1 client.
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v1 client: %s", err)
	}

	// Only name and enterprise project ID are supported.
	createOpts := v1groups.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
	securityGroup, err := v1groups.Create(v1Client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Security Group: %s", err)
	}
	d.SetId(securityGroup.ID)

	if _, ok := d.GetOk("description"); ok {
		// The v1 API does not support creating and updating methods for description parameters.
		err = resourceNetworkingSecGroupUpdateV2(d, cfg, region)
		if err != nil {
			return diag.Errorf("error updating description of Security group (%s): %s", d.Id(), err)
		}
	}

	// Delete the default security group rules if it has been requested.
	if d.Get("delete_default_rules").(bool) {
		for _, rule := range securityGroup.SecurityGroupRules {
			if err := v1rules.Delete(v1Client, rule.ID).ExtractErr(); err != nil {
				return diag.Errorf("there was a problem deleting a default security group rule: %s", err)
			}
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		v2Client, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating networking v2 client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(v2Client, "security-groups", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of security group %q: %s", d.Id(), tagErr)
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v1 client: %s", err)
	}
	v3Client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	v1Resp, err := v1groups.Get(v1Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Security group")
	}

	log.Printf("[DEBUG] Retrieved Security Group (%s) by v1 client: %v", d.Id(), v1Resp)

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
		log.Printf("[DEBUG] Retrieved Security Group (%s) by v3 client: %v", d.Id(), v3Resp)
		rules, err := flattenSecurityGroupRulesV3(v3Resp.SecurityGroupRules)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr,
			d.Set("rules", rules),                    // Override the configuration of the rules list.
			d.Set("name", v3Resp.Name),               // Override the name
			d.Set("description", v3Resp.Description), // Override the description
			d.Set("created_at", v3Resp.CreatedAt),
			d.Set("updated_at", v3Resp.UpdatedAt),
		)
	}

	// If the query process returns an error, either because the specified region does not exist or the v3 API is
	// not released, or other reasons, skip the setting.

	// save VirtualPrivateCloudV2 tags
	v2Client, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v2 client: %s", err)
	}
	if resourceTags, err := tags.Get(v2Client, "security-groups", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for security group (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of security group (%s): %s", d.Id(), err)
	}
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
				re := regexp.MustCompile(`^(\d+)(?:\-(\d+))?$`)
				rangeSet := re.FindStringSubmatch(rule.MultiPort)
				if len(rangeSet) < 3 {
					log.Printf("[DEBUG] Regular result for port range (%v) not as expected, should be 3, but %d.",
						rule.MultiPort, len(rule.MultiPort))
					return sgRules, fmt.Errorf("the parameter format of the 'ports' is invalid")
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	description := d.Get("description").(string)
	name := d.Get("name").(string)
	updateOpts := v3groups.UpdateOpts{
		Name:        name,
		Description: &description,
	}

	log.Printf("[DEBUG] Updating SecGroup %s with options: %#v", d.Id(), updateOpts)
	_, err = v3groups.Update(client, d.Id(), updateOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// The v1 API does not support creating and updating description parameters.
			err = resourceNetworkingSecGroupUpdateV2(d, cfg, region)
			if err != nil {
				return diag.Errorf("error updating description of security group (%s): %s", d.Id(), err)
			}
		} else {
			return diag.Errorf("error updating security group (%s): %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		v2Client, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating networking v2 client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(v2Client, d, "security-groups", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of security group %s: %s", d.Id(), tagErr)
		}
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupUpdateV2(d *schema.ResourceData, config *config.Config, region string) error {
	v2Client, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating networking v2 client: %s", err)
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
	cfg := meta.(*config.Config)
	client, err := cfg.NetworkingV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking v1 client: %s", err)
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
		return diag.Errorf("error deleting security group (%s): %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupDelete(client *golangsdk.ServiceClient, secGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete Security Group %s.", secGroupId)

		r, err := v1groups.Get(client, secGroupId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = v1groups.Delete(client, secGroupId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted Security Group %s", secGroupId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return r, "ACTIVE", nil
			}
			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] Security Group %s still active", secGroupId)
		return r, "ACTIVE", nil
	}
}
