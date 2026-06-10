package apig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type (
	ApiType         string
	RequestMethod   string
	ApiAuthType     string
	ParamLocation   string
	ParamType       string
	MatchMode       string
	InvacationType  string
	EffectiveMode   string
	ConditionSource string
	ConditionType   string
	ParameterType   string
	SystemParamType string
	BackendType     string
	AppCodeAuthType string
	ProtocolType    string
)

const (
	ApiTypePublic  ApiType = "Public"
	ApiTypePrivate ApiType = "Private"

	RequestMethodGet     RequestMethod = "GET"
	RequestMethodPost    RequestMethod = "POST"
	RequestMethodPut     RequestMethod = "PUT"
	RequestMethodDelete  RequestMethod = "DELETE"
	RequestMethodHead    RequestMethod = "HEAD"
	RequestMethodPatch   RequestMethod = "PATCH"
	RequestMethodOptions RequestMethod = "OPTIONS"
	RequestMethodAny     RequestMethod = "ANY"

	ApiAuthTypeNone       ApiAuthType = "NONE"
	ApiAuthTypeApp        ApiAuthType = "APP"
	ApiAuthTypeIam        ApiAuthType = "IAM"
	ApiAuthTypeAuthorizer ApiAuthType = "AUTHORIZER"

	ParamLocationPath   ParamLocation = "PATH"
	ParamLocationHeader ParamLocation = "HEADER"
	ParamLocationQuery  ParamLocation = "QUERY"

	ParamTypeString ParamType = "STRING"
	ParamTypeNumber ParamType = "NUMBER"

	MatchModePrefix MatchMode = "Prefix"
	MatchModeExact  MatchMode = "Exact"

	InvacationTypeAsync InvacationType = "async"
	InvacationTypeSync  InvacationType = "sync"

	EffectiveModeAll EffectiveMode = "ALL"
	EffectiveModeAny EffectiveMode = "ANY"

	ConditionTypeEqual      ConditionType = "Equal"
	ConditionTypeEnumerated ConditionType = "Enumerated"
	ConditionTypeMatching   ConditionType = "Matching"

	ParameterTypeRequest  ParameterType = "REQUEST"
	ParameterTypeConstant ParameterType = "CONSTANT"
	ParameterTypeSystem   ParameterType = "SYSTEM"

	SystemParamTypeFrontend SystemParamType = "frontend"
	SystemParamTypeBackend  SystemParamType = "backend"
	SystemParamTypeInternal SystemParamType = "internal"

	BackendTypeHttp     BackendType = "HTTP"
	BackendTypeFunction BackendType = "FUNCTION"
	BackendTypeMock     BackendType = "MOCK"

	AppCodeAuthTypeDisable AppCodeAuthType = "DISABLE"
	AppCodeAuthTypeEnable  AppCodeAuthType = "HEADER"

	ProtocolTypeTCP   ProtocolType = "TCP"
	ProtocolTypeHTTP  ProtocolType = "HTTP"
	ProtocolTypeHTTPS ProtocolType = "HTTPS"
	ProtocolTypeBoth  ProtocolType = "BOTH"
	ProtocolTypeGPRCS ProtocolType = "GRPCS"

	vpcChannelEnabled  int = 1
	vpcChannelDisabled int = 2
)

var (
	matching = map[string]string{
		string(MatchModePrefix): "SWA",
		string(MatchModeExact):  "NORMAL",
	}
	conditionType = map[string]string{
		string(ConditionTypeEqual):      "exact",
		string(ConditionTypeEnumerated): "enum",
		string(ConditionTypeMatching):   "pattern",
	}
)

var apiNonUpdatableParams = []string{
	"instance_id",
	"group_id",
	"mock",
	"func_graph",
	"web",
}

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis
func ResourceApigAPIV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiCreate,
		ReadContext:   resourceApiRead,
		UpdateContext: resourceApiUpdate,
		DeleteContext: resourceApiDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApiImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(apiNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the API is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the API belongs.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API group to which the API belongs.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The API type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The API name.`,
			},
			"request_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request method of the API.`,
			},
			"request_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request address.`,
			},
			"request_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request protocol of the API request.`,
			},

			// Optional parameters.
			"security_authentication": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(ApiAuthTypeNone),
				Description: `The security authentication mode of the API request.`,
			},
			"simple_authentication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the authentication of the application code is enabled.",
			},
			"authorizer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the authorizer to which the API request used.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of tags configuration.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The content type of the request body.",
			},
			"is_send_fg_body_base64": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to perform Base64 encoding on the body for interaction with FunctionGraph.",
			},
			"request_params": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 50,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the request parameter.",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether this parameter is required.",
						},
						"passthrough": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to transparently transfer the parameter.",
						},
						"enumeration": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The enumerated value.",
						},
						"location": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(ParamLocationPath),
							Description: "Where this parameter is located.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(ParamTypeString),
							Description: "The parameter type.",
						},
						"maximum": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The maximum value or length (string parameter) for parameter.",
						},
						"minimum": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The minimum value or length (string parameter) for parameter.",
						},
						"example": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The parameter example.",
						},
						"default": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The default value of the parameter.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The parameter description.",
						},
						"valid_enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Whether to enable the parameter validation.",
						},
						"orchestrations": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The list of orchestration rules that parameter used.",
						},
					},
				},
				Description: "The configurations of the front-end parameters.",
			},
			"backend_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    50,
				Elem:        backendParamSchemaResource(),
				Set:         resourceBackendParamtersHash,
				Description: "The configurations of the backend parameters.",
			},
			"body_description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The description of the API request body, which can be an example request body, media " +
					"type or parameters.",
			},
			"cors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether CORS is supported.",
			},
			"sampling_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sampling strategy of the link trace.",
			},
			"sampling_param": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sampling parameter of the link trace.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API description.",
			},
			"matching": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     MatchModeExact,
				Description: "The matching mode of the API.",
			},
			"response_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the custom response that API used.",
			},
			"success_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The example response for a successful request.",
			},
			"failure_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The example response for a failure request.",
			},
			"mock": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"func_graph", "web"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The custom status code of the mock response.",
						},
						"response": {
							Type:     schema.TypeString,
							Optional: true,
							Description: utils.SchemaDesc(
								"The response content of the mock.",
								utils.SchemaDescInput{
									Required: true,
								},
							),
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The mock backend details.",
			},
			"func_graph": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_urn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URN of the FunctionGraph function.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The version of the FunctionGraph function.",
						},
						"function_alias_urn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The alias URN of the FunctionGraph function.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The network architecture (framework) type of the FunctionGraph function.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The request protocol of the FunctionGraph function.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5000,
							Description: "The timeout for API requests to backend service.",
						},
						"invocation_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(InvacationTypeSync),
							Description: "The invocation type.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The FunctionGraph backend details.",
			},
			"web": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The backend request path.",
						},
						"host_header": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"web.0.backend_address"},
							Description:   "The proxy host header.",
						},
						"vpc_channel_id": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"web.0.backend_address"},
							Description:  "The VPC channel ID.",
						},
						"backend_address": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "The backend service address, which consists of a domain name or IP " +
								"address, and a port number.",
						},
						"request_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend request method of the API.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(ProtocolTypeHTTPS),
							Description: "The web protocol type of the API request.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5000,
							Description: "The timeout for API requests to backend service.",
						},
						"retry_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     -1,
							Description: "The number of retry attempts to request the backend service.",
						},
						"ssl_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to enable two-way authentication.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The web backend details.",
			},
			"mock_policy": {
				Type:          schema.TypeList,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"func_graph", "web", "func_graph_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The backend policy name.",
						},
						"conditions": {
							Type:        schema.TypeSet,
							Required:    true,
							MaxItems:    5,
							Elem:        policyConditionSchemaResource(),
							Description: "The policy conditions.",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The custom status code of the mock response.",
						},
						"response": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The response content of the mock.",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(EffectiveModeAny),
							Description: "The effective mode of the backend policy.",
						},
						"backend_params": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        backendParamSchemaResource(),
							Set:         resourceBackendParamtersHash,
							Description: "The configuration list of backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The mock policy backends.",
			},
			"func_graph_policy": {
				Type:          schema.TypeList,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "web", "mock_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the backend policy.",
						},
						"function_urn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URN of the FunctionGraph function.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The version of the FunctionGraph function.",
						},
						"function_alias_urn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The alias URN of the FunctionGraph function.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The network (framework) type of the FunctionGraph function.",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The request protocol of the FunctionGraph function.",
						},
						"conditions": {
							Type:        schema.TypeSet,
							Required:    true,
							MaxItems:    5,
							Elem:        policyConditionSchemaResource(),
							Description: "The policy conditions.",
						},
						"invocation_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(InvacationTypeSync),
							Description: "The invocation mode of the FunctionGraph function.",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(EffectiveModeAny),
							Description: "The effective mode of the backend policy.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5000,
							Description: "The timeout for API requests to backend service.",
						},
						"backend_params": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        backendParamSchemaResource(),
							Set:         resourceBackendParamtersHash,
							Description: "The configaiton list of the backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
						// Deprecated arguments
						"invocation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Description: utils.SchemaDesc(
								`The invocation mode of the FunctionGraph function.`,
								utils.SchemaDescInput{
									Required:   true,
									Deprecated: true,
								},
							),
						},
					},
				},
				Description: "The policy backends of the FunctionGraph function.",
			},
			"web_policy": {
				Type:          schema.TypeList,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "func_graph", "mock_policy", "func_graph_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the web policy.",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The backend request address.",
						},
						"request_method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The backend request method of the API.",
						},
						"conditions": {
							Type:        schema.TypeSet,
							Required:    true,
							MaxItems:    5,
							Elem:        policyConditionSchemaResource(),
							Description: "The policy conditions.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The proxy host header.",
						},
						"vpc_channel_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The VPC channel ID.",
						},
						"backend_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend service address",
						},
						"request_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The backend request protocol.",
						},
						"effective_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     string(EffectiveModeAny),
							Description: "The effective mode of the backend policy.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5000,
							Description: "The timeout for API requests to backend service.",
						},
						"retry_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     -1,
							Description: "The number of retry attempts to request the backend service.",
						},
						"backend_params": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        backendParamSchemaResource(),
							Set:         resourceBackendParamtersHash,
							Description: "The configuration list of the backend parameters.",
						},
						"authorizer_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend custom authorization.",
						},
					},
				},
				Description: "The web policy backends.",
			},

			// Attributes.
			"registered_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The registered time of the API.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the API.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},

			// Internal attributes.
			"request_params_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the request parameter.",
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of request parameters that used to reorder the 'request_params' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"func_graph_policy_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the function graph policy.",
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of function graph policies that used to reorder the 'func_graph_policy' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"web_policy_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the web policy.",
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of web policies that used to reorder the 'web_policy' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"mock_policy_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the mock policy.",
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of mock policies that used to reorder the 'mock_policy' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func resourceBackendParamtersHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["type"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
	}
	if m["name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	return hashcode.String(buf.String())
}

func policyConditionSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the backend policy.",
			},
			"param_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The request parameter name.",
			},
			"sys_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The gateway built-in parameter name.",
			},
			"cookie_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cookie parameter name.",
			},
			"frontend_authorizer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The frontend authentication parameter name.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "param",
				Description: "The type of the backend policy.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(ConditionTypeEqual),
				Description: "The condition type.",
			},
			"mapped_param_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a parameter generated after orchestration.",
			},
			"mapped_param_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of a parameter generated after orchestration.",
			},
		},
	}
}

func backendParamSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter type.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter name.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Where the parameter is located.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the parameter",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the parameter.",
			},
			"system_param_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildApiType(t string) int {
	switch t {
	case string(ApiTypePublic):
		return 1
	default:
		return 2 // Private
	}
}

func isObjectEnabled(isEnabled bool) int {
	if isEnabled {
		return vpcChannelEnabled
	}
	return vpcChannelDisabled
}

func buildMockStructure(mocks []interface{}) map[string]interface{} {
	if len(mocks) < 1 {
		return nil
	}

	mockMap := mocks[0].(map[string]interface{})
	return map[string]interface{}{
		"status_code":    mockMap["status_code"],
		"result_content": utils.ValueIgnoreEmpty(mockMap["response"]),
		"authorizer_id":  utils.ValueIgnoreEmpty(mockMap["authorizer_id"]),
	}
}

func buildFuncGraphStructure(funcGraphs []interface{}) map[string]interface{} {
	if len(funcGraphs) < 1 {
		return nil
	}

	funcMap := funcGraphs[0].(map[string]interface{})
	return map[string]interface{}{
		"function_urn":    funcMap["function_urn"],
		"alias_urn":       funcMap["function_alias_urn"],
		"network_type":    funcMap["network_type"],
		"timeout":         funcMap["timeout"],
		"invocation_type": funcMap["invocation_type"],
		"version":         funcMap["version"],
		"authorizer_id":   utils.ValueIgnoreEmpty(funcMap["authorizer_id"]),
		"req_protocol":    funcMap["request_protocol"],
	}
}

func buildWebStructure(webs []interface{}) map[string]interface{} {
	if len(webs) < 1 {
		return nil
	}

	webMap := webs[0].(map[string]interface{})
	webResp := map[string]interface{}{
		"req_uri":           webMap["path"],
		"req_method":        webMap["request_method"],
		"req_protocol":      webMap["request_protocol"],
		"timeout":           webMap["timeout"],
		"enable_client_ssl": webMap["ssl_enable"],
		"authorizer_id":     utils.ValueIgnoreEmpty(webMap["authorizer_id"]),
		"retry_count":       strconv.Itoa(webMap["retry_count"].(int)),
	}
	// If vpc_channel_id is empty, the backend address is used.
	if chanId, ok := webMap["vpc_channel_id"]; ok && chanId != "" {
		webResp["vpc_channel_status"] = vpcChannelEnabled
		webResp["vpc_channel_info"] = map[string]interface{}{
			"vpc_channel_id":         chanId,
			"vpc_channel_proxy_host": webMap["host_header"],
		}
	} else {
		webResp["vpc_channel_status"] = vpcChannelDisabled
		webResp["url_domain"] = webMap["backend_address"]
	}

	return webResp
}

func buildRequestParameters(requests []interface{}) []map[string]interface{} {
	if len(requests) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(requests))
	for _, v := range requests {
		paramMap := v.(map[string]interface{})
		paramType := paramMap["type"].(string)
		param := map[string]interface{}{
			"type":           paramType,
			"name":           paramMap["name"],
			"required":       isObjectEnabled(paramMap["required"].(bool)),
			"location":       paramMap["location"],
			"remark":         utils.ValueIgnoreEmpty(paramMap["description"]),
			"enumerations":   utils.ValueIgnoreEmpty(paramMap["enumeration"]),
			"pass_through":   isObjectEnabled(paramMap["passthrough"].(bool)),
			"default_value":  utils.ValueIgnoreEmpty(paramMap["default"]),
			"sample_value":   paramMap["example"],
			"valid_enable":   paramMap["valid_enable"],
			"orchestrations": utils.ValueIgnoreEmpty(utils.ExpandToStringList(paramMap["orchestrations"].([]interface{}))),
		}
		switch paramType {
		case string(ParamTypeNumber):
			param["max_num"] = utils.ValueIgnoreEmpty(paramMap["maximum"])
			param["min_num"] = utils.ValueIgnoreEmpty(paramMap["minimum"])
		case string(ParamTypeString):
			param["max_size"] = utils.ValueIgnoreEmpty(paramMap["maximum"])
			param["min_size"] = utils.ValueIgnoreEmpty(paramMap["minimum"])
		}
		result = append(result, param)
	}
	return result
}

func buildBackendParameterValue(origin, value, paramAuthType string) string {
	// The internal parameters of the system parameters include as below:
	internalParams := []string{
		"sourceIp", "stage", "apiId", "appId", "requestId", "serverAddr", "serverName", "handleTime", "providerAppId",
	}

	if origin == "SYSTEM" {
		if paramAuthType == string(SystemParamTypeFrontend) || paramAuthType == string(SystemParamTypeBackend) {
			// The fornt-end or backend format is used to construct.
			return fmt.Sprintf("$context.authorizer.%s.%s", paramAuthType, value)
		}
		if utils.StrSliceContains(internalParams, value) {
			// If the system parameters are configured as internal parameters, the internal format is used to construct.
			return fmt.Sprintf("$context.%s", value)
		}
	}
	return value
}

// For backend API, the parameters contains request parameters and constant parameters.
func buildBackendParameters(backends *schema.Set) ([]map[string]interface{}, error) {
	if backends.Len() < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, backends.Len())
	for _, v := range backends.List() {
		origin := utils.PathSearch("type", v, "").(string)
		if origin == string(ParameterTypeSystem) && utils.PathSearch("system_param_type", v, "").(string) == "" {
			return nil, fmt.Errorf("The 'system_param_type' must set if parameter type is 'SYSTEM'")
		}
		param := map[string]interface{}{
			"origin":   origin,
			"name":     utils.PathSearch("name", v, ""),
			"location": utils.PathSearch("location", v, ""),
			"value": buildBackendParameterValue(origin, utils.PathSearch("value", v, "").(string),
				utils.PathSearch("system_param_type", v, "").(string)),
		}

		if origin != string(ParameterTypeRequest) {
			param["remark"] = utils.ValueIgnoreEmpty(utils.PathSearch("description", v, ""))
		}
		result = append(result, param)
	}

	return result, nil
}

func buildPolicyConditions(conditions *schema.Set) []map[string]interface{} {
	if conditions.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, conditions.Len())
	for _, v := range conditions.List() {
		source := utils.PathSearch("source", v, "param").(string)

		condition := map[string]interface{}{
			"condition_value":                utils.PathSearch("value", v, ""),
			"req_param_name":                 utils.ValueIgnoreEmpty(utils.PathSearch("param_name", v, "")),
			"sys_param_name":                 utils.ValueIgnoreEmpty(utils.PathSearch("sys_name", v, "")),
			"cookie_param_name":              utils.ValueIgnoreEmpty(utils.PathSearch("cookie_name", v, "")),
			"frontend_authorizer_param_name": utils.ValueIgnoreEmpty(utils.PathSearch("frontend_authorizer_name", v, "")),
			"condition_origin":               source,
			"mapped_param_name":              utils.ValueIgnoreEmpty(utils.PathSearch("mapped_param_name", v, "")),
			"mapped_param_location":          utils.ValueIgnoreEmpty(utils.PathSearch("mapped_param_location", v, "")),
		}

		conType := utils.PathSearch("type", v, string(ConditionTypeEqual)).(string)
		// If the input of the condition type is invalid, keep the condition parameter omitted and the API will throw an
		// error.
		if vt, ok := conditionType[conType]; ok {
			condition["condition_type"] = vt
		}

		result = append(result, condition)
	}
	return result
}

func buildMockPolicy(policies []interface{}) ([]map[string]interface{}, error) {
	if len(policies) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		params, err := buildBackendParameters(utils.PathSearch("backend_params", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set))
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"authorizer_id":  utils.ValueIgnoreEmpty(utils.PathSearch("authorizer_id", policy, nil)),
			"name":           utils.PathSearch("name", policy, nil),
			"status_code":    utils.PathSearch("status_code", policy, nil),
			"result_content": utils.PathSearch("response", policy, nil),
			"effect_mode":    utils.PathSearch("effective_mode", policy, nil),
			"conditions":     buildPolicyConditions(utils.PathSearch("conditions", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"backend_params": params,
		})
	}
	return result, nil
}

func buildInvocationType(invocationType, invocationMode string) string {
	if invocationMode != "" {
		return invocationMode
	}

	return invocationType
}

func buildFuncGraphPolicy(policies []interface{}) ([]map[string]interface{}, error) {
	if len(policies) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		params, err := buildBackendParameters(utils.PathSearch("backend_params", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set))
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"authorizer_id": utils.ValueIgnoreEmpty(utils.PathSearch("authorizer_id", policy, nil)),
			"name":          utils.PathSearch("name", policy, nil),
			"function_urn":  utils.PathSearch("function_urn", policy, nil),
			"alias_urn":     utils.PathSearch("function_alias_urn", policy, nil),
			"invocation_type": buildInvocationType(utils.PathSearch("invocation_type", policy, "").(string),
				utils.PathSearch("invocation_mode", policy, "").(string)),
			"effect_mode":    utils.PathSearch("effective_mode", policy, nil),
			"network_type":   utils.PathSearch("network_type", policy, nil),
			"req_protocol":   utils.PathSearch("request_protocol", policy, nil),
			"timeout":        utils.PathSearch("timeout", policy, nil),
			"version":        utils.PathSearch("version", policy, nil),
			"conditions":     buildPolicyConditions(utils.PathSearch("conditions", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"backend_params": params,
		})
	}
	return result, nil
}

func buildApigAPIWebPolicy(policies []interface{}) ([]map[string]interface{}, error) {
	if len(policies) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		params, err := buildBackendParameters(utils.PathSearch("backend_params", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set))
		if err != nil {
			return nil, err
		}
		wp := map[string]interface{}{
			"authorizer_id":  utils.ValueIgnoreEmpty(utils.PathSearch("authorizer_id", policy, nil)),
			"name":           utils.PathSearch("name", policy, nil),
			"req_protocol":   utils.PathSearch("request_protocol", policy, nil),
			"req_method":     utils.PathSearch("request_method", policy, nil),
			"req_uri":        utils.PathSearch("path", policy, nil),
			"effect_mode":    utils.PathSearch("effective_mode", policy, string(EffectiveModeAny)),
			"retry_count":    strconv.Itoa(utils.PathSearch("retry_count", policy, 0).(int)),
			"timeout":        utils.PathSearch("timeout", policy, 0),
			"url_domain":     utils.PathSearch("host_header", policy, nil),
			"conditions":     buildPolicyConditions(utils.PathSearch("conditions", policy, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			"backend_params": params,
		}
		if chanId := utils.PathSearch("vpc_channel_id", policy, ""); chanId != "" {
			wp["vpc_channel_info"] = map[string]interface{}{
				"vpc_channel_id":         utils.PathSearch("vpc_channel_id", policy, ""),
				"vpc_channel_proxy_host": utils.PathSearch("host_header", policy, ""),
			}
			wp["vpc_channel_status"] = vpcChannelEnabled
		} else {
			wp["vpc_channel_status"] = vpcChannelDisabled
		}
		result = append(result, wp)
	}
	return result, nil
}

func buildApiBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	var (
		authType  = d.Get("security_authentication").(string)
		matchMode = d.Get("matching").(string)
	)

	parsedMatchMode, ok := matching[matchMode]
	if !ok {
		return nil, fmt.Errorf("invalid match mode: '%s'", matchMode)
	}

	result := map[string]interface{}{
		"type":                   buildApiType(d.Get("type").(string)),
		"authorizer_id":          utils.ValueIgnoreEmpty(d.Get("authorizer_id")),
		"group_id":               d.Get("group_id"),
		"name":                   d.Get("name"),
		"req_protocol":           d.Get("request_protocol"),
		"req_method":             d.Get("request_method"),
		"req_uri":                d.Get("request_path"),
		"cors":                   d.Get("cors"),
		"auth_type":              authType,
		"match_mode":             parsedMatchMode,
		"remark":                 utils.ValueIgnoreEmpty(d.Get("description")),
		"body_remark":            utils.ValueIgnoreEmpty(d.Get("body_description")),
		"result_normal_sample":   utils.ValueIgnoreEmpty(d.Get("success_response")),
		"result_failure_sample":  utils.ValueIgnoreEmpty(d.Get("failure_response")),
		"response_id":            utils.ValueIgnoreEmpty(d.Get("response_id")),
		"req_params":             buildRequestParameters(d.Get("request_params").([]interface{})),
		"tags":                   utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(d.Get("tags").(*schema.Set))),
		"content_type":           utils.ValueIgnoreEmpty(d.Get("content_type")),
		"is_send_fg_body_base64": d.Get("is_send_fg_body_base64"),
	}

	isSimpleAuthEnabled := d.Get("simple_authentication").(bool)
	if authType == string(ApiAuthTypeApp) {
		if isSimpleAuthEnabled {
			result["auth_opt"] = map[string]interface{}{
				"app_code_auth_type": string(AppCodeAuthTypeEnable),
			}
		} else {
			result["auth_opt"] = map[string]interface{}{
				"app_code_auth_type": string(AppCodeAuthTypeDisable),
			}
		}
	} else if isSimpleAuthEnabled {
		return nil, errors.New("the security authentication must be 'APP' if simple authentication is true")
	}

	// build backend (one of the mock, function graph and web) server and related policies.
	if m, ok := d.GetOk("mock"); ok {
		result["backend_type"] = string(BackendTypeMock)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return nil, err
		}
		result["backend_params"] = params
		result["mock_info"] = buildMockStructure(m.([]interface{}))
		policy, err := buildMockPolicy(d.Get("mock_policy").([]interface{}))
		if err != nil {
			return nil, err
		}
		result["policy_mocks"] = policy
	} else if fg, ok := d.GetOk("func_graph"); ok {
		result["backend_type"] = string(BackendTypeFunction)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return nil, err
		}
		result["backend_params"] = params
		result["func_info"] = buildFuncGraphStructure(fg.([]interface{}))
		policy, err := buildFuncGraphPolicy(d.Get("func_graph_policy").([]interface{}))
		if err != nil {
			return nil, err
		}
		result["policy_functions"] = policy
	} else {
		result["backend_type"] = string(BackendTypeHttp)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return nil, err
		}
		result["backend_params"] = params
		result["backend_api"] = buildWebStructure(d.Get("web").([]interface{}))
		policy, err := buildApigAPIWebPolicy(d.Get("web_policy").([]interface{}))
		if err != nil {
			return nil, err
		}
		result["policy_https"] = policy
	}

	if d.Get("sampling_strategy").(string) != "" && d.Get("sampling_param").(string) != "" {
		result["trace_enabled"] = true
		result["sampling_strategy"] = d.Get("sampling_strategy").(string)
		result["sampling_param"] = d.Get("sampling_param").(string)
	}

	log.Printf("[DEBUG] The API body params is : %+v", result)
	return result, nil
}

func createApi(client *golangsdk.ServiceClient, instanceId string, body map[string]interface{}) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(body),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func buildSliceParamOrderByElementName(requestParams []interface{}) []interface{} {
	if len(requestParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(requestParams))
	for _, v := range requestParams {
		name := utils.PathSearch("name", v, "").(string)
		if name != "" {
			result = append(result, map[string]interface{}{
				"name": name,
			})
		}
	}
	return result
}

func updateAllOriginParameters(d *schema.ResourceData) error {
	var (
		rawConfig            = d.GetRawConfig()
		requestParamsOrigin  = buildSliceParamOrderByElementName(utils.GetNestedObjectFromRawConfig(rawConfig, "request_params").([]interface{}))
		funcGraphPolicyOrder = buildSliceParamOrderByElementName(utils.GetNestedObjectFromRawConfig(rawConfig, "func_graph_policy").([]interface{}))
		webPolicyOrder       = buildSliceParamOrderByElementName(utils.GetNestedObjectFromRawConfig(rawConfig, "web_policy").([]interface{}))
		mockPolicyOrder      = buildSliceParamOrderByElementName(utils.GetNestedObjectFromRawConfig(rawConfig, "mock_policy").([]interface{}))
	)

	if len(requestParamsOrigin) > 0 {
		d.Set("request_params_order", requestParamsOrigin)
	}
	if len(funcGraphPolicyOrder) > 0 {
		d.Set("func_graph_policy_order", funcGraphPolicyOrder)
	}
	if len(webPolicyOrder) > 0 {
		d.Set("web_policy_order", webPolicyOrder)
	}
	if len(mockPolicyOrder) > 0 {
		d.Set("mock_policy_order", mockPolicyOrder)
	}
	return nil
}

func resourceApiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	body, err := buildApiBodyParams(d)
	if err != nil {
		return diag.Errorf("unable to build the API create opts: %s", err)
	}

	respBody, err := createApi(client, instanceId, body)
	if err != nil {
		return diag.Errorf("error creating API: %s", err)
	}
	resourceId := utils.PathSearch("id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the API ID from the API response")
	}
	d.SetId(resourceId)

	if err = updateAllOriginParameters(d); err != nil {
		return diag.Errorf("error updating all origin parameters: %s", err)
	}

	return resourceApiRead(ctx, d, meta)
}

func GetApiById(client *golangsdk.ServiceClient, instanceId, apiId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{api_id}", apiId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func analyseBackendParameterValue(origin, value string) (paramType, paramValue string) {
	log.Printf("[ERROR] The value of the backend parameter is: %s", value)
	if origin == string(ParameterTypeSystem) {
		// Backend parameter types include internal parameters and authorizer parameters, and the authorizer parameter
		// types include front-end parameters and backend parameters.
		regex := regexp.MustCompile(`\$context\.authorizer\.(frontend|backend)\.([\w-]+)`)
		result := regex.FindStringSubmatch(value)
		if len(result) == 3 {
			paramType = result[1]
			paramValue = result[2]
			return
		}

		regex = regexp.MustCompile(`\$context\.([\w-]+)`)
		result = regex.FindStringSubmatch(value)
		if len(result) == 2 {
			paramType = string(SystemParamTypeInternal)
			paramValue = result[1]
			return
		}
		log.Printf("[ERROR] The system parameter format is invalid, want '$context.xxx' (internal parameter), "+
			"'$context.authorizer.frontend.xxx' or '$context.authorizer.frontend.xxx', but '%s'.", value)
		return
	}
	// custom backend parameter
	paramValue = value
	return
}

func flattenBackendParameters(backendParams []interface{}) []map[string]interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(backendParams))
	for _, item := range backendParams {
		origin := utils.PathSearch("origin", item, "").(string)
		value := utils.PathSearch("value", item, "").(string)
		paramAuthType, paramValue := analyseBackendParameterValue(origin, value)
		param := map[string]interface{}{
			"type":     origin,
			"name":     utils.PathSearch("name", item, nil),
			"location": utils.PathSearch("location", item, nil),
			"value":    paramValue,
		}
		if paramAuthType != "" {
			param["system_param_type"] = paramAuthType
		}
		if origin != string(ParameterTypeRequest) {
			param["description"] = utils.PathSearch("remark", item, nil)
		}
		result = append(result, param)
	}
	return result
}

func analyseConditionType(conType string) string {
	for k, v := range conditionType {
		if v == conType {
			return k
		}
	}
	return ""
}

func parseApiType(apiType int) string {
	return map[int]string{
		1: "Public",
		2: "Private",
	}[apiType]
}

func analyseApiMatchMode(mode string) string {
	for k, v := range matching {
		if v == mode {
			return k
		}
	}
	return ""
}

func analyseAppSimpleAuth(authOpt interface{}) bool {
	// HEADER: AppCode authentication is enabled and the AppCode is located in the header.
	return utils.PathSearch("app_code_auth_type", authOpt, "").(string) == string(AppCodeAuthTypeEnable)
}

func parseObjectEnabled(objStatus interface{}) bool {
	var status int
	switch v := objStatus.(type) {
	case float64:
		status = int(v)
	case int:
		status = v
	default:
		return false
	}
	if status == vpcChannelEnabled {
		return true
	}
	if status != vpcChannelDisabled {
		log.Printf("[DEBUG] unexpected object value, want '1'(yes) or '2'(no), but got '%d'", status)
	}
	return false
}

func orderRequestParamsByRequestParamsOrder(reqParams []interface{}, requestParamsOrigin []interface{}) []interface{} {
	if len(requestParamsOrigin) < 1 {
		return reqParams
	}

	sortedReqParams := make([]interface{}, 0, len(reqParams))
	requestParamsCopy := reqParams

	for _, requestParamOrigin := range requestParamsOrigin {
		nameOrigin := utils.PathSearch("name", requestParamOrigin, "").(string)
		for index, requestParam := range requestParamsCopy {
			if utils.PathSearch("name", requestParam, "").(string) != nameOrigin {
				continue
			}
			// Add the found request parameter to the sorted request parameters list.
			sortedReqParams = append(sortedReqParams, requestParamsCopy[index])
			// Remove the processed request parameter from the original request parameters array.
			requestParamsCopy = append(requestParamsCopy[:index], requestParamsCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted request parameters to the end of the sorted list.
	sortedReqParams = append(sortedReqParams, requestParamsCopy...)
	return sortedReqParams
}

func flattenApiRequestParams(reqParams []interface{}, requestParamsOrder []interface{}) []map[string]interface{} {
	if len(reqParams) < 1 {
		return nil
	}

	if len(requestParamsOrder) > 0 {
		reqParams = orderRequestParamsByRequestParamsOrder(reqParams, requestParamsOrder)
	}

	result := make([]map[string]interface{}, 0, len(reqParams))
	for _, v := range reqParams {
		paramType := utils.PathSearch("type", v, "").(string)
		param := map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"location":       utils.PathSearch("location", v, nil),
			"type":           paramType,
			"required":       parseObjectEnabled(utils.PathSearch("required", v, nil)),
			"passthrough":    parseObjectEnabled(utils.PathSearch("pass_through", v, nil)),
			"enumeration":    utils.PathSearch("enumerations", v, nil),
			"example":        utils.PathSearch("sample_value", v, nil),
			"default":        utils.PathSearch("default_value", v, nil),
			"description":    utils.PathSearch("remark", v, nil),
			"valid_enable":   utils.PathSearch("valid_enable", v, nil),
			"orchestrations": utils.PathSearch("orchestrations", v, nil),
		}
		switch paramType {
		case string(ParamTypeNumber):
			param["maximum"] = utils.PathSearch("max_num", v, nil)
			param["minimum"] = utils.PathSearch("min_num", v, nil)
		case string(ParamTypeString):
			param["maximum"] = utils.PathSearch("max_size", v, nil)
			param["minimum"] = utils.PathSearch("min_size", v, nil)
		}
		result = append(result, param)
	}
	return result
}

func flattenMockStructure(mockResp interface{}) []map[string]interface{} {
	if mockResp == nil {
		return nil
	}
	if utils.PathSearch("status_code", mockResp, nil) == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"status_code":   utils.PathSearch("status_code", mockResp, nil),
			"response":      utils.PathSearch("result_content", mockResp, nil),
			"authorizer_id": utils.PathSearch("authorizer_id", mockResp, nil),
		},
	}
}

func flattenFuncGraphStructure(funcResp interface{}) []map[string]interface{} {
	if funcResp == nil {
		return nil
	}
	if utils.PathSearch("function_urn", funcResp, "") == "" {
		return nil
	}

	return []map[string]interface{}{
		{
			"function_urn":       utils.PathSearch("function_urn", funcResp, nil),
			"function_alias_urn": utils.PathSearch("alias_urn", funcResp, nil),
			"timeout":            utils.PathSearch("timeout", funcResp, nil),
			"invocation_type":    utils.PathSearch("invocation_type", funcResp, nil),
			"network_type":       utils.PathSearch("network_type", funcResp, nil),
			"request_protocol":   utils.PathSearch("req_protocol", funcResp, nil),
			"version":            utils.PathSearch("version", funcResp, nil),
			"authorizer_id":      utils.PathSearch("authorizer_id", funcResp, nil),
		},
	}
}

func flattenWebStructure(webResp interface{}, sslEnabled bool) []map[string]interface{} {
	if webResp == nil {
		return nil
	}
	if utils.PathSearch("req_uri", webResp, "") == "" {
		return nil
	}

	result := map[string]interface{}{
		"path":             utils.PathSearch("req_uri", webResp, nil),
		"request_method":   utils.PathSearch("req_method", webResp, nil),
		"request_protocol": utils.PathSearch("req_protocol", webResp, nil),
		"timeout":          utils.PathSearch("timeout", webResp, nil),
		"ssl_enable":       sslEnabled,
		"authorizer_id":    utils.PathSearch("authorizer_id", webResp, nil),
	}
	retryCount := utils.PathSearch("retry_count", webResp, "").(string)
	if retryCount != "" {
		result["retry_count"] = utils.StringToInt(&retryCount)
	}
	vpcChannelId := utils.PathSearch("vpc_channel_info.vpc_channel_id", webResp, "")
	if vpcChannelId != "" {
		result["vpc_channel_id"] = vpcChannelId
		result["host_header"] = utils.PathSearch("vpc_channel_info.vpc_channel_proxy_host", webResp, nil)
	} else {
		result["backend_address"] = utils.PathSearch("url_domain", webResp, nil)
	}

	return []map[string]interface{}{
		result,
	}
}

func flattenPolicyConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, v := range conditions {
		result = append(result, map[string]interface{}{
			"source":                   utils.PathSearch("condition_origin", v, nil),
			"param_name":               utils.PathSearch("req_param_name", v, nil),
			"sys_name":                 utils.PathSearch("sys_param_name", v, nil),
			"cookie_name":              utils.PathSearch("cookie_param_name", v, nil),
			"frontend_authorizer_name": utils.PathSearch("frontend_authorizer_param_name", v, nil),
			"type":                     analyseConditionType(utils.PathSearch("condition_type", v, "").(string)),
			"value":                    utils.PathSearch("condition_value", v, nil),
			"mapped_param_name":        utils.PathSearch("mapped_param_name", v, nil),
			"mapped_param_location":    utils.PathSearch("mapped_param_location", v, nil),
		})
	}
	return result
}

func orderFuncGraphPolicyByFuncGraphPolicyOrder(policies []interface{}, funcGraphPolicyOrigin []interface{}) []interface{} {
	if len(funcGraphPolicyOrigin) < 1 {
		return policies
	}

	sortedPolicies := make([]interface{}, 0, len(policies))
	funcGraphPolicyCopy := policies

	for _, policyOrigin := range funcGraphPolicyOrigin {
		nameOrigin := utils.PathSearch("name", policyOrigin, "").(string)
		for index, funcGraphPolicy := range funcGraphPolicyCopy {
			if utils.PathSearch("name", funcGraphPolicy, "").(string) != nameOrigin {
				continue
			}
			// Add the found func graph policy to the sorted func graph policies list.
			sortedPolicies = append(sortedPolicies, funcGraphPolicyCopy[index])
			// Remove the processed func graph policy from the original func graph policies array.
			funcGraphPolicyCopy = append(funcGraphPolicyCopy[:index], funcGraphPolicyCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted func graph policies to the end of the sorted list.
	sortedPolicies = append(sortedPolicies, funcGraphPolicyCopy...)
	return sortedPolicies
}

func flattenFuncGraphPolicy(policies []interface{}, funcGraphPolicyOrder []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	if len(funcGraphPolicyOrder) > 0 {
		policies = orderFuncGraphPolicyByFuncGraphPolicyOrder(policies, funcGraphPolicyOrder)
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", policy, nil),
			"function_urn":       utils.PathSearch("function_urn", policy, nil),
			"function_alias_urn": utils.PathSearch("alias_urn", policy, nil),
			"version":            utils.PathSearch("version", policy, nil),
			"network_type":       utils.PathSearch("network_type", policy, nil),
			"request_protocol":   utils.PathSearch("req_protocol", policy, nil),
			"invocation_type":    utils.PathSearch("invocation_type", policy, nil),
			"effective_mode":     utils.PathSearch("effect_mode", policy, nil),
			"timeout":            utils.PathSearch("timeout", policy, nil),
			"authorizer_id":      utils.PathSearch("authorizer_id", policy, nil),
			"backend_params": flattenBackendParameters(utils.PathSearch("backend_params", policy,
				make([]interface{}, 0)).([]interface{})),
			"conditions": flattenPolicyConditions(utils.PathSearch("conditions", policy,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func orderWebPolicyByWebPolicyOrder(policies []interface{}, webPolicyOrigin []interface{}) []interface{} {
	if len(webPolicyOrigin) < 1 {
		return policies
	}

	sortedPolicies := make([]interface{}, 0, len(policies))
	webPolicyCopy := policies

	for _, policyOrigin := range webPolicyOrigin {
		nameOrigin := utils.PathSearch("name", policyOrigin, "").(string)
		for index, webPolicy := range webPolicyCopy {
			if utils.PathSearch("name", webPolicy, "").(string) != nameOrigin {
				continue
			}
			sortedPolicies = append(sortedPolicies, webPolicyCopy[index])
			webPolicyCopy = append(webPolicyCopy[:index], webPolicyCopy[index+1:]...)
		}
	}
	sortedPolicies = append(sortedPolicies, webPolicyCopy...)
	return sortedPolicies
}

func flattenWebPolicy(policies []interface{}, webPolicyOrder []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	if len(webPolicyOrder) > 0 {
		policies = orderWebPolicyByWebPolicyOrder(policies, webPolicyOrder)
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		wp := map[string]interface{}{
			"name":             utils.PathSearch("name", policy, nil),
			"request_protocol": utils.PathSearch("req_protocol", policy, nil),
			"request_method":   utils.PathSearch("req_method", policy, nil),
			"effective_mode":   utils.PathSearch("effect_mode", policy, nil),
			"path":             utils.PathSearch("req_uri", policy, nil),
			"timeout":          utils.PathSearch("timeout", policy, nil),
			"authorizer_id":    utils.PathSearch("authorizer_id", policy, nil),
			"backend_params": flattenBackendParameters(utils.PathSearch("backend_params", policy,
				make([]interface{}, 0)).([]interface{})),
			"conditions": flattenPolicyConditions(utils.PathSearch("conditions", policy,
				make([]interface{}, 0)).([]interface{})),
		}
		retryCount := utils.PathSearch("retry_count", policy, "").(string)
		if retryCount != "" {
			wp["retry_count"] = utils.StringToInt(&retryCount)
		}
		// which policy use backend address or vpc channel.
		vpcChannelId := utils.PathSearch("vpc_channel_info.vpc_channel_id", policy, "")
		if vpcChannelId != "" {
			wp["vpc_channel_id"] = vpcChannelId
			wp["host_header"] = utils.PathSearch("vpc_channel_info.vpc_channel_proxy_host", policy, nil)
		} else {
			wp["backend_address"] = utils.PathSearch("url_domain", policy, nil)
		}

		result = append(result, wp)
	}

	return result
}

func orderMockPolicyByMockPolicyOrder(policies []interface{}, mockPolicyOrigin []interface{}) []interface{} {
	if len(mockPolicyOrigin) < 1 {
		return policies
	}

	sortedPolicies := make([]interface{}, 0, len(policies))
	mockPolicyCopy := policies

	for _, policyOrigin := range mockPolicyOrigin {
		nameOrigin := utils.PathSearch("name", policyOrigin, "").(string)
		for index, mockPolicy := range mockPolicyCopy {
			if utils.PathSearch("name", mockPolicy, "").(string) != nameOrigin {
				continue
			}
			// Add the found mock policy to the sorted mock policies list.
			sortedPolicies = append(sortedPolicies, mockPolicyCopy[index])
			// Remove the processed mock policy from the original mock policies array.
			mockPolicyCopy = append(mockPolicyCopy[:index], mockPolicyCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted mock policies to the end of the sorted list.
	sortedPolicies = append(sortedPolicies, mockPolicyCopy...)
	return sortedPolicies
}

func flattenMockPolicy(policies []interface{}, mockPolicyOrigin []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	if len(mockPolicyOrigin) > 0 {
		policies = orderMockPolicyByMockPolicyOrder(policies, mockPolicyOrigin)
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"name":           utils.PathSearch("name", policy, nil),
			"status_code":    utils.PathSearch("status_code", policy, nil),
			"response":       utils.PathSearch("result_content", policy, nil),
			"effective_mode": utils.PathSearch("effect_mode", policy, nil),
			"authorizer_id":  utils.PathSearch("authorizer_id", policy, nil),
			"backend_params": flattenBackendParameters(utils.PathSearch("backend_params", policy,
				make([]interface{}, 0)).([]interface{})),
			"conditions": flattenPolicyConditions(utils.PathSearch("conditions", policy,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func resourceApiRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		instanceId            = d.Get("instance_id").(string)
		apiId                 = d.Id()
		requestParamsOrigin   = d.Get("request_params_order").([]interface{})
		funcGraphPolicyOrigin = d.Get("func_graph_policy_order").([]interface{})
		webPolicyOrigin       = d.Get("web_policy_order").([]interface{})
		mockPolicyOrigin      = d.Get("mock_policy_order").([]interface{})
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := GetApiById(client, instanceId, apiId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying API (%s)", apiId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("group_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("authorizer_id", utils.PathSearch("authorizer_id", respBody, nil)),
		d.Set("tags", utils.PathSearch("tags", respBody, nil)),
		d.Set("content_type", utils.PathSearch("content_type", respBody, nil)),
		d.Set("is_send_fg_body_base64", utils.PathSearch("is_send_fg_body_base64", respBody, nil)),
		d.Set("request_protocol", utils.PathSearch("req_protocol", respBody, nil)),
		d.Set("request_method", utils.PathSearch("req_method", respBody, nil)),
		d.Set("request_path", utils.PathSearch("req_uri", respBody, nil)),
		d.Set("security_authentication", utils.PathSearch("auth_type", respBody, nil)),
		d.Set("cors", utils.PathSearch("cors", respBody, nil)),
		d.Set("sampling_strategy", utils.PathSearch("sampling_strategy", respBody, "").(string)),
		d.Set("sampling_param", utils.PathSearch("sampling_param", respBody, "").(string)),
		d.Set("description", utils.PathSearch("remark", respBody, nil)),
		d.Set("body_description", utils.PathSearch("body_remark", respBody, nil)),
		d.Set("success_response", utils.PathSearch("result_normal_sample", respBody, nil)),
		d.Set("failure_response", utils.PathSearch("result_failure_sample", respBody, nil)),
		d.Set("response_id", utils.PathSearch("response_id", respBody, nil)),
		d.Set("type", parseApiType(int(utils.PathSearch("type", respBody, float64(0)).(float64)))),
		d.Set("request_params", flattenApiRequestParams(utils.PathSearch("req_params", respBody,
			make([]interface{}, 0)).([]interface{}), requestParamsOrigin)),
		d.Set("backend_params", flattenBackendParameters(utils.PathSearch("backend_params", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("matching", analyseApiMatchMode(utils.PathSearch("match_mode", respBody, "").(string))),
		d.Set("simple_authentication", analyseAppSimpleAuth(utils.PathSearch("auth_opt", respBody, nil))),
		d.Set("mock", flattenMockStructure(utils.PathSearch("mock_info", respBody, nil))),
		d.Set("func_graph", flattenFuncGraphStructure(utils.PathSearch("func_info", respBody, nil))),
		d.Set("func_graph_policy", flattenFuncGraphPolicy(utils.PathSearch("policy_functions", respBody,
			make([]interface{}, 0)).([]interface{}), funcGraphPolicyOrigin)),
		d.Set("web", flattenWebStructure(utils.PathSearch("backend_api", respBody, nil), d.Get("web.0.ssl_enable").(bool))),
		d.Set("web_policy", flattenWebPolicy(utils.PathSearch("policy_https", respBody,
			make([]interface{}, 0)).([]interface{}), webPolicyOrigin)),
		d.Set("mock_policy", flattenMockPolicy(utils.PathSearch("policy_mocks", respBody,
			make([]interface{}, 0)).([]interface{}), mockPolicyOrigin)),
		d.Set("registered_at", utils.PathSearch("register_time", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", respBody, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving API fields: %s", err)
	}
	return nil
}

func updateApi(client *golangsdk.ServiceClient, instanceId, apiId string, body map[string]interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{api_id}", apiId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(body),
		OkCodes:  []int{200},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceApiUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	body, err := buildApiBodyParams(d)
	if err != nil {
		return diag.Errorf("unable to build the API updateOpts: %s", err)
	}

	if err = updateApi(client, instanceId, apiId, body); err != nil {
		return diag.Errorf("error updating API (%s): %s", apiId, err)
	}

	if err = updateAllOriginParameters(d); err != nil {
		return diag.Errorf("error updating all origin parameters: %s", err)
	}

	return resourceApiRead(ctx, d, meta)
}

func deleteApi(client *golangsdk.ServiceClient, instanceId, apiId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}"
	unpublishPath := client.Endpoint + httpUrl
	unpublishPath = strings.ReplaceAll(unpublishPath, "{project_id}", client.ProjectID)
	unpublishPath = strings.ReplaceAll(unpublishPath, "{instance_id}", instanceId)
	unpublishPath = strings.ReplaceAll(unpublishPath, "{api_id}", apiId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", unpublishPath, &opt)
	return err
}

func resourceApiDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	if err = deleteApi(client, instanceId, apiId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting API (%s)", apiId))
	}

	return nil
}

// the value of queryParams must start with '&' character.
func listApis(client *golangsdk.ServiceClient, instanceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis?limit={limit}"
		offset  = 0
		limit   = 500
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	if len(queryParams) > 0 {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apisList := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, apisList...)
		if len(apisList) < limit {
			break
		}
		offset += len(apisList)
	}

	return result, nil
}

func getApiIdByName(client *golangsdk.ServiceClient, instanceId, name string) (string, error) {
	queryParams := fmt.Sprintf("&name=%s", name)
	apiRecords, err := listApis(client, instanceId, queryParams)
	if err != nil {
		return "", fmt.Errorf("error retrieving APIs: %s", err)
	}

	apiId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", name), apiRecords, "").(string)
	if apiId == "" {
		return "", fmt.Errorf("unable to find the API (%s) form APIG service", name)
	}
	return apiId, nil
}

func resourceApiImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	instanceId := parts[0]
	apiId, err := getApiIdByName(client, instanceId, parts[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(apiId)

	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
