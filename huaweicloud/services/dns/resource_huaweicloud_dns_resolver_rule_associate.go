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
	"github.com/chnsz/golangsdk/openstack/dns/v2/associate"
	"github.com/chnsz/golangsdk/openstack/dns/v2/resolverrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DNS POST /v2.1/resolverrules/{resolverrule_id}/associaterouter
// @API DNS POST /v2.1/resolverrules/{resolverrule_id}/disassociaterouter
// @API DNS GET /v2.1/resolverrules/{resolverrule_id}
func ResourceDNSResolverRuleAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSResolverRuleAssociateCreate,
		ReadContext:   resourceDNSResolverRuleAssociateRead,
		DeleteContext: resourceDNSResolverRuleAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resolver_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSResolverRuleAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	ruleID := d.Get("resolver_rule_id").(string)
	vpcID := d.Get("vpc_id").(string)
	opts := associate.RouterOpts{
		RouterID: vpcID,
	}

	body, err := associate.Associate(dnsClient, ruleID, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS resolver rule associate: %s", err)
	}

	id := fmt.Sprintf("%s/%s", ruleID, vpcID)
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for DNS resolver rule associate (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      waitForDNSResolverRuleAssociate(dnsClient, ruleID, body.RouterID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS resolver rule associate (%s) to become ACTIVE for creation: %s",
			d.Id(), err)
	}

	return resourceDNSResolverRuleAssociateRead(ctx, d, meta)
}

func resourceDNSResolverRuleAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	arr := strings.Split(d.Id(), "/")
	if len(arr) != 2 {
		return diag.Errorf("error getting resolver rule ID and VPC ID, DNS resolver rule associate ID: %s", d.Id())
	}
	ruleID := arr[0]
	vpcID := arr[1]

	rule, err := resolverrule.Get(dnsClient, ruleID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS resolver rule")
	}

	for _, r := range rule.Routers {
		if r.RouterID == vpcID {
			mErr := multierror.Append(nil,
				d.Set("vpc_id", vpcID),
				d.Set("resolver_rule_id", ruleID),
				d.Set("status", r.Status),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
}

func resourceDNSResolverRuleAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	arr := strings.Split(d.Id(), "/")
	if len(arr) != 2 {
		return diag.Errorf("error getting resolver rule ID and VPC ID, DNS resolver rule associate resource ID: %s", d.Id())
	}
	ruleID := arr[0]
	vpcID := arr[1]

	opts := associate.RouterOpts{
		RouterID: vpcID,
	}
	body, err := associate.DisAssociate(dnsClient, ruleID, opts).Extract()

	if err != nil {
		return diag.Errorf("error deleting DNS resolver rule associate: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS resolver rule associate (%s) to become DELETED", d.Id())

	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		// we allow to try to delete ERROR resolver rule associate
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      waitForDNSResolverRuleAssociate(dnsClient, ruleID, body.RouterID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for DNS resolver rule associate (%s) to delete: %s",
			d.Id(), err)
	}

	return nil
}

func waitForDNSResolverRuleAssociate(client *golangsdk.ServiceClient, resolverRuleID string, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rule, err := resolverrule.Get(client, resolverRuleID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil, "DELETED", nil
			}
			return nil, "", err
		}

		for _, router := range rule.Routers {
			if router.RouterID == vpcID {
				log.Printf("[DEBUG] DNS resolver rule associate (%s) current status: %s", resolverRuleID, router.Status)
				return router, parseStatus(router.Status), nil
			}
		}

		return rule.Routers, "DELETED", nil
	}
}
