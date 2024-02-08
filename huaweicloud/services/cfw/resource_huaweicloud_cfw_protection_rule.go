// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW POST /v1/{project_id}/acl-rule
// @API CFW DELETE /v1/{project_id}/acl-rule/{id}
// @API CFW PUT /v1/{project_id}/acl-rule/{id}
// @API CFW GET /v1/{project_id}/acl-rules
// @API CFW PUT /v1/{project_id}/acl-rule/order/{id}
func ResourceProtectionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectionRuleCreate,
		UpdateContext: resourceProtectionRuleUpdate,
		ReadContext:   resourceProtectionRuleRead,
		DeleteContext: resourceProtectionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProtectionRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The rule name.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The protected object ID`,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"the input is invalid"),
			},
			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `The rule type.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"action_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The action type.`,
			},
			"address_type": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `The address type.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"sequence": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        ProtectionRuleOrderRuleAclDtoSchema(),
				Required:    true,
				Description: `The sequence configuration.`,
			},
			"service": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        ProtectionRuleRuleServiceDtoSchema(),
				Required:    true,
				Description: `The service configuration.`,
			},
			"source": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        ProtectionRuleRuleAddressDtoSchema(),
				Required:    true,
				Description: `The source configuration.`,
			},
			"destination": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        ProtectionRuleRuleAddressDtoSchema(),
				Required:    true,
				Description: `The destination configuration.`,
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `The rule status.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"long_connect_enable": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `Whether to support persistent connections.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"long_connect_time_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The persistent connection duration (hour).`,
			},
			"long_connect_time_minute": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The persistent connection duration (minute).`,
			},
			"long_connect_time_second": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The persistent Connection Duration (second).`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description.`,
			},
			"direction": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  `The direction.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func ProtectionRuleOrderRuleAclDtoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"dest_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the rule that the added rule will follow.`,
			},
			"top": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Whether to pin on top.`,
			},
		},
	}
	return &sc
}

func ProtectionRuleRuleServiceDtoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The service input type.`,
			},
			"dest_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The destination port.`,
			},
			"protocol": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The protocol type.`,
			},
			"service_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The service group ID.`,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
					"the input is invalid"),
			},
			"service_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The service group name.`,
			},
			"source_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The source port.`,
			},
		},
	}
	return &sc
}

func ProtectionRuleRuleAddressDtoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The Source type.`,
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IP address.`,
			},
			"address_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the associated IP address group.`,
			},
			"address_set_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IP address group name.`,
			},
			"address_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The address type.`,
			},
			"domain_address_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the domain name address.`,
			},
		},
	}
	return &sc
}

func resourceProtectionRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createProtectionRule: Create a CFW Protection Rule.
	var (
		createProtectionRuleHttpUrl = "v1/{project_id}/acl-rule"
		createProtectionRuleProduct = "cfw"
	)
	createProtectionRuleClient, err := conf.NewServiceClient(createProtectionRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule Client: %s", err)
	}

	createProtectionRulePath := createProtectionRuleClient.Endpoint + createProtectionRuleHttpUrl
	createProtectionRulePath = strings.ReplaceAll(createProtectionRulePath, "{project_id}", createProtectionRuleClient.ProjectID)

	createProtectionRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createProtectionRuleOpt.JSONBody = utils.RemoveNil(buildCreateProtectionRuleBodyParams(d))
	createProtectionRuleResp, err := createProtectionRuleClient.Request("POST", createProtectionRulePath, &createProtectionRuleOpt)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule: %s", err)
	}

	createProtectionRuleRespBody, err := utils.FlattenResponse(createProtectionRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.rules[0].id", createProtectionRuleRespBody)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule: ID is not found in API response: %s", err)
	}
	d.SetId(id.(string))

	return resourceProtectionRuleRead(ctx, d, meta)
}

func buildCreateProtectionRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id": utils.ValueIngoreEmpty(d.Get("object_id")),
		"type":      d.Get("type"),
		"rules":     buildCreateProtectionRulesOpts(d),
	}
	return bodyParams
}

func buildCreateProtectionRulesOpts(d *schema.ResourceData) []map[string]interface{} {
	params := map[string]interface{}{
		"action_type":              d.Get("action_type"),
		"address_type":             d.Get("address_type"),
		"description":              utils.ValueIngoreEmpty(d.Get("description")),
		"direction":                d.Get("direction"),
		"long_connect_enable":      d.Get("long_connect_enable"),
		"long_connect_time_hour":   utils.ValueIngoreEmpty(d.Get("long_connect_time_hour")),
		"long_connect_time_minute": utils.ValueIngoreEmpty(d.Get("long_connect_time_minute")),
		"long_connect_time_second": utils.ValueIngoreEmpty(d.Get("long_connect_time_second")),
		"name":                     utils.ValueIngoreEmpty(d.Get("name")),
		"sequence":                 buildCreateProtectionRuleRequestBodyOrderRuleAclDto(d.Get("sequence")),
		"service":                  buildCreateProtectionRuleRequestBodyRuleServiceDto(d.Get("service")),
		"source":                   buildCreateProtectionRuleRequestBodyRuleAddressDto(d.Get("source")),
		"destination":              buildCreateProtectionRuleRequestBodyRuleAddressDto(d.Get("destination")),
		"status":                   d.Get("status"),
	}

	return []map[string]interface{}{params}
}

func buildCreateProtectionRuleRequestBodyOrderRuleAclDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"dest_rule_id": utils.ValueIngoreEmpty(raw["dest_rule_id"]),
			"top":          raw["top"],
		}
		return params
	}
	return nil
}

func buildCreateProtectionRuleRequestBodyRuleServiceDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"dest_port":        utils.ValueIngoreEmpty(raw["dest_port"]),
			"protocol":         utils.ValueIngoreEmpty(raw["protocol"]),
			"service_set_id":   utils.ValueIngoreEmpty(raw["service_set_id"]),
			"service_set_name": utils.ValueIngoreEmpty(raw["service_set_name"]),
			"source_port":      utils.ValueIngoreEmpty(raw["source_port"]),
			"type":             raw["type"],
		}
		return params
	}
	return nil
}

func buildCreateProtectionRuleRequestBodyRuleAddressDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"address":             utils.ValueIngoreEmpty(raw["address"]),
			"address_set_id":      utils.ValueIngoreEmpty(raw["address_set_id"]),
			"address_set_name":    utils.ValueIngoreEmpty(raw["address_set_name"]),
			"address_type":        utils.ValueIngoreEmpty(raw["address_type"]),
			"domain_address_name": utils.ValueIngoreEmpty(raw["domain_address_name"]),
			"type":                raw["type"],
		}
		return params
	}
	return nil
}

func resourceProtectionRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getProtectionRule: Query the CFW Protection Rule detail
	var (
		getProtectionRuleHttpUrl = "v1/{project_id}/acl-rules"
		getProtectionRuleProduct = "cfw"
	)
	getProtectionRuleClient, err := conf.NewServiceClient(getProtectionRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule Client: %s", err)
	}

	getProtectionRulePath := getProtectionRuleClient.Endpoint + getProtectionRuleHttpUrl
	getProtectionRulePath = strings.ReplaceAll(getProtectionRulePath, "{project_id}", getProtectionRuleClient.ProjectID)

	getProtectionRulequeryParams := buildGetProtectionRuleQueryParams(d)
	getProtectionRulePath += getProtectionRulequeryParams

	getPotectionRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProtectionRuleResp, err := getProtectionRuleClient.Request("GET", getProtectionRulePath, &getPotectionRulesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving protection rule")
	}

	getProtectionRuleRespBody, err := utils.FlattenResponse(getProtectionRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := jmespath.Search("data.records", getProtectionRuleRespBody)
	if err != nil {
		diag.Errorf("error parsing data.records from response= %#v", getProtectionRuleRespBody)
	}

	rule, err := FilterRules(rules.([]interface{}), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving protection rule")
	}

	// the params 'sequence' and 'type 'not not returned
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("action_type", utils.PathSearch("action_type", rule, nil)),
		d.Set("address_type", utils.PathSearch("address_type", rule, nil)),
		d.Set("description", utils.PathSearch("description", rule, nil)),
		d.Set("direction", utils.PathSearch("direction", rule, nil)),
		d.Set("long_connect_enable", utils.PathSearch("long_connect_enable", rule, nil)),
		d.Set("long_connect_time_hour", utils.PathSearch("long_connect_time_hour", rule, nil)),
		d.Set("long_connect_time_minute", utils.PathSearch("long_connect_time_minute", rule, nil)),
		d.Set("long_connect_time_second", utils.PathSearch("long_connect_time_second", rule, nil)),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("service", flattenGetProtectionRuleResponseBodyRuleServiceDto(rule)),
		d.Set("source", flattenGetProtectionRuleResponseBodyRuleSourceAddressDto(rule)),
		d.Set("destination", flattenGetProtectionRuleResponseBodyRuleDestinationAddressDto(rule)),
		d.Set("status", utils.PathSearch("status", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func FilterRules(rules []interface{}, id string) (interface{}, error) {
	if len(rules) != 0 {
		for _, v := range rules {
			rule := v.(map[string]interface{})
			if rule["rule_id"] == id {
				return v, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func flattenGetProtectionRuleResponseBodyRuleServiceDto(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("service", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing service from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"dest_port":        utils.PathSearch("dest_port", curJson, nil),
			"protocol":         utils.PathSearch("protocol", curJson, nil),
			"service_set_id":   utils.PathSearch("service_set_id", curJson, nil),
			"service_set_name": utils.PathSearch("service_set_name", curJson, nil),
			"source_port":      utils.PathSearch("source_port", curJson, nil),
			"type":             utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func flattenGetProtectionRuleResponseBodyRuleSourceAddressDto(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("source", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing source from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"address":             utils.PathSearch("address", curJson, nil),
			"address_set_id":      utils.PathSearch("address_set_id", curJson, nil),
			"address_set_name":    utils.PathSearch("address_set_name", curJson, nil),
			"address_type":        utils.PathSearch("address_type", curJson, nil),
			"domain_address_name": utils.PathSearch("domain_address_name", curJson, nil),
			"type":                utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func flattenGetProtectionRuleResponseBodyRuleDestinationAddressDto(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("destination", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing destination from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"address":             utils.PathSearch("address", curJson, nil),
			"address_set_id":      utils.PathSearch("address_set_id", curJson, nil),
			"address_set_name":    utils.PathSearch("address_set_name", curJson, nil),
			"address_type":        utils.PathSearch("address_type", curJson, nil),
			"domain_address_name": utils.PathSearch("domain_address_name", curJson, nil),
			"type":                utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func buildGetProtectionRuleQueryParams(d *schema.ResourceData) string {
	res := "?offset=0&limit=1024"
	res = fmt.Sprintf("%s&object_id=%v", res, d.Get("object_id"))

	return res
}

func resourceProtectionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	updateProtectionRulehasChanges := []string{
		"action_type",
		"address_type",
		"description",
		"destination",
		"direction",
		"long_connect_enable",
		"long_connect_time_hour",
		"long_connect_time_minute",
		"long_connect_time_second",
		"name",
		"service",
		"source",
		"status",
		"type",
	}

	var (
		updateProtectionRuleHttpUrl      = "v1/{project_id}/acl-rule/{id}"
		updateProtectionRuleOrderHttpUrl = "v1/{project_id}/acl-rule/order/{id}"
		updateProtectionRuleProduct      = "cfw"
	)
	updateProtectionRuleClient, err := conf.NewServiceClient(updateProtectionRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule Client: %s", err)
	}

	if d.HasChanges(updateProtectionRulehasChanges...) {
		// updateProtectionRule: Update the configuration of CFW Protection Rule
		updateProtectionRulePath := updateProtectionRuleClient.Endpoint + updateProtectionRuleHttpUrl
		updateProtectionRulePath = strings.ReplaceAll(updateProtectionRulePath, "{project_id}", updateProtectionRuleClient.ProjectID)
		updateProtectionRulePath = strings.ReplaceAll(updateProtectionRulePath, "{id}", d.Id())

		updateProtectionRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateProtectionRuleOpt.JSONBody = utils.RemoveNil(buildUpdateProtectionRuleBodyParams(d))
		_, err = updateProtectionRuleClient.Request("PUT", updateProtectionRulePath, &updateProtectionRuleOpt)
		if err != nil {
			return diag.Errorf("error updating ProtectionRule: %s", err)
		}
	}

	if d.HasChange("sequence") {
		// updateProtectionRuleOrder: Update the order of CFW Protection Rule
		updateProtectionRuleOrderPath := updateProtectionRuleClient.Endpoint + updateProtectionRuleOrderHttpUrl
		updateProtectionRuleOrderPath = strings.ReplaceAll(updateProtectionRuleOrderPath, "{project_id}", updateProtectionRuleClient.ProjectID)
		updateProtectionRuleOrderPath = strings.ReplaceAll(updateProtectionRuleOrderPath, "{id}", d.Id())

		updateProtectionRuleOrderOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateProtectionRuleOrderOpt.JSONBody = utils.RemoveNil(buildUpdateProtectionRuleRequestBodyOrderRuleAclDto(d.Get("sequence")))
		_, err = updateProtectionRuleClient.Request("PUT", updateProtectionRuleOrderPath, &updateProtectionRuleOrderOpt)
		if err != nil {
			return diag.Errorf("error updating protection rule order: %s", err)
		}
	}

	return resourceProtectionRuleRead(ctx, d, meta)
}

func buildUpdateProtectionRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action_type":              d.Get("action_type"),
		"address_type":             d.Get("address_type"),
		"description":              utils.ValueIngoreEmpty(d.Get("description")),
		"direction":                d.Get("direction"),
		"long_connect_enable":      d.Get("long_connect_enable"),
		"long_connect_time_hour":   utils.ValueIngoreEmpty(d.Get("long_connect_time_hour")),
		"long_connect_time_minute": utils.ValueIngoreEmpty(d.Get("long_connect_time_minute")),
		"long_connect_time_second": utils.ValueIngoreEmpty(d.Get("long_connect_time_second")),
		"name":                     utils.ValueIngoreEmpty(d.Get("name")),
		"service":                  buildUpdateProtectionRuleRequestBodyRuleServiceDto(d.Get("service")),
		"source":                   buildUpdateProtectionRuleRequestBodyRuleAddressDto(d.Get("source")),
		"destination":              buildUpdateProtectionRuleRequestBodyRuleAddressDto(d.Get("destination")),
		"status":                   d.Get("status"),
		"type":                     d.Get("type"),
	}
	return bodyParams
}

func buildUpdateProtectionRuleRequestBodyOrderRuleAclDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"dest_rule_id": utils.ValueIngoreEmpty(raw["dest_rule_id"]),
			"top":          raw["top"],
		}
		return params
	}
	return nil
}

func buildUpdateProtectionRuleRequestBodyRuleServiceDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"dest_port":        utils.ValueIngoreEmpty(raw["dest_port"]),
			"protocol":         utils.ValueIngoreEmpty(raw["protocol"]),
			"service_set_id":   utils.ValueIngoreEmpty(raw["service_set_id"]),
			"service_set_name": utils.ValueIngoreEmpty(raw["service_set_name"]),
			"source_port":      utils.ValueIngoreEmpty(raw["source_port"]),
			"type":             raw["type"],
		}
		return params
	}
	return nil
}

func buildUpdateProtectionRuleRequestBodyRuleAddressDto(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"address":             utils.ValueIngoreEmpty(raw["address"]),
			"address_set_id":      utils.ValueIngoreEmpty(raw["address_set_id"]),
			"address_set_name":    utils.ValueIngoreEmpty(raw["address_set_name"]),
			"address_type":        utils.ValueIngoreEmpty(raw["address_type"]),
			"domain_address_name": utils.ValueIngoreEmpty(raw["domain_address_name"]),
			"type":                raw["type"],
		}
		return params
	}
	return nil
}

func resourceProtectionRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteProtectionRule: Delete an existing CFW Protection Rule
	var (
		deleteProtectionRuleHttpUrl = "v1/{project_id}/acl-rule/{id}"
		deleteProtectionRuleProduct = "cfw"
	)
	deleteProtectionRuleClient, err := conf.NewServiceClient(deleteProtectionRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating ProtectionRule Client: %s", err)
	}

	deleteProtectionRulePath := deleteProtectionRuleClient.Endpoint + deleteProtectionRuleHttpUrl
	deleteProtectionRulePath = strings.ReplaceAll(deleteProtectionRulePath, "{project_id}", deleteProtectionRuleClient.ProjectID)
	deleteProtectionRulePath = strings.ReplaceAll(deleteProtectionRulePath, "{id}", d.Id())

	deleteProtectionRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteProtectionRuleClient.Request("DELETE", deleteProtectionRulePath, &deleteProtectionRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting ProtectionRule: %s", err)
	}

	return nil
}

func resourceProtectionRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <object_id>/<id>")
	}

	d.Set("object_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
