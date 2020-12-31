package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/keypairs"
)

func resourceIecKeypair() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecKeypairCreate,
		Read:   resourceIecKeypairRead,
		Delete: resourceIecKeypairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceIecKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	createOpts := keypairs.CreateOpts{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("public_key").(string),
	}

	log.Printf("[DEBUG] Create iec keypair Options: %#v", createOpts)
	kp, err := keypairs.Create(iecClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud iec keypair: %s", err)
	}

	d.SetId(kp.Name)

	return resourceIecKeypairRead(d, meta)
}

func resourceIecKeypairRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	kp, err := keypairs.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "keypair")
	}

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("fingerprint", kp.Fingerprint)

	return nil
}

func resourceIecKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	err = keypairs.Delete(iecClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud iec keypair: %s", err)
	}
	d.SetId("")
	return nil
}
