package apig

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}
func DataSourceApi() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the API belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the API.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the API.`,
			},
			"request_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request method of the API.",
			},
			"request_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request address of the API.",
			},
			"request_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request protocol of the API.",
			},
			"security_authentication": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security authentication mode of the API request.",
			},
			"simple_authentication": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the authentication of the application code is enabled.",
			},
			"authorizer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the authorizer to which the API request used.`,
			},
			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of tags configuration.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The group ID corresponding to the API.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The group name corresponding to the API.`,
			},
			"group_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of group corresponding to the API.`,
			},
			"request_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the request parameter.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the request parameter.",
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this parameter is required.",
						},
						"passthrough": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to transparently transfer the parameter.",
						},
						"enumeration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enumerated value.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Where this parameter is located.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter type.",
						},
						"maximum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value or length (string parameter) for parameter.",
						},
						"minimum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value or length (string parameter) for parameter.",
						},
						"example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter example.",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default value of the parameter.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter description.",
						},
						"valid_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable the parameter validation.",
						},
					},
				},
				Description: "The configuration list of the front-end parameters.",
			},
			"backend_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        backendParamSchemaDataSource(),
				Description: "The configuration list of the backend parameters.",
			},
			"body_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the API request body.",
			},
			"cors": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether CORS is supported.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the API.",
			},
			"matching": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The matching mode of the API.",
			},
			"response_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The response ID of the corresponding APIG group.",
			},
			"success_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The example response for a successful request.",
			},
			"failure_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The example response for a failure request.",
			},
			"mock": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the mock backend configuration.",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The custom status code of the mock response.",
						},
						"response": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The response of the mock backend configuration.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The mock backend details.",
			},
			"mock_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the mock backend policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend policy name.",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The custom status code of the mock response.",
						},
						"response": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The response of the backend policy.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        policyConditionSchemaDataSource(),
							Description: "The policy conditions.",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The effective mode of the backend policy.",
						},
						"backend_params": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        backendParamSchemaDataSource(),
							Description: "The configuration list of backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The policy backends of the mock.",
			},
			"func_graph": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the FunctionGraph backend configuration.",
						},
						"function_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URN of the FunctionGraph function.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the FunctionGraph function.",
						},
						"function_alias_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias URN of the FunctionGraph function.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network architecture (framework) type of the FunctionGraph function.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request protocol of the FunctionGraph function.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout for API requests to backend service.",
						},
						"invocation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The invocation type.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The FunctionGraph backend details.",
			},
			"func_graph_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the FunctionGraph backend policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the backend policy.",
						},
						"function_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URN of the FunctionGraph function.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the FunctionGraph function.",
						},
						"function_alias_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias URN of the FunctionGraph function.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network (framework) type of the FunctionGraph function.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request protocol of the FunctionGraph function.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        policyConditionSchemaDataSource(),
							Description: "The policy conditions.",
						},
						"invocation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The invocation mode of the FunctionGraph function.",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The effective mode of the backend policy.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout for API requests to backend service.",
						},
						"backend_params": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        backendParamSchemaDataSource(),
							Description: "The configaiton list of the backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The policy backends of the FunctionGraph.",
			},
			"web": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend configuration.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend request path.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The proxy host header.",
						},
						"vpc_channel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC channel ID.",
						},
						"backend_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service address.",
						},
						"request_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend request method of the API.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The web protocol type of the API request.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout for API requests to backend service.",
						},
						"retry_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of retry attempts to request the backend service.",
						},
						"ssl_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable two-way authentication.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The web backend details.",
			},
			"web_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the web policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the web policy.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend request address.",
						},
						"request_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend request method of the API.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend request protocol.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        policyConditionSchemaDataSource(),
							Description: "The policy conditions.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The proxy host header.",
						},
						"vpc_channel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC channel ID.",
						},
						"backend_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The backend service address",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The effective mode of the backend policy.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout for API requests to backend service.",
						},
						"retry_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of retry attempts to request the backend service.",
						},
						"backend_params": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        backendParamSchemaDataSource(),
							Description: "The configuration list of the backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The policy backends of the web.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"publish_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of publish corresponding to the API.`,
			},
			"backend_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backend type of the API.`,
			},
			// Formats the time according to the local computer's time.
			// The format is `yyyy-MM-ddTHH:mm:ss{timezone}`, e.g. `2006-01-02 15:04:05+08:00`.
			"registered_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registered time of the API, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the API, in RFC3339 format.",
			},
			"published_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The published time of the API, in RFC3339 format.",
			},
		},
	}
}

func policyConditionSchemaDataSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the backend policy condition.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value corresponding to the parameter name.",
			},
			"param_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request parameter name.",
			},
			"sys_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gateway built-in parameter name.",
			},
			"cookie_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cookie parameter name.",
			},
			"frontend_authorizer_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The frontend authentication parameter name.",
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the backend policy.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The condition type of the backend policy.",
			},
			"request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the corresponding request parameter.",
			},
			"request_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location of the corresponding request parameter.",
			},
		},
	}
}

func backendParamSchemaDataSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the backend parameter.",
			},
			"request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the corresponding request parameter.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of parameter.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of parameter.",
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Where the parameter is located.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the parameter.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the constant or system parameter.",
			},
			"system_param_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the system parameter.",
			},
		},
	}
}

func flattenDataSourceRequestParams(reqParams []interface{}) []map[string]interface{} {
	result := flattenApiRequestParams(reqParams, nil)
	if result == nil {
		return nil
	}

	for i, item := range reqParams {
		result[i]["id"] = utils.PathSearch("id", item, nil)
	}
	return result
}

func flattenDataSourceBackendParams(backendParams []interface{}) []map[string]interface{} {
	result := flattenBackendParameters(backendParams)
	if result == nil {
		return nil
	}

	for i, item := range backendParams {
		result[i]["id"] = utils.PathSearch("id", item, nil)
		result[i]["request_id"] = utils.PathSearch("req_param_id", item, nil)
	}
	return result
}

func flattenDataSourcePolicyConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		result[i] = map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"source":                   utils.PathSearch("condition_origin", v, nil),
			"param_name":               utils.PathSearch("req_param_name", v, nil),
			"sys_name":                 utils.PathSearch("sys_param_name", v, nil),
			"cookie_name":              utils.PathSearch("cookie_param_name", v, nil),
			"frontend_authorizer_name": utils.PathSearch("frontend_authorizer_param_name", v, nil),
			"type":                     analyseConditionType(utils.PathSearch("condition_type", v, "").(string)),
			"value":                    utils.PathSearch("condition_value", v, nil),
			"request_id":               utils.PathSearch("req_param_id", v, nil),
			"request_location":         utils.PathSearch("req_param_location", v, nil),
		}
	}
	return result
}

func flattenDataSourceMock(mockResp interface{}) []map[string]interface{} {
	result := flattenMockStructure(mockResp)
	if result == nil {
		return nil
	}

	result[0]["id"] = utils.PathSearch("id", mockResp, nil)
	return result
}

func flattenDataSourceMockPolicies(policies []interface{}) []map[string]interface{} {
	result := flattenMockPolicy(policies, nil)
	if result == nil {
		return nil
	}

	for i, policy := range policies {
		result[i]["id"] = utils.PathSearch("id", policy, nil)
		result[i]["conditions"] = flattenDataSourcePolicyConditions(utils.PathSearch("conditions", policy,
			make([]interface{}, 0)).([]interface{}))
		result[i]["backend_params"] = flattenDataSourceBackendParams(utils.PathSearch("backend_params", policy,
			make([]interface{}, 0)).([]interface{}))
	}
	return result
}

func flattenDataSourceFuncGraph(funcResp interface{}) []map[string]interface{} {
	result := flattenFuncGraphStructure(funcResp)
	if result == nil {
		return nil
	}

	result[0]["id"] = utils.PathSearch("id", funcResp, nil)
	return result
}

func flattenDataSourceFuncGraphPolicies(policies []interface{}) []map[string]interface{} {
	result := flattenFuncGraphPolicy(policies, nil)
	if result == nil {
		return nil
	}

	for i, policy := range policies {
		result[i]["id"] = utils.PathSearch("id", policy, nil)
		result[i]["conditions"] = flattenDataSourcePolicyConditions(utils.PathSearch("conditions", policy,
			make([]interface{}, 0)).([]interface{}))
		result[i]["backend_params"] = flattenDataSourceBackendParams(utils.PathSearch("backend_params", policy,
			make([]interface{}, 0)).([]interface{}))
	}
	return result
}

func flattenDataSourceWeb(webResp interface{}, sslEnabled bool) []map[string]interface{} {
	result := flattenWebStructure(webResp, sslEnabled)
	if result == nil {
		return nil
	}

	result[0]["id"] = utils.PathSearch("id", webResp, nil)
	return result
}

func flattenDataSourceWebPolicies(policies []interface{}) []map[string]interface{} {
	result := flattenWebPolicy(policies, nil)
	if result == nil {
		return nil
	}

	for i, policy := range policies {
		result[i]["id"] = utils.PathSearch("id", policy, nil)
		result[i]["conditions"] = flattenDataSourcePolicyConditions(utils.PathSearch("conditions", policy,
			make([]interface{}, 0)).([]interface{}))
		result[i]["backend_params"] = flattenDataSourceBackendParams(utils.PathSearch("backend_params", policy,
			make([]interface{}, 0)).([]interface{}))
	}
	return result
}

func dataSourceApiRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := GetApiById(client, instanceId, apiId)
	if err != nil {
		return diag.Errorf("error querying API (%s): %s", apiId, err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("type", parseApiType(int(utils.PathSearch("type", respBody, float64(0)).(float64)))),
		d.Set("request_protocol", utils.PathSearch("req_protocol", respBody, nil)),
		d.Set("request_method", utils.PathSearch("req_method", respBody, nil)),
		d.Set("request_path", utils.PathSearch("req_uri", respBody, nil)),
		d.Set("security_authentication", utils.PathSearch("auth_type", respBody, nil)),
		d.Set("authorizer_id", utils.PathSearch("authorizer_id", respBody, nil)),
		d.Set("tags", utils.PathSearch("tags", respBody, nil)),
		d.Set("request_params", flattenDataSourceRequestParams(utils.PathSearch("req_params", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("backend_params", flattenDataSourceBackendParams(utils.PathSearch("backend_params", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("group_id", utils.PathSearch("group_id", respBody, nil)),
		d.Set("group_name", utils.PathSearch("group_name", respBody, nil)),
		d.Set("group_version", utils.PathSearch("group_version", respBody, nil)),
		d.Set("env_id", utils.PathSearch("run_env_id", respBody, nil)),
		d.Set("env_name", utils.PathSearch("run_env_name", respBody, nil)),
		d.Set("publish_id", utils.PathSearch("publish_id", respBody, nil)),
		d.Set("backend_type", utils.PathSearch("backend_type", respBody, nil)),
		d.Set("cors", utils.PathSearch("cors", respBody, nil)),
		d.Set("body_description", utils.PathSearch("body_remark", respBody, nil)),
		d.Set("description", utils.PathSearch("remark", respBody, nil)),
		d.Set("matching", analyseApiMatchMode(utils.PathSearch("match_mode", respBody, "").(string))),
		d.Set("response_id", utils.PathSearch("response_id", respBody, nil)),
		d.Set("success_response", utils.PathSearch("result_normal_sample", respBody, nil)),
		d.Set("failure_response", utils.PathSearch("result_failure_sample", respBody, nil)),
		d.Set("simple_authentication", analyseAppSimpleAuth(utils.PathSearch("auth_opt", respBody, nil))),
		d.Set("mock", flattenDataSourceMock(utils.PathSearch("mock_info", respBody, nil))),
		d.Set("mock_policy", flattenDataSourceMockPolicies(utils.PathSearch("policy_mocks", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("func_graph", flattenDataSourceFuncGraph(utils.PathSearch("func_info", respBody, nil))),
		d.Set("func_graph_policy", flattenDataSourceFuncGraphPolicies(utils.PathSearch("policy_functions", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("web", flattenDataSourceWeb(utils.PathSearch("backend_api", respBody, nil), d.Get("web.0.ssl_enable").(bool))),
		d.Set("web_policy", flattenDataSourceWebPolicies(utils.PathSearch("policy_https", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("registered_at", flattenTimeToRFC3339(utils.PathSearch("register_time", respBody, "").(string))),
		d.Set("updated_at", flattenTimeToRFC3339(utils.PathSearch("update_time", respBody, "").(string))),
		d.Set("published_at", flattenPulishTime(utils.PathSearch("publish_time", respBody, "").(string))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
