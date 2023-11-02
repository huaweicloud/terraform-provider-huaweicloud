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
)

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
				Elem: &schema.Resource{
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
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntBetween(0, 65535),
							Required:     true,
							ForceNew:     true,
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
				},
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
	}
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

func updateWafDedicatedDomainPolicyHost(d *schema.ResourceData, cfg *config.Config) error {
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
	log.Printf("[DEBUG] Bind WAF dedicated domain %s to policy %s", d.Id(), newPolicyId)

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
		d.Set("compliance_certification", flattenComplianceCertificationAttribute(dm)),
		d.Set("traffic_identifier", flattenTrafficIdentifierAttribute(dm)),
		d.Set("alarm_page", flattenAlarmPageAttribute(dm)),
	)

	if dm.Flag["pci_3ds"] != "" {
		pci3ds, err := strconv.ParseBool(dm.Flag["pci_3ds"])
		if err != nil {
			log.Printf("[WARN] error parse bool pci 3ds, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_3ds", pci3ds))
	}

	if dm.Flag["pci_dss"] != "" {
		pciDss, err := strconv.ParseBool(dm.Flag["pci_dss"])
		if err != nil {
			log.Printf("[WARN] error parse bool pci dss, %s", err)
		}
		mErr = multierror.Append(mErr, d.Set("pci_dss", pciDss))
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
		if err := updateWafDedicatedDomainPolicyHost(d, cfg); err != nil {
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
		return diag.Errorf("error deleting WAF dedicated domain: %s", err)
	}
	return nil
}
