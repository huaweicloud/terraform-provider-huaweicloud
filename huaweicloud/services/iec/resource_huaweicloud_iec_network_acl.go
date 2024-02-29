package iec

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC POST /v1/firewalls
// @API IEC PUT /v1/firewalls/{firewall_id}
// @API IEC DELETE /v1/firewalls/{firewall_id}
// @API IEC GET /v1/firewalls/{firewall_id}
func ResourceNetworkACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkACLCreate,
		ReadContext:   resourceNetworkACLRead,
		UpdateContext: resourceNetworkACLUpdate,
		DeleteContext: resourceNetworkACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"networks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"inbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"outbound_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkACLNetworks(d *schema.ResourceData) []firewalls.ReqSubnet {
	rawNetworks := d.Get("networks").(*schema.Set).List()
	networkOpts := make([]firewalls.ReqSubnet, len(rawNetworks))

	for i, val := range rawNetworks {
		raw := val.(map[string]interface{})
		networkOpts[i] = firewalls.ReqSubnet{
			ID:    raw["subnet_id"].(string),
			VpcID: raw["vpc_id"].(string),
		}
	}
	return networkOpts
}

func resourceNetworkACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	createOpts := firewalls.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] create IEC network ACL: %#v", createOpts)
	group, err := firewalls.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC network ACL: %s", err)
	}
	d.SetId(group.ID)

	// Associate subnets
	subnetsOpts := resourceNetworkACLNetworks(d)
	if len(subnetsOpts) > 0 {
		updateOpts := firewalls.UpdateOpts{
			Name:    d.Get("name").(string),
			Subnets: &subnetsOpts,
		}
		log.Printf("[DEBUG] attempt to associate IEC network ACL with subnets: %#v", updateOpts)
		_, err := firewalls.Update(iecClient, group.ID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error associating subnets with IEC network ACL: %s", err)
		}
	}

	return resourceNetworkACLRead(ctx, d, meta)
}

func resourceNetworkACLRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	var mErr *multierror.Error
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	fwGroup, err := firewalls.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IEC network ACL")
	}

	log.Printf("[DEBUG] read IEC network ACL %s: %#v", d.Id(), fwGroup)
	mErr = multierror.Append(
		mErr,
		d.Set("name", fwGroup.Name),
		d.Set("status", fwGroup.Status),
		d.Set("description", fwGroup.Description),
		d.Set("inbound_rules", getFirewallRuleIDs(fwGroup.IngressFWPolicy)),
		d.Set("outbound_rules", getFirewallRuleIDs(fwGroup.EgressFWPolicy)),
	)
	var networkSet []map[string]interface{}
	for _, val := range fwGroup.Subnets {
		subnet := make(map[string]interface{})
		subnet["vpc_id"] = val.VpcID
		subnet["subnet_id"] = val.ID
		networkSet = append(networkSet, subnet)
	}
	if err = d.Set("networks", networkSet); err != nil {
		return diag.Errorf("saving IEC networks failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	if d.HasChanges("name", "description", "networks") {
		opts := firewalls.UpdateOpts{
			Name: d.Get("name").(string),
		}

		if d.HasChange("description") {
			desc := d.Get("description").(string)
			opts.Description = &desc
		}

		if d.HasChange("networks") {
			subnetsOpts := resourceNetworkACLNetworks(d)
			opts.Subnets = &subnetsOpts
		}

		log.Printf("[DEBUG] updating IEC network ACL with id %s: %#v", d.Id(), opts)
		_, err := firewalls.Update(iecClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.Errorf("error updating IEC network ACL: %s", err)
		}
	}

	return resourceNetworkACLRead(ctx, d, meta)
}

func resourceNetworkACLDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	// Unbind subents before deleting the network acl.
	rawNetworks := d.Get("networks").(*schema.Set).List()
	if len(rawNetworks) > 0 {
		unbindSubnets := make([]firewalls.ReqSubnet, 0)
		opts := firewalls.UpdateOpts{
			Name:    d.Get("name").(string),
			Subnets: &unbindSubnets,
		}
		err = firewalls.Update(iecClient, d.Id(), opts).Err
		if err != nil {
			return diag.Errorf("error disassociating all subents with IEC network ACL: %s", err)
		}
	}

	err = firewalls.Delete(iecClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting IEC firewall: %s", err)
	}

	return nil
}
