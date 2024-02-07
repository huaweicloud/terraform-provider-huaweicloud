package dns

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dns/v2/resolverrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DNS DELETE /v2.1/resolverrules/{resolverrule_id}
// @API DNS GET /v2.1/resolverrules/{resolverrule_id}
// @API DNS PUT /v2.1/resolverrules/{resolverrule_id}
// @API DNS POST /v2.1/resolverrules
func ResourceDNSResolverRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSResolverRuleCreate,
		ReadContext:   resourceDNSResolverRuleRead,
		UpdateContext: resourceDNSResolverRuleUpdate,
		DeleteContext: resourceDNSResolverRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSResolverRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	ipAddressList := d.Get("ip_addresses").([]interface{})
	ipAddresses := make([]resolverrule.IPAddress, len(ipAddressList))
	for i, a := range ipAddressList {
		address := a.(map[string]interface{})
		ipAddresses[i] = resolverrule.IPAddress{
			IP: address["ip"].(string),
		}
	}
	opts := resolverrule.CreateOpts{
		Name:        d.Get("name").(string),
		DomainName:  d.Get("domain_name").(string),
		EndpointID:  d.Get("endpoint_id").(string),
		IPAddresses: ipAddresses,
	}
	rule, err := resolverrule.Create(dnsClient, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS resolver rule: %s", err)
	}

	d.SetId(rule.ID)
	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) to become available", rule.ID)
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      waitForDNSResolverRule(dnsClient, rule.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS resolver rule (%s) to become ACTIVE for creation: %s",
			rule.ID, err)
	}
	return resourceDNSResolverRuleRead(ctx, d, meta)
}

func resourceDNSResolverRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	id := d.Id()
	body, err := resolverrule.Get(dnsClient, id).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS resolver rule")
	}

	rule := body.ResolverRule
	vpcs := make([]map[string]interface{}, len(rule.Routers))
	for i, r := range rule.Routers {
		vpcs[i] = map[string]interface{}{
			"vpc_id":     r.RouterID,
			"vpc_region": r.RouterRegion,
			"status":     r.Status,
		}
	}

	ipAddresses := make([]map[string]interface{}, len(rule.IPAddresses))
	for i, r := range rule.IPAddresses {
		ipAddresses[i] = map[string]interface{}{
			"ip": r.IP,
		}
	}
	mErr := multierror.Append(nil,
		d.Set("name", rule.Name),
		d.Set("domain_name", rule.DomainName),
		d.Set("endpoint_id", rule.EndpointID),
		d.Set("status", rule.Status),
		d.Set("rule_type", rule.RuleType),
		d.Set("vpcs", vpcs),
		d.Set("ip_addresses", ipAddresses),
		d.Set("created_at", rule.CreatedAt),
		d.Set("updated_at", rule.UpdatedAt),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSResolverRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	name := d.Get("name").(string)
	address := d.Get("ip_addresses").([]interface{})
	ipAddressList := make([]resolverrule.IPAddress, len(address))
	for i, a := range address {
		addr := a.(map[string]interface{})
		ipAddressList[i] = resolverrule.IPAddress{
			IP: addr["ip"].(string),
		}
	}

	opts := resolverrule.UpdateOpts{
		Name:        name,
		IPAddresses: ipAddressList,
	}
	_, err = resolverrule.Update(dnsClient, d.Id(), opts).Extract()
	if err != nil {
		return diag.Errorf("error updating DNS resolver rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      waitForDNSResolverRule(dnsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS resolver rule (%s) to become ACTIVE for update: %s",
			d.Id(), err)
	}

	return resourceDNSResolverRuleRead(ctx, d, meta)
}

func resourceDNSResolverRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	err = resolverrule.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS resolver rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) to become DELETED", d.Id())
	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		// we allow to try to delete ERROR resolver rule
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      waitForDNSResolverRule(dnsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS resolver rule (%s) to delete: %s",
			d.Id(), err)
	}
	return nil
}

func waitForDNSResolverRule(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rule, err := resolverrule.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return rule, "DELETED", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] DNS resolver rule (%s) current status: %s", rule.ID, rule.Status)
		return rule, parseStatus(rule.Status), nil
	}
}
