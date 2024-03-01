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
// @API VPC PUT /v1/ports/{port_id}
// @API VPC GET /v1/ports/{port_id}
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

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
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
				ForceNew: true,
				Computed: true,
			},

			"public_ip": {
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
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, cfg),
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

	// Once the LoadBalancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	if lb.VipPortID != "" {
		networkingClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v1 Client: %s", err)
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
		d.Set("public_ip", publicIp),
	)

	// Get any security groups on the VIP Port
	if lb.VipPortID != "" {
		networkingClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v1 Client: %s", err)
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

	if d.HasChanges("name", "description", "admin_state_up") {
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

		// Wait for LoadBalancer to become active before continuing
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBV2LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating LoadBalancer %s with options: %#v", d.Id(), updateOpts)
		//lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err = loadbalancers.Update(elbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})

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
			networkingClient, err := cfg.NetworkingV1Client(region)
			if err != nil {
				return diag.Errorf("error creating VPC V1 Client: %s", err)
			}

			if err := resourceLoadBalancerV2SecurityGroups(networkingClient, vipPortID, d); err != nil {
				return diag.FromErr(err)
			}
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
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = loadbalancers.Delete(elbClient, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

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
