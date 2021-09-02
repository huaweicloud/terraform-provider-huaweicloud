package waf

import (
	"time"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	PROTOCOL_HTTP  = "HTTP"
	PROTOCOL_HTTPS = "HTTPS"
)

func ResourceWafDomainV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafDomainV1Create,
		Read:   resourceWafDomainV1Read,
		Update: resourceWafDomainV1Update,
		Delete: resourceWafDomainV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								PROTOCOL_HTTP, PROTOCOL_HTTPS,
							}, false),
						},
						"server_protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								PROTOCOL_HTTP, PROTOCOL_HTTPS,
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
				},
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
				ForceNew: true,
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

	logp.Printf("[DEBUG] build WAF domain ServerOpts: %#v", serverOpts)
	return serverOpts
}

func resourceWafDomainV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	proxy := d.Get("proxy").(bool)

	createOpts := domains.CreateOpts{
		HostName:        d.Get("domain").(string),
		CertificateId:   d.Get("certificate_id").(string),
		CertificateName: d.Get("certificate_name").(string),
		Servers:         buildWafDomainServers(d),
		Proxy:           &proxy,
	}
	logp.Printf("[DEBUG] CreateOpts: %#v", createOpts)

	domain, err := domains.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error creating WAF Domain: %s", err)
	}

	logp.Printf("[DEBUG] Waf domain created: %#v", domain)
	d.SetId(domain.Id)

	if v, ok := d.GetOk("policy_id"); ok {
		policyID := v.(string)
		hosts := []string{d.Id()}
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts: hosts,
		}

		logp.Printf("[DEBUG] Bind Waf domain %s to policy %s", d.Id(), policyID)
		_, err = policies.UpdateHosts(wafClient, policyID, updateHostsOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating WAF Policy Hosts: %s", err)
		}

		// delete the policy that was auto-created by domain
		err = policies.Delete(wafClient, domain.PolicyId).ExtractErr()
		if err != nil {
			logp.Printf("[WARN] error deleting WAF Policy %s: %s", domain.PolicyId, err)
		}
	}

	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	n, err := domains.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error obtain WAF domain information")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("domain", n.HostName),
		d.Set("certificate_id", n.CertificateId),
		d.Set("certificate_name", n.CertificateName),
		d.Set("policy_id", n.PolicyId),
		d.Set("proxy", n.Proxy),
		d.Set("protect_status", n.ProtectStatus),
		d.Set("access_status", n.AccessStatus),
		d.Set("protocol", n.Protocol),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("error setting WAF fields: %s", err)
	}

	servers := make([]map[string]interface{}, len(n.Servers))
	for i, server := range n.Servers {
		servers[i] = map[string]interface{}{
			"client_protocol": server.FrontProtocol,
			"server_protocol": server.BackProtocol,
			"address":         server.Address,
			"port":            server.Port,
		}
	}
	d.Set("server", servers)

	return nil
}

func resourceWafDomainV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	if d.HasChanges("certificate_id", "server", "proxy") {
		proxy := d.Get("proxy").(bool)

		updateOpts := domains.UpdateOpts{
			CertificateId:   d.Get("certificate_id").(string),
			CertificateName: d.Get("certificate_name").(string),
			Servers:         buildWafDomainServers(d),
			Proxy:           &proxy,
		}

		logp.Printf("[DEBUG] updateOpts: %#v", updateOpts)

		_, err = domains.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating WAF Domain: %s", err)
		}
	}
	if d.HasChanges("proxy_id") {
		oVal, nVal := d.GetChange("proxy_id")
		policyId := nVal.(string)
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts: []string{policyId},
		}

		logp.Printf("[DEBUG] Bind Waf domain %s to policy %s", d.Id(), policyId)
		_, err = policies.UpdateHosts(wafClient, policyId, updateHostsOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating WAF Policy Hosts: %s", err)
		}

		// delete the old policy
		err = policies.Delete(wafClient, oVal.(string)).ExtractErr()
		if err != nil {
			logp.Printf("[WARN] error deleting WAF Policy %s: %s", oVal.(string), err)
		}
	}
	return resourceWafDomainV1Read(d, meta)
}

func resourceWafDomainV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	delOpts := domains.DeleteOpts{
		KeepPolicy: d.Get("keep_policy").(bool),
	}
	logp.Printf("[DEBUG] delete WAF Domain: %#v", d.Get("keep_policy"))
	err = domains.Delete(wafClient, d.Id(), delOpts).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting WAF Domain: %s", err)
	}

	d.SetId("")
	return nil
}
