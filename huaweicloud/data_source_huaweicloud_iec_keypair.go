package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/keypairs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceIECKeypair() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIECKeypairRead,

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

func dataSourceIECKeypairRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	keyName := d.Get("name").(string)
	kp, err := keypairs.Get(iecClient, keyName).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return fmtp.Errorf("Your query returned no results. " +
				"Please change your search criteria and try again.")
		}
		return fmtp.Errorf("fetching IEC keypair error: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved IEC keypair %s: %+v", kp.Name, kp)
	d.SetId(kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("fingerprint", kp.Fingerprint)

	return nil
}
