package cdn

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var domainRuleNonUpdatableParams = []string{"name"}

// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/rule
// @API CDN PUT /v1.0/cdn/configuration/domains/{domain_name}/rules/full-update

// Please apply for whitelist permission before using this resource.
// The API used by this resource is an offline document and has no official website. So no documentation is provided yet.
func ResourceDomainRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainRuleCreate,
		ReadContext:   resourceDomainRuleRead,
		UpdateContext: resourceDomainRuleUpdate,
		DeleteContext: resourceDomainRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainRuleImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(domainRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain name.",
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The rule name. The valid length is limit from `1` to `50`.",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The rule status. Valid values are **on** and **off**.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The rule priority. The valid value is limit from 1 to 100.`,
						},
						"conditions": ruleConditions(),
						"actions":    ruleActions(),
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "The list of rules.",
			},

			// Internal parameter
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func ruleConditionsMatch() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"logic": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the logical operator. Valid values are **and** and **or**.",
				},
				"criteria": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsJSON,
					DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
						return utils.JSONStringsEqual(o, n)
					},
					Description: "Specifies the match criteria list in JSON format.",
				},
			},
		},
		Description: "Specifies the match configuration.",
	}
}

func ruleConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"match": ruleConditionsMatch(),
			},
		},
	}
}

func ruleActionsHttpResponseHeader() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action": {
					Type:     schema.TypeString,
					Required: true,
					Description: `Specifies the operation type of setting HTTP response header.
Valid values are **set** and **delete**.`,
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Set HTTP response header parameters.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Set the value of HTTP response header parameters.",
				},
			},
		},
	}
}

func ruleActionsAccessControl() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the access control type. Valid values are **block** and **trust**.",
				},
			},
		},
	}
}

func ruleActionsRequestUrlRewrite() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"redirect_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the redirect URL.",
				},
				"execution_mode": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the execution mode. Valid values are **redirect** and **break**.",
				},
				"redirect_status_code": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the redirect status code. Valid values are `301`, `302`, `303`, and `307`.",
				},
				"redirect_host": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the redirect host.",
				},
			},
		},
	}
}

func ruleActionsCacheRule() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ttl": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Specifies the cache expiration time.",
				},
				"ttl_unit": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the cache expiration time unit. Valid values: **s**, **m**, **h**, and **d**",
				},
				"follow_origin": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the cache expiration time source. Valid values: **off**, **on**, and **min_ttl**",
				},
				"force_cache": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Specifies whether to enable forced caching. Valid values are **on** and **off**.",
				},
			},
		},
	}
}

func ruleActionsOriginRequestUrlRewrite() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rewrite_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the rewrite type. Valid values are **simple** and **wildcard**.",
				},
				"target_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the target URL.",
				},
				"source_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the URL to be rewritten back to the source.",
				},
			},
		},
	}
}

func ruleActionsFlexibleOrigin() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"priority": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Specifies the origin priority. Valid value ranges from `1` to `100`.",
				},
				"weight": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Specifies the weight. Valid value ranges from `1` to `100`.",
				},
				"sources_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the source type. Valid values are: **ipaddr**, **domain**, and **obs_bucket**.",
				},
				"ip_or_domain": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the origin IP or domain name.",
				},
				"origin_protocol": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the origin protocol.",
				},
				"obs_bucket_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the OBS bucket type. Valid values are **private** and **public**.",
				},
				"http_port": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the HTTP port. Ranges from `1` to `65,535`. Defaults to `80`.",
				},
				"https_port": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the HTTPS port. Ranges from `1` to `65,535`. Defaults to `443`.",
				},
				"host_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Specifies the host name.",
				},
			},
		},
	}
}

func ruleActionsOriginRequestHeader() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"action": {
					Type:     schema.TypeString,
					Required: true,
					Description: `Specifies the back-to-origin request header setting type.
Valid values are **delete** and **set**`,
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Set back-to-origin request header parameters.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Set the value of the return-to-origin request header parameter.",
				},
			},
		},
	}
}

func ruleActions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Description: "Specifies a list of actions to be performed when the rules are met",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"http_response_header":       ruleActionsHttpResponseHeader(),
				"access_control":             ruleActionsAccessControl(),
				"request_url_rewrite":        ruleActionsRequestUrlRewrite(),
				"cache_rule":                 ruleActionsCacheRule(),
				"origin_request_url_rewrite": ruleActionsOriginRequestUrlRewrite(),
				"flexible_origin":            ruleActionsFlexibleOrigin(),
				"origin_request_header":      ruleActionsOriginRequestHeader(),
			},
		},
	}
}

func buildConditionsMatchBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"logic":    rawMap["logic"],
		"criteria": utils.StringToJsonArray(rawMap["criteria"].(string)),
	}
}

func buildDomainRuleConditionsBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"match": buildConditionsMatchBodyParams(rawMap["match"].([]interface{})),
	}
}

func buildActionsHttpResponseHeaderBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"action": rawMap["action"],
			"name":   rawMap["name"],
			"value":  utils.ValueIgnoreEmpty(rawMap["value"]),
		}
	}

	return rst
}

func buildActionsAccessControlBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"type": rawMap["type"],
	}
}

func buildActionsRequestUrlRewriteBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"redirect_url":         rawMap["redirect_url"],
		"execution_mode":       rawMap["execution_mode"],
		"redirect_status_code": utils.ValueIgnoreEmpty(rawMap["redirect_status_code"]),
		"redirect_host":        utils.ValueIgnoreEmpty(rawMap["redirect_host"]),
	}
}

func buildActionsCacheRuleBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"ttl":           utils.ValueIgnoreEmpty(rawMap["ttl"]),
		"ttl_unit":      rawMap["ttl_unit"],
		"follow_origin": utils.ValueIgnoreEmpty(rawMap["follow_origin"]),
		"force_cache":   utils.ValueIgnoreEmpty(rawMap["force_cache"]),
	}
}

func buildActionsOriginRequestUrlRewriteBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"rewrite_type": rawMap["rewrite_type"],
		"source_url":   utils.ValueIgnoreEmpty(rawMap["source_url"]),
		"target_url":   rawMap["target_url"],
	}
}

func buildActionsFlexibleOriginBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"priority":        rawMap["priority"],
			"weight":          rawMap["weight"],
			"sources_type":    rawMap["sources_type"],
			"ip_or_domain":    rawMap["ip_or_domain"],
			"obs_bucket_type": utils.ValueIgnoreEmpty(rawMap["obs_bucket_type"]),
			"http_port":       utils.ValueIgnoreEmpty(rawMap["http_port"]),
			"https_port":      utils.ValueIgnoreEmpty(rawMap["https_port"]),
			"origin_protocol": rawMap["origin_protocol"],
			"host_name":       utils.ValueIgnoreEmpty(rawMap["host_name"]),
		}
	}

	return rst
}

func buildActionsOriginRequestHeaderBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"action": rawMap["action"],
			"name":   rawMap["name"],
			"value":  utils.ValueIgnoreEmpty(rawMap["value"]),
		}
	}

	return rst
}

func buildDomainRuleActionsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"http_response_header":       buildActionsHttpResponseHeaderBodyParams(rawMap["http_response_header"].(*schema.Set).List()),
			"access_control":             buildActionsAccessControlBodyParams(rawMap["access_control"].([]interface{})),
			"request_url_rewrite":        buildActionsRequestUrlRewriteBodyParams(rawMap["request_url_rewrite"].([]interface{})),
			"cache_rule":                 buildActionsCacheRuleBodyParams(rawMap["cache_rule"].([]interface{})),
			"origin_request_url_rewrite": buildActionsOriginRequestUrlRewriteBodyParams(rawMap["origin_request_url_rewrite"].([]interface{})),
			"flexible_origin":            buildActionsFlexibleOriginBodyParams(rawMap["flexible_origin"].(*schema.Set).List()),
			"origin_request_header":      buildActionsOriginRequestHeaderBodyParams(rawMap["origin_request_header"].(*schema.Set).List()),
		}
	}

	return rst
}

// API restrictions must be written to an uuid value, otherwise an error will be reported.
// The server recommends using 32-bit numbers or letters.
func buildUpdateCdnDomainRuleID() string {
	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		log.Printf("[ERROR] error generating uuid: %s", err)
	}

	return strings.ReplaceAll(generateUUID, "-", "")
}

func buildUpdateCdnDomainRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	ruleArray := d.Get("rules").(*schema.Set).List()
	rstArray := make([]map[string]interface{}, len(ruleArray))
	for i, v := range ruleArray {
		rawMap := v.(map[string]interface{})
		rstArray[i] = map[string]interface{}{
			"name":       rawMap["name"],
			"status":     rawMap["status"],
			"priority":   rawMap["priority"],
			"conditions": buildDomainRuleConditionsBodyParams(rawMap["conditions"].([]interface{})),
			"actions":    buildDomainRuleActionsBodyParams(rawMap["actions"].(*schema.Set).List()),
			"rule_id":    buildUpdateCdnDomainRuleID(),
		}
	}

	rst := map[string]interface{}{
		"rules": rstArray,
	}

	return utils.RemoveNil(rst)
}

func updateCdnDomainRule(client *golangsdk.ServiceClient, d *schema.ResourceData, jsonBody map[string]interface{}) (interface{}, error) {
	requestPath := client.Endpoint + "v1.0/cdn/configuration/domains/{domain_name}/rules/full-update"
	requestPath = strings.ReplaceAll(requestPath, "{domain_name}", d.Get("name").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
		OkCodes:          []int{200, 201, 204},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceDomainRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	respBody, err := updateCdnDomainRule(client, d, buildUpdateCdnDomainRuleBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating CDN domain rule: %s", err)
	}

	// When the API reports an error, the response status code is still `200`.
	// The response example at this time is as follows:
	// {"error": {"error_code": "CDN.0105","error_msg": "The acceleration domain name does not exist."}}
	errorCode := utils.PathSearch("error.error_code", respBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", respBody, "").(string)
		return diag.Errorf("error creating CDN domain rule, error_code: %s; error_msg: %s", errorCode, errorMsg)
	}

	d.SetId(domainName)

	if err := waitingForCdnDomainStatusOnline(ctx, client, d.Get("name").(string), cfg.GetEnterpriseProjectID(d),
		d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) to become online in create operation: %s", domainName, err)
	}

	return resourceDomainRuleRead(ctx, d, meta)
}

// When the domain name does not exist, the API response status code is still `200`.
// The response example at this time is as follows:
// {"error": {"error_code": "CDN.0105","error_msg": "The acceleration domain name does not exist."}}
func parseQueryApiErrorMsg(respBody interface{}) (interface{}, error) {
	errorCode := utils.PathSearch("error.error_code", respBody, "").(string)
	if errorCode == "CDN.0105" {
		return nil, golangsdk.ErrDefault404{}
	}

	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", respBody, "").(string)
		return nil, fmt.Errorf("error_code: %s; error_msg: %s", errorCode, errorMsg)
	}

	rawArray := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func QueryCdnDomainRule(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1.0/cdn/configuration/domains/{domain_name}/rule"
	requestPath = strings.ReplaceAll(requestPath, "{domain_name}", domainName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return parseQueryApiErrorMsg(respBody)
}

func flattenConditionsMatchAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rst := map[string]interface{}{
		"logic":    utils.PathSearch("logic", respBody, nil),
		"criteria": utils.JsonToString(utils.PathSearch("criteria", respBody, nil)),
	}

	return []interface{}{rst}
}

func flattenRuleConditionsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rst := map[string]interface{}{
		"match": flattenConditionsMatchAttribute(utils.PathSearch("match", respBody, nil)),
	}

	return []interface{}{rst}
}

func flattenActionsHttpResponseHeaderAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("http_response_header", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"action": utils.PathSearch("action", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"value":  utils.PathSearch("value", v, nil),
		}
	}

	return rst
}

func flattenActionsAccessControlAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("access_control", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rst := map[string]interface{}{
		"type": utils.PathSearch("type", rawMap, nil),
	}

	return []interface{}{rst}
}

func flattenActionsRequestUrlRewriteAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("request_url_rewrite", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rst := map[string]interface{}{
		"redirect_url":         utils.PathSearch("redirect_url", rawMap, nil),
		"execution_mode":       utils.PathSearch("execution_mode", rawMap, nil),
		"redirect_status_code": utils.PathSearch("redirect_status_code", rawMap, nil),
		"redirect_host":        utils.PathSearch("redirect_host", rawMap, nil),
	}

	return []interface{}{rst}
}

func flattenActionsCacheRuleAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("cache_rule", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rst := map[string]interface{}{
		"ttl":           utils.PathSearch("ttl", rawMap, nil),
		"ttl_unit":      utils.PathSearch("ttl_unit", rawMap, nil),
		"follow_origin": utils.PathSearch("follow_origin", rawMap, nil),
		"force_cache":   utils.PathSearch("force_cache", rawMap, nil),
	}

	return []interface{}{rst}
}

func flattenActionsOriginRequestUrlRewriteAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("origin_request_url_rewrite", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rst := map[string]interface{}{
		"rewrite_type": utils.PathSearch("rewrite_type", rawMap, nil),
		"source_url":   utils.PathSearch("source_url", rawMap, nil),
		"target_url":   utils.PathSearch("target_url", rawMap, nil),
	}

	return []interface{}{rst}
}

func flattenActionsFlexibleOriginAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("flexible_origin", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"priority":        utils.PathSearch("priority", v, nil),
			"weight":          utils.PathSearch("weight", v, nil),
			"sources_type":    utils.PathSearch("sources_type", v, nil),
			"ip_or_domain":    utils.PathSearch("ip_or_domain", v, nil),
			"obs_bucket_type": utils.PathSearch("obs_bucket_type", v, nil),
			"http_port":       utils.PathSearch("http_port", v, nil),
			"https_port":      utils.PathSearch("https_port", v, nil),
			"origin_protocol": utils.PathSearch("origin_protocol", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
		}
	}

	return rst
}

func flattenActionsOriginRequestHeaderAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("origin_request_header", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"action": utils.PathSearch("action", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"value":  utils.PathSearch("value", v, nil),
		}
	}

	return rst
}

func flattenRuleActionsAttribute(rawArray []interface{}) []interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"http_response_header":       flattenActionsHttpResponseHeaderAttribute(v),
			"access_control":             flattenActionsAccessControlAttribute(v),
			"request_url_rewrite":        flattenActionsRequestUrlRewriteAttribute(v),
			"cache_rule":                 flattenActionsCacheRuleAttribute(v),
			"origin_request_url_rewrite": flattenActionsOriginRequestUrlRewriteAttribute(v),
			"flexible_origin":            flattenActionsFlexibleOriginAttribute(v),
			"origin_request_header":      flattenActionsOriginRequestHeaderAttribute(v),
		}
	}

	return rst
}

func flattenDomainRuleAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rawArray := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"name":       utils.PathSearch("name", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"priority":   utils.PathSearch("priority", v, nil),
			"conditions": flattenRuleConditionsAttribute(utils.PathSearch("conditions", v, nil)),
			"actions":    flattenRuleActionsAttribute(utils.PathSearch("actions", v, make([]interface{}, 0)).([]interface{})),
			"rule_id":    utils.PathSearch("rule_id", v, nil),
		}
	}

	return rst
}

func resourceDomainRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		domainName = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainRule, err := QueryCdnDomainRule(client, domainName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CDN domain rule")
	}

	mErr := multierror.Append(nil,
		d.Set("name", domainName),
		d.Set("rules", flattenDomainRuleAttribute(domainRule)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDomainRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	respBody, err := updateCdnDomainRule(client, d, buildUpdateCdnDomainRuleBodyParams(d))
	if err != nil {
		return diag.Errorf("error updating CDN domain rule: %s", err)
	}

	// When the API reports an error, the response status code is still `200`.
	// The response example at this time is as follows:
	// {"error": {"error_code": "CDN.0105","error_msg": "The acceleration domain name does not exist."}}
	errorCode := utils.PathSearch("error.error_code", respBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", respBody, "").(string)
		return diag.Errorf("error updating CDN domain rule, error_code: %s; error_msg: %s", errorCode, errorMsg)
	}

	if err := waitingForCdnDomainStatusOnline(ctx, client, d.Get("name").(string), cfg.GetEnterpriseProjectID(d),
		d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) to become online in update operation: %s", d.Id(), err)
	}
	return resourceDomainRuleRead(ctx, d, meta)
}

func resourceDomainRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}
	// The deletion operation is equivalent to configuring the rule to be empty.
	bodyParams := map[string]interface{}{
		"rules": make([]map[string]interface{}, 0),
	}

	respBody, err := updateCdnDomainRule(client, d, bodyParams)
	if err != nil {
		return diag.Errorf("error deleting CDN domain rule: %s", err)
	}

	// When the API reports an error, the response status code is still `200`.
	// The response example at this time is as follows:
	// {"error": {"error_code": "CDN.0105","error_msg": "The acceleration domain name does not exist."}}
	errorCode := utils.PathSearch("error.error_code", respBody, "").(string)
	if errorCode == "CDN.0105" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", respBody, "").(string)
		return diag.Errorf("error deleting CDN domain rule, error_code: %s; error_msg: %s", errorCode, errorMsg)
	}

	if err := waitingForCdnDomainStatusOnline(ctx, client, d.Get("name").(string), cfg.GetEnterpriseProjectID(d),
		d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) to become online in delete operation: %s", d.Id(), err)
	}

	return nil
}

func resourceDomainRuleImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("name", d.Id())
}
