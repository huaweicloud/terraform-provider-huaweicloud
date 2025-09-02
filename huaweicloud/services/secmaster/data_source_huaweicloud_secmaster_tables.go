package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/tables
func DataSourceTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the directory type.",
			},
			"table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the table ID.",
			},
			"table_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the table alias.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the table name.",
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the attribute fields for sorting.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sorting order. Supported values are **ASC** and **DESC**.",
			},
			"exists": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the table exists.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The tables list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project ID.",
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The workspace ID.",
						},
						"table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table ID.",
						},
						"pipe_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The pipe ID.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table name.",
						},
						"table_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table alias.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table description.",
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The directory group.",
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The directory type.",
						},
						"lock_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table lock status.",
						},
						"process_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The processing status.",
						},
						"process_error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table processing error.",
						},
						"format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table format.",
						},
						"rw_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table read/write type.",
						},
						"owner_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner type.",
						},
						"data_layering": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data layering.",
						},
						"data_classification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data classification.",
						},
						"schema": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The table schema.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"columns": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The table columns list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The table column name.",
												},
												"column_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The column field type.",
												},
												"column_type_setting": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The table column type setting.",
												},
												"column_data_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The column field data type.",
												},
												"column_data_type_setting": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The table column data type setting.",
												},
												"nullable": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the column is nullable.",
												},
												"array": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the column is an array.",
												},
												"depth": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The depth.",
												},
												"parent_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The parent name.",
												},
												"own_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The own name.",
												},
												"column_display_setting": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The table column display setting.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mapping_required": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether mapping is required.",
															},
															"group_sequence_number": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The group sequence number.",
															},
															"intra_group_sequence_number": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The intra-group sequence number.",
															},
															"value_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The value type.",
															},
															"value_qualified": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The qualified value.",
															},
															"display_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The display name.",
															},
															"display_description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The display description.",
															},
															"group_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The group name.",
															},
														},
													},
												},
												"column_sequence_number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The column sequence number.",
												},
											},
										},
									},
									"primary_key": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The table primary key list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"partition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The table partition list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"watermark_column": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The table watermark column.",
									},
									"watermark_interval": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The table watermark delay interval.",
									},
									"time_filter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The table time filter column.",
									},
								},
							},
						},
						"storage_setting": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The table storage setting.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_index": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The application index.",
									},
									"application_topic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The application topic.",
									},
									"application_data_class_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The application data class ID.",
									},
									"streaming_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The streaming bandwidth (MB/s).",
									},
									"streaming_partition": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The streaming partition.",
									},
									"streaming_retention_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The streaming retention size.",
									},
									"streaming_dataspace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The streaming dataspace ID.",
									},
									"index_storage_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The index storage period.",
									},
									"index_storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The index storage size.",
									},
									"index_shards": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The index shards number.",
									},
									"index_replicas": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The index replicas number.",
									},
									"lake_storage_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The data lake storage period.",
									},
									"lake_partition_setting": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The data lake partition setting.",
									},
									"lake_expiration_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expired status of data lake.",
									},
								},
							},
						},
						"display_setting": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The table display setting.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"columns": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The application index.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"column_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The table column name.",
												},
												"column_alias": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The table column alias.",
												},
												"display_by_default": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Is it displayed by default.",
												},
											},
										},
									},
									"format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The table display settings.",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time (timestamp in milliseconds).",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time (timestamp in milliseconds).",
						},
						"delete_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The deletion time (timestamp in milliseconds).",
						},
					},
				},
			},
		},
	}
}

func buildTablesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("table_id"); ok {
		queryParams = fmt.Sprintf("%s&table_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("table_alias"); ok {
		queryParams = fmt.Sprintf("%s&table_alias=%v", queryParams, v)
	}
	if v, ok := d.GetOk("table_name"); ok {
		queryParams = fmt.Sprintf("%s&table_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}
	queryParams = fmt.Sprintf("%s&exists=%v", queryParams, d.Get("exists").(bool))

	return queryParams
}

func dataSourceTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/siem/tables"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath += buildTablesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster tables: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", requestRespBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 {
			break
		}

		result = append(result, recordsResp...)
		offset += len(recordsResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenTablesRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTablesRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"project_id":          utils.PathSearch("project_id", v, nil),
			"workspace_id":        utils.PathSearch("workspace_id", v, nil),
			"table_id":            utils.PathSearch("table_id", v, nil),
			"pipe_id":             utils.PathSearch("pipe_id", v, nil),
			"table_name":          utils.PathSearch("table_name", v, nil),
			"table_alias":         utils.PathSearch("table_alias", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"directory":           utils.PathSearch("directory", v, nil),
			"category":            utils.PathSearch("category", v, nil),
			"lock_status":         utils.PathSearch("lock_status", v, nil),
			"process_status":      utils.PathSearch("process_status", v, nil),
			"process_error":       utils.PathSearch("process_error", v, nil),
			"format":              utils.PathSearch("format", v, nil),
			"rw_type":             utils.PathSearch("rw_type", v, nil),
			"owner_type":          utils.PathSearch("owner_type", v, nil),
			"data_layering":       utils.PathSearch("data_layering", v, nil),
			"data_classification": utils.PathSearch("data_classification", v, nil),
			"schema":              flattenTablesSchema(utils.PathSearch("schema", v, nil)),
			"storage_setting":     flattenTablesStorageSetting(utils.PathSearch("storage_setting", v, nil)),
			"display_setting":     flattenTablesDisplaySetting(utils.PathSearch("display_setting", v, nil)),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"update_time":         utils.PathSearch("update_time", v, nil),
			"delete_time":         utils.PathSearch("delete_time", v, nil),
		})
	}

	return rst
}

func flattenTablesSchema(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"columns":            flattenTablesSchemaColumns(utils.PathSearch("columns", resp, make([]interface{}, 0)).([]interface{})),
			"primary_key":        utils.ExpandToStringList(utils.PathSearch("primary_key", resp, make([]interface{}, 0)).([]interface{})),
			"partition":          utils.ExpandToStringList(utils.PathSearch("partition", resp, make([]interface{}, 0)).([]interface{})),
			"watermark_column":   utils.PathSearch("watermark_column", resp, nil),
			"watermark_interval": utils.PathSearch("watermark_interval", resp, nil),
			"time_filter":        utils.PathSearch("time_filter", resp, nil),
		},
	}
}

func flattenTablesSchemaColumns(columnsResp []interface{}) []interface{} {
	if len(columnsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columnsResp))
	for _, v := range columnsResp {
		rst = append(rst, map[string]interface{}{
			"column_name":              utils.PathSearch("column_name", v, nil),
			"column_type":              utils.PathSearch("column_type", v, nil),
			"column_type_setting":      utils.PathSearch("column_type_setting", v, nil),
			"column_data_type":         utils.PathSearch("column_data_type", v, nil),
			"column_data_type_setting": utils.PathSearch("column_data_type_setting", v, nil),
			"nullable":                 utils.PathSearch("nullable", v, nil),
			"array":                    utils.PathSearch("array", v, nil),
			"depth":                    utils.PathSearch("depth", v, nil),
			"parent_name":              utils.PathSearch("parent_name", v, nil),
			"own_name":                 utils.PathSearch("own_name", v, nil),
			"column_display_setting":   flattenTablesColumnDisplaySetting(utils.PathSearch("column_display_setting", v, nil)),
			"column_sequence_number":   utils.PathSearch("column_sequence_number", v, nil),
		})
	}

	return rst
}

func flattenTablesColumnDisplaySetting(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"mapping_required":            utils.PathSearch("mapping_required", resp, nil),
			"group_sequence_number":       utils.PathSearch("group_sequence_number", resp, nil),
			"intra_group_sequence_number": utils.PathSearch("intra_group_sequence_number", resp, nil),
			"value_type":                  utils.PathSearch("value_type", resp, nil),
			"value_qualified":             utils.PathSearch("value_qualified", resp, nil),
			"display_name":                utils.PathSearch("display_name", resp, nil),
			"display_description":         utils.PathSearch("display_description", resp, nil),
			"group_name":                  utils.PathSearch("group_name", resp, nil),
		},
	}
}

func flattenTablesStorageSetting(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"application_index":         utils.PathSearch("application_index", resp, nil),
			"application_topic":         utils.PathSearch("application_topic", resp, nil),
			"application_data_class_id": utils.PathSearch("application_data_class_id", resp, nil),
			"streaming_bandwidth":       utils.PathSearch("streaming_bandwidth", resp, nil),
			"streaming_partition":       utils.PathSearch("streaming_partition", resp, nil),
			"streaming_retention_size":  utils.PathSearch("streaming_retention_size", resp, nil),
			"streaming_dataspace_id":    utils.PathSearch("streaming_dataspace_id", resp, nil),
			"index_storage_period":      utils.PathSearch("index_storage_period", resp, nil),
			"index_storage_size":        utils.PathSearch("index_storage_size", resp, nil),
			"index_shards":              utils.PathSearch("index_shards", resp, nil),
			"index_replicas":            utils.PathSearch("index_replicas", resp, nil),
			"lake_storage_period":       utils.PathSearch("lake_storage_period", resp, nil),
			"lake_partition_setting":    utils.PathSearch("lake_partition_setting", resp, nil),
			"lake_expiration_status":    utils.PathSearch("lake_expiration_status", resp, nil),
		},
	}
}

func flattenTablesDisplaySetting(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	// Return an empty list structure for display_setting since the API docs don't specify the exact fields
	return []interface{}{
		map[string]interface{}{},
	}
}
