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
	"github.com/jmespath/go-jmespath"

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
				ForceNew: true,
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
		},
	}
}

func hostAccessConfigDeatilSchema(parent string) *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"paths": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"black_paths": {
				Type:     schema.TypeList,
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

	id, err := jmespath.Search("access_config_id", createHostAccessConfigRespBody)
	if err != nil {
		return diag.Errorf("error creating host access config: ID is not found in API response")
	}
	d.SetId(id.(string))

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
			"paths":            utils.ValueIgnoreEmpty(raw["paths"]),
			"black_paths":      utils.ValueIgnoreEmpty(raw["black_paths"]),
			"format":           buildHostAccessConfigFormatRequestBody(raw),
			"windows_log_info": buildHostAccessConfigWindowsLogInfoRequestBody(raw["windows_log_info"]),
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
		d.Set("access_type", utils.PathSearch("access_config_type", listHostAccessConfigRespBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_info.log_group_id", listHostAccessConfigRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_info.log_stream_id", listHostAccessConfigRespBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_info.log_group_name", listHostAccessConfigRespBody, nil)),
		d.Set("log_stream_name", utils.PathSearch("log_info.log_stream_name", listHostAccessConfigRespBody, nil)),
		d.Set("host_group_ids", utils.PathSearch("host_group_info.host_group_id_list", listHostAccessConfigRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("access_config_tag", listHostAccessConfigRespBody, nil))),
		d.Set("access_config", flattenHostAccessConfigDetail(listHostAccessConfigRespBody)),
	)

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

func resourceHostAccessConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateHostAccessConfigChanges := []string{
		"access_config",
		"host_group_ids",
		"tags",
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
		"access_config_detail": buildHostAccessConfigDeatilRequestBody(d.Get("access_config")),
		"host_group_info":      buildUpdateHostGroupInfoRequestBody(d),
		"access_config_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
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
