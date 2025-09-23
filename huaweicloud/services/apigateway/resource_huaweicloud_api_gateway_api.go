package apigateway

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG DELETE /v1.0/apigw/apis/{id}
// @API APIG GET /v1.0/apigw/apis/{id}
// @API APIG PUT /v1.0/apigw/apis/{id}
// @API APIG POST /v1.0/apigw/apis
func ResourceAPI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGatewayAPICreate,
		ReadContext:   resourceAPIGatewayAPIRead,
		UpdateContext: resourceAPIGatewayAPIUpdate,
		DeleteContext: resourceAPIGatewayAPIDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"visibility": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"APP", "IAM", "NONE",
				}, false),
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HTTPS",
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP", "HTTPS", "BOTH",
				}, false),
			},
			"request_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backend_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP", "FUNCTION", "MOCK",
				}, false),
			},
			"http_backend": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"url_domain": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"http_backend.0.vpc_channel"},
						},
						"vpc_channel": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  50000,
						},
					},
				},
			},
			"mock_backend": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"function_backend": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_urn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"invocation_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"sync", "async",
							}, true),
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  50000,
						},
					},
				},
			},
			"request_parameter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"required": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"backend_parameter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "REQUEST",
							DiffSuppressFunc: utils.SuppressCaseDiffs(),
							ValidateFunc: validation.StringInSlice([]string{
								"REQUEST", "CONSTANT", "SYSTEM",
							}, true),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cors": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"example_success_response": {
				Type:     schema.TypeString,
				Required: true,
			},
			"example_failure_response": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAPIGatewayAPICreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	createOpts, err := buildApiParameter(d)
	if err != nil {
		return diag.Errorf("error creating API Gateway parameter: %s", err)
	}

	log.Printf("[DEBUG] create API options: %#v", createOpts)
	v, err := apis.Create(apigwClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating API Gateway API: %s", err)
	}

	// Store the ID now
	d.SetId(v.Id)

	return resourceAPIGatewayAPIRead(ctx, d, meta)
}

func resourceAPIGatewayAPIRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	v, err := apis.Get(apigwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "API Gateway API")
	}

	log.Printf("[DEBUG] retrieved API Gateway API %s: %+v", d.Id(), v)

	mErr := multierror.Append(
		d.Set("group_id", v.GroupId),
		d.Set("group_name", v.GroupName),
		d.Set("name", v.Name),
		d.Set("description", v.Remark),
		d.Set("tags", v.Tags),
		d.Set("version", v.Version),
		d.Set("visibility", v.Type),
		d.Set("auth_type", v.AuthType),
		d.Set("request_protocol", v.ReqProtocol),
		d.Set("request_method", v.ReqMethod),
		d.Set("request_uri", v.ReqUri),
		d.Set("backend_type", v.BackendType),
		d.Set("example_success_response", v.ResultNormalSample),
		d.Set("example_failure_response", v.ResultFailureSample),
		d.Set("cors", v.Cors),
	)

	var requestParameters []map[string]interface{}
	for _, val := range v.ReqParams {
		parameters := make(map[string]interface{})
		parameters["name"] = val.Name
		parameters["location"] = val.Location
		parameters["type"] = val.Type
		if val.Required == 1 {
			parameters["required"] = true
		} else if val.Required == 2 {
			parameters["required"] = false
		}
		parameters["default"] = val.DefaultValue
		parameters["description"] = val.Remark
		requestParameters = append(requestParameters, parameters)
	}

	mErr = multierror.Append(mErr, d.Set("request_parameter", requestParameters))
	var backendParameters []map[string]interface{}
	for _, val := range v.BackendParams {
		parameters := make(map[string]interface{})
		parameters["name"] = val.Name
		parameters["location"] = val.Location
		parameters["value"] = val.Value
		parameters["type"] = val.Origin
		parameters["description"] = val.Remark
		backendParameters = append(backendParameters, parameters)
	}

	mErr = multierror.Append(mErr, d.Set("backend_parameter", backendParameters))
	backend := make([]map[string]interface{}, 0, 1)
	switch v.BackendType {
	case "HTTP":
		httpInfo := map[string]interface{}{
			"protocol":    v.BackendInfo.Protocol,
			"method":      v.BackendInfo.Method,
			"uri":         v.BackendInfo.Uri,
			"url_domain":  v.BackendInfo.UrlDomain,
			"vpc_channel": v.BackendInfo.VpcInfo,
			"timeout":     v.BackendInfo.Timeout,
		}
		backend = append(backend, httpInfo)
		mErr = multierror.Append(mErr, d.Set("http_backend", backend))
	case "FUNCTION":
		functionInfo := map[string]interface{}{
			"function_urn":    v.FunctionInfo.FunctionUrn,
			"invocation_type": v.FunctionInfo.InvocationType,
			"version":         v.FunctionInfo.Version,
			"timeout":         v.FunctionInfo.Timeout,
		}
		backend = append(backend, functionInfo)
		mErr = multierror.Append(mErr, d.Set("function_backend", backend))
	case "MOCK":
		mockInfo := map[string]interface{}{
			"result_content": v.MockInfo.ResultContent,
			"version":        v.MockInfo.Version,
			"description":    v.MockInfo.Remark,
		}
		backend = append(backend, mockInfo)
		mErr = multierror.Append(mErr, d.Set("mock_backend", backend))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAPIGatewayAPIUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	updateOpts, err := buildApiParameter(d)
	if err != nil {
		return diag.Errorf("error creating API Gateway options: %s", err)
	}

	log.Printf("[DEBUG] update API options: %#v", updateOpts)
	_, err = apis.Update(apigwClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating API Gateway API: %s", err)
	}

	return resourceAPIGatewayAPIRead(ctx, d, meta)
}

func resourceAPIGatewayAPIDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	if err := apis.Delete(apigwClient, d.Id()).ExtractErr(); err != nil {
		return common.CheckDeletedDiag(d, err, "API apis")
	}

	return nil
}

func buildApiParameter(d *schema.ResourceData) (*apis.CreateOpts, error) {
	backendType := d.Get("backend_type").(string)
	v := d.Get("tags").(*schema.Set)
	apiTags := buildApiTags(v)

	opts := &apis.CreateOpts{
		GroupId:             d.Get("group_id").(string),
		Name:                d.Get("name").(string),
		Remark:              d.Get("description").(string),
		Type:                d.Get("visibility").(int),
		AuthType:            d.Get("auth_type").(string),
		ReqProtocol:         d.Get("request_protocol").(string),
		ReqMethod:           d.Get("request_method").(string),
		ReqUri:              d.Get("request_uri").(string),
		BackendType:         backendType,
		Tags:                apiTags,
		ResultNormalSample:  d.Get("example_success_response").(string),
		ResultFailureSample: d.Get("example_failure_response").(string),
		Cors:                d.Get("cors").(bool),
	}

	switch backendType {
	case "HTTP":
		httpBackend := buildHttpBackendParam(d)
		if httpBackend == nil {
			return nil, fmt.Errorf("the argument \"http_backend\" is required under HTTP backend type")
		}
		opts.BackendOpts = *httpBackend
	case "FUNCTION":
		funcBackend := buildFunctionBackendParam(d)
		if funcBackend == nil {
			return nil, fmt.Errorf("the argument \"function_backend\" is required under FUNCTION backend type")
		}
		opts.FunctionOpts = *funcBackend
	case "MOCK":
		mockBackend := buildMockBackendParam(d)
		if mockBackend == nil {
			return nil, fmt.Errorf("the argument \"mock_backend\" is required under MOCK backend type")
		}
		opts.MockOpts = *mockBackend
	}
	opts.ReqParams = buildRequestParameters(d)
	opts.BackendParams = buildBackendParameters(d)

	return opts, nil
}

func buildApiTags(v *schema.Set) []string {
	var tags []string
	for _, v := range v.List() {
		tags = append(tags, v.(string))
	}
	return tags
}

func buildHttpBackendParam(d *schema.ResourceData) *apis.BackendOpts {
	raw := d.Get("http_backend").([]interface{})

	if len(raw) == 1 {
		httpBackend := &apis.BackendOpts{}
		cfg := raw[0].(map[string]interface{})
		httpBackend.Protocol = cfg["protocol"].(string)
		httpBackend.Method = cfg["method"].(string)
		httpBackend.Uri = cfg["uri"].(string)
		httpBackend.Timeout = cfg["timeout"].(int)

		if v, ok := cfg["vpc_channel"]; ok && v.(string) != "" {
			httpBackend.VpcStatus = 1
			httpBackend.VpcInfo.VpcId = v.(string)
		} else {
			httpBackend.VpcStatus = 2
			httpBackend.UrlDomain = cfg["url_domain"].(string)
		}
		return httpBackend
	}

	return nil
}

func buildFunctionBackendParam(d *schema.ResourceData) *apis.FunctionOpts {
	raw := d.Get("function_backend").([]interface{})

	if len(raw) == 1 {
		funcBackend := &apis.FunctionOpts{}
		cfg := raw[0].(map[string]interface{})
		funcBackend.FunctionUrn = cfg["function_urn"].(string)
		funcBackend.InvocationType = cfg["invocation_type"].(string)
		funcBackend.Version = cfg["version"].(string)
		funcBackend.Timeout = cfg["timeout"].(int)
		return funcBackend
	}

	return nil
}

func buildMockBackendParam(d *schema.ResourceData) *apis.MockOpts {
	raw := d.Get("mock_backend").([]interface{})

	// all parameters of mock_backend are optional
	mockBackend := &apis.MockOpts{}
	if len(raw) == 1 {
		cfg := raw[0].(map[string]interface{})
		mockBackend.ResultContent = cfg["result_content"].(string)
		mockBackend.Version = cfg["version"].(string)
		mockBackend.Remark = cfg["description"].(string)
	}

	return mockBackend
}

func buildRequestParameters(d *schema.ResourceData) []apis.RequestParameter {
	var requestList []apis.RequestParameter

	rawParams := d.Get("request_parameter").([]interface{})
	for i := range rawParams {
		parameter := rawParams[i].(map[string]interface{})
		request := apis.RequestParameter{
			Name:     parameter["name"].(string),
			Location: parameter["location"].(string),
			Type:     parameter["type"].(string),
			Remark:   parameter["description"].(string),
			// disable validity check
			ValidEnable: 2,
		}
		if parameter["required"].(bool) {
			request.Required = 1
		} else {
			request.Required = 2
			// the default value is used when the input parameter was optional
			request.DefaultValue = parameter["default"].(string)
		}
		requestList = append(requestList, request)
	}
	return requestList
}

func buildBackendParameters(d *schema.ResourceData) []apis.BackendParameter {
	var backendList []apis.BackendParameter

	rawParams := d.Get("backend_parameter").([]interface{})
	for i := range rawParams {
		parameter := rawParams[i].(map[string]interface{})
		request := apis.BackendParameter{
			Name:     parameter["name"].(string),
			Location: parameter["location"].(string),
			Origin:   parameter["type"].(string),
			Value:    parameter["value"].(string),
			Remark:   parameter["description"].(string),
		}

		backendList = append(backendList, request)
	}
	return backendList
}
