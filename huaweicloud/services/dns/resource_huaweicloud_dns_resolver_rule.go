package dns

import (
	"context"
	"fmt"
	"log"
	"strings"
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
func ResourceResolverRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResolverRuleCreate,
		ReadContext:   resourceResolverRuleRead,
		UpdateContext: resourceResolverRuleUpdate,
		DeleteContext: resourceResolverRuleDelete,
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
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the DNS endpoint to which the resolver rule belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resolver rule name.`,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
					return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
				},
				Description: `The domain name.`,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The IP of the IP address.`,
						},
					},
				},
				Description: `The IP address list of the DNS resolver rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resolver rule.`,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rule type of the resolver rule.`,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID.`,
						},
						"vpc_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region of the VPC.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the VPC.`,
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the resolver rule.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the resolver rule.`,
			},
		},
	}
}

func buildIpAddresses(ipAddresses *schema.Set) []resolverrule.IPAddress {
	rest := make([]resolverrule.IPAddress, 0)
	for _, v := range ipAddresses.List() {
		ipAddress := v.(map[string]interface{})
		ip := ipAddress["ip"].(string)
		if ip == "" {
			continue
		}

		rest = append(rest, resolverrule.IPAddress{
			IP: ip,
		})
	}

	return rest
}

func resourceResolverRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	opts := resolverrule.CreateOpts{
		Name:        d.Get("name").(string),
		DomainName:  d.Get("domain_name").(string),
		EndpointID:  d.Get("endpoint_id").(string),
		IPAddresses: buildIpAddresses(d.Get("ip_addresses").(*schema.Set)),
	}
	rule, err := resolverrule.Create(dnsClient, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS resolver rule: %s", err)
	}

	resolverRuleId := rule.ID
	if resolverRuleId == "" {
		return diag.Errorf("unable to find resolver rule ID from API response")
	}

	d.SetId(resolverRuleId)

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) status to become active", resolverRuleId)
	err = waitForResolverRuleStatusCompleted(ctx, dnsClient, resolverRuleId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for DNS resolver rule (%s) to create: %s", resolverRuleId, err)
	}

	return resourceResolverRuleRead(ctx, d, meta)
}

func waitForResolverRuleStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, resolverRuleId string,
	timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshResolverRuleStatus(client, resolverRuleId, false),
		Timeout:      timeOut,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func GetResolverRuleById(client *golangsdk.ServiceClient, resolverRuleId string) (*resolverrule.ResolverRule, error) {
	respBody, err := resolverrule.Get(client, resolverRuleId).Extract()
	if err != nil {
		return nil, err
	}

	resolverRule := respBody.ResolverRule
	// When the resolver rule has been deleted, the status code of calling the query interface is not 404 but 200 and
	// the status value is `DELETED`.
	if resolverRule.Status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return &resolverRule, nil
}

func resourceResolverRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	resolverRuleId := d.Id()
	rule, err := GetResolverRuleById(dnsClient, resolverRuleId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS resolver rule")
	}

	mErr := multierror.Append(nil,
		d.Set("name", rule.Name),
		d.Set("domain_name", rule.DomainName),
		d.Set("endpoint_id", rule.EndpointID),
		d.Set("ip_addresses", flattenIpAddresses(rule.IPAddresses)),
		// Attributes.
		d.Set("status", rule.Status),
		d.Set("rule_type", rule.RuleType),
		d.Set("vpcs", flattenResolverRuleVpcs(rule.Routers)),
		d.Set("created_at", rule.CreatedAt),
		d.Set("updated_at", rule.UpdatedAt),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenResolverRuleVpcs(routers []resolverrule.Router) []map[string]interface{} {
	vpcs := make([]map[string]interface{}, len(routers))
	for i, r := range routers {
		vpcs[i] = map[string]interface{}{
			"vpc_id":     r.RouterID,
			"vpc_region": r.RouterRegion,
			"status":     r.Status,
		}
	}
	return vpcs
}

func flattenIpAddresses(ipAddresses []resolverrule.IPAddress) []map[string]interface{} {
	rest := make([]map[string]interface{}, len(ipAddresses))
	for i, r := range ipAddresses {
		rest[i] = map[string]interface{}{
			"ip": r.IP,
		}
	}
	return rest
}

func resourceResolverRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	opts := resolverrule.UpdateOpts{
		Name:        d.Get("name").(string),
		IPAddresses: buildIpAddresses(d.Get("ip_addresses").(*schema.Set)),
	}
	resolverRuleId := d.Id()
	_, err = resolverrule.Update(dnsClient, resolverRuleId, opts).Extract()
	if err != nil {
		return diag.Errorf("error updating DNS resolver rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) status to become active", resolverRuleId)
	err = waitForResolverRuleStatusCompleted(ctx, dnsClient, resolverRuleId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error waiting for DNS resolver rule (%s) to update: %s", resolverRuleId, err)
	}

	return resourceResolverRuleRead(ctx, d, meta)
}

func resourceResolverRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	resolverRuleId := d.Id()
	err = resolverrule.Delete(dnsClient, resolverRuleId).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS resolver rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) to become DELETED", resolverRuleId)
	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		// Allows deletion of  resolver rule with status `ACTIVE` and `ERROR`.
		Pending:      []string{"PENDING"},
		Refresh:      refreshResolverRuleStatus(dnsClient, resolverRuleId, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DNS resolver rule (%s) to be deleted: %s", resolverRuleId, err)
	}
	return nil
}

func refreshResolverRuleStatus(client *golangsdk.ServiceClient, resolverRuleId string, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resolverRule, err := GetResolverRuleById(client, resolverRuleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "", err
		}

		status := resolverRule.Status
		if !isDelete {
			if status == "ACTIVE" {
				return resolverRule, "COMPLETED", nil
			}

			if status == "ERROR" {
				return resolverRule, "ERROR", fmt.Errorf("unexpect status (%s)", status)
			}
		}

		return resolverRule, "PENDING", nil
	}
}
