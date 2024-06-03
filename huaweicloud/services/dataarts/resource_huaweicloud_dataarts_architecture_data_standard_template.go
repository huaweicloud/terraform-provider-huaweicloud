// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArts
// ---------------------------------------------------------------

package dataarts

import (
	"context"
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

// @API DataArtsStudio POST /v2/{project_id}/design/standards/templates/action
// @API DataArtsStudio GET /v2/{project_id}/design/standards/templates
// @API DataArtsStudio DELETE /v2/{project_id}/design/standards/templates
func ResourceDataStandardTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataStandardTemplateCreate,
		ReadContext:   resourceDataStandardTemplateRead,
		UpdateContext: resourceDataStandardTemplateUpdate,
		DeleteContext: resourceDataStandardTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureDataStandardTemplateImportState,
		},

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
				ForceNew: true,
			},
			"optional_fields": {
				Type:        schema.TypeSet,
				Elem:        dataStandardTemplateOptionalFieldSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the optional fields of the data standard template to be activated.`,
			},
			"custom_fields": {
				Type:        schema.TypeSet,
				Elem:        dataStandardTemplateCustomFieldSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom fields of the data standard template to be added.`,
			},
		},
	}
}

func dataStandardTemplateOptionalFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the field.`,
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is required.`,
			},
			"searchable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is search supported.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the optional field.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the field.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of the field.`,
			},
		},
	}
	return &sc
}

func dataStandardTemplateCustomFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the field.`,
			},
			"optional_values": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the optional values of the field.`,
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is required.`,
			},
			"searchable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is search supported.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the custom field.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the field.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of the field.`,
			},
		},
	}
	return &sc
}

func resourceDataStandardTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDataStandardTemplate: create DataArts Architecture data standard template
	var (
		createDataStandardTemplateProduct = "dataarts"
	)
	createDataStandardTemplateClient, err := cfg.NewServiceClient(createDataStandardTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getDataStandardTemplateRespBody, err := getDataStandardTemplate(d, createDataStandardTemplateClient)
	if err != nil {
		return diag.FromErr(err)
	}
	hasTemplate := utils.PathSearch("data.value.hasTemplate", getDataStandardTemplateRespBody, false).(bool)
	if hasTemplate {
		return diag.Errorf("a DataArts Architecture data standard template has exist")
	}

	systemDefaultFields := utils.PathSearch("data.value.preFields_system_default", getDataStandardTemplateRespBody,
		make([]interface{}, 0)).([]interface{})
	optionalFields := utils.PathSearch("data.value.preFields_optional", getDataStandardTemplateRespBody,
		make([]interface{}, 0)).([]interface{})

	// init the template
	initDataStandardTemplateRespBody, err := initDataStandardTemplate(d, createDataStandardTemplateClient,
		buildInitDataStandardTemplateBodyParams(d, systemDefaultFields, optionalFields))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("workspace_id").(string))

	// for the optional fields, the value of required, searchable can not be set for the first initialization operation,
	// it will be modified by subsequent initialization operations, so it is necessary to init again to set required,
	// searchable of the optional fields
	rawOptionalFields := d.Get("optional_fields").(*schema.Set).List()
	if len(rawOptionalFields) > 0 {
		fields := utils.PathSearch("data.value", initDataStandardTemplateRespBody, make([]interface{}, 0)).([]interface{})
		templateNameToIdMap := buildTemplateNameToIdMap(fields)
		updateOptionalFields := buildUpdateDataStandardTemplateRequestBodyOptionalField(
			d.Get("optional_fields").(*schema.Set).List(), nil, templateNameToIdMap)
		_, err = initDataStandardTemplate(d, createDataStandardTemplateClient,
			map[string]interface{}{"fields": updateOptionalFields})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDataStandardTemplateRead(ctx, d, meta)
}

func buildInitDataStandardTemplateBodyParams(d *schema.ResourceData, systemDefaultFields,
	optionalFields []interface{}) map[string]interface{} {
	fields := make([]map[string]interface{}, 0)
	fields = append(fields, buildSystemDefaultOrOptionalFields(systemDefaultFields, true)...)
	fields = append(fields, buildSystemDefaultOrOptionalFields(optionalFields, false)...)
	fields = append(fields, buildCreateDataStandardTemplateRequestBodyCustomField(
		d.Get("custom_fields").(*schema.Set).List())...)
	return map[string]interface{}{"fields": fields}
}

func buildSystemDefaultOrOptionalFields(fields []interface{}, isSystemField bool) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		if raw, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"fd_name":    raw["fd_name"],
				"actived":    isSystemField,
				"required":   isSystemField,
				"searchable": isSystemField,
			})
		}
	}
	return rst
}

func buildCreateDataStandardTemplateRequestBodyCustomField(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst = append(rst, map[string]interface{}{
					"fd_name":         raw["fd_name"],
					"actived":         true,
					"optional_values": utils.ValueIgnoreEmpty(raw["optional_values"]),
					"required":        utils.ValueIgnoreEmpty(raw["required"]),
					"searchable":      utils.ValueIgnoreEmpty(raw["searchable"]),
				})
			}
		}
		return rst
	}
	return nil
}

func resourceDataStandardTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDataStandardTemplate: query DataArts Architecture data standard template
	var (
		getDataStandardTemplateProduct = "dataarts"
	)
	getDataStandardTemplateClient, err := cfg.NewServiceClient(getDataStandardTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getDataStandardTemplateRespBody, err := getDataStandardTemplate(d, getDataStandardTemplateClient)
	if err != nil {
		return diag.FromErr(err)
	}

	hasTemplate := utils.PathSearch("data.value.hasTemplate", getDataStandardTemplateRespBody, false).(bool)
	if !hasTemplate {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("optional_fields", flattenGetDataStandardTemplateResponseBodyOptionalField(getDataStandardTemplateRespBody)),
		d.Set("custom_fields", flattenGetDataStandardTemplateResponseBodyCustomField(getDataStandardTemplateRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDataStandardTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateDataStandardTemplate: update DataArts Architecture data standard template
	var (
		updateDataStandardTemplateProduct = "dataarts"
	)
	updateDataStandardTemplateClient, err := cfg.NewServiceClient(updateDataStandardTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	dataStandardTemplateRespBody, err := getDataStandardTemplate(d, updateDataStandardTemplateClient)
	if err != nil {
		return diag.FromErr(err)
	}

	fields := utils.PathSearch("data.value.allFields", dataStandardTemplateRespBody, make([]interface{}, 0)).([]interface{})
	templateNameToIdMap := buildTemplateNameToIdMap(fields)

	oldOptionalRaw, newOptionalRaw := d.GetChange("optional_fields")
	addOptionalRaw := newOptionalRaw.(*schema.Set).Difference(oldOptionalRaw.(*schema.Set))
	deleteOptionalRaw := oldOptionalRaw.(*schema.Set).Difference(newOptionalRaw.(*schema.Set))

	oldCustomRaw, newCustomRaw := d.GetChange("custom_fields")
	addCustomRaw := newCustomRaw.(*schema.Set).Difference(oldCustomRaw.(*schema.Set))
	deleteCustomRaw := oldCustomRaw.(*schema.Set).Difference(newCustomRaw.(*schema.Set))

	updateOptionalFields := buildUpdateDataStandardTemplateRequestBodyOptionalField(addOptionalRaw.List(),
		deleteOptionalRaw.List(), templateNameToIdMap)
	updateCustomFields, deleteCustomFields := buildUpdateDataStandardTemplateRequestBodyCustomField(addCustomRaw.List(),
		deleteCustomRaw.List(), templateNameToIdMap)
	if len(updateOptionalFields)+len(updateCustomFields) > 0 {
		updateTemplateRequestBody := buildUpdateDataStandardTemplateBodyParams(updateOptionalFields, updateCustomFields)
		_, err = initDataStandardTemplate(d, updateDataStandardTemplateClient, updateTemplateRequestBody)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(deleteCustomFields) > 0 {
		err = deleteDataStandardTemplate(d, updateDataStandardTemplateClient, deleteCustomFields, templateNameToIdMap)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDataStandardTemplateRead(ctx, d, meta)
}

func buildUpdateDataStandardTemplateBodyParams(updateOptionalFields, updateCustomFields []interface{}) map[string]interface{} {
	fields := make([]interface{}, 0)
	fields = append(fields, updateOptionalFields...)
	fields = append(fields, updateCustomFields...)
	return map[string]interface{}{"fields": fields}
}

func buildUpdateDataStandardTemplateRequestBodyOptionalField(addOptionalFields, deleteOptionalFields []interface{},
	templateNameToIdMap map[string]string) []interface{} {
	updateOptionalFieldMap := make(map[string]interface{})
	rst := make([]interface{}, 0)
	for _, v := range addOptionalFields {
		if raw, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"fd_name":    raw["fd_name"],
				"actived":    true,
				"required":   raw["required"],
				"searchable": raw["searchable"],
				"id":         templateNameToIdMap[raw["fd_name"].(string)],
			})
			updateOptionalFieldMap[raw["fd_name"].(string)] = v
		}
	}
	for _, v := range deleteOptionalFields {
		if raw, ok := v.(map[string]interface{}); ok {
			if _, ok = updateOptionalFieldMap[raw["fd_name"].(string)]; ok {
				continue
			}
			rst = append(rst, map[string]interface{}{
				"fd_name":    raw["fd_name"],
				"actived":    false,
				"required":   false,
				"searchable": false,
				"id":         templateNameToIdMap[raw["fd_name"].(string)],
			})
		}
	}
	return rst
}

func buildUpdateDataStandardTemplateRequestBodyCustomField(addCustomFields, deleteCustomFields []interface{},
	templateNameToIdMap map[string]string) ([]interface{}, []interface{}) {
	updateCustomFieldMap := make(map[string]interface{})
	updateFields := make([]interface{}, len(addCustomFields))
	for i, v := range addCustomFields {
		raw := v.(map[string]interface{})
		updateFields[i] = map[string]interface{}{
			"fd_name":         raw["fd_name"],
			"actived":         true,
			"optional_values": raw["optional_values"],
			"required":        raw["required"],
			"searchable":      raw["searchable"],
			"id":              templateNameToIdMap[raw["fd_name"].(string)],
		}
		updateCustomFieldMap[raw["fd_name"].(string)] = v
	}
	deleteFields := make([]interface{}, 0)
	for _, v := range deleteCustomFields {
		raw := v.(map[string]interface{})
		if _, ok := updateCustomFieldMap[raw["fd_name"].(string)]; ok {
			continue
		}
		deleteFields = append(deleteFields, v)
	}
	return updateFields, deleteFields
}

func buildTemplateNameToIdMap(fields []interface{}) map[string]string {
	rst := make(map[string]string)
	for _, v := range fields {
		id := utils.PathSearch("id", v, "").(string)
		fdName := utils.PathSearch("fd_name", v, "").(string)
		rst[fdName] = id
	}
	return rst
}

func getDataStandardTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	// getDataStandardTemplate: query DataArts Architecture data standard template
	var (
		getDataStandardTemplateHttpUrl = "v2/{project_id}/design/standards/templates"
	)

	getDataStandardTemplatePath := client.Endpoint + getDataStandardTemplateHttpUrl
	getDataStandardTemplatePath = strings.ReplaceAll(getDataStandardTemplatePath, "{project_id}", client.ProjectID)

	getDataStandardTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	getDataStandardTemplateResp, err := client.Request("GET", getDataStandardTemplatePath, &getDataStandardTemplateOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture data standard template: %s", err)
	}

	getDataStandardTemplateRespBody, err := utils.FlattenResponse(getDataStandardTemplateResp)
	if err != nil {
		return nil, err
	}
	return getDataStandardTemplateRespBody, nil
}

func initDataStandardTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient,
	params map[string]interface{}) (interface{}, error) {
	var (
		createDataStandardTemplateHttpUrl = "v2/{project_id}/design/standards/templates/action?action-id=init"
	)

	initDataStandardTemplatePath := client.Endpoint + createDataStandardTemplateHttpUrl
	initDataStandardTemplatePath = strings.ReplaceAll(initDataStandardTemplatePath, "{project_id}", client.ProjectID)

	initDataStandardTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	initDataStandardTemplateOpt.JSONBody = utils.RemoveNil(params)
	initDataStandardTemplateResp, err := client.Request("POST", initDataStandardTemplatePath, &initDataStandardTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error initing DataArts Architecture data standard template: %s", err)
	}
	initDataStandardTemplateRespBody, err := utils.FlattenResponse(initDataStandardTemplateResp)
	if err != nil {
		return nil, err
	}
	return initDataStandardTemplateRespBody, nil
}

func flattenGetDataStandardTemplateResponseBodyOptionalField(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("data.value.optional", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		if !utils.PathSearch("actived", v, false).(bool) {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"fd_name":    utils.PathSearch("fd_name", v, nil),
			"required":   utils.PathSearch("required", v, nil),
			"searchable": utils.PathSearch("searchable", v, nil),
			"id":         utils.PathSearch("id", v, nil),
			"created_by": utils.PathSearch("create_by", v, nil),
			"updated_by": utils.PathSearch("update_by", v, nil),
			"created_at": utils.PathSearch("create_time", v, nil),
			"updated_at": utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}

func flattenGetDataStandardTemplateResponseBodyCustomField(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("data.value.custom", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"fd_name":         utils.PathSearch("fd_name", v, nil),
			"optional_values": utils.PathSearch("optional_values", v, nil),
			"required":        utils.PathSearch("required", v, nil),
			"searchable":      utils.PathSearch("searchable", v, nil),
			"id":              utils.PathSearch("id", v, nil),
			"created_by":      utils.PathSearch("create_by", v, nil),
			"updated_by":      utils.PathSearch("update_by", v, nil),
			"created_at":      utils.PathSearch("create_time", v, nil),
			"updated_at":      utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}

func resourceDataStandardTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDataStandardTemplateField: delete DataArts Architecture data standard template field
	var (
		deleteDataStandardTemplateFieldProduct = "dataarts"
	)
	deleteDataStandardTemplateFieldClient, err := cfg.NewServiceClient(deleteDataStandardTemplateFieldProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	dataStandardTemplateRespBody, err := getDataStandardTemplate(d, deleteDataStandardTemplateFieldClient)
	if err != nil {
		return diag.FromErr(err)
	}

	allFields := utils.PathSearch("data.value.allFields", dataStandardTemplateRespBody,
		make([]interface{}, 0)).([]interface{})
	err = deleteDataStandardTemplate(d, deleteDataStandardTemplateFieldClient, allFields, make(map[string]string))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteDataStandardTemplate(d *schema.ResourceData, client *golangsdk.ServiceClient, deleteFields []interface{},
	templateNameToIdMap map[string]string) error {
	var (
		deleteDataStandardTemplateFieldHttpUrl = "v2/{project_id}/design/standards/templates"
	)
	deleteDataStandardTemplateFieldPath := client.Endpoint + deleteDataStandardTemplateFieldHttpUrl
	deleteDataStandardTemplateFieldPath = strings.ReplaceAll(deleteDataStandardTemplateFieldPath, "{project_id}",
		client.ProjectID)

	deleteDataStandardTemplateFieldOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	dataStandardTemplateDeleteParams := buildDeleteDataStandardTemplateRequestParam(deleteFields, templateNameToIdMap)
	deleteDataStandardTemplateFieldPath += dataStandardTemplateDeleteParams

	_, err := client.Request("DELETE", deleteDataStandardTemplateFieldPath, &deleteDataStandardTemplateFieldOpt)
	if err != nil {
		return fmt.Errorf("error deleting DataArts Architecture data standard template: %s", err)
	}
	return nil
}

func buildDeleteDataStandardTemplateRequestParam(deleteFields []interface{}, templateNameToIdMap map[string]string) string {
	if len(deleteFields) == 0 {
		return ""
	}

	rst := ""
	for _, v := range deleteFields {
		raw := v.(map[string]interface{})
		if value, ok := templateNameToIdMap[raw["fd_name"].(string)]; ok {
			rst = fmt.Sprintf("%s,%s", rst, value)
		} else {
			rst = fmt.Sprintf("%s,%s", rst, raw["id"].(string))
		}
	}
	return fmt.Sprintf("?ids=%s", rst)
}

func resourceArchitectureDataStandardTemplateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	mErr := multierror.Append(nil,
		d.Set("workspace_id", d.Id()),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
