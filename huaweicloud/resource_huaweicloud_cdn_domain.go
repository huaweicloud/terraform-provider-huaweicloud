package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	cdnv1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v1/model"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
	Type:     schema.TypeList,
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
				ValidateFunc: validation.StringInSlice([]string{
					"full_url", "ignore_url_params", "del_args", "reserve_args",
				}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

func resourceCdnDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCdnDomainV1Create,
		Read:   resourceCdnDomainV1Read,
		Update: resourceCdnDomainV1Update,
		Delete: resourceCdnDomainV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
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
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
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
						"https_settings":             &httpsConfig,
						"retrieval_request_header":   &requestAndResponseHeader,
						"http_response_header":       &requestAndResponseHeader,
						"url_signing":                &authOpts,
						"force_redirect":             &forceRedirectAndCompress,
						"compress":                   &forceRedirectAndCompress,
						"cache_url_parameter_filter": &cacheUrlParameterFilter,
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"ttl_type": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
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

type WaitDomainStatus struct {
	ID      string
	Penging []string
	Target  []string
	Opts    *domains.ExtensionOpts
}

func getDomainSources(d *schema.ResourceData) []domains.SourcesOpts {
	var sourceRequests []domains.SourcesOpts

	sources := d.Get("sources").([]interface{})
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

func buildHttpsOpts(rawHttps []interface{}) *model.HttpPutBody {
	if len(rawHttps) != 1 {
		return nil
	}

	https := rawHttps[0].(map[string]interface{})
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

func buildUrlAuthOpts(rawUrlAuth []interface{}) *model.UrlAuth {
	if len(rawUrlAuth) != 1 {
		return nil
	}

	urlAuth := rawUrlAuth[0].(map[string]interface{})

	status := "off"
	if urlAuth["enabled"].(bool) {
		status = "on"
	}
	urlAuthOpts := model.UrlAuth{
		Status:     status,
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
	status := "off"
	if forceRedirect["enabled"].(bool) {
		status = "on"
	}
	forceRedirectOpts := model.ForceRedirectConfig{
		Status: status,
		Type:   utils.StringIgnoreEmpty(forceRedirect["type"].(string)),
	}

	return &forceRedirectOpts
}

func buildCompressOpts(rawCompress []interface{}) *model.Compress {
	if len(rawCompress) != 1 {
		return nil
	}

	compress := rawCompress[0].(map[string]interface{})
	status := "off"
	if compress["enabled"].(bool) {
		status = "on"
	}
	compressOpts := model.Compress{
		Status: status,
		Type:   utils.StringIgnoreEmpty(compress["type"].(string)),
	}

	return &compressOpts
}

func buildCacheUrlParameterFilterOpts(rawCacheUrlParameterFilter []interface{}) *model.CacheUrlParameterFilter {
	if len(rawCacheUrlParameterFilter) != 1 {
		return nil
	}

	cacheUrlParameterFilter := rawCacheUrlParameterFilter[0].(map[string]interface{})
	CacheUrlParameterFilterOpts := model.CacheUrlParameterFilter{
		Value: utils.StringIgnoreEmpty(cacheUrlParameterFilter["value"].(string)),
		Type:  utils.StringIgnoreEmpty(cacheUrlParameterFilter["type"].(string)),
	}

	return &CacheUrlParameterFilterOpts
}

func configOrUpdateConfigs(hcCdnClient *cdnv1.CdnClient, rawConfigs []interface{}, domainName, epsId string) error {
	if len(rawConfigs) != 1 {
		return nil
	}

	configs := rawConfigs[0].(map[string]interface{})
	ipv6Accelerate := 0
	if configs["ipv6_enable"].(bool) {
		ipv6Accelerate = 1
	}

	configsOpts := model.Configs{
		Https:                   buildHttpsOpts(configs["https_settings"].([]interface{})),
		OriginRequestHeader:     buildOriginRequestHeaderOpts(configs["retrieval_request_header"].([]interface{})),
		HttpResponseHeader:      buildHttpResponseHeaderOpts(configs["http_response_header"].([]interface{})),
		UrlAuth:                 buildUrlAuthOpts(configs["url_signing"].([]interface{})),
		OriginProtocol:          utils.StringIgnoreEmpty(configs["origin_protocol"].(string)),
		ForceRedirect:           buildForceRedirectOpts(configs["force_redirect"].([]interface{})),
		Compress:                buildCompressOpts(configs["compress"].([]interface{})),
		CacheUrlParameterFilter: buildCacheUrlParameterFilterOpts(configs["cache_url_parameter_filter"].([]interface{})),
		Ipv6Accelerate:          utils.Int32(int32(ipv6Accelerate)),
	}

	req := model.UpdateDomainFullConfigRequest{
		DomainName:          domainName,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		Body: &model.ModifyDomainConfigRequestBody{
			Configs: &configsOpts,
		},
	}

	_, err := hcCdnClient.UpdateDomainFullConfig(&req)
	if err != nil {
		return err
	}

	return nil
}

func buildCacheConfigRulesOpts(rawRules []interface{}) *[]model.Rules {
	if len(rawRules) < 1 {
		return nil
	}

	rulesOpts := make([]model.Rules, len(rawRules))
	for i, v := range rawRules {
		rule := v.(map[string]interface{})
		rulesOpts[i] = model.Rules{
			RuleType: int32(rule["rule_type"].(int)),
			Content:  utils.StringIgnoreEmpty(rule["content"].(string)),
			Ttl:      int32(rule["ttl"].(int)),
			TtlType:  int32(rule["ttl_type"].(int)),
			Priority: int32(rule["priority"].(int)),
		}
	}

	return &rulesOpts
}

func configOrUpdateCacheConfigOpts(hcCdnClient *cdnv1.CdnClient, rawCacheConfig []interface{}, domainId, epsId string) error {
	if len(rawCacheConfig) != 1 {
		return nil
	}
	cacheConfig := rawCacheConfig[0].(map[string]interface{})

	cacheConfigOpts := model.CacheConfigRequestBody{
		CacheConfig: &model.CacheConfigRequest{
			FollowOrigin: utils.Bool(cacheConfig["follow_origin"].(bool)),
			Rules:        buildCacheConfigRulesOpts(cacheConfig["rules"].([]interface{})),
		},
	}

	req := model.UpdateCacheRulesRequest{
		DomainId:            domainId,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
		Body:                &cacheConfigOpts,
	}

	_, err := hcCdnClient.UpdateCacheRules(&req)
	if err != nil {
		return err
	}

	return nil
}

func resourceCdnDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cdnClient, err := config.CdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	createOpts := &domains.CreateOpts{
		DomainName:          d.Get("name").(string),
		BusinessType:        d.Get("type").(string),
		Sources:             getDomainSources(d),
		ServiceArea:         d.Get("service_area").(string),
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := domains.Create(cdnClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN Domain: %s", err)
	}

	// Wait for CDN domain to become active again before continuing
	opts := getResourceExtensionOpts(d, config)
	timeout := d.Timeout(schema.TimeoutCreate)
	logp.Printf("[INFO] Waiting for CDN domain %s to become online.", v.ID)
	err = waitDomainOnlin(cdnClient, v.ID, opts, timeout)
	if err != nil {
		return err
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceCdnDomainV1Update(d, meta)
}

func waitforCDNV1DomainStatus(c *golangsdk.ServiceClient, waitstatus *WaitDomainStatus, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    waitstatus.Penging,
		Target:     waitstatus.Target,
		Refresh:    resourceCDNV1DomainRefreshFunc(c, waitstatus.ID, waitstatus.Opts),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error waiting for CDN domain %s to become %s: %s",
			waitstatus.ID, waitstatus.Target, err)
	}
	return nil
}

func resourceCDNV1DomainRefreshFunc(c *golangsdk.ServiceClient, id string, opts *domains.ExtensionOpts) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		domain, err := domains.Get(c, id, opts).Extract()
		if err != nil {
			return nil, "", err
		}

		// return DomainStatus attribute of CDN domain resource
		return domain, domain.DomainStatus, nil
	}
}

func flattenHttpsAttrs(https *model.HttpGetBody, privateKey string) []map[string]interface{} {
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
		"https_enabled":      https.HttpsStatus != nil && *https.HttpsStatus == "on",
		"http2_enabled":      https.Http2Status != nil && *https.Http2Status == "on",
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
		"enabled":     urlAuth.Status == "on",
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
		"enabled": forceRedirect.Status == "on",
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
		"enabled": compress.Status == "on",
	}

	return []map[string]interface{}{compressAttrs}
}

func flattenCacheUrlParameterFilterAttrs(cacheUrlParameterFilter *model.CacheUrlParameterFilter) []map[string]interface{} {
	if cacheUrlParameterFilter == nil {
		return nil
	}

	cacheUrlParameterFilterAttrs := map[string]interface{}{
		"value": cacheUrlParameterFilter.Value,
		"type":  cacheUrlParameterFilter.Type,
	}

	return []map[string]interface{}{cacheUrlParameterFilterAttrs}
}

func getConfigsAttrs(hcCdnClient *cdnv1.CdnClient, domainName, epsId, privateKey, urlAuthKey string) ([]map[string]interface{}, error) {
	req := model.ShowDomainFullConfigRequest{
		DomainName:          domainName,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
	}
	resp, err := hcCdnClient.ShowDomainFullConfig(&req)
	if err != nil {
		return nil, err
	}

	if resp.Configs == nil {
		return nil, fmtp.Errorf("unbale to find the configs of domain: %s", domainName)
	}

	configs := resp.Configs
	configsAttrs := map[string]interface{}{
		"https_settings":             flattenHttpsAttrs(configs.Https, privateKey),
		"retrieval_request_header":   flattenOriginRequestHeaderAttrs(configs.OriginRequestHeader),
		"http_response_header":       flattenHttpResponseHeaderAttrs(configs.HttpResponseHeader),
		"url_signing":                flattenUrlAuthAttrs(configs.UrlAuth, urlAuthKey),
		"origin_protocol":            configs.OriginProtocol,
		"force_redirect":             flattenForceRedirectAttrs(configs.ForceRedirect),
		"compress":                   flattenCompressAttrs(configs.Compress),
		"cache_url_parameter_filter": flattenCacheUrlParameterFilterAttrs(configs.CacheUrlParameterFilter),
		"ipv6_enable":                configs.Ipv6Accelerate != nil && *configs.Ipv6Accelerate == 1,
	}

	return []map[string]interface{}{configsAttrs}, nil
}

func getCacheAttrs(hcCdnClient *cdnv1.CdnClient, domainId, epsId string) ([]map[string]interface{}, error) {
	req := model.ShowCacheRulesRequest{
		DomainId:            domainId,
		EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
	}
	resp, err := hcCdnClient.ShowCacheRules(&req)
	if err != nil {
		return nil, err
	}

	if resp.CacheConfig == nil {
		return nil, fmtp.Errorf("unbale to find the cache config of domain: %s", domainId)
	}

	cacheConfig := resp.CacheConfig
	cacheAttrs := map[string]interface{}{
		"follow_origin": cacheConfig.FollowOrigin,
	}

	if cacheConfig.Rules == nil {
		return nil, fmtp.Errorf("unbale to find the cache config rules of domain: %s", domainId)
	}
	rules := make([]map[string]interface{}, len(*cacheConfig.Rules))
	for i, v := range *cacheConfig.Rules {
		rules[i] = map[string]interface{}{
			"rule_type": v.RuleType,
			"content":   v.Content,
			"ttl":       v.Ttl,
			"ttl_type":  v.TtlType,
			"priority":  v.Priority,
		}
	}

	cacheAttrs["rules"] = rules

	return []map[string]interface{}{cacheAttrs}, nil
}

func resourceCdnDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cdnClient, err := config.CdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	hcCdnClient, err := config.HcCdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	id := d.Id()
	epsId := GetEnterpriseProjectID(d, config)

	opts := getResourceExtensionOpts(d, config)
	v, err := domains.Get(cdnClient, id, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error reading CDN Domain: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved CDN domain %s: %+v", id, v)

	d.Set("name", v.DomainName)
	d.Set("type", v.BusinessType)
	d.Set("cname", v.CName)
	d.Set("domain_status", v.DomainStatus)
	d.Set("service_area", v.ServiceArea)

	// set sources
	sources := make([]map[string]interface{}, len(v.Sources))
	for i, source := range v.Sources {
		sources[i] = make(map[string]interface{})
		sources[i]["origin"] = source.IporDomain
		sources[i]["origin_type"] = source.OriginType
		sources[i]["active"] = source.ActiveStandby
	}
	d.Set("sources", sources)

	privateKey := d.Get("configs.0.https_settings.0.private_key").(string)
	urlAuthKey := d.Get("configs.0.url_signing.0.key").(string)
	configAttrs, err := getConfigsAttrs(hcCdnClient, v.DomainName, epsId, privateKey, urlAuthKey)
	if err != nil {
		return fmtp.Errorf("Error reading CDN Domain configs settings: %s", err)
	}

	d.Set("configs", configAttrs)

	cacheAttrs, err := getCacheAttrs(hcCdnClient, id, epsId)
	if err != nil {
		return fmtp.Errorf("Error reading CDN Domain cache settings: %s", err)
	}
	d.Set("cache_settings", cacheAttrs)

	return nil
}

func resourceCdnDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cdnClient, err := config.CdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	hcCdnClient, err := config.HcCdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	id := d.Id()
	domainName := d.Get("name").(string)
	epsId := GetEnterpriseProjectID(d, config)
	opts := getResourceExtensionOpts(d, config)
	timeout := d.Timeout(schema.TimeoutCreate)

	if d.HasChange("sources") && !d.IsNewResource() {
		opts := getResourceExtensionOpts(d, config)
		updateOpts := &domains.OriginOpts{
			Sources: getDomainSources(d),
		}

		if err = domains.Origin(cdnClient, id, opts, updateOpts).Err; err != nil {
			return fmtp.Errorf("Error updating CDN Domain orgin servers: %s", err)
		}

		// Wait for CDN domain to become active again before continuing
		logp.Printf("[INFO] Waiting for CDN domain %s to become online.", id)
		err = waitDomainOnlin(cdnClient, id, opts, timeout)
		if err != nil {
			return err
		}
	}

	if d.HasChange("configs") {
		err = configOrUpdateConfigs(hcCdnClient, d.Get("configs").([]interface{}), domainName, epsId)
		if err != nil {
			return fmtp.Errorf("Error updating CDN Domain configs settings: %s", err)
		}

		// Wait for CDN domain to become active again before continuing
		logp.Printf("[INFO] Waiting for CDN domain %s to become online.", id)
		err = waitDomainOnlin(cdnClient, id, opts, timeout)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cache_settings") {
		err = configOrUpdateCacheConfigOpts(hcCdnClient, d.Get("cache_settings").([]interface{}), id, epsId)
		if err != nil {
			return fmtp.Errorf("Error updating CDN Domain cache settings: %s", err)
		}

		// Wait for CDN domain to become active again before continuing
		logp.Printf("[INFO] Waiting for CDN domain %s to become online.", id)
		err = waitDomainOnlin(cdnClient, id, opts, timeout)
		if err != nil {
			return err
		}
	}

	return resourceCdnDomainV1Read(d, meta)
}

func waitDomainOnlin(cdnClient *golangsdk.ServiceClient, id string, opts *domains.ExtensionOpts, timeout time.Duration) error {
	wait := &WaitDomainStatus{
		ID:      id,
		Penging: []string{"configuring"},
		Target:  []string{"online"},
		Opts:    opts,
	}
	err := waitforCDNV1DomainStatus(cdnClient, wait, timeout)
	if err != nil {
		return fmtp.Errorf("Error waiting for CDN domain %s to become online: %s", id, err)
	}

	return nil
}

func resourceCdnDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	cdnClient, err := config.CdnV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN v1 client: %s", err)
	}

	id := d.Id()
	opts := getResourceExtensionOpts(d, config)
	timeout := d.Timeout(schema.TimeoutCreate)

	if d.Get("domain_status").(string) == "online" {
		// make sure the status has changed to offline
		logp.Printf("[INFO] Disable CDN domain %s.", id)
		if err = domains.Disable(cdnClient, id, opts).Err; err != nil {
			return fmtp.Errorf("Error disable  HuaweiCloud CDN Domain %s: %s", id, err)
		}

		logp.Printf("[INFO] Waiting for disabling CDN domain %s.", id)
		wait := &WaitDomainStatus{
			ID:      id,
			Penging: []string{"configuring", "online"},
			Target:  []string{"offline"},
			Opts:    opts,
		}

		err = waitforCDNV1DomainStatus(cdnClient, wait, timeout)
		if err != nil {
			return fmtp.Errorf("Error waiting for CDN domain %s to become offline: %s", id, err)
		}
	}

	logp.Printf("[INFO] Waiting for deleting CDN domain %s.", id)
	_, err = domains.Delete(cdnClient, id, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error deleting CDN Domain %s: %s", id, err)
	}

	// an API issue will be raised in ForceNew scene, so wait for a while
	time.Sleep(3 * time.Second) //lintignore:R018

	d.SetId("")
	return nil
}

func getResourceExtensionOpts(d *schema.ResourceData, config *config.Config) *domains.ExtensionOpts {
	epsID := GetEnterpriseProjectID(d, config)
	if epsID != "" {
		return &domains.ExtensionOpts{
			EnterpriseProjectId: epsID,
		}
	}

	return nil
}
