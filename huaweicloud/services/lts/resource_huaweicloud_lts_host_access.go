// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS DELETE /v3/{project_id}/lts/access-config
// @API LTS POST /v3/{project_id}/lts/access-config
// @API LTS PUT /v3/{project_id}/lts/access-config
// @API LTS POST /v3/{project_id}/lts/access-config-list
func ResourceHostAccessConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostAccessConfigCreate,
		UpdateContext: resourceHostAccessConfigUpdate,
		ReadContext:   resourceHostAccessConfigRead,
		DeleteContext: resourceHostAccessConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: hostAccessConfigResourceImportState,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem:     hostAccessConfigDeatilSchema("access_config.0"),
			},
			"host_group_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"processor_type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"processors"},
				Description:  `The type of the ICAgent structuring parsing.`,
			},
			"processors": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"processor_type"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The type of the parser.`,
						},
						"detail": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The configuration of the parser, in JSON format.`,
						},
					},
				},
				Description: `The list of the ICAgent structuring parsing rules.`,
			},
			"demo_log": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The example log of the ICAgent structuring parsing.`,
			},
			"demo_fields": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the parsed field.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the parsed field.`,
						},
					},
				},
				Description: `The list of the parsed fields of the example log`,
			},
			"binary_collect": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to allow collection of binary log files.`,
			},
			"encoding_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encoding format log file.`,
			},
			// If not specified, the API defaults to true.
			"incremental_collect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to collect logs incrementally.`,
			},
			"log_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable log splitting.`,
			},
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the host access, in RFC3339 format.`,
			},
		},
	}
}

func hostAccessConfigDeatilSchema(parent string) *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"black_paths": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"single_log_format": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     hostAccessConfigSingleLogFormatSchema(),
				ExactlyOneOf: []string{
					fmt.Sprintf("%s.multi_log_format", parent),
				},
			},
			"multi_log_format": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem:     hostAccessConfigMultiLogFormatSchema(),
			},
			"windows_log_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     hostAccessConfigWindowsLogInfoSchema(),
				Optional: true,
				Computed: true,
			},
			"custom_key_value": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The custom key/value pairs of the host access.`,
			},
			"system_fields": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of system built-in fields of the host access.`,
			},
			"repeat_collect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to allow repeated flie collection.`,
			},
		},
	}
	return &sc
}

func hostAccessConfigSingleLogFormatSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func hostAccessConfigMultiLogFormatSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func hostAccessConfigWindowsLogInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"categorys": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"event_level": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"time_offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"time_offset_unit": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceHostAccessConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHostAccessConfigHttpUrl = "v3/{project_id}/lts/access-config"
		createHostAccessConfigProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(createHostAccessConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createHostAccessConfigPath := ltsClient.Endpoint + createHostAccessConfigHttpUrl
	createHostAccessConfigPath = strings.ReplaceAll(createHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	createHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	createHostAccessConfigOpt.JSONBody = utils.RemoveNil(buildCreateHostAccessConfigBodyParams(d))
	createHostAccessConfigResp, err := ltsClient.Request("POST", createHostAccessConfigPath, &createHostAccessConfigOpt)
	if err != nil {
		return diag.Errorf("error creating host access config: %s", err)
	}

	createHostAccessConfigRespBody, err := utils.FlattenResponse(createHostAccessConfigResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessId := utils.PathSearch("access_config_id", createHostAccessConfigRespBody, "").(string)
	if accessId == "" {
		return diag.Errorf("unable to find the LTS host access ID from the API response")
	}
	d.SetId(accessId)

	return resourceHostAccessConfigRead(ctx, d, meta)
}

func buildCreateHostAccessConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	logInfoOpts := map[string]interface{}{
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
	}

	bodyParams := map[string]interface{}{
		"access_config_type":   "AGENT",
		"access_config_name":   d.Get("name"),
		"access_config_detail": buildHostAccessConfigDeatilRequestBody(d.Get("access_config")),
		"log_info":             logInfoOpts,
		"host_group_info":      buildHostGroupInfoRequestBody(d),
		"access_config_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"processor_type":       utils.ValueIgnoreEmpty(d.Get("processor_type")),
		"processors":           buildHostAccessProcessors(d.Get("processors").([]interface{})),
		"demo_log":             utils.ValueIgnoreEmpty(d.Get("demo_log")),
		"demo_fields":          buildHostAccessDemoFields(d.Get("demo_fields").(*schema.Set)),
		"binary_collect":       d.Get("binary_collect"),
		"encoding_format":      utils.ValueIgnoreEmpty(d.Get("encoding_format")),
		"incremental_collect":  d.Get("incremental_collect"),
		"log_split":            d.Get("log_split"),
	}
	return bodyParams
}

func buildHostGroupInfoRequestBody(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("host_group_ids"); ok {
		return map[string]interface{}{
			"host_group_id_list": utils.ExpandToStringList(v.([]interface{})),
		}
	}

	return make(map[string]interface{})
}

func buildHostAccessConfigDeatilRequestBody(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"paths":            utils.ValueIgnoreEmpty(raw["paths"].(*schema.Set).List()),
			"black_paths":      utils.ValueIgnoreEmpty(raw["black_paths"].(*schema.Set).List()),
			"format":           buildHostAccessConfigFormatRequestBody(raw),
			"windows_log_info": buildHostAccessConfigWindowsLogInfoRequestBody(raw["windows_log_info"]),
			"custom_key_value": utils.ValueIgnoreEmpty(raw["custom_key_value"]),
			"system_fields":    utils.ValueIgnoreEmpty(raw["system_fields"].(*schema.Set).List()),
			"repeat_collect":   raw["repeat_collect"],
		}
		return params
	}
	return nil
}

func buildSingleLogFormatRequestBody(rawParams map[string]interface{}) map[string]interface{} {
	mode := rawParams["mode"].(string)
	value := rawParams["value"].(string)

	// get the current timestamp if value is not specified in system mode
	if mode == "system" && value == "" {
		value = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	return map[string]interface{}{
		"mode":  mode,
		"value": value,
	}
}

func buildHostAccessConfigFormatRequestBody(rawParams map[string]interface{}) map[string]interface{} {
	log.Printf("[DEBUG] single_log_format: %#v", rawParams["single_log_format"])
	log.Printf("[DEBUG] multi_log_format: %#v", rawParams["multi_log_format"])

	if singleRaw, ok := rawParams["single_log_format"]; ok {
		if rawArray, ok := singleRaw.([]interface{}); ok {
			if len(rawArray) > 0 {
				raw, ok := rawArray[0].(map[string]interface{})
				if ok && len(raw) > 0 {
					return map[string]interface{}{
						"single": buildSingleLogFormatRequestBody(raw),
					}
				}
			}
		}
	}

	if multiRaw, ok := rawParams["multi_log_format"]; ok {
		if rawArray, ok := multiRaw.([]interface{}); ok {
			if len(rawArray) > 0 {
				raw, ok := rawArray[0].(map[string]interface{})
				if ok && len(raw) > 0 {
					return map[string]interface{}{
						"multi": raw,
					}
				}
			}
		}
	}

	return nil
}

func buildHostAccessConfigWindowsLogInfoRequestBody(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		timeOffsetOpts := map[string]interface{}{
			"offset": raw["time_offset"],
			"unit":   raw["time_offset_unit"],
		}
		params := map[string]interface{}{
			"categorys":   raw["categorys"],
			"event_level": raw["event_level"],
			"time_offset": timeOffsetOpts,
		}
		return params
	}
	return nil
}

func buildHostAccessProcessors(processors []interface{}) []map[string]interface{} {
	if len(processors) < 1 || processors[0] == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(processors))
	for _, processor := range processors {
		result = append(result, map[string]interface{}{
			"type":   utils.ValueIgnoreEmpty(utils.PathSearch("type", processor, nil)),
			"detail": utils.StringToJson(utils.PathSearch("detail", processor, "").(string)),
		})
	}

	return result
}

func buildHostAccessDemoFields(demoFields *schema.Set) []map[string]interface{} {
	if demoFields.Len() < 1 {
		return nil
	}
	result := make([]map[string]interface{}, 0, demoFields.Len())
	for _, demoField := range demoFields.List() {
		result = append(result, map[string]interface{}{
			"field_name":  utils.ValueIgnoreEmpty(utils.PathSearch("name", demoField, nil)),
			"field_value": utils.ValueIgnoreEmpty(utils.PathSearch("value", demoField, nil)),
		})
	}

	return result
}

func resourceHostAccessConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listHostAccessConfigHttpUrl = "v3/{project_id}/lts/access-config-list"
		listHostAccessConfigProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(listHostAccessConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listHostAccessConfigPath := ltsClient.Endpoint + listHostAccessConfigHttpUrl
	listHostAccessConfigPath = strings.ReplaceAll(listHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	name := d.Get("name").(string)
	listHostAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}
	listHostAccessConfigResp, err := ltsClient.Request("POST", listHostAccessConfigPath, &listHostAccessConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving host access config")
	}

	listHostAccessConfigRespBody, err := utils.FlattenResponse(listHostAccessConfigResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listHostAccessConfigRespBody = utils.PathSearch(jsonPath, listHostAccessConfigRespBody, nil)
	if listHostAccessConfigRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	// update the resource ID for import scenario
	configID := utils.PathSearch("access_config_id", listHostAccessConfigRespBody, "")
	if configID != "" {
		d.SetId(configID.(string))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("access_config_name", listHostAccessConfigRespBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_info.log_group_id", listHostAccessConfigRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_info.log_stream_id", listHostAccessConfigRespBody, nil)),
		d.Set("host_group_ids", utils.PathSearch("host_group_info.host_group_id_list", listHostAccessConfigRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("access_config_tag", listHostAccessConfigRespBody, nil))),
		d.Set("access_config", flattenHostAccessConfigDetail(listHostAccessConfigRespBody)),
		d.Set("processor_type", utils.PathSearch("processor_type", listHostAccessConfigRespBody, nil)),
		d.Set("demo_log", utils.PathSearch("demo_log", listHostAccessConfigRespBody, nil)),
		d.Set("demo_fields",
			flattenHostAccessDemoFields(utils.PathSearch("demo_fields", listHostAccessConfigRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("binary_collect", utils.PathSearch("binary_collect", listHostAccessConfigRespBody, nil)),
		d.Set("encoding_format", utils.PathSearch("encoding_format", listHostAccessConfigRespBody, nil)),
		d.Set("incremental_collect", utils.PathSearch("incremental_collect", listHostAccessConfigRespBody, nil)),
		d.Set("log_split", utils.PathSearch("log_split", listHostAccessConfigRespBody, nil)),
		// Attributes.
		d.Set("access_type", utils.PathSearch("access_config_type", listHostAccessConfigRespBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_info.log_group_name", listHostAccessConfigRespBody, nil)),
		d.Set("log_stream_name", utils.PathSearch("log_info.log_stream_name", listHostAccessConfigRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", listHostAccessConfigRespBody,
			float64(0)).(float64))/1000, false)))

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHostAccessConfigDetail(resp interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"paths":             utils.PathSearch("access_config_detail.paths", resp, nil),
			"black_paths":       utils.PathSearch("access_config_detail.black_paths", resp, nil),
			"single_log_format": flattenHostAccessConfigLogFormat(utils.PathSearch("access_config_detail.format.single", resp, nil)),
			"multi_log_format":  flattenHostAccessConfigLogFormat(utils.PathSearch("access_config_detail.format.multi", resp, nil)),
			"windows_log_info":  flattenHostAccessConfigWindowsLogInfo(utils.PathSearch("access_config_detail.windows_log_info", resp, nil)),
			"custom_key_value":  utils.PathSearch("access_config_detail.custom_key_value", resp, nil),
			"system_fields":     utils.PathSearch("access_config_detail.system_fields", resp, nil),
			"repeat_collect":    utils.PathSearch("access_config_detail.repeat_collect", resp, nil),
		},
	}
}

func flattenHostAccessConfigLogFormat(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"mode":  utils.PathSearch("mode", resp, nil),
			"value": utils.PathSearch("value", resp, nil),
		},
	}
}

func flattenHostAccessConfigWindowsLogInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"categorys":        utils.PathSearch("categorys", resp, nil),
			"event_level":      utils.PathSearch("event_level", resp, nil),
			"time_offset":      utils.PathSearch("time_offset.offset", resp, nil),
			"time_offset_unit": utils.PathSearch("time_offset.unit", resp, nil),
		},
	}
}

func flattenHostAccessDemoFields(demoFields []interface{}) []map[string]interface{} {
	if len(demoFields) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(demoFields))
	for i, demoField := range demoFields {
		result[i] = map[string]interface{}{
			"name":  utils.PathSearch("field_name", demoField, nil),
			"value": utils.PathSearch("field_value", demoField, nil),
		}
	}
	return result
}

func resourceHostAccessConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateHostAccessConfigChanges := []string{
		"name",
		"access_config",
		"host_group_ids",
		"tags",
		"processor_type",
		"processors",
		"demo_log",
		"demo_fields",
		"encoding_format",
		"incremental_collect",
		"log_split",
	}

	if d.HasChanges(updateHostAccessConfigChanges...) {
		var (
			updateHostAccessConfigHttpUrl = "v3/{project_id}/lts/access-config"
			updateHostAccessConfigProduct = "lts"
		)
		ltsClient, err := cfg.NewServiceClient(updateHostAccessConfigProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS client: %s", err)
		}

		updateHostAccessConfigPath := ltsClient.Endpoint + updateHostAccessConfigHttpUrl
		updateHostAccessConfigPath = strings.ReplaceAll(updateHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

		updateHostAccessConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}

		updateHostAccessConfigOpt.JSONBody = utils.RemoveNil(buildUpdateHostAccessConfigBodyParams(d))
		_, err = ltsClient.Request("PUT", updateHostAccessConfigPath, &updateHostAccessConfigOpt)
		if err != nil {
			return diag.Errorf("error updating host access config: %s", err)
		}
	}
	return resourceHostAccessConfigRead(ctx, d, meta)
}

func buildUpdateHostGroupInfoRequestBody(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("host_group_ids"); ok {
		return map[string]interface{}{
			"host_group_id_list": utils.ExpandToStringList(v.([]interface{})),
		}
	}

	return make(map[string]interface{})
}

func buildUpdateHostAccessConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_config_id":     d.Id(),
		"access_config_name":   d.Get("name"),
		"access_config_detail": buildUpdateHostAccessConfigDeatilRequestBody(d.Get("access_config").([]interface{})),
		"host_group_info":      buildUpdateHostGroupInfoRequestBody(d),
		"access_config_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"processor_type":       utils.ValueIgnoreEmpty(d.Get("processor_type")),
		"processors":           utils.ValueIgnoreEmpty(buildHostAccessProcessors(d.Get("processors").([]interface{}))),
		"demo_log":             utils.ValueIgnoreEmpty(d.Get("demo_log")),
		"demo_fields":          utils.ValueIgnoreEmpty(buildHostAccessDemoFields(d.Get("demo_fields").(*schema.Set))),
		"encoding_format":      utils.ValueIgnoreEmpty(d.Get("encoding_format")),
		"incremental_collect":  d.Get("incremental_collect"),
		"log_split":            d.Get("log_split"),
	}
	return bodyParams
}

func buildUpdateHostAccessConfigDeatilRequestBody(accessConfig []interface{}) map[string]interface{} {
	if len(accessConfig) == 0 {
		return nil
	}

	raw := accessConfig[0].(map[string]interface{})
	return map[string]interface{}{
		"paths":            utils.ValueIgnoreEmpty(raw["paths"].(*schema.Set).List()),
		"black_paths":      utils.ValueIgnoreEmpty(raw["black_paths"].(*schema.Set).List()),
		"format":           buildHostAccessConfigFormatRequestBody(raw),
		"windows_log_info": buildHostAccessConfigWindowsLogInfoRequestBody(raw["windows_log_info"]),
		"repeat_collect":   raw["repeat_collect"],
	}
}

func resourceHostAccessConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHostAccessConfigHttpUrl = "v3/{project_id}/lts/access-config"
		deleteHostAccessConfigProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(deleteHostAccessConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deleteHostAccessConfigPath := ltsClient.Endpoint + deleteHostAccessConfigHttpUrl
	deleteHostAccessConfigPath = strings.ReplaceAll(deleteHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	deleteHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	deleteHostAccessConfigOpt.JSONBody = buildDeleteHostAccessConfigBodyParams(d)
	_, err = ltsClient.Request("DELETE", deleteHostAccessConfigPath, &deleteHostAccessConfigOpt)
	if err != nil {
		return diag.Errorf("error deleting host access config: %s", err)
	}

	return nil
}

func buildDeleteHostAccessConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_config_id_list": []string{d.Id()},
	}
	return bodyParams
}

func hostAccessConfigResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	mErr := multierror.Append(nil,
		d.Set("name", name),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
