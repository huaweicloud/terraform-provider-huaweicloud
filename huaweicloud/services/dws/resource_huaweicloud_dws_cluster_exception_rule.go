package dws

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	clusterExceptionRuleNonUpdatableParams = []string{
		"cluster_id",
		"name",
	}
	clusterExceptionRuleJSONParamKeys = []string{
		"configurations",
	}
)

// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/workload/rules
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/workload/rules
// @API DWS PUT /v1/{project_id}/clusters/{cluster_id}/workload/rules/{rule_name}
// @API DWS DELETE /v1/{project_id}/clusters/{cluster_id}/workload/rules/{rule_name}
func ResourceClusterExceptionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterExceptionRuleCreate,
		ReadContext:   resourceClusterExceptionRuleRead,
		UpdateContext: resourceClusterExceptionRuleUpdate,
		DeleteContext: resourceClusterExceptionRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterExceptionRuleNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterExceptionRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the exception rule is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the exception rule belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the exception rule.`,
			},
			"configurations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the exception rule configuration.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The value of the exception rule configuration.`,
						},
					},
				},
				DiffSuppressFunc: utils.SuppressObjectSliceDiffs(),
				Description:      `The list of exception rule configurations.`,
			},

			// Internal parameters.
			"configurations_origin": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the exception rule configuration.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the exception rule configuration.",
						},
					},
				},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'configurations'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildClusterExceptionRuleConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, v := range configurations {
		result = append(result, map[string]interface{}{
			"rule_key":   utils.PathSearch("key", v, nil),
			"rule_value": utils.PathSearch("value", v, nil),
		})
	}
	return result
}

func buildClusterExceptionRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rule_name":    d.Get("name"),
		"except_rules": buildClusterExceptionRuleConfigurations(d.Get("configurations").([]interface{})),
	}
}

func createClusterExceptionRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/workload/rules"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         utils.RemoveNil(buildClusterExceptionRuleBodyParams(d)),
		OkCodes:          []int{200},
	}

	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func resourceClusterExceptionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		ruleName  = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = createClusterExceptionRule(client, d)
	if err != nil {
		return diag.Errorf("error creating exception rule (%s): %s", ruleName, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, ruleName))

	// If the request is successful, obtain the values of all JSON parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, clusterExceptionRuleJSONParamKeys)
	if err != nil {
		// Don't report an error if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceClusterExceptionRuleRead(ctx, d, meta)
}

func flattenClusterExceptionRuleConfigurations(configurations map[string]interface{}) []interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(configurations))
	for k, v := range configurations {
		result = append(result, map[string]interface{}{
			"key":   k,
			"value": fmt.Sprintf("%v", v), // Convert the value to a string to avoid the difference in the type of the value.
		})
	}

	return result
}

func listClusterExceptionRules(client *golangsdk.ServiceClient, clusterId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/clusters/{cluster_id}/workload/rules?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 {
		listPath += fmt.Sprintf("&%s", queryParams[0])
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}
		offset += len(items)
	}

	return result, nil
}

func getClusterExceptionRuleByName(client *golangsdk.ServiceClient, clusterId, ruleName string) (interface{}, error) {
	// Reduce the number of queries by specifying query filter parameters (rule_name is a fuzzy query parameter).
	queryParams := fmt.Sprintf("rule_name=%s", ruleName)

	rules, err := listClusterExceptionRules(client, clusterId, queryParams)
	if err != nil {
		return nil, err
	}

	// Exception rules are matched precisely using JMES expressions and based on the rule_name.
	result := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", ruleName), rules, nil)
	if result == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/clusters/{cluster_id}/workload/rules",
				RequestId: "NONE",
				Body:      fmt.Appendf(nil, "the exception rule (%s) is not found", ruleName),
			},
		}
	}

	return result, nil
}

func GetClusterExceptionRuleConfigurations(client *golangsdk.ServiceClient, clusterId, ruleName string, configurationsOrigin []interface{},
	keepRemoteState bool) (interface{}, error) {
	rule, err := getClusterExceptionRuleByName(client, clusterId, ruleName)
	if err != nil {
		return nil, err
	}

	return orderExceptionRuleExceptionRulesByConfigurationsOrigin(
		flattenClusterExceptionRuleConfigurations(utils.PathSearch("except_rules", rule, make(map[string]interface{})).(map[string]interface{})),
		configurationsOrigin, keepRemoteState), nil
}

func orderExceptionRuleExceptionRulesByConfigurationsOrigin(rules, configurationsOrigin []interface{}, keepRemoteState bool) []interface{} {
	if len(configurationsOrigin) < 1 {
		return rules
	}

	sortedConfigurations := make([]interface{}, 0, len(rules))
	rulesCopy := rules

	for _, ruleOrigin := range configurationsOrigin {
		ruleKeyOrigin := utils.PathSearch("key", ruleOrigin, "").(string)
		for index, rule := range rulesCopy {
			ruleKey := utils.PathSearch("key", rule, "").(string)
			if ruleKey != ruleKeyOrigin {
				continue
			}
			// Add the found rule to the sorted rules list.
			sortedConfigurations = append(sortedConfigurations, rulesCopy[index])
			// Remove the processed rule from the original array.
			rulesCopy = append(rulesCopy[:index], rulesCopy[index+1:]...)
			break
		}
	}

	if keepRemoteState {
		sortedConfigurations = append(sortedConfigurations, rulesCopy...)
	}
	return sortedConfigurations
}

func resourceClusterExceptionRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		region               = cfg.GetRegion(d)
		clusterId            = d.Get("cluster_id").(string)
		ruleName             = d.Get("name").(string)
		configurationsOrigin = d.Get("configurations_origin").([]interface{})
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	configurations, err := GetClusterExceptionRuleConfigurations(client, clusterId, ruleName, configurationsOrigin, false)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS exception rule configurations")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Attributes.
		d.Set("configurations", configurations),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildClusterExceptionRuleUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	var (
		rawConfig      = d.GetRawConfig()
		configurations = utils.GetNestedObjectFromRawConfig(rawConfig, "configurations").([]interface{})
	)

	return map[string]interface{}{
		"rule_name":    d.Get("name"),
		"except_rules": buildClusterExceptionRuleConfigurations(configurations),
	}
}

func updateClusterExceptionRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		clusterId = d.Get("cluster_id").(string)
		ruleName  = d.Get("name").(string)
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/workload/rules/{rule_name}"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{rule_name}", ruleName)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		JSONBody:         utils.RemoveNil(buildClusterExceptionRuleUpdateBodyParams(d)),
		OkCodes:          []int{200},
	}

	_, err := client.Request("PUT", updatePath, &updateOpts)
	return err
}

func resourceClusterExceptionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		ruleName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChange("configurations") {
		if err := updateClusterExceptionRule(client, d); err != nil {
			return diag.Errorf("error updating DWS exception rule (%s): %s", ruleName, err)
		}

		// If the request is successful, obtain the values of all JSON parameters first and save them to the corresponding
		// '_origin' attributes for subsequent determination and construction of the request body during next updates.
		// And whether corresponding parameters are changed, the origin values must be refreshed.
		err = utils.RefreshObjectParamOriginValues(d, clusterExceptionRuleJSONParamKeys)
		if err != nil {
			// Don't report an error if origin refresh fails
			log.Printf("[WARN] Unable to refresh the origin values: %s", err)
		}
	}

	return resourceClusterExceptionRuleRead(ctx, d, meta)
}

func deleteClusterExceptionRule(client *golangsdk.ServiceClient, clusterId, ruleName string) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/workload/rules/{rule_name}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{rule_name}", ruleName)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		OkCodes:          []int{200},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpts)
	return err
}

func resourceClusterExceptionRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		ruleName  = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err := deleteClusterExceptionRule(client, clusterId, ruleName); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DWS exception rule (%s)", ruleName))
	}

	return nil
}

func resourceClusterExceptionRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
		d.Set("name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
