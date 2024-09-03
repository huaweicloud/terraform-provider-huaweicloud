package cdn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"

	cdnv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{"name"}

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
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http2_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tls_version": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
				Computed:         true,
			},
			"ocsp_stapling_status": {
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
			"sign_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"inherit_config": {
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
						"inherit_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"inherit_time_type": {
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
			},
			"sign_arg": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"backup_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
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

var forceRedirect = schema.Schema{
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			// Cloud will configure this field to `302` by default
			"redirect_code": {
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

var compress = schema.Schema{
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
			"file_type": {
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
	Type:        schema.TypeList,
	Optional:    true,
	Computed:    true,
	MaxItems:    1,
	Description: "schema: Deprecated; Field `cache_url_parameter_filter` will be offline soon, use `cache_settings` instead",
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

var quic = schema.Schema{
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
		},
	},
}

var referer = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
				Computed:         true,
			},
			"include_empty": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var videoSeek = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_video_seek": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_flv_by_time_seek": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"start_parameter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_parameter": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var requestLimitRules = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"limit_rate_after": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_rate_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var errorCodeCache = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	},
}

var ipFilter = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var originRequestUrlRewrite = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var userAgentFilter = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Computed: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ua_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
		},
	},
}

var errorCodeRedirectRules = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_code": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_link": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	},
}

var hsts = schema.Schema{
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
			"max_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"include_subdomains": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var accessAreaFilters = schema.Schema{
	Type:        schema.TypeSet,
	Optional:    true,
	Description: "schema: Internal; Specifies the geographic access control rules.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "schema: Internal; Specifies the the blacklist and whitelist rule type.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "schema: Internal; Specifies the content type.",
			},
			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "schema: Internal; Specifies the areas, separated by commas.",
			},
			"content_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Internal; Specifies the content value.",
			},
			"exception_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Internal; Specifies the IP addresses exception in access control, separated by commas.",
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
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
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
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"obs_bucket_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// Cloud will configure this field to `on` by default
						"slice_etag_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						// Cloud will configure this field to `30` by default
						"origin_receive_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						// Cloud will configure this field to `off` by default
						"origin_follow302_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"https_settings":             &httpsConfig,
						"retrieval_request_header":   &requestAndResponseHeader,
						"http_response_header":       &requestAndResponseHeader,
						"url_signing":                &authOpts,
						"force_redirect":             &forceRedirect,
						"compress":                   &compress,
						"cache_url_parameter_filter": &cacheUrlParameterFilter,
						"ip_frequency_limit":         &ipFrequencyLimit,
						"websocket":                  &websocket,
						"flexible_origin":            &flexibleOrigin,
						"remote_auth":                &remoteAuth,
						"quic":                       &quic,
						"referer":                    &referer,
						"video_seek":                 &videoSeek,
						"request_limit_rules":        &requestLimitRules,
						"error_code_cache":           &errorCodeCache,
						"ip_filter":                  &ipFilter,
						"origin_request_url_rewrite": &originRequestUrlRewrite,
						"user_agent_filter":          &userAgentFilter,
						"error_code_redirect_rules":  &errorCodeRedirectRules,
						"hsts":                       &hsts,
						"access_area_filter":         &accessAreaFilters,
					},
				},
			},

			// The cloud will create a rule for `cache_settings` by default, so its value will not be set when querying.
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
									"url_parameter_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url_parameter_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "schema: Internal",
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

func buildHTTPSStatusOpts(enable bool) string {
	if enable {
		return "on"
	}
	return "off"
}

func buildHTTP2StatusOpts(enable bool) string {
	if enable {
		return "on"
	}
	// Currently, European sites do not support this parameter, so we will handle it this way for the time being.
	return ""
}

func buildHTTPSOpts(rawHTTPS []interface{}) *model.HttpPutBody {
	if len(rawHTTPS) != 1 {
		return nil
	}

	https := rawHTTPS[0].(map[string]interface{})
	httpsOpts := model.HttpPutBody{
		HttpsStatus:        utils.String(buildHTTPSStatusOpts(https["https_enabled"].(bool))),
		CertificateName:    utils.StringIgnoreEmpty(https["certificate_name"].(string)),
		CertificateValue:   utils.StringIgnoreEmpty(https["certificate_body"].(string)),
		PrivateKey:         utils.StringIgnoreEmpty(https["private_key"].(string)),
		CertificateSource:  utils.Int32(int32(https["certificate_source"].(int))),
		CertificateType:    utils.StringIgnoreEmpty(https["certificate_type"].(string)),
		Http2Status:        utils.StringIgnoreEmpty(buildHTTP2StatusOpts(https["http2_enabled"].(bool))),
		TlsVersion:         utils.StringIgnoreEmpty(https["tls_version"].(string)),
		OcspStaplingStatus: utils.StringIgnoreEmpty(https["ocsp_stapling_status"].(string)),
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
		Status:        parseFunctionEnabledStatus(urlAuth["enabled"].(bool)),
		Type:          utils.StringIgnoreEmpty(urlAuth["type"].(string)),
		SignMethod:    utils.StringIgnoreEmpty(urlAuth["sign_method"].(string)),
		MatchType:     utils.StringIgnoreEmpty(urlAuth["match_type"].(string)),
		InheritConfig: buildInheritConfigOpts(urlAuth["inherit_config"].([]interface{})),
		SignArg:       utils.StringIgnoreEmpty(urlAuth["sign_arg"].(string)),
		Key:           utils.StringIgnoreEmpty(urlAuth["key"].(string)),
		BackupKey:     utils.StringIgnoreEmpty(urlAuth["backup_key"].(string)),
		TimeFormat:    utils.StringIgnoreEmpty(urlAuth["time_format"].(string)),
		ExpireTime:    utils.Int32(int32(urlAuth["expire_time"].(int))),
	}

	return &urlAuthOpts
}

func buildInheritConfigOpts(rwaInheritConfig []interface{}) *model.InheritConfig {
	if len(rwaInheritConfig) != 1 {
		return nil
	}

	inheritConfig := rwaInheritConfig[0].(map[string]interface{})
	inheritConfigOpts := model.InheritConfig{
		Status:          parseFunctionEnabledStatus(inheritConfig["enabled"].(bool)),
		InheritType:     utils.StringIgnoreEmpty(inheritConfig["inherit_type"].(string)),
		InheritTimeType: utils.StringIgnoreEmpty(inheritConfig["inherit_time_type"].(string)),
	}

	return &inheritConfigOpts
}

func buildForceRedirectOpts(rawForceRedirect []interface{}) *model.ForceRedirectConfig {
	if len(rawForceRedirect) != 1 {
		return nil
	}

	forceRedirect := rawForceRedirect[0].(map[string]interface{})
	forceRedirectOpts := model.ForceRedirectConfig{
		Status:       parseFunctionEnabledStatus(forceRedirect["enabled"].(bool)),
		Type:         utils.StringIgnoreEmpty(forceRedirect["type"].(string)),
		RedirectCode: utils.Int32IgnoreEmpty(int32(forceRedirect["redirect_code"].(int))),
	}

	return &forceRedirectOpts
}

func buildCompressOpts(rawCompress []interface{}) *model.Compress {
	if len(rawCompress) != 1 {
		return nil
	}

	compress := rawCompress[0].(map[string]interface{})
	compressOpts := model.Compress{
		Status:   parseFunctionEnabledStatus(compress["enabled"].(bool)),
		Type:     utils.StringIgnoreEmpty(compress["type"].(string)),
		FileType: utils.StringIgnoreEmpty(compress["file_type"].(string)),
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

func buildRemoteAuthRulesOpts(rawRemoteAuthRules []interface{}) *model.RemoteAuthRule {
	if len(rawRemoteAuthRules) != 1 {
		return nil
	}

	remoteAuthRule := rawRemoteAuthRules[0].(map[string]interface{})
	remoteAuthRuleOpts := model.RemoteAuthRule{
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

func buildQUICOpts(rawQuic []interface{}) *model.Quic {
	if len(rawQuic) != 1 {
		return nil
	}

	quic := rawQuic[0].(map[string]interface{})
	quicOpts := model.Quic{
		Status: parseFunctionEnabledStatus(quic["enabled"].(bool)),
	}

	return &quicOpts
}

func buildRefererOpts(rawReferer []interface{}) *model.RefererConfig {
	if len(rawReferer) != 1 {
		return nil
	}

	referer := rawReferer[0].(map[string]interface{})
	refererOpts := model.RefererConfig{
		Type:         referer["type"].(string),
		Value:        utils.String(referer["value"].(string)),
		IncludeEmpty: utils.Bool(referer["include_empty"].(bool)),
	}

	return &refererOpts
}

func buildVideoSeekOpts(rawVideoSeek []interface{}) *model.VideoSeek {
	if len(rawVideoSeek) != 1 {
		return nil
	}

	videoSeek := rawVideoSeek[0].(map[string]interface{})
	videoSeekOpts := model.VideoSeek{
		EnableVideoSeek:     videoSeek["enable_video_seek"].(bool),
		EnableFlvByTimeSeek: utils.Bool(videoSeek["enable_flv_by_time_seek"].(bool)),
		StartParameter:      utils.String(videoSeek["start_parameter"].(string)),
		EndParameter:        utils.String(videoSeek["end_parameter"].(string)),
	}

	return &videoSeekOpts
}

func buildRequestLimitRulesOpts(rawRequestLimitRules []interface{}) *[]model.RequestLimitRules {
	if len(rawRequestLimitRules) < 1 {
		// Define an empty array to clear all request limit rules
		rst := make([]model.RequestLimitRules, 0)
		return &rst
	}

	requestLimitRulesOpts := make([]model.RequestLimitRules, len(rawRequestLimitRules))
	for i, v := range rawRequestLimitRules {
		ruleMap := v.(map[string]interface{})
		ruleOpt := model.RequestLimitRules{
			Priority:       int32(ruleMap["priority"].(int)),
			MatchType:      ruleMap["match_type"].(string),
			MatchValue:     utils.String(ruleMap["match_value"].(string)),
			Type:           ruleMap["type"].(string),
			LimitRateAfter: int64(ruleMap["limit_rate_after"].(int)),
			LimitRateValue: int32(ruleMap["limit_rate_value"].(int)),
		}
		requestLimitRulesOpts[i] = ruleOpt
	}
	return &requestLimitRulesOpts
}

func buildErrorCodeCacheOpts(rawErrorCodeCache []interface{}) *[]model.ErrorCodeCache {
	if len(rawErrorCodeCache) < 1 {
		// Define an empty array to clear all error code cache
		rst := make([]model.ErrorCodeCache, 0)
		return &rst
	}

	errorCodeCacheOpts := make([]model.ErrorCodeCache, len(rawErrorCodeCache))
	for i, v := range rawErrorCodeCache {
		cacheMap := v.(map[string]interface{})
		cacheOpt := model.ErrorCodeCache{
			Code: utils.Int32(int32(cacheMap["code"].(int))),
			Ttl:  utils.Int32(int32(cacheMap["ttl"].(int))),
		}
		errorCodeCacheOpts[i] = cacheOpt
	}
	return &errorCodeCacheOpts
}

func buildIpFilterOpts(rawIpFilter []interface{}) *model.IpFilter {
	if len(rawIpFilter) != 1 {
		return nil
	}

	ipFilter := rawIpFilter[0].(map[string]interface{})
	ipFilterOpts := model.IpFilter{
		Type:  ipFilter["type"].(string),
		Value: utils.String(ipFilter["value"].(string)),
	}

	return &ipFilterOpts
}

func buildOriginRequestUrlRewriteOpts(rawOriginRequestUrlRewrite []interface{}) *[]model.OriginRequestUrlRewrite {
	if len(rawOriginRequestUrlRewrite) < 1 {
		// Define an empty array to clear all origin request url rewrite
		rst := make([]model.OriginRequestUrlRewrite, 0)
		return &rst
	}

	originRequestUrlRewriteOpts := make([]model.OriginRequestUrlRewrite, len(rawOriginRequestUrlRewrite))
	for i, v := range rawOriginRequestUrlRewrite {
		urlMap := v.(map[string]interface{})
		urlOpt := model.OriginRequestUrlRewrite{
			Priority:  int32(urlMap["priority"].(int)),
			MatchType: urlMap["match_type"].(string),
			TargetUrl: urlMap["target_url"].(string),
			SourceUrl: utils.StringIgnoreEmpty(urlMap["source_url"].(string)),
		}
		originRequestUrlRewriteOpts[i] = urlOpt
	}
	return &originRequestUrlRewriteOpts
}

func buildUserAgentFilterOpts(rawUserAgentFilter []interface{}) *model.UserAgentFilter {
	if len(rawUserAgentFilter) != 1 {
		return nil
	}

	userAgentFilter := rawUserAgentFilter[0].(map[string]interface{})
	userAgentFilterOpts := model.UserAgentFilter{
		Type:   userAgentFilter["type"].(string),
		UaList: utils.ExpandToStringListPointer(userAgentFilter["ua_list"].(*schema.Set).List()),
	}

	return &userAgentFilterOpts
}

func buildErrorCodeRedirectRules(errorCodeRedirectRules []interface{}) *[]model.ErrorCodeRedirectRules {
	if len(errorCodeRedirectRules) < 1 {
		// Define an empty array to clear all error code redirect rules
		rst := make([]model.ErrorCodeRedirectRules, 0)
		return &rst
	}

	errorCodeRedirectRulesOpts := make([]model.ErrorCodeRedirectRules, len(errorCodeRedirectRules))
	for i, v := range errorCodeRedirectRules {
		ruleMap := v.(map[string]interface{})
		ruleOpt := model.ErrorCodeRedirectRules{
			ErrorCode:  int32(ruleMap["error_code"].(int)),
			TargetCode: int32(ruleMap["target_code"].(int)),
			TargetLink: ruleMap["target_link"].(string),
		}
		errorCodeRedirectRulesOpts[i] = ruleOpt
	}
	return &errorCodeRedirectRulesOpts
}

func buildHstsOpts(rawHsts []interface{}) *model.Hsts {
	if len(rawHsts) != 1 {
		return nil
	}

	hsts := rawHsts[0].(map[string]interface{})
	hstsOpts := model.Hsts{
		Status:            parseFunctionEnabledStatus(hsts["enabled"].(bool)),
		MaxAge:            utils.Int32(int32(hsts["max_age"].(int))),
		IncludeSubdomains: utils.StringIgnoreEmpty(hsts["include_subdomains"].(string)),
	}

	return &hstsOpts
}

func buildAccessAreaFilters(accessAreaFilters []interface{}) *[]model.AccessAreaFilter {
	if len(accessAreaFilters) < 1 {
		// Define an empty array to clear all access area filters
		rst := make([]model.AccessAreaFilter, 0)
		return &rst
	}

	accessAreaFiltersOpts := make([]model.AccessAreaFilter, len(accessAreaFilters))
	for i, v := range accessAreaFilters {
		filterMap := v.(map[string]interface{})
		filterOpt := model.AccessAreaFilter{
			Type:         utils.String(filterMap["type"].(string)),
			ContentType:  utils.String(filterMap["content_type"].(string)),
			Area:         utils.String(filterMap["area"].(string)),
			ContentValue: utils.StringIgnoreEmpty(filterMap["content_value"].(string)),
			ExceptionIp:  utils.StringIgnoreEmpty(filterMap["exception_ip"].(string)),
		}
		accessAreaFiltersOpts[i] = filterOpt
	}
	return &accessAreaFiltersOpts
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
			Weight:              utils.Int32IgnoreEmpty(int32(source["weight"].(int))),
			ObsBucketType:       utils.StringIgnoreEmpty(source["obs_bucket_type"].(string)),
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
			FollowOrigin:      utils.StringIgnoreEmpty(parseFunctionEnabledStatus(followOrigin)),
			MatchType:         utils.StringIgnoreEmpty(parseCacheRuleType(rule["rule_type"].(string))),
			MatchValue:        utils.StringIgnoreEmpty(rule["content"].(string)),
			Ttl:               utils.Int32(int32(rule["ttl"].(int))),
			TtlUnit:           parseCacheTTLUnits(rule["ttl_type"].(string)),
			Priority:          int32(rule["priority"].(int)),
			UrlParameterType:  utils.StringIgnoreEmpty(rule["url_parameter_type"].(string)),
			UrlParameterValue: utils.StringIgnoreEmpty(rule["url_parameter_value"].(string)),
		}
	}
	return &result
}

func buildIpv6AccelerateOpts(ipv6Enable bool) *int32 {
	ipv6Accelerate := 0
	if ipv6Enable {
		ipv6Accelerate = 1
	}
	return utils.Int32(int32(ipv6Accelerate))
}

// buildUpdateDomainFullConfigsOpts Build CDN domain config opts from field `configs`
func buildUpdateDomainFullConfigsOpts(configsOpts *model.Configs, configs map[string]interface{}, d *schema.ResourceData) {
	if d.HasChange("configs.0.ipv6_enable") {
		configsOpts.Ipv6Accelerate = buildIpv6AccelerateOpts(configs["ipv6_enable"].(bool))
	}
	if d.HasChange("configs.0.range_based_retrieval_enabled") {
		retrievalEnabled := configs["range_based_retrieval_enabled"].(bool)
		configsOpts.OriginRangeStatus = utils.String(parseFunctionEnabledStatus(retrievalEnabled))
	}
	if d.HasChange("configs.0.description") {
		configsOpts.Remark = utils.String(configs["description"].(string))
	}
	if d.HasChange("configs.0.slice_etag_status") {
		configsOpts.SliceEtagStatus = utils.StringIgnoreEmpty(configs["slice_etag_status"].(string))
	}
	if d.HasChange("configs.0.origin_receive_timeout") {
		configsOpts.OriginReceiveTimeout = utils.Int32IgnoreEmpty(int32(configs["origin_receive_timeout"].(int)))
	}
	if d.HasChange("configs.0.origin_follow302_status") {
		configsOpts.OriginFollow302Status = utils.StringIgnoreEmpty(configs["origin_follow302_status"].(string))
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
	if d.HasChange("configs.0.quic") {
		configsOpts.Quic = buildQUICOpts(configs["quic"].([]interface{}))
	}
	if d.HasChange("configs.0.referer") {
		configsOpts.Referer = buildRefererOpts(configs["referer"].([]interface{}))
	}
	if d.HasChange("configs.0.video_seek") {
		configsOpts.VideoSeek = buildVideoSeekOpts(configs["video_seek"].([]interface{}))
	}
	if d.HasChange("configs.0.request_limit_rules") {
		configsOpts.RequestLimitRules = buildRequestLimitRulesOpts(configs["request_limit_rules"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.error_code_cache") {
		configsOpts.ErrorCodeCache = buildErrorCodeCacheOpts(configs["error_code_cache"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.ip_filter") {
		configsOpts.IpFilter = buildIpFilterOpts(configs["ip_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.origin_request_url_rewrite") {
		originRequestUrlRewrites := configs["origin_request_url_rewrite"].(*schema.Set).List()
		configsOpts.OriginRequestUrlRewrite = buildOriginRequestUrlRewriteOpts(originRequestUrlRewrites)
	}
	if d.HasChange("configs.0.user_agent_filter") {
		configsOpts.UserAgentFilter = buildUserAgentFilterOpts(configs["user_agent_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.error_code_redirect_rules") {
		errorCodeRedirectRules := configs["error_code_redirect_rules"].(*schema.Set).List()
		configsOpts.ErrorCodeRedirectRules = buildErrorCodeRedirectRules(errorCodeRedirectRules)
	}
	if d.HasChange("configs.0.hsts") {
		configsOpts.Hsts = buildHstsOpts(configs["hsts"].([]interface{}))
	}
	if d.HasChange("configs.0.access_area_filter") {
		accessAreaFilters := configs["access_area_filter"].(*schema.Set).List()
		configsOpts.AccessAreaFilter = buildAccessAreaFilters(accessAreaFilters)
	}
}

func updateDomainFullConfigs(client *cdnv2.CdnClient, cfg *config.Config, d *schema.ResourceData) error {
	// When the configs configuration is empty, the interface will report an error.
	// Make fields `business_type` and `service_area` are configured by default.
	configsOpts := model.Configs{
		BusinessType: utils.StringIgnoreEmpty(d.Get("type").(string)),
		ServiceArea:  utils.StringIgnoreEmpty(d.Get("service_area").(string)),
	}
	if d.HasChange("sources") {
		configsOpts.Sources = buildSourcesOpts(d.Get("sources").(*schema.Set).List())
	}

	if d.HasChange("configs") {
		rawConfigs := d.Get("configs").([]interface{})
		if len(rawConfigs) > 0 && rawConfigs[0] != nil {
			buildUpdateDomainFullConfigsOpts(&configsOpts, rawConfigs[0].(map[string]interface{}), d)
		}
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

	opts := buildResourceExtensionOpts(d, cfg)
	if err := waitingForStatusOnline(ctx, cdnClient, d, d.Timeout(schema.TimeoutCreate), opts); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) creation to become online: %s", d.Id(), err)
	}
	return resourceCdnDomainUpdate(ctx, d, meta)
}

func waitingForStatusOnline(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, opts *domains.ExtensionOpts) error {
	domainName := d.Get("name").(string)
	unexpectedStatus := []string{"offline", "configure_failed", "check_failed", "deleting"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domain, err := domains.GetByName(client, domainName, opts).Extract()
			if err != nil {
				return nil, "ERROR", err
			}

			if domain == nil {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: Domain is not found in API response")
			}

			status := domain.DomainStatus
			if status == "online" {
				return domain, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return domain, status, nil
			}
			return domain, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForStatusOffline(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, opts *domains.ExtensionOpts) error {
	domainName := d.Get("name").(string)
	unexpectedStatus := []string{"online", "configure_failed", "check_failed", "deleting"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domain, err := domains.GetByName(client, domainName, opts).Extract()
			if err != nil {
				return nil, "ERROR", err
			}

			if domain == nil {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: Domain is not found in API response")
			}

			status := domain.DomainStatus
			if status == "offline" {
				return domain, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return domain, status, nil
			}
			return domain, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
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

// flattenHTTPSAttrs Field `privateKey` is not returned in the details interface.
// The value of the field `certificateBody` will be modified by the cloud, resulting in inconsistency with the local value.
func flattenHTTPSAttrs(https *model.HttpGetBody, privateKey, certificateBody string) []map[string]interface{} {
	if https == nil {
		return nil
	}
	httpsAttrs := map[string]interface{}{
		"https_status":         https.HttpsStatus,
		"certificate_name":     https.CertificateName,
		"certificate_body":     certificateBody,
		"private_key":          privateKey,
		"certificate_source":   https.CertificateSource,
		"certificate_type":     https.CertificateType,
		"http2_status":         https.Http2Status,
		"tls_version":          https.TlsVersion,
		"ocsp_stapling_status": https.OcspStaplingStatus,
		"https_enabled":        analyseFunctionEnabledStatusPtr(https.HttpsStatus),
		"http2_enabled":        analyseFunctionEnabledStatusPtr(https.Http2Status),
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

func flattenUrlAuthAttrs(urlAuth *model.UrlAuthGetBody, urlAuthKey, urlAuthBackupKey string) []map[string]interface{} {
	if urlAuth == nil {
		return nil
	}

	urlAuthAttrs := map[string]interface{}{
		"enabled":        analyseFunctionEnabledStatus(urlAuth.Status),
		"status":         urlAuth.Status,
		"type":           urlAuth.Type,
		"sign_method":    urlAuth.SignMethod,
		"match_type":     urlAuth.MatchType,
		"inherit_config": flattenInheritConfigAttrs(urlAuth.InheritConfig),
		"sign_arg":       urlAuth.SignArg,
		"key":            urlAuthKey,
		"backup_key":     urlAuthBackupKey,
		"time_format":    urlAuth.TimeFormat,
		"expire_time":    urlAuth.ExpireTime,
	}

	return []map[string]interface{}{urlAuthAttrs}
}

func flattenInheritConfigAttrs(inheritConfig *model.InheritConfigQuery) []map[string]interface{} {
	if inheritConfig == nil {
		return nil
	}

	inheritConfigAttrs := map[string]interface{}{
		"enabled":           analyseFunctionEnabledStatus(inheritConfig.Status),
		"status":            inheritConfig.Status,
		"inherit_type":      inheritConfig.InheritType,
		"inherit_time_type": inheritConfig.InheritTimeType,
	}

	return []map[string]interface{}{inheritConfigAttrs}
}

func flattenForceRedirectAttrs(forceRedirect *model.ForceRedirectConfig) []map[string]interface{} {
	if forceRedirect == nil {
		return nil
	}

	forceRedirectAttrs := map[string]interface{}{
		"status":        forceRedirect.Status,
		"type":          forceRedirect.Type,
		"enabled":       analyseFunctionEnabledStatus(forceRedirect.Status),
		"redirect_code": forceRedirect.RedirectCode,
	}

	return []map[string]interface{}{forceRedirectAttrs}
}

func flattenCompressAttrs(compress *model.Compress) []map[string]interface{} {
	if compress == nil {
		return nil
	}

	compressAttrs := map[string]interface{}{
		"status":    compress.Status,
		"type":      compress.Type,
		"file_type": compress.FileType,
		"enabled":   analyseFunctionEnabledStatus(compress.Status),
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

func flattenRemoteAuthRulesAttrs(remoteAuthRule *model.RemoteAuthRule) []map[string]interface{} {
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

func flattenQUICAttrs(quic *model.Quic) []map[string]interface{} {
	if quic == nil {
		return nil
	}

	quicAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(quic.Status),
	}

	return []map[string]interface{}{quicAttrs}
}

func flattenRefererAttrs(referer *model.RefererConfig) []map[string]interface{} {
	if referer == nil {
		return nil
	}

	refererAttrs := map[string]interface{}{
		"type":          referer.Type,
		"value":         referer.Value,
		"include_empty": referer.IncludeEmpty,
	}

	return []map[string]interface{}{refererAttrs}
}

func flattenVideoSeekAttrs(videoSeek *model.VideoSeek) []map[string]interface{} {
	if videoSeek == nil {
		// When closing `video_seek`, the API response body will not return the information of this field.
		// In order to avoid plan problems in terraform, a default value is added.
		return []map[string]interface{}{{
			"enable_video_seek": false,
		}}
	}

	return []map[string]interface{}{{
		"enable_video_seek":       videoSeek.EnableVideoSeek,
		"enable_flv_by_time_seek": videoSeek.EnableFlvByTimeSeek,
		"start_parameter":         videoSeek.StartParameter,
		"end_parameter":           videoSeek.EndParameter,
	}}
}

func flattenRequestLimitRulesAttrs(requestLimitRules *[]model.RequestLimitRules) []map[string]interface{} {
	if requestLimitRules == nil || len(*requestLimitRules) == 0 {
		return nil
	}

	requestLimitRulesAttrs := make([]map[string]interface{}, len(*requestLimitRules))
	for i, v := range *requestLimitRules {
		requestLimitRulesAttrs[i] = map[string]interface{}{
			"priority":         v.Priority,
			"match_type":       v.MatchType,
			"match_value":      v.MatchValue,
			"type":             v.Type,
			"limit_rate_after": v.LimitRateAfter,
			"limit_rate_value": v.LimitRateValue,
		}
	}
	return requestLimitRulesAttrs
}

func flattenErrorCodeCacheAttrs(errorCodeCache *[]model.ErrorCodeCache) []map[string]interface{} {
	if errorCodeCache == nil || len(*errorCodeCache) == 0 {
		return nil
	}

	errorCodeCacheAttrs := make([]map[string]interface{}, len(*errorCodeCache))
	for i, v := range *errorCodeCache {
		errorCodeCacheAttrs[i] = map[string]interface{}{
			"code": v.Code,
			"ttl":  v.Ttl,
		}
	}
	return errorCodeCacheAttrs
}

func flattenIpFilterAttrs(ipFilter *model.IpFilter) []map[string]interface{} {
	if ipFilter == nil {
		return nil
	}

	ipFilterAttrs := map[string]interface{}{
		"type":  ipFilter.Type,
		"value": ipFilter.Value,
	}
	return []map[string]interface{}{ipFilterAttrs}
}

func flattenOriginRequestUrlRewriteAttrs(originRequestUrlRewrite *[]model.OriginRequestUrlRewrite) []map[string]interface{} {
	if originRequestUrlRewrite == nil || len(*originRequestUrlRewrite) == 0 {
		return nil
	}

	originRequestUrlRewriteAttrs := make([]map[string]interface{}, len(*originRequestUrlRewrite))
	for i, v := range *originRequestUrlRewrite {
		originRequestUrlRewriteAttrs[i] = map[string]interface{}{
			"priority":   v.Priority,
			"match_type": v.MatchType,
			"target_url": v.TargetUrl,
			"source_url": v.SourceUrl,
		}
	}
	return originRequestUrlRewriteAttrs
}

func flattenUserAgentFilterAttrs(userAgentFilter *model.UserAgentFilter) []map[string]interface{} {
	if userAgentFilter == nil {
		return nil
	}

	userAgentFilterAttrs := map[string]interface{}{
		"type": userAgentFilter.Type,
	}
	if uaList := userAgentFilter.UaList; uaList != nil {
		userAgentFilterAttrs["ua_list"] = *uaList
	}

	return []map[string]interface{}{userAgentFilterAttrs}
}

func flattenErrorCodeRedirectRulesAttrs(errorCodeRedirectRules *[]model.ErrorCodeRedirectRules) []map[string]interface{} {
	if errorCodeRedirectRules == nil || len(*errorCodeRedirectRules) == 0 {
		return nil
	}

	errorCodeRedirectRulesAttrs := make([]map[string]interface{}, len(*errorCodeRedirectRules))
	for i, v := range *errorCodeRedirectRules {
		errorCodeRedirectRulesAttrs[i] = map[string]interface{}{
			"error_code":  v.ErrorCode,
			"target_code": v.TargetCode,
			"target_link": v.TargetLink,
		}
	}
	return errorCodeRedirectRulesAttrs
}

func flattenHstsAttrs(hsts *model.HstsQuery) []map[string]interface{} {
	if hsts == nil {
		return nil
	}

	hstsAttrs := map[string]interface{}{
		"enabled":            analyseFunctionEnabledStatus(hsts.Status),
		"max_age":            hsts.MaxAge,
		"include_subdomains": hsts.IncludeSubdomains,
	}

	return []map[string]interface{}{hstsAttrs}
}

func flattenAccessAreaFiltersAttrs(accessAreaFilters *[]model.AccessAreaFilter) []map[string]interface{} {
	if accessAreaFilters == nil || len(*accessAreaFilters) == 0 {
		return nil
	}

	accessAreaFiltersAttrs := make([]map[string]interface{}, len(*accessAreaFilters))
	for i, v := range *accessAreaFilters {
		accessAreaFiltersAttrs[i] = map[string]interface{}{
			"type":          v.Type,
			"content_type":  v.ContentType,
			"area":          v.Area,
			"content_value": v.ContentValue,
			"exception_ip":  v.ExceptionIp,
		}
	}
	return accessAreaFiltersAttrs
}

func flattenSourcesAttrs(sources *[]model.SourcesConfigResponseBody) []map[string]interface{} {
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
			"weight":                  v.Weight,
			"obs_bucket_type":         v.ObsBucketType,
		}
	}

	return sourcesAttrs
}

func flattenConfigAttrs(configsResp *model.ConfigsGetBody, d *schema.ResourceData) []map[string]interface{} {
	privateKey := d.Get("configs.0.https_settings.0.private_key").(string)
	certificateBody := d.Get("configs.0.https_settings.0.certificate_body").(string)
	urlAuthKey := d.Get("configs.0.url_signing.0.key").(string)
	urlAuthBackupKey := d.Get("configs.0.url_signing.0.backup_key").(string)

	configsAttrs := map[string]interface{}{
		"https_settings":                flattenHTTPSAttrs(configsResp.Https, privateKey, certificateBody),
		"retrieval_request_header":      flattenOriginRequestHeaderAttrs(configsResp.OriginRequestHeader),
		"http_response_header":          flattenHttpResponseHeaderAttrs(configsResp.HttpResponseHeader),
		"url_signing":                   flattenUrlAuthAttrs(configsResp.UrlAuth, urlAuthKey, urlAuthBackupKey),
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
		"description":                   configsResp.Remark,
		"slice_etag_status":             configsResp.SliceEtagStatus,
		"origin_receive_timeout":        configsResp.OriginReceiveTimeout,
		"origin_follow302_status":       configsResp.OriginFollow302Status,
		"quic":                          flattenQUICAttrs(configsResp.Quic),
		"referer":                       flattenRefererAttrs(configsResp.Referer),
		"video_seek":                    flattenVideoSeekAttrs(configsResp.VideoSeek),
		"request_limit_rules":           flattenRequestLimitRulesAttrs(configsResp.RequestLimitRules),
		"error_code_cache":              flattenErrorCodeCacheAttrs(configsResp.ErrorCodeCache),
		"ip_filter":                     flattenIpFilterAttrs(configsResp.IpFilter),
		"origin_request_url_rewrite":    flattenOriginRequestUrlRewriteAttrs(configsResp.OriginRequestUrlRewrite),
		"user_agent_filter":             flattenUserAgentFilterAttrs(configsResp.UserAgentFilter),
		"error_code_redirect_rules":     flattenErrorCodeRedirectRulesAttrs(configsResp.ErrorCodeRedirectRules),
		"hsts":                          flattenHstsAttrs(configsResp.Hsts),
		"access_area_filter":            flattenAccessAreaFiltersAttrs(configsResp.AccessAreaFilter),
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

func queryAndFlattenDomainTags(cdnClient *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]string, error) {
	tags, err := domains.GetTags(cdnClient, d.Id())
	if err != nil {
		return nil, fmt.Errorf("error retrieving CDN domain tags: %s", err)
	}
	return utils.TagsToMap(tags), nil
}

func resourceCdnDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cdnClient, err := cfg.CdnV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v1 client: %s", err)
	}

	hcCdnClient, err := cfg.HcCdnV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	domainName := d.Get("domain_name").(string)
	if domainName == "" {
		domainName = d.Get("name").(string)
	}

	v, err := domains.GetByName(cdnClient, domainName, buildResourceExtensionOpts(d, cfg)).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, parseDetailResponseError(err), "error retrieving CDN domain")
	}

	if v == nil {
		return diag.Errorf("error retrieving CDN domain: Domain is not found in API response")
	}

	// Backfield the ID when executing the import operation
	d.SetId(v.ID)
	configsResp, err := queryDomainFullConfig(hcCdnClient, cfg, d, v.DomainName)
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := queryAndFlattenDomainTags(cdnClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", v.DomainName),
		d.Set("type", v.BusinessType),
		d.Set("cname", v.CName),
		d.Set("domain_status", v.DomainStatus),
		d.Set("service_area", v.ServiceArea),
		d.Set("sources", flattenSourcesAttrs(configsResp.Sources)),
		d.Set("configs", flattenConfigAttrs(configsResp, d)),
		d.Set("tags", tags),
		d.Set("domain_name", v.DomainName),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// When the domain name does not exist, the response body example of the details interface is as follows:
// {"error": {"error_code": "CDN.0170","error_msg": "domain not exist!"}}
func parseDetailResponseError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode, errorCodeErr := jmespath.Search("error.error_code", apiError)
		if errorCodeErr != nil || errorCode == nil {
			return err
		}

		if errorCode.(string) == "CDN.0170" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func resourceCdnDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	if d.HasChanges("sources", "configs", "cache_settings", "type", "service_area") || d.IsNewResource() {
		err = updateDomainFullConfigs(hcCdnClient, cfg, d)
		if err != nil {
			return diag.Errorf("error updating CDN domain configs settings: %s", err)
		}

		cdnClient, err := cfg.CdnV1Client(region)
		if err != nil {
			return diag.Errorf("error creating CDN v1 client: %s", err)
		}
		opts := buildResourceExtensionOpts(d, cfg)
		if err := waitingForStatusOnline(ctx, cdnClient, d, d.Timeout(schema.TimeoutUpdate), opts); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become online: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err := updateDomainTags(hcCdnClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "cdn",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
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

		if err := waitingForStatusOffline(ctx, cdnClient, d, d.Timeout(schema.TimeoutDelete), opts); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become offline: %s", d.Id(), err)
		}
	}

	_, err = domains.Delete(cdnClient, d.Id(), opts).Extract()
	if err != nil {
		// When the domain does not exist, the deletion API will report an error and return the following information:
		// {"error": {"error_code": "CDN.0000","error_msg": "domain is null or more than one."}}.
		// The error code "CDN.0000" indicates an internal system error and cannot be used to prove that the resource
		// no longer exists, so the logic of checkDeleted is not added.
		return diag.Errorf("error deleting CDN domain (%s): %s", d.Id(), err)
	}

	if err := waitingForDomainDeleted(ctx, cdnClient, d, d.Timeout(schema.TimeoutDelete), opts); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func waitingForDomainDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, opts *domains.ExtensionOpts) error {
	domainName := d.Get("name").(string)
	unexpectedStatus := []string{"online", "offline", "configuring", "configure_failed", "checking", "check_failed"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domain, err := domains.GetByName(client, domainName, opts).Extract()
			if err != nil {
				parseErr := parseDeleteDetailResponseError(err)
				if _, ok := parseErr.(golangsdk.ErrDefault404); ok {
					return "success", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			if domain == nil {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: Domain is not found in API response")
			}

			status := domain.DomainStatus
			if utils.StrSliceContains(unexpectedStatus, status) {
				return domain, status, nil
			}
			return domain, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// When the deletion interface is successfully called, in the error information responded by the query details interface,
// the following two situations need to be processed as 404:
// {"error": {"error_code": "CDN.0170","error_msg": "domain not exist!"}}
// {"error": {"error_code": "CDN.00010182","error_msg": "The resource is not belong to the enterprise project."}}
func parseDeleteDetailResponseError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode, errorCodeErr := jmespath.Search("error.error_code", apiError)
		if errorCodeErr != nil || errorCode == nil {
			return err
		}

		if errorCode.(string) == "CDN.0170" || errorCode.(string) == "CDN.00010182" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func buildResourceExtensionOpts(d *schema.ResourceData, cfg *config.Config) *domains.ExtensionOpts {
	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		return &domains.ExtensionOpts{
			EnterpriseProjectId: epsID,
		}
	}

	return nil
}

func resourceCDNDomainImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("domain_name", d.Id())
}
