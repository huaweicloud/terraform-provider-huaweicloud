package dataarts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/service/apis
// @API DataArtsStudio GET /v1/{project_id}/service/apis/{api_id}
func DataSourceDataServiceApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the API is located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the API belongs.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of DLM engine.`,
			},

			// Query arguments
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The API ID to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The API name to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The API type to be queried.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The API description to be queried.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The API creator to be queried.`,
			},
			"datatable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data table name used by API to be queried.`,
			},

			// Attributes
			"apis": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisElem(),
				Description: `All APIs that match the filter parameters.`,
			},
		},
	}
}

func dataApisElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			// Public attributes.
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API ID, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the API.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request protocol of the API.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API request path.`,
			},
			"request_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request type of the API.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API reviewer.`,
			},
			"datasource_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisDataSourceConfigElemSchema(),
				Description: `The configuration of the API data source.`,
			},
			"request_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisRequestParamsElemSchema(),
				Description: `The parameters of the API request.`,
			},
			"backend_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisBackendConfigElemSchema(),
				Description: `The configuration of the API backend.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API creator.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the API.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the API.`,
			},
			// Attribute for shared API.
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API status.`,
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host configuration for shared API.`,
			},
			// Attribute for exclusvie API.
			"hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The host configuration for exclusive API.`,
			},
		},
	}
	return &sc
}

func dataApisDataSourceConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data source.`,
			},
			"connection_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data connection.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data connection.`,
			},
			"database": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database.`,
			},
			"datatable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data table.`,
			},
			"table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data table.`,
			},
			"queue": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the DLI queue.`,
			},
			"access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access mode for the data.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SQL statements in script access type.`,
			},
			"backend_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisDataSourceConfigBackendParamsElemSchema(),
				Description: `The backend parameters of the API.`,
			},
			"response_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisDataSourceConfigResponseParamsElemSchema(),
				Description: `The response parameters of the API.`,
			},
			"order_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisDataSourceConfigOrderParamsElemSchema(),
				Description: `The order parameters of the API.`,
			},
		},
	}
	return &sc
}

func dataApisDataSourceConfigBackendParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the backend parameter.`,
			},
			"mapping": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the mapping parameter.`,
			},
			"condition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The condition character.`,
			},
		},
	}
	return &sc
}

func dataApisDataSourceConfigResponseParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the response parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the response parameter.`,
			},
			"field": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bound table field for the response parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the response parameter.`,
			},
			"example_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The example value of the response parameter.`,
			},
		},
	}
	return &sc
}

func dataApisDataSourceConfigOrderParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the order parameter.`,
			},
			"field": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The corresponding parameter field for the order parameter.`,
			},
			"optional": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this order parameter is the optional parameter.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sort type of the order parameter.`,
			},
			"order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The order of the sorting parameters.`,
			},
		},
	}
	return &sc
}

func dataApisRequestParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the request parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The position of the request parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the request parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the request parameter.`,
			},
			"necessary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this parameter is the required parameter.`,
			},
			"example_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The example value of the request parameter.`,
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default value of the request parameter.`,
			},
		},
	}
	return &sc
}

func dataApisBackendConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the backend request.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protocol of the backend request.`,
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backend host.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The backend timeout.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backend path.`,
			},
			"backend_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisBackendConfigBackendParamsElemSchema(),
				Description: `The backend parameters of the API.`,
			},
			"constant_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataApisBackendConfigConstantParamsElemSchema(),
				Description: `The backend constant parameters of the API.`,
			},
		},
	}
	return &sc
}

func dataApisBackendConfigBackendParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the request parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The position of the request parameter.`,
			},
			"backend_param_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the corresponding backend parameter.`,
			},
		},
	}
	return &sc
}

func dataApisBackendConfigConstantParamsElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the constant parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the constant parameter.`,
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The position of the constant parameter.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the constant parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the constant parameter.`,
			},
		},
	}
	return &sc
}

func buildDataServiceApisQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, apiName)
	}
	if apiType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&api_type=%v", res, apiType)
	}
	if apiDesc, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, apiDesc)
	}
	if apiCreateUser, ok := d.GetOk("create_user"); ok {
		res = fmt.Sprintf("%s&create_user=%v", res, apiCreateUser)
	}
	if tableName, ok := d.GetOk("datatable"); ok {
		res = fmt.Sprintf("%s&table_name=%v", res, tableName)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func GetDataServiceApis(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		workspaceId = d.Get("workspace_id").(string)
		dlmType     = d.Get("dlm_type").(string)
		httpUrl     = "v1/{project_id}/service/apis"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildDataServiceApisQueryParams(d)

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
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	apiRecords := utils.PathSearch("records", respBody, make([]interface{}, 0))
	if apiId, ok := d.GetOk("api_id"); ok {
		apiRecords = utils.PathSearch(fmt.Sprintf("[?id == '%v']", apiId), apiRecords, make([]interface{}, 0))
	}
	result := make([]interface{}, 0, len(apiRecords.([]interface{})))
	for _, apiRecord := range apiRecords.([]interface{}) {
		apiId := utils.PathSearch("id", apiRecord, "id_not_found").(string)
		apiDetail, err := GetDataServiceApiDetail(client, workspaceId, dlmType, apiId)
		if err != nil {
			log.Printf("error querying API configuration by its ID (%s)", apiId)
			continue
		}
		result = append(result, map[string]interface{}{
			"id": apiId,
			// Public attributes.
			"name":         utils.PathSearch("name", apiDetail, nil),
			"type":         utils.PathSearch("type", apiDetail, nil),
			"description":  utils.PathSearch("description", apiDetail, nil),
			"protocol":     utils.PathSearch("protocol", apiDetail, nil),
			"path":         utils.PathSearch("path", apiDetail, nil),
			"request_type": utils.PathSearch("request_type", apiDetail, nil),
			"manager":      utils.PathSearch("manager", apiDetail, nil),
			"datasource_config": flattenApiDataSourceConfig(client, workspaceId, utils.PathSearch("datasource_config", apiDetail,
				make(map[string]interface{})).(map[string]interface{})),
			"request_params": flattenApiRequestParams(utils.PathSearch("request_paras", apiDetail, make([]interface{}, 0)).([]interface{})),
			"backend_config": flattenApiBackendConfig(utils.PathSearch("backend_config", apiDetail, make([]interface{}, 0)).([]interface{})),
			"create_user":    utils.PathSearch("create_user", apiRecord, nil),
			"created_at":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", apiDetail, float64(0)).(float64))/1000, false),
			"updated_at":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", apiDetail, float64(0)).(float64))/1000, false),
			// Attribute for shared API.
			"group_id": utils.PathSearch("group_id", apiDetail, nil),
			"status":   utils.PathSearch("status", apiDetail, nil),
			"host":     utils.PathSearch("host", apiDetail, nil),
			// Attribute for exclusvie API.
			"hosts": flattenExclusiveApiHosts(utils.PathSearch("hosts", apiDetail, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result, nil
}

func GetDataServiceApiDetail(client *golangsdk.ServiceClient, workspaceId, dlmType, apiId string) (interface{}, error) {
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

func dataSourceDataServiceApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	apiRecords, err := GetDataServiceApis(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error getting Data Service API (%s) for DataArts Studio", d.Id()))
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("apis", apiRecords),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
