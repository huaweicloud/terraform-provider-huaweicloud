package dataarts

import (
	"context"
	"errors"
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

var (
	architectureDimensionErrCodes = []string{
		"DLG.6007", // dimension not found
		"DLG.0818", // workspace not found
	}
	architectureDimensionNonUpdatableParams = []string{
		"workspace_id",
	}
)

// @API DataArtsStudio POST /v2/{project_id}/design/dimensions
// @API DataArtsStudio GET /v2/{project_id}/design/dimensions/{id}
// @API DataArtsStudio PUT /v2/{project_id}/design/dimensions
// @API DataArtsStudio DELETE /v2/{project_id}/design/dimensions
func ResourceArchitectureDimension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureDimensionCreate,
		ReadContext:   resourceArchitectureDimensionRead,
		UpdateContext: resourceArchitectureDimensionUpdate,
		DeleteContext: resourceArchitectureDimensionDelete,

		CustomizeDiff: config.FlexibleForceNew(architectureDimensionNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureDimensionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dimension is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the dimension belongs.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Chinese name of the dimension.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The English name of the dimension.`,
			},
			"l3_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the business subject to which the dimension belongs.`,
			},
			"dimension_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the dimension.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The asset owner of the dimension.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the dimension.`,
			},
			"code_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The code table ID of the dimension.`,
			},

			// This parameter is only valid when deleting the dimension.
			"is_delete_physical_table": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to delete the physical table corresponding to the dimension.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of attributes of the dimension.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_en": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The English name of the attribute.`,
						},
						"name_ch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The Chinese name of the attribute.`,
						},
						"data_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The data type of the attribute.`,
						},
						"is_primary_key": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether the attribute is the primary key.`,
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The sequence number of the attribute.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The description of the attribute.`,
						},
						"is_biz_primary": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the attribute is the business primary key.`,
						},
						"is_partition_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the attribute is the partition key.`,
						},
						"not_null": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the attribute is not null.`,
						},

						// Attributes.
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the attribute.`,
						},
						"domain_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain type of the attribute.`,
						},
						"code_table_field_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The code table field ID of the attribute.`,
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the attribute.`,
						},
						"data_type_extend": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type extension of the attribute.`,
						},
						"stand_row_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated data standard ID of the attribute.`,
						},
						"stand_row_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated data standard name of the attribute.`,
						},
						"quality_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The quality information of the attribute.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"secrecy_levels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The secrecy levels of the attribute.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The publishing status of the attribute.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the attribute, in RFC3339 format.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the attribute, in RFC3339 format.`,
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alias of the attribute.`,
						},
						"self_defined_fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The self-defined fields of the attribute.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"datasource": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The data source configuration of the dimension.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dw_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the data connection.`,
						},
						"dw_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the data connection.`,
						},
						"biz_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The business type of the data source.`,
						},
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The name of the database corresponding to the data connection.`,
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The queue name corresponding to the DLI data connection.`,
						},

						// Attributes.
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data source.`,
						},
						"dw_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data connection.`,
						},
						"schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database schema.`,
						},
						"biz_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The business ID of the data source.`,
						},
					},
				},
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publishing status of the dimension.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account name of the user who created the dimension.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the dimension, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the dimension, in RFC3339 format.`,
			},
			"l1_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the L1 subject to which the dimension belongs.`,
			},
			"l2_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the L2 subject to which the dimension belongs.`,
			},
			"l1_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the L1 subject to which the dimension belongs.`,
			},
			"l2_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the L2 subject to which the dimension belongs.`,
			},
			"l3_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the L3 subject to which the dimension belongs.`,
			},
			"table_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The table type of the dimension.`,
			},
			"distribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distribute type of the dimension.`,
			},
			"distribute_column": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distribute column of the dimension.`,
			},
			"compression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The compression type of the dimension.`,
			},
			"obs_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS location of the dimension.`,
			},
			"pre_combine_field": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The pre-combine field of the dimension.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the dimension.`,
			},
			"configs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The configs of the dimension.`,
			},
			"env_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The environment type of the dimension.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The model ID of the dimension.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the dimension.`,
			},
			"dev_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development environment version of the dimension.`,
			},
			"prod_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The production environment version of the dimension.`,
			},
			"dev_version_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development environment version name of the dimension.`,
			},
			"prod_version_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The production environment version name of the dimension.`,
			},
			"model": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The model information of the dimension.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the workspace.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the workspace.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the workspace.`,
						},
						"is_physical": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether it is a physical table.`,
						},
						"frequent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether it is frequently used.`,
						},
						"top": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether it is a top-level governance.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data governance level.`,
						},
						"dw_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data warehouse type.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the workspace, in RFC3339 format.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the workspace, in RFC3339 format.`,
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the workspace.`,
						},
						"update_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The updater of the workspace.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace type.`,
						},
						"biz_catalog_ids": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated business catalog IDs.`,
						},
						"databases": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The database names.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"table_model_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The table model prefix.`,
						},
					},
				},
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateArchitectureDimensionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name_ch":        d.Get("name_ch"),
		"name_en":        d.Get("name_en"),
		"l3_id":          d.Get("l3_id"),
		"dimension_type": d.Get("dimension_type"),
		"owner":          d.Get("owner"),
		"description":    d.Get("description"),
		"code_table_id":  d.Get("code_table_id"),
		"attributes":     buildArchitectureDimensionAttributes(d.Get("attributes").([]interface{})),
		"datasource":     buildArchitectureDimensionDatasource(d.Get("datasource").([]interface{})),
	}
}

func buildArchitectureDimensionAttributes(attributes []interface{}) []map[string]interface{} {
	if len(attributes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(attributes))
	for _, item := range attributes {
		result = append(result, map[string]interface{}{
			"name_en":          utils.PathSearch("name_en", item, nil),
			"name_ch":          utils.PathSearch("name_ch", item, nil),
			"description":      utils.PathSearch("description", item, nil),
			"data_type":        utils.PathSearch("data_type", item, nil),
			"is_primary_key":   utils.PathSearch("is_primary_key", item, nil),
			"is_biz_primary":   utils.PathSearch("is_biz_primary", item, nil),
			"is_partition_key": utils.PathSearch("is_partition_key", item, nil),
			"not_null":         utils.PathSearch("not_null", item, nil),
			"ordinal":          utils.PathSearch("ordinal", item, nil),
		})
	}

	return result
}

func buildArchitectureDimensionDatasource(datasource []interface{}) map[string]interface{} {
	if len(datasource) == 0 {
		return nil
	}

	datasourceItem := datasource[0]
	return map[string]interface{}{
		"biz_type":   utils.PathSearch("biz_type", datasourceItem, nil),
		"dw_id":      utils.PathSearch("dw_id", datasourceItem, nil),
		"dw_type":    utils.PathSearch("dw_type", datasourceItem, nil),
		"db_name":    utils.PathSearch("db_name", datasourceItem, nil),
		"queue_name": utils.PathSearch("queue_name", datasourceItem, nil),
	}
}

func createArchitectureDimension(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	var (
		httpURL = "v2/{project_id}/design/dimensions"
	)

	createPath := client.Endpoint + httpURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildCreateArchitectureDimensionBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	dimensionId := utils.PathSearch("data.value.id", respBody, "").(string)
	if dimensionId == "" {
		return "", errors.New("unable to find dimension ID from the API response")
	}

	return dimensionId, nil
}

func resourceArchitectureDimensionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	dimensionId, err := createArchitectureDimension(client, d)
	if err != nil {
		return diag.Errorf("error creating architecture dimension: %s", err)
	}

	d.SetId(dimensionId)

	return resourceArchitectureDimensionRead(ctx, d, meta)
}

func GetArchitectureDimensionById(client *golangsdk.ServiceClient, workspaceId, dimensionId string) (interface{}, error) {
	var (
		httpURL = "v2/{project_id}/design/dimensions/{id}"
	)

	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", dimensionId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(workspaceId),
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	resourceValue := utils.PathSearch("data.value", respBody, make(map[string]interface{})).(map[string]interface{})
	if len(resourceValue) == 0 {
		errMsg := utils.PathSearch("errors|[0].error_msg", respBody, "the dimension not found")
		requestId := utils.PathSearch("errors|[0].request_id", respBody, "")
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/design/dimensions/{id}",
				RequestId: requestId.(string),
				Body:      []byte(fmt.Sprintf("error message: %s. dimension ID: %s", errMsg, dimensionId)),
			},
		}
	}

	return resourceValue, nil
}

func resourceArchitectureDimensionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		workspaceId           = d.Get("workspace_id").(string)
		isDeletePhysicalTable = d.Get("is_delete_physical_table").(bool)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetArchitectureDimensionById(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", architectureDimensionErrCodes...),
			"error retrieving dimension",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("workspace_id", workspaceId),
		d.Set("is_delete_physical_table", isDeletePhysicalTable),
		d.Set("name_ch", utils.PathSearch("name_ch", respBody, nil)),
		d.Set("name_en", utils.PathSearch("name_en", respBody, nil)),
		d.Set("l3_id", utils.PathSearch("l3_id", respBody, nil)),
		d.Set("dimension_type", utils.PathSearch("dimension_type", respBody, nil)),
		d.Set("owner", utils.PathSearch("owner", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_by", utils.PathSearch("create_by", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("l1_id", utils.PathSearch("l1_id", respBody, nil)),
		d.Set("l2_id", utils.PathSearch("l2_id", respBody, nil)),
		d.Set("l1_name", utils.PathSearch("l1", respBody, nil)),
		d.Set("l2_name", utils.PathSearch("l2", respBody, nil)),
		d.Set("l3_name", utils.PathSearch("l3", respBody, nil)),
		d.Set("table_type", utils.PathSearch("table_type", respBody, nil)),
		d.Set("distribute", utils.PathSearch("distribute", respBody, nil)),
		d.Set("distribute_column", utils.PathSearch("distribute_column", respBody, nil)),
		d.Set("compression", utils.PathSearch("compression", respBody, nil)),
		d.Set("obs_location", utils.PathSearch("obs_location", respBody, nil)),
		d.Set("pre_combine_field", utils.PathSearch("pre_combine_field", respBody, nil)),
		d.Set("alias", utils.PathSearch("alias", respBody, nil)),
		d.Set("configs", utils.PathSearch("configs", respBody, nil)),
		d.Set("env_type", utils.PathSearch("env_type", respBody, nil)),
		d.Set("model_id", utils.PathSearch("model_id", respBody, nil)),
		d.Set("update_by", utils.PathSearch("update_by", respBody, nil)),
		d.Set("code_table_id", utils.PathSearch("code_table_id", respBody, nil)),
		d.Set("dev_version", utils.PathSearch("dev_version", respBody, nil)),
		d.Set("prod_version", utils.PathSearch("prod_version", respBody, nil)),
		d.Set("dev_version_name", utils.PathSearch("dev_version_name", respBody, nil)),
		d.Set("prod_version_name", utils.PathSearch("prod_version_name", respBody, nil)),
		d.Set("attributes", flattenArchitectureDimensionAttributes(utils.PathSearch("attributes", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("datasource", flattenArchitectureDimensionDatasource(utils.PathSearch("datasource", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("model", flattenDimensionModel(utils.PathSearch("model", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenArchitectureDimensionAttributes(attributes []interface{}) []interface{} {
	if len(attributes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(attributes))
	for _, v := range attributes {
		result = append(result, map[string]interface{}{
			"name_en":             utils.PathSearch("name_en", v, nil),
			"name_ch":             utils.PathSearch("name_ch", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"data_type":           utils.PathSearch("data_type", v, nil),
			"is_primary_key":      utils.PathSearch("is_primary_key", v, nil),
			"is_biz_primary":      utils.PathSearch("is_biz_primary", v, nil),
			"is_partition_key":    utils.PathSearch("is_partition_key", v, nil),
			"not_null":            utils.PathSearch("not_null", v, nil),
			"ordinal":             utils.PathSearch("ordinal", v, nil),
			"id":                  utils.PathSearch("id", v, nil),
			"domain_type":         utils.PathSearch("domain_type", v, nil),
			"code_table_field_id": utils.PathSearch("code_table_field_id", v, nil),
			"create_by":           utils.PathSearch("create_by", v, nil),
			"data_type_extend":    utils.PathSearch("data_type_extend", v, nil),
			"stand_row_id":        utils.PathSearch("stand_row_id", v, nil),
			"stand_row_name":      utils.PathSearch("stand_row_name", v, nil),
			"quality_infos":       utils.PathSearch("quality_infos", v, nil),
			"secrecy_levels":      utils.PathSearch("secrecy_levels", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"alias":               utils.PathSearch("alias", v, nil),
			"self_defined_fields": utils.PathSearch("self_defined_fields", v, nil),
		})
	}

	return result
}

func flattenArchitectureDimensionDatasource(datasource map[string]interface{}) []map[string]interface{} {
	if len(datasource) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":         utils.PathSearch("id", datasource, nil),
			"biz_type":   utils.PathSearch("biz_type", datasource, nil),
			"biz_id":     utils.PathSearch("biz_id", datasource, nil),
			"dw_id":      utils.PathSearch("dw_id", datasource, nil),
			"dw_type":    utils.PathSearch("dw_type", datasource, nil),
			"dw_name":    utils.PathSearch("dw_name", datasource, nil),
			"db_name":    utils.PathSearch("db_name", datasource, nil),
			"queue_name": utils.PathSearch("queue_name", datasource, nil),
			"schema":     utils.PathSearch("schema", datasource, nil),
		},
	}
}

func flattenDimensionModel(model map[string]interface{}) []map[string]interface{} {
	if len(model) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                 utils.PathSearch("id", model, nil),
			"name":               utils.PathSearch("name", model, nil),
			"description":        utils.PathSearch("description", model, nil),
			"is_physical":        utils.PathSearch("is_physical", model, nil),
			"frequent":           utils.PathSearch("frequent", model, nil),
			"top":                utils.PathSearch("top", model, nil),
			"level":              utils.PathSearch("level", model, nil),
			"dw_type":            utils.PathSearch("dw_type", model, nil),
			"create_time":        utils.PathSearch("create_time", model, nil),
			"update_time":        utils.PathSearch("update_time", model, nil),
			"create_by":          utils.PathSearch("create_by", model, nil),
			"update_by":          utils.PathSearch("update_by", model, nil),
			"type":               utils.PathSearch("type", model, nil),
			"biz_catalog_ids":    utils.PathSearch("biz_catalog_ids", model, nil),
			"databases":          utils.PathSearch("databases", model, nil),
			"table_model_prefix": utils.PathSearch("table_model_prefix", model, nil),
		},
	}
}

func buildUpdateArchitectureDimensionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":             d.Id(),
		"name_ch":        d.Get("name_ch"),
		"name_en":        d.Get("name_en"),
		"l3_id":          d.Get("l3_id"),
		"dimension_type": d.Get("dimension_type"),
		"owner":          d.Get("owner"),
		"description":    d.Get("description"),
		"code_table_id":  d.Get("code_table_id"),
		"attributes":     buildArchitectureDimensionAttributes(d.Get("attributes").([]interface{})),
		"datasource":     buildArchitectureDimensionDatasource(d.Get("datasource").([]interface{})),
	}
}

func updateArchitectureDimension(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpURL := "v2/{project_id}/design/dimensions"
	updatePath := client.Endpoint + httpURL
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildUpdateArchitectureDimensionBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceArchitectureDimensionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatedParameters := []string{
		"name_ch",
		"name_en",
		"l3_id",
		"dimension_type",
		"owner",
		"description",
		"code_table_id",
		"attributes",
		"datasource",
	}

	if d.HasChanges(updatedParameters...) {
		if err = updateArchitectureDimension(client, d); err != nil {
			return diag.Errorf("error updating dimension (%s): %s", d.Id(), err)
		}
	}

	return resourceArchitectureDimensionRead(ctx, d, meta)
}

func deleteArchitectureDimension(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpURL = "v2/{project_id}/design/dimensions"
	)

	deletePath := client.Endpoint + httpURL
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	// delete parameters
	bodyParams := map[string]interface{}{
		"ids": []string{d.Id()},
	}
	if d.Get("is_delete_physical_table").(bool) {
		bodyParams["del_types"] = "PHYSICAL_TABLE"
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(bodyParams),
	}
	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceArchitectureDimensionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := deleteArchitectureDimension(client, d); err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", architectureDimensionErrCodes[0]),
			fmt.Sprintf("error deleting dimension (%s)", d.Id()),
		)
	}

	return nil
}

func resourceArchitectureDimensionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
