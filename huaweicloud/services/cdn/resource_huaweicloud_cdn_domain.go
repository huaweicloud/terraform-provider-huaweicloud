package cdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

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
			"scm_certificate_id": {
				Type:     schema.TypeString,
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
			// This field has a default value, so Computed is added.
			"include_empty": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

var requestUrlRewrite = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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
						"match_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"redirect_status_code": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"redirect_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	},
}

var browserCacheRules = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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
						"match_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cache_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ttl_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

var sni = schema.Schema{
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
			"server_name": {
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

var accessAreaFilters = schema.Schema{
	Type:        schema.TypeSet,
	Optional:    true,
	Description: "Specifies the geographic access control rules.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the blacklist and whitelist rule type.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the content type.",
			},
			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the areas, separated by commas.",
			},
			"content_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the content value.",
			},
			"exception_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the IP addresses exception in access control, separated by commas.",
			},
		},
	},
}

var clientCert = schema.Schema{
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
			"trusted_cert": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hosts": {
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

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableParams),
			config.MergeDefaultTags(),
		),

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
						"sni":                        &sni,
						"request_url_rewrite":        &requestUrlRewrite,
						"browser_cache_rules":        &browserCacheRules,
						"access_area_filter":         &accessAreaFilters,
						"client_cert":                &clientCert,
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

func parseFunctionEnabledStatus(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
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

func buildCreateCdnDomainSourcesBodyParams(d *schema.ResourceData) []interface{} {
	sources := d.Get("sources").(*schema.Set).List()
	rst := make([]interface{}, 0, len(sources))

	for _, v := range sources {
		sourceMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"ip_or_domain":   sourceMap["origin"],
			"origin_type":    sourceMap["origin_type"],
			"active_standby": sourceMap["active"],
		})
	}
	return rst
}

func buildCreateCdnDomainBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_name":           d.Get("name"),
		"business_type":         d.Get("type"),
		"sources":               buildCreateCdnDomainSourcesBodyParams(d),
		"service_area":          utils.ValueIgnoreEmpty(d.Get("service_area")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return map[string]interface{}{
		"domain": bodyParams,
	}
}

func buildCdnDomainQueryParams(epsID string) string {
	if epsID == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsID)
}

func ReadCdnDomainDetail(client *golangsdk.ServiceClient, domainName, epsID string) (interface{}, error) {
	getPath := client.Endpoint + "v1.0/cdn/configuration/domains/{domain_name}"
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)
	getPath += buildCdnDomainQueryParams(epsID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func waitingForCdnDomainStatusOnline(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, cfg *config.Config) error {
	var (
		domainName       = d.Get("name").(string)
		epsID            = cfg.GetEnterpriseProjectID(d)
		unexpectedStatus = []string{"offline", "configure_failed", "check_failed", "deleting"}
	)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domainResp, err := ReadCdnDomainDetail(client, domainName, epsID)
			if err != nil {
				return nil, "ERROR", err
			}

			domainStatus := utils.PathSearch("domain.domain_status", domainResp, "").(string)
			if domainStatus == "" {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: domain_status is not found in API response")
			}

			if domainStatus == "online" {
				return domainResp, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, domainStatus) {
				return domainResp, domainStatus, nil
			}
			return domainResp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCdnDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/cdn/domains"
		product = "cdn"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCdnDomainBodyParams(d, cfg)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CDN domain: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}
	// Even if there is an error when creating an API, the status code for the API response is always 200.
	errorCode := utils.PathSearch("error.error_code", createRespBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", createRespBody, "").(string)
		return diag.Errorf("error creating CDN domain: error_code (%s), error_msg (%s)", errorCode, errorMsg)
	}

	id := utils.PathSearch("domain.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CDN domain: ID is not found in API response")
	}
	d.SetId(id)

	if err := waitingForCdnDomainStatusOnline(ctx, client, d, d.Timeout(schema.TimeoutCreate), cfg); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) creation to become online: %s", d.Id(), err)
	}

	return resourceCdnDomainUpdate(ctx, d, meta)
}

func analyseFunctionEnabledStatus(enabledStatus string) bool {
	return enabledStatus == "on"
}

// In order to prevent the user from modifying the field `name` and causing the resource not to be queried, first try
// to obtain the value of the attribute field `domain_name`.
func getDomainName(d *schema.ResourceData) string {
	if domainName := d.Get("domain_name").(string); domainName != "" {
		return domainName
	}
	return d.Get("name").(string)
}

func queryCdnDomainFullConfig(client *golangsdk.ServiceClient, domainName, epsID string) (interface{}, error) {
	getPath := client.Endpoint + "v1.1/cdn/configuration/domains/{domain_name}/configs"
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)
	getPath += buildCdnDomainQueryParams(epsID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func queryAndFlattenCdnDomainTags(client *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]interface{}, error) {
	getTagPath := client.Endpoint + "v1.0/cdn/configuration/tags"
	getTagPath += fmt.Sprintf("?resource_id=%s", d.Id())
	getTagOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTagResp, err := client.Request("GET", getTagPath, &getTagOpt)
	if err != nil {
		return nil, err
	}

	getTagRespBody, err := utils.FlattenResponse(getTagResp)
	if err != nil {
		return nil, err
	}
	return utils.FlattenTagsToMap(utils.PathSearch("tags", getTagRespBody, nil)), nil
}

func flattenSourcesActive(priority int64) interface{} {
	if priority == 70 {
		return 1
	}
	return 0
}

func flattenSourcesAttributes(configResp interface{}) []interface{} {
	if configResp == nil {
		return nil
	}

	curJson := utils.PathSearch("configs.sources", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		sourceMap := v.(map[string]interface{})

		var priority float64
		if float64Value, ok := sourceMap["priority"].(float64); ok {
			priority = float64Value
		}

		var obsWebHostingStatus string
		if stringValue, ok := sourceMap["obs_web_hosting_status"].(string); ok {
			obsWebHostingStatus = stringValue
		}

		rst = append(rst, map[string]interface{}{
			"origin":                  sourceMap["origin_addr"],
			"origin_type":             sourceMap["origin_type"],
			"active":                  flattenSourcesActive(int64(priority)),
			"obs_web_hosting_enabled": analyseFunctionEnabledStatus(obsWebHostingStatus),
			"http_port":               sourceMap["http_port"],
			"https_port":              sourceMap["https_port"],
			"retrieval_host":          sourceMap["host_name"],
			"weight":                  sourceMap["weight"],
			"obs_bucket_type":         sourceMap["obs_bucket_type"],
		})
	}
	return rst
}

// flattenHTTPSAttributes Field `privateKey` is not returned in the details interface.
// The value of the field `certificateBody` will be modified by the cloud, resulting in inconsistency with the local value.
func flattenHTTPSAttributes(configResp interface{}, privateKey, certificateBody string) []map[string]interface{} {
	curJson := utils.PathSearch("configs.https", configResp, nil)
	if curJson == nil {
		return nil
	}

	httpsMap := curJson.(map[string]interface{})
	var httpsStatus string
	if stringValue, ok := httpsMap["https_status"].(string); ok {
		httpsStatus = stringValue
	}

	var http2Status string
	if stringValue, ok := httpsMap["http2_status"].(string); ok {
		http2Status = stringValue
	}

	httpsAttributes := map[string]interface{}{
		"https_status":         httpsMap["https_status"],
		"certificate_name":     httpsMap["certificate_name"],
		"certificate_body":     certificateBody,
		"private_key":          privateKey,
		"certificate_source":   httpsMap["certificate_source"],
		"scm_certificate_id":   httpsMap["scm_certificate_id"],
		"certificate_type":     httpsMap["certificate_type"],
		"http2_status":         httpsMap["http2_status"],
		"tls_version":          httpsMap["tls_version"],
		"ocsp_stapling_status": httpsMap["ocsp_stapling_status"],
		"https_enabled":        analyseFunctionEnabledStatus(httpsStatus),
		"http2_enabled":        analyseFunctionEnabledStatus(http2Status),
	}

	return []map[string]interface{}{httpsAttributes}
}

func flattenOriginRequestHeaderAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.origin_request_header", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		headerMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"name":   headerMap["name"],
			"value":  headerMap["value"],
			"action": headerMap["action"],
		})
	}
	return rst
}

func flattenHttpResponseHeaderAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.http_response_header", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		headerMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"name":   headerMap["name"],
			"value":  headerMap["value"],
			"action": headerMap["action"],
		})
	}
	return rst
}

func flattenInheritConfigAttributes(inheritConfigResp interface{}) []map[string]interface{} {
	if inheritConfigResp == nil {
		return nil
	}

	inheritConfigRespMap := inheritConfigResp.(map[string]interface{})
	var status string
	if stringValue, ok := inheritConfigRespMap["status"].(string); ok {
		status = stringValue
	}

	inheritConfigAttrs := map[string]interface{}{
		"enabled":           analyseFunctionEnabledStatus(status),
		"status":            inheritConfigRespMap["status"],
		"inherit_type":      inheritConfigRespMap["inherit_type"],
		"inherit_time_type": inheritConfigRespMap["inherit_time_type"],
	}

	return []map[string]interface{}{inheritConfigAttrs}
}

func flattenUrlAuthAttributes(configResp interface{}, urlAuthKey, urlAuthBackupKey string) []map[string]interface{} {
	curJson := utils.PathSearch("configs.url_auth", configResp, nil)
	if curJson == nil {
		return nil
	}

	authMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := authMap["status"].(string); ok {
		status = stringValue
	}

	urlAuthAttrs := map[string]interface{}{
		"enabled":        analyseFunctionEnabledStatus(status),
		"status":         authMap["status"],
		"type":           authMap["type"],
		"sign_method":    authMap["sign_method"],
		"match_type":     authMap["match_type"],
		"inherit_config": flattenInheritConfigAttributes(authMap["inherit_config"]),
		"sign_arg":       authMap["sign_arg"],
		"key":            urlAuthKey,
		"backup_key":     urlAuthBackupKey,
		"time_format":    authMap["time_format"],
		"expire_time":    authMap["expire_time"],
	}
	return []map[string]interface{}{urlAuthAttrs}
}

func flattenForceRedirectAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.force_redirect", configResp, nil)
	if curJson == nil {
		return nil
	}

	forceRedirectMap := curJson.(map[string]interface{})

	var status string
	if stringValue, ok := forceRedirectMap["status"].(string); ok {
		status = stringValue
	}

	forceRedirectAttrs := map[string]interface{}{
		"status":        forceRedirectMap["status"],
		"type":          forceRedirectMap["type"],
		"enabled":       analyseFunctionEnabledStatus(status),
		"redirect_code": forceRedirectMap["redirect_code"],
	}

	return []map[string]interface{}{forceRedirectAttrs}
}

func flattenCompressAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.compress", configResp, nil)
	if curJson == nil {
		return nil
	}

	compressMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := compressMap["status"].(string); ok {
		status = stringValue
	}

	compressAttrs := map[string]interface{}{
		"status":    compressMap["status"],
		"type":      compressMap["type"],
		"file_type": compressMap["file_type"],
		"enabled":   analyseFunctionEnabledStatus(status),
	}

	return []map[string]interface{}{compressAttrs}
}

func flattenCacheUrlParameterFilterAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.cache_url_parameter_filter", configResp, nil)
	if curJson == nil {
		return nil
	}

	filterMap := curJson.(map[string]interface{})
	cacheUrlParameterFilterAttrs := map[string]interface{}{
		"value": filterMap["value"],
		"type":  filterMap["type"],
	}

	return []map[string]interface{}{cacheUrlParameterFilterAttrs}
}

func flattenIpFrequencyLimitAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.ip_frequency_limit", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := rawMap["status"].(string); ok {
		status = stringValue
	}

	ipFrequencyLimitAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(status),
		"qps":     rawMap["qps"],
	}

	return []map[string]interface{}{ipFrequencyLimitAttrs}
}

func flattenWebsocketAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.websocket", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := rawMap["status"].(string); ok {
		status = stringValue
	}

	websocketAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(status),
		"timeout": rawMap["timeout"],
	}

	return []map[string]interface{}{websocketAttrs}
}

func flattenFlexibleOriginBackSourceAttributes(backSources interface{}) []interface{} {
	if backSources == nil {
		return nil
	}

	curArray := backSources.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"sources_type":    rawMap["sources_type"],
			"ip_or_domain":    rawMap["ip_or_domain"],
			"obs_bucket_type": rawMap["obs_bucket_type"],
			"http_port":       rawMap["http_port"],
			"https_port":      rawMap["https_port"],
		})
	}
	return rst
}

func flattenFlexibleOriginAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.flexible_origin", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		flexibleOriginMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"match_type":    flexibleOriginMap["match_type"],
			"match_pattern": flexibleOriginMap["match_pattern"],
			"priority":      flexibleOriginMap["priority"],
			"back_sources":  flattenFlexibleOriginBackSourceAttributes(flexibleOriginMap["back_sources"]),
		})
	}
	return rst
}

func flattenCustomArgsAttributes(customArgs interface{}) []interface{} {
	if customArgs == nil {
		return nil
	}

	curArray := customArgs.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"type":  rawMap["type"],
			"key":   rawMap["key"],
			"value": rawMap["value"],
		})
	}
	return rst
}

func flattenRemoteAuthRulesAttributes(remoteAuthRules interface{}) []map[string]interface{} {
	if remoteAuthRules == nil {
		return nil
	}

	rawMap := remoteAuthRules.(map[string]interface{})
	remoteAuthRuleAttrs := map[string]interface{}{
		"auth_server":              rawMap["auth_server"],
		"request_method":           rawMap["request_method"],
		"file_type_setting":        rawMap["file_type_setting"],
		"specified_file_type":      rawMap["specified_file_type"],
		"reserve_args_setting":     rawMap["reserve_args_setting"],
		"reserve_args":             rawMap["reserve_args"],
		"add_custom_args_rules":    flattenCustomArgsAttributes(rawMap["add_custom_args_rules"]),
		"reserve_headers_setting":  rawMap["reserve_headers_setting"],
		"add_custom_headers_rules": flattenCustomArgsAttributes(rawMap["add_custom_headers_rules"]),
		"auth_success_status":      rawMap["auth_success_status"],
		"auth_failed_status":       rawMap["auth_failed_status"],
		"response_status":          rawMap["response_status"],
		"timeout":                  rawMap["timeout"],
		"timeout_action":           rawMap["timeout_action"],
		"reserve_headers":          rawMap["reserve_headers"],
	}
	return []map[string]interface{}{remoteAuthRuleAttrs}
}

func flattenRemoteAuthAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.remote_auth", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	var remoteAuthentication string
	if stringValue, ok := rawMap["remote_authentication"].(string); ok {
		remoteAuthentication = stringValue
	}

	remoteAuthAttrs := map[string]interface{}{
		"enabled":           analyseFunctionEnabledStatus(remoteAuthentication),
		"remote_auth_rules": flattenRemoteAuthRulesAttributes(rawMap["remote_auth_rules"]),
	}
	return []map[string]interface{}{remoteAuthAttrs}
}

func flattenIpv6EnableAttributes(configResp interface{}) bool {
	ipv6Accelerate := utils.PathSearch("configs.ipv6_accelerate", configResp, float64(0)).(float64)
	return int64(ipv6Accelerate) == 1
}

func flattenQUICAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.quic", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := rawMap["status"].(string); ok {
		status = stringValue
	}

	quicAttrs := map[string]interface{}{
		"enabled": analyseFunctionEnabledStatus(status),
	}

	return []map[string]interface{}{quicAttrs}
}

func flattenRefererAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.referer", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	refererAttrs := map[string]interface{}{
		"type":          rawMap["type"],
		"value":         rawMap["value"],
		"include_empty": rawMap["include_empty"],
	}

	return []map[string]interface{}{refererAttrs}
}

func flattenVideoSeekAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.video_seek", configResp, nil)
	if curJson == nil {
		// When closing `video_seek`, the API response body will not return the information of this field.
		// In order to avoid plan problems in terraform, a default value is added.
		return []map[string]interface{}{{
			"enable_video_seek": false,
		}}
	}

	rawMap := curJson.(map[string]interface{})
	return []map[string]interface{}{{
		"enable_video_seek":       rawMap["enable_video_seek"],
		"enable_flv_by_time_seek": rawMap["enable_flv_by_time_seek"],
		"start_parameter":         rawMap["start_parameter"],
		"end_parameter":           rawMap["end_parameter"],
	}}
}

func flattenRequestLimitRulesAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.request_limit_rules", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"priority":         rawMap["priority"],
			"match_type":       rawMap["match_type"],
			"match_value":      rawMap["match_value"],
			"type":             rawMap["type"],
			"limit_rate_after": rawMap["limit_rate_after"],
			"limit_rate_value": rawMap["limit_rate_value"],
		})
	}
	return rst
}

func flattenErrorCodeCacheAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.error_code_cache", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"code": rawMap["code"],
			"ttl":  rawMap["ttl"],
		})
	}
	return rst
}

func flattenIpFilterAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.ip_filter", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	ipFilterAttrs := map[string]interface{}{
		"type":  rawMap["type"],
		"value": rawMap["value"],
	}
	return []map[string]interface{}{ipFilterAttrs}
}

func flattenOriginRequestUrlRewriteAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.origin_request_url_rewrite", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"priority":   rawMap["priority"],
			"match_type": rawMap["match_type"],
			"target_url": rawMap["target_url"],
			"source_url": rawMap["source_url"],
		})
	}
	return rst
}

func flattenUserAgentFilterAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.user_agent_filter", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	userAgentFilterAttrs := map[string]interface{}{
		"type":          rawMap["type"],
		"include_empty": fmt.Sprintf("%v", rawMap["include_empty"]),
	}

	if uaList, ok := rawMap["ua_list"].([]interface{}); ok {
		userAgentFilterAttrs["ua_list"] = uaList
	}
	return []map[string]interface{}{userAgentFilterAttrs}
}

func flattenErrorCodeRedirectRulesAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.error_code_redirect_rules", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"error_code":  rawMap["error_code"],
			"target_code": rawMap["target_code"],
			"target_link": rawMap["target_link"],
		})
	}
	return rst
}

func flattenHstsAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.hsts", configResp, nil)
	if curJson == nil {
		return nil
	}

	rawMap := curJson.(map[string]interface{})
	var status string
	if stringValue, ok := rawMap["status"].(string); ok {
		status = stringValue
	}

	hstsAttrs := map[string]interface{}{
		"enabled":            analyseFunctionEnabledStatus(status),
		"max_age":            rawMap["max_age"],
		"include_subdomains": rawMap["include_subdomains"],
	}

	return []map[string]interface{}{hstsAttrs}
}

func flattenSniAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.sni", configResp, nil)
	if curJson == nil {
		return nil
	}

	hstsAttrs := map[string]interface{}{
		"enabled":     analyseFunctionEnabledStatus(utils.PathSearch("status", curJson, "").(string)),
		"server_name": utils.PathSearch("server_name", curJson, nil),
		"status":      utils.PathSearch("status", curJson, nil),
	}

	return []map[string]interface{}{hstsAttrs}
}

func flattenRequestUrlRewriteAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.request_url_rewrite", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"condition":            flattenRequestUrlRewriteConditionAttributes(utils.PathSearch("condition", v, nil)),
			"redirect_url":         utils.PathSearch("redirect_url", v, nil),
			"execution_mode":       utils.PathSearch("execution_mode", v, nil),
			"redirect_status_code": utils.PathSearch("redirect_status_code", v, nil),
			"redirect_host":        utils.PathSearch("redirect_host", v, nil),
		})
	}
	return rst
}

func flattenRequestUrlRewriteConditionAttributes(conditionResp interface{}) []interface{} {
	if conditionResp == nil {
		return nil
	}

	conditionAttribute := map[string]interface{}{
		"match_type":  utils.PathSearch("match_type", conditionResp, nil),
		"match_value": utils.PathSearch("match_value", conditionResp, nil),
		"priority":    utils.PathSearch("priority", conditionResp, nil),
	}

	return []interface{}{conditionAttribute}
}

func flattenBrowserCacheRulesAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.browser_cache_rules", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"condition":  flattenBrowserCacheRulesConditionAttributes(utils.PathSearch("condition", v, nil)),
			"cache_type": utils.PathSearch("cache_type", v, nil),
			"ttl":        utils.PathSearch("ttl", v, nil),
			"ttl_unit":   utils.PathSearch("ttl_unit", v, nil),
		})
	}
	return rst
}

func flattenBrowserCacheRulesConditionAttributes(conditionResp interface{}) []interface{} {
	if conditionResp == nil {
		return nil
	}

	conditionAttribute := map[string]interface{}{
		"match_type":  utils.PathSearch("match_type", conditionResp, nil),
		"match_value": utils.PathSearch("match_value", conditionResp, nil),
		"priority":    utils.PathSearch("priority", conditionResp, nil),
	}

	return []interface{}{conditionAttribute}
}

func flattenAccessAreaFiltersAttributes(configResp interface{}) []interface{} {
	curJson := utils.PathSearch("configs.access_area_filter", configResp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"type":          rawMap["type"],
			"content_type":  rawMap["content_type"],
			"area":          rawMap["area"],
			"content_value": rawMap["content_value"],
			"exception_ip":  rawMap["exception_ip"],
		})
	}
	return rst
}

func flattenClientCertAttributes(configResp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("configs.client_cert", configResp, nil)
	if curJson == nil {
		return nil
	}

	hstsAttrs := map[string]interface{}{
		"enabled":      analyseFunctionEnabledStatus(utils.PathSearch("status", curJson, "").(string)),
		"trusted_cert": utils.PathSearch("trusted_cert", curJson, nil),
		"hosts":        utils.PathSearch("hosts", curJson, nil),
		"status":       utils.PathSearch("status", curJson, nil),
	}

	return []map[string]interface{}{hstsAttrs}
}

func flattenConfigAttributes(configResp interface{}, d *schema.ResourceData) []map[string]interface{} {
	if configResp == nil {
		return nil
	}

	privateKey := d.Get("configs.0.https_settings.0.private_key").(string)
	certificateBody := d.Get("configs.0.https_settings.0.certificate_body").(string)
	urlAuthKey := d.Get("configs.0.url_signing.0.key").(string)
	urlAuthBackupKey := d.Get("configs.0.url_signing.0.backup_key").(string)

	configsAttrs := map[string]interface{}{
		"https_settings":                flattenHTTPSAttributes(configResp, privateKey, certificateBody),
		"retrieval_request_header":      flattenOriginRequestHeaderAttributes(configResp),
		"http_response_header":          flattenHttpResponseHeaderAttributes(configResp),
		"url_signing":                   flattenUrlAuthAttributes(configResp, urlAuthKey, urlAuthBackupKey),
		"origin_protocol":               utils.PathSearch("configs.origin_protocol", configResp, nil),
		"force_redirect":                flattenForceRedirectAttributes(configResp),
		"compress":                      flattenCompressAttributes(configResp),
		"cache_url_parameter_filter":    flattenCacheUrlParameterFilterAttributes(configResp),
		"ip_frequency_limit":            flattenIpFrequencyLimitAttributes(configResp),
		"websocket":                     flattenWebsocketAttributes(configResp),
		"flexible_origin":               flattenFlexibleOriginAttributes(configResp),
		"remote_auth":                   flattenRemoteAuthAttributes(configResp),
		"ipv6_enable":                   flattenIpv6EnableAttributes(configResp),
		"range_based_retrieval_enabled": analyseFunctionEnabledStatus(utils.PathSearch("configs.origin_range_status", configResp, "").(string)),
		"description":                   utils.PathSearch("configs.remark", configResp, nil),
		"slice_etag_status":             utils.PathSearch("configs.slice_etag_status", configResp, nil),
		"origin_receive_timeout":        utils.PathSearch("configs.origin_receive_timeout", configResp, nil),
		"origin_follow302_status":       utils.PathSearch("configs.origin_follow302_status", configResp, nil),
		"quic":                          flattenQUICAttributes(configResp),
		"referer":                       flattenRefererAttributes(configResp),
		"video_seek":                    flattenVideoSeekAttributes(configResp),
		"request_limit_rules":           flattenRequestLimitRulesAttributes(configResp),
		"error_code_cache":              flattenErrorCodeCacheAttributes(configResp),
		"ip_filter":                     flattenIpFilterAttributes(configResp),
		"origin_request_url_rewrite":    flattenOriginRequestUrlRewriteAttributes(configResp),
		"user_agent_filter":             flattenUserAgentFilterAttributes(configResp),
		"error_code_redirect_rules":     flattenErrorCodeRedirectRulesAttributes(configResp),
		"hsts":                          flattenHstsAttributes(configResp),
		"sni":                           flattenSniAttributes(configResp),
		"request_url_rewrite":           flattenRequestUrlRewriteAttributes(configResp),
		"browser_cache_rules":           flattenBrowserCacheRulesAttributes(configResp),
		"access_area_filter":            flattenAccessAreaFiltersAttributes(configResp),
		"client_cert":                   flattenClientCertAttributes(configResp),
	}
	return []map[string]interface{}{configsAttrs}
}

func resourceCdnDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "cdn"
		domainName = getDomainName(d)
		epsID      = cfg.GetEnterpriseProjectID(d)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	domainResp, err := ReadCdnDomainDetail(client, domainName, epsID)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", "CDN.0170"),
			"error retrieving CDN domain")
	}

	// Backfield the ID when executing the import operation
	if id := utils.PathSearch("domain.id", domainResp, "").(string); id != "" {
		d.SetId(id)
	}

	configResp, err := queryCdnDomainFullConfig(client, domainName, epsID)
	if err != nil {
		return diag.Errorf("error retrieving CDN domain full config: %s", err)
	}

	tags, err := queryAndFlattenCdnDomainTags(client, d)
	if err != nil {
		return diag.Errorf("error retrieving CDN domain tags: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("domain.domain_name", domainResp, nil)),
		d.Set("type", utils.PathSearch("domain.business_type", domainResp, nil)),
		d.Set("cname", utils.PathSearch("domain.cname", domainResp, nil)),
		d.Set("domain_status", utils.PathSearch("domain.domain_status", domainResp, nil)),
		d.Set("service_area", utils.PathSearch("domain.service_area", domainResp, nil)),
		d.Set("sources", flattenSourcesAttributes(configResp)),
		d.Set("configs", flattenConfigAttributes(configResp, d)),
		d.Set("tags", tags),
		d.Set("domain_name", utils.PathSearch("domain.domain_name", domainResp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCdnDomainSourcesPriorityOpts(active int) interface{} {
	if active == 1 {
		return 70
	}
	return 30
}

func buildCdnDomainSourcesOpts(rawSources []interface{}) []interface{} {
	if len(rawSources) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawSources))
	for _, v := range rawSources {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"origin_addr":            rawMap["origin"],
			"origin_type":            rawMap["origin_type"],
			"priority":               buildCdnDomainSourcesPriorityOpts(rawMap["active"].(int)),
			"obs_web_hosting_status": parseFunctionEnabledStatus(rawMap["obs_web_hosting_enabled"].(bool)),
			"http_port":              utils.ValueIgnoreEmpty(rawMap["http_port"]),
			"https_port":             utils.ValueIgnoreEmpty(rawMap["https_port"]),
			"host_name":              utils.ValueIgnoreEmpty(rawMap["retrieval_host"]),
			"weight":                 utils.ValueIgnoreEmpty(rawMap["weight"]),
			"obs_bucket_type":        utils.ValueIgnoreEmpty(rawMap["obs_bucket_type"]),
		})
	}
	return rst
}

func buildCdnDomainIpv6AccelerateOpts(ipv6Enable bool) interface{} {
	if ipv6Enable {
		return 1
	}
	return 0
}

func buildCdnDomainHTTPSStatusOpts(enable bool) string {
	if enable {
		return "on"
	}
	return "off"
}

func buildCdnDomainHTTP2StatusOpts(enable bool) string {
	if enable {
		return "on"
	}
	// Currently, European sites do not support this parameter, so we will handle it this way for the time being.
	return ""
}

func buildCdnDomainHTTPSOpts(rawHTTPS []interface{}) map[string]interface{} {
	if len(rawHTTPS) != 1 {
		return nil
	}

	https := rawHTTPS[0].(map[string]interface{})
	rst := map[string]interface{}{
		"https_status":         buildCdnDomainHTTPSStatusOpts(https["https_enabled"].(bool)),
		"certificate_name":     utils.ValueIgnoreEmpty(https["certificate_name"]),
		"certificate_value":    utils.ValueIgnoreEmpty(https["certificate_body"]),
		"private_key":          utils.ValueIgnoreEmpty(https["private_key"]),
		"certificate_source":   https["certificate_source"],
		"scm_certificate_id":   utils.ValueIgnoreEmpty(https["scm_certificate_id"]),
		"certificate_type":     utils.ValueIgnoreEmpty(https["certificate_type"]),
		"tls_version":          utils.ValueIgnoreEmpty(https["tls_version"]),
		"ocsp_stapling_status": utils.ValueIgnoreEmpty(https["ocsp_stapling_status"]),
	}

	// The API restriction field "http2_status" is only configurable if HTTPS is enabled.
	if https["https_enabled"].(bool) {
		rst["http2_status"] = utils.ValueIgnoreEmpty(buildCdnDomainHTTP2StatusOpts(https["http2_enabled"].(bool)))
	}

	return rst
}

func buildCdnDomainOriginRequestHeaderOpts(rawOriginRequestHeader []interface{}) []interface{} {
	if len(rawOriginRequestHeader) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawOriginRequestHeader))
	for _, v := range rawOriginRequestHeader {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"name":   rawMap["name"],
			"value":  utils.ValueIgnoreEmpty(rawMap["value"]),
			"action": rawMap["action"],
		})
	}
	return rst
}

func buildCdnDomainHttpResponseHeaderOpts(rawHttpResponseHeader []interface{}) []interface{} {
	if len(rawHttpResponseHeader) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawHttpResponseHeader))
	for _, v := range rawHttpResponseHeader {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"name":   rawMap["name"],
			"value":  utils.ValueIgnoreEmpty(rawMap["value"]),
			"action": rawMap["action"],
		})
	}
	return rst
}

func buildCdnDomainInheritConfigOpts(rwaInheritConfig []interface{}) map[string]interface{} {
	if len(rwaInheritConfig) != 1 {
		return nil
	}

	inheritConfig := rwaInheritConfig[0].(map[string]interface{})
	return map[string]interface{}{
		"status":            parseFunctionEnabledStatus(inheritConfig["enabled"].(bool)),
		"inherit_type":      utils.ValueIgnoreEmpty(inheritConfig["inherit_type"]),
		"inherit_time_type": utils.ValueIgnoreEmpty(inheritConfig["inherit_time_type"]),
	}
}

func buildCdnDomainUrlAuthOpts(rawUrlAuth []interface{}) map[string]interface{} {
	if len(rawUrlAuth) != 1 {
		return nil
	}

	urlAuth := rawUrlAuth[0].(map[string]interface{})
	return map[string]interface{}{
		"status":         parseFunctionEnabledStatus(urlAuth["enabled"].(bool)),
		"type":           utils.ValueIgnoreEmpty(urlAuth["type"]),
		"sign_method":    utils.ValueIgnoreEmpty(urlAuth["sign_method"]),
		"match_type":     utils.ValueIgnoreEmpty(urlAuth["match_type"]),
		"inherit_config": buildCdnDomainInheritConfigOpts(urlAuth["inherit_config"].([]interface{})),
		"sign_arg":       utils.ValueIgnoreEmpty(urlAuth["sign_arg"]),
		"key":            utils.ValueIgnoreEmpty(urlAuth["key"]),
		"backup_key":     utils.ValueIgnoreEmpty(urlAuth["backup_key"]),
		"time_format":    utils.ValueIgnoreEmpty(urlAuth["time_format"]),
		"expire_time":    urlAuth["expire_time"],
	}
}

func buildCdnDomainForceRedirectOpts(rawForceRedirect []interface{}) map[string]interface{} {
	if len(rawForceRedirect) != 1 {
		return nil
	}

	forceRedirect := rawForceRedirect[0].(map[string]interface{})
	return map[string]interface{}{
		"status":        parseFunctionEnabledStatus(forceRedirect["enabled"].(bool)),
		"type":          utils.ValueIgnoreEmpty(forceRedirect["type"]),
		"redirect_code": utils.ValueIgnoreEmpty(forceRedirect["redirect_code"]),
	}
}

func buildCdnDomainCompressOpts(rawCompress []interface{}) map[string]interface{} {
	if len(rawCompress) != 1 {
		return nil
	}

	compress := rawCompress[0].(map[string]interface{})
	return map[string]interface{}{
		"status":    parseFunctionEnabledStatus(compress["enabled"].(bool)),
		"type":      utils.ValueIgnoreEmpty(compress["type"]),
		"file_type": utils.ValueIgnoreEmpty(compress["file_type"]),
	}
}

func buildCdnDomainCacheUrlParameterFilterOpts(rawCacheUrlParameterFilter []interface{}) map[string]interface{} {
	if len(rawCacheUrlParameterFilter) != 1 {
		return nil
	}

	cacheUrlParameterFilter := rawCacheUrlParameterFilter[0].(map[string]interface{})
	return map[string]interface{}{
		"value": utils.ValueIgnoreEmpty(cacheUrlParameterFilter["value"]),
		"type":  utils.ValueIgnoreEmpty(cacheUrlParameterFilter["type"]),
	}
}

func buildCdnDomainIpFrequencyLimitOpts(rawIpFrequencyLimit []interface{}) map[string]interface{} {
	if len(rawIpFrequencyLimit) != 1 {
		return nil
	}

	ipFrequencyLimit := rawIpFrequencyLimit[0].(map[string]interface{})
	return map[string]interface{}{
		"status": parseFunctionEnabledStatus(ipFrequencyLimit["enabled"].(bool)),
		"qps":    utils.ValueIgnoreEmpty(ipFrequencyLimit["qps"]),
	}
}

func buildCdnDomainWebsocketOpts(rawWebsocket []interface{}) map[string]interface{} {
	if len(rawWebsocket) != 1 {
		return nil
	}

	websocket := rawWebsocket[0].(map[string]interface{})
	return map[string]interface{}{
		"status":  parseFunctionEnabledStatus(websocket["enabled"].(bool)),
		"timeout": websocket["timeout"],
	}
}

func buildCdnDomainFlexibleOriginBackSourceOpts(rawBackSources []interface{}) []map[string]interface{} {
	if len(rawBackSources) != 1 {
		return nil
	}

	backSource := rawBackSources[0].(map[string]interface{})
	backSourceOpts := map[string]interface{}{
		"sources_type":    backSource["sources_type"],
		"ip_or_domain":    backSource["ip_or_domain"],
		"obs_bucket_type": utils.ValueIgnoreEmpty(backSource["obs_bucket_type"]),
		"http_port":       utils.ValueIgnoreEmpty(backSource["http_port"]),
		"https_port":      utils.ValueIgnoreEmpty(backSource["https_port"]),
	}
	return []map[string]interface{}{backSourceOpts}
}

func buildCdnDomainFlexibleOriginOpts(rawFlexibleOrigins []interface{}) []interface{} {
	if len(rawFlexibleOrigins) < 1 {
		// Define an empty array to clear all flexible origins
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(rawFlexibleOrigins))
	for _, v := range rawFlexibleOrigins {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"match_type":    rawMap["match_type"],
			"match_pattern": rawMap["match_pattern"],
			"priority":      rawMap["priority"],
			"back_sources":  buildCdnDomainFlexibleOriginBackSourceOpts(rawMap["back_sources"].([]interface{})),
		})
	}
	return rst
}

func buildCdnDomainCustomArgsOpts(rawCustomArgs []interface{}) []interface{} {
	if len(rawCustomArgs) < 1 {
		// Define an empty array to clear all custom args
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(rawCustomArgs))
	for _, v := range rawCustomArgs {
		argMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"type":  argMap["type"],
			"key":   argMap["key"],
			"value": argMap["value"],
		})
	}
	return rst
}

func buildCdnDomainRemoteAuthRulesOpts(rawRemoteAuthRules []interface{}) map[string]interface{} {
	if len(rawRemoteAuthRules) != 1 {
		return nil
	}

	remoteAuthRule := rawRemoteAuthRules[0].(map[string]interface{})
	return map[string]interface{}{
		"auth_server":              remoteAuthRule["auth_server"],
		"request_method":           remoteAuthRule["request_method"],
		"file_type_setting":        remoteAuthRule["file_type_setting"],
		"specified_file_type":      utils.ValueIgnoreEmpty(remoteAuthRule["specified_file_type"]),
		"reserve_args_setting":     remoteAuthRule["reserve_args_setting"],
		"reserve_args":             utils.ValueIgnoreEmpty(remoteAuthRule["reserve_args"]),
		"add_custom_args_rules":    buildCdnDomainCustomArgsOpts(remoteAuthRule["add_custom_args_rules"].(*schema.Set).List()),
		"reserve_headers_setting":  remoteAuthRule["reserve_headers_setting"],
		"add_custom_headers_rules": buildCdnDomainCustomArgsOpts(remoteAuthRule["add_custom_headers_rules"].(*schema.Set).List()),
		"auth_success_status":      remoteAuthRule["auth_success_status"],
		"auth_failed_status":       remoteAuthRule["auth_failed_status"],
		"response_status":          remoteAuthRule["response_status"],
		"timeout":                  remoteAuthRule["timeout"],
		"timeout_action":           remoteAuthRule["timeout_action"],
		"reserve_headers":          utils.ValueIgnoreEmpty(remoteAuthRule["reserve_headers"]),
	}
}

func buildCdnDomainRemoteAuthOpts(rawRemoteAuth []interface{}) map[string]interface{} {
	if len(rawRemoteAuth) != 1 {
		return nil
	}

	remoteAuth := rawRemoteAuth[0].(map[string]interface{})
	return map[string]interface{}{
		"remote_authentication": parseFunctionEnabledStatus(remoteAuth["enabled"].(bool)),
		"remote_auth_rules":     buildCdnDomainRemoteAuthRulesOpts(remoteAuth["remote_auth_rules"].([]interface{})),
	}
}

func buildCdnDomainQUICOpts(rawQuic []interface{}) map[string]interface{} {
	if len(rawQuic) != 1 {
		return nil
	}

	quic := rawQuic[0].(map[string]interface{})
	return map[string]interface{}{
		"status": parseFunctionEnabledStatus(quic["enabled"].(bool)),
	}
}

func buildCdnDomainRefererOpts(rawReferer []interface{}) map[string]interface{} {
	if len(rawReferer) != 1 {
		return nil
	}

	referer := rawReferer[0].(map[string]interface{})
	return map[string]interface{}{
		"type":          referer["type"],
		"value":         referer["value"],
		"include_empty": referer["include_empty"],
	}
}

func buildCdnDomainVideoSeekOpts(rawVideoSeek []interface{}) map[string]interface{} {
	if len(rawVideoSeek) != 1 {
		return nil
	}

	videoSeek := rawVideoSeek[0].(map[string]interface{})
	return map[string]interface{}{
		"enable_video_seek":       videoSeek["enable_video_seek"],
		"enable_flv_by_time_seek": videoSeek["enable_flv_by_time_seek"],
		"start_parameter":         videoSeek["start_parameter"],
		"end_parameter":           videoSeek["end_parameter"],
	}
}

func buildCdnDomainRequestLimitRulesOpts(rawRequestLimitRules []interface{}) []interface{} {
	if len(rawRequestLimitRules) < 1 {
		// Define an empty array to clear all request limit rules
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(rawRequestLimitRules))
	for _, v := range rawRequestLimitRules {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"priority":         rawMap["priority"],
			"match_type":       rawMap["match_type"],
			"match_value":      rawMap["match_value"],
			"type":             rawMap["type"],
			"limit_rate_after": rawMap["limit_rate_after"],
			"limit_rate_value": rawMap["limit_rate_value"],
		})
	}
	return rst
}

func buildCdnDomainErrorCodeCacheOpts(rawErrorCodeCache []interface{}) []interface{} {
	if len(rawErrorCodeCache) < 1 {
		// Define an empty array to clear all error code cache
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(rawErrorCodeCache))
	for _, v := range rawErrorCodeCache {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"code": rawMap["code"],
			"ttl":  rawMap["ttl"],
		})
	}
	return rst
}

func buildCdnDomainIpFilterOpts(rawIpFilter []interface{}) map[string]interface{} {
	if len(rawIpFilter) != 1 {
		return nil
	}

	ipFilter := rawIpFilter[0].(map[string]interface{})
	return map[string]interface{}{
		"type":  ipFilter["type"],
		"value": ipFilter["value"],
	}
}

func buildCdnDomainOriginRequestUrlRewriteOpts(rawOriginRequestUrlRewrite []interface{}) []interface{} {
	if len(rawOriginRequestUrlRewrite) < 1 {
		// Define an empty array to clear all origin request url rewrite
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(rawOriginRequestUrlRewrite))
	for _, v := range rawOriginRequestUrlRewrite {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"priority":   rawMap["priority"],
			"match_type": rawMap["match_type"],
			"target_url": rawMap["target_url"],
			"source_url": utils.ValueIgnoreEmpty(rawMap["source_url"]),
		})
	}
	return rst
}

// This method is used to handle three scenarios: passing true, passing false and not passing.
func buildCdnDomainUserAgentFilterIncludeEmptyOpts(includeEmpty string) interface{} {
	switch includeEmpty {
	case "true":
		return true
	case "false":
		return false
	}
	return nil
}

func buildCdnDomainUserAgentFilterOpts(rawUserAgentFilter []interface{}) map[string]interface{} {
	if len(rawUserAgentFilter) != 1 {
		return nil
	}

	userAgentFilter := rawUserAgentFilter[0].(map[string]interface{})
	return map[string]interface{}{
		"type":          userAgentFilter["type"],
		"ua_list":       utils.ExpandToStringList(userAgentFilter["ua_list"].(*schema.Set).List()),
		"include_empty": buildCdnDomainUserAgentFilterIncludeEmptyOpts(userAgentFilter["include_empty"].(string)),
	}
}

func buildCdnDomainErrorCodeRedirectRules(errorCodeRedirectRules []interface{}) []interface{} {
	if len(errorCodeRedirectRules) < 1 {
		// Define an empty array to clear all error code redirect rules
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(errorCodeRedirectRules))
	for _, v := range errorCodeRedirectRules {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"error_code":  rawMap["error_code"],
			"target_code": rawMap["target_code"],
			"target_link": rawMap["target_link"],
		})
	}
	return rst
}

func buildCdnDomainHstsOpts(rawHsts []interface{}) map[string]interface{} {
	if len(rawHsts) != 1 {
		return nil
	}

	hsts := rawHsts[0].(map[string]interface{})
	return map[string]interface{}{
		"status":             parseFunctionEnabledStatus(hsts["enabled"].(bool)),
		"max_age":            hsts["max_age"],
		"include_subdomains": utils.ValueIgnoreEmpty(hsts["include_subdomains"]),
	}
}

func buildCdnDomainSniOpts(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"status":      parseFunctionEnabledStatus(rawMap["enabled"].(bool)),
		"server_name": utils.ValueIgnoreEmpty(rawMap["server_name"]),
	}
}

func buildCdnDomainRequestUrlRewrite(rawArray []interface{}) interface{} {
	if len(rawArray) < 1 {
		// Define an empty array to clear all request url rewrite
		return make([]interface{}, 0)
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"condition":            buildRequestUrlRewriteCondition(rawMap["condition"].([]interface{})),
			"redirect_url":         rawMap["redirect_url"],
			"execution_mode":       rawMap["execution_mode"],
			"redirect_status_code": utils.ValueIgnoreEmpty(rawMap["redirect_status_code"]),
			"redirect_host":        utils.ValueIgnoreEmpty(rawMap["redirect_host"]),
		})
	}
	return rst
}

func buildRequestUrlRewriteCondition(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"match_type":  rawMap["match_type"],
		"match_value": utils.ValueIgnoreEmpty(rawMap["match_value"]),
		"priority":    rawMap["priority"],
	}
}

func buildCdnDomainBrowserCacheRules(rawArray []interface{}) interface{} {
	if len(rawArray) < 1 {
		// Define an empty array to clear all request url rewrite
		return make([]interface{}, 0)
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"condition":  buildBrowserCacheRulesCondition(rawMap["condition"].([]interface{})),
			"cache_type": rawMap["cache_type"],
			"ttl":        utils.ValueIgnoreEmpty(rawMap["ttl"]),
			"ttl_unit":   utils.ValueIgnoreEmpty(rawMap["ttl_unit"]),
		})
	}
	return rst
}

func buildBrowserCacheRulesCondition(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"match_type":  rawMap["match_type"],
		"match_value": utils.ValueIgnoreEmpty(rawMap["match_value"]),
		"priority":    rawMap["priority"],
	}
}

func buildCdnDomainAccessAreaFilters(accessAreaFilters []interface{}) []interface{} {
	if len(accessAreaFilters) < 1 {
		// Define an empty array to clear all access area filters
		return make([]interface{}, 0)
	}

	rst := make([]interface{}, 0, len(accessAreaFilters))
	for _, v := range accessAreaFilters {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"type":          rawMap["type"],
			"content_type":  rawMap["content_type"],
			"area":          rawMap["area"],
			"content_value": utils.ValueIgnoreEmpty(rawMap["content_value"]),
			"exception_ip":  utils.ValueIgnoreEmpty(rawMap["exception_ip"]),
		})
	}
	return rst
}

func buildCdnDomainClientCertOpts(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"status":       parseFunctionEnabledStatus(rawMap["enabled"].(bool)),
		"trusted_cert": utils.ValueIgnoreEmpty(rawMap["trusted_cert"]),
		"hosts":        utils.ValueIgnoreEmpty(rawMap["hosts"]),
	}
}

// nolint
func buildUpdateCdnDomainFullConfigsOpts(bodyParams map[string]interface{}, configs map[string]interface{}, d *schema.ResourceData) {
	if d.HasChange("configs.0.ipv6_enable") {
		bodyParams["ipv6_accelerate"] = buildCdnDomainIpv6AccelerateOpts(configs["ipv6_enable"].(bool))
	}
	if d.HasChange("configs.0.range_based_retrieval_enabled") {
		retrievalEnabled := configs["range_based_retrieval_enabled"].(bool)
		bodyParams["origin_range_status"] = parseFunctionEnabledStatus(retrievalEnabled)
	}
	if d.HasChange("configs.0.description") {
		bodyParams["remark"] = configs["description"]
	}
	if d.HasChange("configs.0.slice_etag_status") {
		bodyParams["slice_etag_status"] = utils.ValueIgnoreEmpty(configs["slice_etag_status"])
	}
	if d.HasChange("configs.0.origin_receive_timeout") {
		bodyParams["origin_receive_timeout"] = utils.ValueIgnoreEmpty(configs["origin_receive_timeout"])
	}
	if d.HasChange("configs.0.origin_follow302_status") {
		bodyParams["origin_follow302_status"] = utils.ValueIgnoreEmpty(configs["origin_follow302_status"])
	}
	if d.HasChange("configs.0.https_settings") {
		bodyParams["https"] = buildCdnDomainHTTPSOpts(configs["https_settings"].([]interface{}))
	}
	if d.HasChange("configs.0.retrieval_request_header") {
		bodyParams["origin_request_header"] = buildCdnDomainOriginRequestHeaderOpts(configs["retrieval_request_header"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.http_response_header") {
		bodyParams["http_response_header"] = buildCdnDomainHttpResponseHeaderOpts(configs["http_response_header"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.url_signing") {
		bodyParams["url_auth"] = buildCdnDomainUrlAuthOpts(configs["url_signing"].([]interface{}))
	}
	if d.HasChange("configs.0.origin_protocol") {
		bodyParams["origin_protocol"] = utils.ValueIgnoreEmpty(configs["origin_protocol"])
	}
	if d.HasChange("configs.0.force_redirect") {
		bodyParams["force_redirect"] = buildCdnDomainForceRedirectOpts(configs["force_redirect"].([]interface{}))
	}
	if d.HasChange("configs.0.compress") {
		bodyParams["compress"] = buildCdnDomainCompressOpts(configs["compress"].([]interface{}))
	}
	if d.HasChange("configs.0.cache_url_parameter_filter") {
		bodyParams["cache_url_parameter_filter"] = buildCdnDomainCacheUrlParameterFilterOpts(configs["cache_url_parameter_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.ip_frequency_limit") {
		bodyParams["ip_frequency_limit"] = buildCdnDomainIpFrequencyLimitOpts(configs["ip_frequency_limit"].([]interface{}))
	}
	if d.HasChange("configs.0.websocket") {
		bodyParams["websocket"] = buildCdnDomainWebsocketOpts(configs["websocket"].([]interface{}))
	}
	if d.HasChange("configs.0.flexible_origin") {
		bodyParams["flexible_origin"] = buildCdnDomainFlexibleOriginOpts(configs["flexible_origin"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.remote_auth") {
		bodyParams["remote_auth"] = buildCdnDomainRemoteAuthOpts(configs["remote_auth"].([]interface{}))
	}
	if d.HasChange("configs.0.quic") {
		bodyParams["quic"] = buildCdnDomainQUICOpts(configs["quic"].([]interface{}))
	}
	if d.HasChange("configs.0.referer") {
		bodyParams["referer"] = buildCdnDomainRefererOpts(configs["referer"].([]interface{}))
	}
	if d.HasChange("configs.0.video_seek") {
		bodyParams["video_seek"] = buildCdnDomainVideoSeekOpts(configs["video_seek"].([]interface{}))
	}
	if d.HasChange("configs.0.request_limit_rules") {
		bodyParams["request_limit_rules"] = buildCdnDomainRequestLimitRulesOpts(configs["request_limit_rules"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.error_code_cache") {
		bodyParams["error_code_cache"] = buildCdnDomainErrorCodeCacheOpts(configs["error_code_cache"].(*schema.Set).List())
	}
	if d.HasChange("configs.0.ip_filter") {
		bodyParams["ip_filter"] = buildCdnDomainIpFilterOpts(configs["ip_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.origin_request_url_rewrite") {
		originRequestUrlRewrites := configs["origin_request_url_rewrite"].(*schema.Set).List()
		bodyParams["origin_request_url_rewrite"] = buildCdnDomainOriginRequestUrlRewriteOpts(originRequestUrlRewrites)
	}
	if d.HasChange("configs.0.user_agent_filter") {
		bodyParams["user_agent_filter"] = buildCdnDomainUserAgentFilterOpts(configs["user_agent_filter"].([]interface{}))
	}
	if d.HasChange("configs.0.error_code_redirect_rules") {
		errorCodeRedirectRules := configs["error_code_redirect_rules"].(*schema.Set).List()
		bodyParams["error_code_redirect_rules"] = buildCdnDomainErrorCodeRedirectRules(errorCodeRedirectRules)
	}
	if d.HasChange("configs.0.hsts") {
		bodyParams["hsts"] = buildCdnDomainHstsOpts(configs["hsts"].([]interface{}))
	}
	if d.HasChange("configs.0.sni") {
		bodyParams["sni"] = buildCdnDomainSniOpts(configs["sni"].([]interface{}))
	}
	if d.HasChange("configs.0.request_url_rewrite") {
		requestUrlRewrite := configs["request_url_rewrite"].(*schema.Set).List()
		bodyParams["request_url_rewrite"] = buildCdnDomainRequestUrlRewrite(requestUrlRewrite)
	}
	if d.HasChange("configs.0.browser_cache_rules") {
		browserCacheRules := configs["browser_cache_rules"].(*schema.Set).List()
		bodyParams["browser_cache_rules"] = buildCdnDomainBrowserCacheRules(browserCacheRules)
	}
	if d.HasChange("configs.0.access_area_filter") {
		accessAreaFilters := configs["access_area_filter"].(*schema.Set).List()
		bodyParams["access_area_filter"] = buildCdnDomainAccessAreaFilters(accessAreaFilters)
	}
	if d.HasChange("configs.0.client_cert") {
		bodyParams["client_cert"] = buildCdnDomainClientCertOpts(configs["client_cert"].([]interface{}))
	}
}

func buildCdnDomainCacheRules(followOrigin bool, rules []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"follow_origin":       utils.ValueIgnoreEmpty(parseFunctionEnabledStatus(followOrigin)),
			"match_type":          utils.ValueIgnoreEmpty(parseCacheRuleType(rawMap["rule_type"].(string))),
			"match_value":         utils.ValueIgnoreEmpty(rawMap["content"]),
			"ttl":                 rawMap["ttl"],
			"ttl_unit":            parseCacheTTLUnits(rawMap["ttl_type"].(string)),
			"priority":            rawMap["priority"],
			"url_parameter_type":  utils.ValueIgnoreEmpty(rawMap["url_parameter_type"]),
			"url_parameter_value": utils.ValueIgnoreEmpty(rawMap["url_parameter_value"]),
		})
	}
	return rst
}

func updateCdnDomainFullConfigs(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	// When the configs configuration is empty, the interface will report an error.
	// Make fields `business_type` and `service_area` are configured by default.
	bodyParams := map[string]interface{}{
		"business_type": utils.ValueIgnoreEmpty(d.Get("type")),
		"service_area":  utils.ValueIgnoreEmpty(d.Get("service_area")),
	}

	if d.HasChange("sources") {
		bodyParams["sources"] = buildCdnDomainSourcesOpts(d.Get("sources").(*schema.Set).List())
	}

	if d.HasChange("configs") {
		rawConfigs := d.Get("configs").([]interface{})
		if len(rawConfigs) > 0 && rawConfigs[0] != nil {
			buildUpdateCdnDomainFullConfigsOpts(bodyParams, rawConfigs[0].(map[string]interface{}), d)
		}
	}

	if d.HasChange("cache_settings") {
		cacheSettings := d.Get("cache_settings").([]interface{})
		if len(cacheSettings) > 0 {
			cacheSetting := cacheSettings[0].(map[string]interface{})
			bodyParams["cache_rules"] = buildCdnDomainCacheRules(cacheSetting["follow_origin"].(bool), cacheSetting["rules"].(*schema.Set).List())
		}
	}

	updatePath := client.Endpoint + "v1.1/cdn/configuration/domains/{domain_name}/configs"
	updatePath = strings.ReplaceAll(updatePath, "{domain_name}", d.Get("name").(string))
	updatePath += buildCdnDomainQueryParams(cfg.GetEnterpriseProjectID(d))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202, 204},
		JSONBody: map[string]interface{}{
			"configs": utils.RemoveNil(bodyParams),
		},
	}
	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func createCdnDomainTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	createPath := client.Endpoint + "v1.0/cdn/configuration/tags"
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202, 204},
		JSONBody: map[string]interface{}{
			"resource_id": d.Id(),
			"tags":        utils.ExpandResourceTags(tags),
		},
	}
	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func deleteCdnDomainTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tags map[string]interface{}) error {
	tagKeys := make([]string, 0, len(tags))
	for k := range tags {
		tagKeys = append(tagKeys, k)
	}

	deletePath := client.Endpoint + "v1.0/cdn/configuration/tags/batch-delete"
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202, 204},
		JSONBody: map[string]interface{}{
			"resource_id": d.Id(),
			"tags":        tagKeys,
		},
	}
	_, err := client.Request("POST", deletePath, &deleteOpt)
	return err
}

func updateCdnDomainTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oTagsRaw, nTagsRaw := d.GetChange("tags")
	oTagsMap := oTagsRaw.(map[string]interface{})
	nTagsMap := nTagsRaw.(map[string]interface{})

	if len(oTagsMap) > 0 {
		err := deleteCdnDomainTags(client, d, oTagsMap)
		if err != nil {
			return fmt.Errorf("error deleting CDN domain tags: %s", err)
		}
	}

	if len(nTagsMap) > 0 {
		err := createCdnDomainTags(client, d, nTagsMap)
		if err != nil {
			return fmt.Errorf("error creating CDN domain tags: %s", err)
		}
	}
	return nil
}

func resourceCdnDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if d.HasChanges("sources", "configs", "cache_settings", "type", "service_area") || d.IsNewResource() {
		err = updateCdnDomainFullConfigs(client, cfg, d)
		if err != nil {
			return diag.Errorf("error updating CDN domain configs settings: %s", err)
		}

		if err := waitingForCdnDomainStatusOnline(ctx, client, d, d.Timeout(schema.TimeoutUpdate), cfg); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become online: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err := updateCdnDomainTags(client, d); err != nil {
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

func deleteCdnDomain(client *golangsdk.ServiceClient, d *schema.ResourceData, epsID string) error {
	deletePath := client.Endpoint + "v1.0/cdn/domains/{domain_id}"
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", d.Id())
	deletePath += buildCdnDomainQueryParams(epsID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return err
	}
	// Even if there is an error when deleting an API, the status code for the API response is always 200.
	errorCode := utils.PathSearch("error.error_code", deleteRespBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", deleteRespBody, "").(string)
		return fmt.Errorf("error_code (%s), error_msg (%s)", errorCode, errorMsg)
	}

	return err
}

func waitingForCdnDomainDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, epsID string) error {
	var (
		domainName               = d.Get("name").(string)
		resourceNotFoundErrCodes = []string{
			"CDN.0170",     // Domain not exist.
			"CDN.00010182", // The resource is not belong to the enterprise project.
		}
	)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domainResp, err := ReadCdnDomainDetail(client, domainName, epsID)
			if err != nil {
				parseErr := common.ConvertExpected400ErrInto404Err(err, "error.error_code", resourceNotFoundErrCodes...)
				if _, ok := parseErr.(golangsdk.ErrDefault404); ok {
					return "success", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return domainResp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func disableCdnDomain(client *golangsdk.ServiceClient, d *schema.ResourceData, epsID string) error {
	disablePath := client.Endpoint + "v1.0/cdn/domains/{domain_id}/disable"
	disablePath = strings.ReplaceAll(disablePath, "{domain_id}", d.Id())
	disablePath += buildCdnDomainQueryParams(epsID)
	disableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("PUT", disablePath, &disableOpt)
	if err != nil {
		return nil
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}
	// Even if there is an error when disabling an API, the status code for the API response is always 200.
	errorCode := utils.PathSearch("error.error_code", respBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("error.error_msg", respBody, "").(string)
		return fmt.Errorf("error_code (%s), error_msg (%s)", errorCode, errorMsg)
	}

	return err
}

func waitingForCdnDomainStatusOffline(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, epsID string) error {
	domainName := d.Get("name").(string)
	unexpectedStatus := []string{"online", "configure_failed", "check_failed", "deleting"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domainResp, err := ReadCdnDomainDetail(client, domainName, epsID)
			if err != nil {
				return nil, "ERROR", err
			}

			domainStatus := utils.PathSearch("domain.domain_status", domainResp, "").(string)
			if domainStatus == "" {
				return nil, "ERROR", fmt.Errorf("error retrieving CDN domain: domain_status is not found in API response")
			}

			if domainStatus == "offline" {
				return domainResp, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, domainStatus) {
				return domainResp, domainStatus, nil
			}
			return domainResp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCdnDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
		epsID   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if d.Get("domain_status").(string) == "online" {
		// make sure the status has changed to offline before deleting it.
		if err := disableCdnDomain(client, d, epsID); err != nil {
			return diag.Errorf("error disabling CDN domain %s: %s", d.Id(), err)
		}

		if err := waitingForCdnDomainStatusOffline(ctx, client, d, d.Timeout(schema.TimeoutDelete), epsID); err != nil {
			return diag.Errorf("error waiting for CDN domain (%s) update to become offline: %s", d.Id(), err)
		}
	}

	if err := deleteCdnDomain(client, d, epsID); err != nil {
		// When the domain does not exist, the deletion API will report an error and return the following information:
		// {"error": {"error_code": "CDN.0000","error_msg": "domain is null or more than one."}}.
		// The error code "CDN.0000" indicates an internal system error and cannot be used to prove that the resource
		// no longer exists, so the logic of checkDeleted is not added.
		return diag.Errorf("error deleting CDN domain (%s): %s", d.Id(), err)
	}

	if err := waitingForCdnDomainDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete), epsID); err != nil {
		return diag.Errorf("error waiting for CDN domain (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func resourceCDNDomainImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("domain_name", d.Id())
}
