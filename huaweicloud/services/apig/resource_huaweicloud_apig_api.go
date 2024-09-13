package apig

import (
	"bytes"
	"context"
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
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

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

	ConditionSourceParam              ConditionSource = "param"
	ConditionSourceSource             ConditionSource = "source"
	ConditionSourceSystem             ConditionSource = "system"
	ConditionSourceCookie             ConditionSource = "cookie"
	ConditionSourceFrontendAuthorizer ConditionSource = "frontend_authorizer"

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

	strBoolEnabled  int = 1
	strBoolDisabled int = 2
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

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the API is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance to which the API belongs.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the API group to which the API belongs.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ApiTypePublic),
					string(ApiTypePrivate),
				}, false),
				Description: "The API type.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The API name.",
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RequestMethodGet),
					string(RequestMethodPost),
					string(RequestMethodPut),
					string(RequestMethodDelete),
					string(RequestMethodHead),
					string(RequestMethodPatch),
					string(RequestMethodOptions),
					string(RequestMethodAny),
				}, false),
				Description: "The request method of the API.",
			},
			"request_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The request address.",
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ProtocolTypeHTTP),
					string(ProtocolTypeHTTPS),
					string(ProtocolTypeBoth),
					string(ProtocolTypeGPRCS),
				}, false),
				Description: "The request protocol of the API request.",
			},
			"security_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ApiAuthTypeNone),
				ValidateFunc: validation.StringInSlice([]string{
					string(ApiAuthTypeNone),
					string(ApiAuthTypeApp),
					string(ApiAuthTypeIam),
					string(ApiAuthTypeAuthorizer),
				}, false),
				Description: "The security authentication mode of the API request.",
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
			"request_params": {
				Type:     schema.TypeSet,
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
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(ParamLocationPath),
							ValidateFunc: validation.StringInSlice([]string{
								string(ParamLocationPath),
								string(ParamLocationHeader),
								string(ParamLocationQuery),
							}, false),
							Description: "Where this parameter is located.",
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(ParamTypeString),
							ValidateFunc: validation.StringInSlice([]string{
								string(ParamTypeString),
								string(ParamTypeNumber),
							}, false),
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API description.",
			},
			"matching": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  MatchModeExact,
				ValidateFunc: validation.StringInSlice([]string{
					string(MatchModePrefix),
					string(MatchModeExact),
				}, false),
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
				ForceNew:     true,
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
				ForceNew: true,
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
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(InvacationTypeSync),
							ValidateFunc: validation.StringInSlice([]string{
								string(InvacationTypeAsync),
								string(InvacationTypeSync),
							}, false),
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
				ForceNew: true,
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(RequestMethodGet),
								string(RequestMethodPost),
								string(RequestMethodPut),
								string(RequestMethodDelete),
								string(RequestMethodHead),
								string(RequestMethodPatch),
								string(RequestMethodOptions),
								string(RequestMethodAny),
							}, false),
							Description: "The backend request method of the API.",
						},
						"request_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(ProtocolTypeHTTPS),
							ValidateFunc: validation.StringInSlice([]string{
								string(ProtocolTypeHTTP),
								string(ProtocolTypeHTTPS),
							}, false),
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
				Type:          schema.TypeSet,
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
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(EffectiveModeAny),
							ValidateFunc: validation.StringInSlice([]string{
								string(EffectiveModeAll),
								string(EffectiveModeAny),
							}, false),
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
				Type:          schema.TypeSet,
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
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(InvacationTypeSync),
							ValidateFunc: validation.StringInSlice([]string{
								string(InvacationTypeAsync),
								string(InvacationTypeSync),
							}, false),
							Description: "The invocation mode of the FunctionGraph function.",
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(EffectiveModeAny),
							ValidateFunc: validation.StringInSlice([]string{
								string(EffectiveModeAll),
								string(EffectiveModeAny),
							}, false),
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
							ValidateFunc: validation.StringInSlice([]string{
								string(InvacationTypeAsync),
								string(InvacationTypeSync),
							}, false),
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
				Type:          schema.TypeSet,
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(RequestMethodGet),
								string(RequestMethodPost),
								string(RequestMethodPut),
								string(RequestMethodDelete),
								string(RequestMethodHead),
								string(RequestMethodPatch),
								string(RequestMethodOptions),
								string(RequestMethodAny),
							}, false),
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
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(ProtocolTypeHTTP),
								string(ProtocolTypeHTTPS),
							}, false),
							Description: "The backend request protocol.",
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(EffectiveModeAny),
							ValidateFunc: validation.StringInSlice([]string{
								string(EffectiveModeAll),
								string(EffectiveModeAny),
							}, false),
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
			"registered_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registered time of the API.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the API.",
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ConditionSourceParam),
				ValidateFunc: validation.StringInSlice([]string{
					string(ConditionSourceParam),
					string(ConditionSourceSource),
					string(ConditionSourceSystem),
					string(ConditionSourceCookie),
					string(ConditionSourceFrontendAuthorizer),
				}, false),
				Description: "The type of the backend policy.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ConditionTypeEqual),
				ValidateFunc: validation.StringInSlice([]string{
					string(ConditionTypeEqual),
					string(ConditionTypeEnumerated),
					string(ConditionTypeMatching),
				}, false),
				Description: "The condition type.",
			},
		},
	}
}

func backendParamSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParameterTypeRequest),
					string(ParameterTypeConstant),
					string(ParameterTypeSystem),
				}, false),
				Description: "The parameter type.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The parameter name.",
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ParamLocationPath),
					string(ParamLocationQuery),
					string(ParamLocationHeader),
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					string(SystemParamTypeInternal),
					string(SystemParamTypeFrontend),
					string(SystemParamTypeBackend),
				}, false),
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
		return strBoolEnabled
	}
	return strBoolDisabled
}

func buildMockStructure(mocks []interface{}) *apis.Mock {
	if len(mocks) < 1 {
		return nil
	}

	mockMap := mocks[0].(map[string]interface{})
	return &apis.Mock{
		StatusCode:    mockMap["status_code"].(int),
		ResultContent: utils.String(mockMap["response"].(string)),
		AuthorizerId:  utils.String(mockMap["authorizer_id"].(string)),
	}
}

func buildFuncGraphStructure(funcGraphs []interface{}) *apis.FuncGraph {
	if len(funcGraphs) < 1 {
		return nil
	}

	funcMap := funcGraphs[0].(map[string]interface{})
	return &apis.FuncGraph{
		FunctionUrn:      funcMap["function_urn"].(string),
		FunctionAliasUrn: funcMap["function_alias_urn"].(string),
		NetworkType:      funcMap["network_type"].(string),
		Timeout:          funcMap["timeout"].(int),
		InvocationType:   funcMap["invocation_type"].(string),
		Version:          funcMap["version"].(string),
		AuthorizerId:     utils.String(funcMap["authorizer_id"].(string)),
		RequestProtocol:  funcMap["request_protocol"].(string),
	}
}

func buildWebStructure(webs []interface{}) *apis.Web {
	if len(webs) < 1 {
		return nil
	}

	var (
		webMap  = webs[0].(map[string]interface{})
		webResp = apis.Web{
			ReqURI:          webMap["path"].(string),
			ReqMethod:       webMap["request_method"].(string),
			ReqProtocol:     webMap["request_protocol"].(string),
			Timeout:         webMap["timeout"].(int),
			ClientSslEnable: utils.Bool(webMap["ssl_enable"].(bool)),
			AuthorizerId:    utils.String(webMap["authorizer_id"].(string)),
			RetryCount:      utils.String(strconv.Itoa(webMap["retry_count"].(int))),
		}
	)
	// If vpc_channel_id is empty, the backend address is used.
	if chanId, ok := webMap["vpc_channel_id"]; ok && chanId != "" {
		webResp.VpcChannelStatus = strBoolEnabled
		webResp.VpcChannelInfo = &apis.VpcChannel{
			VpcChannelId:        chanId.(string),
			VpcChannelProxyHost: webMap["host_header"].(string),
		}
	} else {
		webResp.VpcChannelStatus = strBoolDisabled
		webResp.DomainURL = webMap["backend_address"].(string)
	}

	return &webResp
}

func buildRequestParameters(requests *schema.Set) []apis.ReqParamBase {
	if requests.Len() < 1 {
		return nil
	}

	result := make([]apis.ReqParamBase, requests.Len())
	for i, v := range requests.List() {
		paramMap := v.(map[string]interface{})
		paramType := paramMap["type"].(string)
		param := apis.ReqParamBase{
			Type:         paramType,
			Name:         paramMap["name"].(string),
			Required:     isObjectEnabled(paramMap["required"].(bool)),
			Location:     paramMap["location"].(string),
			Description:  utils.String(paramMap["description"].(string)),
			Enumerations: utils.String(paramMap["enumeration"].(string)),
			PassThrough:  isObjectEnabled(paramMap["passthrough"].(bool)),
			DefaultValue: utils.String(paramMap["default"].(string)),
			SampleValue:  paramMap["example"].(string),
			ValidEnable:  paramMap["valid_enable"].(int),
		}
		switch paramType {
		case string(ParamTypeNumber):
			param.MaxNum = utils.Int(paramMap["maximum"].(int))
			param.MinNum = utils.Int(paramMap["minimum"].(int))
		case string(ParamTypeString):
			param.MaxSize = utils.Int(paramMap["maximum"].(int))
			param.MinSize = utils.Int(paramMap["minimum"].(int))
		}
		result[i] = param
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
func buildBackendParameters(backends *schema.Set) ([]apis.BackendParamBase, error) {
	result := make([]apis.BackendParamBase, backends.Len())
	for i, v := range backends.List() {
		pm := v.(map[string]interface{})
		origin := pm["type"].(string)
		if origin == string(ParameterTypeSystem) && pm["system_param_type"].(string) == "" {
			return nil, fmt.Errorf("The 'system_param_type' must set if parameter type is 'SYSTEM'")
		}
		param := apis.BackendParamBase{
			Origin:   origin,
			Name:     pm["name"].(string),
			Location: pm["location"].(string),
			Value:    buildBackendParameterValue(origin, pm["value"].(string), pm["system_param_type"].(string)),
		}

		if origin != string(ParameterTypeRequest) {
			param.Description = utils.String(pm["description"].(string))
		}
		result[i] = param
	}

	return result, nil
}

func buildPolicyConditions(conditions *schema.Set) []apis.APIConditionBase {
	if conditions.Len() < 1 {
		return nil
	}

	result := make([]apis.APIConditionBase, conditions.Len())
	for i, v := range conditions.List() {
		cm := v.(map[string]interface{})
		condition := apis.APIConditionBase{
			ReqParamName:                cm["param_name"].(string),
			SysParamName:                cm["sys_name"].(string),
			CookieParamName:             cm["cookie_name"].(string),
			FrontendAuthorizerParamName: cm["frontend_authorizer_name"].(string),
			ConditionOrigin:             cm["source"].(string),
			ConditionValue:              cm["value"].(string),
		}
		conType := cm["type"].(string)
		// If the input of the condition type is invalid, keep the condition parameter omitted and the API will throw an
		// error.
		if v, ok := conditionType[conType]; ok {
			condition.ConditionType = v
		}
		result[i] = condition
	}
	return result
}

func buildMockPolicy(policies *schema.Set) ([]apis.PolicyMock, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyMock, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		result[i] = apis.PolicyMock{
			AuthorizerId:  utils.String(pm["authorizer_id"].(string)),
			Name:          pm["name"].(string),
			StatusCode:    pm["status_code"].(int),
			ResultContent: pm["response"].(string),
			EffectMode:    pm["effective_mode"].(string),
			Conditions:    buildPolicyConditions(pm["conditions"].(*schema.Set)),
			BackendParams: params,
		}
	}
	return result, nil
}

func buildInvocationType(invocationType, invocationMode string) string {
	if invocationMode != "" {
		return invocationMode
	}

	return invocationType
}

func buildFuncGraphPolicy(policies *schema.Set) ([]apis.PolicyFuncGraph, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyFuncGraph, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		result[i] = apis.PolicyFuncGraph{
			AuthorizerId:     utils.String(pm["authorizer_id"].(string)),
			Name:             pm["name"].(string),
			FunctionUrn:      pm["function_urn"].(string),
			FunctionAliasUrn: pm["function_alias_urn"].(string),
			InvocationType:   buildInvocationType(pm["invocation_type"].(string), pm["invocation_mode"].(string)),
			EffectMode:       pm["effective_mode"].(string),
			NetworkType:      pm["network_type"].(string),
			RequestProtocol:  pm["request_protocol"].(string),
			Timeout:          pm["timeout"].(int),
			Version:          pm["version"].(string),
			Conditions:       buildPolicyConditions(pm["conditions"].(*schema.Set)),
			BackendParams:    params,
		}
	}
	return result, nil
}

func buildApigAPIWebPolicy(policies *schema.Set) ([]apis.PolicyWeb, error) {
	if policies.Len() < 1 {
		return nil, nil
	}

	result := make([]apis.PolicyWeb, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		params, err := buildBackendParameters(pm["backend_params"].(*schema.Set))
		if err != nil {
			return nil, err
		}
		wp := apis.PolicyWeb{
			AuthorizerId:  utils.String(pm["authorizer_id"].(string)),
			Name:          pm["name"].(string),
			ReqProtocol:   pm["request_protocol"].(string),
			ReqMethod:     pm["request_method"].(string),
			ReqURI:        pm["path"].(string),
			EffectMode:    pm["effective_mode"].(string),
			RetryCount:    utils.String(strconv.Itoa(pm["retry_count"].(int))),
			Timeout:       pm["timeout"].(int),
			DomainURL:     pm["host_header"].(string),
			Conditions:    buildPolicyConditions(pm["conditions"].(*schema.Set)),
			BackendParams: params,
		}
		if chanId, ok := pm["vpc_channel_id"]; ok {
			if chanId != "" {
				wp.VpcChannelInfo = &apis.VpcChannel{
					VpcChannelId:        pm["vpc_channel_id"].(string),
					VpcChannelProxyHost: pm["host_header"].(string),
				}
				wp.VpcChannelStatus = strBoolEnabled
			} else {
				wp.VpcChannelStatus = strBoolDisabled
			}
		}
		result[i] = wp
	}
	return result, nil
}

func buildApiCreateOpts(d *schema.ResourceData) (apis.APIOpts, error) {
	authType := d.Get("security_authentication").(string)
	opt := apis.APIOpts{
		Type:                buildApiType(d.Get("type").(string)),
		AuthorizerId:        d.Get("authorizer_id").(string),
		GroupId:             d.Get("group_id").(string),
		Name:                d.Get("name").(string),
		ReqProtocol:         d.Get("request_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqURI:              d.Get("request_path").(string),
		Cors:                utils.Bool(d.Get("cors").(bool)),
		AuthType:            authType,
		MatchMode:           d.Get("matching").(string),
		Description:         utils.String(d.Get("description").(string)),
		BodyDescription:     utils.String(d.Get("body_description").(string)),
		ResultNormalSample:  utils.String(d.Get("success_response").(string)),
		ResultFailureSample: utils.String(d.Get("failure_response").(string)),
		ResponseId:          d.Get("response_id").(string),
		ReqParams:           buildRequestParameters(d.Get("request_params").(*schema.Set)),
		Tags:                utils.ExpandToStringListBySet(d.Get("tags").((*schema.Set))),
	}
	// build match mode
	matchMode := d.Get("matching").(string)
	v, ok := matching[matchMode]
	if !ok {
		return opt, fmt.Errorf("invalid match mode: '%s'", matchMode)
	}
	opt.MatchMode = v

	isSimpleAuthEnabled := d.Get("simple_authentication").(bool)
	if authType == string(ApiAuthTypeApp) {
		if isSimpleAuthEnabled {
			opt.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: string(AppCodeAuthTypeEnable),
			}
		} else {
			opt.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: string(AppCodeAuthTypeDisable),
			}
		}
	} else if isSimpleAuthEnabled {
		return opt, fmt.Errorf("the security authentication must be 'APP' if simple authentication is true")
	}

	// build backend (one of the mock, function graph and web) server and related policies.
	if m, ok := d.GetOk("mock"); ok {
		opt.BackendType = string(BackendTypeMock)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.BackendParams = params
		opt.MockInfo = buildMockStructure(m.([]interface{}))
		policy, err := buildMockPolicy(d.Get("mock_policy").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.PolicyMocks = policy
	} else if fg, ok := d.GetOk("func_graph"); ok {
		opt.BackendType = string(BackendTypeFunction)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.BackendParams = params
		opt.FuncInfo = buildFuncGraphStructure(fg.([]interface{}))
		policy, err := buildFuncGraphPolicy(d.Get("func_graph_policy").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.PolicyFunctions = policy
	} else {
		opt.BackendType = string(BackendTypeHttp)
		params, err := buildBackendParameters(d.Get("backend_params").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.BackendParams = params
		opt.WebInfo = buildWebStructure(d.Get("web").([]interface{}))
		policy, err := buildApigAPIWebPolicy(d.Get("web_policy").(*schema.Set))
		if err != nil {
			return opt, err
		}
		opt.PolicyWebs = policy
	}

	log.Printf("[DEBUG] The API Opts is : %+v", opt)
	return opt, nil
}

func resourceApiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt, err := buildApiCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to build the API create opts: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apis.Create(client, instanceId, opt).Extract()
	if err != nil {
		return diag.Errorf("error creating API: %s", err)
	}
	d.SetId(resp.ID)

	return resourceApiRead(ctx, d, meta)
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

func flattenBackendParameters(backendParams []apis.BackendParamResp) []map[string]interface{} {
	if len(backendParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(backendParams))
	for i, v := range backendParams {
		origin := v.Origin
		paramAuthType, paramValue := analyseBackendParameterValue(v.Origin, v.Value)
		param := map[string]interface{}{
			"type":     origin,
			"name":     v.Name,
			"location": v.Location,
			"value":    paramValue,
		}
		if paramAuthType != "" {
			param["system_param_type"] = paramAuthType
		}
		if origin != string(ParameterTypeRequest) {
			param["description"] = v.Description
		}
		result[i] = param
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

func analyseApiType(t int) string {
	apiType := map[int]string{
		1: "Public",
		2: "Private",
	}
	if v, ok := apiType[t]; ok {
		return v
	}
	return ""
}

func analyseApiMatchMode(mode string) string {
	for k, v := range matching {
		if v == mode {
			return k
		}
	}
	return ""
}

func analyseAppSimpleAuth(opt apis.AuthOpt) bool {
	// HEADER: AppCode authentication is enabled and the AppCode is located in the header.
	return opt.AppCodeAuthType == string(AppCodeAuthTypeEnable)
}

func parseObjectEnabled(objStatus int) bool {
	if objStatus == strBoolEnabled {
		return true
	}
	if objStatus != strBoolDisabled {
		log.Printf("[DEBUG] unexpected object value, want '1'(yes) or '2'(no), but got '%d'", objStatus)
	}
	return false
}

func flattenApiRequestParams(reqParams []apis.ReqParamResp) []map[string]interface{} {
	if len(reqParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
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

func flattenMockStructure(mockResp apis.Mock) []map[string]interface{} {
	if mockResp == (apis.Mock{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"status_code":   mockResp.StatusCode,
			"response":      mockResp.ResultContent,
			"authorizer_id": mockResp.AuthorizerId,
		},
	}
}

func flattenFuncGraphStructure(funcResp apis.FuncGraph) []map[string]interface{} {
	if funcResp == (apis.FuncGraph{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"function_urn":       funcResp.FunctionUrn,
			"function_alias_urn": funcResp.FunctionAliasUrn,
			"timeout":            funcResp.Timeout,
			"invocation_type":    funcResp.InvocationType,
			"network_type":       funcResp.NetworkType,
			"request_protocol":   funcResp.RequestProtocol,
			"version":            funcResp.Version,
			"authorizer_id":      funcResp.AuthorizerId,
		},
	}
}

func flattenWebStructure(webResp apis.Web, sslEnabled bool) []map[string]interface{} {
	if webResp == (apis.Web{}) {
		return nil
	}

	result := map[string]interface{}{
		"path":             webResp.ReqURI,
		"request_method":   webResp.ReqMethod,
		"request_protocol": webResp.ReqProtocol,
		"timeout":          webResp.Timeout,
		"ssl_enable":       sslEnabled,
		"authorizer_id":    webResp.AuthorizerId,
		"retry_count":      utils.StringToInt(webResp.RetryCount),
	}
	if webResp.VpcChannelInfo.VpcChannelId != "" {
		result["vpc_channel_id"] = webResp.VpcChannelInfo.VpcChannelId
		result["host_header"] = webResp.VpcChannelInfo.VpcChannelProxyHost
	} else {
		result["backend_address"] = webResp.DomainURL
	}

	return []map[string]interface{}{
		result,
	}
}

func flattenPolicyConditions(conditions []apis.APIConditionBase) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		result[i] = map[string]interface{}{
			"source":                   v.ConditionOrigin,
			"param_name":               v.ReqParamName,
			"sys_name":                 v.SysParamName,
			"cookie_name":              v.CookieParamName,
			"frontend_authorizer_name": v.FrontendAuthorizerParamName,
			"type":                     analyseConditionType(v.ConditionType),
			"value":                    v.ConditionValue,
		}
	}
	return result
}

func flattenMockPolicy(policies []apis.PolicyMockResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"name":           policy.Name,
			"status_code":    policy.StatusCode,
			"response":       policy.ResultContent,
			"effective_mode": policy.EffectMode,
			"authorizer_id":  policy.AuthorizerId,
			"backend_params": flattenBackendParameters(policy.BackendParams),
			"conditions":     flattenPolicyConditions(policy.Conditions),
		}
	}

	return result
}

func flattenFuncGraphPolicy(policies []apis.PolicyFuncGraphResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"name":               policy.Name,
			"function_urn":       policy.FunctionUrn,
			"function_alias_urn": policy.FunctionAliasUrn,
			"version":            policy.Version,
			"network_type":       policy.NetworkType,
			"request_protocol":   policy.RequestProtocol,
			"invocation_type":    policy.InvocationType,
			"effective_mode":     policy.EffectMode,
			"timeout":            policy.Timeout,
			"authorizer_id":      policy.AuthorizerId,
			"backend_params":     flattenBackendParameters(policy.BackendParams),
			"conditions":         flattenPolicyConditions(policy.Conditions),
		}
	}

	return result
}

func flattenWebPolicy(policies []apis.PolicyWebResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		retryCount := policy.RetryCount
		wp := map[string]interface{}{
			"name":             policy.Name,
			"request_protocol": policy.ReqProtocol,
			"request_method":   policy.ReqMethod,
			"effective_mode":   policy.EffectMode,
			"path":             policy.ReqURI,
			"timeout":          policy.Timeout,
			"retry_count":      utils.StringToInt(&retryCount),
			"authorizer_id":    policy.AuthorizerId,
			"backend_params":   flattenBackendParameters(policy.BackendParams),
			"conditions":       flattenPolicyConditions(policy.Conditions),
		}
		// which policy use backend address or vpc channel.
		if policy.VpcChannelInfo.VpcChannelId != "" {
			wp["vpc_channel_id"] = policy.VpcChannelInfo.VpcChannelId
			wp["host_header"] = policy.VpcChannelInfo.VpcChannelProxyHost
		} else {
			wp["backend_address"] = policy.DomainURL
		}

		result[i] = wp
	}

	return result
}

func resourceApiRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Id()
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	resp, err := apis.Get(client, instanceId, apiId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated API")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_id", resp.GroupId),
		d.Set("name", resp.Name),
		d.Set("authorizer_id", resp.AuthorizerId),
		d.Set("tags", resp.Tags),
		d.Set("request_protocol", resp.ReqProtocol),
		d.Set("request_method", resp.ReqMethod),
		d.Set("request_path", resp.ReqURI),
		d.Set("security_authentication", resp.AuthType),
		d.Set("cors", resp.Cors),
		d.Set("description", resp.Description),
		d.Set("body_description", resp.BodyDescription),
		d.Set("success_response", resp.ResultNormalSample),
		d.Set("failure_response", resp.ResultFailureSample),
		d.Set("response_id", resp.ResponseId),
		d.Set("type", analyseApiType(resp.Type)),
		d.Set("request_params", flattenApiRequestParams(resp.ReqParams)),
		d.Set("backend_params", flattenBackendParameters(resp.BackendParams)),
		d.Set("matching", analyseApiMatchMode(resp.MatchMode)),
		d.Set("simple_authentication", analyseAppSimpleAuth(resp.AuthOpt)),
		d.Set("mock", flattenMockStructure(resp.MockInfo)),
		d.Set("mock_policy", flattenMockPolicy(resp.PolicyMocks)),
		d.Set("func_graph", flattenFuncGraphStructure(resp.FuncInfo)),
		d.Set("func_graph_policy", flattenFuncGraphPolicy(resp.PolicyFunctions)),
		d.Set("web", flattenWebStructure(resp.WebInfo, d.Get("web.0.ssl_enable").(bool))),
		d.Set("web_policy", flattenWebPolicy(resp.PolicyWebs)),
		d.Set("registered_at", resp.RegisterTime),
		d.Set("updated_at", resp.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving API fields: %s", err)
	}
	return nil
}

func resourceApiUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Id()
	)
	opt, err := buildApiCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to build the API updateOpts: %s", err)
	}
	_, err = apis.Update(client, instanceId, apiId, opt).Extract()
	if err != nil {
		return diag.Errorf("error updating API (%s): %s", apiId, err)
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
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Id()
	)
	if err = deleteApi(client, instanceId, apiId); err != nil {
		return diag.Errorf("unable to delete the API (%s): %s", apiId, err)
	}

	return nil
}

// GetApigAPIIdByName is a method to get a specifies API ID from a APIG instance by name.
func GetApiIdByName(client *golangsdk.ServiceClient, instanceId, name string) (string, error) {
	opt := apis.ListOpts{
		Name: name, // Fuzzy search (reduce the time cost of the traversal)
	}
	pages, err := apis.List(client, instanceId, opt).AllPages()
	if err != nil {
		return "", fmt.Errorf("error retrieving APIs: %s", err)
	}
	apiRecords, err := apis.ExtractApis(pages)
	if err != nil {
		return "", err
	}
	for _, apiRecord := range apiRecords {
		if apiRecord.Name == name {
			return apiRecord.ID, nil
		}
	}
	return "", fmt.Errorf("unable to find the API (%s) form APIG service", name)
}

func resourceApiImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}
	name := parts[1]
	instanceId := parts[0]
	apiId, err := GetApiIdByName(client, instanceId, name)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(apiId)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
