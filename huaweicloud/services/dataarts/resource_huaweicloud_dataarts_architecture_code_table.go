package dataarts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// API: DataArtsStudio GET /v2/{project_id}/design/code-tables
// API: DataArtsStudio POST /v2/{project_id}/design/code-tables
// API: DataArtsStudio DELETE /v2/{project_id}/design/code-tables
// API: DataArtsStudio PUT /v2/{project_id}/design/code-tables/{id}
// API: DataArtsStudio GET /v2/{project_id}/design/code-tables/{id}
// API: DataArtsStudio GET /v2/{project_id}/design/code-tables/{id}/values
// API: DataArtsStudio PUT /v2/{project_id}/design/code-tables/{id}/values
func ResourceArchitectureCodeTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureCodeTableCreate,
		ReadContext:   resourceArchitectureCodeTableRead,
		UpdateContext: resourceArchitectureCodeTableUpdate,
		DeleteContext: resourceArchitectureCodeTableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCodeTableImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the user group is located.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The workspace ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the code table.",
			},
			"code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The code of the code table.",
			},
			"directory_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The directory ID of the code table.",
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of a field.",
						},
						"code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The code of a field.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of a field.",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The values of a field.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of a field.",
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ordinal of a field.",
						},
						"table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the code table.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the field.",
						},
					},
				},
				Description: "The fields information of the code table.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the code table.",
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The directory path of the code table.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the code table.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the code table was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the code table was updated.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the code table.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the code table.",
			},
		},
	}
}

func resourceArchitectureCodeTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createCodeTableHttpUrl = "v2/{project_id}/design/code-tables"
		createCodeTableProduct = "dataarts"
	)
	createCodeTableClient, err := cfg.NewServiceClient(createCodeTableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createCodeTablePath := createCodeTableClient.Endpoint + createCodeTableHttpUrl
	createCodeTablePath = strings.ReplaceAll(createCodeTablePath, "{project_id}", createCodeTableClient.ProjectID)

	createCodeTableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	createCodeTableOpt.JSONBody = utils.RemoveNil(buildCreateCodeTableParams(d))
	createCodeTableResp, err := createCodeTableClient.Request("POST", createCodeTablePath, &createCodeTableOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Architecture code table: %s", err)
	}

	createCodeTableRespBody, err := utils.FlattenResponse(createCodeTableResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.value.id", createCodeTableRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating DataArts Architecture code table: %s is not found in API response", "id")
	}

	d.SetId(id.(string))

	yes := hasFieldValues(d)
	if yes {
		err = insertCodeTableFieldValues(createCodeTableClient, d)
		if err != nil {
			return diag.Errorf("error adding field values to DataArts Architecture code table: %s", err)
		}
	}

	return resourceArchitectureCodeTableRead(ctx, d, meta)
}

func buildCreateCodeTableParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name_ch":           d.Get("name"),
		"name_en":           d.Get("code"),
		"directory_id":      d.Get("directory_id"),
		"code_table_fields": utils.ValueIngoreEmpty(buildCreateCodeTableFieldParams(d)),
		"description":       utils.ValueIngoreEmpty(d.Get("description")),
	}
}

func buildCreateCodeTableFieldParams(d *schema.ResourceData) []map[string]interface{} {
	rawFields := d.Get("fields").([]interface{})
	fields := make([]map[string]interface{}, len(rawFields))
	for i, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		field := make(map[string]interface{})
		field["ordinal"] = i + 1
		field["name_ch"] = fieldMap["name"].(string)
		field["name_en"] = fieldMap["code"].(string)
		field["data_type"] = fieldMap["type"].(string)
		field["description"] = utils.ValueIngoreEmpty(fieldMap["description"].(string))
		fields[i] = field
	}
	return fields
}

func insertCodeTableFieldValues(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	ids, err := getCodeTableFieldID(client, d)
	if err != nil {
		return err
	}
	var insertCodeTableFieldValuesHttpUrl = "v2/{project_id}/design/code-tables/{id}/values"
	insertCodeTableFieldValuesPath := client.Endpoint + insertCodeTableFieldValuesHttpUrl
	insertCodeTableFieldValuesPath = strings.ReplaceAll(insertCodeTableFieldValuesPath, "{project_id}", client.ProjectID)
	insertCodeTableFieldValuesPath = strings.ReplaceAll(insertCodeTableFieldValuesPath, "{id}", d.Id())

	insertCodeTableFieldValuesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	insertCodeTableFieldValuesOpt.JSONBody = utils.RemoveNil(buildInsertCodeTableFieldValuesParams(d, ids))
	_, err = client.Request("PUT", insertCodeTableFieldValuesPath, &insertCodeTableFieldValuesOpt)

	if err != nil {
		return err
	}
	return nil
}

func getCodeTableFieldID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	getCodeTableResp, err := readCodeTable(client, d)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture code table: %s", err)
	}

	createCodeTableRespBody, err := utils.FlattenResponse(getCodeTableResp)
	if err != nil {
		return nil, err
	}

	codeTable := utils.PathSearch("data.value", createCodeTableRespBody, nil)

	fieldsIDs, err := jmespath.Search("code_table_fields[*].id", codeTable)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture code table fields: %s is not found in API response", "id")
	}

	return fieldsIDs.([]interface{}), nil
}

func hasFieldValues(d *schema.ResourceData) bool {
	count := false
	rawFields := d.Get("fields").([]interface{})
	for _, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		if _, ok := fieldMap["values"]; ok {
			count = true
			break
		}
	}
	return count
}

func buildInsertCodeTableFieldValuesParams(d *schema.ResourceData, ids []interface{}) map[string]interface{} {
	rawFields := d.Get("fields").([]interface{})
	fields := make([]map[string]interface{}, len(rawFields))
	for i, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		values := buildCodeTableFieldValuesParams(fieldMap, ids[i].(string))
		if values != nil {

			field := make(map[string]interface{})
			field["id"] = ids[i].(string)
			field["code_table_id"] = d.Id()
			field["ordinal"] = i + 1
			field["name_ch"] = fieldMap["name"].(string)
			field["name_en"] = fieldMap["code"].(string)
			field["data_type"] = fieldMap["type"].(string)
			field["description"] = utils.ValueIngoreEmpty(fieldMap["description"].(string))
			field["code_table_field_values"] = utils.ValueIngoreEmpty(values)
			fields[i] = field
		}
	}
	return map[string]interface{}{
		"to_add": fields,
	}
}

func buildCodeTableFieldValuesParams(fieldMap map[string]interface{}, id string) []map[string]interface{} {
	rawValues := fieldMap["values"].([]interface{})

	if len(rawValues) > 0 {
		values := make([]map[string]interface{}, len(rawValues))
		for i, rawValue := range rawValues {
			value := make(map[string]interface{})
			value["ordinal"] = i + 1
			value["fd_id"] = id
			value["fd_value"] = rawValue.(string)
			values[i] = value
		}
		return values
	}
	return nil
}

func resourceArchitectureCodeTableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)
	var mErr *multierror.Error
	getCodeTableProduct := "dataarts"

	getCodeTableClient, err := cfg.NewServiceClient(getCodeTableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getCodeTableResp, err := readCodeTable(getCodeTableClient, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataArts Architecture code table")
	}

	getCodeTableRespBody, err := utils.FlattenResponse(getCodeTableResp)
	if err != nil {
		return diag.FromErr(err)
	}

	codeTable := utils.PathSearch("data.value", getCodeTableRespBody, nil)
	fields, _, err := getCodeTableFieldsAndValues(getCodeTableClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("name", utils.PathSearch("name_ch", codeTable, nil)),
		d.Set("code", utils.PathSearch("name_en", codeTable, nil)),
		d.Set("directory_id", utils.PathSearch("directory_id", codeTable, nil)),
		d.Set("directory_path", utils.PathSearch("directory_path", codeTable, nil)),
		d.Set("fields", fields),
		d.Set("description", utils.PathSearch("description", codeTable, nil)),
		d.Set("created_by", utils.PathSearch("create_by", codeTable, nil)),
		d.Set("created_at", utils.PathSearch("create_time", codeTable, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", codeTable, nil)),
		d.Set("status", utils.PathSearch("status", codeTable, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getCodeTableFieldsAndValues(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]map[string]interface{}, []map[string]interface{}, error) {
	var getCodeTableFieldValuesHttpUrl = "v2/{project_id}/design/code-tables/{id}/values"

	getCodeTableFieldValuesPath := client.Endpoint + getCodeTableFieldValuesHttpUrl
	getCodeTableFieldValuesPath = strings.ReplaceAll(getCodeTableFieldValuesPath, "{project_id}", client.ProjectID)
	getCodeTableFieldValuesPath = strings.ReplaceAll(getCodeTableFieldValuesPath, "{id}", d.Id())

	getCodeTableFieldValuesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	getCodeTableFieldValuesResp, err := client.Request("GET", getCodeTableFieldValuesPath, &getCodeTableFieldValuesOpt)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving DataArts Architecture code table fields: %s", err)
	}

	getCodeTableFieldValuesRespBody, err := utils.FlattenResponse(getCodeTableFieldValuesResp)
	if err != nil {
		return nil, nil, err
	}

	curJson := utils.PathSearch("data.value.records", getCodeTableFieldValuesRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	fields := make([]map[string]interface{}, len(curArray))
	values := make([]map[string]interface{}, len(curArray))
	for i, val := range curArray {
		field := make(map[string]interface{})
		value := make(map[string]interface{})
		field["id"] = utils.PathSearch("id", val, nil)
		field["ordinal"] = utils.PathSearch("ordinal", val, nil)
		value["fd_id"] = field["id"]
		value["values"] = make([]interface{}, 0)
		value["ordinal"] = make([]interface{}, 0)
		value["value_ids"] = make([]interface{}, 0)
		field["table_id"] = utils.PathSearch("code_table_id", val, nil)
		field["name"] = utils.PathSearch("name_ch", val, nil)
		field["code"] = utils.PathSearch("name_en", val, nil)
		field["type"] = utils.PathSearch("data_type", val, nil)
		field["description"] = utils.PathSearch("description", val, nil)
		fields[i] = field
		values[i] = value
	}

	offset := 0
	has_values := false

	getCodeTableFieldValuesPath = getCodeTableFieldValuesPath + fmt.Sprintf("?limit=50&offset=%v", offset)
	for {

		getCodeTableFieldValuesResp, err = client.Request("GET", getCodeTableFieldValuesPath, &getCodeTableFieldValuesOpt)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving DataArts Architecture code table fields: %s", err)
		}

		getCodeTableFieldValuesRespBody, err = utils.FlattenResponse(getCodeTableFieldValuesResp)
		if err != nil {
			return nil, nil, err
		}

		curJson := utils.PathSearch("data.value.records", getCodeTableFieldValuesRespBody, make([]interface{}, 0))
		curArray := curJson.([]interface{})
		has_values = false

		// To build the value structure.
		for _, val := range curArray {
			for j, value := range values {
				if utils.PathSearch("id", val, nil) == value["fd_id"] {
					rawValues := utils.PathSearch("code_table_field_values[*].fd_value", val, nil)
					rawValueIDs := utils.PathSearch("code_table_field_values[*].id", val, nil)
					rawValueOrdinals := utils.PathSearch("code_table_field_values[*].ordinal", val, nil)

					v := rawValues.([]interface{})
					IDs := rawValueIDs.([]interface{})
					ordinals := rawValueOrdinals.([]interface{})
					// Get the values of the field.
					for _, cell := range v {
						values[j]["values"] = append(values[j]["values"].([]interface{}), cell)
						has_values = true
					}
					// Get the IDs of the values.
					for _, id := range IDs {
						values[j]["value_ids"] = append(values[j]["value_ids"].([]interface{}), id)
					}
					// Get the ordinals of the values.
					for _, ordinal := range ordinals {
						values[j]["ordinal"] = append(values[j]["ordinal"].([]interface{}), ordinal)
					}
				}
			}
		}
		if !has_values {
			break
		}
		offset += 50
		index := strings.Index(getCodeTableFieldValuesPath, "offset")
		getCodeTableFieldValuesPath = fmt.Sprintf("%soffset=%v", getCodeTableFieldValuesPath[:index], offset)
	}

	for i, field := range fields {
		for j, value := range values {
			if field["id"] == value["fd_id"] {
				fields[i]["values"] = values[j]["values"]
			}
		}
	}

	return fields, values, nil
}

func resourceArchitectureCodeTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateCodeTableHttpUrl = "v2/{project_id}/design/code-tables/{id}"
		updateCodeTableProduct = "dataarts"
	)
	updateCodeTableClient, err := cfg.NewServiceClient(updateCodeTableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updateCodeTablePath := updateCodeTableClient.Endpoint + updateCodeTableHttpUrl
	updateCodeTablePath = strings.ReplaceAll(updateCodeTablePath, "{project_id}", updateCodeTableClient.ProjectID)
	updateCodeTablePath = strings.ReplaceAll(updateCodeTablePath, "{id}", d.Id())

	updateCodeTableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	updateCodeTableOpt.JSONBody = utils.RemoveNil(buildUpdateCodeTableParams(d))
	_, err = updateCodeTableClient.Request("PUT", updateCodeTablePath, &updateCodeTableOpt)

	if err != nil {
		return diag.Errorf("error updating DataArts Architecture code table: %s", err)
	}

	err = deleteCodeTableFieldValues(updateCodeTableClient, d)
	if err != nil {
		return diag.Errorf("error deleting field values of DataArts Architecture code table: %s", err)
	}

	err = insertCodeTableFieldValues(updateCodeTableClient, d)
	if err != nil {
		return diag.Errorf("error adding field values to DataArts Architecture code table: %s", err)
	}
	return resourceArchitectureCodeTableRead(ctx, d, meta)
}

func deleteCodeTableFieldValues(client *golangsdk.ServiceClient, d *schema.ResourceData) error {

	var codeTableFieldValuesHttpUrl = "v2/{project_id}/design/code-tables/{id}/values"
	codeTableFieldValuesPath := client.Endpoint + codeTableFieldValuesHttpUrl
	codeTableFieldValuesPath = strings.ReplaceAll(codeTableFieldValuesPath, "{project_id}", client.ProjectID)
	codeTableFieldValuesPath = strings.ReplaceAll(codeTableFieldValuesPath, "{id}", d.Id())

	codeTableFieldValuesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}
	rawFields, rawValues, err := getCodeTableFieldsAndValues(client, d)
	if err != nil {
		return err
	}

	n := countFieldNeedToDeleteValues(rawFields)

	// Delete values only the field has values.
	if n > 0 {
		fields := make([]map[string]interface{}, n)
		i := 0
		for j, rawField := range rawFields {

			values := buildCodeTableFieldValuesToDeleteParams(rawValues[j])
			// Skip the field if there is no value.
			if values != nil {
				field := make(map[string]interface{})
				field["id"] = rawField["id"]
				field["code_table_id"] = rawField["table_id"]
				field["name_ch"] = rawField["name"]
				field["name_en"] = rawField["code"]
				field["ordinal"] = rawField["ordinal"]
				field["data_type"] = rawField["type"]
				field["code_table_field_values"] = values
				fields[i] = field
				i++

			}
		}

		codeTableFieldValuesOpt.JSONBody = utils.RemoveNil(map[string]interface{}{"to_remove": fields})
		_, err = client.Request("PUT", codeTableFieldValuesPath, &codeTableFieldValuesOpt)
		if err != nil {
			return err
		}
	}
	return nil
}

func countFieldNeedToDeleteValues(fields []map[string]interface{}) int {
	count := 0
	for _, field := range fields {
		if len(field["values"].([]interface{})) > 0 {
			count++
		}
	}
	return count
}

func buildCodeTableFieldValuesToDeleteParams(val interface{}) []map[string]interface{} {
	valuesJson := val.(map[string]interface{})["values"]
	valuesArray := valuesJson.([]interface{})

	if len(valuesArray) == 0 {
		return nil
	}

	valueIDsJson := val.(map[string]interface{})["value_ids"]
	valueIDsArray := valueIDsJson.([]interface{})
	valuesOrdinalsJson := val.(map[string]interface{})["ordinal"]
	valuesOrdinalsArray := valuesOrdinalsJson.([]interface{})
	values := make([]map[string]interface{}, len(valuesArray))
	for i := range valuesArray {
		value := make(map[string]interface{})
		value["ordinal"] = valuesOrdinalsArray[i]
		value["id"] = valueIDsArray[i].(string)
		values[i] = value
	}
	return values
}

func buildUpdateCodeTableParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":                d.Id(),
		"name_ch":           d.Get("name"),
		"name_en":           d.Get("code"),
		"directory_id":      d.Get("directory_id"),
		"code_table_fields": utils.ValueIngoreEmpty(buildUpdateCodeTableFieldParams(d)),
		"description":       utils.ValueIngoreEmpty(d.Get("description")),
	}
}

func buildUpdateCodeTableFieldParams(d *schema.ResourceData) []map[string]interface{} {
	rawFields := d.Get("fields").([]interface{})
	fields := make([]map[string]interface{}, len(rawFields))
	for i, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		field := make(map[string]interface{})
		field["id"] = fieldMap["id"]
		field["code_table_id"] = d.Id()
		field["ordinal"] = i + 1
		field["name_ch"] = fieldMap["name"].(string)
		field["name_en"] = fieldMap["code"].(string)
		field["data_type"] = fieldMap["type"].(string)
		field["description"] = utils.ValueIngoreEmpty(fieldMap["description"].(string))
		fields[i] = field
	}
	return fields
}

func resourceArchitectureCodeTableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteCodeTableHttpUrl = "v2/{project_id}/design/code-tables"
		deleteCodeTableProduct = "dataarts"
	)
	deleteCodeTableClient, err := cfg.NewServiceClient(deleteCodeTableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deleteCodeTablePath := deleteCodeTableClient.Endpoint + deleteCodeTableHttpUrl
	deleteCodeTablePath = strings.ReplaceAll(deleteCodeTablePath, "{project_id}", deleteCodeTableClient.ProjectID)
	deleteCodeTablePath = strings.ReplaceAll(deleteCodeTablePath, "{id}", d.Id())

	deleteCodeTableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         buildCodeTableDeleteBodyParams(d),
	}

	_, err = deleteCodeTableClient.Request("DELETE", deleteCodeTablePath, &deleteCodeTableOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Architecture code table: %s", err)
	}

	_, err = readCodeTable(deleteCodeTableClient, d)
	if err == nil {
		return diag.Errorf("error deleting DataArts Architecture code table")
	}

	return common.CheckDeletedDiag(d, parseCodeTableError(err), "error deleting DataArts Architecture code table")
}

func buildCodeTableDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ids": []string{d.Id()},
	}
}

func readCodeTable(client *golangsdk.ServiceClient, d *schema.ResourceData) (*http.Response, error) {
	getPath := client.Endpoint + "v2/{project_id}/design/code-tables/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	return client.Request("GET", getPath, &getOpt)
}

func resourceCodeTableImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s'", d.Id())
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := parts[0]
	name := parts[1]
	getCodeTableHttpUrl := "v2/{project_id}/design/code-tables?name={name}"
	getCodeTableProduct := "dataarts"

	getCodeTableClient, err := cfg.NewServiceClient(getCodeTableProduct, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getCodeTablePath := getCodeTableClient.Endpoint + getCodeTableHttpUrl
	getCodeTablePath = strings.ReplaceAll(getCodeTablePath, "{project_id}", getCodeTableClient.ProjectID)
	getCodeTablePath = strings.ReplaceAll(getCodeTablePath, "{name}", name)

	getCodeTableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getCodeTableResp, err := getCodeTableClient.Request("GET", getCodeTablePath, &getCodeTableOpt)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error query DataArts Architecture code table: %s", err)
	}

	getDirectoryRespBody, err := utils.FlattenResponse(getCodeTableResp)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	jsonPaths := fmt.Sprintf("data.value.records[?name_ch=='%s'].id", name)
	codeTableID := utils.PathSearch(jsonPaths, getDirectoryRespBody, make([]interface{}, 0)).([]interface{})
	if len(codeTableID) == 0 {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving DataArts Architecture code table")
	}
	d.SetId(codeTableID[0].(string))
	d.Set("workspace_id", workspaceID)
	return []*schema.ResourceData{d}, nil
}

func parseCodeTableError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode, errorCodeErr := jmespath.Search("errors|[0].error_code", apiError)
		if errorCodeErr != nil {
			return err
		}

		if errorCode == "DLG.6022" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}
