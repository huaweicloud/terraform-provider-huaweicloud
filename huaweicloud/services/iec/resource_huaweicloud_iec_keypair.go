package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/iec/v1/keypairs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC DELETE /v1/os-keypairs/{keypair_name}
// @API IEC GET /v1/os-keypairs/{keypair_name}
// @API IEC POST /v1/os-keypairs
func ResourceKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeypairCreate,
		ReadContext:   resourceKeypairRead,
		DeleteContext: resourceKeypairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKeypairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	createOpts := keypairs.CreateOpts{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
	}

	log.Printf("[DEBUG] Create IEC keypair options: %#v", createOpts)
	kp, err := keypairs.Create(iecClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC keypair: %s", err)
	}

	d.SetId(kp.Name)

	return resourceKeypairRead(ctx, d, meta)
}

func resourceKeypairRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	kp, err := keypairs.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IEC keypair")
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", kp.Name),
		d.Set("public_key", kp.PublicKey),
		d.Set("fingerprint", kp.Fingerprint),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IEC keypair: %s", err)
	}
	return nil
}

func resourceKeypairDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	err = keypairs.Delete(iecClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting IEC keypair: %s", err)
	}
	return nil
}
