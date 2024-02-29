package iec

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC POST /v1/vpcs
// @API IEC DELETE /v1/vpcs/{vpc_id}
// @API IEC GET /v1/vpcs/{vpc_id}
// @API IEC PUT /v1/vpcs/{vpc_id}
func ResourceVpc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcCreate,
		ReadContext:   resourceVpcRead,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SYSTEM",
			},
			"subnet_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVpcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name: d.Get("name").(string),
		Cidr: d.Get("cidr").(string),
		Mode: d.Get("mode").(string),
	}

	n, err := vpcs.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC VPC: %s", err)
	}

	log.Printf("[INFO] IEC VPC ID: %s", n.ID)
	d.SetId(n.ID)

	return resourceVpcRead(ctx, d, meta)
}

func resourceVpcRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	n, err := vpcs.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IEC VPC")
	}
	mErr := multierror.Append(
		nil,
		d.Set("name", n.Name),
		d.Set("cidr", n.Cidr),
		d.Set("mode", n.Mode),
		d.Set("subnet_num", n.SubnetNum),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVpcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	var updateOpts vpcs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("cidr") {
		updateOpts.Cidr = d.Get("cidr").(string)
	}

	_, err = vpcs.Update(iecClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating IEC VPC: %s", err)
	}

	return resourceVpcRead(ctx, d, meta)
}

func resourceVpcDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	// lintignore:R018
	time.Sleep(3 * time.Second) // Prevent delete failure

	err = vpcs.Delete(iecClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IEC VPC")
	}

	return nil
}
