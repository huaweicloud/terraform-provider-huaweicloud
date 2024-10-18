package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/code-tables
// @API DataArtsStudio POST /v2/{project_id}/design/code-tables
// @API DataArtsStudio DELETE /v2/{project_id}/design/code-tables
// @API DataArtsStudio PUT /v2/{project_id}/design/code-tables/{id}
// @API DataArtsStudio GET /v2/{project_id}/design/code-tables/{id}
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
				Description: "The region in which to create the resource.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of DataArts Studio workspace.",
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

	tableId := utils.PathSearch("data.value.id", createCodeTableRespBody, "").(string)
	if tableId == "" {
		return diag.Errorf("unable to find the DataArts Architecture code table ID from the API response")
	}

	d.SetId(tableId)

	return resourceArchitectureCodeTableRead(ctx, d, meta)
}

func buildCreateCodeTableParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name_ch":           d.Get("name"),
		"name_en":           d.Get("code"),
		"directory_id":      d.Get("directory_id"),
		"code_table_fields": buildCreateCodeTableFieldParams(d),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func buildCreateCodeTableFieldParams(d *schema.ResourceData) []map[string]interface{} {
	rawFields := d.Get("fields").([]interface{})
	fields := make([]map[string]interface{}, len(rawFields))
	for i, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		field := map[string]interface{}{
			"ordinal":     i + 1,
			"name_ch":     fieldMap["name"].(string),
			"name_en":     fieldMap["code"].(string),
			"data_type":   fieldMap["type"].(string),
			"description": utils.ValueIgnoreEmpty(fieldMap["description"].(string)),
		}
		fields[i] = field
	}
	return fields
}

func resourceArchitectureCodeTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.6022"),
			"error retrieving DataArts Architecture code table")
	}

	getCodeTableRespBody, err := utils.FlattenResponse(getCodeTableResp)
	if err != nil {
		return diag.FromErr(err)
	}

	codeTable := utils.PathSearch("data.value", getCodeTableRespBody, nil)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("name", utils.PathSearch("name_ch", codeTable, nil)),
		d.Set("code", utils.PathSearch("name_en", codeTable, nil)),
		d.Set("directory_id", utils.PathSearch("directory_id", codeTable, nil)),
		d.Set("directory_path", utils.PathSearch("directory_path", codeTable, nil)),
		d.Set("fields", flattenCodeTableFields(codeTable)),
		d.Set("description", utils.PathSearch("description", codeTable, nil)),
		d.Set("created_by", utils.PathSearch("create_by", codeTable, nil)),
		d.Set("created_at", utils.PathSearch("create_time", codeTable, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", codeTable, nil)),
		d.Set("status", utils.PathSearch("status", codeTable, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCodeTableFields(codeTable interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("code_table_fields", codeTable, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	fields := make([]map[string]interface{}, len(curArray))

	for i, val := range curArray {
		field := map[string]interface{}{
			"id":          utils.PathSearch("id", val, nil),
			"ordinal":     utils.PathSearch("ordinal", val, nil),
			"name":        utils.PathSearch("name_ch", val, nil),
			"code":        utils.PathSearch("name_en", val, nil),
			"type":        utils.PathSearch("data_type", val, nil),
			"description": utils.PathSearch("description", val, nil),
		}
		fields[i] = field
	}
	return fields
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

	return resourceArchitectureCodeTableRead(ctx, d, meta)
}

func buildUpdateCodeTableParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":                d.Id(),
		"name_ch":           d.Get("name"),
		"name_en":           d.Get("code"),
		"directory_id":      d.Get("directory_id"),
		"code_table_fields": buildUpdateCodeTableFieldParams(d),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func buildUpdateCodeTableFieldParams(d *schema.ResourceData) []map[string]interface{} {
	rawFields := d.Get("fields").([]interface{})
	fields := make([]map[string]interface{}, len(rawFields))
	for i, rawField := range rawFields {
		fieldMap := rawField.(map[string]interface{})
		field := map[string]interface{}{
			"id":            fieldMap["id"],
			"code_table_id": d.Id(),
			"ordinal":       i + 1,
			"name_ch":       fieldMap["name"].(string),
			"name_en":       fieldMap["code"].(string),
			"data_type":     fieldMap["type"].(string),
			"description":   utils.ValueIgnoreEmpty(fieldMap["description"].(string)),
		}
		fields[i] = field
	}
	return fields
}

func resourceArchitectureCodeTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = readCodeTable(deleteCodeTableClient, d)
	if err == nil {
		return diag.Errorf("error deleting DataArts Architecture code table")
	}

	return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.6022"),
		"error deleting DataArts Architecture code table")
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
		return []*schema.ResourceData{d}, fmt.Errorf("error querying DataArts Architecture code table: %s", err)
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
