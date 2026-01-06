package coc

import (
	"context"
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

var documentExecuteNonUpdatableParams = []string{"document_id", "version", "sys_tags", "sys_tags.*.key",
	"sys_tags.*.value", "target_parameter_name", "targets", "targets.*.key", "targets.*.values", "document_type",
	"description"}

// @API COC POST /v1/documents/{document_id}
// @API COC GET /v1/executions/{execution_id}
func ResourceDocumentExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDocumentExecuteCreate,
		ReadContext:   resourceDocumentExecuteRead,
		UpdateContext: resourceDocumentExecuteUpdate,
		DeleteContext: resourceDocumentExecuteDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(documentExecuteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"document_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     documentExecuteKeyValueSchema(),
			},
			"sys_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     documentExecuteKeyValueSchema(),
			},
			"target_parameter_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"values": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"document_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"execution_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     documentExecuteKeyValueComputedSchema(),
			},
			"document_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"document_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsComputedSchema(),
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func documentExecuteKeyValueSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func documentExecuteKeyValueComputedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDocumentExecuteCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version":               utils.ValueIgnoreEmpty(d.Get("version")),
		"parameters":            buildDocumentExecuteKeyValueCreateOpts(d.Get("parameters")),
		"sys_tags":              buildDocumentExecuteKeyValueCreateOpts(d.Get("sys_tags")),
		"target_parameter_name": utils.ValueIgnoreEmpty(d.Get("target_parameter_name")),
		"targets":               buildDocumentExecuteTargetsValuesCreateOpts(d.Get("targets")),
		"document_type":         utils.ValueIgnoreEmpty(d.Get("document_type")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func buildDocumentExecuteKeyValueCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"key":   utils.ValueIgnoreEmpty(raw["key"]),
				"value": utils.ValueIgnoreEmpty(raw["value"]),
			}
		}
		return params
	}

	return nil
}

func buildDocumentExecuteTargetsValuesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"key":    utils.ValueIgnoreEmpty(raw["key"]),
				"values": parseJson(raw["values"].(string)),
			}
		}
		return params
	}

	return nil
}

func resourceDocumentExecuteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/documents/{document_id}"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	documentID := d.Get("document_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{document_id}", documentID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDocumentExecuteCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing the COC document (%s): %s", documentID, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening COC document execution: %s", err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC document execution ID from the API response")
	}

	d.SetId(id)

	return resourceDocumentExecuteRead(ctx, d, meta)
}

func resourceDocumentExecuteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	getRespBody, err := GetDocumentExecution(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving document execution")
	}
	if utils.PathSearch("data", getRespBody, nil) == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving document execution")
	}

	mErr := multierror.Append(
		nil,
		d.Set("document_id", utils.PathSearch("data.document_id", getRespBody, nil)),
		d.Set("target_parameter_name", utils.PathSearch("data.target_parameter_name", getRespBody, nil)),
		d.Set("targets", flattenCocDocumentExecuteTargets(utils.PathSearch("data.targets", getRespBody, nil))),
		d.Set("description", utils.PathSearch("data.description", getRespBody, nil)),
		d.Set("document_name", utils.PathSearch("data.document_name", getRespBody, nil)),
		d.Set("execution_parameters", flattenCocDocumentExecuteKeyValues(
			utils.PathSearch("data.parameters", getRespBody, nil))),
		d.Set("document_version_id", utils.PathSearch("data.document_version_id", getRespBody, nil)),
		d.Set("version", utils.PathSearch("data.document_version", getRespBody, nil)),
		d.Set("start_time", utils.PathSearch("data.start_time", getRespBody, nil)),
		d.Set("end_time", utils.PathSearch("data.end_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", getRespBody, nil)),
		d.Set("creator", utils.PathSearch("data.creator", getRespBody, nil)),
		d.Set("status", utils.PathSearch("data.status", getRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("data.tags", getRespBody, nil))),
		d.Set("type", utils.PathSearch("data.type", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDocumentExecution(client *golangsdk.ServiceClient, executionID string) (interface{}, error) {
	httpUrl := "v1/executions/{execution_id}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{execution_id}", executionID)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	readDocumentExecutionResp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}
	readDocumentExecutionRespBody, err := utils.FlattenResponse(readDocumentExecutionResp)
	if err != nil {
		return nil, err
	}
	return readDocumentExecutionRespBody, nil
}

func flattenCocDocumentExecuteKeyValues(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"key":   utils.PathSearch("key", raw, nil),
				"value": utils.PathSearch("value", raw, nil),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func flattenCocDocumentExecuteTargets(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"key":    utils.PathSearch("key", raw, nil),
				"values": utils.JsonToString(utils.PathSearch("values", raw, nil)),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func resourceDocumentExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDocumentExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting document execute operation resource is not supported. The document execute operation resource" +
		" is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
