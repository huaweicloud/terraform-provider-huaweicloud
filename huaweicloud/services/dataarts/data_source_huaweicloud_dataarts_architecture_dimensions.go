package dataarts

import (
	"context"
	"fmt"
	"net/url"
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

// @API DataArtsStudio GET /v2/{project_id}/design/dimensions
func DataSourceArchitectureDimensions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureDimensionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dimensions are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the dimensions belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name or code of the dimension to be fuzzy queried.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the dimension to be exactly queried.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The English name of the dimension to be exactly queried.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the dimension to be queried.`,
			},
			"approver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The approver of the dimension to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The publishing status of the dimension to be queried.`,
			},
			"l2_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The subject domain L2 ID to which the dimension belongs.`,
			},
			"dimension_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the dimension to be queried.`,
			},
			"biz_catalog_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business catalog ID to which the dimension belongs.`,
			},
			"fact_logic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The fact table ID of the dimension to be queried.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the dimension to be queried, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the dimension to be queried, in RFC3339 format.`,
			},
			"derivative_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The derivative indicator ID list for querying dimensions.`,
			},

			// Attributes.
			"dimensions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionsElem(),
				Description: `The list of dimensions that matched filter parameters.`,
			},
		},
	}
}

func dataArchitectureDimensionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the dimension.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the dimension.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the dimension.`,
			},
			"dimension_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the dimension.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the dimension.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset owner of the dimension.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publishing status of the dimension.`,
			},
			"code_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code table ID of the dimension.`,
			},
			"l1_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L1 subject domain grouping ID of the dimension.`,
			},
			"l2_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L2 subject domain ID of the dimension.`,
			},
			"l3_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L3 business object ID of the dimension.`,
			},
			"l1_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L1 subject domain grouping name of the dimension.`,
			},
			"l2_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L2 subject domain name of the dimension.`,
			},
			"l3_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The L3 business object name of the dimension.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the dimension.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the dimension.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the dimension, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the dimension, in RFC3339 format.`,
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
			"obs_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS location of the dimension.`,
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
			"datasource": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionDatasource(),
				Description: `The data source configuration of the dimension.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionAttributes(),
				Description: `The list of attributes of the dimension.`,
			},
			"hierarchies": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionHierarchies(),
				Description: `The hierarchy attributes of the dimension.`,
			},
			"code_table": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionCodeTable(),
				Description: `The code table information of the dimension.`,
			},
			"model": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionModel(),
				Description: `The model information of the dimension.`,
			},
		},
	}
}

func dataArchitectureDimensionDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data source.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business type of the data source.`,
			},
			"biz_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business ID of the data source.`,
			},
			"dw_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data connection.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data connection.`,
			},
			"dw_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data connection.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database corresponding to the data connection.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The queue name corresponding to the DLI data connection.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database schema.`,
			},
		},
	}
}

func dataArchitectureDimensionAttributes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the attribute.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the attribute.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the attribute.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the attribute.`,
			},
			"data_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type of the attribute.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain type of the attribute.`,
			},
			"data_type_extend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data type extension of the attribute.`,
			},
			"is_primary_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is the primary key.`,
			},
			"is_biz_primary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is the business primary key.`,
			},
			"is_partition_key": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is the partition key.`,
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The sequence number of the attribute.`,
			},
			"not_null": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the attribute is not null.`,
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
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The quality information of the attribute.`,
			},
			"secrecy_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The secrecy levels of the attribute.`,
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
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The self-defined fields of the attribute.`,
			},
		},
	}
}

func dataArchitectureDimensionHierarchies() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the hierarchy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the hierarchy.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the hierarchy.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the hierarchy.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the hierarchy, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the hierarchy, in RFC3339 format.`,
			},
			"attrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDimensionHierarchiesAttr(),
				Description: `The attributes of the hierarchy.`,
			},
		},
	}
}

func dataArchitectureDimensionHierarchiesAttr() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the hierarchy attribute.`,
			},
			"hierarchies_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hierarchy ID of the attribute.`,
			},
			"attr_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attribute ID.`,
			},
			"level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The level of the hierarchy attribute.`,
			},
			"attr_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the referenced attribute.`,
			},
			"attr_name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the referenced attribute.`,
			},
		},
	}
}

func dataArchitectureDimensionCodeTable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the code table.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The English name of the code table.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Chinese name of the code table.`,
			},
			"tb_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The version of the code table.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory ID of the code table.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory path of the code table.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the code table.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publishing status of the code table.`,
			},
		},
	}
}

func dataArchitectureDimensionModel() *schema.Resource {
	return &schema.Resource{
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
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The database names.`,
			},
			"table_model_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The table model prefix.`,
			},
		},
	}
}

func buildArchitectureDimensionsQueryParams(d *schema.ResourceData) string {
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
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("l2_id"); ok {
		res = fmt.Sprintf("%s&l2_id=%v", res, v)
	}
	if v, ok := d.GetOk("dimension_type"); ok {
		res = fmt.Sprintf("%s&dimension_type=%v", res, v)
	}
	if v, ok := d.GetOk("biz_catalog_id"); ok {
		res = fmt.Sprintf("%s&biz_catalog_id=%v", res, v)
	}
	if v, ok := d.GetOk("fact_logic_id"); ok {
		res = fmt.Sprintf("%s&fact_logic_id=%v", res, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%s", res, url.QueryEscape(v.(string)))
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%s", res, url.QueryEscape(v.(string)))
	}
	if v, ok := d.GetOk("derivative_ids"); ok {
		for _, id := range v.([]interface{}) {
			res = fmt.Sprintf("%s&derivative_ids=%v", res, id)
		}
	}

	return res
}

func listArchitectureDimensions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/dimensions?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureDimensionsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, records...)

		if len(records) < limit {
			break
		}
		offset += len(records)
	}

	return result, nil
}

func flattenDimensionsHierarchies(hierarchies []interface{}) []map[string]interface{} {
	if len(hierarchies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(hierarchies))
	for _, hierarchy := range hierarchies {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", hierarchy, nil),
			"name":       utils.PathSearch("name", hierarchy, nil),
			"created_by": utils.PathSearch("create_by", hierarchy, nil),
			"updated_by": utils.PathSearch("update_by", hierarchy, nil),
			"created_at": utils.PathSearch("create_time", hierarchy, nil),
			"updated_at": utils.PathSearch("update_time", hierarchy, nil),
			"attrs":      flattenDimensionsHierarchiesAttrs(utils.PathSearch("attrs", hierarchy, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenDimensionsHierarchiesAttrs(attrs []interface{}) []map[string]interface{} {
	if len(attrs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(attrs))
	for _, attr := range attrs {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", attr, nil),
			"hierarchies_id": utils.PathSearch("hierarchies_id", attr, nil),
			"attr_id":        utils.PathSearch("attr_id", attr, nil),
			"level":          utils.PathSearch("level", attr, nil),
			"attr_name_en":   utils.PathSearch("attr_name_en", attr, nil),
			"attr_name_ch":   utils.PathSearch("attr_name_ch", attr, nil),
		})
	}
	return result
}

func flattenDimensionsCodeTable(codeTable map[string]interface{}) []map[string]interface{} {
	if len(codeTable) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":             utils.PathSearch("id", codeTable, nil),
			"name_en":        utils.PathSearch("name_en", codeTable, nil),
			"name_ch":        utils.PathSearch("name_ch", codeTable, nil),
			"tb_version":     utils.PathSearch("tb_version", codeTable, nil),
			"directory_id":   utils.PathSearch("directory_id", codeTable, nil),
			"directory_path": utils.PathSearch("directory_path", codeTable, nil),
			"description":    utils.PathSearch("description", codeTable, nil),
			"status":         utils.PathSearch("status", codeTable, nil),
		},
	}
}

func flattenDimensionsDatasource(datasource map[string]interface{}) []map[string]interface{} {
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

func flattenDimensionsAttributes(attributes []interface{}) []map[string]interface{} {
	if len(attributes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(attributes))
	for _, attr := range attributes {
		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("id", attr, nil),
			"name_en":             utils.PathSearch("name_en", attr, nil),
			"name_ch":             utils.PathSearch("name_ch", attr, nil),
			"description":         utils.PathSearch("description", attr, nil),
			"data_type":           utils.PathSearch("data_type", attr, nil),
			"domain_type":         utils.PathSearch("domain_type", attr, nil),
			"data_type_extend":    utils.PathSearch("data_type_extend", attr, nil),
			"is_primary_key":      utils.PathSearch("is_primary_key", attr, nil),
			"is_biz_primary":      utils.PathSearch("is_biz_primary", attr, nil),
			"is_partition_key":    utils.PathSearch("is_partition_key", attr, nil),
			"ordinal":             utils.PathSearch("ordinal", attr, nil),
			"not_null":            utils.PathSearch("not_null", attr, nil),
			"code_table_field_id": utils.PathSearch("code_table_field_id", attr, nil),
			"create_by":           utils.PathSearch("create_by", attr, nil),
			"stand_row_id":        utils.PathSearch("stand_row_id", attr, nil),
			"stand_row_name":      utils.PathSearch("stand_row_name", attr, nil),
			"quality_infos":       utils.PathSearch("quality_infos", attr, nil),
			"secrecy_levels":      utils.PathSearch("secrecy_levels", attr, nil),
			"status":              utils.PathSearch("status", attr, nil),
			"create_time":         utils.PathSearch("create_time", attr, nil),
			"update_time":         utils.PathSearch("update_time", attr, nil),
			"alias":               utils.PathSearch("alias", attr, nil),
			"self_defined_fields": utils.PathSearch("self_defined_fields", attr, nil),
		})
	}
	return result
}

func flattenDimensionsModel(model map[string]interface{}) []map[string]interface{} {
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

func flattenArchitectureDimensions(dimensions []interface{}) []map[string]interface{} {
	if len(dimensions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", dimension, nil),
			"name_ch":           utils.PathSearch("name_ch", dimension, nil),
			"name_en":           utils.PathSearch("name_en", dimension, nil),
			"dimension_type":    utils.PathSearch("dimension_type", dimension, nil),
			"description":       utils.PathSearch("description", dimension, nil),
			"owner":             utils.PathSearch("owner", dimension, nil),
			"status":            utils.PathSearch("status", dimension, nil),
			"code_table_id":     utils.PathSearch("code_table_id", dimension, nil),
			"l1_id":             utils.PathSearch("l1_id", dimension, nil),
			"l2_id":             utils.PathSearch("l2_id", dimension, nil),
			"l3_id":             utils.PathSearch("l3_id", dimension, nil),
			"l1_name":           utils.PathSearch("l1", dimension, nil),
			"l2_name":           utils.PathSearch("l2", dimension, nil),
			"l3_name":           utils.PathSearch("l3", dimension, nil),
			"created_by":        utils.PathSearch("create_by", dimension, nil),
			"updated_by":        utils.PathSearch("update_by", dimension, nil),
			"created_at":        utils.PathSearch("create_time", dimension, nil),
			"updated_at":        utils.PathSearch("update_time", dimension, nil),
			"table_type":        utils.PathSearch("table_type", dimension, nil),
			"distribute":        utils.PathSearch("distribute", dimension, nil),
			"distribute_column": utils.PathSearch("distribute_column", dimension, nil),
			"obs_location":      utils.PathSearch("obs_location", dimension, nil),
			"alias":             utils.PathSearch("alias", dimension, nil),
			"configs":           utils.PathSearch("configs", dimension, nil),
			"env_type":          utils.PathSearch("env_type", dimension, nil),
			"model_id":          utils.PathSearch("model_id", dimension, nil),
			"dev_version":       utils.PathSearch("dev_version", dimension, nil),
			"prod_version":      utils.PathSearch("prod_version", dimension, nil),
			"dev_version_name":  utils.PathSearch("dev_version_name", dimension, nil),
			"prod_version_name": utils.PathSearch("prod_version_name", dimension, nil),
			"datasource": flattenDimensionsDatasource(
				utils.PathSearch("datasource", dimension, make(map[string]interface{})).(map[string]interface{})),
			"attributes":  flattenDimensionsAttributes(utils.PathSearch("attributes", dimension, make([]interface{}, 0)).([]interface{})),
			"hierarchies": flattenDimensionsHierarchies(utils.PathSearch("hierarchies", dimension, make([]interface{}, 0)).([]interface{})),
			"code_table": flattenDimensionsCodeTable(
				utils.PathSearch("code_table", dimension, make(map[string]interface{})).(map[string]interface{})),
			"model": flattenDimensionsModel(
				utils.PathSearch("model", dimension, make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return result
}

func dataSourceArchitectureDimensionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	dimensions, err := listArchitectureDimensions(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture dimensions: %s", err)
	}

	d.SetId(uuid.New().String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("dimensions", flattenArchitectureDimensions(dimensions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
