package dataarts

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var tableModelErrCodes = []string{
	"DLG.0818", // Workspace not found.
	"DLG.6019", // Resource not found.
	"DLG.3902", // Resource ID value is incorrect.
}

// @API DataArtsStudio POST /v2/{project_id}/design/table-model
// @API DataArtsStudio DELETE /v2/{project_id}/design/table-model
// @API DataArtsStudio GET /v2/{project_id}/design/table-model
// @API DataArtsStudio PUT /v2/{project_id}/design/table-model
func ResourceArchitectureTableModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTableModelCreate,
		ReadContext:   resourceTableModelRead,
		UpdateContext: resourceTableModelUpdate,
		DeleteContext: resourceTableModelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataArtsStudioImportState,
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
			"model_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dw_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subject_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_table_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_en": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_type_extend": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stand_row_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"related_logic_attr_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"related_logic_attr_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"related_logic_attr_name_en": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"stand_row_name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"ordinal": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_partition_key": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_primary_key": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"not_null": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_foreign_key": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"extend_field": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return hashcode.String(m["name"].(string) + m["name_en"].(string) + m["data_type"].(string))
				},
			},
			"del_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dw_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"obs_location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_table_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related_logic_table_model_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related_logic_model_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dirty_out_database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dirty_out_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dirty_out_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partition_conf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mappings": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_field_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"target_field_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"target_field_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"source_field_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updated_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updated_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_table_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"target_table_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_table_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"source_table_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_tables": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table1_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"table2_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"join_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"table1_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"table2_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"join_fields": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field1_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"field2_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"field1_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"field2_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"source_fields": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_ids": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"field_names": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"target_field_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"transform_expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"changed": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"src_model_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"src_model_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"view_text": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dirty_out_switch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"table_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"compression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"distribute_column": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"data_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dlf_task_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"use_recently_partition": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reversed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"parent_table_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"related_logic_model_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"related_logic_table_model_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dw_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parent_table_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tb_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logic_tb_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"catalog_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_related_physical_table": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_related_logic_table": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_partition": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"physical_table_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dev_physical_table_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"technical_asset_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"business_asset_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_data_link_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_quality_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"summary_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTableModelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	createTableModelHttpUrl := "v2/{project_id}/design/table-model"
	createTableModelProduct := "dataarts"
	createTableModelClient, err := cfg.NewServiceClient(createTableModelProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	createTableModelPath := createTableModelClient.Endpoint + createTableModelHttpUrl
	createTableModelPath = strings.ReplaceAll(createTableModelPath, "{project_id}", createTableModelClient.ProjectID)
	createTableModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateTableModelBodyParams(d)),
	}

	createTableModelResp, err := createTableModelClient.Request("POST", createTableModelPath, &createTableModelOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createTableModelRespBody, err := utils.FlattenResponse(createTableModelResp)
	if err != nil {
		return diag.FromErr(err)
	}

	modelId := utils.PathSearch("data.value.id", createTableModelRespBody, "").(string)
	if modelId == "" {
		return diag.Errorf("unable to find the DataArts Architecture table model ID from the API response")
	}
	d.SetId(modelId)

	return resourceTableModelRead(ctx, d, meta)
}

func buildCreateOrUpdateTableModelBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"model_id":                       d.Get("model_id"),
		"dw_type":                        d.Get("dw_type"),
		"biz_catalog_id":                 d.Get("subject_id"),
		"tb_name":                        d.Get("physical_table_name"),
		"logic_tb_name":                  d.Get("table_name"),
		"description":                    d.Get("description"),
		"attributes":                     buildTableModelRequestBodyAttributes(d.Get("attributes").(*schema.Set)),
		"relations":                      buildTableModelRequestBodyRelations(d.Get("relations")),
		"mappings":                       buildTableModelRequestBodyMappings(d.Get("mappings")),
		"configs":                        utils.ValueIgnoreEmpty(d.Get("configs")),
		"dw_id":                          utils.ValueIgnoreEmpty(d.Get("dw_id")),
		"parent_table_id":                utils.ValueIgnoreEmpty(d.Get("parent_table_id")),
		"related_logic_table_id":         utils.ValueIgnoreEmpty(d.Get("related_logic_table_model_id")),
		"related_logic_table_model_id":   utils.ValueIgnoreEmpty(d.Get("related_logic_model_id")),
		"parent_table_name":              utils.ValueIgnoreEmpty(d.Get("parent_table_name")),
		"related_logic_table_name":       utils.ValueIgnoreEmpty(d.Get("related_logic_table_model_name")),
		"related_logic_table_model_name": utils.ValueIgnoreEmpty(d.Get("related_logic_model_name")),
		"dw_name":                        utils.ValueIgnoreEmpty(d.Get("dw_name")),
		"owner":                          utils.ValueIgnoreEmpty(d.Get("owner")),
		"table_type":                     utils.ValueIgnoreEmpty(d.Get("table_type")),
		"compression":                    utils.ValueIgnoreEmpty(d.Get("compression")),
		"code":                           utils.ValueIgnoreEmpty(d.Get("code")),
		"distribute":                     utils.ValueIgnoreEmpty(d.Get("distribute")),
		"distribute_column":              utils.ValueIgnoreEmpty(d.Get("distribute_column")),
		"data_format":                    utils.ValueIgnoreEmpty(d.Get("data_format")),
		"dlf_task_id":                    utils.ValueIgnoreEmpty(d.Get("dlf_task_id")),
		"use_recently_partition":         utils.ValueIgnoreEmpty(d.Get("use_recently_partition")),
		"reversed":                       d.Get("reversed"),
		"db_name":                        utils.ValueIgnoreEmpty(d.Get("db_name")),
		"queue_name":                     utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"schema":                         utils.ValueIgnoreEmpty(d.Get("schema")),
		"obs_location":                   utils.ValueIgnoreEmpty(d.Get("obs_location")),
		"dirty_out_switch":               utils.ValueIgnoreEmpty(d.Get("dirty_out_switch")),
		"dirty_out_database":             utils.ValueIgnoreEmpty(d.Get("dirty_out_database")),
		"dirty_out_prefix":               utils.ValueIgnoreEmpty(d.Get("dirty_out_prefix")),
		"dirty_out_suffix":               utils.ValueIgnoreEmpty(d.Get("dirty_out_suffix")),
		"partition_conf":                 utils.ValueIgnoreEmpty(d.Get("partition_conf")),
	}
	return bodyParams
}

func buildTableModelRequestBodyAttributes(rawParams *schema.Set) []map[string]interface{} {
	if rawParams.Len() == 0 {
		return nil
	}
	attributes := make([]map[string]interface{}, rawParams.Len())
	for _, val := range rawParams.List() {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"name_ch":                    raw["name"],
			"name_en":                    raw["name_en"],
			"data_type":                  raw["data_type"],
			"data_type_extend":           utils.ValueIgnoreEmpty(raw["data_type_extend"]),
			"description":                utils.ValueIgnoreEmpty(raw["description"]),
			"stand_row_id":               utils.ValueIgnoreEmpty(raw["stand_row_id"]),
			"related_logic_attr_id":      utils.ValueIgnoreEmpty(raw["related_logic_attr_id"]),
			"stand_row_name":             utils.ValueIgnoreEmpty(raw["stand_row_name"]),
			"related_logic_attr_name":    utils.ValueIgnoreEmpty(raw["related_logic_attr_name"]),
			"related_logic_attr_name_en": utils.ValueIgnoreEmpty(raw["related_logic_attr_name_en"]),
			"ordinal":                    utils.ValueIgnoreEmpty(raw["ordinal"]),
			"is_partition_key":           utils.ValueIgnoreEmpty(raw["is_partition_key"]),
			"is_primary_key":             utils.ValueIgnoreEmpty(raw["is_primary_key"]),
			"is_foreign_key":             utils.ValueIgnoreEmpty(raw["is_foreign_key"]),
			"extend_field":               utils.ValueIgnoreEmpty(raw["extend_field"]),
			"not_null":                   utils.ValueIgnoreEmpty(raw["not_null"]),
			"code":                       utils.ValueIgnoreEmpty(raw["code"]),
		}
		attributes = append(attributes, params)
	}
	return attributes
}

func buildTableModelRequestBodyRelations(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	relations := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"name":              raw["name"],
			"source_table_id":   utils.ValueIgnoreEmpty(raw["source_table_id"]),
			"source_table_name": utils.ValueIgnoreEmpty(raw["source_table_name"]),
			"source_type":       utils.ValueIgnoreEmpty(raw["source_type"]),
			"target_table_id":   utils.ValueIgnoreEmpty(raw["target_table_id"]),
			"target_table_name": utils.ValueIgnoreEmpty(raw["target_table_name"]),
			"target_type":       utils.ValueIgnoreEmpty(raw["target_type"]),
			"role":              utils.ValueIgnoreEmpty(raw["role"]),
			"mappings":          buildTableModelRequestBodyRelationsMappings(raw["mappings"]),
		}
		relations = append(relations, params)
	}
	return relations
}

func buildTableModelRequestBodyRelationsMappings(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	mappings := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		mapping := map[string]interface{}{
			"source_field_id":   utils.ValueIgnoreEmpty(raw["source_field_id"]),
			"source_field_name": utils.ValueIgnoreEmpty(raw["source_field_name"]),
			"target_field_id":   utils.ValueIgnoreEmpty(raw["target_field_id"]),
			"target_field_name": utils.ValueIgnoreEmpty(raw["target_field_name"]),
		}
		mappings = append(mappings, mapping)
	}
	return mappings
}

func buildTableModelRequestBodyMappings(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	mappings := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"name":           raw["name"],
			"src_model_id":   raw["src_model_id"],
			"src_model_name": utils.ValueIgnoreEmpty(raw["src_model_name"]),
			"view_text":      utils.ValueIgnoreEmpty(raw["view_text"]),
			"source_tables":  buildTableModelRequestBodyMappingsSourceTables(raw["source_tables"]),
			"source_fields":  buildTableModelRequestBodyMappingsSourceFields(raw["source_fields"]),
		}
		mappings = append(mappings, params)
	}
	return mappings
}

func buildTableModelRequestBodyMappingsSourceTables(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	sourceTables := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		sourceTable := map[string]interface{}{
			"join_type":     utils.ValueIgnoreEmpty(raw["join_type"]),
			"table1_id":     utils.ValueIgnoreEmpty(raw["table1_id"]),
			"table2_id":     utils.ValueIgnoreEmpty(raw["table2_id"]),
			"table1_name":   utils.ValueIgnoreEmpty(raw["table1_name"]),
			"table2_name":   utils.ValueIgnoreEmpty(raw["table2_name"]),
			"source_tables": buildTableModelRequestBodyMappingsSourceTablesJoinFields(raw["source_tables"]),
		}
		sourceTables = append(sourceTables, sourceTable)
	}
	return sourceTables
}

func buildTableModelRequestBodyMappingsSourceTablesJoinFields(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	joinFields := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		joinField := map[string]interface{}{
			"field1_id":   utils.ValueIgnoreEmpty(raw["field1_id"]),
			"field2_id":   utils.ValueIgnoreEmpty(raw["field2_id"]),
			"field1_name": utils.ValueIgnoreEmpty(raw["field1_name"]),
			"field2_name": utils.ValueIgnoreEmpty(raw["field2_name"]),
		}
		joinFields = append(joinFields, joinField)
	}
	return joinFields
}

func buildTableModelRequestBodyMappingsSourceFields(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	sourceFileds := make([]map[string]interface{}, len(rawArray))
	for _, val := range rawArray {
		raw := val.(map[string]interface{})
		sourceField := map[string]interface{}{
			"field_ids":            utils.ValueIgnoreEmpty(raw["field_ids"]),
			"field_names":          utils.ValueIgnoreEmpty(raw["field_names"]),
			"target_field_name":    utils.ValueIgnoreEmpty(raw["target_field_name"]),
			"transform_expression": utils.ValueIgnoreEmpty(raw["transform_expression"]),
		}
		sourceFileds = append(sourceFileds, sourceField)
	}
	return sourceFileds
}

func resourceTableModelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)
	getTableModelHttpUrl := "v2/{project_id}/design/table-model/{id}?latest=true"
	getTableModelProduct := "dataarts"
	getTableModelClient, err := cfg.NewServiceClient(getTableModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	getTableModelPath := getTableModelClient.Endpoint + getTableModelHttpUrl
	getTableModelPath = strings.ReplaceAll(getTableModelPath, "{project_id}", getTableModelClient.ProjectID)
	getTableModelPath = strings.ReplaceAll(getTableModelPath, "{id}", d.Id())
	getTableModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}

	getTableModelResp, err := getTableModelClient.Request("GET", getTableModelPath, &getTableModelOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", tableModelErrCodes...),
			"error retrieving table model")
	}
	getTableModelRespBody, err := utils.FlattenResponse(getTableModelResp)
	if err != nil {
		return diag.FromErr(err)
	}

	tableModel := utils.PathSearch("data.value", getTableModelRespBody, nil)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceID),
		d.Set("model_id", utils.PathSearch("model_id", tableModel, nil)),
		d.Set("dw_type", utils.PathSearch("dw_type", tableModel, nil)),
		d.Set("subject_id", utils.PathSearch("biz_catalog_id", tableModel, nil)),
		d.Set("physical_table_name", utils.PathSearch("tb_name", tableModel, nil)),
		d.Set("table_name", utils.PathSearch("logic_tb_name", tableModel, nil)),
		d.Set("description", utils.PathSearch("description", tableModel, nil)),
		d.Set("configs", utils.PathSearch("configs", tableModel, nil)),
		d.Set("dw_id", utils.PathSearch("dw_id", tableModel, nil)),
		d.Set("parent_table_id", utils.PathSearch("parent_table_id", tableModel, nil)),
		d.Set("related_logic_table_model_id", utils.PathSearch("related_logic_table_id", tableModel, nil)),
		d.Set("related_logic_model_id", utils.PathSearch("related_logic_table_model_id", tableModel, nil)),
		d.Set("owner", utils.PathSearch("owner", tableModel, nil)),
		d.Set("table_type", utils.PathSearch("table_type", tableModel, nil)),
		d.Set("compression", utils.PathSearch("compression", tableModel, nil)),
		d.Set("code", utils.PathSearch("code", tableModel, nil)),
		d.Set("distribute", utils.PathSearch("distribute", tableModel, nil)),
		d.Set("distribute_column", utils.PathSearch("distribute_column", tableModel, nil)),
		d.Set("parent_table_name", utils.PathSearch("parent_table_name", tableModel, nil)),
		d.Set("parent_table_code", utils.PathSearch("parent_table_code", tableModel, nil)),
		d.Set("related_logic_model_name", utils.PathSearch("related_logic_table_model_name", tableModel, nil)),
		d.Set("related_logic_table_model_name", utils.PathSearch("related_logic_table_name", tableModel, nil)),
		d.Set("data_format", utils.PathSearch("data_format", tableModel, nil)),
		d.Set("dlf_task_id", utils.PathSearch("dlf_task_id", tableModel, nil)),
		d.Set("use_recently_partition", utils.PathSearch("use_recently_partition", tableModel, false)),
		d.Set("reversed", utils.PathSearch("reversed", tableModel, false)),
		d.Set("db_name", utils.PathSearch("db_name", tableModel, nil)),
		d.Set("queue_name", utils.PathSearch("queue_name", tableModel, nil)),
		d.Set("schema", utils.PathSearch("schema", tableModel, nil)),
		d.Set("obs_location", utils.PathSearch("obs_location", tableModel, nil)),
		d.Set("dw_name", utils.PathSearch("dw_name", tableModel, nil)),
		d.Set("dirty_out_switch", utils.PathSearch("dirty_out_switch", tableModel, false)),
		d.Set("dirty_out_database", utils.PathSearch("dirty_out_database", tableModel, nil)),
		d.Set("dirty_out_prefix", utils.PathSearch("dirty_out_prefix", tableModel, nil)),
		d.Set("dirty_out_suffix", utils.PathSearch("dirty_out_suffix", tableModel, nil)),
		d.Set("partition_conf", utils.PathSearch("partition_conf", tableModel, nil)),
		d.Set("extend_info", utils.PathSearch("extend_info", tableModel, nil)),
		d.Set("tb_guid", utils.PathSearch("tb_guid", tableModel, nil)),
		d.Set("logic_tb_guid", utils.PathSearch("logic_tb_guid", tableModel, nil)),
		d.Set("status", utils.PathSearch("status", tableModel, nil)),
		d.Set("catalog_path", utils.PathSearch("catalog_path", tableModel, nil)),
		d.Set("created_at", utils.PathSearch("create_time", tableModel, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", tableModel, nil)),
		d.Set("created_by", utils.PathSearch("create_by", tableModel, nil)),
		d.Set("updated_by", utils.PathSearch("update_by", tableModel, nil)),
		d.Set("has_related_physical_table", utils.PathSearch("has_related_physical_table", tableModel, false)),
		d.Set("has_related_logic_table", utils.PathSearch("has_related_logic_table", tableModel, false)),
		d.Set("is_partition", utils.PathSearch("is_partition", tableModel, false)),
		d.Set("attributes", flattenAttributes(utils.PathSearch("attributes", tableModel, make([]interface{}, 0)))),
		d.Set("relations", flattenRelations(utils.PathSearch("relations", tableModel, make([]interface{}, 0)))),
		d.Set("mappings", flattenMappings(utils.PathSearch("mappings", tableModel, make([]interface{}, 0)))),
		d.Set("physical_table_status", utils.PathSearch("physical_table", tableModel, nil)),
		d.Set("dev_physical_table_status", utils.PathSearch("dev_physical_table", tableModel, nil)),
		d.Set("technical_asset_status", utils.PathSearch("technical_asset", tableModel, nil)),
		d.Set("business_asset_status", utils.PathSearch("business_asset", tableModel, nil)),
		d.Set("meta_data_link_status", utils.PathSearch("meta_data_link", tableModel, nil)),
		d.Set("data_quality_status", utils.PathSearch("data_quality", tableModel, nil)),
		d.Set("summary_status", utils.PathSearch("summary_status", tableModel, nil)),
		d.Set("env_type", utils.PathSearch("env_type", tableModel, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting table model fields: %s", err)
	}
	return nil
}

func flattenAttributes(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	i := 1
	for _, val := range rawArray {
		if utils.PathSearch("table_model_id", val, "") != "" {
			continue
		}
		params := map[string]interface{}{
			"name":                       utils.PathSearch("name_ch", val, nil),
			"name_en":                    utils.PathSearch("name_en", val, nil),
			"data_type":                  utils.PathSearch("data_type", val, nil),
			"data_type_extend":           utils.PathSearch("data_type_extend", val, nil),
			"description":                utils.PathSearch("description", val, nil),
			"stand_row_id":               utils.PathSearch("stand_row_id", val, nil),
			"related_logic_attr_id":      utils.PathSearch("related_logic_attr_id", val, nil),
			"is_partition_key":           utils.PathSearch("is_partition_key", val, false),
			"is_primary_key":             utils.PathSearch("is_primary_key", val, false),
			"is_foreign_key":             utils.PathSearch("is_foreign_key", val, false),
			"extend_field":               utils.PathSearch("extend_field", val, false),
			"not_null":                   utils.PathSearch("not_null", val, false),
			"code":                       utils.PathSearch("code", val, nil),
			"id":                         utils.PathSearch("id", val, nil),
			"related_logic_attr_name":    utils.PathSearch("related_logic_attr_name", val, nil),
			"related_logic_attr_name_en": utils.PathSearch("related_logic_attr_name_en", val, nil),
			"stand_row_name":             utils.PathSearch("stand_row_name", val, nil),
			"domain_type":                utils.PathSearch("domain_type", val, nil),
			"created_at":                 utils.PathSearch("create_time", val, nil),
			"updated_at":                 utils.PathSearch("update_time", val, nil),
			"ordinal":                    strconv.Itoa(i),
		}
		i++
		rst = append(rst, params)
	}
	return rst
}

func flattenRelations(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"name":              utils.PathSearch("name", val, nil),
			"source_type":       utils.PathSearch("source_type", val, nil),
			"target_table_id":   utils.PathSearch("target_table_id", val, nil),
			"target_table_name": utils.PathSearch("target_table_name", val, nil),
			"target_type":       utils.PathSearch("target_type", val, nil),
			"source_table_id":   utils.PathSearch("source_table_id", val, nil),
			"source_table_name": utils.PathSearch("source_table_name", val, nil),
			"role":              utils.PathSearch("role", val, nil),
			"id":                utils.PathSearch("id", val, nil),
			"created_at":        utils.PathSearch("create_time", val, nil),
			"updated_at":        utils.PathSearch("update_time", val, nil),
			"created_by":        utils.PathSearch("create_by", val, nil),
			"updated_by":        utils.PathSearch("update_by", val, nil),
			"mappings":          flattenRelationsMappings(utils.PathSearch("mappings", val, make([]interface{}, 0))),
		}
		rst[i] = params
	}
	return rst
}

func flattenRelationsMappings(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"source_field_name": utils.PathSearch("source_field_name", val, nil),
			"target_field_id":   utils.PathSearch("target_field_id", val, nil),
			"target_field_name": utils.PathSearch("target_field_name", val, nil),
			"source_field_id":   utils.PathSearch("source_field_id", val, nil),
			"id":                utils.PathSearch("id", val, nil),
			"created_at":        utils.PathSearch("create_time", val, nil),
			"updated_at":        utils.PathSearch("update_time", val, nil),
			"created_by":        utils.PathSearch("create_by", val, nil),
			"updated_by":        utils.PathSearch("update_by", val, nil),
		}
		rst[i] = params
	}
	return rst
}

func flattenMappings(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"name":           utils.PathSearch("name", val, nil),
			"src_model_id":   utils.PathSearch("src_model_id", val, nil),
			"src_model_name": utils.PathSearch("src_model_name", val, nil),
			"id":             utils.PathSearch("id", val, nil),
			"view_text":      utils.PathSearch("view_text", val, nil),
			"created_at":     utils.PathSearch("create_time", val, nil),
			"updated_at":     utils.PathSearch("update_time", val, nil),
			"created_by":     utils.PathSearch("create_by", val, nil),
			"updated_by":     utils.PathSearch("update_by", val, nil),
			"source_tables":  flattenMappingsSourceTables(utils.PathSearch("source_tables", val, make([]interface{}, 0))),
			"source_fields":  flattenMappingsSourceFields(utils.PathSearch("source_fields", val, make([]interface{}, 0))),
		}
		rst[i] = params
	}
	return rst
}

func flattenMappingsSourceTables(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"join_type":   utils.PathSearch("join_type", val, nil),
			"table1_id":   utils.PathSearch("table1_id", val, nil),
			"table2_id":   utils.PathSearch("table2_id", val, nil),
			"table1_name": utils.PathSearch("table1_name", val, nil),
			"table2_name": utils.PathSearch("table2_name", val, nil),
			"join_fields": flattenMappingsSourceTablesJoinFields(utils.PathSearch("join_fields", val, make([]interface{}, 0))),
		}
		rst[i] = params
	}
	return rst
}

func flattenMappingsSourceTablesJoinFields(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"field1_id":   utils.PathSearch("field1_id", val, nil),
			"field2_id":   utils.PathSearch("field2_id", val, nil),
			"field1_name": utils.PathSearch("field1_name", val, nil),
			"field2_name": utils.PathSearch("field2_name", val, nil),
		}
		rst[i] = params
	}
	return rst
}

func flattenMappingsSourceFields(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"field_ids":            utils.PathSearch("field_ids", val, nil),
			"field_names":          utils.PathSearch("field_names", val, nil),
			"target_field_name":    utils.PathSearch("target_field_name", val, nil),
			"transform_expression": utils.PathSearch("transform_expression", val, nil),
			"changed":              utils.PathSearch("changed", val, false),
		}
		rst[i] = params
	}
	return rst
}

func resourceTableModelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateTableModelHttpUrl := "v2/{project_id}/design/table-model"
	updateTableModelProduct := "dataarts"
	updateTableModelClient, err := cfg.NewServiceClient(updateTableModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	updateTableModelPath := updateTableModelClient.Endpoint + updateTableModelHttpUrl
	updateTableModelPath = strings.ReplaceAll(updateTableModelPath, "{project_id}", updateTableModelClient.ProjectID)

	updateBody := utils.RemoveNil(buildCreateOrUpdateTableModelBodyParams(d))
	updateBody["id"] = d.Id()

	updateTableModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         updateBody,
	}
	_, err = updateTableModelClient.Request("PUT", updateTableModelPath, &updateTableModelOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTableModelRead(ctx, d, meta)
}

func resourceTableModelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	deleteTableModelHttpUrl := "v2/{project_id}/design/table-model"
	deleteTableModelProduct := "dataarts"
	deleteTableModelClient, err := cfg.NewServiceClient(deleteTableModelProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	deleteTableModelPath := deleteTableModelClient.Endpoint + deleteTableModelHttpUrl
	deleteTableModelPath = strings.ReplaceAll(deleteTableModelPath, "{project_id}", deleteTableModelClient.ProjectID)
	deleteTableModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"ids":      []string{d.Id()},
			"del_type": utils.ValueIgnoreEmpty(d.Get("del_type")),
		}),
	}

	deleteTableModelResp, err := deleteTableModelClient.Request("DELETE", deleteTableModelPath, &deleteTableModelOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	deleteTableModelRespBody, err := utils.FlattenResponse(deleteTableModelResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// if successfully deleting table model, data.value in return will be 1, else will be 0
	if utils.PathSearch("data.value", deleteTableModelRespBody, 0) != 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error deleting table model")
	}

	return nil
}
