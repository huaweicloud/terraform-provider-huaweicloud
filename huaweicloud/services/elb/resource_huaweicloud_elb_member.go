package elb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API ELB POST /v3/{project_id}/elb/pools/{pool_id}/members
// @API ELB GET /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
// @API ELB PUT /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
// @API ELB DELETE /v3/{project_id}/elb/pools/{pool_id}/members/{member_id}
func ResourceMemberV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV3Create,
		ReadContext:   resourceMemberV3Read,
		UpdateContext: resourceMemberV3Update,
		DeleteContext: resourceMemberV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceELBMemberImport,
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
				Optional: true,
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 {
						errors = append(errors, fmt.Errorf(
							"only numbers greater than 0 are supported values for 'weight'"))
					}
					return
				},
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IPv4 or IPv6 subnet ID of the subnet in which to access the member",
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMemberV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createOpts := pools.CreateMemberOpts{
		Name:         d.Get("name").(string),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		Weight:       d.Get("weight").(int),
	}

	// Must omit if not set
	if v, ok := d.GetOk("subnet_id"); ok {
		createOpts.SubnetID = v.(string)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	poolID := d.Get("pool_id").(string)
	member, err := pools.CreateMember(elbClient, poolID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating member: %s", err)
	}

	d.SetId(member.ID)

	return resourceMemberV3Read(ctx, d, meta)
}

func resourceMemberV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	member, err := pools.GetMember(elbClient, d.Get("pool_id").(string), d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "member")
	}

	log.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	mErr := multierror.Append(nil,
		d.Set("name", member.Name),
		d.Set("weight", member.Weight),
		d.Set("subnet_id", member.SubnetID),
		d.Set("address", member.Address),
		d.Set("protocol_port", member.ProtocolPort),
		d.Set("region", cfg.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB member fields: %s", err)
	}

	return nil
}

func resourceMemberV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var updateOpts pools.UpdateMemberOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("weight") {
		updateOpts.Weight = d.Get("weight").(int)
	}

	log.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)
	poolID := d.Get("pool_id").(string)
	_, err = pools.UpdateMember(elbClient, poolID, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update member %s: %s", d.Id(), err)
	}

	return resourceMemberV3Read(ctx, d, meta)
}

func resourceMemberV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	poolID := d.Get("pool_id").(string)
	err = pools.DeleteMember(elbClient, poolID, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete member %s: %s", d.Id(), err)
	}
	return nil
}

func resourceELBMemberImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for member. Format must be <pool_id>/<member_id>")
		return nil, err
	}

	poolID := parts[0]
	memberID := parts[1]

	d.SetId(memberID)
	d.Set("pool_id", poolID)

	return []*schema.ResourceData{d}, nil
}
