package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/app-center/app-rules
// @API Workspace GET /v1/{project_id}/app-center/app-rules
// @API Workspace PATCH /v1/{project_id}/app-center/app-rules/{rule_id}
// @API Workspace DELETE /v1/{project_id}/app-center/app-rules/{rule_id}
func ResourceApplicationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationRuleCreate,
		ReadContext:   resourceApplicationRuleRead,
		UpdateContext: resourceApplicationRuleUpdate,
		DeleteContext: resourceApplicationRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application rule is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the application rule.`,
			},
			"detail": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        workspaceApplicationDetailSchema(),
				Description: `The detail of the application rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the application rule.`,
			},
		},
	}
}

func workspaceApplicationDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The scope of the application rule.`,
			},
			"product_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workspaceApplicationProductRule(),
				Description: `The detail of the product rule.`,
			},
			"path_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workspaceApplicationPathRule(),
				Description: `The detail of the path rule.`,
			},
		},
	}
}

func workspaceApplicationProductRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"identify_condition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The identify condition of the product rule.`,
			},
			"publisher": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, o string, n string, _ *schema.ResourceData) bool {
					return o == "*" && n == ""
				},
				Description: `The publisher of the product.`,
			},
			"product_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, o string, n string, _ *schema.ResourceData) bool {
					return o == "*" && n == ""
				},
				Description: `The name of the product.`,
			},
			"process_name": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, o string, n string, _ *schema.ResourceData) bool {
					return o == "*" && n == ""
				},
				Description: `The process name of the product.`,
			},
			"support_os": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     `Windows`,
				Description: `The list of the supported operating system types.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version of the product rule.`,
			},
			"product_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version of the product.`,
			},
		},
	}
}

func workspaceApplicationPathRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The path where the product is installed.`,
			},
		},
	}
}

func buildApplicationPathRule(pathRules []interface{}) map[string]interface{} {
	if len(pathRules) < 1 {
		return nil
	}

	return map[string]interface{}{
		"path": utils.PathSearch("path", pathRules[0], nil),
	}
}

func buildApplicationProductRule(productRules []interface{}) map[string]interface{} {
	if len(productRules) < 1 {
		return nil
	}

	return map[string]interface{}{
		"identify_condition": utils.PathSearch("identify_condition", productRules[0], nil),
		"publisher":          utils.PathSearch("publisher", productRules[0], nil),
		"product_name":       utils.PathSearch("product_name", productRules[0], nil),
		"process_name":       utils.PathSearch("process_name", productRules[0], nil),
		"support_os":         utils.PathSearch("support_os", productRules[0], "Windows"),
		"version":            utils.ValueIgnoreEmpty(utils.PathSearch("version", productRules[0], nil)),
		"product_version":    utils.ValueIgnoreEmpty(utils.PathSearch("product_version", productRules[0], nil)),
	}
}

func buildApplicationRuleDetail(rules []interface{}) map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	return map[string]interface{}{
		"scope":        utils.PathSearch("scope", rules[0], nil),
		"product_rule": buildApplicationProductRule(utils.PathSearch("product_rule", rules[0], make([]interface{}, 0)).([]interface{})),
		"path_rule":    buildApplicationPathRule(utils.PathSearch("path_rule", rules[0], make([]interface{}, 0)).([]interface{})),
	}
}

func buildCreateApplicationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"rule":        buildApplicationRuleDetail(d.Get("detail").([]interface{})),
		"description": d.Get("description"),
	}
}

func resourceApplicationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := "v1/{project_id}/app-center/app-rules"
	createPath = client.Endpoint + createPath
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationRuleBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace application rule: %s", err)
	}
	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleApplicationId := utils.PathSearch("id", respBody, "").(string)
	if ruleApplicationId == "" {
		return diag.Errorf("unable to find application rule ID from the API response")
	}
	d.SetId(ruleApplicationId)

	return resourceApplicationRuleRead(ctx, d, meta)
}

// GetApplicationRuleById is a method is used to get the application rule.
func GetApplicationRuleById(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	applicationRules, err := listApplicationRules(client)
	if err != nil {
		return nil, err
	}

	applicationRule := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", ruleId), applicationRules, nil)
	if applicationRule == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return applicationRule, nil
}

func buildApplicationRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	return res
}

func listApplicationRules(client *golangsdk.ServiceClient, d ...*schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-center/app-rules?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	if len(d) != 0 {
		listPathWithLimit += buildApplicationRulesQueryParams(d[0])
	}

	opt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		applicationRules := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applicationRules...)
		if len(applicationRules) < limit {
			break
		}
		offset += len(applicationRules)
	}

	return result, nil
}

func flattenApplicationPathRule(pathRule interface{}) []interface{} {
	if pathRule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"path": utils.PathSearch("path", pathRule, nil),
		},
	}
}

func flattenApplicationProductRule(productRule interface{}) []interface{} {
	if productRule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"identify_condition": utils.PathSearch("identify_condition", productRule, nil),
			"publisher":          utils.PathSearch("publisher", productRule, nil),
			"product_name":       utils.PathSearch("product_name", productRule, nil),
			"process_name":       utils.PathSearch("process_name", productRule, nil),
			"support_os":         utils.PathSearch("support_os", productRule, "Windows"),
			"version":            utils.PathSearch("version", productRule, nil),
			"product_version":    utils.PathSearch("product_version", productRule, nil),
		},
	}
}

func flattenApplicationRuleDetail(ruleDetail interface{}) []interface{} {
	if ruleDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"scope":        utils.PathSearch("scope", ruleDetail, nil),
			"product_rule": flattenApplicationProductRule(utils.PathSearch("product_rule", ruleDetail, nil)),
			"path_rule":    flattenApplicationPathRule(utils.PathSearch("path_rule", ruleDetail, nil)),
		},
	}
}

func resourceApplicationRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	applicationRule, err := GetApplicationRuleById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying Workspace application rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", applicationRule, nil)),
		d.Set("detail", flattenApplicationRuleDetail(utils.PathSearch("rule", applicationRule, nil))),
		d.Set("description", utils.PathSearch("description", applicationRule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateApplicationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"rule":        buildApplicationRuleDetail(d.Get("detail").([]interface{})),
		"description": d.Get("description"),
	}
}

func updateApplicationRule(client *golangsdk.ServiceClient, applicationRuleId string, params map[string]interface{}) error {
	updatePath := "v1/{project_id}/app-center/app-rules/{rule_id}"
	updatePath = client.Endpoint + updatePath
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", applicationRuleId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}
	_, err := client.Request("PATCH", updatePath, &updateOpt)
	return err
}

func resourceApplicationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = updateApplicationRule(client, d.Id(), buildUpdateApplicationRuleBodyParams(d))
	if err != nil {
		return diag.Errorf("error updating Workspace application rule (%s): %s", d.Id(), err)
	}

	return resourceApplicationRuleRead(ctx, d, meta)
}

func resourceApplicationRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		deletePath        = "v1/{project_id}/app-center/app-rules/{rule_id}"
		applicationRuleId = d.Id()
	)

	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace application client: %s", err)
	}

	deletePath = client.Endpoint + deletePath
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", applicationRuleId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting Workspace application rule (%s)", applicationRuleId))
	}

	return nil
}
