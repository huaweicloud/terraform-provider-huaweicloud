package lts

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var streamIndexConfigNonUpdatableParams = []string{"group_id", "stream_id"}

// @API LTS POST /v1.0/{project_id}/groups/{group_id}/stream/{stream_id}/index/config
// @API LTS GET /v1.0/{project_id}/groups/{group_id}/stream/{stream_id}/index/config
func ResourceStreamIndexConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resoureStreamIndexConfigurationCreate,
		ReadContext:   resoureStreamIndexConfigurationRead,
		UpdateContext: resoureStreamIndexConfigurationUpdate,
		DeleteContext: resoureStreamIndexConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: streamIndexConfigurationResourceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(streamIndexConfigNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group to which the index configuration belongs.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream to which the index configuration belongs.`,
			},
			"full_text_index": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: `The full-text index configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tokenizer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The custom delimiter.`,
						},
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Whether to enable the full-text index.`,
						},
						"case_sensitive": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether letters are case-sensitive.`,
						},
						"include_chinese": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Whether to include Chinese.`,
						},
						"ascii": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `The list of the ASCII delimiters.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of the index fields.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the field.`,
						},
						"field_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the field.`,
						},
						"tokenizer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The custom delimiter.`,
						},
						"field_analysis_alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The alias name of the field.`,
						},
						"quick_analysis": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to enable quick analysis.`,
						},
						"case_sensitive": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether letters are case sensitive.`,
						},
						"include_chinese": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to include Chinese.`,
						},
						"ascii": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `The list of the ASCII delimiters.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"lts_sub_fields_info_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The list of of the JSON fields.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The name of the field.`,
									},
									"field_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the field.`,
									},
									"field_analysis_alias": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The alias name of the field.`,
									},
									"quick_analysis": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether to enable quick analysis.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resoureStreamIndexConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		groupId  = d.Get("group_id").(string)
		streamId = d.Get("stream_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	resp, err := streamIndexConfiguration(client, groupId, buildStreamIndexConfigurationBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating index configuration for log stream (%s) under log group (%s): %s",
			streamId, groupId, err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", groupId, streamId))

	return resoureStreamIndexConfigurationRead(ctx, d, meta)
}

func streamIndexConfiguration(client *golangsdk.ServiceClient, groupId string, params map[string]interface{}) (*http.Response, error) {
	httpUrl := "v1.0/{project_id}/groups/{group_id}/stream/{stream_id}/index/config"
	indexConfigPath := client.Endpoint + httpUrl
	indexConfigPath = strings.ReplaceAll(indexConfigPath, "{project_id}", client.ProjectID)
	indexConfigPath = strings.ReplaceAll(indexConfigPath, "{group_id}", groupId)
	indexConfigPath = strings.ReplaceAll(indexConfigPath, "{stream_id}", params["logStreamId"].(string))
	indexConfigOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(params),
	}

	return client.Request("POST", indexConfigPath, &indexConfigOpts)
}

func buildStreamIndexConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"logStreamId":   d.Get("stream_id"),
		"fullTextIndex": buildStreamfullTextIndex(d.Get("full_text_index").([]interface{})),
		"fields":        buildStreamIndexConfigFields(d.Get("fields").([]interface{})),
	}
}

func buildStreamfullTextIndex(fullTextIndexes []interface{}) map[string]interface{} {
	// For API, the `fullTextIndex` and `fullTextIndex.tokenizer` fields is required.
	fullTextIndex := utils.PathSearch("[0]", fullTextIndexes, nil)
	return map[string]interface{}{
		// If tokenizer not specified, defaults to empty string, otherwise the request will fail.
		"tokenizer":      utils.PathSearch("tokenizer", fullTextIndex, ""),
		"enable":         utils.PathSearch("enable", fullTextIndex, nil),
		"caseSensitive":  utils.PathSearch("case_sensitive", fullTextIndex, nil),
		"includeChinese": utils.PathSearch("include_chinese", fullTextIndex, nil),
		"ascii": utils.ValueIgnoreEmpty(utils.PathSearch("ascii", fullTextIndex,
			schema.NewSet(schema.HashString, nil)).(*schema.Set).List()),
	}
}

func buildStreamIndexConfigFields(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(fields))
	for i, v := range fields {
		rest[i] = map[string]interface{}{
			"fieldName":      utils.PathSearch("field_name", v, nil),
			"fieldType":      utils.PathSearch("field_type", v, nil),
			"caseSensitive":  utils.PathSearch("case_sensitive", v, nil),
			"includeChinese": utils.PathSearch("include_chinese", v, nil),
			// The `tokenizer` cannot be ignored, otherwise the request will fail.
			"tokenizer":          utils.PathSearch("tokenizer", v, nil),
			"quickAnalysis":      utils.PathSearch("quick_analysis", v, nil),
			"ascii":              utils.ValueIgnoreEmpty(utils.PathSearch("ascii", v, schema.NewSet(schema.HashString, nil)).(*schema.Set).List()),
			"fieldAnalysisAlias": utils.ValueIgnoreEmpty(utils.PathSearch("field_analysis_alias", v, nil)),
			"ltsSubFieldsInfoList": buildStreamIndexConfigltsSubFields(utils.PathSearch("lts_sub_fields_info_list", v,
				make([]interface{}, 0)).([]interface{})),
		}
	}
	return rest
}

func buildStreamIndexConfigltsSubFields(ltsFields []interface{}) []map[string]interface{} {
	if len(ltsFields) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(ltsFields))
	for i, v := range ltsFields {
		rest[i] = map[string]interface{}{
			"fieldName":          utils.PathSearch("field_name", v, nil),
			"fieldType":          utils.PathSearch("field_type", v, nil),
			"quickAnalysis":      utils.PathSearch("quick_analysis", v, nil),
			"fieldAnalysisAlias": utils.ValueIgnoreEmpty(utils.PathSearch("field_analysis_alias", v, nil)),
		}
	}
	return rest
}

func GetStreamIndexConfiguration(client *golangsdk.ServiceClient, groupId, streamId string) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/groups/{group_id}/stream/{stream_id}/index/config"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
	getPath = strings.ReplaceAll(getPath, "{stream_id}", streamId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func resoureStreamIndexConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		groupId  = d.Get("group_id").(string)
		streamId = d.Get("stream_id").(string)
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	indexConfig, err := GetStreamIndexConfiguration(client, groupId, streamId)
	if err != nil {
		// SVCSTG.ALS.200201: The log group or log stream does not exist.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SVCSTG.ALS.200201"),
			fmt.Sprintf("error retrieving index configuration of the log stream (%s)", streamId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("full_text_index", flattenFullTextIndex(utils.PathSearch("fullTextIndex", indexConfig, nil))),
		d.Set("fields", flattenIndexConfigurationFields(utils.PathSearch("fields", indexConfig, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFullTextIndex(fullTextIndex interface{}) []map[string]interface{} {
	if fullTextIndex == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"tokenizer":       utils.PathSearch("tokenizer", fullTextIndex, nil),
			"enable":          utils.PathSearch("enable", fullTextIndex, nil),
			"case_sensitive":  utils.PathSearch("caseSensitive", fullTextIndex, nil),
			"include_chinese": utils.PathSearch("includeChinese", fullTextIndex, nil),
			"ascii":           utils.PathSearch("ascii", fullTextIndex, nil),
		},
	}
}

func flattenIndexConfigurationFields(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(fields))
	for i, v := range fields {
		rest[i] = map[string]interface{}{
			"field_name":           utils.PathSearch("fieldName", v, nil),
			"field_type":           utils.PathSearch("fieldType", v, nil),
			"case_sensitive":       utils.PathSearch("caseSensitive", v, nil),
			"include_chinese":      utils.PathSearch("includeChinese", v, nil),
			"tokenizer":            utils.PathSearch("tokenizer", v, nil),
			"quick_analysis":       utils.PathSearch("quickAnalysis", v, nil),
			"ascii":                utils.PathSearch("ascii", v, nil),
			"field_analysis_alias": utils.PathSearch("fieldAnalysisAlias", v, nil),
			"lts_sub_fields_info_list": flattenIndexConfigLtsSubFields(utils.PathSearch("ltsSubFieldsInfoList", v,
				make([]interface{}, 0)).([]interface{})),
		}
	}
	return rest
}

func flattenIndexConfigLtsSubFields(ltsSubFields []interface{}) []map[string]interface{} {
	if len(ltsSubFields) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(ltsSubFields))
	for i, v := range ltsSubFields {
		rest[i] = map[string]interface{}{
			"field_name":           utils.PathSearch("fieldName", v, nil),
			"field_type":           utils.PathSearch("fieldType", v, nil),
			"quick_analysis":       utils.PathSearch("quickAnalysis", v, nil),
			"field_analysis_alias": utils.PathSearch("fieldAnalysisAlias", v, nil),
		}
	}
	return rest
}

func resoureStreamIndexConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		groupId  = d.Get("group_id").(string)
		streamId = d.Get("stream_id").(string)
	)
	client, err := cfg.NewServiceClient("lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	_, err = streamIndexConfiguration(client, groupId, buildStreamIndexConfigurationBodyParams(d))
	if err != nil {
		return diag.Errorf("error updating index configuration for log stream (%s) under log group (%s): %s",
			streamId, groupId, err)
	}

	return resoureStreamIndexConfigurationRead(ctx, d, meta)
}

func resoureStreamIndexConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not initialize the currently configured index, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func streamIndexConfigurationResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<group_id>/<stream_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("group_id", parts[0]),
		d.Set("stream_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
