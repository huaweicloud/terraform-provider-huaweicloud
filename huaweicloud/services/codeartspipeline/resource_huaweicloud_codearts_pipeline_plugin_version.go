package codeartspipeline

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pluginNonUpdatableParams = []string{
	"plugin_name", "business_type", "business_type_display_name", "runtime_attribution", "version", "is_private",
}

// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/create
// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/create-draft
// @API CodeArtsPipeline GET /v1/{domain_id}/agent-plugin/detail
// @API CodeArtsPipeline GET /v1/{domain_id}/agent-plugin/query-all
// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/edit-draft
// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/publish-draft
// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/update-info
// @API CodeArtsPipeline DELETE /v1/{domain_id}/agent-plugin/delete-draft
// @API CodeArtsPipeline DELETE /v3/{domain_id}/extension/info/delete

// Creating a customized plugin version at the first time will create a plugin automatically,
// and when all versions are deleted, the plugin will be deleted at the same time.
func ResourceCodeArtsPipelinePluginVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelinePluginVersionCreate,
		ReadContext:   resourcePipelinePluginVersionRead,
		UpdateContext: resourcePipelinePluginVersionUpdate,
		DeleteContext: resourcePipelinePluginVersionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePipelinePluginVersionImportStateFunc,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(pluginNonUpdatableParams),
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				o, n := d.GetChange("is_formal")
				if o.(bool) && !n.(bool) {
					return errors.New("unsupport to change the published plugin version into draft")
				}
				if o.(bool) && d.HasChanges("version_description", "execution_info", "input_info") {
					return errors.New("the published plugin version infos are not supported to change")
				}

				return nil
			},
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the plugin name.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the display name.`,
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the service type.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the version.`,
			},
			"execution_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the execution information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inner_execution_info": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the inner execution information.`,
						},
					},
				},
			},
			"runtime_attribution": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the runtime attributes.`,
			},
			"version_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version description.`,
			},
			"icon_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the icon URL.`,
			},
			"business_type_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the display name of service type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the basic plugin description.`,
			},
			"is_private": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the plugin is private.`,
			},
			"input_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the input information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the name.`,
						},
						"default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the default value.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the type.`,
						},
						"layout_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the style information..`,
						},
					},
				},
			},
			"is_formal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the plugin is formal.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"maintainers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the maintenance engineer.`,
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID.`,
			},
			"op_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the operator.`,
			},
			"op_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the operation time.`,
			},
			"plugin_attribution": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the plugin attribution.`,
			},
			"plugin_composition_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the combination extension type.`,
			},
		},
	}
}

func resourcePipelinePluginVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/agent-plugin/create-draft"
	if d.Get("is_formal").(bool) {
		httpUrl = "v1/{domain_id}/agent-plugin/create"
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelinePluginVersionBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline plugin version: %s", err)
	}

	d.SetId(d.Get("plugin_name").(string) + "/" + d.Get("version").(string))

	return resourcePipelinePluginVersionRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelinePluginVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	// some attrs have to send empty value, otherwise, will cause unknown error
	bodyParams := map[string]interface{}{
		"plugin_name":                d.Get("plugin_name"),
		"display_name":               d.Get("display_name"),
		"business_type":              d.Get("business_type"),
		"version":                    d.Get("version"),
		"execution_info":             buildPipelinePluginVersionExecutionInfo(d),
		"input_info":                 buildPipelinePluginVersionInputInfo(d),
		"runtime_attribution":        d.Get("runtime_attribution"),
		"icon_url":                   d.Get("icon_url"),
		"description":                d.Get("description"),
		"version_description":        d.Get("version_description"),
		"business_type_display_name": utils.ValueIgnoreEmpty(d.Get("business_type_display_name")),
		"is_private":                 utils.ValueIgnoreEmpty(d.Get("is_private")),

		// have to send `unique_id` for updating
		"unique_id": utils.ValueIgnoreEmpty(d.Get("unique_id")),
	}

	return bodyParams
}

func buildPipelinePluginVersionExecutionInfo(d *schema.ResourceData) interface{} {
	rawRarams := d.Get("execution_info").([]interface{})
	params := rawRarams[0].(map[string]interface{})
	return map[string]interface{}{
		"inner_execution_info": utils.StringToJson(params["inner_execution_info"].(string)),
	}
}

func buildPipelinePluginVersionInputInfo(d *schema.ResourceData) interface{} {
	rawRarams := d.Get("input_info").([]interface{})
	if len(rawRarams) == 0 {
		return make([]interface{}, 0)
	}

	params := make([]map[string]interface{}, 0, len(rawRarams))
	for _, v := range rawRarams {
		if param, ok := v.(map[string]interface{}); ok {
			rst := map[string]interface{}{
				"name":           utils.ValueIgnoreEmpty(param["name"]),
				"default_value":  param["default_value"],
				"type":           utils.ValueIgnoreEmpty(param["type"]),
				"layout_content": utils.ValueIgnoreEmpty(param["layout_content"]),
			}
			params = append(params, rst)
		}
	}

	return params
}

func resourcePipelinePluginVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/agent-plugin/detail?version={version}&plugin_name={plugin_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getPath = strings.ReplaceAll(getPath, "{version}", d.Get("version").(string))
	getPath = strings.ReplaceAll(getPath, "{plugin_name}", d.Get("plugin_name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DEVPIPE.30011001"),
			"error retrieving CodeArts Pipeline plugin version")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	basicInfos, err := getPipelinePluginBasicInfos(client, d, cfg.DomainID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline plugin")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugin_name", utils.PathSearch("plugin_name", getRespBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", getRespBody, nil)),
		d.Set("op_user", utils.PathSearch("op_user", getRespBody, nil)),
		d.Set("op_time", utils.PathSearch("op_time", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
		d.Set("unique_id", utils.PathSearch("unique_id", getRespBody, nil)),
		d.Set("version_description", utils.PathSearch("version_description", getRespBody, nil)),
		d.Set("is_formal", utils.PathSearch("version_attribution", getRespBody, "").(string) == "formal"),
		d.Set("plugin_composition_type", utils.PathSearch("plugin_composition_type", getRespBody, nil)),
		d.Set("plugin_attribution", utils.PathSearch("plugin_attribution", getRespBody, nil)),
		d.Set("input_info", flattenPipelinePluginVersionInputInfo(getRespBody)),
		d.Set("runtime_attribution", utils.PathSearch("runtime_attribution", getRespBody, nil)),
		d.Set("business_type", utils.PathSearch("business_type", basicInfos, nil)),
		d.Set("icon_url", utils.PathSearch("icon_url", basicInfos, nil)),
		d.Set("business_type_display_name", utils.PathSearch("business_type_display_name", basicInfos, nil)),
		d.Set("description", utils.PathSearch("description", basicInfos, nil)),
		d.Set("maintainers", utils.PathSearch("maintainers", basicInfos, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPipelinePluginBasicInfos(client *golangsdk.ServiceClient, d *schema.ResourceData, domainID string) (interface{}, error) {
	var (
		getHttpUrl = "v1/{domain_id}/agent-plugin/query-all?limit=1&offset=0"
		pluginName = d.Get("plugin_name")
	)
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"plugin_name":   pluginName,
			"business_type": []string{"Build", "Gate", "Deploy", "Test", "Normal"},
		},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	plugin := utils.PathSearch("data[0]", getRespBody, nil)
	if plugin == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "POST",
				URL:       "/v1/{domain_id}/agent-plugin/query-all",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the pipeline plugin (%s) does not exist", pluginName)),
			},
		}
	}

	return plugin, nil
}

func flattenPipelinePluginVersionInputInfo(resp interface{}) []interface{} {
	paramsList, ok := utils.PathSearch("input_info", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, v := range paramsList {
			params := v.(map[string]interface{})
			m := map[string]interface{}{
				"name":           utils.PathSearch("name", params, nil),
				"default_value":  utils.PathSearch("default_value", params, nil),
				"type":           utils.PathSearch("type", params, nil),
				"layout_content": utils.PathSearch("layout_content", params, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourcePipelinePluginVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	if d.HasChanges("display_name", "icon_url", "description") {
		httpUrl := "v1/{domain_id}/agent-plugin/update-info"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{domain_id}", cfg.DomainID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdatePipelinePluginBasicInfosBodyParams(d)),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts Pipeline plugin basic infos: %s", err)
		}
	}

	if d.HasChanges("version_description", "execution_info", "input_info") {
		httpUrl := "v1/{domain_id}/agent-plugin/edit-draft"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{domain_id}", cfg.DomainID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelinePluginVersionBodyParams(d)),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts Pipeline plugin version infos: %s", err)
		}
	}

	// publishing the version after updating plugin version infos, for published version infos can not be updated
	if d.HasChange("is_formal") {
		if err := publishPluginDraft(client, d, cfg.DomainID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePipelinePluginVersionRead(ctx, d, meta)
}

func buildUpdatePipelinePluginBasicInfosBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"plugin_name":                d.Get("plugin_name"),
		"business_type":              d.Get("business_type"),
		"display_name":               d.Get("display_name"),
		"runtime_attribution":        utils.ValueIgnoreEmpty(d.Get("runtime_attribution")),
		"icon_url":                   utils.ValueIgnoreEmpty(d.Get("icon_url")),
		"description":                utils.ValueIgnoreEmpty(d.Get("description")),
		"business_type_display_name": utils.ValueIgnoreEmpty(d.Get("business_type_display_name")),
		"is_private":                 utils.ValueIgnoreEmpty(d.Get("is_private")),
	}

	return bodyParams
}

func publishPluginDraft(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	httpUrl := "v1/{domain_id}/agent-plugin/publish-draft"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", domainId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"plugin_name":        d.Get("plugin_name"),
			"display_name":       d.Get("display_name"),
			"plugin_attribution": "custom",
			"version":            d.Get("version"),
		},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error publishing plugin version: %s", err)
	}

	return nil
}

func resourcePipelinePluginVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/agent-plugin/delete-draft?version={version}&plugin_name={plugin_name}"
	if d.Get("is_formal").(bool) {
		httpUrl = "v3/{domain_id}/extension/info/delete?version={version}&plugin_name={plugin_name}&type=single"
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{version}", d.Get("version").(string))
	deletePath = strings.ReplaceAll(deletePath, "{plugin_name}", d.Get("plugin_name").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DEVPIPE.30011001"),
			"error deleting CodeArts Pipeline plugin version")
	}

	return nil
}

func resourcePipelinePluginVersionImportStateFunc(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<plugin_name>/<version>', but got '%s'", d.Id())
	}

	if err := d.Set("plugin_name", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving plugin name: %s", err)
	}
	if err := d.Set("version", parts[1]); err != nil {
		return nil, fmt.Errorf("error saving version: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
