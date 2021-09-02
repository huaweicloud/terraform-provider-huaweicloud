package waf

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	// protectStatusEnable 1: protection status enabled.
	protectStatusEnable = 1
)

func ResourceWafDedicatedDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafDedicatedDomainV1Create,
		Read:   resourceWafDedicatedDomainV1Read,
		Update: resourceWafDedicatedDomainV1Update,
		Delete: resourceWafDedicatedDomainV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cipher": {
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

// getCertificateNameById get certificate id from HuaweiCloud.
func getCertificateNameById(d *schema.ResourceData, meta interface{}) (string, error) {
	config := meta.(*config.Config)
	c, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return "", fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}
	if _, ok := d.GetOk("certificate_id"); ok {
		r, err := certificates.Get(c, d.Get("certificate_id").(string)).Extract()
		if err != nil {
			return "", fmtp.Errorf("error obtain WAF certificate name according id: %s, error:",
				d.Get("certificate_id").(string), err)
		}
		return r.Name, nil
	}
	return "", nil
}

// buildCreatePremiumHostOpts build the options for creating premium domains.
func buildCreatePremiumHostOpts(d *schema.ResourceData, meta interface{}) (*domains.CreateOpts, error) {
	certName, err := getCertificateNameById(d, meta)
	if err != nil {
		return nil, err
	}

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

	proxy := d.Get("proxy").(bool)
	opts := domains.CreateOpts{
		CertificateId:   d.Get("certificate_id").(string),
		CertificateName: certName,
		HostName:        d.Get("domain").(string),
		Proxy:           &proxy,
		PolicyId:        d.Get("policy_id").(string),
		Servers:         serverOpts,
	}

	logp.Printf("[DEBUG] build WAF dedicated domain creation options: %#v", serverOpts)
	return &opts, nil
}

// resourceWafDedicatedDomainV1Create create a premium domain name in HuaweiCloud.
func resourceWafDedicatedDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafDedicatedClient, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF dedicated Client: %s", err)
	}

	createOpts, err := buildCreatePremiumHostOpts(d, meta)
	if err != nil {
		return err
	}
	logp.Printf("[DEBUG] The options of creating WAF dedicated domain: %#v", createOpts)

	domain, err := domains.Create(wafDedicatedClient, *createOpts)
	if err != nil {
		return fmtp.Errorf("error creating WAF Domain: %s", err)
	}
	logp.Printf("[DEBUG] Waf dedicated domain has been created: %#v", domain)
	d.SetId(domain.Id)

	if d.Get("protect_status").(int) != protectStatusEnable {
		_, err = domains.UpdateProtectStatus(wafDedicatedClient, d.Id(), d.Get("protect_status").(int))
		if err != nil {
			// If the protection status update fails, it will be managed by terraform, and only print log.
			logp.Printf("[ERROR] error change the protection status of WAF dedicate domain %s: %s", d.Id(), err)
		}
	}

	return resourceWafDedicatedDomainV1Read(d, meta)
}

// buildDomainServerAttributes build the 'server' attribute after querying a domain.
func buildDomainServerAttribute(domain *domains.PremiumHost) []map[string]interface{} {
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

// buildDomainServerAttributes build the 'traffic_identifier' attribute after querying a domain.
func buildComplianceCertificationAttribute(domain *domains.PremiumHost) map[string]interface{} {
	f := domain.Flag

	pciDss, _ := strconv.ParseBool(f["pci_dss"])
	pci3ds, _ := strconv.ParseBool(f["pci_3ds"])
	return map[string]interface{}{
		"pci_dss": pciDss,
		"pci_3ds": pci3ds,
	}
}

// buildDomainServerAttributes build the 'compliance_certification' attribute after querying a domain.
func buildTrafficIdentifierAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.TrafficMark
	return map[string]interface{}{
		"ip_tag":      strings.Join(t.Sip, ","),
		"session_tag": t.Cookie,
		"user_tag":    t.Params,
	}
}

// buildAlarmPageAttribute build the 'alarm_page' attribute after querying a domain.
func buildAlarmPageAttribute(domain *domains.PremiumHost) map[string]interface{} {
	t := domain.BlockPage
	return map[string]interface{}{
		"template_name": t.Template,
		"redirect_url":  t.RedirectUrl,
	}
}

// resourceWafDedicatedDomainV1Read query a domain detail by id.
func resourceWafDedicatedDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	dm, err := domains.Get(wafClient, d.Id())
	if err != nil {
		return common.CheckDeleted(d, err, "Error obtain WAF dedicated domain information")
	}
	logp.Printf("[DEBUG] Get the WAF dedicated domain : %#v", dm)

	servers := buildDomainServerAttribute(dm)
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("domain", dm.HostName),
		d.Set("server", servers),
		d.Set("certificate_id", dm.CertificateId),
		d.Set("certificate_name", dm.CertificateName),
		d.Set("policy_id", dm.PolicyId),
		d.Set("proxy", dm.Proxy),
		d.Set("protect_status", dm.ProtectStatus),
		d.Set("access_status", dm.AccessStatus),
		d.Set("protocol", dm.Protocol),
		d.Set("tls", dm.Tls),
		d.Set("cipher", dm.Cipher),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.Errorf("error setting WAF fields: %s", err)
	}

	// The resources of compliance_certification, alarm_page and traffic_identifier may be empty.
	d.Set("compliance_certification", buildComplianceCertificationAttribute(dm))
	d.Set("traffic_identifier", buildTrafficIdentifierAttribute(dm))
	d.Set("alarm_page", buildAlarmPageAttribute(dm))

	return nil
}

// resourceWafDedicatedDomainV1Update modify some fields of domain: certificate_id, proxy, protect_status and policy_id.
func resourceWafDedicatedDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafDedicatedClient, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF dedicated Client: %s", err)
	}
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	if d.HasChanges("certificate_id", "proxy") {
		proxy := d.Get("proxy").(bool)

		certName, err := getCertificateNameById(d, meta)
		if err != nil {
			return err
		}

		updateOpts := domains.UpdatePremiumHostOpts{
			CertificateId:   d.Get("certificate_id").(string),
			CertificateName: certName,
			Proxy:           &proxy,
			Tls:             d.Get("tls").(string),
			Cipher:          d.Get("cipher").(string),
		}
		logp.Printf("[DEBUG] The options of update dedicated domain: %#v", updateOpts)

		_, err = domains.Update(wafDedicatedClient, d.Id(), updateOpts)
		if err != nil {
			return fmtp.Errorf("error updating WAF dedicated Domain: %s", err)
		}
	}

	if d.HasChanges("protect_status") {
		_, err = domains.UpdateProtectStatus(wafDedicatedClient, d.Id(), d.Get("protect_status").(int))
		if err != nil {
			return fmtp.Errorf("[ERROR] error change the protection status of WAF dedicate domain %s: %s",
				d.Id(), err)
		}
	}

	if d.HasChanges("policy_id") {
		oVal, nVal := d.GetChange("policy_id")
		policyId := nVal.(string)
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts: []string{d.Id()},
		}
		logp.Printf("[DEBUG] Bind Waf dedicated domain %s to policy %s", d.Id(), policyId)

		_, err = policies.UpdateHosts(wafClient, policyId, updateHostsOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating WAF Policy Hosts: %s", err)
		}

		// delete the old policy
		err = policies.Delete(wafClient, oVal.(string)).ExtractErr()
		if err != nil {
			// If other domains are using this policy, the deletion will fail.
			logp.Printf("[WARN] error deleting WAF Policy %s: %s", oVal.(string), err)
		}
	}
	return resourceWafDedicatedDomainV1Read(d, meta)
}

// resourceWafDedicatedDomainV1Delete delete a domain by id.
func resourceWafDedicatedDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafDedicatedClient, err := config.WafDedicatedV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	logp.Printf("[DEBUG] Delete WAF dedicated domain(keep_policy: %v).", d.Get("keep_policy"))
	_, err = domains.Delete(wafDedicatedClient, d.Id(), d.Get("keep_policy").(bool))
	if err != nil {
		return fmtp.Errorf("error deleting WAF dedicated domain: %s", err)
	}

	d.SetId("")
	return nil
}
