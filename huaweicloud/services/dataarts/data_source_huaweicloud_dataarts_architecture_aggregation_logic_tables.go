package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/aggregation-logic-tables
func DataSourceArchitectureAggregationLogicTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureAggregationLogicTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the aggregation logic tables are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the aggregation logic tables belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name or English name of the aggregation logic table to be fuzzy queried.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the aggregation logic table to be exactly queried.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The English name of the aggregation logic table to be exactly queried.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the aggregation logic table to be queried.`,
			},
			"approver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The approver of the aggregation logic table to be queried.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The owner of the aggregation logic table to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The publishing status of the aggregation logic table to be queried.`,
			},
			"sync_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The synchronization status of the aggregation logic table to be queried.`,
			},
			"sync_key": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of synchronization task types of the aggregation logic table.`,
			},
			"biz_catalog_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business catalog ID to which the aggregation logic table belongs.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the modification time for the aggregation logic table, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the modification time for the aggregation logic table, in RFC3339 format.`,
			},
			"auto_generate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the aggregation logic table is auto-generated.`,
			},

			// Attributes.
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the aggregation logic table.`,
						},
						"tb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The physical table name in English of the aggregation logic table.`,
						},
						"tb_logic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name in Chinese of the aggregation logic table.`,
						},
						"dw_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data connection.`,
						},
						"dw_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data connection.`,
						},
						"dw_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the data connection.`,
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The asset owner of the aggregation logic table.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the aggregation logic table.`,
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alias of the aggregation logic table.`,
						},
						"queue_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The queue name of the DLI data connection.`,
						},
						"schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The schema name of the aggregation logic table.`,
						},
						"table_attributes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        architectureAggregationLogicTablesAttributeSchema(),
							Description: `The list of attributes of the aggregation logic table.`,
						},
						"table_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the aggregation logic table.`,
						},
						"distribute": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distribution mode of the database.`,
						},
						"distribute_column": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The column used for hash distribution.`,
						},
						"compression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The compression level of the DWS table.`,
						},
						"pre_combine_field": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The column used for table combine or versioning.`,
						},
						"obs_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS storage path of the external table.`,
						},
						"configs": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The additional configuration of the aggregation logic table, in JSON format.`,
						},
						"dimension_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the statistical dimension of the derived indicator.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dimension group of the derivative metric.`,
						},
						"group_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dimension group code of the derivative metric.`,
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL statement of the aggregation logic table.`,
						},
						"partition_conf": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The partition expression of the aggregation logic table.`,
						},
						"dirty_out_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable dirty data output switch.`,
						},
						"dirty_out_database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The output database of the dirty data.`,
						},
						"dirty_out_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The prefix of the dirty data table.`,
						},
						"dirty_out_suffix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The suffix of the dirty data table.`,
						},
						"secret_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the data secrecy level.`,
						},
						"self_defined_fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        architectureAggregationLogicTablesSelfDefinedFieldSchema(),
							Description: `The list of user-defined extended fields for the aggregation logic table.`,
						},
						"tb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The internal table ID of the aggregation logic table.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The publishing status of the aggregation logic table.`,
						},
						"model_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the model to which the aggregation logic table belongs.`,
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the aggregation logic table.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the aggregation logic table, in RFC3339 format.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the aggregation logic table, in RFC3339 format.`,
						},
						"env_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The publishing environment type of the aggregation logic table.`,
						},
						"physical_table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the table creation in the production environment.`,
						},
						"dev_physical_table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the table creation in the development environment.`,
						},
						"technical_asset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization status of the technical asset.`,
						},
						"business_asset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization status of the business asset.`,
						},
						"meta_data_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the asset association.`,
						},
						"tb_guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The technical catalog asset GUID after the aggregation logic table is published.`,
						},
						"tb_logic_guid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The business catalog asset GUID after the aggregation logic table is published.`,
						},
						"data_quality": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the data quality job creation.`,
						},
						"quality_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data quality.`,
						},
						"dlf_task": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the DLF task.`,
						},
						"dlf_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the DLF task.`,
						},
						"publish_to_dlm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the DLM API generation.`,
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API after publishing.`,
						},
						"summary_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The synchronization status of the summary.`,
						},
						"reversed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the aggregation logic table is reversed.`,
						},
						"table_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The version of the aggregation logic table.`,
						},
						"dev_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The development environment version of the aggregation logic table.`,
						},
						"prod_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The production environment version of the aggregation logic table.`,
						},
						"dev_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The development environment version name of the aggregation logic table.`,
						},
						"prod_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The production environment version name of the aggregation logic table.`,
						},
						"approval_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        architectureAggregationLogicTablesApprovalInfoSchema(),
							Description: `The approval information of the aggregation logic table.`,
						},
					},
				},
				Description: `The list of aggregation logic tables that match the filter parameters.`,
			},
		},
	}
}

func architectureAggregationLogicTablesAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the attribute.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the attribute.`,
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
			"attribute_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The configuration type of the attribute.`,
			},
			"is_primary_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is the primary key.`,
			},
			"is_partition_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is used as a partition key.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is not null.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the attribute.`,
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
			"ref_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the object referenced by the attribute.`,
			},
			"ref_name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the object associated with the attribute.`,
			},
			"ref_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the object associated with the attribute.`,
			},
			"stand_row_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data standard associated with the attribute.`,
			},
			"stand_row_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data standard associated with the attribute.`,
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The sequence number of the attribute.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alias of the attribute.`,
			},
			"secrecy_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        architectureAggregationLogicTablesAttributeSecrecyLevelSchema(),
				Description: `The list of secrecy levels associated with the attribute.`,
			},
		},
	}
}

func architectureAggregationLogicTablesAttributeSecrecyLevelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the secrecy level.`,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the secrecy level.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the secrecy level.`,
			},
			"slevel": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The secrecy level number.`,
			},
		},
	}
}

func architectureAggregationLogicTablesSelfDefinedFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese display name of the custom extended field.`,
			},
			"fd_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the custom extended field.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the custom extended field requires a value.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the custom extended field.`,
			},
		},
	}
}

func architectureAggregationLogicTablesApprovalInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the approval for the aggregation logic table.`,
			},
			"approver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approver of the aggregation logic table.`,
			},
			"approval_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval status of the aggregation logic table.`,
			},
			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval message for the aggregation logic table.`,
			},
			"approval_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval time for the aggregation logic table, in RFC3339 format.`,
			},
		},
	}
}

func buildArchitectureAggregationLogicTablesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("name_ch"); ok {
		res = fmt.Sprintf("%s&name_ch=%v", res, v)
	}
	if v, ok := d.GetOk("name_en"); ok {
		res = fmt.Sprintf("%s&name_en=%v", res, v)
	}

	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}

	if v, ok := d.GetOk("approver"); ok {
		res = fmt.Sprintf("%s&approver=%v", res, v)
	}

	if v, ok := d.GetOk("owner"); ok {
		res = fmt.Sprintf("%s&owner=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("sync_status"); ok {
		res = fmt.Sprintf("%s&sync_status=%v", res, v)
	}

	if v, ok := d.GetOk("sync_key"); ok {
		for _, item := range v.([]interface{}) {
			res = fmt.Sprintf("%s&sync_key=%v", res, item)
		}
	}

	if v, ok := d.GetOk("biz_catalog_id"); ok {
		res = fmt.Sprintf("%s&biz_catalog_id=%v", res, v)
	}

	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%v", res, v)
	}

	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	if d.Get("auto_generate").(bool) {
		res = fmt.Sprintf("%s&auto_generate=true", res)
	}

	return res
}

func listArchitectureAggregationLogicTables(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/aggregation-logic-tables?limit={limit}"
		offset  = 0
		// Defaults to `50`. The maximum value is `100`.
		limit  = 100
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureAggregationLogicTablesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		tables := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tables...)

		if len(tables) < limit {
			break
		}

		offset += len(tables)
	}

	return result, nil
}

func dataSourceArchitectureAggregationLogicTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	tables, err := listArchitectureAggregationLogicTables(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture aggregation logic tables: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tables", flattenArchitectureAggregationLogicTables(tables)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenArchitectureAggregationLogicTables(tables []interface{}) []interface{} {
	if len(tables) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", table, nil),
			"tb_name":       utils.PathSearch("tb_name", table, nil),
			"tb_logic_name": utils.PathSearch("tb_logic_name", table, nil),
			"dw_id":         utils.PathSearch("dw_id", table, nil),
			"dw_name":       utils.PathSearch("dw_name", table, nil),
			"dw_type":       utils.PathSearch("dw_type", table, nil),
			"db_name":       utils.PathSearch("db_name", table, nil),
			"owner":         utils.PathSearch("owner", table, nil),
			"description":   utils.PathSearch("description", table, nil),
			"alias":         utils.PathSearch("alias", table, nil),
			"queue_name":    utils.PathSearch("queue_name", table, nil),
			"schema":        utils.PathSearch("schema", table, nil),
			"table_attributes": flattenArchitectureAggregationLogicTablesAttributes(
				utils.PathSearch("table_attributes", table, make([]interface{}, 0)).([]interface{})),
			"table_type":         utils.PathSearch("table_type", table, nil),
			"distribute":         utils.PathSearch("distribute", table, nil),
			"distribute_column":  utils.PathSearch("distribute_column", table, nil),
			"compression":        utils.PathSearch("compression", table, nil),
			"pre_combine_field":  utils.PathSearch("pre_combine_field", table, nil),
			"obs_location":       utils.PathSearch("obs_location", table, nil),
			"configs":            utils.PathSearch("configs", table, nil),
			"dimension_group":    utils.PathSearch("dimension_group", table, nil),
			"group_name":         utils.PathSearch("group_name", table, nil),
			"group_code":         utils.PathSearch("group_code", table, nil),
			"sql":                utils.PathSearch("sql", table, nil),
			"partition_conf":     utils.PathSearch("partition_conf", table, nil),
			"dirty_out_switch":   utils.PathSearch("dirty_out_switch", table, false),
			"dirty_out_database": utils.PathSearch("dirty_out_database", table, nil),
			"dirty_out_prefix":   utils.PathSearch("dirty_out_prefix", table, nil),
			"dirty_out_suffix":   utils.PathSearch("dirty_out_suffix", table, nil),
			"secret_type":        utils.PathSearch("secret_type", table, nil),
			"self_defined_fields": flattenArchitectureAggregationLogicTablesSelfDefinedFields(
				utils.PathSearch("self_defined_fields", table, make([]interface{}, 0)).([]interface{})),
			"tb_id":              utils.PathSearch("tb_id", table, nil),
			"status":             utils.PathSearch("status", table, nil),
			"model_id":           utils.PathSearch("model_id", table, nil),
			"create_by":          utils.PathSearch("create_by", table, nil),
			"create_time":        utils.PathSearch("create_time", table, nil),
			"update_time":        utils.PathSearch("update_time", table, nil),
			"env_type":           utils.PathSearch("env_type", table, nil),
			"physical_table":     utils.PathSearch("physical_table", table, nil),
			"dev_physical_table": utils.PathSearch("dev_physical_table", table, nil),
			"technical_asset":    utils.PathSearch("technical_asset", table, nil),
			"business_asset":     utils.PathSearch("business_asset", table, nil),
			"meta_data_link":     utils.PathSearch("meta_data_link", table, nil),
			"tb_guid":            utils.PathSearch("tb_guid", table, nil),
			"tb_logic_guid":      utils.PathSearch("tb_logic_guid", table, nil),
			"data_quality":       utils.PathSearch("data_quality", table, nil),
			"quality_id":         utils.PathSearch("quality_id", table, nil),
			"dlf_task":           utils.PathSearch("dlf_task", table, nil),
			"dlf_task_id":        utils.PathSearch("dlf_task_id", table, nil),
			"publish_to_dlm":     utils.PathSearch("publish_to_dlm", table, nil),
			"api_id":             utils.PathSearch("api_id", table, nil),
			"summary_status":     utils.PathSearch("summary_status", table, nil),
			"reversed":           utils.PathSearch("reversed", table, false),
			"table_version":      utils.PathSearch("table_version", table, 0),
			"dev_version":        utils.PathSearch("dev_version", table, nil),
			"prod_version":       utils.PathSearch("prod_version", table, nil),
			"dev_version_name":   utils.PathSearch("dev_version_name", table, nil),
			"prod_version_name":  utils.PathSearch("prod_version_name", table, nil),
			"approval_info": flattenArchitectureAggregationLogicTablesApprovalInfo(utils.PathSearch("approval_info",
				table, nil)),
		})
	}

	return result
}

func flattenArchitectureAggregationLogicTablesAttributes(attributes []interface{}) []interface{} {
	if len(attributes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(attributes))
	for _, attr := range attributes {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", attr, nil),
			"name_ch":          utils.PathSearch("name_ch", attr, nil),
			"name_en":          utils.PathSearch("name_en", attr, nil),
			"data_type":        utils.PathSearch("data_type", attr, nil),
			"attribute_type":   utils.PathSearch("attribute_type", attr, nil),
			"is_primary_key":   utils.PathSearch("is_primary_key", attr, nil),
			"is_partition_key": utils.PathSearch("is_partition_key", attr, nil),
			"not_null":         utils.PathSearch("not_null", attr, nil),
			"description":      utils.PathSearch("description", attr, nil),
			"data_type_extend": utils.PathSearch("data_type_extend", attr, nil),
			"domain_type":      utils.PathSearch("domain_type", attr, nil),
			"ref_id":           utils.PathSearch("ref_id", attr, nil),
			"ref_name_ch":      utils.PathSearch("ref_name_ch", attr, nil),
			"ref_name_en":      utils.PathSearch("ref_name_en", attr, nil),
			"stand_row_id":     utils.PathSearch("stand_row_id", attr, nil),
			"stand_row_name":   utils.PathSearch("stand_row_name", attr, nil),
			"ordinal":          utils.PathSearch("ordinal", attr, 0),
			"alias":            utils.PathSearch("alias", attr, nil),
			"secrecy_levels": flattenArchitectureAggregationLogicTablesAttributeSecrecyLevels(
				utils.PathSearch("secrecy_levels", attr, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenArchitectureAggregationLogicTablesAttributeSecrecyLevels(levels []interface{}) []interface{} {
	if len(levels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(levels))
	for _, level := range levels {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("secrecyLevel_id", level, nil),
			"uuid":   utils.PathSearch("uuid", level, nil),
			"name":   utils.PathSearch("secrecyLevel_name", level, nil),
			"slevel": utils.PathSearch("slevel", level, nil),
		})
	}
	return result
}

func flattenArchitectureAggregationLogicTablesSelfDefinedFields(fields []interface{}) []interface{} {
	if len(fields) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		result = append(result, map[string]interface{}{
			"fd_name_ch": utils.PathSearch("fd_name_ch", field, nil),
			"fd_name_en": utils.PathSearch("fd_name_en", field, nil),
			"not_null":   utils.PathSearch("not_null", field, nil),
			"fd_value":   utils.PathSearch("fd_value", field, nil),
		})
	}
	return result
}

func flattenArchitectureAggregationLogicTablesApprovalInfo(approval interface{}) []interface{} {
	if approval == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":              utils.PathSearch("id", approval, nil),
			"approver":        utils.PathSearch("approver", approval, nil),
			"approval_status": utils.PathSearch("approval_status", approval, nil),
			"msg":             utils.PathSearch("msg", approval, nil),
			"approval_time":   utils.PathSearch("approval_time", approval, nil),
		},
	}
}
