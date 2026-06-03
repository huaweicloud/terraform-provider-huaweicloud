package das

import (
	"context"
	"fmt"
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

var sqlLimitRuleNonUpdatableParams = []string{
	"instance_id",
	"database_name",
	"sql_type",
	"pattern",
	"max_concurrency",
	"max_waiting",
	"his_sql_limit_switch",
}

// @API DAS POST /v3/{project_id}/instances/{instance_id}/sql-limit/rules
// @API DAS GET /v3/{project_id}/instances/{instance_id}/sql-limit/rules
// @API DAS DELETE /v3/{project_id}/instances/{instance_id}/sql-limit/rules
func ResourceSqlLimitRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlLimitRuleCreate,
		ReadContext:   resourceSqlLimitRuleRead,
		UpdateContext: resourceSqlLimitRuleUpdate,
		DeleteContext: resourceSqlLimitRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(sqlLimitRuleNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSqlLimitRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the SQL limit rule is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the database instance.`,
			},
			"sql_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The SQL type.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The SQL limit rule pattern.`,
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum concurrency.`,
			},

			// Optional parameters.
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name.`,
			},
			"max_waiting": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum waiting time.`,
			},
			"his_sql_limit_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the historical SQL limit switch.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildCreateSqlLimitRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	rule := map[string]interface{}{
		"sql_type":        d.Get("sql_type").(string),
		"pattern":         d.Get("pattern").(string),
		"max_concurrency": d.Get("max_concurrency").(int),
	}

	if v, ok := d.GetOk("max_waiting"); ok && v.(int) > 0 {
		rule["max_waiting"] = v
	}
	if v, ok := d.GetOk("his_sql_limit_switch"); ok {
		rule["his_sql_limit_switch"] = v
	}

	result := map[string]interface{}{
		// The API only supports MySQL
		"datastore_type":  "MySQL",
		"sql_limit_rules": []interface{}{rule},
	}

	if v, ok := d.GetOk("database_name"); ok {
		result["database_name"] = v
	}

	return result
}

func createSqlLimitRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/sql-limit/rules"
		instanceId = d.Get("instance_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildCreateSqlLimitRuleBodyParams(d),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceSqlLimitRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	err = createSqlLimitRule(client, d)
	if err != nil {
		return diag.Errorf("error creating DAS SQL limit rule: %s", err)
	}

	rules, err := GetSqlLimitRules(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.Errorf("error retrieving DAS SQL limit rules: %s", err)
	}
	if len(rules) == 0 {
		return diag.Errorf("no SQL limit rules found after creation")
	}

	ruleId := utils.PathSearch("id", rules[0], "").(string)
	d.SetId(ruleId)

	return resourceSqlLimitRuleRead(ctx, d, meta)
}

// GetSqlLimitRules queries all SQL limit rules for the instance.
func GetSqlLimitRules(client *golangsdk.ServiceClient, instanceId string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/sql-limit/rules"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	// The API only supports MySQL
	getPath = fmt.Sprintf("%s?datastore_type=%s", getPath, "MySQL")

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

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("sql_limit_rules", respBody, make([]interface{}, 0)).([]interface{})
	return rules, nil
}

func flattenSqlLimitRules(rule interface{}) map[string]interface{} {
	if rule == nil {
		return nil
	}

	return map[string]interface{}{
		"sql_type":             utils.PathSearch("sql_type", rule, nil),
		"pattern":              utils.PathSearch("pattern", rule, nil),
		"max_concurrency":      utils.PathSearch("max_concurrency", rule, nil),
		"max_waiting":          utils.PathSearch("max_waiting", rule, nil),
		"his_sql_limit_switch": utils.PathSearch("his_sql_limit_switch", rule, nil),
	}
}

func resourceSqlLimitRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	rules, err := GetSqlLimitRules(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DAS.200114"),
			"error retrieving DAS SQL limit rule")
	}
	if len(rules) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DAS SQL limit rule")
	}

	ruleId := d.Id()
	rule := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", ruleId), rules, nil)
	if rule == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DAS SQL limit rule")
	}

	ruleResp := flattenSqlLimitRules(rule)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sql_type", ruleResp["sql_type"]),
		d.Set("pattern", ruleResp["pattern"]),
		d.Set("max_concurrency", ruleResp["max_concurrency"]),
		d.Set("his_sql_limit_switch", ruleResp["his_sql_limit_switch"]),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSqlLimitRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func deleteSqlLimitRule(client *golangsdk.ServiceClient, d *schema.ResourceData, ruleId string) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/sql-limit/rules"
		instanceId = d.Get("instance_id").(string)
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	bodyParams := map[string]interface{}{
		// The API only supports MySQL
		"datastore_type":     "MySQL",
		"sql_limit_rule_ids": []string{ruleId},
	}

	if v, ok := d.GetOk("database_name"); ok {
		bodyParams["database_name"] = v
	}

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: bodyParams,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceSqlLimitRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	ruleId := d.Id()

	err = deleteSqlLimitRule(client, d, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DAS SQL limit rule")
	}

	return nil
}

func resourceSqlLimitRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format, expected '<instance_id>/<rule_id>', got '%s'", d.Id())
	}

	instanceId := parts[0]
	ruleId := parts[1]
	d.SetId(ruleId)

	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
