package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	protectStatusEnable = 1

	defaultBlockPageTemplate  = "default"
	customBlockPageTemplate   = "custom"
	redirectBlockPageTemplate = "redirect"
)

// @API WAF GET /v1/{project_id}/waf/instance/{instance_id}
// @API WAF PUT /v1/{project_id}/waf/instance/{instance_id}
// @API WAF DELETE /v1/{project_id}/waf/instance/{instance_id}
// @API WAF POST /v1/{project_id}/waf/instance
// @API WAF PUT /v1/{project_id}/waf/instance/{instance_id}/protect-status
// @API WAF PUT /v1/{project_id}/waf/instance/{instance_id}/access-status
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}
// @API WAF PATCH /v1/{project_id}/waf/policy/{policy_id}
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
				Default:  "prePaid",
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
			"pci_3ds": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"pci_dss": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"tls", "cipher"},
			},
			"tls": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"traffic_mark": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     domainTrafficMarkSchema(),
			},
			"website_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"forward_header_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Optional: true,
				Computed: true,
			},
			// Due to environmental and permission restrictions, the API test for the `access_status` parameter has not
			// been effective, and there is currently no testing in the test cases.
			"access_status": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"access_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func domainCustomPageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"http_return_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block_page_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func domainTimeoutSettingSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"connection_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"write_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func domainTrafficMarkSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_tags": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"session_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func domainServerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Required",
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
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
		Description:         d.Get("description").(string),
		ForwardHeaderMap:    buildHostForwardHeaderMapOpts(d),
		LbAlgorithm:         d.Get("lb_algorithm").(string),
		WebTag:              d.Get("website_name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
}

func buildHostForwardHeaderMapOpts(d *schema.ResourceData) map[string]string {
	if v, ok := d.GetOk("forward_header_map"); ok {
		return utils.ExpandToStringMap(v.(map[string]interface{}))
	}
	return nil
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
			Type:          server["type"].(string),
			Weight:        server["weight"].(int),
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
			"type":            server.Type,
			"weight":          server.Weight,
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

func flattenDomainTrafficMark(dm *domains.Domain) []map[string]interface{} {
	trafficMark := dm.TrafficMark
	return []map[string]interface{}{
		{
			"ip_tags":     trafficMark.Sip,
			"session_tag": trafficMark.Cookie,
			"user_tag":    trafficMark.Params,
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

func buildDomainHostFlag(d *schema.ResourceData) (*domains.Flag, error) {
	pci3ds := d.Get("pci_3ds").(bool)
	pciDss := d.Get("pci_dss").(bool)
	if !pci3ds && !pciDss {
		return nil, nil
	}

	// required tls="TLS v1.2" && cipher="cipher_2"
	if d.Get("tls").(string) != "TLS v1.2" || d.Get("cipher").(string) != "cipher_2" {
		return nil, fmt.Errorf("pci_3ds and pci_dss must be used together with tls and cipher. " +
			"Tls must be set to TLS v1.2, and cipher must be set to cipher_2")
	}
	return &domains.Flag{
		Pci3ds: strconv.FormatBool(pci3ds),
		PciDss: strconv.FormatBool(pciDss),
	}, nil
}

func updatePremiumHostTrafficMarkOpts(d *schema.ResourceData) *domains.TrafficMark {
	if v, ok := d.GetOk("traffic_mark"); ok {
		rawArray, isArray := v.([]interface{})
		if !isArray || len(rawArray) == 0 {
			return nil
		}

		raw, isMap := rawArray[0].(map[string]interface{})
		if !isMap {
			return nil
		}

		return &domains.TrafficMark{
			Sip:    utils.ExpandToStringList(raw["ip_tags"].([]interface{})),
			Cookie: raw["session_tag"].(string),
			Params: raw["user_tag"].(string),
		}
	}
	return nil
}

func updateWafDomain(wafClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	// Check whether ipv6_enable is valid.
	servers := buildWafDomainServers(d)
	ipv6Enable := d.Get("ipv6_enable").(bool)
	for _, server := range servers {
		if server.Type == "ipv6" && !ipv6Enable {
			return fmt.Errorf("when type in server contains IPv6 address, `ipv6_enable` should be configured to true")
		}
	}

	updateOpts := domains.UpdateOpts{
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	// Fields "certificate_id", "proxy", and "ipv6_enable" are valid only when they are used together with fields "server" in the update interface
	if d.HasChanges("certificate_id", "server", "proxy", "ipv6_enable") {
		updateOpts.CertificateId = d.Get("certificate_id").(string)
		updateOpts.CertificateName = d.Get("certificate_name").(string)
		updateOpts.Servers = servers
		updateOpts.Proxy = utils.Bool(d.Get("proxy").(bool))
		updateOpts.Ipv6Enable = utils.Bool(ipv6Enable)
	}

	if d.HasChanges("custom_page", "redirect_url") {
		updateOpts.BlockPage = buildUpdateDomainBlockPageOpts(d)
	}

	if d.HasChange("http2_enable") {
		updateOpts.Http2Enable = utils.Bool(d.Get("http2_enable").(bool))
	}

	if d.HasChange("timeout_settings") {
		updateOpts.TimeoutConfig = buildUpdateDomainTimeoutSettingOpts(d)
	}

	if d.HasChange("description") && !d.IsNewResource() {
		updateOpts.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("forward_header_map") && !d.IsNewResource() {
		updateOpts.ForwardHeaderMap = buildHostForwardHeaderMapOpts(d)
	}

	if d.HasChange("lb_algorithm") && !d.IsNewResource() {
		updateOpts.LbAlgorithm = utils.String(d.Get("lb_algorithm").(string))
	}

	if d.HasChange("website_name") && !d.IsNewResource() {
		updateOpts.WebTag = utils.String(d.Get("website_name").(string))
	}

	if d.HasChanges("tls", "cipher", "pci_3ds", "pci_dss") {
		updateOpts.Tls = d.Get("tls").(string)
		updateOpts.Cipher = d.Get("cipher").(string)
		// `pci_3ds` and `pci_dss` must be used together with `tls` and `cipher`.
		if d.HasChanges("pci_3ds", "pci_dss") {
			flag, err := buildDomainHostFlag(d)
			if err != nil {
				return err
			}
			updateOpts.Flag = flag
		}
	}

	if d.HasChange("traffic_mark") {
		updateOpts.TrafficMark = updatePremiumHostTrafficMarkOpts(d)
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

	if d.Get("protect_status").(int) != protectStatusEnable {
		if err := updateWafDomainProtectStatus(wafClient, d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("access_status"); ok {
		if err := updateWafDomainAccessStatus(wafClient, d); err != nil {
			return diag.Errorf("error update WAF domain `access_status` in creation operation: %s", err)
		}
	}

	return resourceWafDomainRead(ctx, d, meta)
}

func updateWafDomainProtectStatus(wafClient *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config) error {
	protectStatus := d.Get("protect_status").(int)
	epsID := cfg.GetEnterpriseProjectID(d)
	_, err := domains.UpdateProtectStatus(wafClient, protectStatus, d.Id(), epsID)
	if err != nil {
		return fmt.Errorf("error updating WAF domain protect status: %s", err)
	}
	return nil
}

func updateWafDomainAccessStatus(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl           = "v1/{project_id}/waf/instance/{instance_id}/access-status"
		inputAccessStatus = d.Get("access_status").(int)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: map[string]interface{}{
			"access_status": inputAccessStatus,
		},
	}

	requestResp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error update WAF domain access status: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	respAccessStatus := utils.PathSearch("access_status", respBody, float64(0)).(float64)
	if int(respAccessStatus) != inputAccessStatus {
		return fmt.Errorf("failed to update WAF domain access status: expected value: %v,"+
			" actual value: %v", inputAccessStatus, respAccessStatus)
	}

	return nil
}

func resourceWafDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	dm, err := domains.GetWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		// If the domain does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF domain")
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
		d.Set("access_code", dm.AccessCode),
		d.Set("protocol", dm.Protocol),
		d.Set("server", flattenDomainServerAttrs(dm)),
		d.Set("custom_page", flattenDomainCustomPage(dm)),
		d.Set("redirect_url", dm.BlockPage.RedirectUrl),
		d.Set("http2_enable", dm.Http2Enable),
		d.Set("timeout_settings", flattenDomainTimeoutSetting(dm)),
		d.Set("description", dm.Description),
		d.Set("forward_header_map", dm.ForwardHeaderMap),
		d.Set("lb_algorithm", dm.LbAlgorithm),
		d.Set("website_name", dm.WebTag),
		d.Set("cipher", dm.Cipher),
		d.Set("tls", dm.Tls),
		d.Set("traffic_mark", flattenDomainTrafficMark(dm)),
	)

	if dm.Flag.Pci3ds != "" {
		mErr = multierror.Append(mErr, d.Set("pci_3ds", utils.StringToBool(dm.Flag.Pci3ds)))
	}

	if dm.Flag.PciDss != "" {
		mErr = multierror.Append(mErr, d.Set("pci_dss", utils.StringToBool(dm.Flag.PciDss)))
	}

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

	if d.HasChanges("protect_status") {
		if err := updateWafDomainProtectStatus(wafClient, d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("access_status") {
		if err := updateWafDomainAccessStatus(wafClient, d); err != nil {
			return diag.Errorf("error update WAF domain `access_status` in update operation: %s", err)
		}
	}

	return resourceWafDomainRead(ctx, d, meta)
}

func updateWafDomainPolicyHost(d *schema.ResourceData, cfg *config.Config) error {
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating WAF client: %s", err)
	}

	oVal, nVal := d.GetChange("policy_id")
	newPolicyId := nVal.(string)
	oldPolicyId := oVal.(string)

	epsID := cfg.GetEnterpriseProjectID(d)
	updateHostsOpts := policies.UpdateHostsOpts{
		Hosts:               []string{d.Id()},
		EnterpriseProjectId: epsID,
	}
	log.Printf("[DEBUG] Bind WAF domain %s to policy %s", d.Id(), newPolicyId)

	if _, err := policies.UpdateHosts(client, newPolicyId, updateHostsOpts).Extract(); err != nil {
		return fmt.Errorf("error updating WAF policy hosts: %s", err)
	}

	if err := policies.DeleteWithEpsID(client, oldPolicyId, epsID).ExtractErr(); err != nil {
		// If other domains are using this policy, the deletion will fail.
		log.Printf("[WARN] error deleting WAF policy %s: %s", oldPolicyId, err)
	}
	return nil
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
		// If the domain does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF domain")
	}

	return nil
}
