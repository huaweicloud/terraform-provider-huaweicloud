package apig

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

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/orchestrations
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}
func ResourceOrchestrationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrchestrationRuleCreate,
		ReadContext:   resourceOrchestrationRuleRead,
		UpdateContext: resourceOrchestrationRuleUpdate,
		DeleteContext: resourceOrchestrationRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOrchestrationRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the orchestration rule is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the orchestration rule belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the orchestration rule.",
			},
			"strategy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the orchestration rule.",
			},
			"is_preprocessing": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether rule is a preprocessing rule.",
			},
			"mapped_param": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The parameter configuration after orchestration, in JSON format.",
				ValidateFunc: validation.StringIsJSON,
			},
			"map": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
				},
				Description: "The list of orchestration mapping rules, each item should be in JSON format.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the orchestration rule, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the orchestration rule, in RFC3339 format.",
			},
		},
	}
}

func buildOrchestrationRuleMapList(mapList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(mapList))

	for _, val := range mapList {
		if strVal, ok := val.(string); ok && strVal != "" {
			result = append(result, utils.StringToJson(strVal))
		}
	}

	return result
}

func buildOrchestrationRuleModifyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"orchestration_name":         d.Get("name"),
		"orchestration_strategy":     d.Get("strategy"),
		"orchestration_mapped_param": utils.StringToJson(d.Get("mapped_param").(string)),
		"is_preprocessing":           d.Get("is_preprocessing"),
		"orchestration_map":          buildOrchestrationRuleMapList(d.Get("map").([]interface{})),
	}
}

func resourceOrchestrationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations"
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildOrchestrationRuleModifyBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating orchestration rule under dedicated instance (%s): %s", instanceId, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orchestrationId := utils.PathSearch("orchestration_id", respBody, "").(string)
	if orchestrationId == "" {
		return diag.Errorf("unable to find the orchestration rule ID from the API response")
	}
	d.SetId(orchestrationId)

	return resourceOrchestrationRuleRead(ctx, d, meta)
}

func flattenOrchestrationRuleMap(mapList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(mapList))

	for _, val := range mapList {
		result = append(result, utils.JsonToString(val))
	}

	return result
}

func GetOrchestrationRuleById(client *golangsdk.ServiceClient, instanceId, ruleId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}"
	)

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{instance_id}", instanceId)
	queryPath = strings.ReplaceAll(queryPath, "{orchestration_id}", ruleId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceOrchestrationRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		ruleId     = d.Id()
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	rule, err := GetOrchestrationRuleById(client, instanceId, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying orchestration rule (%s)", ruleId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("orchestration_name", rule, nil)),
		d.Set("strategy", utils.PathSearch("orchestration_strategy", rule, nil)),
		d.Set("mapped_param", utils.JsonToString(utils.PathSearch("orchestration_mapped_param", rule, nil))),
		d.Set("is_preprocessing", utils.PathSearch("is_preprocessing", rule, nil)),
		d.Set("map", flattenOrchestrationRuleMap(utils.PathSearch("orchestration_map", rule, make([]interface{}, 0)).([]interface{}))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("orchestration_create_time",
			rule, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("orchestration_update_time",
			rule, "").(string))/1000, false)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving orchestration rule fields: %s", err)
	}
	return nil
}

func resourceOrchestrationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}"
		instanceId      = d.Get("instance_id").(string)
		orchestrationId = d.Id()
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{orchestration_id}", orchestrationId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildOrchestrationRuleModifyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating orchestration rule under dedicated instance (%s): %s", instanceId, err)
	}

	return resourceOrchestrationRuleRead(ctx, d, meta)
}

func resourceOrchestrationRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v2/{project_id}/apigw/instances/{instance_id}/orchestrations/{orchestration_id}"
		instanceId      = d.Get("instance_id").(string)
		orchestrationId = d.Id()
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{orchestration_id}", orchestrationId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting orchestration rule (%s)", orchestrationId))
	}
	return nil
}

func resourceOrchestrationRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
