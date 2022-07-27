package apig

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	protocolHTTP  = "HTTP"
	protocolHTTPS = "HTTPS"
	protocolBOTH  = "BOTH"

	methodGET     = "GET"
	methodPOST    = "POST"
	methodPUT     = "PUT"
	methodDELETE  = "DELETE"
	methodHEAD    = "HEAD"
	methodPATCH   = "PATCH"
	methodOPTIONS = "OPTIONS"
	methodANY     = "ANY"
)

var (
	matching = map[string]string{
		"Prefix": "SWA",
		"Exact":  "NORMAL",
	}
	conditionType = map[string]string{
		"Equal":      "exact",
		"Enumerated": "enum",
		"Matching":   "pattern",
	}
)

func ResourceApigAPIV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigAPIV2Create,
		Read:   resourceApigAPIV2Read,
		Update: resourceApigAPIV2Update,
		Delete: resourceApigAPIV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigAPIResourceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Public", "Private",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z_0-9]{2,63})$"),
					"The name consists of 3 to 64 characters and only letters, digits, underscore (_) and chinese "+
						"characters are allowed. The name must start with a letter or chinese character."),
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					methodGET, methodPOST, methodPUT, methodDELETE, methodHEAD, methodPATCH, methodOPTIONS, methodANY,
				}, false),
			},
			"request_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					protocolHTTP, protocolHTTPS, protocolBOTH,
				}, false),
			},
			"request_params": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z0-9-_\\.]{0,31})$"),
								"The name can contains of 1 to 32 characters and start with a letter."+
									"Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed."),
						},
						"required": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "PATH",
							ValidateFunc: validation.StringInSlice([]string{
								"PATH", "HEADER", "QUERY",
							}, false),
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "STRING",
							ValidateFunc: validation.StringInSlice([]string{
								"STRING", "NUMBER",
							}, false),
						},
						"maximum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minimum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"example": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The example value contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The default value contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
								"The description contain a maximum of 255 characters, "+
									"and the angle brackets (< and >) are not allowed."),
						},
					},
				},
			},
			"backend_params": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     backendParamSchemaResource(),
				Set:      resourceBackendParamtersHash,
			},
			"security_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE", "APP", "IAM", "AUTHORIZER",
				}, false),
			},
			"simple_authentication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"authorizer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"body_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"cors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"matching": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Exact",
				ValidateFunc: validation.StringInSlice([]string{
					"Exact", "Prefix",
				}, false),
			},
			"response_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"success_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
			},
			"failure_response": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 20480),
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
						"response": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 2048),
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"invocation_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "sync",
							ValidateFunc: validation.StringInSlice([]string{
								"async", "sync",
							}, false),
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
							Type:     schema.TypeString,
							Required: true,
						},
						"host_header": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"web.0.backend_address"},
						},
						"vpc_channel_id": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"web.0.backend_address"},
						},
						"backend_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								methodGET, methodPOST, methodPUT, methodDELETE, methodHEAD, methodPATCH, methodOPTIONS, methodANY,
							}, false),
						},
						"request_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  protocolHTTPS,
							ValidateFunc: validation.StringInSlice([]string{
								protocolHTTP, protocolHTTPS,
							}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"mock_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"func_graph", "web", "func_graph_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z0-9_]{2,63})$"),
								"The name can contains of 3 to 64 characters and start with a letter."+
									"Only letters, digits and underscores (_) are allowed."),
						},
						"conditions": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 5,
							Elem:     policyConditionSchemaResource(),
						},
						"response": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(8, 2048),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ANY",
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     backendParamSchemaResource(),
							Set:      resourceBackendParamtersHash,
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"func_graph_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "web", "mock_policy", "web_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z0-9_]{2,63})$"),
								"The name can contains of 3 to 64 characters and start with a letter."+
									"Only letters, digits and underscores (_) are allowed."),
						},
						"function_urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"conditions": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 5,
							Elem:     policyConditionSchemaResource(),
						},
						"invocation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "sync",
							ValidateFunc: validation.StringInSlice([]string{
								"async", "sync",
							}, false),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ANY",
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     backendParamSchemaResource(),
							Set:      resourceBackendParamtersHash,
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"web_policy": {
				Type:          schema.TypeSet,
				MaxItems:      5,
				Optional:      true,
				ConflictsWith: []string{"mock", "func_graph", "mock_policy", "func_graph_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z0-9_]{2,63})$"),
								"The name can contains of 3 to 64 characters and start with a letter."+
									"Only letters, digits and underscores (_) are allowed."),
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"request_method": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								methodGET, methodPOST, methodPUT, methodDELETE, methodHEAD, methodPATCH, methodOPTIONS, methodANY,
							}, false),
						},
						"conditions": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 5,
							Elem:     policyConditionSchemaResource(),
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpc_channel_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backend_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								protocolHTTP, protocolHTTPS,
							}, false),
						},
						"effective_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ANY",
							ValidateFunc: validation.StringInSlice([]string{
								"ALL", "ANY",
							}, false),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5000,
							ValidateFunc: validation.IntBetween(1, 600000),
						},
						"backend_params": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     backendParamSchemaResource(),
							Set:      resourceBackendParamtersHash,
						},
						"authorizer_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"register_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func policyConditionSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"param_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "param",
				ValidateFunc: validation.StringInSlice([]string{
					"param", "source",
				}, false),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice([]string{
					"Equal", "Enumerated", "Matching",
				}, false),
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
					"REQUEST", "CONSTANT", "SYSTEM",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([A-Za-z][A-Za-z0-9-_\\.]{0,31})$"),
					"The name can contains of 1 to 32 characters and start with a letter."+
						"Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed."),
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PATH", "QUERY", "HEADER",
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The value contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{1,255}$"),
					"The description contain a maximum of 255 characters, "+
						"and the angle brackets (< and >) are not allowed."),
			},
		},
	}
}

func buildAPIType(t string) (int, error) {
	apiType := map[string]int{
		"Public":  1,
		"Private": 2,
	}
	if v, ok := apiType[t]; ok {
		return v, nil
	}
	return -1, fmtp.Errorf("The API type (%s) is invalid", t)
}

func buildAPIMock(mock []interface{}) *apis.Mock {
	mockMap := mock[0].(map[string]interface{})
	cont := mockMap["response"].(string)
	authId := mockMap["authorizer_id"].(string)

	return &apis.Mock{
		ResultContent: &cont,
		AuthorizerId:  &authId,
	}
}

func buildAPIFuncGraph(funcGraph []interface{}) *apis.FuncGraph {
	funcMap := funcGraph[0].(map[string]interface{})
	authId := funcMap["authorizer_id"].(string)

	return &apis.FuncGraph{
		Timeout:        funcMap["timeout"].(int),
		InvocationType: funcMap["invocation_type"].(string),
		FunctionUrn:    funcMap["function_urn"].(string),
		Version:        funcMap["version"].(string),
		AuthorizerId:   &authId,
	}
}

func buildApigAPIWeb(web []interface{}) *apis.Web {
	webMap := web[0].(map[string]interface{})
	sslEnable := webMap["ssl_enable"].(bool)
	authId := webMap["authorizer_id"].(string)

	result := apis.Web{
		ReqURI:          webMap["path"].(string),
		ReqMethod:       webMap["request_method"].(string),
		ReqProtocol:     webMap["request_protocol"].(string),
		Timeout:         webMap["timeout"].(int),
		ClientSslEnable: &sslEnable,
		AuthorizerId:    &authId,
	}
	// If vpc_channel_id is empty, the backend address is used.
	if chanId, ok := webMap["vpc_channel_id"]; ok && chanId != "" {
		result.VpcChannelStatus = 1
		result.VpcChannelInfo = &apis.VpcChannel{
			VpcChannelId:        chanId.(string),
			VpcChannelProxyHost: webMap["host_header"].(string),
		}
	} else {
		result.VpcChannelStatus = 2
		result.DomainURL = webMap["backend_address"].(string)
	}

	return &result
}

func buildAPIRequestParameters(requests *schema.Set) []apis.ReqParamBase {
	result := make([]apis.ReqParamBase, requests.Len())
	for i, v := range requests.List() {
		paramMap := v.(map[string]interface{})
		paramType := paramMap["type"].(string)
		desc := paramMap["description"].(string)
		param := apis.ReqParamBase{
			Name:        paramMap["name"].(string),
			Location:    paramMap["location"].(string),
			Description: &desc,
		}
		if paramType == "NUMBER" {
			param.MaxNum = golangsdk.IntToPointer(paramMap["maximum"].(int))
			param.MinNum = golangsdk.IntToPointer(paramMap["minimum"].(int))
		} else if paramType == "STRING" {
			param.MaxSize = golangsdk.IntToPointer(paramMap["maximum"].(int))
			param.MinSize = golangsdk.IntToPointer(paramMap["minimum"].(int))
		}
		param.Type = paramType

		if paramMap["required"].(bool) {
			param.Required = 1
		} else {
			param.Required = 2
		}
		result[i] = param
	}
	return result
}

func buildAPIBackendParameterValue(origin, value string) string {
	// The internal parameters of the system parameters include as below:
	internalParams := []string{
		"sourceIp", "stage", "apiId", "requestId", "serverAddr", "serverName", "handleTime", "providerAppId",
	}
	if origin == "SYSTEM" {
		// If the system parameters are configured as internal parameters, the internal format is used to construct.
		if utils.StrSliceContains(internalParams, value) {
			return fmt.Sprintf("$context.%s", value)
		}
		// The fornt-end format is used to construct.
		return fmt.Sprintf("$context.authorizer.frontend.%s", value)
	}
	return value
}

// For backend API, the parameters contains request parameters and constant parameters.
func buildAPIBackendParameters(backends *schema.Set) []apis.BackendParamBase {
	result := make([]apis.BackendParamBase, backends.Len())
	for i, v := range backends.List() {
		paramMap := v.(map[string]interface{})
		origin := paramMap["type"].(string)
		param := apis.BackendParamBase{
			Origin:   origin,
			Name:     paramMap["name"].(string),
			Location: paramMap["location"].(string),
			Value:    buildAPIBackendParameterValue(origin, paramMap["value"].(string)),
		}
		if origin != "REQUEST" {
			desc := paramMap["description"].(string)
			param.Description = &desc
		}
		result[i] = param
	}

	return result
}

func buildAPIPolicyConditions(conditions *schema.Set) ([]apis.APIConditionBase, error) {
	result := make([]apis.APIConditionBase, conditions.Len())
	for i, v := range conditions.List() {
		conditionMap := v.(map[string]interface{})
		condition := apis.APIConditionBase{
			ReqParamName:    conditionMap["param_name"].(string),
			ConditionOrigin: conditionMap["source"].(string),
			ConditionValue:  conditionMap["value"].(string),
		}
		conType := conditionMap["type"].(string)
		if v, ok := conditionType[conType]; ok {
			condition.ConditionType = v
		} else {
			return result, fmtp.Errorf("The condition type is invalid")
		}
		result[i] = condition
	}
	return result, nil
}

func buildAPIMockPolicy(mocks *schema.Set) ([]apis.PolicyMock, error) {
	result := make([]apis.PolicyMock, mocks.Len())
	for i, policy := range mocks.List() {
		pm := policy.(map[string]interface{})
		condition, err := buildAPIPolicyConditions(pm["conditions"].(*schema.Set))
		if err != nil {
			return result, err
		}
		result[i] = apis.PolicyMock{
			Name:          pm["name"].(string),
			ResultContent: pm["response"].(string),
			EffectMode:    pm["effective_mode"].(string),
			Conditions:    condition,
			BackendParams: buildAPIBackendParameters(pm["backend_params"].(*schema.Set)),
		}
	}
	return result, nil
}

func buildAPIFuncGraphPolicy(policies *schema.Set) ([]apis.PolicyFuncGraph, error) {
	result := make([]apis.PolicyFuncGraph, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		condition, err := buildAPIPolicyConditions(pm["conditions"].(*schema.Set))
		if err != nil {
			return result, err
		}
		result[i] = apis.PolicyFuncGraph{
			Name:           pm["name"].(string),
			FunctionUrn:    pm["function_urn"].(string),
			InvocationType: pm["invocation_mode"].(string),
			EffectMode:     pm["effective_mode"].(string),
			Timeout:        pm["timeout"].(int),
			Conditions:     condition,
			BackendParams:  buildAPIBackendParameters(pm["backend_params"].(*schema.Set)),
		}
	}
	return result, nil
}

func buildApigAPIWebPolicy(policies *schema.Set) ([]apis.PolicyWeb, error) {
	result := make([]apis.PolicyWeb, policies.Len())
	for i, policy := range policies.List() {
		pm := policy.(map[string]interface{})
		condition, err := buildAPIPolicyConditions(pm["conditions"].(*schema.Set))
		if err != nil {
			return result, err
		}
		wp := apis.PolicyWeb{
			Name:          pm["name"].(string),
			ReqProtocol:   pm["request_protocol"].(string),
			ReqMethod:     pm["request_method"].(string),
			ReqURI:        pm["path"].(string),
			EffectMode:    pm["effective_mode"].(string),
			Timeout:       pm["timeout"].(int),
			DomainURL:     pm["host_header"].(string),
			Conditions:    condition,
			BackendParams: buildAPIBackendParameters(pm["backend_params"].(*schema.Set)),
		}
		if chanId, ok := pm["vpc_channel_id"]; ok {
			if chanId != "" {
				wp.VpcChannelInfo = &apis.VpcChannel{
					VpcChannelId:        pm["vpc_channel_id"].(string),
					VpcChannelProxyHost: pm["host_header"].(string),
				}
				wp.VpcChannelStatus = 1
			} else {
				wp.VpcChannelStatus = 2
			}
		}
		result[i] = wp
	}
	return result, nil
}

func buildApigAPIParameters(d *schema.ResourceData) (apis.APIOpts, error) {
	apiType, err := buildAPIType(d.Get("type").(string))
	if err != nil {
		return apis.APIOpts{}, err
	}
	cors := d.Get("cors").(bool)
	desc := d.Get("description").(string)
	bodyDesc := d.Get("body_description").(string)
	successResp := d.Get("success_response").(string)
	failureResp := d.Get("failure_response").(string)
	authType := d.Get("security_authentication").(string)
	opt := apis.APIOpts{
		Type:                apiType,
		AuthorizerId:        d.Get("authorizer_id").(string),
		GroupId:             d.Get("group_id").(string),
		Name:                d.Get("name").(string),
		ReqProtocol:         d.Get("request_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqURI:              d.Get("request_path").(string),
		Cors:                &cors,
		AuthType:            authType,
		MatchMode:           d.Get("matching").(string),
		Description:         &desc,
		BodyDescription:     &bodyDesc,
		ResultNormalSample:  &successResp,
		ResultFailureSample: &failureResp,
		ResponseId:          d.Get("response_id").(string),
	}
	// build match mode
	v, ok := matching[d.Get("matching").(string)]
	if !ok {
		return opt, fmtp.Errorf("Unable to extract match mode")
	}
	opt.MatchMode = v

	isSimpleAuthEnabled := d.Get("simple_authentication").(bool)
	if authType == "APP" {
		if isSimpleAuthEnabled {
			opt.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: "HEADER",
			}
		} else {
			opt.AuthOpt = &apis.AuthOpt{
				AppCodeAuthType: "DISABLE",
			}
		}
	} else if isSimpleAuthEnabled {
		return opt, fmtp.Errorf("The security authentication must be 'APP' if simple authentication is true")
	}

	opt.ReqParams = buildAPIRequestParameters(d.Get("request_params").(*schema.Set))
	opt.BackendParams = buildAPIBackendParameters(d.Get("backend_params").(*schema.Set))

	// build backend (one of the mock, function graph and web) server and related policies.
	if m, ok := d.GetOk("mock"); ok {
		opt.BackendType = "MOCK"
		opt.MockInfo = buildAPIMock(m.([]interface{}))
		mp := d.Get("mock_policy").(*schema.Set)
		policy, err := buildAPIMockPolicy(mp)
		if err != nil {
			return opt, err
		}
		opt.PolicyMocks = policy
	} else if fg, ok := d.GetOk("func_graph"); ok {
		opt.BackendType = "FUNCTION"
		opt.FuncInfo = buildAPIFuncGraph(fg.([]interface{}))
		fgp := d.Get("func_graph_policy").(*schema.Set)
		policy, err := buildAPIFuncGraphPolicy(fgp)
		if err != nil {
			return opt, err
		}
		opt.PolicyFunctions = policy
	} else {
		opt.BackendType = protocolHTTP
		web := d.Get("web").([]interface{})
		opt.WebInfo = buildApigAPIWeb(web)
		wp := d.Get("web_policy").(*schema.Set)
		policy, err := buildApigAPIWebPolicy(wp)
		if err != nil {
			return opt, err
		}
		opt.PolicyWebs = policy
	}

	logp.Printf("[DEBUG] The API Opts is : %+v", opt)
	return opt, nil
}

func resourceApigAPIV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opt, err := buildApigAPIParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to build the Api parameter: %s", err)
	}

	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := apis.Create(client, instanceId, opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating APIG v2 API: %s", err)
	}
	d.SetId(resp.ID)

	return resourceApigAPIV2Read(d, meta)
}

func getAPIBackendParameterValue(origin, value string) (string, error) {
	if origin == "SYSTEM" {
		// The system parameters include gateway parameters and front-end parameters.
		if ok, err := regexp.MatchString(`\$context\.authorizer\.frontend\.`, value); ok && err == nil {
			regex, err := regexp.Compile(`\$context\.authorizer\.frontend\.(.*)`)
			if err != nil {
				return "", err
			}
			result := regex.FindStringSubmatch(value)
			fmt.Println(result)
			if len(result) < 2 {
				return "", fmtp.Errorf("Wrong front-end parameter format: %s", value)
			}
			return result[1], nil
		}
		if ok, err := regexp.MatchString(`\$context\.`, value); ok && err == nil {
			regex, err := regexp.Compile(`\$context\.(.*)`)
			if err != nil {
				return "", err
			}
			result := regex.FindStringSubmatch(value)
			fmt.Println(result)
			if len(result) < 2 {
				return "", fmtp.Errorf("Wrong gateway parameter format: %s", value)
			}
			return result[1], nil
		}
		return "", fmtp.Errorf("[ERROR] The parameter %s is not in the format of a gateway or front-end", value)
	}
	return value, nil
}

func getAPIBackendParameters(backendParams []apis.BackendParamResp) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(backendParams))
	for i, v := range backendParams {
		origin := v.Origin
		value, err := getAPIBackendParameterValue(v.Origin, v.Value)
		if err != nil {
			return result, err
		}
		param := map[string]interface{}{
			"type":     origin,
			"name":     v.Name,
			"location": v.Location,
			"value":    value,
		}
		if origin != "REQUEST" {
			param["description"] = v.Description
		}
		result[i] = param
	}
	return result, nil
}

func getAPIPolicyConditions(conditions []apis.APIConditionBase) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(conditions))
	for i, v := range conditions {
		conType, err := analysisConditionType(v.ConditionType)
		if err != nil {
			return result, err
		}
		result[i] = map[string]interface{}{
			"source":     v.ConditionOrigin,
			"param_name": v.ReqParamName,
			"type":       conType,
			"value":      v.ConditionValue,
		}
	}
	return result, nil
}

func analysisConditionType(conType string) (string, error) {
	for k, v := range conditionType {
		if v == conType {
			return k, nil
		}
	}
	return "", fmtp.Errorf("The condition type (%s) is invalid", conType)
}

func setApigAPIType(d *schema.ResourceData, t int) error {
	apiType := map[int]string{
		1: "Public",
		2: "Private",
	}
	if v, ok := apiType[t]; ok {
		return d.Set("type", v)
	}
	return fmtp.Errorf("The API type (%d) is invalid", t)
}

func setApigAPIMatchMode(d *schema.ResourceData, mode string) error {
	for k, v := range matching {
		if v == mode {
			return d.Set("matching", k)
		}
	}
	return fmtp.Errorf("The API matching mode is invalid", mode)
}

func setApigAPIAppSimpleAuth(d *schema.ResourceData, opt apis.AuthOpt) error {
	// HEADER: AppCode authentication is enabled and the AppCode is located in the header.
	if opt.AppCodeAuthType == "HEADER" {
		return d.Set("simple_authentication", true)
	}
	return d.Set("simple_authentication", false)
}

func setApigAPIReqParams(d *schema.ResourceData, reqParams []apis.ReqParamResp) error {
	result := make([]map[string]interface{}, len(reqParams))
	for i, v := range reqParams {
		param := map[string]interface{}{
			"name":        v.Name,
			"location":    v.Location,
			"type":        v.Type,
			"example":     v.SampleValue,
			"default":     v.DefaultValue,
			"description": v.Description,
		}
		if v.Type == "NUMBER" {
			param["maximum"] = v.MaxNum
			param["minimum"] = v.MinNum
		} else if v.Type == "STRING" {
			param["maximum"] = v.MaxSize
			param["minimum"] = v.MinSize
		}
		if v.Required == 1 {
			param["required"] = true
		} else if v.Required == 2 {
			param["required"] = false
		} else {
			return fmtp.Errorf("Invalid value of 'required'")
		}
		result[i] = param
	}
	return d.Set("request_params", result)
}

func setApigAPIMock(d *schema.ResourceData, mockResp apis.Mock) error {
	if mockResp == (apis.Mock{}) {
		return d.Set("mock", nil)
	}
	result := []map[string]interface{}{
		{
			"response":      mockResp.ResultContent,
			"authorizer_id": mockResp.AuthorizerId,
		},
	}
	return d.Set("mock", result)
}

func setApigAPIFuncGraph(d *schema.ResourceData, funcResp apis.FuncGraph) error {
	if funcResp == (apis.FuncGraph{}) {
		return d.Set("func_graph", nil)
	}
	result := []map[string]interface{}{
		{
			"function_urn":    funcResp.FunctionUrn,
			"timeout":         funcResp.Timeout,
			"invocation_type": funcResp.InvocationType,
			"version":         funcResp.Version,
			"authorizer_id":   funcResp.AuthorizerId,
		},
	}
	return d.Set("func_graph", result)
}

func setApigAPIWeb(d *schema.ResourceData, webResp apis.Web) error {
	if webResp == (apis.Web{}) {
		return d.Set("web", nil)
	}
	result := make([]map[string]interface{}, 1)
	web := map[string]interface{}{
		"path":             webResp.ReqURI,
		"request_method":   webResp.ReqMethod,
		"request_protocol": webResp.ReqProtocol,
		"timeout":          webResp.Timeout,
		"ssl_enable":       d.Get("web.0.ssl_enable"),
		"authorizer_id":    webResp.AuthorizerId,
	}
	if webResp.VpcChannelInfo.VpcChannelId != "" {
		web["vpc_channel_id"] = webResp.VpcChannelInfo.VpcChannelId
		web["host_header"] = webResp.VpcChannelInfo.VpcChannelProxyHost
	} else {
		web["backend_address"] = webResp.DomainURL
	}
	result[0] = web

	return d.Set("web", result)
}

func setApigAPIMockPolicy(d *schema.ResourceData, policies []apis.PolicyMockResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		mp := map[string]interface{}{
			"name":           policy.Name,
			"response":       policy.ResultContent,
			"effective_mode": policy.EffectMode,
			"authorizer_id":  policy.AuthorizerId,
		}
		backendParams, err := getAPIBackendParameters(policy.BackendParams)
		if err != nil {
			return err
		}
		mp["backend_params"] = backendParams
		condition, err := getAPIPolicyConditions(policy.Conditions)
		if err != nil {
			return fmtp.Errorf("Error setting policy (%s): %s", policy.Name, err)
		}
		mp["conditions"] = condition

		result[i] = mp
	}

	return d.Set("mock_policy", result)
}

func setApigAPIFuncGraphPolicy(d *schema.ResourceData, policies []apis.PolicyFuncGraphResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		fgp := map[string]interface{}{
			"name":            policy.Name,
			"function_urn":    policy.FunctionUrn,
			"version":         policy.Version,
			"invocation_mode": policy.InvocationType,
			"effective_mode":  policy.EffectMode,
			"timeout":         policy.Timeout,
			"authorizer_id":   policy.AuthorizerId,
		}
		backendParams, err := getAPIBackendParameters(policy.BackendParams)
		if err != nil {
			return err
		}
		fgp["backend_params"] = backendParams
		condition, err := getAPIPolicyConditions(policy.Conditions)
		if err != nil {
			return fmtp.Errorf("Error setting policy (%s): %s", policy.Name, err)
		}
		fgp["conditions"] = condition

		result[i] = fgp
	}

	return d.Set("func_graph_policy", result)
}

func setApigAPIWebPolicy(d *schema.ResourceData, policies []apis.PolicyWebResp) error {
	result := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		wp := map[string]interface{}{
			"name":             policy.Name,
			"request_protocol": policy.ReqProtocol,
			"request_method":   policy.ReqMethod,
			"effective_mode":   policy.EffectMode,
			"path":             policy.ReqURI,
			"timeout":          policy.Timeout,
			"authorizer_id":    policy.AuthorizerId,
		}
		// which policy use backend address or vpc channel.
		if policy.VpcChannelInfo.VpcChannelId != "" {
			wp["vpc_channel_id"] = policy.VpcChannelInfo.VpcChannelId
			wp["host_header"] = policy.VpcChannelInfo.VpcChannelProxyHost
		} else {
			wp["backend_address"] = policy.DomainURL
		}
		backendParams, err := getAPIBackendParameters(policy.BackendParams)
		if err != nil {
			return err
		}
		wp["backend_params"] = backendParams
		condition, err := getAPIPolicyConditions(policy.Conditions)
		if err != nil {
			return fmtp.Errorf("Error setting policy (%s): %s", policy.Name, err)
		}
		wp["conditions"] = condition

		result[i] = wp
	}

	return d.Set("web_policy", result)
}

func setApigAPIParameters(d *schema.ResourceData, config *config.Config, resp *apis.APIResp) error {
	backendParams, err := getAPIBackendParameters(resp.BackendParams)
	if err != nil {
		return err
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("group_id", resp.GroupId),
		d.Set("name", resp.Name),
		d.Set("authorizer_id", resp.AuthorizerId),
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
		d.Set("backend_params", backendParams),
		setApigAPIType(d, resp.Type),
		setApigAPIReqParams(d, resp.ReqParams),
		setApigAPIMatchMode(d, resp.MatchMode),
		setApigAPIAppSimpleAuth(d, resp.AuthOpt),
		setApigAPIMock(d, resp.MockInfo),
		setApigAPIMockPolicy(d, resp.PolicyMocks),
		setApigAPIFuncGraph(d, resp.FuncInfo),
		setApigAPIFuncGraphPolicy(d, resp.PolicyFunctions),
		setApigAPIWeb(d, resp.WebInfo),
		setApigAPIWebPolicy(d, resp.PolicyWebs),
		d.Set("register_time", resp.RegisterTime),
		d.Set("update_time", resp.UpdateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceApigAPIV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := apis.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "error getting API information from server")
	}

	return setApigAPIParameters(d, config, resp)
}

func resourceApigAPIV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	opt, err := buildApigAPIParameters(d)
	if err != nil {
		return fmtp.Errorf("Unable to build the Api parameter: %s", err)
	}
	_, err = apis.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating APIG v2 API: %s", err)
	}

	return resourceApigAPIV2Read(d, meta)
}

func resourceApigAPIV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	if err = apis.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Unable to delete the APIG v2 dedicated instance (%s): %s", d.Id(), err)
	}
	d.SetId("")

	return nil
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

// GetApigAPIIdByName is a method to get a specifies API ID from a APIG instance by name.
func GetApigAPIIdByName(client *golangsdk.ServiceClient, instanceId, name string) (string, error) {
	opt := apis.ListOpts{
		Name: name,
	}
	pages, err := apis.List(client, instanceId, opt).AllPages()
	if err != nil {
		return "", fmtp.Errorf("Error retrieving APIs: %s", err)
	}
	resp, err := apis.ExtractApis(pages)
	if len(resp) < 1 {
		return "", fmtp.Errorf("Unable to find the API (%s) form server: %s", name, err)
	}
	return resp[0].ID, nil
}

func resourceApigAPIResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<name>")
	}
	name := parts[1]
	instanceId := parts[0]
	id, err := GetApigAPIIdByName(client, instanceId, name)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(id)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
