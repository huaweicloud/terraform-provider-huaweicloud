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

var (
	apiResourceNotFoundCodes = []string{
		"DLM.4001", // Instance or workspace does not exist.
		"DLM.4018", // API does not exist (during API detail query).
	}
	apiResourceNotFoundCodesForDelete = []string{
		"DLM.4001", // Instance or workspace does not exist.
		"DLM.4197", // API does not exist (during delete API operation).
	}
)

// @API DataArtsStudio POST /v1/{project_id}/service/apis
// @API DataArtsStudio GET /v2/{project_id}/{connection_id}/datatables
// @API DataArtsStudio GET /v1/{project_id}/data-connections/{data_connection_id}
// @API DataArtsStudio GET /v1/{project_id}/service/apis/{api_id}
// @API DataArtsStudio PUT /v1/{project_id}/service/apis/{api_id}
// @API DataArtsStudio POST /v1/{project_id}/service/apis/batch-delete
func ResourceDataServiceApi() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiCreate,
		ReadContext:   resourceDataServiceApiRead,
		UpdateContext: resourceDataServiceApiUpdate,
		DeleteContext: resourceDataServiceApiDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataServiceApiImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the API is located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the workspace to which the API belongs.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The type of DLM engine.`,
			},

			// Arguments
			"catalog_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the catalog where the API is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the API.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The API type.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The authentication type.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request protocol of the API.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The API request path.`,
			},
			"request_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request type of the API.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The API reviewer.`,
			},
			"datasource_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        apiDataSourceConfigElemSchema(),
				Description: `The configuration of the API data source.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the API.`,
			},
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The visibility to the catalog of API.`,
			},
			"request_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        apiRequestParamsElemSchema(),
				Description: `The parameters of the API request.`,
			},
			"backend_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        apiBackendConfigElemSchema(),
				Description: `The configuration of the API backend.`,
			},

			// Public attributes
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the API, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the API, in RFC3339 format.`,
			},
			// Attribute for shared API.
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the group to which the API belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API status.`,
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API host configuration, for shared type.`,
			},
			// Attribute for exclusvie API.
			"hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        apiHostsSchema(),
				Description: `The API host configuration, for exclusive type.`,
			},
		},
	}
}

func apiDataSourceConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the data source.`,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The ID of the data connection.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database.`,
			},
			"datatable": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the data table.`,
			},
			"queue": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the DLI queue.`,
			},
			"access_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The access mode for the data.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SQL statements in script access type.`,
			},
			"backend_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        apiDataSourceConfigBackendParamsElemSchema(),
				Description: `The backend parameters of the API.`,
			},
			"response_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        apiDataSourceConfigResponseParamsElemSchema(),
				Description: `The response parameters of the API.`,
			},
			"order_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        apiDataSourceConfigOrderParamsElemSchema(),
				Description: `The order parameters of the API.`,
			},
			// These following parameters can be obtained through API request (for DataArts Studio service) calls.
			// However, considering that the API request calls may fail in some regions, corresponding parameters are
			// provided to ensure the creation of resources in this special case.
			"connection_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The name of the data connection.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
			"table_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The ID of the data table.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
	return &sc
}

func apiDataSourceConfigBackendParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the backend parameter.`,
			},
			"mapping": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the mapping parameter.`,
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The condition character.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
		},
	}
	return &sc
}

func apiDataSourceConfigResponseParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the response parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the response parameter.`,
			},
			"field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The bound table field for the response parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the response parameter.`,
			},
			"example_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The example value of the response parameter.`,
			},
		},
	}
	return &sc
}

func apiDataSourceConfigOrderParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the order parameter.`,
			},
			"field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The corresponding parameter field for the order parameter.`,
			},
			"optional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether this order parameter is the optional parameter.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The sort type of the order parameter.`,
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The order of the sorting parameters.`,
			},
		},
	}
	return &sc
}

func apiRequestParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the request parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The position of the request parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the request parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the request parameter.`,
			},
			"necessary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether this parameter is the required parameter.`,
			},
			"example_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The example value of the request parameter.`,
			},
			"default_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The default value of the request parameter.`,
			},
		},
	}
	return &sc
}

func apiBackendConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the backend request.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The protocol of the backend request.`,
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The backend host.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The backend path.`,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The backend timeout.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"backend_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        apiBackendConfigBackendParamsElemSchema(),
				Description: `The backend parameters of the API.`,
			},
			"constant_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        apiBackendConfigConstantParamsElemSchema(),
				Description: `The backend constant parameters of the API.`,
			},
		},
	}
	return &sc
}

func apiBackendConfigBackendParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the request parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The position of the request parameter.`,
			},
			"backend_param_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the corresponding backend parameter.`,
			},
		},
	}
	return &sc
}

func apiBackendConfigConstantParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the constant parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the constant parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The position of the constant parameter.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of the constant parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the constant parameter.`,
			},
		},
	}
	return &sc
}

func apiHostsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster ID to which the API belongs.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster name to which the API belongs.`,
			},
			"intranet_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The intranet address.`,
			},
			"external_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The exrernal address.`,
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of gateway damains.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func buildModifyDataServiceApiBodyParams(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// Required
		"catalog_id":        d.Get("catalog_id"),
		"name":              d.Get("name"),
		"api_type":          d.Get("type"),
		"auth_type":         d.Get("auth_type"),
		"protocol":          d.Get("protocol"),
		"path":              d.Get("path"),
		"request_type":      d.Get("request_type"),
		"manager":           d.Get("manager"),
		"datasource_config": buildApiDataSourceConfigBodyParams(client, workspaceId, d.Get("datasource_config").([]interface{})),
		// Optional
		"description":    d.Get("description"),
		"visibility":     d.Get("visibility"),
		"request_paras":  buildApiRequestParamsBodyParams(d.Get("request_params").(*schema.Set)),
		"backend_config": utils.RemoveNil(buildApiBackendConfigBodyParams(d.Get("backend_config").([]interface{}))),
	}
	return bodyParams
}

func buildApiRequestParamsBodyParams(params *schema.Set) []interface{} {
	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", val, nil),
			"position":      utils.PathSearch("position", val, nil),
			"type":          utils.PathSearch("type", val, nil),
			"description":   utils.PathSearch("description", val, nil),
			"necessary":     utils.PathSearch("necessary", val, false),
			"example_value": utils.PathSearch("example_value", val, nil),
			"default_value": utils.PathSearch("default_value", val, nil),
		})
	}
	return result
}

func getDataConnectionById(client *golangsdk.ServiceClient, workspaceId, connectionId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/data-connections/{data_connection_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{data_connection_id}", connectionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error querying DataArts Studio connection by its ID (%s): %s", connectionId, err)
	}
	return utils.FlattenResponse(requestResp)
}

func getDataSourceDatatableById(client *golangsdk.ServiceClient, workspaceId, connectionId, databaseName, tableName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/{connection_id}/datatables?database_name={database_name}&table_name={table_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", connectionId)
	getPath = strings.ReplaceAll(getPath, "{database_name}", databaseName)
	getPath = strings.ReplaceAll(getPath, "{table_name}", tableName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error querying data table by its name (%s) under the database (%s): %s", tableName, databaseName, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("tables[0]", respBody, nil), nil
}

func buildApiDataSourceConfigBodyParams(client *golangsdk.ServiceClient, workspaceId string, datasourceConfig []interface{}) map[string]interface{} {
	if len(datasourceConfig) < 1 {
		return nil
	}

	cfgDetail := datasourceConfig[0]
	connectionId := utils.PathSearch("connection_id", cfgDetail, "").(string)
	databaseName := utils.PathSearch("database", cfgDetail, "").(string)
	tableName := utils.PathSearch("datatable", cfgDetail, "").(string)
	result := map[string]interface{}{
		"type":          utils.PathSearch("type", cfgDetail, nil),
		"connection_id": connectionId,
		"database":      databaseName,
		"datatable":     tableName,
		"queue":         utils.ValueIgnoreEmpty(utils.PathSearch("queue", cfgDetail, nil)),
		"access_mode":   utils.ValueIgnoreEmpty(utils.PathSearch("access_mode", cfgDetail, nil)),
		"pagination":    "DEFAULT",
		"sql":           utils.ValueIgnoreEmpty(utils.PathSearch("sql", cfgDetail, nil)),
		"backend_paras": buildApiDataSourceConfigBackendParamsBodyParams(utils.PathSearch("backend_params",
			cfgDetail, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		"response_paras": buildApiDataSourceConfigResponseParamsBodyParams(utils.PathSearch("response_params",
			cfgDetail, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		"order_paras": buildApiDataSourceConfigOrderParamsBodyParams(utils.PathSearch("order_params",
			cfgDetail, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
	}

	// The parameter 'datasource_config.0.connection_name' is required.
	// It's not friendly to require user to enter the connection name after having already entered the connection ID.
	if connectionName := utils.PathSearch("connection_name", cfgDetail, "").(string); connectionName != "" {
		result["connection_name"] = connectionName
	} else {
		connectionDetail, err := getDataConnectionById(client, workspaceId, connectionId)
		if err != nil {
			log.Printf("[ERROR] unable to build data source configuration: %s", err)
		} else {
			result["connection_name"] = utils.PathSearch("dw_name", connectionDetail, "")
		}
	}

	// The value format of parameter 'datasource_config.0.table_id' is 'Native-{connection_id}-{database}-{datatable}',
	// but the format maybe change in the future and this parameter is required.
	// Query is the safest way.
	if tableId := utils.PathSearch("table_id", cfgDetail, "").(string); tableId != "" {
		result["table_id"] = tableId
	} else {
		tableDetail, err := getDataSourceDatatableById(client, workspaceId, connectionId, databaseName, tableName)
		if err != nil {
			log.Printf("[ERROR] unable to build data source configuration: %s", err)
		} else {
			result["table_id"] = utils.PathSearch("table_id", tableDetail, "")
		}
	}

	return result
}

func buildApiDataSourceConfigBackendParamsBodyParams(params *schema.Set) []interface{} {
	if params.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		result = append(result, map[string]interface{}{
			"name":      utils.PathSearch("name", val, nil),
			"mapping":   utils.PathSearch("mapping", val, nil),
			"condition": utils.PathSearch("condition", val, nil),
		})
	}
	return result
}

func buildApiDataSourceConfigResponseParamsBodyParams(params *schema.Set) []interface{} {
	if params.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", val, nil),
			"type":          utils.PathSearch("type", val, nil),
			"field":         utils.PathSearch("field", val, nil),
			"description":   utils.PathSearch("description", val, nil),
			"example_value": utils.PathSearch("example_value", val, nil),
		})
	}
	return result
}

func buildApiDataSourceConfigOrderParamsBodyParams(params *schema.Set) []interface{} {
	if params.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		// The default integer triggers the structure change.
		result = append(result, map[string]interface{}{
			"name":     utils.StringIgnoreEmpty(utils.PathSearch("name", val, "").(string)),
			"field":    utils.StringIgnoreEmpty(utils.PathSearch("field", val, "").(string)),
			"optional": utils.PathSearch("optional", val, nil),
			"sort":     utils.StringIgnoreEmpty(utils.PathSearch("sort", val, "").(string)),
			"order":    utils.IntIgnoreEmpty(utils.PathSearch("order", val, 0).(int)),
		})
	}
	return result
}

func buildApiBackendConfigBodyParams(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	// The default integer triggers the structure change.
	return map[string]interface{}{
		"type":     utils.ValueIgnoreEmpty(utils.PathSearch("type", params[0], "").(string)),
		"protocol": utils.ValueIgnoreEmpty(utils.PathSearch("protocol", params[0], "").(string)),
		"host":     utils.ValueIgnoreEmpty(utils.PathSearch("host", params[0], "").(string)),
		"path":     utils.ValueIgnoreEmpty(utils.PathSearch("path", params[0], "").(string)),
		"timeout":  utils.ValueIgnoreEmpty(utils.PathSearch("timeout", params[0], 0).(int)),
		"backend_paras": buildApiBackendConfigBackendParamsBodyParams(utils.PathSearch("backend_params",
			params[0], schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		"constant_paras": buildApiBackendConfigConstantParamsBodyParams(utils.PathSearch("constant_params",
			params[0], schema.NewSet(schema.HashString, nil)).(*schema.Set)),
	}
}

func buildApiBackendConfigBackendParamsBodyParams(params *schema.Set) []interface{} {
	if params.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", val, nil),
			"position":          utils.PathSearch("position", val, nil),
			"backend_para_name": utils.PathSearch("backend_param_name", val, nil),
		})
	}
	return result
}

func buildApiBackendConfigConstantParamsBodyParams(params *schema.Set) []interface{} {
	if params.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, params.Len())
	for _, val := range params.List() {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", val, nil),
			"type":        utils.PathSearch("type", val, nil),
			"position":    utils.PathSearch("position", val, nil),
			"description": utils.PathSearch("description", val, nil),
			"value":       utils.PathSearch("value", val, nil),
		})
	}
	return result
}

func resourceDataServiceApiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/service/apis"
		workspaceId = d.Get("workspace_id").(string)
		dlmType     = d.Get("dlm_type").(string)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     dlmType,
		},
		JSONBody: buildModifyDataServiceApiBodyParams(client, d, workspaceId),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating API: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	apiId := utils.PathSearch("id", respBody, "").(string)
	if apiId == "" {
		return diag.Errorf("unable to find the DataArts DataService API ID from the API response")
	}
	d.SetId(apiId)

	return resourceDataServiceApiRead(ctx, d, meta)
}

func GetDataServiceApi(client *golangsdk.ServiceClient, workspaceId, dlmType, apiId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/service/apis/{api_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{api_id}", apiId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     dlmType,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", apiResourceNotFoundCodes...)
	}
	return utils.FlattenResponse(requestResp)
}

func flattenExclusiveApiHosts(hosts []interface{}) []interface{} {
	if len(hosts) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(hosts))
	for _, host := range hosts {
		result = append(result, map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", host, nil),
			"instance_name": utils.PathSearch("instance_name", host, nil),
			"intranet_host": utils.PathSearch("intranet_host", host, nil),
			"external_host": utils.PathSearch("external_host", host, nil),
			"domains":       utils.PathSearch("domains", host, nil),
		})
	}
	return result
}

func flattenApiRequestParams(requestParams []interface{}) []interface{} {
	if len(requestParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(requestParams))
	for _, requestParam := range requestParams {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", requestParam, nil),
			"position":      utils.PathSearch("position", requestParam, nil),
			"type":          utils.PathSearch("type", requestParam, nil),
			"description":   utils.PathSearch("description", requestParam, nil),
			"necessary":     utils.PathSearch("necessary", requestParam, nil),
			"example_value": utils.PathSearch("example_value", requestParam, nil),
			"default_value": utils.PathSearch("default_value", requestParam, nil),
		})
	}
	return result
}

func getDataServiceDatatableById(client *golangsdk.ServiceClient, workspaceId, connectionId, databaseName, tableId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/{connection_id}/datatables?database_name={database_name}&limit=16"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", connectionId)
	getPath = strings.ReplaceAll(getPath, "{database_name}", databaseName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}
	var offset float64
	for {
		getPathWithOffset := fmt.Sprintf("%s&offset=%v", getPath, offset)
		requestResp, err := client.Request("GET", getPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error querying datatables (%s) under the database (%s): %s", tableId, databaseName, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		totalTables := utils.PathSearch("length(tables)", respBody, float64(0)).(float64)
		if totalTables < 1 {
			break
		}
		if tableDetail := utils.PathSearch(fmt.Sprintf("tables[?table_id=='%s']|[0]", tableId), respBody, nil); tableDetail != nil {
			return tableDetail, nil
		}
		offset += totalTables
	}
	return nil, fmt.Errorf("unable to find the datatable by its ID (%s)", tableId)
}

func flattenApiDataSourceConfig(client *golangsdk.ServiceClient, workspaceId string, dataSourceConfig interface{}) []map[string]interface{} {
	if dataSourceConfig == nil {
		return nil
	}

	connectionId := utils.PathSearch("connection_id", dataSourceConfig, "").(string)
	databaseName := utils.PathSearch("database", dataSourceConfig, "").(string)
	tableId := utils.PathSearch("table_id", dataSourceConfig, "").(string)
	result := map[string]interface{}{
		"type":            utils.PathSearch("type", dataSourceConfig, nil),
		"connection_name": utils.PathSearch("connection_name", dataSourceConfig, nil),
		"connection_id":   connectionId,
		"database":        databaseName,
		"table_id":        tableId,
		"queue":           utils.PathSearch("queue", dataSourceConfig, nil),
		"access_mode":     utils.PathSearch("access_mode", dataSourceConfig, nil),
		"sql":             utils.PathSearch("sql", dataSourceConfig, nil),
		"backend_params": flattenApiDataSourceConfigBackendParams(utils.PathSearch("backend_paras",
			dataSourceConfig, make([]interface{}, 0)).([]interface{})),
		"response_params": flattenApiDataSourceConfigResponseParams(utils.PathSearch("response_paras",
			dataSourceConfig, make([]interface{}, 0)).([]interface{})),
		"order_params": flattenApiDataSourceConfigOrderParams(utils.PathSearch("order_paras",
			dataSourceConfig, make([]interface{}, 0)).([]interface{})),
	}
	if tableName := utils.PathSearch("datatable", dataSourceConfig, "").(string); tableName != "" {
		result["datatable"] = tableName
	} else {
		tableDetail, err := getDataServiceDatatableById(client, workspaceId, connectionId, databaseName, tableId)
		if err != nil {
			log.Printf("[ERROR] error setting datatable name: %s", err)
		}
		result["datatable"] = utils.PathSearch("table_name", tableDetail, nil)
	}

	return []map[string]interface{}{result}
}

func flattenApiDataSourceConfigBackendParams(backendParams []interface{}) []interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(backendParams))
	for _, backendParam := range backendParams {
		result = append(result, map[string]interface{}{
			"name":      utils.PathSearch("name", backendParam, nil),
			"mapping":   utils.PathSearch("mapping", backendParam, nil),
			"condition": utils.PathSearch("condition", backendParam, nil),
		})
	}
	return result
}

func flattenApiDataSourceConfigResponseParams(responseParams []interface{}) []interface{} {
	if len(responseParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(responseParams))
	for _, responseParam := range responseParams {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", responseParam, nil),
			"type":          utils.PathSearch("type", responseParam, nil),
			"field":         utils.PathSearch("field", responseParam, nil),
			"description":   utils.PathSearch("description", responseParam, nil),
			"example_value": utils.PathSearch("example_value", responseParam, nil),
		})
	}
	return result
}

func flattenApiDataSourceConfigOrderParams(orderParams []interface{}) []interface{} {
	if len(orderParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(orderParams))
	for _, orderParam := range orderParams {
		result = append(result, map[string]interface{}{
			"name":     utils.PathSearch("name", orderParam, nil),
			"field":    utils.PathSearch("field", orderParam, nil),
			"optional": utils.PathSearch("optional", orderParam, nil),
			"sort":     utils.PathSearch("sort", orderParam, nil),
			"order":    utils.PathSearch("order", orderParam, nil),
		})
	}
	return result
}

func flattenApiBackendConfig(backendConfig interface{}) []map[string]interface{} {
	if backendConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":     utils.PathSearch("type", backendConfig, nil),
			"protocol": utils.PathSearch("protocol", backendConfig, nil),
			"host":     utils.PathSearch("host", backendConfig, nil),
			"timeout":  utils.PathSearch("timeout", backendConfig, nil),
			"path":     utils.PathSearch("path", backendConfig, nil),
			"backend_params": flattenApiBackendConfigBackendParams(utils.PathSearch("backend_paras", backendConfig,
				make([]interface{}, 0)).([]interface{})),
			"constant_params": flattenApiBackendConfigConstantParams(utils.PathSearch("constant_paras", backendConfig,
				make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenApiBackendConfigBackendParams(backendParams []interface{}) []interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(backendParams))
	for _, backendParam := range backendParams {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", backendParam, nil),
			"position":           utils.PathSearch("position", backendParam, nil),
			"backend_param_name": utils.PathSearch("backend_para_name", backendParam, nil),
		})
	}
	return result
}

func flattenApiBackendConfigConstantParams(backendParams []interface{}) []interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(backendParams))
	for _, backendParam := range backendParams {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", backendParam, nil),
			"type":        utils.PathSearch("type", backendParam, nil),
			"position":    utils.PathSearch("position", backendParam, nil),
			"description": utils.PathSearch("description", backendParam, nil),
			"value":       utils.PathSearch("value", backendParam, nil),
		})
	}
	return result
}

func resourceDataServiceApiRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		workspaceId = d.Get("workspace_id").(string)
		dlmType     = d.Get("dlm_type").(string)
		apiId       = d.Id()
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	apiRecord, err := GetDataServiceApi(client, workspaceId, dlmType, apiId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error getting Data Service API (%s) for DataArts Studio service", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Public arguments.
		d.Set("name", utils.PathSearch("name", apiRecord, nil)),
		d.Set("type", utils.PathSearch("type", apiRecord, nil)),
		d.Set("description", utils.PathSearch("description", apiRecord, nil)),
		d.Set("protocol", utils.PathSearch("protocol", apiRecord, nil)),
		d.Set("path", utils.PathSearch("path", apiRecord, nil)),
		d.Set("request_type", utils.PathSearch("request_type", apiRecord, nil)),
		d.Set("manager", utils.PathSearch("manager", apiRecord, nil)),
		d.Set("datasource_config", flattenApiDataSourceConfig(client, workspaceId, utils.PathSearch("datasource_config", apiRecord,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("request_params", flattenApiRequestParams(utils.PathSearch("request_paras", apiRecord, make([]interface{}, 0)).([]interface{}))),
		d.Set("backend_config", flattenApiBackendConfig(utils.PathSearch("backend_config", apiRecord, make([]interface{}, 0)).([]interface{}))),
		// Public attributes.
		d.Set("create_user", utils.PathSearch("create_user", apiRecord, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", apiRecord, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", apiRecord, float64(0)).(float64))/1000, false)),
		// Attribute for shared API.
		d.Set("group_id", utils.PathSearch("group_id", apiRecord, nil)),
		d.Set("status", utils.PathSearch("status", apiRecord, nil)),
		d.Set("host", utils.PathSearch("host", apiRecord, nil)),
		// Attribute for exclusvie API.
		d.Set("hosts", flattenExclusiveApiHosts(utils.PathSearch("hosts", apiRecord, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDataServiceApiUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/service/apis/{api_id}"
		apiId       = d.Id()
		workspaceId = d.Get("workspace_id").(string)
		dlmType     = d.Get("dlm_type").(string)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{api_id}", apiId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     dlmType,
		},
		JSONBody: buildModifyDataServiceApiBodyParams(client, d, workspaceId),
	}

	_, err = client.Request("PUT", createPath, &opt)
	if err != nil {
		return diag.Errorf("error updating API (%s): %s", apiId, err)
	}

	return resourceDataServiceApiRead(ctx, d, meta)
}

func resourceDataServiceApiDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/service/apis/batch-delete"
		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Id()
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	// Due to the lack of an interface for deleting a single API, the corresponding function can only be implemented
	// through the batch deletion interface, and the batch deletion interface cannot be sent to the server
	// at the same time, which will cause an error to be returned.
	config.MutexKV.Lock(workspaceId)
	defer config.MutexKV.Unlock(workspaceId)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     d.Get("dlm_type").(string),
		},
		JSONBody: []string{apiId},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", apiResourceNotFoundCodesForDelete...),
			fmt.Sprintf("error deleting API (%s)", apiId))
	}
	return nil
}

func resourceDataServiceApiImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be '<workspace_id>/<dlm_type>/<id>' or "+
			"'<workspace_id>/<id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(nil, d.Set("workspace_id", parts[0]))
	if len(parts) == 2 {
		d.SetId(parts[1])
	}
	if len(parts) == 3 {
		mErr = multierror.Append(mErr, d.Set("dlm_type", parts[1]))
		d.SetId(parts[2])
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
