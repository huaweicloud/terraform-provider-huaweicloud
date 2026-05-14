package dataarts

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	architectureAggregationLogicTableErrCodes = []string{
		"DLG.0818", // Workspace not found.
		"DLG.6018", // Resource not found.
	}
	architectureAggregationLogicTableNonUpdatableParams = []string{
		"workspace_id",
	}
)

// @API DataArtsStudio POST /v2/{project_id}/design/aggregation-logic-tables
// @API DataArtsStudio GET /v2/{project_id}/design/aggregation-logic-tables/{id}
// @API DataArtsStudio PUT /v2/{project_id}/design/aggregation-logic-tables
// @API DataArtsStudio DELETE /v2/{project_id}/design/aggregation-logic-tables
func ResourceArchitectureAggregationLogicTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureAggregationLogicTableCreate,
		ReadContext:   resourceArchitectureAggregationLogicTableRead,
		UpdateContext: resourceArchitectureAggregationLogicTableUpdate,
		DeleteContext: resourceArchitectureAggregationLogicTableDelete,

		CustomizeDiff: config.FlexibleForceNew(architectureAggregationLogicTableNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureAggregationLogicTableImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the aggregation logic table is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the aggregation logic table belongs.`,
			},
			"tb_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The physical table name in English of the aggregation logic table.`,
			},
			"tb_logic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The display name in Chinese of the aggregation logic table.`,
			},
			"l3_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the business subject to which the aggregation logic table belongs.`,
			},
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
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database corresponding to the data connection.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The asset owner of the aggregation logic table.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the aggregation logic table.`,
			},

			// Optional parameters.
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The alias of the aggregation logic table.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The queue name corresponding to the DLI data connection.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database schema.`,
			},
			"table_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_ch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The Chinese name of the attribute.`,
						},
						"name_en": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The English name of the attribute.`,
						},
						"data_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The data type of the attribute.`,
						},
						"attribute_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The configuration type of the attribute.`,
						},
						"is_primary_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the attribute is the primary key.`,
						},
						"is_partition_key": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the attribute is used as a partition key.`,
						},
						"not_null": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the attribute is not null.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the attribute.`,
						},
						"data_type_extend": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The data type extend field of attribute.`,
						},
						"ref_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the object referenced by the attribute.`,
						},
						"stand_row_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the data standard associated with the attribute.`,
						},
						"alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The alias of the attribute.`,
						},
						"secrecy_levels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
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
								},
							},
							Description: `The list of secrecy levels associated with the attribute.`,
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
					},
				},
				Description: `The list of attributes of the aggregation logic table.`,
			},
			"table_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the database table.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the model to which the aggregation logic table belongs.`,
			},
			"distribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The distribution mode of the database.`,
			},
			"distribute_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The column used for hash distribution.`,
			},
			"compression": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The compression level of the DWS table.`,
			},
			"pre_combine_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The column used for record combine or versioning.`,
			},
			"obs_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OBS storage path of the external table.`,
			},
			"configs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The additional configuration of the aggregation logic table, in JSON format.`,
			},
			"dimension_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the statistical dimension of the derived indicator.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SQL statement of the aggregation logic table.`,
			},
			"partition_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The partition expression of the aggregation logic table.`,
			},
			"dirty_out_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable dirty data output switch.`,
			},
			"dirty_out_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The output database of the dirty data.`,
			},
			"dirty_out_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The prefix of the dirty data table.`,
			},
			"dirty_out_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The suffix of the dirty data table.`,
			},
			"secret_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the data secrecy level.`,
			},
			"self_defined_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of user-defined extended fields for the aggregation logic table.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fd_name_ch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The Chinese display name of the custom extended field.`,
						},
						"fd_name_en": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The English name of the custom extended field.`,
						},
						"not_null": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the custom extended field requires a value.`,
						},
						"fd_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the custom extended field.`,
						},
					},
				},
			},
			"del_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The deletion type when the aggregation logic table is deleted.`,
			},

			// Attributes.
			"dw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data connection.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publishing status of the aggregation logic table.`,
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
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account name of the user who created the aggregation logic table.`,
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
				Description: `The status of the table creation in the production environment after the aggregation logic table is published.`,
			},
			"dev_physical_table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the table creation in the development environment after the aggregation logic table is published.`,
			},
			"technical_asset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronization status of the technical asset after the aggregation logic table is published.`,
			},
			"business_asset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronization status of the business asset after the aggregation logic table is published.`,
			},
			"meta_data_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the asset association after the aggregation logic table is published.`,
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
				Description: `The status of the data quality job creation after the aggregation logic table is published.`,
			},
			"quality_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data quality.`,
			},
			"dlf_task": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the DLF task after the aggregation logic table is published.`,
			},
			"dlf_task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the DLF task after the aggregation logic table is published.`,
			},
			"publish_to_dlm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the DLM API generation after the aggregation logic table is published.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the API after the aggregation logic table is published.`,
			},
			"summary_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronization status of the summary after the aggregation logic table is published.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"table_attributes_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_ch": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The Chinese name of the attribute.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
						"name_en": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The English name of the attribute.`,
								utils.SchemaDescInput{Internal: true},
							),
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'table_attributes'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateArchitectureAggregationLogicTableBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"tb_name":       d.Get("tb_name"),
		"tb_logic_name": d.Get("tb_logic_name"),
		"l3_id":         d.Get("l3_id"),
		"dw_id":         d.Get("dw_id"),
		"dw_type":       d.Get("dw_type"),
		"db_name":       d.Get("db_name"),
		"owner":         d.Get("owner"),
		"description":   d.Get("description"),
		// Optional parameters.
		"alias":               utils.ValueIgnoreEmpty(d.Get("alias")),
		"queue_name":          utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"schema":              utils.ValueIgnoreEmpty(d.Get("schema")),
		"table_attributes":    buildArchitectureAggregationLogicTableAttributes(d.Get("table_attributes").([]interface{})),
		"table_type":          utils.ValueIgnoreEmpty(d.Get("table_type")),
		"model_id":            utils.ValueIgnoreEmpty(d.Get("model_id")),
		"distribute":          utils.ValueIgnoreEmpty(d.Get("distribute")),
		"distribute_column":   utils.ValueIgnoreEmpty(d.Get("distribute_column")),
		"compression":         utils.ValueIgnoreEmpty(d.Get("compression")),
		"pre_combine_field":   utils.ValueIgnoreEmpty(d.Get("pre_combine_field")),
		"obs_location":        utils.ValueIgnoreEmpty(d.Get("obs_location")),
		"configs":             utils.ValueIgnoreEmpty(d.Get("configs")),
		"dimension_group":     utils.ValueIgnoreEmpty(d.Get("dimension_group")),
		"sql":                 utils.ValueIgnoreEmpty(d.Get("sql")),
		"dirty_out_switch":    d.Get("dirty_out_switch"),
		"dirty_out_database":  utils.ValueIgnoreEmpty(d.Get("dirty_out_database")),
		"dirty_out_prefix":    utils.ValueIgnoreEmpty(d.Get("dirty_out_prefix")),
		"dirty_out_suffix":    utils.ValueIgnoreEmpty(d.Get("dirty_out_suffix")),
		"secret_type":         utils.ValueIgnoreEmpty(d.Get("secret_type")),
		"self_defined_fields": buildArchitectureAggregationLogicTableSelfDefinedFields(d.Get("self_defined_fields").([]interface{})),
	}
}

func buildArchitectureAggregationLogicTableAttributes(attributes []interface{}) []map[string]interface{} {
	if len(attributes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(attributes))
	for _, item := range attributes {
		result = append(result, map[string]interface{}{
			"name_ch":          utils.PathSearch("name_ch", item, nil),
			"name_en":          utils.PathSearch("name_en", item, nil),
			"data_type":        utils.PathSearch("data_type", item, nil),
			"attribute_type":   utils.ValueIgnoreEmpty(utils.PathSearch("attribute_type", item, nil)),
			"is_primary_key":   utils.PathSearch("is_primary_key", item, nil),
			"is_partition_key": utils.PathSearch("is_partition_key", item, nil),
			"not_null":         utils.PathSearch("not_null", item, nil),
			"description":      utils.ValueIgnoreEmpty(utils.PathSearch("description", item, nil)),
			"data_type_extend": utils.ValueIgnoreEmpty(utils.PathSearch("data_type_extend", item, nil)),
			"ref_id":           utils.ValueIgnoreEmpty(utils.PathSearch("ref_id", item, nil)),
			"stand_row_id":     utils.ValueIgnoreEmpty(utils.PathSearch("stand_row_id", item, nil)),
			"alias":            utils.ValueIgnoreEmpty(utils.PathSearch("alias", item, nil)),
			"secrecy_levels": buildArchitectureAggregationLogicTableAttributeSecrecyLevels(utils.PathSearch("secrecy_levels",
				item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func buildArchitectureAggregationLogicTableAttributeSecrecyLevels(secrecyLevels []interface{}) []map[string]interface{} {
	if len(secrecyLevels) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(secrecyLevels))
	for _, v := range secrecyLevels {
		result = append(result, map[string]interface{}{
			"id": utils.ValueIgnoreEmpty(utils.PathSearch("id", v, nil)),
		})
	}

	return result
}

func buildArchitectureAggregationLogicTableSelfDefinedFields(selfDefinedFields []interface{}) []map[string]interface{} {
	if len(selfDefinedFields) == 0 || selfDefinedFields[0] == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(selfDefinedFields))
	for _, item := range selfDefinedFields {
		result = append(result, map[string]interface{}{
			"fd_name_ch": utils.ValueIgnoreEmpty(utils.PathSearch("fd_name_ch", item, nil)),
			"fd_name_en": utils.ValueIgnoreEmpty(utils.PathSearch("fd_name_en", item, nil)),
			"not_null":   utils.PathSearch("not_null", item, nil),
			"fd_value":   utils.ValueIgnoreEmpty(utils.PathSearch("fd_value", item, nil)),
		})
	}
	return result
}

func createArchitectureAggregationLogicTable(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	createPath := client.Endpoint + "v2/{project_id}/design/aggregation-logic-tables"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildCreateArchitectureAggregationLogicTableBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return "", err
	}

	reapBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	aggregationLogicTableId := utils.PathSearch("data.value.id", reapBody, "").(string)
	if aggregationLogicTableId == "" {
		return "", errors.New("unable to find aggregation logic table ID from the API response")
	}

	return aggregationLogicTableId, err
}

func resourceArchitectureAggregationLogicTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	aggregationLogicTableId, err := createArchitectureAggregationLogicTable(client, d)
	if err != nil {
		return diag.Errorf("error creating Architecture aggregation logic table: %s", err)
	}

	d.SetId(aggregationLogicTableId)

	if attributes, ok := d.GetOk("table_attributes"); ok && len(attributes.([]interface{})) > 0 {
		// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
		// '_origin' attributes for subsequent determination and construction of the request body during next updates.
		// And whether corresponding parameters are changed, the origin values must be refreshed.
		err = d.Set("table_attributes_origin", refreshTableAttributesOrigin(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(),
			"table_attributes")))
		if err != nil {
			// Don't report an error if origin refresh fails
			log.Printf("[WARN] Unable to refresh the origin values: %s", err)
		}
	}

	if _, ok := d.GetOk("partition_conf"); ok {
		if err = updateArchitectureAggregationLogicTable(client, d); err != nil {
			return diag.Errorf("error updating aggregation logic table (%s): %s", d.Id(), err)
		}
	}

	return resourceArchitectureAggregationLogicTableRead(ctx, d, meta)
}

func refreshTableAttributesOrigin(scriptConfigTableAttributes interface{}) interface{} {
	if scriptConfigTableAttributes == nil {
		return nil
	}

	tableAttributes := scriptConfigTableAttributes.([]interface{})
	originTableAttributes := make([]interface{}, 0, len(tableAttributes))
	for _, item := range tableAttributes {
		originTableAttributes = append(originTableAttributes, map[string]interface{}{
			"name_ch": utils.PathSearch("name_ch", item, nil),
			"name_en": utils.PathSearch("name_en", item, nil),
		})
	}

	return originTableAttributes
}

func GetArchitectureAggregationLogicTableById(client *golangsdk.ServiceClient, workspaceId, aggregationLogicTableId string) (interface{}, error) {
	httpURL := "v2/{project_id}/design/aggregation-logic-tables/{id}?latest=true"
	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", aggregationLogicTableId)
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

	aggregationLogicTable := utils.PathSearch("data.value", respBody, make(map[string]interface{})).(map[string]interface{})
	if len(aggregationLogicTable) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/design/aggregation-logic-tables/{id}?latest=true",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the aggregation logic table (%s) does not exist", aggregationLogicTableId)),
			},
		}
	}

	return aggregationLogicTable, nil
}

func resourceArchitectureAggregationLogicTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceId := d.Get("workspace_id").(string)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	respBody, err := GetArchitectureAggregationLogicTableById(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", architectureAggregationLogicTableErrCodes...),
			"error retrieving aggregation logic table",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		// Required parameters.
		d.Set("workspace_id", workspaceId),
		d.Set("tb_name", utils.PathSearch("tb_name", respBody, nil)),
		d.Set("tb_logic_name", utils.PathSearch("tb_logic_name", respBody, nil)),
		d.Set("dw_id", utils.PathSearch("dw_id", respBody, nil)),
		d.Set("dw_type", utils.PathSearch("dw_type", respBody, nil)),
		d.Set("db_name", utils.PathSearch("db_name", respBody, nil)),
		d.Set("owner", utils.PathSearch("owner", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		// Optional parameters.
		d.Set("alias", utils.PathSearch("alias", respBody, nil)),
		d.Set("queue_name", utils.PathSearch("queue_name", respBody, nil)),
		d.Set("schema", utils.PathSearch("schema", respBody, nil)),
		d.Set("table_attributes", orderArchitectureAggregationLogicTableAttributesByAttributesOrigin(
			utils.PathSearch("table_attributes", respBody, make([]interface{}, 0)).([]interface{}),
			d.Get("table_attributes_origin").([]interface{})),
		),
		d.Set("table_type", utils.PathSearch("table_type", respBody, nil)),
		d.Set("distribute", utils.PathSearch("distribute", respBody, nil)),
		d.Set("distribute_column", utils.PathSearch("distribute_column", respBody, nil)),
		d.Set("compression", utils.PathSearch("compression", respBody, nil)),
		d.Set("pre_combine_field", utils.PathSearch("pre_combine_field", respBody, nil)),
		d.Set("obs_location", utils.PathSearch("obs_location", respBody, nil)),
		d.Set("configs", utils.PathSearch("configs", respBody, nil)),
		d.Set("dimension_group", utils.PathSearch("dimension_group", respBody, nil)),
		d.Set("sql", utils.PathSearch("sql", respBody, nil)),
		d.Set("partition_conf", utils.PathSearch("partition_conf", respBody, nil)),
		d.Set("dirty_out_switch", utils.PathSearch("dirty_out_switch", respBody, false)),
		d.Set("dirty_out_database", utils.PathSearch("dirty_out_database", respBody, nil)),
		d.Set("dirty_out_prefix", utils.PathSearch("dirty_out_prefix", respBody, nil)),
		d.Set("dirty_out_suffix", utils.PathSearch("dirty_out_suffix", respBody, nil)),
		d.Set("self_defined_fields", flattenArchitectureAggregationLogicTableSelfDefinedFields(utils.PathSearch("self_defined_fields",
			respBody, make([]interface{}, 0)).([]interface{}))),
		// Attributes.
		d.Set("dw_name", utils.PathSearch("dw_name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("group_name", utils.PathSearch("group_name", respBody, nil)),
		d.Set("group_code", utils.PathSearch("group_code", respBody, nil)),
		d.Set("model_id", utils.PathSearch("model_id", respBody, nil)),
		d.Set("created_by", utils.PathSearch("create_by", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
		d.Set("env_type", utils.PathSearch("env_type", respBody, nil)),
		d.Set("physical_table", utils.PathSearch("physical_table", respBody, nil)),
		d.Set("dev_physical_table", utils.PathSearch("dev_physical_table", respBody, nil)),
		d.Set("technical_asset", utils.PathSearch("technical_asset", respBody, nil)),
		d.Set("business_asset", utils.PathSearch("business_asset", respBody, nil)),
		d.Set("meta_data_link", utils.PathSearch("meta_data_link", respBody, nil)),
		d.Set("tb_guid", utils.PathSearch("tb_guid", respBody, nil)),
		d.Set("tb_logic_guid", utils.PathSearch("tb_logic_guid", respBody, nil)),
		d.Set("data_quality", utils.PathSearch("data_quality", respBody, nil)),
		d.Set("quality_id", utils.PathSearch("quality_id", respBody, nil)),
		d.Set("dlf_task", utils.PathSearch("dlf_task", respBody, nil)),
		d.Set("dlf_task_id", utils.PathSearch("dlf_task_id", respBody, nil)),
		d.Set("publish_to_dlm", utils.PathSearch("publish_to_dlm", respBody, nil)),
		d.Set("api_id", utils.PathSearch("api_id", respBody, nil)),
		d.Set("summary_status", utils.PathSearch("summary_status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func orderArchitectureAggregationLogicTableAttributesByAttributesOrigin(attributes, attributesOrigin []interface{}) []interface{} {
	if len(attributesOrigin) < 1 {
		return flattenArchitectureAggregationLogicTableAttributes(attributes)
	}

	sortedAttributes := make([]interface{}, 0, len(attributes))
	attributesCopy := attributes
	for _, attributeOrigin := range attributesOrigin {
		nameChOrigin := utils.PathSearch("name_ch", attributeOrigin, "").(string)
		nameEnOrigin := utils.PathSearch("name_en", attributeOrigin, "").(string)
		for index, attr := range attributesCopy {
			if utils.PathSearch("name_ch", attr, "").(string) == nameChOrigin &&
				utils.PathSearch("name_en", attr, "").(string) == nameEnOrigin {
				// Add the found attribute to the sorted attributes list.
				sortedAttributes = append(sortedAttributes, attr)
				// Remove the processed attribute from the original attributes array.
				attributesCopy = append(attributesCopy[:index], attributesCopy[index+1:]...)
				break
			}
		}
	}

	// Add any remaining unsorted attibutes to the end of the sorted list.
	sortedAttributes = append(sortedAttributes, attributesCopy...)
	return flattenArchitectureAggregationLogicTableAttributes(sortedAttributes)
}

func flattenArchitectureAggregationLogicTableAttributes(attributes []interface{}) []interface{} {
	if len(attributes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(attributes))
	for _, v := range attributes {
		resolvedSl, _ := utils.PathSearch("secrecy_levels", v, make([]interface{}, 0)).([]interface{})
		result = append(result, map[string]interface{}{
			"name_ch":          utils.PathSearch("name_ch", v, nil),
			"name_en":          utils.PathSearch("name_en", v, nil),
			"data_type":        utils.PathSearch("data_type", v, nil),
			"attribute_type":   utils.PathSearch("attribute_type", v, nil),
			"is_primary_key":   utils.PathSearch("is_primary_key", v, nil),
			"is_partition_key": utils.PathSearch("is_partition_key", v, nil),
			"not_null":         utils.PathSearch("not_null", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"data_type_extend": utils.PathSearch("data_type_extend", v, nil),
			"ref_id":           utils.PathSearch("ref_id", v, nil),
			"stand_row_id":     utils.PathSearch("stand_row_id", v, nil),
			"alias":            utils.PathSearch("alias", v, nil),
			"secrecy_levels":   flattenArchitectureAggregationLogicTableAttributeSecrecyLevels(resolvedSl),
			// Attributes.
			"id":             utils.PathSearch("id", v, nil),
			"domain_type":    utils.PathSearch("domain_type", v, nil),
			"ref_name_ch":    utils.PathSearch("ref_name_ch", v, nil),
			"ref_name_en":    utils.PathSearch("ref_name_en", v, nil),
			"stand_row_name": utils.PathSearch("stand_row_name", v, nil),
			"ordinal":        utils.PathSearch("ordinal", v, nil),
		})
	}

	return result
}

func flattenArchitectureAggregationLogicTableAttributeSecrecyLevels(secrecyLevels []interface{}) []map[string]interface{} {
	if len(secrecyLevels) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(secrecyLevels))
	for _, v := range secrecyLevels {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("secrecyLevel_id", v, nil),
			// Attributes.
			"uuid": utils.PathSearch("uuid", v, nil),
			"name": utils.PathSearch("secrecyLevel_name", v, nil),
		})
	}

	return result
}

func flattenArchitectureAggregationLogicTableSelfDefinedFields(selfDefinedFields []interface{}) []map[string]interface{} {
	if len(selfDefinedFields) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(selfDefinedFields))
	for _, v := range selfDefinedFields {
		result = append(result, map[string]interface{}{
			"fd_name_ch": utils.PathSearch("fd_name_ch", v, nil),
			"fd_name_en": utils.PathSearch("fd_name_en", v, nil),
			"not_null":   utils.PathSearch("not_null", v, false),
			"fd_value":   utils.PathSearch("fd_value", v, nil),
		})
	}
	return result
}

func buildUpdateArchitectureAggregationLogicTableBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"id":            d.Id(),
		"tb_name":       d.Get("tb_name"),
		"tb_logic_name": d.Get("tb_logic_name"),
		"l3_id":         d.Get("l3_id"),
		"dw_id":         d.Get("dw_id"),
		"dw_type":       d.Get("dw_type"),
		"db_name":       d.Get("db_name"),
		"owner":         d.Get("owner"),
		"description":   d.Get("description"),
		// Optional parameters.
		"alias":               utils.ValueIgnoreEmpty(d.Get("alias")),
		"queue_name":          utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"schema":              utils.ValueIgnoreEmpty(d.Get("schema")),
		"table_attributes":    buildArchitectureAggregationLogicTableAttributes(d.Get("table_attributes").([]interface{})),
		"table_type":          utils.ValueIgnoreEmpty(d.Get("table_type")),
		"model_id":            utils.ValueIgnoreEmpty(d.Get("model_id")),
		"distribute":          utils.ValueIgnoreEmpty(d.Get("distribute")),
		"distribute_column":   utils.ValueIgnoreEmpty(d.Get("distribute_column")),
		"compression":         utils.ValueIgnoreEmpty(d.Get("compression")),
		"pre_combine_field":   utils.ValueIgnoreEmpty(d.Get("pre_combine_field")),
		"obs_location":        d.Get("obs_location"),
		"configs":             utils.ValueIgnoreEmpty(d.Get("configs")),
		"dimension_group":     utils.ValueIgnoreEmpty(d.Get("dimension_group")),
		"sql":                 utils.ValueIgnoreEmpty(d.Get("sql")),
		"partition_conf":      d.Get("partition_conf"),
		"dirty_out_switch":    d.Get("dirty_out_switch"),
		"dirty_out_database":  utils.ValueIgnoreEmpty(d.Get("dirty_out_database")),
		"dirty_out_prefix":    utils.ValueIgnoreEmpty(d.Get("dirty_out_prefix")),
		"dirty_out_suffix":    utils.ValueIgnoreEmpty(d.Get("dirty_out_suffix")),
		"secret_type":         utils.ValueIgnoreEmpty(d.Get("secret_type")),
		"self_defined_fields": buildArchitectureAggregationLogicTableSelfDefinedFields(d.Get("self_defined_fields").([]interface{})),
	}
}

func updateArchitectureAggregationLogicTable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpURL := "v2/{project_id}/design/aggregation-logic-tables"
	updatePath := client.Endpoint + httpURL
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildUpdateArchitectureAggregationLogicTableBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceArchitectureAggregationLogicTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	if err := updateArchitectureAggregationLogicTable(client, d); err != nil {
		return diag.Errorf("error updating aggregation logic table (%s): %s", d.Id(), err)
	}

	if d.HasChange("table_attributes") {
		// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
		// '_origin' attributes for subsequent determination and construction of the request body during next updates.
		// And whether corresponding parameters are changed, the origin values must be refreshed.
		err = d.Set("table_attributes_origin", refreshTableAttributesOrigin(utils.GetNestedObjectFromRawConfig(d.GetRawConfig(),
			"table_attributes")))
		if err != nil {
			// Don't report an error if origin refresh fails
			log.Printf("[WARN] Unable to refresh the origin values: %s", err)
		}
	}

	return resourceArchitectureAggregationLogicTableRead(ctx, d, meta)
}

func deleteArchitectureAggregationLogicTable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpURL := "v2/{project_id}/design/aggregation-logic-tables"
	deletePath := client.Endpoint + httpURL
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"ids":      []string{d.Id()},
			"del_type": utils.ValueIgnoreEmpty(d.Get("del_type")),
		}),
	}
	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceArchitectureAggregationLogicTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	if err := deleteArchitectureAggregationLogicTable(client, d); err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", architectureAggregationLogicTableErrCodes[0]),
			fmt.Sprintf("error deleting aggregation logic table (%s)", d.Id()),
		)
	}

	return nil
}

func resourceArchitectureAggregationLogicTableImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
