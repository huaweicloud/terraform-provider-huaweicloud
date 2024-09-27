package secmaster

import (
	"context"
	"encoding/json"
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

const (
	PlaybookRuleNotExistsCode = "SecMaster.20048004"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}
func ResourcePlaybookRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookRuleCreate,
		UpdateContext: resourcePlaybookRuleUpdate,
		ReadContext:   resourcePlaybookRuleRead,
		DeleteContext: resourcePlaybookRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePlaybookRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the workspace to which the playbook rule belongs.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the playbook version ID of the rule.`,
			},
			"expression_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the expression type of the rule.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Elem:        playbookRuleConditionItemSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the conditions of the rule.`,
			},
			"logics": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the logics of the rule.`,
			},
			"cron": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the cron expression.`,
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the schedule type.`,
			},
			"start_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook start type.`,
			},
			"end_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook end type.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook end time.`,
			},
			"repeat_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the repeat range.`,
			},
			"only_once": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the repeat range.`,
			},
			"execution_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the execution type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the playbook version.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updated time of the playbook version.`,
			},
		},
	}
}

func playbookRuleConditionItemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the condition name.`,
			},
			"detail": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the condition detail.`,
			},
			"data": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the condition data.`,
			},
		},
	}
	return &sc
}

func resourcePlaybookRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPlaybookRule: Create a SecMaster playbook.
	var (
		createPlaybookRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules"
		createPlaybookRuleProduct = "secmaster"
	)
	createPlaybookRuleClient, err := cfg.NewServiceClient(createPlaybookRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPlaybookRulePath := createPlaybookRuleClient.Endpoint + createPlaybookRuleHttpUrl
	createPlaybookRulePath = strings.ReplaceAll(createPlaybookRulePath, "{project_id}", createPlaybookRuleClient.ProjectID)
	createPlaybookRulePath = strings.ReplaceAll(createPlaybookRulePath, "{workspace_id}", d.Get("workspace_id").(string))
	createPlaybookRulePath = strings.ReplaceAll(createPlaybookRulePath, "{version_id}", d.Get("version_id").(string))

	createPlaybookRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	ruleJson, err := json.Marshal(buildPlaybookRuleBodyParams(d))
	if err != nil {
		return diag.Errorf("error converting rule object to json")
	}
	createPlaybookRuleOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"rule": string(ruleJson),
	})
	createPlaybookRuleResp, err := createPlaybookRuleClient.Request("POST", createPlaybookRulePath, &createPlaybookRuleOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster playbook rule: %s", err)
	}

	createPlaybookRuleRespBody, err := utils.FlattenResponse(createPlaybookRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createPlaybookRuleRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster playbook rule: ID is not found in API response")
	}
	d.SetId(id)

	return resourcePlaybookRuleRead(ctx, d, meta)
}

func buildPlaybookRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expression_type": utils.ValueIgnoreEmpty(d.Get("expression_type")),
		"conditions":      buildPlaybookRuleRequestBodyConditionItem(d.Get("conditions")),
		"logics":          utils.ValueIgnoreEmpty(d.Get("logics")),
		"cron":            utils.ValueIgnoreEmpty(d.Get("cron")),
		"schedule_type":   utils.ValueIgnoreEmpty(d.Get("schedule_type")),
		"start_type":      utils.ValueIgnoreEmpty(d.Get("start_type")),
		"end_type":        utils.ValueIgnoreEmpty(d.Get("end_type")),
		"end_time":        utils.ValueIgnoreEmpty(d.Get("end_time")),
		"repeat_range":    utils.ValueIgnoreEmpty(d.Get("repeat_range")),
		"only_once":       utils.ValueIgnoreEmpty(d.Get("only_once")),
		"execution_type":  utils.ValueIgnoreEmpty(d.Get("execution_type")),
	}
	return bodyParams
}

func buildPlaybookRuleRequestBodyConditionItem(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":   utils.ValueIgnoreEmpty(raw["name"]),
				"detail": utils.ValueIgnoreEmpty(raw["detail"]),
				"data":   utils.ValueIgnoreEmpty(raw["data"]),
			}
		}
		return rst
	}
	return nil
}

func resourcePlaybookRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPlaybookRule: Query the SecMaster playbook detail
	var (
		getPlaybookRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}"
		getPlaybookRuleProduct = "secmaster"
	)
	getPlaybookRuleClient, err := cfg.NewServiceClient(getPlaybookRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookRulePath := getPlaybookRuleClient.Endpoint + getPlaybookRuleHttpUrl
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{project_id}", getPlaybookRuleClient.ProjectID)
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{workspace_id}", d.Get("workspace_id").(string))
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{version_id}", d.Get("version_id").(string))
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{id}", d.Id())

	getPlaybookRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookRuleResp, err := getPlaybookRuleClient.Request("GET", getPlaybookRulePath, &getPlaybookRuleOpt)
	if err != nil {
		// SecMaster.20010001: the workspace ID not found
		// SecMaster.20048004：the playbook rule ID not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", PlaybookRuleNotExistsCode)
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster playbook rule")
	}

	getPlaybookRuleRespBody, err := utils.FlattenResponse(getPlaybookRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleJson := utils.PathSearch("data.rule", getPlaybookRuleRespBody, "")
	var rule interface{}
	err = json.Unmarshal([]byte(ruleJson.(string)), &rule)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("expression_type", utils.PathSearch("expression_type", rule, nil)),
		d.Set("conditions", flattenGetPlaybookRuleResponseBodyConditionItem(rule)),
		d.Set("logics", utils.PathSearch("logics", rule, nil)),
		d.Set("cron", utils.PathSearch("cron", rule, nil)),
		d.Set("schedule_type", utils.PathSearch("schedule_type", rule, nil)),
		d.Set("start_type", utils.PathSearch("start_type", rule, nil)),
		d.Set("end_type", utils.PathSearch("end_type", rule, nil)),
		d.Set("end_time", utils.PathSearch("end_time", rule, nil)),
		d.Set("repeat_range", utils.PathSearch("repeat_range", rule, nil)),
		d.Set("only_once", utils.PathSearch("only_once", rule, nil)),
		d.Set("execution_type", utils.PathSearch("execution_type", rule, nil)),
		d.Set("created_at", utils.PathSearch("data.create_time", getPlaybookRuleRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("data.update_time", getPlaybookRuleRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetPlaybookRuleResponseBodyConditionItem(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":   utils.PathSearch("name", v, nil),
			"detail": utils.PathSearch("detail", v, nil),
			"data":   utils.PathSearch("data", v, nil),
		})
	}
	return rst
}

func resourcePlaybookRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePlaybookRuleChanges := []string{
		"expression_type",
		"conditions",
		"logics",
		"cron",
		"schedule_type",
		"start_type",
		"end_type",
		"end_time",
		"repeat_range",
		"only_once",
		"execution_type",
	}

	if d.HasChanges(updatePlaybookRuleChanges...) {
		// updatePlaybookRule: Update the configuration of SecMaster playbook
		var (
			updatePlaybookRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}"
			updatePlaybookRuleProduct = "secmaster"
		)
		updatePlaybookRuleClient, err := cfg.NewServiceClient(updatePlaybookRuleProduct, region)
		if err != nil {
			return diag.Errorf("error creating SecMaster client: %s", err)
		}

		updatePlaybookRulePath := updatePlaybookRuleClient.Endpoint + updatePlaybookRuleHttpUrl
		updatePlaybookRulePath = strings.ReplaceAll(updatePlaybookRulePath, "{project_id}", updatePlaybookRuleClient.ProjectID)
		updatePlaybookRulePath = strings.ReplaceAll(updatePlaybookRulePath, "{workspace_id}", d.Get("workspace_id").(string))
		updatePlaybookRulePath = strings.ReplaceAll(updatePlaybookRulePath, "{version_id}", d.Get("version_id").(string))
		updatePlaybookRulePath = strings.ReplaceAll(updatePlaybookRulePath, "{id}", d.Id())

		updatePlaybookRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		ruleJson, err := json.Marshal(buildPlaybookRuleBodyParams(d))
		if err != nil {
			return diag.Errorf("error converting rule object to json sting")
		}
		updatePlaybookRuleOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
			"rule": string(ruleJson),
		})
		_, err = updatePlaybookRuleClient.Request("PUT", updatePlaybookRulePath, &updatePlaybookRuleOpt)
		if err != nil {
			return diag.Errorf("error updating SecMaster playbook rule: %s", err)
		}
	}
	return resourcePlaybookRuleRead(ctx, d, meta)
}

func resourcePlaybookRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePlaybookRule: Delete an existing SecMaster playbook
	var (
		deletePlaybookRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}"
		deletePlaybookRuleProduct = "secmaster"
	)
	deletePlaybookRuleClient, err := cfg.NewServiceClient(deletePlaybookRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePlaybookRulePath := deletePlaybookRuleClient.Endpoint + deletePlaybookRuleHttpUrl
	deletePlaybookRulePath = strings.ReplaceAll(deletePlaybookRulePath, "{project_id}", deletePlaybookRuleClient.ProjectID)
	deletePlaybookRulePath = strings.ReplaceAll(deletePlaybookRulePath, "{workspace_id}", d.Get("workspace_id").(string))
	deletePlaybookRulePath = strings.ReplaceAll(deletePlaybookRulePath, "{version_id}", d.Get("version_id").(string))
	deletePlaybookRulePath = strings.ReplaceAll(deletePlaybookRulePath, "{id}", d.Id())

	deletePlaybookRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deletePlaybookRuleClient.Request("DELETE", deletePlaybookRulePath, &deletePlaybookRuleOpt)
	if err != nil {
		// SecMaster.20010001: the workspace ID not found
		// SecMaster.20048004：the playbook rule ID not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", PlaybookRuleNotExistsCode)
		return common.CheckDeletedDiag(d, err, "error deleting SecMaster playbook rule")
	}

	return nil
}

func resourcePlaybookRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_version_id>/<playbook_rule_id>")
	}

	d.SetId(parts[2])

	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
		d.Set("version_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
