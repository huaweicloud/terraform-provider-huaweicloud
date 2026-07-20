package dsc

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
	scanRuleNonUpdatableParams = []string{"category"}

	scanRuleNotFoundErrCodes = []string{
		"dsc.10000009", // Current instance does not exist.
	}
)

// @API DSC POST /v1/{project_id}/scan-rules
// @API DSC GET /v1/{project_id}/scan-rules
// @API DSC GET /v1/{project_id}/scan-rules/{rule_id}
// @API DSC PUT /v1/{project_id}/scan-rules/{rule_id}
// @API DSC DELETE /v1/{project_id}/scan-rules/{rule_id}
func ResourceScanRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScanRuleCreate,
		ReadContext:   resourceScanRuleRead,
		UpdateContext: resourceScanRuleUpdate,
		DeleteContext: resourceScanRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(scanRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the scan rule is located.`,
			},

			// Required parameters.
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the scan rule.`,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the scan rule.`,
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The category of the scan rule.`,
			},
			"logic_operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The logic operator of the scan rule.`,
			},
			"match_rate": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The match rate of the scan rule.`,
			},
			"min_match": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The minimum match count of the scan rule.`,
			},
			"content": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        scanRuleContentSchema(),
				Description: `The content list of the scan rule.`,
			},

			// Optional parameters.
			"rule_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the scan rule.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        scanRuleTemplateSchema(),
				Description: `The template list associated with the scan rule.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal attributes.
			"templates_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The template ID associated with the rule.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'templates'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func scanRuleContentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"effective_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The effective mode of the rule content.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application location of the rule content.`,
			},
			"rule_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The content of the rule.`,
			},
		},
	}
}

func scanRuleTemplateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The template ID associated with the rule.`,
			},
			"classification_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The classification ID associated with the rule.`,
			},
			"security_level_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The security level ID associated with the rule.`,
			},

			// Optional parameters.
			"is_used": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Whether the rule is enabled in the template.`,
			},

			// Attributes.
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the rule template.`,
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the template.`,
			},
			"classification_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The classification name associated with the rule.`,
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The color number corresponding to the security level.`,
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the security level.`,
			},
		},
	}
}

func buildScanRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := utils.RemoveNil(map[string]interface{}{
		// Required parameters.
		"rule_name":      d.Get("rule_name"),
		"rule_type":      d.Get("rule_type"),
		"category":       d.Get("category"),
		"logic_operator": d.Get("logic_operator"),
		"match_rate":     d.Get("match_rate"),
		"min_match":      d.Get("min_match"),
		"content":        buildScanRuleContentBodyParams(d.Get("content").([]interface{})),
		// Optional parameters.
		"rule_desc": utils.ValueIgnoreEmpty(d.Get("rule_desc")),
		"templates": buildScanRuleTemplatesBodyParams(d.Get("templates").([]interface{})),
	})

	// When 'templates' field is not specified, the API must receive an empty list.
	if _, ok := params["templates"]; !ok {
		params["templates"] = make([]map[string]interface{}, 0)
	}

	return params
}

func buildScanRuleContentBodyParams(contents []interface{}) []map[string]interface{} {
	if len(contents) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(contents))
	for _, v := range contents {
		rst = append(rst, map[string]interface{}{
			"effective_mode": utils.PathSearch("effective_mode", v, nil),
			"location":       utils.PathSearch("location", v, nil),
			"rule_content":   utils.PathSearch("rule_content", v, nil),
		})
	}
	return rst
}

func buildScanRuleTemplatesBodyParams(templates []interface{}) []map[string]interface{} {
	if len(templates) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(templates))
	for _, v := range templates {
		rst = append(rst, map[string]interface{}{
			"template_id":       utils.PathSearch("template_id", v, nil),
			"classification_id": utils.PathSearch("classification_id", v, nil),
			"security_level_id": utils.PathSearch("security_level_id", v, nil),
			"is_used":           utils.ValueIgnoreEmpty(convertStringToBool(utils.PathSearch("is_used", v, "").(string))),
		})
	}

	return rst
}

func convertStringToBool(str string) interface{} {
	if str == "" {
		return nil
	}

	boolValue, err := strconv.ParseBool(str)
	if err != nil {
		log.Printf("[ERROR] unable to convert string (%s) to boolean: %s", str, err)
		return nil
	}

	return boolValue
}

func createScanRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/scan-rules"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildScanRuleBodyParams(d),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceScanRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		ruleName = d.Get("rule_name").(string)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createScanRule(client, d)
	if err != nil {
		return diag.Errorf("error creating scan rule (%s): %s", ruleName, err)
	}

	respBody, err := listScanRules(client)
	if err != nil {
		return diag.Errorf("error getting the scan rules: %s", err)
	}

	ruleId := utils.PathSearch(fmt.Sprintf("scan_rules_list[?rule_name=='%s']|[0].rule_id", ruleName), respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the ID to the scan rule (%s) from API response", ruleName)
	}

	d.SetId(ruleId)

	err = d.Set("templates_origin", refreshScanRuleTemplateOrigin(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "templates")))
	if err != nil {
		return diag.Errorf("unable to refresh 'templates_origin' values: %s", err)
	}

	return resourceScanRuleRead(ctx, d, meta)
}

func refreshScanRuleTemplateOrigin(scriptTemplates interface{}) []interface{} {
	templates, ok := scriptTemplates.([]interface{})
	if !ok {
		return nil
	}

	rest := make([]interface{}, 0, len(templates))
	for _, v := range templates {
		rest = append(rest, map[string]interface{}{
			"template_id": utils.PathSearch("template_id", v, nil),
		})
	}

	return rest
}

func listScanRules(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v1/{project_id}/scan-rules"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func GetScanRuleById(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/scan-rules/{rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", ruleId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceScanRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ruleId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := GetScanRuleById(client, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", scanRuleNotFoundErrCodes...),
			fmt.Sprintf("error retrieving scan rule (%s)", ruleId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rule_name", utils.PathSearch("rule_name", respBody, nil)),
		d.Set("rule_type", utils.PathSearch("rule_type", respBody, nil)),
		d.Set("category", utils.PathSearch("category", respBody, nil)),
		d.Set("logic_operator", utils.PathSearch("logic_operator", respBody, nil)),
		d.Set("match_rate", utils.PathSearch("match_rate", respBody, nil)),
		d.Set("min_match", utils.PathSearch("min_match", respBody, nil)),
		d.Set("content", flattenScanRuleContent(utils.PathSearch("content", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("rule_desc", utils.PathSearch("rule_desc", respBody, nil)),
		d.Set("templates", flattenScanRuleTemplates(
			utils.PathSearch("templates", respBody, make([]interface{}, 0)).([]interface{}),
			d.Get("templates_origin").([]interface{})),
		),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScanRuleContent(contents []interface{}) []map[string]interface{} {
	if len(contents) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(contents))
	for _, v := range contents {
		rst = append(rst, map[string]interface{}{
			"effective_mode": utils.PathSearch("effective_mode", v, nil),
			"location":       utils.PathSearch("location", v, nil),
			"rule_content":   utils.PathSearch("rule_content", v, nil),
		})
	}

	return rst
}

func flattenScanRuleTemplates(templates []interface{}, templatesOrigin []interface{}) []map[string]interface{} {
	if len(templates) < 1 {
		return nil
	}

	associatedTemplates := make([]map[string]interface{}, 0, len(templates))
	for _, v := range templates {
		associatedTemplates = append(associatedTemplates, map[string]interface{}{
			// Required parameters.
			"template_id":       utils.PathSearch("template_id", v, nil),
			"classification_id": utils.PathSearch("classification_id", v, nil),
			"security_level_id": utils.PathSearch("security_level_id", v, nil),
			"is_used":           strconv.FormatBool(utils.PathSearch("is_used", v, false).(bool)),
			// Attributes.
			"id":                   utils.PathSearch("id", v, nil),
			"classification_name":  utils.PathSearch("classification_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"template_name":        utils.PathSearch("template_name", v, nil),
		})
	}

	return orderScanRuleTemplatesByTemplatesOrigin(associatedTemplates, templatesOrigin)
}

func orderScanRuleTemplatesByTemplatesOrigin(templates []map[string]interface{}, templatesOrigin []interface{}) []map[string]interface{} {
	if len(templatesOrigin) < 1 {
		return templates
	}

	sortedTemplates := make([]map[string]interface{}, 0, len(templates))
	templatesCopy := templates
	for _, templateOrigin := range templatesOrigin {
		templateIdOrigin := utils.PathSearch("template_id", templateOrigin, "").(string)
		for index, template := range templatesCopy {
			if utils.PathSearch("template_id", template, "").(string) != templateIdOrigin {
				continue
			}

			sortedTemplates = append(sortedTemplates, templatesCopy[index])
			templatesCopy = append(templatesCopy[:index], templatesCopy[index+1:]...)
			break
		}
	}

	// Add any remaining unsorted templates to the end of the sorted list.
	sortedTemplates = append(sortedTemplates, templatesCopy...)
	return sortedTemplates
}

func updateScanRule(client *golangsdk.ServiceClient, ruleId string, params map[string]interface{}) error {
	httpUrl := "v1/{project_id}/scan-rules/{rule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", ruleId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceScanRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ruleId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		err = updateScanRule(client, ruleId, buildScanRuleBodyParams(d))
		if err != nil {
			return diag.Errorf("error updating scan rule (%s): %s", ruleId, err)
		}

		err = d.Set("templates_origin", refreshScanRuleTemplateOrigin(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "templates")))
		if err != nil {
			return diag.Errorf("unable to refresh 'templates_origin' values: %s", err)
		}
	}

	return resourceScanRuleRead(ctx, d, meta)
}

func deleteScanRule(client *golangsdk.ServiceClient, ruleId string) error {
	httpUrl := "v1/{project_id}/scan-rules/{rule_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", ruleId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceScanRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteScanRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", scanRuleNotFoundErrCodes...),
			fmt.Sprintf("error deleting scan rule (%v)", d.Get("rule_name")))
	}

	return nil
}
