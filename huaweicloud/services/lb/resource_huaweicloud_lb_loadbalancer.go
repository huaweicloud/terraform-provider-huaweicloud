package lb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"
	lbv3 "github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v2/{project_id}/elb/loadbalancers
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB POST /v2.0/{project_id}/loadbalancers/{loadbalancer_id}/tags/action
// @API ELB GET /v2.0/{project_id}/loadbalancers/{loadbalancer_id}/tags
// @API ELB PUT /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB DELETE /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB POST /v3/{project_id}/elb/loadbalancers/change-charge-mode
// @API VPC PUT /v2.0/ports/{port_id}
// @API VPC GET /v2.0/ports/{port_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV2Create,
		ReadContext:   resourceLoadBalancerV2Read,
		UpdateContext: resourceLoadBalancerV2Update,
		DeleteContext: resourceLoadBalancerV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vip_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the IPv4 subnet ID of the subnet where the load balancer works",
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": common.TagsSchema(),

			"loadbalancer_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "schema: Deprecated",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "schema: Deprecated",
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// charge info: charging_mode, period_unit, period, auto_renew
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
			},

			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},

			"auto_renew": common.SchemaAutoRenewUpdatable(nil),

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Deprecated",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "schema: Deprecated",
			},
		},
	}
}

func resourceLoadBalancerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.LoadBalancerClient(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2 Client: %s", err)
	}

	// client for changing charging mode
	elbV3Client, err := cfg.ElbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	// client for setting tags
	elbV2Client, err := cfg.ElbV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2.0 Client: %s", err)
	}

	var lbProvider string
	if v, ok := d.GetOk("loadbalancer_provider"); ok {
		lbProvider = v.(string)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := loadbalancers.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		VipSubnetID:         d.Get("vip_subnet_id").(string),
		TenantID:            d.Get("tenant_id").(string),
		VipAddress:          d.Get("vip_address").(string),
		AdminStateUp:        &adminStateUp,
		Flavor:              d.Get("flavor").(string),
		Provider:            lbProvider,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		ProtectionStatus:    d.Get("protection_status").(string),
		ProtectionReason:    d.Get("protection_reason").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	lb, err := loadbalancers.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating LoadBalancer: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(ctx, elbClient, lb.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// set the ID on the resource
	d.SetId(lb.ID)

	// change to pre-paid mode
	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}
		changeChargingModeOpts := lbv3.ChangeChargingModeOpts{
			LoadBalancerIds: []string{d.Id()},
			ChargingMode:    "prepaid",
			PrepaidOptions: lbv3.PrepaidOptions{
				PeriodType: d.Get("period_unit").(string),
				PeriodNum:  d.Get("period").(int),
				AutoRenew:  d.Get("auto_renew").(string),
				AutoPay:    true,
			},
		}
		orderId, err := lbv3.ChangeChargingMode(elbV3Client, changeChargingModeOpts).Extract()
		if err != nil {
			return diag.Errorf("error changing charging mode of load-balancer(%s): %s", d.Id(), err)
		}

		// wait for order complete
		if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	// Once the LoadBalancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	if lb.VipPortID != "" {
		networkingClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2 Client: %s", err)
		}

		if err := resourceLoadBalancerV2SecurityGroups(networkingClient, lb.VipPortID, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(elbV2Client, "loadbalancers", lb.ID, tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of LoadBalancer %s: %s", lb.ID, tagErr)
		}
	}

	return resourceLoadBalancerV2Read(ctx, d, meta)
}

func resourceLoadBalancerV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.LoadBalancerClient(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2 Client: %s", err)
	}

	// client for fetching tags
	elbV2Client, err := cfg.ElbV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2.0 Client: %s", err)
	}

	lb, err := loadbalancers.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving LoadBalancer")
	}

	log.Printf("[DEBUG] Retrieved LoadBalancer %s: %#v", d.Id(), lb)

	var publicIp string
	if len(lb.PublicIps) > 0 {
		publicIp = lb.PublicIps[0].PublicIpAddress
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("vip_subnet_id", lb.VipSubnetID),
		d.Set("tenant_id", lb.TenantID),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_port_id", lb.VipPortID),
		d.Set("admin_state_up", lb.AdminStateUp),
		d.Set("flavor", lb.Flavor),
		d.Set("loadbalancer_provider", lb.Provider),
		d.Set("enterprise_project_id", lb.EnterpriseProjectID),
		d.Set("protection_status", lb.ProtectionStatus),
		d.Set("protection_reason", lb.ProtectionReason),
		d.Set("public_ip", publicIp),
		d.Set("charge_mode", lb.ChargeMode),
		d.Set("frozen_scene", lb.FrozenScene),
		d.Set("created_at", lb.CreatedAt),
		d.Set("updated_at", lb.UpdatedAt),
	)

	if lb.BillingInfo != "" {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	} else {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "postPaid"))
	}

	// Get any security groups on the VIP Port
	if lb.VipPortID != "" {
		networkingClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2 Client: %s", err)
		}

		port, err := ports.Get(networkingClient, lb.VipPortID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr, d.Set("security_group_ids", port.SecurityGroups))
	}

	// fetch tags
	if resourceTags, err := tags.Get(elbV2Client, "loadbalancers", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] fetching tags of ELB LoadBalancer failed: %s", err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting LoadBalancer fields: %s", err)
	}

	return nil
}

func resourceLoadBalancerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.LoadBalancerClient(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2 Client: %s", err)
	}
	elbV3Client, err := cfg.ElbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChanges("name", "description", "admin_state_up", "protection_status", "protection_reason") {
		var updateOpts loadbalancers.UpdateOpts
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			desc := d.Get("description").(string)
			updateOpts.Description = &desc
		}
		if d.HasChange("admin_state_up") {
			asu := d.Get("admin_state_up").(bool)
			updateOpts.AdminStateUp = &asu
		}
		if d.HasChange("protection_status") {
			updateOpts.ProtectionStatus = d.Get("protection_status").(string)
		}
		if d.HasChange("protection_reason") {
			protectionReason := d.Get("protection_reason").(string)
			updateOpts.ProtectionReason = &protectionReason
		}

		// Wait for LoadBalancer to become active before continuing
		err = waitForLBV2LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating LoadBalancer %s with options: %#v", d.Id(), updateOpts)
		// lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err = loadbalancers.Update(elbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("error updating loadbalancer: %s", err)
		}

		// Wait for LoadBalancer to become active before continuing
		err = waitForLBV2LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Security Groups get updated separately
	if d.HasChange("security_group_ids") {
		vipPortID := d.Get("vip_port_id").(string)
		if vipPortID != "" {
			networkingClient, err := cfg.NetworkingV2Client(region)
			if err != nil {
				return diag.Errorf("error creating VPC V2 Client: %s", err)
			}

			if err := resourceLoadBalancerV2SecurityGroups(networkingClient, vipPortID, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// update charging mode
	if d.HasChange("charging_mode") {
		if d.Get("charging_mode").(string) == "postPaid" {
			return diag.Errorf("error updating the charging mode of the load-balancer (%s): %s", d.Id(),
				"only support changing post-paid load-balancer to pre-paid")
		}
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}
		changeChargingModeOpts := lbv3.ChangeChargingModeOpts{
			LoadBalancerIds: []string{d.Id()},
			ChargingMode:    "prepaid",
			PrepaidOptions: lbv3.PrepaidOptions{
				PeriodType: d.Get("period_unit").(string),
				PeriodNum:  d.Get("period").(int),
				AutoRenew:  d.Get("auto_renew").(string),
				AutoPay:    true,
			},
		}
		orderId, err := lbv3.ChangeChargingMode(elbV3Client, changeChargingModeOpts).Extract()
		if err != nil {
			return diag.Errorf("error changing charging mode of load-balancer(%s): %s", d.Id(), err)
		}

		// wait for order complete
		if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeout); err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the load-balancer (%s): %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		elbV2Client, err := cfg.ElbV2Client(region)
		if err != nil {
			return diag.Errorf("error creating ELB v2.0 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(elbV2Client, d, "loadbalancers", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of LoadBalancer:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "loadbalancers",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLoadBalancerV2Read(ctx, d, meta)
}

func resourceLoadBalancerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.LoadBalancerClient(region)
	if err != nil {
		return diag.Errorf("error creating ELB v2 Client: %s", err)
	}

	log.Printf("[DEBUG] Deleting LoadBalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	if d.Get("charging_mode").(string) == "prePaid" {
		// Unsubscribe the prepaid LoadBalancer will automatically delete it
		if err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing ELB LoadBalancer : %s", err)
		}
	} else {
		// lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			err = loadbalancers.Delete(elbClient, d.Id()).ExtractErr()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("error deleting loadbalancer: %s", err)
		}
	}

	// Wait for LoadBalancer to become delete
	pending := []string{"PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}
	err = waitForLBV2LoadBalancer(ctx, elbClient, d.Id(), "DELETED", pending, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceLoadBalancerV2SecurityGroups(networkingClient *golangsdk.ServiceClient, vipPortID string, d *schema.ResourceData) error {
	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroups := resourcePortSecurityGroupsV2(v.(*schema.Set))
		updateOpts := ports.UpdateOpts{
			SecurityGroups: &securityGroups,
		}

		log.Printf("[DEBUG] Adding security groups to LoadBalancer "+
			"VIP Port %s: %#v", vipPortID, updateOpts)

		_, err := ports.Update(networkingClient, vipPortID, updateOpts).Extract()
		if err != nil {
			return err
		}
	}

	return nil
}

func resourcePortSecurityGroupsV2(v *schema.Set) []string {
	var securityGroups []string
	for _, v := range v.List() {
		securityGroups = append(securityGroups, v.(string))
	}
	return securityGroups
}
