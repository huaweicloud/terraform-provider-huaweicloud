package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/keypairs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/os-keypairs/{keypair_name}
func DataSourceKeypair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeypairRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeypairRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	keyName := d.Get("name").(string)
	kp, err := keypairs.Get(iecClient, keyName).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return diag.Errorf("your query returned no results, " +
				"please change your search criteria and try again")
		}
		return diag.Errorf("fetching IEC keypair error: %s", err)
	}

	log.Printf("[DEBUG] Retrieved IEC keypair %s: %+v", kp.Name, kp)
	d.SetId(kp.Name)
	mErr := multierror.Append(
		d.Set("public_key", kp.PublicKey),
		d.Set("fingerprint", kp.Fingerprint),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
