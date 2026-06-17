package secmaster

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

var socPreprocessRuleNonUpdatableParams = []string{"workspace_id", "mapping_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules/search
func ResourceSocPreprocessRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSocPreprocessRuleCreate,
		ReadContext:   resourceSocPreprocessRuleRead,
		UpdateContext: resourceSocPreprocessRuleUpdate,
		DeleteContext: resourceSocPreprocessRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSocPreprocessRuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(socPreprocessRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"preprocess_rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mapper_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mapper_type_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapping_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapper_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mapper_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expression": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSocPreprocessRuleCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"mapping_id":       d.Get("mapping_id"),
		"preprocess_rules": buildSocPreprocessRuleListBodyParams(d.Get("preprocess_rules").(*schema.Set).List()),
	}
}

func buildSocPreprocessRuleListBodyParams(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rule, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":           utils.ValueIgnoreEmpty(rule["name"]),
			"mapper_id":      utils.ValueIgnoreEmpty(rule["mapper_id"]),
			"mapper_type_id": utils.ValueIgnoreEmpty(rule["mapper_type_id"]),
			"action":         utils.ValueIgnoreEmpty(rule["action"]),
			"expression":     utils.ValueIgnoreEmpty(rule["expression"]),
		})
	}

	return result
}

func configSocPreprocessRuleCreate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildSocPreprocessRuleCreateBodyParams(d)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceSocPreprocessRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := configSocPreprocessRuleCreate(client, d); err != nil {
		return diag.Errorf("error creating SecMaster SOC preprocess rule: %s", err)
	}

	d.SetId(d.Get("mapping_id").(string))

	return resourceSocPreprocessRuleRead(ctx, d, meta)
}

func buildSocPreprocessRuleSearchBodyParams(mappingId string, offset int) map[string]interface{} {
	return map[string]interface{}{
		"mapping_id": mappingId,
		"limit":      100,
		"offset":     offset,
	}
}

func GetSocPreprocessRule(client *golangsdk.ServiceClient, workspaceId, mappingId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules/search"
		offset  = 0
		allData = make([]interface{}, 0)
	)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	for {
		requestOpt.JSONBody = buildSocPreprocessRuleSearchBodyParams(mappingId, offset)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		allData = append(allData, data...)
		offset += len(data)
	}

	if len(allData) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return allData, nil
}

func resourceSocPreprocessRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	allData, err := GetSocPreprocessRule(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster SOC preprocess rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("preprocess_rules", flattenSocPreprocessRulesAttr(allData)),
		d.Set("data", flattenSocPreprocessRuleData(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSocPreprocessRulesAttr(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"mapper_id":      utils.PathSearch("mapper_id", v, nil),
			"mapper_type_id": utils.PathSearch("mapper_type_id", v, nil),
			"action":         utils.PathSearch("action", v, nil),
			"expression":     utils.PathSearch("expression", v, nil),
		})
	}

	return rst
}

func flattenSocPreprocessRuleData(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"workspace_id":   utils.PathSearch("workspace_id", v, nil),
			"mapping_id":     utils.PathSearch("mapping_id", v, nil),
			"mapper_id":      utils.PathSearch("mapper_id", v, nil),
			"mapper_type_id": utils.PathSearch("mapper_type_id", v, nil),
			"action":         utils.PathSearch("action", v, nil),
			"expression":     utils.PathSearch("expression", v, nil),
			"creator_id":     utils.PathSearch("creator_id", v, nil),
			"creator_name":   utils.PathSearch("creator_name", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"modifier_id":    utils.PathSearch("modifier_id", v, nil),
			"modifier_name":  utils.PathSearch("modifier_name", v, nil),
		})
	}

	return rst
}

func resourceSocPreprocessRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// This operation is an overriding operation.
	if err := configSocPreprocessRuleCreate(client, d); err != nil {
		return diag.Errorf("error updating SecMaster SOC preprocess rule: %s", err)
	}

	return resourceSocPreprocessRuleRead(ctx, d, meta)
}

func resourceSocPreprocessRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/preprocess-rules"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody: map[string]interface{}{
			"mapping_id":       d.Get("mapping_id"),
			"preprocess_rules": []interface{}{},
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster SOC preprocess rule: %s", err)
	}

	return nil
}

func resourceSocPreprocessRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<id>, but got %s", importId)
	}

	d.SetId(importIdParts[1])
	mErr := multierror.Append(
		d.Set("workspace_id", importIdParts[0]),
		d.Set("mapping_id", importIdParts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
