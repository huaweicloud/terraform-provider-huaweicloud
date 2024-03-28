package cdn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	cdnv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var httpsConfig = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"https_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"certificate_source": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"http2_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"https_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http2_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var requestAndResponseHeader = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"set", "delete",
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var authOpts = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"type_a", "type_b", "type_c1", "type_c2",
				}, false),
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"time_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dec", "hex",
				}, false),
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var forceRedirectAndCompress = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

var cacheUrlParameterFilter = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var ipFrequencyLimit = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var websocket = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var flexibleOrigin = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"back_sources": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sources_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_or_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"obs_bucket_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"http_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"match_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var remoteAuth = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"remote_auth_rules": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"request_method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"file_type_setting": {
							Type:     schema.TypeString,
							Required: true,
						},
						"reserve_args_setting": {
							Type:     schema.TypeString,
							Required: true,
						},
						"reserve_headers_setting": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_success_status": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_failed_status": {
							Type:     schema.TypeString,
							Required: true,
						},
						"response_status": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"timeout_action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"specified_file_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reserve_args": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"reserve_headers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"add_custom_args_rules":    &customArgs,
						"add_custom_headers_rules": &customArgs,
					},
				},
			},
		},
	},
}

var customArgs = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	},
}

// @API CDN POST /v1.0/cdn/domains
// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}
// @API CDN PUT /v1.0/cdn/domains/{domainId}/disable
// @API CDN DELETE /v1.0/cdn/domains/{domainId}
// @API CDN PUT /v1.1/cdn/configuration/domains/{domain_name}/configs
// @API CDN GET /v1.1/cdn/configuration/domains/{domain_name}/configs
// @API CDN POST /v1.0/cdn/configuration/tags/batch-delete
// @API CDN POST /v1.0/cdn/configuration/tags
func ResourceCdnDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnDomainCreate,
		ReadContext:   resourceCdnDomainRead,
		UpdateContext: resourceCdnDomainUpdate,
		DeleteContext: resourceCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCDNDomainImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"web", "download", "video", "wholeSite",
				}, true),
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:     schema.TypeString,
							Required: true,
						},
						"origin_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ipaddr", "domain", "obs_bucket",
							}, true),
						},
						"active": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"obs_web_hosting_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"http_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"retrieval_host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"follow", "http", "https",
							}, false),
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"range_based_retrieval_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"https_settings":             &httpsConfig,
						"retrieval_request_header":   &requestAndResponseHeader,
						"http_response_header":       &requestAndResponseHeader,
						"url_signing":                &authOpts,
						"force_redirect":             &forceRedirectAndCompress,
						"compress":                   &forceRedirectAndCompress,
						"cache_url_parameter_filter": &cacheUrlParameterFilter,
						"ip_frequency_limit":         &ipFrequencyLimit,
						"websocket":                  &websocket,
						"flexible_origin":            &flexibleOrigin,
						"remote_auth":                &remoteAuth,
					},
				},
			},

			"cache_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"follow_origin": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:     schema.TypeString,
										Required: true,
										DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
											// Convert several original types and change parameter types while ensuring
											// that the original configuration is available.
											// Notes: the state file no longer save the original types.
											return parseCacheRuleType(n) == o
										},
									},
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ttl": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
									},
									"ttl_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
										DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
											// Convert several original types and change parameter types while ensuring
											// that the original configuration is available.
											// Notes: the state file no longer save the original types.
											return parseCacheTTLUnits(n) == o
										},
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "schema: Required",
									},
								},
							},
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateDomainSources(d *schema.ResourceData) []domains.SourcesOpts {
	var sourceRequests []domains.SourcesOpts

	sources := d.Get("sources").(*schema.Set).List()
	for i := range sources {
		source := sources[i].(map[string]interface{})
		sourceRequest := domains.SourcesOpts{
			IporDomain:    source["origin"].(string),
			OriginType:    source["origin_type"].(string),
			ActiveStandby: source["active"].(int),
		}
		sourceRequests = append(sourceRequests, sourceRequest)
	}
	return sourceRequests
}

func buildHTTPSOpts(rawHTTPS []interface{}) *model.HttpPutBody {
	if len(rawHTTPS) != 1 {
		return nil
	}

	https := rawHTTPS[0].(map[string]interface{})
	httpsStatus := ""
	if https["https_enabled"].(bool) {
		httpsStatus = "on"
	}
	http2Status := ""
	if https["http2_enabled"].(bool) {
		http2Status = "on"
	}

	httpsOpts := model.HttpPutBody{
		HttpsStatus:       utils.StringIgnoreEmpty(httpsStatus),
		CertificateName:   utils.StringIgnoreEmpty(https["certificate_name"].(string)),
		CertificateValue:  utils.StringIgnoreEmpty(https["certificate_body"].(string)),
		PrivateKey:        utils.StringIgnoreEmpty(https["private_key"].(string)),
		CertificateSource: utils.Int32IgnoreEmpty(int32(https["certificate_source"].(int))),
		Http2Status:       utils.StringIgnoreEmpty(http2Status),
		TlsVersion:        utils.StringIgnoreEmpty(https["tls_version"].(string)),
	}

	return &httpsOpts
}

func buildOriginRequestHeaderOpts(rawOriginRequestHeader []interface{}) *[]model.OriginRequestHeader {
	if len(rawOriginRequestHeader) < 1 {
		return nil
	}

	originRequestHeaderOpts := make([]model.OriginRequestHeader, len(rawOriginRequestHeader))
	for i, v := range rawOriginRequestHeader {
		header := v.(map[string]interface{})
		originRequestHeaderOpts[i] = model.OriginRequestHeader{
			Name:   header["name"].(string),
			Value:  utils.StringIgnoreEmpty(header["value"].(string)),
			Action: header["action"].(string),
		}
	}

	return &originRequestHeaderOpts
}

func buildHttpResponseHeaderOpts(rawHttpResponseHeader []interface{}) *[]model.HttpResponseHeader {
	if len(rawHttpResponseHeader) < 1 {
		return nil
	}

	httpResponseHeaderOpts := make([]model.HttpResponseHeader, len(rawHttpResponseHeader))
	for i, v := range rawHttpResponseHeader {
		header := v.(map[string]interface{})
		httpResponseHeaderOpts[i] = model.HttpResponseHeader{
			Name:   header["name"].(string),
			Value:  utils.StringIgnoreEmpty(header["value"].(string)),
			Action: header["action"].(string),
		}
	}

	return &httpResponseHeaderOpts
}

func parseFunctionEnabledStatus(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}

func buildUrlAuthOpts(rawUrlAuth []interface{}) *model.UrlAuth {
	if len(rawUrlAuth) != 1 {
		return nil
	}

	urlAuth := rawUrlAuth[0].(map[string]interface{})
	urlAuthOpts := model.UrlAuth{
		Status:     parseFunctionEnabledStatus(urlAuth["enabled"].(bool)),
		Type:       utils.StringIgnoreEmpty(urlAuth["type"].(string)),
		Key:        utils.StringIgnoreEmpty(urlAuth["key"].(string)),
		TimeFormat: utils.StringIgnoreEmpty(urlAuth["time_format"].(string)),
		ExpireTime: utils.Int32IgnoreEmpty(int32(urlAuth["expire_time"].(int))),
	}

	return &urlAuthOpts
}

func buildForceRedirectOpts(rawForceRedirect []interface{}) *model.ForceRedirectConfig {
	if len(rawForceRedirect) != 1 {
		return nil
	}

	forceRedirect := rawForceRedirect[0].(map[string]interface{})
	forceRedirectOpts := model.ForceRedirectConfig{
		Status: parseFunctionEnabledStatus(forceRedirect["enabled"].(bool)),
		Type:   utils.StringIgnoreEmpty(forceRedirect["type"].(string)),
	}

	return &forceRedirectOpts
}

func buildCompressOpts(rawCompress []interface{}) *model.Compress {
	if len(rawCompress) != 1 {
		return nil
	}

	compress := rawCompress[0].(map[string]interface{})
	compressOpts := model.Compress{
		Status: parseFunctionEnabledStatus(compress["enabled"].(bool)),
		Type:   utils.StringIgnoreEmpty(compress["type"].(string)),
	}

	return &compressOpts
}

func buildCacheUrlParameterFilterOpts(rawCacheUrlParameterFilter []interface{}) *model.CacheUrlParameterFilter {
	if len(rawCacheUrlParameterFilter) != 1 {
		return nil
	}

	cacheUrlParameterFilter := rawCacheUrlParameterFilter[0].(map[string]interface{})
	cacheUrlParameterFilterOpts := model.CacheUrlParameterFilter{
		Value: utils.StringIgnoreEmpty(cacheUrlParameterFilter["value"].(string)),
		Type:  utils.StringIgnoreEmpty(cacheUrlParameterFilter["type"].(string)),
	}

	return &cacheUrlParameterFilterOpts
}

func buildIpFrequencyLimitOpts(rawIpFrequencyLimit []interface{}) *model.IpFrequencyLimit {
	if len(rawIpFrequencyLimit) != 1 {
		return nil
	}

	ipFrequencyLimit := rawIpFrequencyLimit[0].(map[string]interface{})
	ipFrequencyLimitOpts := model.IpFrequencyLimit{
		Status: parseFunctionEnabledStatus(ipFrequencyLimit["enabled"].(bool)),
		Qps:    utils.Int32IgnoreEmpty(int32(ipFrequencyLimit["qps"].(int))),
	}

	return &ipFrequencyLimitOpts
}

func buildWebsocketOpts(rawWebsocket []interface{}) *model.WebSocketSeek {
	if len(rawWebsocket) != 1 {
		return nil
	}

	websocket := rawWebsocket[0].(map[string]interface{})
	websocketOpts := model.WebSocketSeek{
		Status:  parseFunctionEnabledStatus(websocket["enabled"].(bool)),
		Timeout: int32(websocket["timeout"].(int)),
	}

	return &websocketOpts
}

func buildFlexibleOriginOpts(rawFlexibleOrigins []interface{}) *[]model.FlexibleOrigins {
	if len(rawFlexibleOrigins) < 1 {
		// Define an empty array to clear all flexible origins
		rst := make([]model.FlexibleOrigins, 0)
		return &rst
	}

	flexibleOriginOpts := make([]model.FlexibleOrigins, len(rawFlexibleOrigins))
	for i, v := range rawFlexibleOrigins {
		originMap := v.(map[string]interface{})
		flexibleOriginOpt := model.FlexibleOrigins{
			MatchType:    originMap["match_type"].(string),
			MatchPattern: originMap["match_pattern"].(string),
			Priority:     int32(originMap["priority"].(int)),
			BackSources:  buildFlexibleOriginBackSourceOpts(originMap["back_sources"].([]interface{})),
		}
		flexibleOriginOpts[i] = flexibleOriginOpt
	}
	return &flexibleOriginOpts
}

func buildFlexibleOriginBackSourceOpts(rawBackSources []interface{}) []model.BackSources {
	if len(rawBackSources) != 1 {
		return nil
	}

	backSource := rawBackSources[0].(map[string]interface{})
	backSourceOpts := model.BackSources{
		SourcesType:   backSource["sources_type"].(string),
		IpOrDomain:    backSource["ip_or_domain"].(string),
		ObsBucketType: utils.StringIgnoreEmpty(backSource["obs_bucket_type"].(string)),
		HttpPort:      utils.Int32IgnoreEmpty(int32(backSource["http_port"].(int))),
		HttpsPort:     utils.Int32IgnoreEmpty(int32(backSource["https_port"].(int))),
	}
	return []model.BackSources{backSourceOpts}
}

func buildRemoteAuthOpts(rawRemoteAuth []interface{}) *model.CommonRemoteAuth {
	if len(rawRemoteAuth) != 1 {
		return nil
	}

	remoteAuth := rawRemoteAuth[0].(map[string]interface{})
	remoteAuthOpts := model.CommonRemoteAuth{
		RemoteAuthentication: parseFunctionEnabledStatus(remoteAuth["enabled"].(bool)),
		RemoteAuthRules:      buildRemoteAuthRulesOpts(remoteAuth["remote_auth_rules"].([]interface{})),
	}
	return &remoteAuthOpts
}

func buildRemoteAuthRulesOpts(rawRemoteAuthRules []interface{}) *model.RemoteAuthRuleVo {
	if len(rawRemoteAuthRules) != 1 {
		return nil
	}

	remoteAuthRule := rawRemoteAuthRules[0].(map[string]interface{})
	remoteAuthRuleOpts := model.RemoteAuthRuleVo{
		AuthServer:            remoteAuthRule["auth_server"].(string),
		RequestMethod:         remoteAuthRule["request_method"].(string),
		FileTypeSetting:       remoteAuthRule["file_type_setting"].(string),
		SpecifiedFileType:     utils.StringIgnoreEmpty(remoteAuthRule["specified_file_type"].(string)),
		ReserveArgsSetting:    remoteAuthRule["reserve_args_setting"].(string),
		ReserveArgs:           utils.StringIgnoreEmpty(remoteAuthRule["reserve_args"].(string)),
		AddCustomArgsRules:    buildCustomArgsOpts(remoteAuthRule["add_custom_args_rules"].(*schema.Set).List()),
		ReserveHeadersSetting: remoteAuthRule["reserve_headers_setting"].(string),
		AddCustomHeadersRules: buildCustomArgsOpts(remoteAuthRule["add_custom_headers_rules"].(*schema.Set).List()),
		AuthSuccessStatus:     remoteAuthRule["auth_success_status"].(string),
		AuthFailedStatus:      remoteAuthRule["auth_failed_status"].(string),
		ResponseStatus:        remoteAuthRule["response_status"].(string),
		Timeout:               int32(remoteAuthRule["timeout"].(int)),
		TimeoutAction:         remoteAuthRule["timeout_action"].(string),
		ReserveHeaders:        utils.StringIgnoreEmpty(remoteAuthRule["reserve_headers"].(string)),
	}
	return &remoteAuthRuleOpts
}

func buildCustomArgsOpts(rawCustomArgs []interface{}) *[]model.CustomArgs {
	if len(rawCustomArgs) < 1 {
		// Define an empty array to clear all custom args
		rst := make([]model.CustomArgs, 0)
		return &rst
	}

	customArgsOpts := make([]model.CustomArgs, len(rawCustomArgs))
	for i, v := range rawCustomArgs {
		argMap := v.(map[string]interface{})
		customArgsOpt := model.CustomArgs{
			Type:  argMap["type"].(string),
			Key:   argMap["key"].(string),
			Value: argMap["value"].(string),
		}
		customArgsOpts[i] = customArgsOpt
	}
	return &customArgsOpts
}

func buildSourcesOpts(rawSources []interface{}) *[]model.SourcesConfig {
	if len(rawSources) < 1 {
		return nil
	}
	sourcesOpts := make([]model.SourcesConfig, len(rawSources))
	for i, v := range rawSources {
		source := v.(map[string]interface{})
		var priority int32
		if source["active"].(int) == 1 {
			priority = 70
		} else {
			priority = 30
		}
		sourcesOpts[i] = model.SourcesConfig{
			OriginAddr:          source["origin"].(string),
			OriginType:          source["origin_type"].(string),
			Priority:            priority,
			ObsWebHostingStatus: utils.String(parseFunctionEnabledStatus(source["obs_web_hosting_enabled"].(bool))),
			HttpPort:            utils.Int32IgnoreEmpty(int32(source["http_port"].(int))),
			HttpsPort:           utils.Int32IgnoreEmpty(int32(source["https_port"].(int))),
			HostName:            utils.StringIgnoreEmpty(source["retrieval_host"].(string)),
		}
	}
	return &sourcesOpts
}

func parseCacheRuleType(ruleType string) string {
	var cacheRuleTypes = map[string]string{
		"0": "all",
		"1": "file_extension",
		"2": "catalog",
		"3": "full_path",
		"5": "home_page",
	}
	if val, ok := cacheRuleTypes[ruleType]; ok {
		return val
	}
	return ruleType
}

func parseCacheTTLUnits(ttlUnit string) string {
	var cacheTTLUnits = map[string]string{
		"1": "s",
		"2": "m",
		"3": "h",
		"4": "d",
	}
	if val, ok := cacheTTLUnits[ttlUnit]; ok {
		return val
	}
	return ttlUnit
}

func buildCacheRules(followOrigin bool, rules []interface{}) *[]model.CacheRules {
	result := make([]model.CacheRules, len(rules))
	for i, val := range rules {
		rule := val.(map[string]interface{})
		result[i] = model.CacheRules{
			FollowOrigin: utils.StringIgnoreEmpty(parseFunctionEnabledStatus(followOrigin)),
			MatchType:    utils.StringIgnoreEmpty(parseCacheRuleType(rule["rule_type"].(string))),
			MatchValue:   utils.StringIgnoreEmpty(rule["content"].(string)),
			Ttl:          utils.Int32IgnoreEmpty(int32(rule["ttl"].(int))),
			TtlUnit:      parseCacheTTLUnits(rule["ttl_type"].(string)),
			Priority:     int32(rule["priority"].(int)),
		}
	}
	return &result
}

func updateDomainFullConfigs(client *cdnv2.CdnClient, cfg *config.Config, d *schema.ResourceData) error {
	rawConfigs := d.Get("configs").([]interface{})
	if len(rawConfigs) < 1 || rawConfigs[0] == nil {
		return nil
	}
	configs := rawConfigs[0].(map[string]interface{})

	ipv6Accelerate := 0
	if configs["ipv6_enable"].(bool) {
		ipv6Accelerate = 1
	}
	configsOpts := model.Configs{
		Sources:           buildSourcesOpts(d.Get("sources").(*schema.Set).List()),
		Ipv6Accelerate:    utils.Int32(int32(ipv6Accelerate)),
		OriginRangeStatus: utils.String(parseFunctionEnabledStatus(configs["range_based_retrieval_enabled"].(bool))),
	}
	if d.HasChange("configs.0.https_settings") {
		configsOpts.Https = buildHTTPSOpts(configs["https_settings"].([]interface{}))
	}
	if d.HasChange("configs.0.retrieval_request_header") {
		configsOpts.OriginRequestHeader = buildOriginRequestHeaderOpts(configs["retrieval_request_header"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.http_response_header") {
		configsOpts.HttpResponseHeader = buildHttpResponseHeaderOpts(configs["http_response_header"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.url_signing") {
		configsOpts.UrlAuth = buildUrlAuthOpts(configs["url_signing"].([]interface{}))
	}
	if d.HasChange("configs.0.origin_protocol") {
		configsOpts.OriginProtocol = utils.StringIgnoreEmpty(configs["origin_protocol"].(string))
	}
	if d.HasChange("configs.0.force_redirect") {
		configsOpts.ForceRedirect = buildForceRedirectOpts(configs["force_redirect"].([]interface{}))
	}
	if d.HasChange("configs.0.compress") {
		configsOpts.Compress = buildCompressOpts(configs["compress"].([]interface{}))
	}
	if d.HasChange("configs.0.cache_url_parameter_filter") {
		configsOpts.CacheUrlParameterFilter = buildCacheUrlParameterFilterOpts(configs["cache_url_parameter_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.ip_frequency_limit") {
		configsOpts.IpFrequencyLimit = buildIpFrequencyLimitOpts(configs["ip_frequency_limit"].([]interface{}))
	}
	if d.HasChange("configs.0.websocket") {
		configsOpts.Websocket = buildWebsocketOpts(configs["websocket"].([]interface{}))
	}
	if d.HasChange("configs.0.flexible_origin") {
		configsOpts.FlexibleOrigin = buildFlexibleOriginOpts(configs["flexible_origin"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.remote_auth") {
		configsOpts.RemoteAuth = buildRemoteAuthOpts(configs["remote_auth"].([]interface{}))
	}

	if d.HasChange("cache_settings") {
		cacheSettings := d.Get("cache_settings").([]interface{})
		if len(cacheSettings) > 0 {
			cacheSetting := cacheSettings[0].(map[string]interface{})
			configsOpts.CacheRules = buildCacheRules(cacheSetting["follow_origin"].(bool), cacheSetting["rules"].(*schema.Set).List())
		}
	}

	req := model.UpdateDomainFullConfigRequest{
		DomainName:          d.Get("name").(string),
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		Body: &model.ModifyDomainConfigRequestBody{
			Configs: &configsOpts,
		},
	}

	_, err := client.UpdateDomainFullConfig(&req)
	if err != nil {
		return err
	}
	return nil
}

func resourceCdnDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	createOpts := &domains.CreateOpts{
		DomainName:          d.Get("name").(string),
		BusinessType:        d.Get("type").(string),
		Sources:             buildCreateDomainSources(d),
		ServiceArea:         d.Get("service_area").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	v, err := domains.Create(cdnClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating CDN domain: %s", err)
	}

	if v.ID == "" {
		return diag.Errorf("error creating CDN domain: ID is not found in API response")
	}
	d.SetId(v.ID)

	hcCdnClient, err := cfg.HcCdnV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}
	requestOpts := buildDomainDetailRequestOpts(d, cfg)
	if err := waitingForStatusOnline(ctx, hcCdnClient, d.Timeout(schema.TimeoutCreate), requestOpts); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) creation to become online: %s", d.Id(), err)
	}
	return resourceCdnDomainUpdate(ctx, d, meta)
}

func waitingForStatusOnline(ctx context.Context, hcCdnClient *cdnv2.CdnClient, timeout time.Duration,
	opts *model.ShowDomainDetailByNameRequest) error {
	unexpectedStatus := []string{"offline", "configure_failed", "check_failed", "deleting"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domain, err := hcCdnClient.ShowDomainDetailByName(opts)
			if err != nil {
				return nil, "ERROR", err
			}

			if domain == nil || domain.Domain == nil {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: Domain is not found in API response")
			}

			status := utils.StringValue(domain.Domain.DomainStatus)
			if status == "online" {
				return domain, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return domain, status, nil
			}
			return domain, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForStatusOffline(ctx context.Context, hcCdnClient *cdnv2.CdnClient, timeout time.Duration,
	opts *model.ShowDomainDetailByNameRequest) error {
	unexpectedStatus := []string{"online", "configure_failed", "check_failed", "deleting"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domain, err := hcCdnClient.ShowDomainDetailByName(opts)
			if err != nil {
				return nil, "ERROR", err
			}

			if domain == nil || domain.Domain == nil {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: Domain is not found in API response")
			}

			status := utils.StringValue(domain.Domain.DomainStatus)
			if status == "offline" {
				return domain, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return domain, status, nil
			}
			return domain, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func analyseFunctionEnabledStatus(enabledStatus string) bool {
	return enabledStatus == "on"
}

func analyseFunctionEnabledStatusPtr(enabledStatus *string) bool {
	return enabledStatus != nil && *enabledStatus == "on"
}

func flattenHTTPSAttrs(https *model.HttpGetBody, privateKey string) []map[string]interface{} {
	if https == nil {
		return nil
	}
	httpsAttrs := map[string]interface{}{
		"https_status":       https.HttpsStatus,
		"certificate_name":   https.CertificateName,
		"certificate_body":   https.CertificateValue,
		"private_key":        privateKey,
		"certificate_source": https.CertificateSource,
		"http2_status":       https.Http2Status,
		"tls_version":        https.TlsVersion,
		"https_enabled":      analyseFunctionEnabledStatusPtr(https.HttpsStatus),
		"http2_enabled":      analyseFunctionEnabledStatusPtr(https.Http2Status),
	}

	return []map[string]interface{}{httpsAttrs}
}

func flattenOriginRequestHeaderAttrs(originRequestHeader *[]model.OriginRequestHeader) []map[string]interface{} {
	if originRequestHeader == nil || len(*originRequestHeader) == 0 {
		return nil
	}

	originRequestHeaderAttrs := make([]map[string]interface{}, len(*originRequestHeader))
	for i, v := range *originRequestHeader {
		originRequestHeaderAttrs[i] = map[string]interface{}{
			"name":   v.Name,
			"value":  v.Value,
			"action": v.Action,
		}
	}

	return originRequestHeaderAttrs
}

func flattenHttpResponseHeaderAttrs(httpResponseHeader *[]model.HttpResponseHeader) []map[string]interface{} {
	if httpResponseHeader == nil || len(*httpResponseHeader) == 0 {
		return nil
	}

	httpResponseHeaderAttrs := make([]map[string]interface{}, len(*httpResponseHeader))
	for i, v := range *httpResponseHeader {
		httpResponseHeaderAttrs[i] = map[string]interface{}{
			"name":   v.Name,
			"value":  v.Value,
			"action": v.Action,
		}
	}

	return httpResponseHeaderAttrs
}

func flattenUrlAuthAttrs(urlAuth *model.UrlAuthGetBody, urlAuthKey string) []map[string]interface{} {
	if urlAuth == nil {
		return nil
	}

	urlAuthAttrs := map[string]interface{}{
		"enabled":     analyseFunctionEnabledStatus(urlAuth.Status),
		"status":      urlAuth.Status,
		"type":        urlAuth.Type,
		"key":         urlAuthKey,
		"time_format": urlAuth.TimeFormat,
		"expire_time": urlAuth.ExpireTime,
	}

	return []map[string]interface{}{urlAuthAttrs}
}

func flattenForceRedirectAttrs(forceRedirect *model.ForceRedirectConfig) []map[string]interface{} {
	if forceRedirect == nil {
		return nil
	}

	forceRedirectAttrs := map[string]interface{}{
		"status":  forceRedirect.Status,
		"type":    forceRedirect.Type,
		"enabled": analyseFunctionEnabledStatus(forceRedirect.Status),
	}

	return []map[string]interface{}{forceRedirectAttrs}
}

func flattenCompressAttrs(compress *model.Compress) []map[string]interface{} {
	if compress == nil {
		return nil
	}

	compressAttrs := map[string]interface{}{
		"status":  compress.Status,
		"type":    compress.Type,
		"enabled": analyseFunctionEnabledStatus(compress.Status),
	}

	return []map[string]interface{}{compressAttrs}
}

func flattenCacheUrlParameterFilterAttrs(cacheUrlParameterFilter *model.CacheUrlParameterFilterGetBody) []map[string]interface{} {
	if cacheUrlParameterFilter == nil {
		return nil
	}

	cacheUrlParameterFilterAttrs := map[string]interface{}{
		"value": cacheUrlParameterFilter.Value,
		"type":  cacheUrlParameterFilter.Type,
	}

	return []map[string]interface{}{cacheUrlParameterFilterAttrs}
}

func flattenIpFrequencyLimitAttrs(ipFrequencyLimit *model.IpFrequencyLimitQuery) []map[string]interface{} {
	if ipFrequencyLimit == nil {
		return nil
	}

	ipFrequencyLimitAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(ipFrequencyLimit.Status),
		"qps":     ipFrequencyLimit.Qps,
	}

	return []map[string]interface{}{ipFrequencyLimitAttrs}
}

func flattenWebsocketAttrs(websocket *model.WebSocketSeek) []map[string]interface{} {
	if websocket == nil {
		return nil
	}

	websocketAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(websocket.Status),
		"timeout": websocket.Timeout,
	}

	return []map[string]interface{}{websocketAttrs}
}

func flattenFlexibleOriginAttrs(flexibleOrigins *[]model.FlexibleOrigins) []map[string]interface{} {
	if flexibleOrigins == nil || len(*flexibleOrigins) == 0 {
		return nil
	}

	flexibleOriginsAttrs := make([]map[string]interface{}, len(*flexibleOrigins))
	for i, v := range *flexibleOrigins {
		flexibleOriginsAttrs[i] = map[string]interface{}{
			"match_type":    v.MatchType,
			"match_pattern": v.MatchPattern,
			"priority":      v.Priority,
			"back_sources":  flattenFlexibleOriginBackSourceAttrs(v.BackSources),
		}
	}
	return flexibleOriginsAttrs
}

func flattenFlexibleOriginBackSourceAttrs(backSources []model.BackSources) []map[string]interface{} {
	if len(backSources) == 0 {
		return nil
	}

	backSourcesAttrs := make([]map[string]interface{}, len(backSources))
	for i, v := range backSources {
		backSourcesAttrs[i] = map[string]interface{}{
			"sources_type":    v.SourcesType,
			"ip_or_domain":    v.IpOrDomain,
			"obs_bucket_type": v.ObsBucketType,
			"http_port":       v.HttpPort,
			"https_port":      v.HttpsPort,
		}
	}
	return backSourcesAttrs
}

func flattenRemoteAuthAttrs(remoteAuth *model.CommonRemoteAuth) []map[string]interface{} {
	if remoteAuth == nil {
		return nil
	}

	remoteAuthAttrs := map[string]interface{}{
		"enabled":           analyseFunctionEnabledStatus(remoteAuth.RemoteAuthentication),
		"remote_auth_rules": flattenRemoteAuthRulesAttrs(remoteAuth.RemoteAuthRules),
	}
	return []map[string]interface{}{remoteAuthAttrs}
}

func flattenRemoteAuthRulesAttrs(remoteAuthRule *model.RemoteAuthRuleVo) []map[string]interface{} {
	if remoteAuthRule == nil {
		return nil
	}

	remoteAuthRuleAttrs := map[string]interface{}{
		"auth_server":              remoteAuthRule.AuthServer,
		"request_method":           remoteAuthRule.RequestMethod,
		"file_type_setting":        remoteAuthRule.FileTypeSetting,
		"specified_file_type":      remoteAuthRule.SpecifiedFileType,
		"reserve_args_setting":     remoteAuthRule.ReserveArgsSetting,
		"reserve_args":             remoteAuthRule.ReserveArgs,
		"add_custom_args_rules":    flattenCustomArgsAttrs(remoteAuthRule.AddCustomArgsRules),
		"reserve_headers_setting":  remoteAuthRule.ReserveHeadersSetting,
		"add_custom_headers_rules": flattenCustomArgsAttrs(remoteAuthRule.AddCustomHeadersRules),
		"auth_success_status":      remoteAuthRule.AuthSuccessStatus,
		"auth_failed_status":       remoteAuthRule.AuthFailedStatus,
		"response_status":          remoteAuthRule.ResponseStatus,
		"timeout":                  remoteAuthRule.Timeout,
		"timeout_action":           remoteAuthRule.TimeoutAction,
		"reserve_headers":          remoteAuthRule.ReserveHeaders,
	}
	return []map[string]interface{}{remoteAuthRuleAttrs}
}

func flattenCustomArgsAttrs(customArgs *[]model.CustomArgs) []map[string]interface{} {
	if customArgs == nil || len(*customArgs) == 0 {
		return nil
	}

	customArgsAttrs := make([]map[string]interface{}, len(*customArgs))
	for i, v := range *customArgs {
		customArgsAttrs[i] = map[string]interface{}{
			"type":  v.Type,
			"key":   v.Key,
			"value": v.Value,
		}
	}
	return customArgsAttrs
}

func flattenSourcesAttrs(sources *[]model.SourcesConfig) []map[string]interface{} {
	if sources == nil || len(*sources) == 0 {
		return nil
	}

	sourcesAttrs := make([]map[string]interface{}, len(*sources))
	for i, v := range *sources {
		var active int
		if v.Priority == 70 {
			active = 1
		}
		sourcesAttrs[i] = map[string]interface{}{
			"origin":                  v.OriginAddr,
			"origin_type":             v.OriginType,
			"active":                  active,
			"obs_web_hosting_enabled": analyseFunctionEnabledStatusPtr(v.ObsWebHostingStatus),
			"http_port":               v.HttpPort,
			"https_port":              v.HttpsPort,
			"retrieval_host":          v.HostName,
		}
	}

	return sourcesAttrs
}

func flattenCacheRulesAttrs(cacheRulesPtr *[]model.CacheRules) []map[string]interface{} {
	if cacheRulesPtr == nil || len(*cacheRulesPtr) == 0 {
		return nil
	}

	cacheRules := *cacheRulesPtr
	sourcesAttrs := make([]map[string]interface{}, len(cacheRules))
	for i, v := range cacheRules {
		sourcesAttrs[i] = map[string]interface{}{
			"rule_type": v.MatchType,
			"content":   v.MatchValue,
			"ttl":       v.Ttl,
			"ttl_type":  v.TtlUnit,
			"priority":  v.Priority,
		}
	}

	return []map[string]interface{}{
		{
			"follow_origin": analyseFunctionEnabledStatus(utils.StringValue(cacheRules[0].FollowOrigin)),
			"rules":         sourcesAttrs,
		},
	}
}

func flattenConfigAttrs(configsResp *model.ConfigsGetBody, d *schema.ResourceData) []map[string]interface{} {
	privateKey := d.Get("configs.0.https_settings.0.private_key").(string)
	urlAuthKey := d.Get("configs.0.url_signing.0.key").(string)

	configsAttrs := map[string]interface{}{
		"https_settings":                flattenHTTPSAttrs(configsResp.Https, privateKey),
		"retrieval_request_header":      flattenOriginRequestHeaderAttrs(configsResp.OriginRequestHeader),
		"http_response_header":          flattenHttpResponseHeaderAttrs(configsResp.HttpResponseHeader),
		"url_signing":                   flattenUrlAuthAttrs(configsResp.UrlAuth, urlAuthKey),
		"origin_protocol":               configsResp.OriginProtocol,
		"force_redirect":                flattenForceRedirectAttrs(configsResp.ForceRedirect),
		"compress":                      flattenCompressAttrs(configsResp.Compress),
		"cache_url_parameter_filter":    flattenCacheUrlParameterFilterAttrs(configsResp.CacheUrlParameterFilter),
		"ip_frequency_limit":            flattenIpFrequencyLimitAttrs(configsResp.IpFrequencyLimit),
		"websocket":                     flattenWebsocketAttrs(configsResp.Websocket),
		"flexible_origin":               flattenFlexibleOriginAttrs(configsResp.FlexibleOrigin),
		"remote_auth":                   flattenRemoteAuthAttrs(configsResp.RemoteAuth),
		"ipv6_enable":                   configsResp.Ipv6Accelerate != nil && *configsResp.Ipv6Accelerate == 1,
		"range_based_retrieval_enabled": analyseFunctionEnabledStatusPtr(configsResp.OriginRangeStatus),
	}
	return []map[string]interface{}{configsAttrs}
}

func queryDomainFullConfig(hcCdnClient *cdnv2.CdnClient, cfg *config.Config, d *schema.ResourceData,
	domainName string) (*model.ConfigsGetBody, error) {
	req := model.ShowDomainFullConfigRequest{
		DomainName:          domainName,
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	resp, err := hcCdnClient.ShowDomainFullConfig(&req)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CDN domain full config: %s", err)
	}

	if resp == nil || resp.Configs == nil {
		return nil, fmt.Errorf("error retrieving CDN domain full config: Config is not found in API response")
	}
	return resp.Configs, nil
}

func queryAndFlattenDomainTags(hcCdnClient *cdnv2.CdnClient, d *schema.ResourceData) (map[string]string, error) {
	tags, err := hcCdnClient.ShowTags(&model.ShowTagsRequest{ResourceId: d.Id()})
	if err != nil {
		return nil, fmt.Errorf("error retrieving CDN domain tags: %s", err)
	}

	if tags.Tags == nil || len(*tags.Tags) == 0 {
		return nil, nil
	}

	tagMap := make(map[string]string, len(*tags.Tags))
	for _, tag := range *tags.Tags {
		if tag.Value != nil {
			tagMap[tag.Key] = *tag.Value
		} else {
			tagMap[tag.Key] = ""
		}
	}
	return tagMap, nil
}

func flattenServiceArea(serviceArea *model.DomainsDetailServiceArea) interface{} {
	if serviceArea == nil {
		return nil
	}
	return serviceArea.Value()
}

func resourceCdnDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	hcCdnClient, err := cfg.HcCdnV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	requestOpts := buildDomainDetailRequestOpts(d, cfg)
	v, err := hcCdnClient.ShowDomainDetailByName(requestOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, parseDetailResponseError(err), "error retrieving CDN domain")
	}

	if v == nil || v.Domain == nil {
		return diag.Errorf("error retrieving CDN domain: Domain is not found in API response")
	}

	domain := *v.Domain
	// Backfield the id when executing the import operation
	d.SetId(*domain.Id)

	configsResp, err := queryDomainFullConfig(hcCdnClient, cfg, d, *domain.DomainName)
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := queryAndFlattenDomainTags(hcCdnClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", domain.DomainName),
		d.Set("type", domain.BusinessType),
		d.Set("cname", domain.Cname),
		d.Set("domain_status", domain.DomainStatus),
		d.Set("service_area", flattenServiceArea(domain.ServiceArea)),
		d.Set("sources", flattenSourcesAttrs(configsResp.Sources)),
		d.Set("configs", flattenConfigAttrs(configsResp, d)),
		d.Set("cache_settings", flattenCacheRulesAttrs(configsResp.CacheRules)),
		d.Set("tags", tags),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// When the domain name does not exist, the response body example of the details interface is as follows:
// {"error": {"error_code": "CDN.0170","error_msg": "domain not exist!"}}
func parseDetailResponseError(err error) error {
	var responseErr *sdkerr.ServiceResponseError
	if errors.As(err, &responseErr) {
		if responseErr.StatusCode == http.StatusBadRequest && responseErr.ErrorCode == "CDN.0170" {
			return golangsdk.ErrDefault404{}
		}
	}
	return err
}

func resourceCdnDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	hcCdnClient, err := cfg.HcCdnV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	if d.HasChanges("sources", "configs", "cache_settings") || d.IsNewResource() {
		err = updateDomainFullConfigs(hcCdnClient, cfg, d)
		if err != nil {
			return diag.Errorf("error updating CDN domain configs settings: %s", err)
		}

		requestOpts := buildDomainDetailRequestOpts(d, cfg)
		if err := waitingForStatusOnline(ctx, hcCdnClient, d.Timeout(schema.TimeoutUpdate), requestOpts); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become online: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err := updateDomainTags(hcCdnClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCdnDomainRead(ctx, d, meta)
}

func updateDomainTags(hcCdnClient *cdnv2.CdnClient, d *schema.ResourceData) error {
	oTagsRaw, nTagsRaw := d.GetChange("tags")
	oTagsMap := oTagsRaw.(map[string]interface{})
	nTagsMap := nTagsRaw.(map[string]interface{})

	if len(oTagsMap) > 0 {
		var tagList []string
		for k := range oTagsMap {
			tagList = append(tagList, k)
		}
		deleteTagsReq := model.BatchDeleteTagsRequest{
			Body: &model.DeleteTagsRequestBody{
				ResourceId: d.Id(),
				Tags:       tagList,
			},
		}
		_, err := hcCdnClient.BatchDeleteTags(&deleteTagsReq)
		if err != nil {
			return fmt.Errorf("error deleting CDN domain tags: %s", err)
		}
	}

	if len(nTagsMap) > 0 {
		tagList := make([]model.TagMap, 0, len(nTagsMap))
		for k, v := range nTagsMap {
			tag := model.TagMap{
				Key:   k,
				Value: utils.String(v.(string)),
			}
			tagList = append(tagList, tag)
		}
		createTagsReq := model.CreateTagsRequest{
			Body: &model.CreateTagsRequestBody{
				ResourceId: d.Id(),
				Tags:       tagList,
			},
		}
		_, err := hcCdnClient.CreateTags(&createTagsReq)
		if err != nil {
			return fmt.Errorf("error creating CDN domain tags: %s", err)
		}
	}
	return nil
}

func resourceCdnDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	opts := buildResourceExtensionOpts(d, cfg)
	if d.Get("domain_status").(string) == "online" {
		// make sure the status has changed to offline before deleting it.
		if err = domains.Disable(cdnClient, d.Id(), opts).Err; err != nil {
			return diag.Errorf("error disable CDN domain %s: %s", d.Id(), err)
		}

		hcCdnClient, err := cfg.HcCdnV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating CDN v2 client: %s", err)
		}
		requestOpts := buildDomainDetailRequestOpts(d, cfg)
		if err := waitingForStatusOffline(ctx, hcCdnClient, d.Timeout(schema.TimeoutDelete), requestOpts); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become offline: %s", d.Id(), err)
		}
	}

	_, err = domains.Delete(cdnClient, d.Id(), opts).Extract()
	if err != nil {
		return diag.Errorf("error deleting CDN domain (%s): %s", d.Id(), err)
	}

	// an API issue will be raised in ForceNew scene, so wait for a while
	time.Sleep(3 * time.Second) // lintignore:R018
	return nil
}

func buildResourceExtensionOpts(d *schema.ResourceData, cfg *config.Config) *domains.ExtensionOpts {
	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		return &domains.ExtensionOpts{
			EnterpriseProjectId: epsID,
		}
	}

	return nil
}

func buildDomainDetailRequestOpts(d *schema.ResourceData, cfg *config.Config) *model.ShowDomainDetailByNameRequest {
	return &model.ShowDomainDetailByNameRequest{
		DomainName:          d.Get("name").(string),
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
}

func resourceCDNDomainImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("name", d.Id())
}
