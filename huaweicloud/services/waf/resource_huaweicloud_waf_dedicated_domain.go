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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

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

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}
// @API WAF PATCH /v1/{project_id}/waf/policy/{policy_id}
// @API WAF GET /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF PUT /v1/{project_id}/premium-waf/host/{host_id}/protect-status
// @API WAF GET /v1/{project_id}/premium-waf/host/{host_id}
// @API WAF PUT /v1/{project_id}/premium-waf/host/{host_id}
// @API WAF DELETE /v1/{project_id}/premium-waf/host/{host_id}
// @API WAF POST /v1/{project_id}/premium-waf/host
func ResourceWafDedicatedDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedDomainCreate,
		ReadContext:   resourceWafDedicatedDomainRead,
		UpdateContext: resourceWafDedicatedDomainUpdate,
		DeleteContext: resourceWafDedicatedDomainDelete,

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
				ForceNew: true,
				MaxItems: 80,
				Elem:     dedicatedDomainServerSchema(),
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"tls": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TLS v1.0", "TLS v1.1", "TLS v1.2",
				}, false),
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cipher_1",
					"cipher_2",
					"cipher_3",
					"cipher_4",
					"cipher_default",
				}, false),
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"website_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"description": {
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
			"connection_protection": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     dedicatedDomainConnectionProtectionSchema(),
			},
			"timeout_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     domainTimeoutSettingSchema(),
			},
			"traffic_mark": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     dedicatedDomainTrafficMarkSchema(),
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_page": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compliance_certification": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
			"traffic_identifier": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dedicatedDomainServerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"server_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	return &sc
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

func dedicatedDomainConnectionProtectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"error_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"initial_downtime": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"multiplier_for_consecutive_breakdowns": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"pending_url_request_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

func dedicatedDomainTrafficMarkSchema() *schema.Resource {
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

func getCertificateNameById(d *schema.ResourceData, cfg *config.Config) (string, error) {
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating WAF client: %s", err)
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		epsID := cfg.GetEnterpriseProjectID(d)
		certificateId := v.(string)
		r, err := certificates.GetWithEpsID(client, certificateId, epsID).Extract()
		if err != nil {
			return "", fmt.Errorf("error retrieving WAF certificate name according ID: %s, error: %s", certificateId, err)
		}
		return r.Name, nil
	}
	return "", nil
}

func buildCreatePremiumHostOpts(d *schema.ResourceData, cfg *config.Config, certName string) *domains.CreateOpts {
	return &domains.CreateOpts{
		CertificateId:       d.Get("certificate_id").(string),
		CertificateName:     certName,
		HostName:            d.Get("domain").(string),
		Proxy:               utils.Bool(d.Get("proxy").(bool)),
		PolicyId:            d.Get("policy_id").(string),
		Servers:             buildCreatePremiumHostServerOpts(d),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		BlockPage:           buildPremiumHostBlockPageOpts(d),
		ForwardHeaderMap:    buildHostForwardHeaderMapOpts(d),
		Description:         d.Get("description").(string),
	}
}

func buildPremiumHostBlockPageOpts(d *schema.ResourceData) *domains.BlockPage {
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

func buildHostForwardHeaderMapOpts(d *schema.ResourceData) map[string]string {
	if v, ok := d.GetOk("forward_header_map"); ok {
		return utils.ExpandToStringMap(v.(map[string]interface{}))
	}
	return nil
}

func buildPremiumHostConnectionProtectionOpts(d *schema.ResourceData) *domains.CircuitBreaker {
	if v, ok := d.GetOk("connection_protection"); ok {
		rawArray, isArray := v.([]interface{})
		if !isArray || len(rawArray) == 0 {
			return nil
		}

		raw, isMap := rawArray[0].(map[string]interface{})
		if !isMap {
			return nil
		}

		return &domains.CircuitBreaker{
			DeadNum:          utils.Int(raw["error_threshold"].(int)),
			DeadRatio:        utils.Float64(raw["error_percentage"].(float64)),
			BlockTime:        utils.Int(raw["initial_downtime"].(int)),
			SuperpositionNum: utils.Int(raw["multiplier_for_consecutive_breakdowns"].(int)),
			SuspendNum:       utils.Int(raw["pending_url_request_threshold"].(int)),
			SusBlockTime:     utils.Int(raw["duration"].(int)),
			Switch:           utils.Bool(raw["status"].(bool)),
		}
	}
	return nil
}

func buildPremiumHostTimeoutSettingOpts(d *schema.ResourceData) *domains.TimeoutConfig {
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

func buildPremiumHostTrafficMarkOpts(d *schema.ResourceData) *domains.TrafficMark {
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

func buildCreatePremiumHostServerOpts(d *schema.ResourceData) []domains.Server {
	servers := d.Get("server").([]interface{})
	serverOpts := make([]domains.Server, len(servers))
	for i, v := range servers {
		s := v.(map[string]interface{})
		serverOpts[i] = domains.Server{
			FrontProtocol: s["client_protocol"].(string),
			BackProtocol:  s["server_protocol"].(string),
			Address:       s["address"].(string),
			Port:          s["port"].(int),
			Type:          s["type"].(string),
			VpcId:         s["vpc_id"].(string),
		}
	}
	return serverOpts
}

func resourceWafDedicatedDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dedicatedClient, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	certName, err := getCertificateNameById(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := buildCreatePremiumHostOpts(d, cfg, certName)
	domain, err := domains.Create(dedicatedClient, *createOpts)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated domain: %s", err)
	}
	d.SetId(domain.Id)

	if err := updateWafDedicatedDomain(dedicatedClient, d, cfg); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("protect_status").(int) != protectStatusEnable {
		if err := updateWafDedicatedDomainProtectStatus(dedicatedClient, d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWafDedicatedDomainRead(ctx, d, meta)
}

func updateWafDedicatedDomainProtectStatus(dedicatedClient *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config) error {
	protectStatus := d.Get("protect_status").(int)
	epsID := cfg.GetEnterpriseProjectID(d)
	_, err := domains.UpdateProtectStatusWithWpsID(dedicatedClient, protectStatus, d.Id(), epsID)
	if err != nil {
		return fmt.Errorf("error updating WAF dedicated domain protect status: %s", err)
	}
	return nil
}

func updateWafDedicatedDomain(dedicatedClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	updateOpts := domains.UpdatePremiumHostOpts{
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	if d.HasChanges("tls", "cipher", "pci_3ds", "pci_dss") {
		updateOpts.Tls = d.Get("tls").(string)
		updateOpts.Cipher = d.Get("cipher").(string)
		// `pci_3ds` and `pci_dss` must be used together with `tls` and `cipher`.
		if d.HasChanges("pci_3ds", "pci_dss") {
			flag, err := buildHostFlag(d)
			if err != nil {
				return err
			}
			updateOpts.Flag = flag
		}
	}

	if d.HasChange("website_name") {
		updateOpts.WebTag = utils.String(d.Get("website_name").(string))
	}

	if d.HasChange("connection_protection") {
		updateOpts.CircuitBreaker = buildPremiumHostConnectionProtectionOpts(d)
	}

	if d.IsNewResource() && updateOpts.CircuitBreaker == nil {
		// When creating new resource, if field `connection_protection` is empty, make configure `switch` to false.
		// Otherwise, when querying the details interface, the corresponding field will be empty.
		updateOpts.CircuitBreaker = &domains.CircuitBreaker{
			Switch: utils.Bool(false),
		}
	}

	if d.HasChange("timeout_settings") {
		updateOpts.TimeoutConfig = buildPremiumHostTimeoutSettingOpts(d)
	}

	if d.HasChange("traffic_mark") {
		updateOpts.TrafficMark = buildPremiumHostTrafficMarkOpts(d)
	}

	if d.HasChange("proxy") && !d.IsNewResource() {
		updateOpts.Proxy = utils.Bool(d.Get("proxy").(bool))
	}

	if d.HasChange("certificate_id") && !d.IsNewResource() {
		if v, ok := d.GetOk("certificate_id"); ok {
			certName, err := getCertificateNameById(d, cfg)
			if err != nil {
				return err
			}
			updateOpts.CertificateName = certName
			updateOpts.CertificateId = v.(string)
		}
	}

	if d.HasChanges("custom_page", "redirect_url") && !d.IsNewResource() {
		updateOpts.BlockPage = buildPremiumHostBlockPageOpts(d)
	}

	if d.HasChange("description") && !d.IsNewResource() {
		updateOpts.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("forward_header_map") && !d.IsNewResource() {
		updateOpts.ForwardHeaderMap = buildHostForwardHeaderMapOpts(d)
	}

	_, err := domains.Update(dedicatedClient, d.Id(), updateOpts)
	if err != nil {
		return fmt.Errorf("error updating WAF dedicated domain: %s", err)
	}
	return nil
}

func buildHostFlag(d *schema.ResourceData) (*domains.Flag, error) {
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

func flattenDomainServerAttribute(domain *domains.PremiumHost) []map[string]interface{} {
	servers := make([]map[string]interface{}, 0, len(domain.Servers))
	for _, s := range domain.Servers {
		servers = append(servers, map[string]interface{}{
			"client_protocol": s.FrontProtocol,
			"server_protocol": s.BackProtocol,
			"address":         s.Address,
			"port":            s.Port,
			"type":            s.Type,
			"vpc_id":          s.VpcId,
		})
	}
	return servers
}

func flattenBlockPageCustomPage(domain *domains.PremiumHost) []map[string]interface{} {
	if domain.BlockPage.Template != customBlockPageTemplate {
		return nil
	}

	customPage := domain.BlockPage.CustomPage
	return []map[string]interface{}{
		{
			"http_return_code": customPage.StatusCode,
			"block_page_type":  customPage.ContentType,
			"page_content":     customPage.Content,
		},
	}
}

func flattenConnectionProtection(domain *domains.PremiumHost) []map[string]interface{} {
	circuitBreaker := domain.CircuitBreaker
	return []map[string]interface{}{
		{
			"error_threshold":                       circuitBreaker.DeadNum,
			"error_percentage":                      circuitBreaker.DeadRatio,
			"initial_downtime":                      circuitBreaker.BlockTime,
			"multiplier_for_consecutive_breakdowns": circuitBreaker.SuperpositionNum,
			"pending_url_request_threshold":         circuitBreaker.SuspendNum,
			"duration":                              circuitBreaker.SusBlockTime,
			"status":                                circuitBreaker.Switch,
		},
	}
}

func flattenTimeoutSetting(domain *domains.PremiumHost) []map[string]interface{} {
	timeoutConfig := domain.TimeoutConfig
	return []map[string]interface{}{
		{
			"connection_timeout": timeoutConfig.ConnectTimeout,
			"read_timeout":       timeoutConfig.ReadTimeout,
			"write_timeout":      timeoutConfig.SendTimeout,
		},
	}
}

func flattenTrafficMark(domain *domains.PremiumHost) []map[string]interface{} {
	trafficMark := domain.TrafficMark
	return []map[string]interface{}{
		{
			"ip_tags":     trafficMark.Sip,
			"session_tag": trafficMark.Cookie,
			"user_tag":    trafficMark.Params,
		},
	}
}

func flattenComplianceCertificationAttribute(domain *domains.PremiumHost) map[string]interface{} {
	f := domain.Flag

	pciDss, _ := strconv.ParseBool(f["pci_dss"])
	pci3ds, _ := strconv.ParseBool(f["pci_3ds"])
	return map[string]interface{}{
		"pci_dss": pciDss,
		"pci_3ds": pci3ds,
	}
}

func flattenTrafficIdentifierAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.TrafficMark
	return map[string]interface{}{
		"ip_tag":      strings.Join(t.Sip, ","),
		"session_tag": t.Cookie,
		"user_tag":    t.Params,
	}
}

func flattenAlarmPageAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.BlockPage
	return map[string]interface{}{
		"template_name": t.Template,
		"redirect_url":  t.RedirectUrl,
	}
}

func resourceWafDedicatedDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dedicatedClient, err := cfg.WafDedicatedV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	dm, err := domains.GetWithEpsID(dedicatedClient, d.Id(), epsID)
	if err != nil {
		// If the dedicated domain does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF dedicated domain")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("domain", dm.HostName),
		d.Set("server", flattenDomainServerAttribute(dm)),
		d.Set("certificate_id", dm.CertificateId),
		d.Set("certificate_name", dm.CertificateName),
		d.Set("policy_id", dm.PolicyId),
		d.Set("proxy", dm.Proxy),
		d.Set("protect_status", dm.ProtectStatus),
		d.Set("access_status", dm.AccessStatus),
		d.Set("protocol", dm.Protocol),
		d.Set("tls", dm.Tls),
		d.Set("cipher", dm.Cipher),
		d.Set("website_name", dm.WebTag),
		d.Set("description", dm.Description),
		d.Set("forward_header_map", dm.ForwardHeaderMap),
		d.Set("custom_page", flattenBlockPageCustomPage(dm)),
		d.Set("redirect_url", dm.BlockPage.RedirectUrl),
		d.Set("connection_protection", flattenConnectionProtection(dm)),
		d.Set("timeout_settings", flattenTimeoutSetting(dm)),
		d.Set("traffic_mark", flattenTrafficMark(dm)),
		d.Set("compliance_certification", flattenComplianceCertificationAttribute(dm)),
		d.Set("traffic_identifier", flattenTrafficIdentifierAttribute(dm)),
		d.Set("alarm_page", flattenAlarmPageAttribute(dm)),
	)

	if dm.Flag["pci_3ds"] != "" {
		mErr = multierror.Append(mErr, d.Set("pci_3ds", utils.StringToBool(dm.Flag["pci_3ds"])))
	}

	if dm.Flag["pci_dss"] != "" {
		mErr = multierror.Append(mErr, d.Set("pci_dss", utils.StringToBool(dm.Flag["pci_dss"])))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafDedicatedDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dedicatedClient, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	if err := updateWafDedicatedDomain(dedicatedClient, d, cfg); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("protect_status") {
		if err := updateWafDedicatedDomainProtectStatus(dedicatedClient, d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("policy_id") {
		if err := updateWafDomainPolicyHost(d, cfg); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceWafDedicatedDomainRead(ctx, d, meta)
}

func resourceWafDedicatedDomainDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dedicatedClient, err := cfg.WafDedicatedV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF dedicated client: %s", err)
	}

	keepPolicy := d.Get("keep_policy").(bool)
	epsID := cfg.GetEnterpriseProjectID(d)
	_, err = domains.DeleteWithEpsID(dedicatedClient, keepPolicy, d.Id(), epsID)
	if err != nil {
		// If the dedicated domain does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF dedicated domain")
	}
	return nil
}
