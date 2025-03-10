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
// @API DNS GET /v2.1/resolverrules/{resolverrule_id}
// @API DNS POST /v2.1/resolverrules/{resolverrule_id}/disassociaterouter
func ResourceResolverRuleAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResolverRuleAssociateCreate,
		ReadContext:   resourceResolverRuleAssociateRead,
		DeleteContext: resourceResolverRuleAssociateDelete,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The DNS resolver rule ID.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID to associate.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resolver rule association.`,
			},
		},
	}
}

func resourceResolverRuleAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dnsClient, err := cfg.DNSV21Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	resolverRuleId := d.Get("resolver_rule_id").(string)
	vpcId := d.Get("vpc_id").(string)
	opts := associate.RouterOpts{
		RouterID: vpcId,
	}

	_, err = associate.Associate(dnsClient, resolverRuleId, opts).Extract()
	if err != nil {
		return diag.Errorf("error associating VPC (%s) to DNS resolver rule (%s): %s", vpcId, resolverRuleId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", resolverRuleId, vpcId))

	log.Printf("[DEBUG] Waiting for DNS resolver rule (%s) associate VPC (%s) to complete", resolverRuleId, vpcId)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForResolverRuleAssociate(dnsClient, resolverRuleId, vpcId, false),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DNS resolver rule (%s) associate VPC (%s) to complete: %s",
			resolverRuleId, vpcId, err)
	}

	return resourceResolverRuleAssociateRead(ctx, d, meta)
}

func getResolverRuleAndVpcId(resourceId string) (resolverRuleId string, vpcId string, err error) {
	parts := strings.Split(resourceId, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource ID format (%s), want '<resolver_rule_id>/<vpc_id>'", resourceId)
	}
	resolverRuleId = parts[0]
	vpcId = parts[1]
	return
}

func resourceResolverRuleAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	resolverRuleId, vpcId, err := getResolverRuleAndVpcId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	associatedVpc, err := GetAssociatedVpcById(dnsClient, resolverRuleId, vpcId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving associated VPC from resolver rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", associatedVpc.RouterID),
		d.Set("resolver_rule_id", resolverRuleId),
		d.Set("status", associatedVpc.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAssociatedVpcById(client *golangsdk.ServiceClient, resolverRuleId, vpcId string) (*resolverrule.Router, error) {
	rule, err := resolverrule.Get(client, resolverRuleId).Extract()
	if err != nil {
		return nil, err
	}

	for _, router := range rule.Routers {
		if router.RouterID == vpcId {
			return &router, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceResolverRuleAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DNSV21Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	resolverRuleId := d.Get("resolver_rule_id").(string)
	vpcId := d.Get("vpc_id").(string)
	opts := associate.RouterOpts{
		RouterID: vpcId,
	}
	_, err = associate.DisAssociate(dnsClient, resolverRuleId, opts).Extract()
	if err != nil {
		return diag.Errorf("error disassociating VPC (%s) from DNS resolver rule (%s): %s", vpcId, resolverRuleId, err)
	}
	log.Printf("[DEBUG] Waiting for disassociating VPC (%s) from DNS resolver rule (%s) to complete", vpcId, resolverRuleId)

	stateConf := &resource.StateChangeConf{
		// we allow to try to delete ERROR resolver rule associate
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      waitForResolverRuleAssociate(dnsClient, resolverRuleId, vpcId, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for disassociating VPC (%s) from DNS resolver rule (%s) to complete: %s",
			vpcId, resolverRuleId, err)
	}

	return nil
}

func waitForResolverRuleAssociate(client *golangsdk.ServiceClient, resolverRuleId string, vpcId string, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		associatedVpc, err := GetAssociatedVpcById(client, resolverRuleId, vpcId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := associatedVpc.Status
		if !isDelete {
			if status == "ACTIVE" {
				return associatedVpc, "COMPLETED", nil
			}

			if status == "ERROR" {
				return associatedVpc, "ERROR", fmt.Errorf("unexpect status (%s)", status)
			}
		}
		return associatedVpc, "PENDING", nil
	}
}
