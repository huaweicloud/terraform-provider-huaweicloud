package dataarts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio PUT /v2/{project_id}/design/code-tables/{id}/values
// @API DataArtsStudio GET /v2/{project_id}/design/code-tables/{id}/values
func ResourceArchitectureCodeTableValues() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureCodeTableValuesCreate,
		ReadContext:   resourceArchitectureCodeTableValuesRead,
		UpdateContext: resourceArchitectureCodeTableValuesUpdate,
		DeleteContext: resourceArchitectureCodeTableValuesDelete,

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
			"table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the code table.",
			},
			"field_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the code table field.",
			},
			"field_ordinal": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ordinal of the code table field.",
			},
			"field_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the code table field.",
			},
			"field_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The code of the code table field.",
			},
			"field_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the code table field.",
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of a field.",
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ordinal of a value.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of a value.",
						},
					},
				},
				Description: "The values list of the code table field.",
			},
		},
	}
}

func resourceArchitectureCodeTableValuesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	product := "dataarts"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath, createOpt := buildCodeTableValuesPathAndOpt(client, d)
	createOpt.JSONBody = utils.RemoveNil(buildCreateCodeTableValuesParams(d))
	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Architecture code table values: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("field_name").(string))

	return resourceArchitectureCodeTableValuesRead(ctx, d, meta)
}

func buildCreateCodeTableValuesParams(d *schema.ResourceData) map[string]interface{} {
	params := make([]map[string]interface{}, 1)
	params[0] = map[string]interface{}{
		"id":                      d.Get("field_id"),
		"ordinal":                 d.Get("field_ordinal"),
		"name_ch":                 d.Get("field_name"),
		"name_en":                 d.Get("field_code"),
		"data_type":               d.Get("field_type"),
		"code_table_field_values": buildCreateCodeTableFieldValuesParams(d),
	}
	return map[string]interface{}{
		"to_add": params,
	}
}

func buildCreateCodeTableFieldValuesParams(d *schema.ResourceData) []map[string]interface{} {
	rawValues := d.Get("values").([]interface{})
	values := make([]map[string]interface{}, len(rawValues))
	fieldID := d.Get("field_id").(string)
	for i, rawValue := range rawValues {
		valueMap := rawValue.(map[string]interface{})
		value := map[string]interface{}{
			"fd_id":    fieldID,
			"fd_value": valueMap["value"].(string),
			"ordinal":  i + 1,
		}
		values[i] = value
	}
	return values
}

func resourceArchitectureCodeTableValuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	var mErr *multierror.Error
	product := "dataarts"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath, getOpt := buildCodeTableValuesPathAndOpt(client, d)
	fieldName := d.Get("field_name").(string)
	values, fieldJson, err := flattenCodeTableFieldValues(client, getOpt,
		getPath, fieldName)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Architecture code table field values: %s", err)
	}

	if len(values) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DataArts Architecture code table field values")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("table_id", utils.PathSearch("code_table_id", fieldJson, nil)),
		d.Set("field_id", utils.PathSearch("id", fieldJson, nil)),
		d.Set("field_ordinal", utils.PathSearch("ordinal", fieldJson, nil)),
		d.Set("field_name", utils.PathSearch("name_ch", fieldJson, nil)),
		d.Set("field_code", utils.PathSearch("name_en", fieldJson, nil)),
		d.Set("field_type", utils.PathSearch("data_type", fieldJson, nil)),
		d.Set("values", values),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCodeTableFieldValues(client *golangsdk.ServiceClient, getCodeTableFieldValuesOpt golangsdk.RequestOpts,
	getCodeTableFieldValuesPath, fieldName string) ([]interface{}, interface{}, error) {
	var fieldJson interface{}
	offset := 0
	values := make([]interface{}, 0)

	for {
		path := fmt.Sprintf("%s?limit=50&offset=%d", getCodeTableFieldValuesPath, offset)
		getCodeTableFieldValuesResp, err := client.Request("GET", path, &getCodeTableFieldValuesOpt)
		if err != nil {
			return nil, nil, err
		}

		getCodeTableFieldValuesRespBody, err := utils.FlattenResponse(getCodeTableFieldValuesResp)
		if err != nil {
			return nil, nil, err
		}

		findFieldExpr := fmt.Sprintf("data.value.records[?name_ch=='%s'] | [0]", fieldName)
		fieldJson = utils.PathSearch(findFieldExpr, getCodeTableFieldValuesRespBody, nil)
		valuesJson := utils.PathSearch("code_table_field_values", fieldJson, make([]interface{}, 0))
		valueArray := valuesJson.([]interface{})

		hasValues := false
		for _, rawValue := range valueArray {
			id := utils.PathSearch("id", rawValue, nil)

			// Filter out the deleted values.
			// If a code table has a field, the deleted values will be cleared.
			// If a code table has multiple fields, the structure of the deleted value will be retained, with the id set to nil.
			if id != nil {
				value := map[string]interface{}{
					"id":      id,
					"ordinal": utils.PathSearch("ordinal", rawValue, nil),
					"value":   utils.PathSearch("fd_value", rawValue, nil),
				}

				values = append(values, value)
				hasValues = true
			}
		}

		if !hasValues {
			log.Print("[INFO] The process of retrieving values has ended.")
			break
		}
		offset += 50
	}
	return values, fieldJson, nil
}

func resourceArchitectureCodeTableValuesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "dataarts"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePath, updateOpt := buildCodeTableValuesPathAndOpt(client, d)
	updateOpt.JSONBody = utils.RemoveNil(buildRemoveCodeTableValuesParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)

	if err != nil {
		return diag.Errorf("error updating DataArts Architecture code table values: %s", err)
	}

	updateOpt.JSONBody = utils.RemoveNil(buildCreateCodeTableValuesParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)

	if err != nil {
		return diag.Errorf("error updating DataArts Architecture code table values: %s", err)
	}

	return resourceArchitectureCodeTableValuesRead(ctx, d, meta)
}

func buildRemoveCodeTableValuesParams(d *schema.ResourceData) map[string]interface{} {
	params := make([]map[string]interface{}, 1)
	params[0] = map[string]interface{}{
		"id":                      d.Get("field_id"),
		"ordinal":                 d.Get("field_ordinal"),
		"name_ch":                 d.Get("field_name"),
		"name_en":                 d.Get("field_code"),
		"data_type":               d.Get("field_type"),
		"code_table_field_values": buildRemoveCodeTableFieldValuesParams(d),
	}
	return map[string]interface{}{
		"to_remove": params,
	}
}

func buildRemoveCodeTableFieldValuesParams(d *schema.ResourceData) []map[string]interface{} {
	rawValues, _ := d.GetChange("values")
	oldValues := rawValues.([]interface{})
	values := make([]map[string]interface{}, len(oldValues))
	for i, oldValue := range oldValues {
		valueMap := oldValue.(map[string]interface{})
		value := map[string]interface{}{
			"id":      valueMap["id"],
			"ordinal": valueMap["ordinal"],
		}
		values[i] = value
	}
	return values
}

func resourceArchitectureCodeTableValuesDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	product := "dataarts"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	path, opt := buildCodeTableValuesPathAndOpt(client, d)
	opt.JSONBody = utils.RemoveNil(buildRemoveCodeTableValuesParams(d))
	_, err = client.Request("PUT", path, &opt)

	if err != nil {
		return diag.Errorf("error deleting DataArts Architecture code table values: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	opt.JSONBody = nil
	fieldName := d.Get("field_name").(string)
	values, _, err := flattenCodeTableFieldValues(client, opt, path, fieldName)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Architecture code table field values: %s", err)
	}

	if len(values) > 0 {
		return diag.Errorf("error deleting DataArts Architecture code table values")
	}
	return nil
}

func buildCodeTableValuesPathAndOpt(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, golangsdk.RequestOpts) {
	httpUrl := "v2/{project_id}/design/code-tables/{id}/values"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{id}", d.Get("table_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	return path, opt
}
