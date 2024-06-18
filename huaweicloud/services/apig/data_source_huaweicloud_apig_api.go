package apig

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

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

func dataSourceApiRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
	)

	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	resp, err := apis.Get(client, instanceId, apiId).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("type", analyseApiType(resp.Type)),
		d.Set("request_protocol", resp.ReqProtocol),
		d.Set("request_method", resp.ReqMethod),
		d.Set("request_path", resp.ReqURI),
		d.Set("security_authentication", resp.AuthType),
		d.Set("authorizer_id", resp.AuthorizerId),
		d.Set("tags", resp.Tags),
		d.Set("request_params", flattenRequestParams(resp.ReqParams)),
		d.Set("backend_params", flattenBackendParams(resp.BackendParams)),
		d.Set("group_id", resp.GroupId),
		d.Set("group_name", resp.GroupName),
		d.Set("group_version", resp.GroupVersion),
		d.Set("env_id", resp.RunEnvId),
		d.Set("env_name", resp.RunEnvName),
		d.Set("publish_id", resp.PublishId),
		d.Set("backend_type", resp.BackendType),
		d.Set("cors", resp.Cors),
		d.Set("body_description", resp.BodyDescription),
		d.Set("description", resp.Description),
		d.Set("matching", analyseApiMatchMode(resp.MatchMode)),
		d.Set("response_id", resp.ResponseId),
		d.Set("success_response", resp.ResultNormalSample),
		d.Set("failure_response", resp.ResultFailureSample),
		d.Set("simple_authentication", analyseAppSimpleAuth(resp.AuthOpt)),
		d.Set("mock", flattenMock(resp.MockInfo)),
		d.Set("mock_policy", flattenMockPolicies(resp.PolicyMocks)),
		d.Set("func_graph", flattenFuncGraph(resp.FuncInfo)),
		d.Set("func_graph_policy", flattenFuncGraphPolicies(resp.PolicyFunctions)),
		d.Set("web", flattenWeb(resp.WebInfo, d.Get("web.0.ssl_enable").(bool))),
		d.Set("web_policy", flattenWebPolicies(resp.PolicyWebs)),
		d.Set("registered_at", flattenTimeToRFC3339(resp.RegisterTime)),
		d.Set("updated_at", flattenTimeToRFC3339(resp.UpdateTime)),
		d.Set("published_at", flattenPulishTime(resp.PublishTime)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRequestParams(reqParams []apis.ReqParamResp) []map[string]interface{} {
	if len(reqParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
			"id":           v.ID,
			"name":         v.Name,
			"location":     v.Location,
			"type":         v.Type,
			"required":     parseObjectEnabled(v.Required),
			"passthrough":  parseObjectEnabled(v.PassThrough),
			"enumeration":  v.Enumerations,
			"example":      v.SampleValue,
			"default":      v.DefaultValue,
			"description":  v.Description,
			"valid_enable": v.ValidEnable,
		}
		switch v.Type {
		case string(ParamTypeNumber):
			param["maximum"] = v.MaxNum
			param["minimum"] = v.MinNum
		case string(ParamTypeString):
			param["maximum"] = v.MaxSize
			param["minimum"] = v.MinSize
		}
		result[i] = param
	}
	return result
}

func flattenBackendParams(backendParams []apis.BackendParamResp) []map[string]interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(backendParams))
	for i, v := range backendParams {
		origin := v.Origin
		paramAuthType, paramValue := analyseBackendParameterValue(v.Origin, v.Value)
		param := map[string]interface{}{
			"id":                v.ID,
			"request_id":        v.ReqParamId,
			"type":              origin,
			"name":              v.Name,
			"location":          v.Location,
			"value":             paramValue,
			"system_param_type": paramAuthType,
			"description":       v.Description,
		}
		result[i] = param
	}
	return result
}

func flattenMock(mockResp apis.Mock) []map[string]interface{} {
	mockInfo := flattenMockStructure(mockResp)
	if mockInfo == nil {
		return nil
	}

	mockInfo[0]["id"] = mockResp.ID
	return mockInfo
}

func flattenMockPolicies(policies []apis.PolicyMockResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"id":             policy.ID,
			"name":           policy.Name,
			"status_code":    policy.StatusCode,
			"response":       policy.ResultContent,
			"effective_mode": policy.EffectMode,
			"authorizer_id":  policy.AuthorizerId,
			"backend_params": flattenBackendParams(policy.BackendParams),
			"conditions":     flattenConditions(policy.Conditions),
		}
	}

	return result
}

func flattenConditions(conditions []apis.APIConditionBase) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		result[i] = map[string]interface{}{
			"id":                       v.ID,
			"source":                   v.ConditionOrigin,
			"param_name":               v.ReqParamName,
			"sys_name":                 v.SysParamName,
			"cookie_name":              v.CookieParamName,
			"frontend_authorizer_name": v.FrontendAuthorizerParamName,
			"type":                     analyseConditionType(v.ConditionType),
			"value":                    v.ConditionValue,
			"request_id":               v.ReqParamId,
			"request_location":         v.ReqParamLocation,
		}
	}
	return result
}

func flattenFuncGraph(funcResp apis.FuncGraph) []map[string]interface{} {
	functionGraphInfo := flattenFuncGraphStructure(funcResp)
	if functionGraphInfo == nil {
		return nil
	}

	functionGraphInfo[0]["id"] = funcResp.ID
	return functionGraphInfo
}

func flattenFuncGraphPolicies(policies []apis.PolicyFuncGraphResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"id":                 policy.ID,
			"name":               policy.Name,
			"function_urn":       policy.FunctionUrn,
			"version":            policy.Version,
			"function_alias_urn": policy.FunctionAliasUrn,
			"network_type":       policy.NetworkType,
			"request_protocol":   policy.RequestProtocol,
			"invocation_type":    policy.InvocationType,
			"effective_mode":     policy.EffectMode,
			"timeout":            policy.Timeout,
			"authorizer_id":      policy.AuthorizerId,
			"backend_params":     flattenBackendParams(policy.BackendParams),
			"conditions":         flattenConditions(policy.Conditions),
		}
	}

	return result
}

func flattenWeb(webResp apis.Web, sslEnabled bool) []map[string]interface{} {
	webInfo := flattenWebStructure(webResp, sslEnabled)
	if webInfo == nil {
		return nil
	}

	webInfo[0]["id"] = webResp.ID
	return webInfo
}

func flattenWebPolicies(policies []apis.PolicyWebResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		retryCount := policy.RetryCount
		policy := map[string]interface{}{
			"id":               policy.ID,
			"name":             policy.Name,
			"request_protocol": policy.ReqProtocol,
			"request_method":   policy.ReqMethod,
			"effective_mode":   policy.EffectMode,
			"path":             policy.ReqURI,
			"host_header":      policy.VpcChannelInfo.VpcChannelProxyHost,
			"vpc_channel_id":   policy.VpcChannelInfo.VpcChannelId,
			"backend_address":  policy.DomainURL,
			"timeout":          policy.Timeout,
			"retry_count":      utils.StringToInt(&retryCount),
			"authorizer_id":    policy.AuthorizerId,
			"backend_params":   flattenBackendParams(policy.BackendParams),
			"conditions":       flattenConditions(policy.Conditions),
		}

		result[i] = policy
	}

	return result
}
