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

var collectorParserNonUpdatableParams = []string{
	"workspace_id", "title", "description", "parser_id",
	"modules",
	"modules.*.name",
	"modules.*.template_id",
	"modules.*.connection_module_id",
	"modules.*.children",
	"modules.*.children.*.name",
	"modules.*.children.*.template_id",
	"modules.*.children.*.connection_module_id",
	"modules.*.children.*.fields",
	"modules.*.children.*.fields.*.name",
	"modules.*.children.*.fields.*.type",
	"modules.*.children.*.fields.*.value",
	"modules.*.children.*.fields.*.other",
	"modules.*.children.*.fields.*.template_field_id",
	"modules.*.children.*.fields.*.connection_module_id",
	"modules.*.fields",
	"modules.*.fields.*.name",
	"modules.*.fields.*.type",
	"modules.*.fields.*.value",
	"modules.*.fields.*.other",
	"modules.*.fields.*.template_field_id",
	"modules.*.fields.*.connection_module_id",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/{parser_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/{parser_id}
func ResourceCollectorParser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorParserCreate,
		UpdateContext: resourceCollectorParserUpdate,
		ReadContext:   resourceCollectorParserRead,
		DeleteContext: resourceCollectorParserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCollectorParserImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(collectorParserNonUpdatableParams),

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
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parser_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     collectorParserModuleSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"channel_refer_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func collectorParserModuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"module_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorParserChildModuleSchema(),
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorParserFieldSchema(),
			},
		},
	}
}

func collectorParserChildModuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorParserFieldSchema(),
			},
		},
	}
}

func collectorParserFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"other": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_field_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildCollectorParserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"title":       d.Get("title"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"parser_id":   utils.ValueIgnoreEmpty(d.Get("parser_id")),
		"modules":     buildParserModulesBodyParams(d.Get("modules").([]interface{})),
	}

	return bodyParams
}

func buildParserModulesBodyParams(modules []interface{}) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(modules))
	for _, v := range modules {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		module := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"template_id":          utils.ValueIgnoreEmpty(rawMap["template_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
			"children":             buildParserChildModulesBodyParams(rawMap["children"].([]interface{})),
			"fields":               buildParserFieldsBodyParams(rawMap["fields"].([]interface{})),
		}
		rst = append(rst, module)
	}

	return rst
}

func buildParserChildModulesBodyParams(modules []interface{}) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(modules))
	for _, v := range modules {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		module := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"template_id":          utils.ValueIgnoreEmpty(rawMap["template_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
			"fields":               buildParserFieldsBodyParams(rawMap["fields"].([]interface{})),
		}
		rst = append(rst, module)
	}

	return rst
}

func buildParserFieldsBodyParams(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		field := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"type":                 utils.ValueIgnoreEmpty(rawMap["type"]),
			"value":                utils.ValueIgnoreEmpty(rawMap["value"]),
			"other":                utils.ValueIgnoreEmpty(rawMap["other"]),
			"template_field_id":    utils.ValueIgnoreEmpty(rawMap["template_field_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
		}
		rst = append(rst, field)
	}

	return rst
}

func resourceCollectorParserCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", d.Get("workspace_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectorParserBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating collector parser: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	parserId := utils.PathSearch("parser_id", respBody, "").(string)
	if parserId == "" {
		return diag.Errorf("unable to find parser ID from the API response")
	}

	d.SetId(parserId)

	return resourceCollectorParserRead(context.Background(), d, meta)
}

func GetCollectorParserById(client *golangsdk.ServiceClient, workspaceId, parserId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/{parser_id}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{project_id}", client.ProjectID)
	readPath = strings.ReplaceAll(readPath, "{workspace_id}", workspaceId)
	readPath = strings.ReplaceAll(readPath, "{parser_id}", parserId)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCollectorParserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/{parser_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{project_id}", client.ProjectID)
	readPath = strings.ReplaceAll(readPath, "{workspace_id}", workspaceId)
	readPath = strings.ReplaceAll(readPath, "{parser_id}", d.Id())

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving collector parser")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceId),
		d.Set("title", utils.PathSearch("title", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("parser_id", utils.PathSearch("parser_id", respBody, nil)),
		d.Set("modules", flattenParserModules(utils.PathSearch("modules", respBody, nil))),
		d.Set("channel_refer_count", utils.PathSearch("channel_refer_count", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCollectorParserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCollectorParserRead(ctx, d, meta)
}

func flattenParserModules(modulesResp interface{}) []interface{} {
	if modulesResp == nil {
		return nil
	}

	modules, ok := modulesResp.([]interface{})
	if !ok || len(modules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(modules))
	for _, v := range modules {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
			"module_id":            utils.PathSearch("module_id", v, nil),
			"children":             flattenParserChildModules(utils.PathSearch("children", v, nil)),
			"fields":               flattenParserFields(utils.PathSearch("fields", v, nil)),
		})
	}

	return rst
}

func flattenParserChildModules(modulesResp interface{}) []interface{} {
	if modulesResp == nil {
		return nil
	}

	modules, ok := modulesResp.([]interface{})
	if !ok || len(modules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(modules))
	for _, v := range modules {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
			"fields":               flattenParserFields(utils.PathSearch("fields", v, nil)),
		})
	}

	return rst
}

func flattenParserFields(fieldsResp interface{}) []interface{} {
	if fieldsResp == nil {
		return nil
	}

	fields, ok := fieldsResp.([]interface{})
	if !ok || len(fields) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(fields))
	for _, v := range fields {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"type":                 utils.PathSearch("type", v, nil),
			"value":                utils.PathSearch("value", v, nil),
			"other":                utils.PathSearch("other", v, nil),
			"template_field_id":    utils.PathSearch("template_field_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
		})
	}

	return rst
}

func resourceCollectorParserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/logstash/parsers/{parser_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{parser_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting collector parser (%s)", d.Id()))
	}

	return nil
}

func resourceCollectorParserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
