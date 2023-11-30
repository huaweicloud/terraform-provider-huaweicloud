package waf

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var PaidType = "prePaid"

const (
	protocolHTTP  = "HTTP"
	protocolHTTPS = "HTTPS"
)

func ResourceWafDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDomainCreate,
		ReadContext:   resourceWafDomainRead,
		UpdateContext: resourceWafDomainUpdate,
		DeleteContext: resourceWafDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     domainServerSchema(),
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keep_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  PaidType,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"custom_page": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     domainCustomPageSchema(),
			},
			"redirect_url": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"custom_page",
				},
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"timeout_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     domainTimeoutSettingSchema(),
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func domainServerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					protocolHTTP, protocolHTTPS,
				}, false),
			},
			"server_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					protocolHTTP, protocolHTTPS,
				}, false),
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 65535),
				Required:     true,
			},
		},
	}
	return &sc
}

func buildCreateDomainHostOpts(d *schema.ResourceData, cfg *config.Config) *domains.CreateOpts {
	return &domains.CreateOpts{
		HostName:            d.Get("domain").(string),
		CertificateId:       d.Get("certificate_id").(string),
		CertificateName:     d.Get("certificate_name").(string),
		Servers:             buildWafDomainServers(d),
		Proxy:               utils.Bool(d.Get("proxy").(bool)),
		PaidType:            d.Get("charging_mode").(string),
		PolicyId:            d.Get("policy_id").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
}

func buildWafDomainServers(d *schema.ResourceData) []domains.ServerOpts {
	servers := d.Get("server").([]interface{})

	serverOpts := make([]domains.ServerOpts, len(servers))
	for i, v := range servers {
		server := v.(map[string]interface{})
		serverOpts[i] = domains.ServerOpts{
			FrontProtocol: server["client_protocol"].(string),
			BackProtocol:  server["server_protocol"].(string),
			Address:       server["address"].(string),
			Port:          server["port"].(int),
		}
	}

	return serverOpts
}

func flattenDomainServerAttrs(dm *domains.Domain) []map[string]interface{} {
	servers := make([]map[string]interface{}, len(dm.Servers))
	for i, server := range dm.Servers {
		servers[i] = map[string]interface{}{
			"client_protocol": server.FrontProtocol,
			"server_protocol": server.BackProtocol,
			"address":         server.Address,
			"port":            server.Port,
		}
	}
	return servers
}

func flattenDomainCustomPage(dm *domains.Domain) []map[string]interface{} {
	if dm.BlockPage.Template != customBlockPageTemplate {
		return nil
	}

	customPage := dm.BlockPage.CustomPage
	return []map[string]interface{}{
		{
			"http_return_code": customPage.StatusCode,
			"block_page_type":  customPage.ContentType,
			"page_content":     customPage.Content,
		},
	}
}

func flattenDomainTimeoutSetting(dm *domains.Domain) []map[string]interface{} {
	timeoutConfig := dm.TimeoutConfig
	return []map[string]interface{}{
		{
			"connection_timeout": timeoutConfig.ConnectTimeout,
			"read_timeout":       timeoutConfig.ReadTimeout,
			"write_timeout":      timeoutConfig.SendTimeout,
		},
	}
}

func buildUpdateDomainBlockPageOpts(d *schema.ResourceData) *domains.BlockPage {
	if v, ok := d.GetOk("redirect_url"); ok {
		return &domains.BlockPage{
			Template:    redirectBlockPageTemplate,
			RedirectUrl: v.(string),
		}
	}

	if v, ok := d.GetOk("custom_page"); ok {
		rawArray, isArray := v.([]interface{})
		if !isArray || len(rawArray) == 0 {
			return nil
		}

		raw, isMap := rawArray[0].(map[string]interface{})
		if !isMap {
			return nil
		}
		return &domains.BlockPage{
			Template: customBlockPageTemplate,
			CustomPage: &domains.CustomPage{
				StatusCode:  raw["http_return_code"].(string),
				ContentType: raw["block_page_type"].(string),
				Content:     raw["page_content"].(string),
			},
		}
	}

	return &domains.BlockPage{
		Template: defaultBlockPageTemplate,
	}
}

func buildUpdateDomainTimeoutSettingOpts(d *schema.ResourceData) *domains.TimeoutConfig {
	if v, ok := d.GetOk("timeout_settings"); ok {
		rawArray, isArray := v.([]interface{})
		if !isArray || len(rawArray) == 0 {
			return nil
		}

		raw, isMap := rawArray[0].(map[string]interface{})
		if !isMap {
			return nil
		}

		return &domains.TimeoutConfig{
			ConnectTimeout: utils.Int(raw["connection_timeout"].(int)),
			ReadTimeout:    utils.Int(raw["read_timeout"].(int)),
			SendTimeout:    utils.Int(raw["write_timeout"].(int)),
		}
	}
	return nil
}

func updateWafDomain(wafClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	updateOpts := domains.UpdateOpts{
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	if d.HasChanges("certificate_id", "server", "proxy") {
		updateOpts.CertificateId = d.Get("certificate_id").(string)
		updateOpts.CertificateName = d.Get("certificate_name").(string)
		updateOpts.Servers = buildWafDomainServers(d)
		updateOpts.Proxy = utils.Bool(d.Get("proxy").(bool))
	}

	if d.HasChanges("custom_page", "redirect_url") {
		updateOpts.BlockPage = buildUpdateDomainBlockPageOpts(d)
	}

	if d.HasChange("http2_enable") {
		updateOpts.Http2Enable = utils.Bool(d.Get("http2_enable").(bool))
	}

	if d.HasChange("ipv6_enable") {
		updateOpts.Ipv6Enable = utils.Bool(d.Get("ipv6_enable").(bool))
	}

	if d.HasChange("timeout_settings") {
		updateOpts.TimeoutConfig = buildUpdateDomainTimeoutSettingOpts(d)
	}

	if _, err := domains.Update(wafClient, d.Id(), updateOpts).Extract(); err != nil {
		return fmt.Errorf("error updating WAF domain: %s", err)
	}

	return nil
}

func resourceWafDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createOpts := buildCreateDomainHostOpts(d, cfg)

	domain, err := domains.Create(wafClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF domain: %s", err)
	}
	d.SetId(domain.Id)

	if err := updateWafDomain(wafClient, d, cfg); err != nil {
		return diag.FromErr(err)
	}

	return resourceWafDomainRead(ctx, d, meta)
}

func resourceWafDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	dm, err := domains.GetWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain WAF domain information")
	}

	// charging_mode not returned by API
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("domain", dm.HostName),
		d.Set("certificate_id", dm.CertificateId),
		d.Set("certificate_name", dm.CertificateName),
		d.Set("policy_id", dm.PolicyId),
		d.Set("proxy", dm.Proxy),
		d.Set("protect_status", dm.ProtectStatus),
		d.Set("access_status", dm.AccessStatus),
		d.Set("protocol", dm.Protocol),
		d.Set("server", flattenDomainServerAttrs(dm)),
		d.Set("custom_page", flattenDomainCustomPage(dm)),
		d.Set("redirect_url", dm.BlockPage.RedirectUrl),
		d.Set("http2_enable", dm.Http2Enable),
		d.Set("timeout_settings", flattenDomainTimeoutSetting(dm)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting WAF domain fields: %s", err)
	}

	return nil
}

func resourceWafDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if err := updateWafDomain(wafClient, d, cfg); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("policy_id") {
		if err := updateWafDomainPolicyHost(d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWafDomainRead(ctx, d, meta)
}

func resourceWafDomainDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	delOpts := domains.DeleteOpts{
		KeepPolicy:          d.Get("keep_policy").(bool),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
	err = domains.Delete(wafClient, d.Id(), delOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting WAF domain: %s", err)
	}

	return nil
}
