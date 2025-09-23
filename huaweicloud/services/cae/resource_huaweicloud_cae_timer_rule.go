package cae

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE POST /v1/{project_id}/cae/timer-rules
// @API CAE GET /v1/{project_id}/cae/timer-rules
// @API CAE PUT /v1/{project_id}/cae/timer-rules/{timer_rule_id}
// @API CAE DELETE /v1/{project_id}/cae/timer-rules/{timer_rule_id}
func ResourceTimerRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTimerRuleCreate,
		ReadContext:   resourceTimerRuleRead,
		UpdateContext: resourceTimerRuleUpdate,
		DeleteContext: resourceTimerRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTimerRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CAE environment.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the timer rule.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the timer rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The status of the timer rule.`,
			},
			"effective_range": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The effective range of the timer rule.`,
			},
			"effective_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The effective policy of the timer rule.`,
			},
			"cron": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cron expression of the timer rule.`,
			},
			"applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the application.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the application.`,
						},
					},
				},
				Description: `The list of the applications in which the timer rule takes effect.`,
			},
			"components": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the component.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the component.`,
						},
					},
				},
				Description: `The list of the components in which the timer rule takes effect.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the enterprise project to which the timer rule belongs.`,
			},
		},
	}
}

func buildTimerRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "TimerRule",
		"spec": map[string]interface{}{
			"name":             d.Get("name"),
			"type":             d.Get("type"),
			"status":           d.Get("status"),
			"effective_range":  d.Get("effective_range"),
			"effective_policy": d.Get("effective_policy"),
			"cron":             d.Get("cron"),
			"apps":             utils.ValueIgnoreEmpty(buildTimerRuleApplications(d.Get("applications").(*schema.Set))),
			"components":       utils.ValueIgnoreEmpty(buildTimerRuleComponents(d.Get("components").(*schema.Set))),
		},
	}
}

func buildTimerRuleApplications(applications *schema.Set) []map[string]interface{} {
	if applications.Len() == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, applications.Len())
	for i, v := range applications.List() {
		rest[i] = map[string]interface{}{
			"app_id":   utils.PathSearch("id", v, nil),
			"app_name": utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
		}
	}
	return rest
}

func buildTimerRuleComponents(components *schema.Set) []map[string]interface{} {
	if components.Len() == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, components.Len())
	for i, v := range components.List() {
		rest[i] = map[string]interface{}{
			"component_id":   utils.PathSearch("id", v, nil),
			"component_name": utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
		}
	}
	return rest
}

func resourceTimerRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/timer-rules"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildTimerRuleBodyParams(d)),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating the timer rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	timerRuleId := utils.PathSearch("items[0].id", respBody, "").(string)
	if timerRuleId == "" {
		return diag.Errorf("unable to find the timer rule ID from the API response")
	}

	d.SetId(timerRuleId)

	return resourceTimerRuleRead(ctx, d, meta)
}

func GetTimerRuleById(client *golangsdk.ServiceClient, environmentId, timerRuleId, epsId string) (interface{}, error) {
	timerRules, err := getTimerRules(client, environmentId, epsId)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...)
	}

	timerRule := utils.PathSearch(fmt.Sprintf("items[?id=='%s']|[0]", timerRuleId), timerRules, nil)
	if timerRule == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return timerRule, nil
}

func getTimerRules(client *golangsdk.ServiceClient, environmentId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/timer-rules"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	getListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, epsId),
	}
	resp, err := client.Request("GET", listPath, &getListOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceTimerRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	timerRule, err := GetTimerRuleById(client, d.Get("environment_id").(string), d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving the timer rule (%s)", d.Get("name").(string)))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", timerRule, nil)),
		d.Set("type", utils.PathSearch("type", timerRule, nil)),
		d.Set("effective_range", utils.PathSearch("effective_range", timerRule, nil)),
		d.Set("effective_policy", utils.PathSearch("effective_policy", timerRule, nil)),
		d.Set("cron", utils.PathSearch("cron", timerRule, nil)),
		d.Set("components", flattenTimerRuleComponents(utils.PathSearch("components", timerRule, make([]interface{}, 0)).([]interface{}))),
		d.Set("applications", flattenTimerRuleApplications(utils.PathSearch("apps", timerRule, make([]interface{}, 0)).([]interface{}))),
		// After the rule is executed, the value of the `status` parameter will be changed to `off`, so in order to prevent the resource
		// from changing. The `status` parameter is not set.
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTimerRuleComponents(components []interface{}) []map[string]interface{} {
	if len(components) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(components))
	for i, v := range components {
		rest[i] = map[string]interface{}{
			"id":   utils.PathSearch("component_id", v, nil),
			"name": utils.PathSearch("component_name", v, nil),
		}
	}
	return rest
}

func flattenTimerRuleApplications(applications []interface{}) []map[string]interface{} {
	if len(applications) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(applications))
	for i, v := range applications {
		rest[i] = map[string]interface{}{
			"id":   utils.PathSearch("app_id", v, nil),
			"name": utils.PathSearch("app_name", v, nil),
		}
	}
	return rest
}

func resourceTimerRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/timer-rules/{timer_rule_id}"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{timer_rule_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildTimerRuleBodyParams(d)),
		OkCodes:          []int{204},
	}
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating the timer rule (%s): %s", d.Get("name").(string), err)
	}

	return resourceTimerRuleRead(ctx, d, meta)
}

func resourceTimerRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/timer-rules/{timer_rule_id}"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{timer_rule_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...),
			fmt.Sprintf("error deleting the timer rule (%s)", d.Get("name").(string)))
	}
	return nil
}

// Since the ID cannot be found on the console, so we need to import by the timer rule name.
func resourceTimerRuleImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg        = meta.(*config.Config)
		importedId = d.Id()
		parts      = strings.Split(importedId, "/")
	)

	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<name>' or "+
			"'<environment_id>/<name>/<enterprise_project_id>', but got '%s'",
			importedId)
	}

	timerRuleName := parts[1]
	environmentId := parts[0]
	mErr := multierror.Append(
		d.Set("environment_id", environmentId),
		d.Set("name", timerRuleName),
	)

	if len(parts) == 3 {
		mErr = multierror.Append(mErr, d.Set("enterprise_project_id", parts[2]))
	}

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	timerRules, err := getTimerRules(client, environmentId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return nil, fmt.Errorf("error retrieving the timer rules: %s", err)
	}

	timerRuleId := utils.PathSearch(fmt.Sprintf("items[?name=='%s']|[0].id", timerRuleName), timerRules, "").(string)
	if timerRuleId == "" {
		return []*schema.ResourceData{d},
			fmt.Errorf("unable to find the ID of the timer rule (%s) from API response : %s", timerRuleName, err)
	}

	d.SetId(timerRuleId)
	return []*schema.ResourceData{d}, nil
}
