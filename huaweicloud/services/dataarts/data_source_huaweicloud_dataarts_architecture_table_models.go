package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/table-model
func DataSourceArchitectureTableModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceArchitectureTableModelsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID of DataArts Architecture.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the model ID to which the table model belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the Chinese or English name of the table model.`,
			},
			"subject_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the subject ID to which the table model belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the table model.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the creator of the table model.`,
			},
			"tables": {
				Type:        schema.TypeList,
				Elem:        TableModelSchema(),
				Computed:    true,
				Description: `All table models that match the filter parameters.`,
			},
		},
	}
}

func TableModelSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the table model.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The model ID corresponding to the table model.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the table model.`,
			},
			"physical_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the table model.`,
			},
			"dw_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the connection corresponding to the table model.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the connection corresponding to the table model.`,
			},
			"subject_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject ID corresponding to the table model.`,
			},
			"catalog_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject path corresponding to the table model.`,
			},
			"table_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The table type of the table model.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the table model.`,
			},
			"configs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The advanced configuration information of the table model, in JSON format.`,
			},
			"parent_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent table ID of table model.`,
			},
			"parent_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent table name of table model.`,
			},
			"parent_table_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent table code of table model.`,
			},
			"related_logic_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical table ID associated with table model.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The owner of the table model.`,
			},
			"compression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The compression type of the table model.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The database name of the table model.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The queue name of the DLI table model.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schema of the DWS table model.`,
			},
			"obs_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS path of the table model.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the attributes of the table model.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the attribute.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the attribute.`,
						},
						"name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The English name of the attribute.`,
						},
						"data_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type of the attribute.`,
						},
						"data_type_extend": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type extend field of attribute.`,
						},
						"domain_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain type of the attribute.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the attribute.`,
						},
						"stand_row_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data standard associated with attribute.`,
						},
						"ordinal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The sequence number of attribute.`,
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The code of the logical attribute associated with attribute.`,
						},
						"extend_field": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `The extend field of the attribute.`,
						},
						"is_foreign_key": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the attribute is foreign key.`,
						},
						"is_primary_key": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the attribute is primary key.`,
						},
						"is_partition_key": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the attribute is partition key.`,
						},
						"not_null": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the attribute is not null.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the attribute, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the attribute, in RFC3339 format.`,
						},
					},
				},
			},
			"distribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attribute distribution mode of the DWS table model.`,
			},
			"distribute_column": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HASH column of the attribute distribution.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code of the logical entity.`,
			},
			"data_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data format of the DLI table model.`,
			},
			"dlf_task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DLF task ID of table model.`,
			},
			"use_recently_partition": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the table model has used latest partition.`,
			},
			"dirty_out_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The dirty data output switch of the table model.`,
			},
			"dirty_out_database": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The database where to record the dirty data.`,
			},
			"dirty_out_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The prefix of the table recording dirty data.`,
			},
			"dirty_out_suffix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The suffix of the table recording dirty data.`,
			},
			"partition_conf": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The condition expression of the partition.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the table model.`,
			},
			"extend_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extend information of the table model.`,
			},
			"is_partition": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether table is the partition table.`,
			},
			"tb_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The globally unique ID generated when publishing the table model.`,
			},
			"logic_tb_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The globally unique ID of the logic table model generated when publishing the table model.`,
			},
			"logic_tb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of logical entity.`,
			},
			"has_related_logic_table": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the table model has associated the logical entities.`,
			},
			"has_related_physical_table": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the logical entity has associated the physical tables.`,
			},
			"tb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data table.`,
			},
			"physical_table_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The physical table status of the table model.`,
			},
			"dev_physical_table_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dev physical table status of the table model.`,
			},
			"technical_asset_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The technical asset status of the table model.`,
			},
			"business_asset_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business asset status of the table model.`,
			},
			"meta_data_link_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The meta data link status the table model.`,
			},
			"data_quality_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data quality status of the table model.`,
			},
			"summary_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The summary status of the table model.`,
			},
			"prod_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The production environment version of the table model.`,
			},
			"dev_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The development environment version of the table model.`,
			},
			"env_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The environment type of the table model.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the table model.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest updater of the table model.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the table model, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the table model, in RFC3339 format.`,
			},
		},
	}
	return &sc
}

func queryArchitectureTableModels(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/table-model?limit=100&model_id={model_id}"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{model_id}", d.Get("model_id").(string))
	queryParams := buildTableModelsQueryParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving DataArts table models: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		tableModels := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(tableModels) < 1 {
			break
		}
		result = append(result, tableModels...)
		offset += len(tableModels)
	}
	return result, nil
}

func buildTableModelsQueryParams(d *schema.ResourceData) string {
	res := ""
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name)
	}
	if subjectId, ok := d.GetOk("subject_id"); ok {
		res = fmt.Sprintf("%s&biz_catalog_id=%v", res, subjectId)
	}
	if status, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, status)
	}

	if createdBy, ok := d.GetOk("created_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, createdBy)
	}

	return res
}

func resourceArchitectureTableModelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	tableModels, err := queryArchitectureTableModels(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tables", flattenGetArchitectureTableModels(tableModels)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetArchitectureTableModels(all []interface{}) []interface{} {
	if all == nil {
		return nil
	}

	result := make([]interface{}, 0, len(all))
	for _, table := range all {
		result = append(result, map[string]interface{}{
			"id":                         utils.PathSearch("id", table, nil),
			"model_id":                   utils.PathSearch("model_id", table, nil),
			"table_name":                 utils.PathSearch("logic_tb_name", table, nil),
			"physical_table_name":        utils.PathSearch("tb_name", table, nil),
			"dw_id":                      utils.PathSearch("dw_id", table, nil),
			"dw_type":                    utils.PathSearch("dw_type", table, nil),
			"subject_id":                 utils.PathSearch("biz_catalog_id", table, nil),
			"catalog_path":               utils.PathSearch("catalog_path", table, nil),
			"table_type":                 utils.PathSearch("table_type", table, nil),
			"description":                utils.PathSearch("description", table, nil),
			"configs":                    utils.PathSearch("configs", table, nil),
			"parent_table_id":            utils.PathSearch("parent_table_id", table, nil),
			"parent_table_name":          utils.PathSearch("parent_table_name", table, nil),
			"parent_table_code":          utils.PathSearch("parent_table_code", table, nil),
			"related_logic_table_id":     utils.PathSearch("related_logic_table_id", table, nil),
			"owner":                      utils.PathSearch("owner", table, nil),
			"compression":                utils.PathSearch("compression", table, nil),
			"db_name":                    utils.PathSearch("db_name", table, nil),
			"queue_name":                 utils.PathSearch("queue_name", table, nil),
			"schema":                     utils.PathSearch("schema", table, nil),
			"obs_location":               utils.PathSearch("obs_location", table, nil),
			"attributes":                 flattenTablesAttrs(utils.PathSearch("attributes", table, make([]interface{}, 0)).([]interface{})),
			"distribute":                 utils.PathSearch("distribute", table, nil),
			"distribute_column":          utils.PathSearch("distribute_column", table, nil),
			"code":                       utils.PathSearch("code", table, nil),
			"data_format":                utils.PathSearch("data_format", table, nil),
			"dlf_task_id":                utils.PathSearch("dlf_task_id", table, nil),
			"use_recently_partition":     utils.PathSearch("use_recently_partition", table, false),
			"dirty_out_switch":           utils.PathSearch("dirty_out_switch", table, false),
			"dirty_out_database":         utils.PathSearch("dirty_out_database", table, nil),
			"dirty_out_prefix":           utils.PathSearch("dirty_out_prefix", table, nil),
			"dirty_out_suffix":           utils.PathSearch("dirty_out_suffix", table, nil),
			"partition_conf":             utils.PathSearch("partition_conf", table, nil),
			"status":                     utils.PathSearch("status", table, nil),
			"extend_info":                utils.PathSearch("extend_info", table, nil),
			"is_partition":               utils.PathSearch("is_partition", table, false),
			"tb_guid":                    utils.PathSearch("tb_guid", table, nil),
			"logic_tb_guid":              utils.PathSearch("logic_tb_guid", table, nil),
			"logic_tb_id":                utils.PathSearch("logic_tb_id", table, nil),
			"has_related_physical_table": utils.PathSearch("has_related_physical_table", table, false),
			"has_related_logic_table":    utils.PathSearch("has_related_logic_table", table, false),
			"tb_id":                      utils.PathSearch("tb_id", table, nil),
			"physical_table_status":      utils.PathSearch("physical_table", table, nil),
			"dev_physical_table_status":  utils.PathSearch("dev_physical_table", table, nil),
			"technical_asset_status":     utils.PathSearch("technical_asset", table, nil),
			"business_asset_status":      utils.PathSearch("business_asset", table, nil),
			"meta_data_link_status":      utils.PathSearch("meta_data_link", table, nil),
			"data_quality_status":        utils.PathSearch("data_quality", table, nil),
			"summary_status":             utils.PathSearch("summary_status", table, nil),
			"prod_version":               utils.PathSearch("prod_version", table, nil),
			"dev_version":                utils.PathSearch("dev_version", table, nil),
			"env_type":                   utils.PathSearch("env_type", table, nil),
			"created_by":                 utils.PathSearch("create_by", table, nil),
			"updated_by":                 utils.PathSearch("update_by", table, nil),
			"created_at":                 formtTimeToRFC3339(utils.PathSearch("create_time", table, "").(string)),
			"updated_at":                 formtTimeToRFC3339(utils.PathSearch("update_time", table, "").(string)),
		})
	}
	return result
}

func flattenTablesAttrs(attrs []interface{}) []map[string]interface{} {
	if len(attrs) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(attrs))
	for i, val := range attrs {
		rst[i] = map[string]interface{}{
			"id":               utils.PathSearch("id", val, nil),
			"name":             utils.PathSearch("name_ch", val, nil),
			"name_en":          utils.PathSearch("name_en", val, nil),
			"data_type":        utils.PathSearch("data_type", val, nil),
			"data_type_extend": utils.PathSearch("data_type_extend", val, nil),
			"domain_type":      utils.PathSearch("domain_type", val, nil),
			"description":      utils.PathSearch("description", val, nil),
			"stand_row_id":     utils.PathSearch("stand_row_id", val, nil),
			"ordinal":          utils.PathSearch("ordinal", val, 0),
			"code":             utils.PathSearch("code", val, nil),
			"extend_field":     utils.PathSearch("extend_field", val, false),
			"is_foreign_key":   utils.PathSearch("is_foreign_key", val, false),
			"is_primary_key":   utils.PathSearch("is_primary_key", val, false),
			"is_partition_key": utils.PathSearch("is_partition_key", val, false),
			"not_null":         utils.PathSearch("not_null", val, false),
			"created_at":       formtTimeToRFC3339(utils.PathSearch("create_time", val, "").(string)),
			"updated_at":       formtTimeToRFC3339(utils.PathSearch("update_time", val, "").(string)),
		}
	}

	return rst
}

// The timeStr formet is "2024-08-01T15:21:58+08:00".
func formtTimeToRFC3339(timeStr string) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(timeStr)/1000, false)
}
