package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/firewalls
// @API VPC GET /v3/{project_id}/vpc/firewalls/{id}
// @API VPC PUT /v3/{project_id}/vpc/firewalls/{id}
// @API VPC PUT /v3/{project_id}/vpc/firewalls/{firewall_id}/insert-rules
// @API VPC PUT /v3/{project_id}/vpc/firewalls/{firewall_id}/remove-rules
// @API VPC PUT /v3/{project_id}/vpc/firewalls/{firewall_id}/associate-subnets
// @API VPC PUT /v3/{project_id}/vpc/firewalls/{firewall_id}/disassociate-subnets
// @API VPC POST /v3/{project_id}/firewalls/{id}/tags/create
// @API VPC POST /v3/{project_id}/firewalls/{id}/tags/delete
// @API VPC GET /v3/{project_id}/firewalls/{id}/tags
// @API VPC DELETE /v3/{project_id}/vpc/firewalls/{id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources/filter
func ResourceNetworkAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkAclCreate,
		ReadContext:   resourceNetworkAclRead,
		UpdateContext: resourceNetworkAclUpdate,
		DeleteContext: resourceNetworkAclDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ingress_rules": {
				Type:     schema.TypeList,
				Elem:     networkAclRuleSchema(),
				Optional: true,
			},
			"egress_rules": {
				Type:     schema.TypeList,
				Elem:     networkAclRuleSchema(),
				Optional: true,
			},
			"associated_subnets": {
				Type:     schema.TypeSet,
				Elem:     networkAclSubnetSchema(),
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"status": {
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
	}
}

func networkAclRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_ip_address_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ip_address_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func networkAclSubnetSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceNetworkAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	createNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls"
	createNetworkAclPath := client.Endpoint + createNetworkAclHttpUrl
	createNetworkAclPath = strings.ReplaceAll(createNetworkAclPath, "{project_id}", client.ProjectID)

	createNetworkAclOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createNetworkAclOpt.JSONBody = utils.RemoveNil(buildCreateNetworkAclBodyParams(d, cfg))
	createNetworkAclResp, err := client.Request("POST", createNetworkAclPath, &createNetworkAclOpt)
	if err != nil {
		return diag.Errorf("error creating network ACL: %s", err)
	}

	createNetworkAclRespBody, err := utils.FlattenResponse(createNetworkAclResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("firewall.id", createNetworkAclRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find ID in API response: %s", err)
	}

	d.SetId(id)

	if v, ok := d.GetOk("ingress_rules"); ok {
		err = networkAclInsertRules(client, v.([]interface{}), "ingress_rules", id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("egress_rules"); ok {
		err = networkAclInsertRules(client, v.([]interface{}), "egress_rules", id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("associated_subnets"); ok {
		err = networkAclAssociatSubnets(client, v.(*schema.Set).List(), id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		err := doTagsAction(client, tagRaw, id, "create")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNetworkAclRead(ctx, d, meta)
}

func networkAclInsertRules(client *golangsdk.ServiceClient, rules []interface{}, ruleType, id string) error {
	networkAclInsertRulesHttpUrl := "v3/{project_id}/vpc/firewalls/{firewall_id}/insert-rules"
	networkAclInsertRulesPath := client.Endpoint + networkAclInsertRulesHttpUrl
	networkAclInsertRulesPath = strings.ReplaceAll(networkAclInsertRulesPath, "{project_id}", client.ProjectID)
	networkAclInsertRulesPath = strings.ReplaceAll(networkAclInsertRulesPath, "{firewall_id}", id)

	networkAclInsertRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	networkAclInsertRulesOpt.JSONBody = utils.RemoveNil(buildNetworkAclInsertRulesBodyParams(rules, ruleType))
	_, err := client.Request("PUT", networkAclInsertRulesPath, &networkAclInsertRulesOpt)
	if err != nil {
		return fmt.Errorf("error inserting rules to network ACL(%s): %s", id, err)
	}

	return nil
}

func networkAclRemoveRules(client *golangsdk.ServiceClient, rules []interface{}, ruleType, id string) error {
	networkAclRemoveRulesHttpUrl := "v3/{project_id}/vpc/firewalls/{firewall_id}/remove-rules"
	networkAclRemoveRulesPath := client.Endpoint + networkAclRemoveRulesHttpUrl
	networkAclRemoveRulesPath = strings.ReplaceAll(networkAclRemoveRulesPath, "{project_id}", client.ProjectID)
	networkAclRemoveRulesPath = strings.ReplaceAll(networkAclRemoveRulesPath, "{firewall_id}", id)

	networkAclRemoveRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	networkAclRemoveRulesOpt.JSONBody = utils.RemoveNil(buildNetworkAclRemoveRulesBodyParams(rules, ruleType))
	_, err := client.Request("PUT", networkAclRemoveRulesPath, &networkAclRemoveRulesOpt)
	if err != nil {
		return fmt.Errorf("error removing rules from network ACL(%s): %s", id, err)
	}

	return nil
}

func networkAclAssociatSubnets(client *golangsdk.ServiceClient, subnets []interface{}, id string) error {
	networkAclAssociatSubnetsHttpUrl := "v3/{project_id}/vpc/firewalls/{firewall_id}/associate-subnets"
	networkAclAssociatSubnetsPath := client.Endpoint + networkAclAssociatSubnetsHttpUrl
	networkAclAssociatSubnetsPath = strings.ReplaceAll(networkAclAssociatSubnetsPath, "{project_id}", client.ProjectID)
	networkAclAssociatSubnetsPath = strings.ReplaceAll(networkAclAssociatSubnetsPath, "{firewall_id}", id)

	networkAclAssociatSubnetsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	networkAclAssociatSubnetsOpt.JSONBody = utils.RemoveNil(buildNetworkAclSubnetsBodyParams(subnets))
	_, err := client.Request("PUT", networkAclAssociatSubnetsPath, &networkAclAssociatSubnetsOpt)
	if err != nil {
		return fmt.Errorf("error associating subnets to network ACL(%s): %s", id, err)
	}

	return nil
}

func networkAclDisassociatSubnets(client *golangsdk.ServiceClient, subnets []interface{}, id string) error {
	networkAclDisassociatSubnetsHttpUrl := "v3/{project_id}/vpc/firewalls/{firewall_id}/disassociate-subnets"
	networkAclDisassociatSubnetsPath := client.Endpoint + networkAclDisassociatSubnetsHttpUrl
	networkAclDisassociatSubnetsPath = strings.ReplaceAll(networkAclDisassociatSubnetsPath, "{project_id}", client.ProjectID)
	networkAclDisassociatSubnetsPath = strings.ReplaceAll(networkAclDisassociatSubnetsPath, "{firewall_id}", id)

	networkAclDisassociatSubnetsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	networkAclDisassociatSubnetsOpt.JSONBody = utils.RemoveNil(buildNetworkAclSubnetsBodyParams(subnets))
	_, err := client.Request("PUT", networkAclDisassociatSubnetsPath, &networkAclDisassociatSubnetsOpt)
	if err != nil {
		return fmt.Errorf("error disassociating subnets from network ACL(%s): %s", id, err)
	}

	return nil
}

func buildCreateNetworkAclBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"firewall": map[string]interface{}{
			"name":                  d.Get("name"),
			"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"admin_state_up":        d.Get("enabled"),
		},
	}
	return bodyParams
}

func buildNetworkAclInsertRulesBodyParams(rules []interface{}, ruleType string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"firewall": map[string]interface{}{
			ruleType: buildNetworkAclInsertRuleBodyParams(rules),
		},
	}

	return bodyParams
}

func buildNetworkAclInsertRuleBodyParams(rules []interface{}) []map[string]interface{} {
	bodyParams := make([]map[string]interface{}, len(rules))
	for i, v := range rules {
		rule := v.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"action":                       rule["action"],
			"protocol":                     rule["protocol"],
			"ip_version":                   rule["ip_version"],
			"description":                  rule["description"],
			"name":                         utils.ValueIgnoreEmpty(rule["name"]),
			"source_ip_address":            utils.ValueIgnoreEmpty(rule["source_ip_address"]),
			"destination_ip_address":       utils.ValueIgnoreEmpty(rule["destination_ip_address"]),
			"source_port":                  utils.ValueIgnoreEmpty(rule["source_port"]),
			"destination_port":             utils.ValueIgnoreEmpty(rule["destination_port"]),
			"source_address_group_id":      utils.ValueIgnoreEmpty(rule["source_ip_address_group_id"]),
			"destination_address_group_id": utils.ValueIgnoreEmpty(rule["destination_ip_address_group_id"]),
		}
	}

	return bodyParams
}

func buildNetworkAclRemoveRulesBodyParams(rules []interface{}, ruleType string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"firewall": map[string]interface{}{
			ruleType: buildNetworkAclRemoveRuleBodyParams(rules),
		},
	}

	return bodyParams
}

func buildNetworkAclRemoveRuleBodyParams(rules []interface{}) []map[string]interface{} {
	bodyParams := make([]map[string]interface{}, len(rules))
	for i, v := range rules {
		rule := v.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"id": rule["rule_id"],
		}
	}

	return bodyParams
}

func buildNetworkAclSubnetsBodyParams(subnets []interface{}) map[string]interface{} {
	bodyParams := make([]map[string]interface{}, len(subnets))
	for i, v := range subnets {
		subnet := v.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"virsubnet_id": subnet["subnet_id"],
		}
	}

	return map[string]interface{}{
		"subnets": bodyParams,
	}
}

func resourceNetworkAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls/" + d.Id()
	getNetworkAclPath := client.Endpoint + getNetworkAclHttpUrl
	getNetworkAclPath = strings.ReplaceAll(getNetworkAclPath, "{project_id}", client.ProjectID)

	getNetworkAclOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getNetworkAclResp, err := client.Request("GET", getNetworkAclPath, &getNetworkAclOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC network ACL")
	}

	getNetworkAclRespBody, err := utils.FlattenResponse(getNetworkAclResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("firewall.name", getNetworkAclRespBody, nil)),
		d.Set("description", utils.PathSearch("firewall.description", getNetworkAclRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("firewall.enterprise_project_id", getNetworkAclRespBody, nil)),
		d.Set("enabled", utils.PathSearch("firewall.admin_state_up", getNetworkAclRespBody, nil)),
		d.Set("ingress_rules", flattenRules(utils.PathSearch("firewall.ingress_rules", getNetworkAclRespBody, nil))),
		d.Set("egress_rules", flattenRules(utils.PathSearch("firewall.egress_rules", getNetworkAclRespBody, nil))),
		d.Set("associated_subnets", flattenSubnets(utils.PathSearch("firewall.associations", getNetworkAclRespBody, nil))),
		d.Set("status", utils.PathSearch("firewall.status", getNetworkAclRespBody, nil)),
		d.Set("created_at", utils.PathSearch("firewall.created_at", getNetworkAclRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("firewall.updated_at", getNetworkAclRespBody, nil)),
	)

	if resourceTags, err := tags.Get(client, "firewalls", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for network ACL (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of network ACL (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRules(rulesRaw interface{}) []map[string]interface{} {
	if rulesRaw == nil {
		return nil
	}

	rules := rulesRaw.([]interface{})
	res := make([]map[string]interface{}, len(rules))
	for i, v := range rules {
		rule := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"rule_id":                         rule["id"],
			"name":                            rule["name"],
			"description":                     rule["description"],
			"action":                          rule["action"],
			"protocol":                        rule["protocol"],
			"ip_version":                      rule["ip_version"],
			"source_ip_address":               rule["source_ip_address"],
			"destination_ip_address":          rule["destination_ip_address"],
			"source_port":                     rule["source_port"],
			"destination_port":                rule["destination_port"],
			"source_ip_address_group_id":      rule["source_address_group_id"],
			"destination_ip_address_group_id": rule["destination_address_group_id"],
		}
	}

	return res
}

func flattenSubnets(subnetsRaw interface{}) []map[string]interface{} {
	if subnetsRaw == nil {
		return nil
	}

	subnets := subnetsRaw.([]interface{})
	res := make([]map[string]interface{}, len(subnets))
	for i, v := range subnets {
		subnet := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"subnet_id": subnet["virsubnet_id"],
		}
	}

	return res
}

func buildUpdateNetworkAclBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"firewall": map[string]interface{}{
			"name":           d.Get("name"),
			"description":    d.Get("description"),
			"admin_state_up": d.Get("enabled"),
		},
	}
	return bodyParams
}

func resourceNetworkAclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	id := d.Id()
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 Client: %s", err)
	}

	if d.HasChanges("name", "description", "enabled") {
		updateNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls/" + id
		updateNetworkAclPath := client.Endpoint + updateNetworkAclHttpUrl
		updateNetworkAclPath = strings.ReplaceAll(updateNetworkAclPath, "{project_id}", client.ProjectID)

		updateNetworkAclOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateNetworkAclOpts.JSONBody = utils.RemoveNil(buildUpdateNetworkAclBodyParams(d))
		_, err = client.Request("PUT", updateNetworkAclPath, &updateNetworkAclOpts)
		if err != nil {
			return diag.Errorf("error updating VPC network ACL: %s", err)
		}
	}

	if d.HasChange("ingress_rules") {
		err = updateRules(client, d, "ingress_rules")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("egress_rules") {
		err = updateRules(client, d, "egress_rules")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("associated_subnets") {
		err = updateAssociatedSubnets(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   id,
			ResourceType: "firewalls",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(client, d, id)
		if tagErr != nil {
			return diag.Errorf("error updating tags of network ACL %s: %s", d.Id(), tagErr)
		}
	}

	return resourceNetworkAclRead(ctx, d, meta)
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		err := doTagsAction(client, oMap, id, "delete")
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		err := doTagsAction(client, nMap, id, "create")
		if err != nil {
			return err
		}
	}

	return nil
}

func doTagsAction(client *golangsdk.ServiceClient, tagsRaw map[string]interface{}, id, action string) error {
	manageTagsHttpUrl := "v3/{project_id}/firewalls/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{project_id}", client.ProjectID)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", id)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{action}", action)
	manageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	manageTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tagsRaw),
	}

	_, err := client.Request("POST", manageTagsPath, &manageTagsOpt)
	if err != nil {
		return fmt.Errorf("unable to %s network ACL tags: %s", action, err)
	}

	return nil
}

func updateRules(client *golangsdk.ServiceClient, d *schema.ResourceData, ruleType string) error {
	id := d.Id()
	oldRules, newRules := d.GetChange(ruleType)
	if len(oldRules.([]interface{})) > 0 {
		err := networkAclRemoveRules(client, oldRules.([]interface{}), ruleType, id)
		if err != nil {
			return fmt.Errorf("error updating rules: %s", err)
		}
	}

	if len(newRules.([]interface{})) > 0 {
		err := networkAclInsertRules(client, newRules.([]interface{}), ruleType, id)
		if err != nil {
			// if failed to insert the new rules, insert the old rules back
			if len(oldRules.([]interface{})) > 0 {
				rollBackErr := networkAclInsertRules(client, oldRules.([]interface{}), ruleType, id)
				if rollBackErr != nil {
					return fmt.Errorf("error updating rules: %s, failed to roll back: %s", err, rollBackErr)
				}
				return fmt.Errorf("error updating rules: %s, it's rolled back", err)
			}
			return fmt.Errorf("error updating rules: %s", err)
		}
	}

	return nil
}

func updateAssociatedSubnets(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	id := d.Id()
	oldSubnets, newSubnets := d.GetChange("associated_subnets")
	if len(oldSubnets.(*schema.Set).List()) > 0 {
		err := networkAclDisassociatSubnets(client, oldSubnets.(*schema.Set).List(), id)
		if err != nil {
			return err
		}
	}

	if len(newSubnets.(*schema.Set).List()) > 0 {
		err := networkAclAssociatSubnets(client, newSubnets.(*schema.Set).List(), id)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceNetworkAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	id := d.Id()
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 Client: %s", err)
	}

	// the associated subnets need to be removed before deleting the ACL
	if d.Get("associated_subnets").(*schema.Set).Len() > 0 {
		err = networkAclDisassociatSubnets(client, d.Get("associated_subnets").(*schema.Set).List(), id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	deleteNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls/" + id
	deleteNetworkAclPath := client.Endpoint + deleteNetworkAclHttpUrl
	deleteNetworkAclPath = strings.ReplaceAll(deleteNetworkAclPath, "{project_id}", client.ProjectID)

	deleteNetworkAclOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteNetworkAclPath, &deleteNetworkAclOpt)
	if err != nil {
		return diag.Errorf("error deleting network ACL: %s", err)
	}

	return nil
}
