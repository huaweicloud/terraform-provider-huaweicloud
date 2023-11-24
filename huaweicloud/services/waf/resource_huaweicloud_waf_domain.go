package waf

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

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

func updateWafDomain(wafClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	if d.HasChanges("certificate_id", "server", "proxy") {
		updateOpts := domains.UpdateOpts{
			CertificateId:       d.Get("certificate_id").(string),
			CertificateName:     d.Get("certificate_name").(string),
			Servers:             buildWafDomainServers(d),
			Proxy:               utils.Bool(d.Get("proxy").(bool)),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}

		if _, err := domains.Update(wafClient, d.Id(), updateOpts).Extract(); err != nil {
			return fmt.Errorf("error updating WAF domain: %s", err)
		}
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

	if v, ok := d.GetOk("policy_id"); ok {
		policyID := v.(string)
		hosts := []string{d.Id()}
		epsID := cfg.GetEnterpriseProjectID(d)
		updateHostsOpts := policies.UpdateHostsOpts{
			Hosts:               hosts,
			EnterpriseProjectId: epsID,
		}

		_, err = policies.UpdateHosts(wafClient, policyID, updateHostsOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF policy Hosts: %s", err)
		}

		// delete the policy that was auto-created by domain
		err = policies.DeleteWithEpsID(wafClient, domain.PolicyId, epsID).ExtractErr()
		if err != nil {
			log.Printf("[WARN] error deleting WAF policy %s: %s", domain.PolicyId, err)
		}
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

	d.SetId("")
	return nil
}
