package huaweicloud

import (
	"fmt"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/keypairs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeKeypairV2Create,
		Read:   resourceComputeKeypairV2Read,
		Delete: resourceComputeKeypairV2Delete,
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"key_file"},
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeKeypairV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	pk, isExist := d.GetOk("public_key")
	createOpts := keypairs.CreateOpts{
		Name:      d.Get("name").(string),
		PublicKey: pk.(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	kp, err := keypairs.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud keypair: %s", err)
	}

	d.SetId(kp.Name)

	if !isExist {
		fp := getKeyFilePath(d)
		if err = utils.WriteToPemFile(fp, kp.PrivateKey); err != nil {
			return fmtp.Errorf("Unable to generate private key: %s", err)
		}
		d.Set("key_file", fp)
	}

	return resourceComputeKeypairV2Read(d, meta)
}

func getKeyFilePath(d *schema.ResourceData) string {
	if path, ok := d.GetOk("key_file"); ok {
		return path.(string)
	}
	keypairName := d.Get("name").(string)
	return fmt.Sprintf("%s.pem", keypairName)
}

func resourceComputeKeypairV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	kp, err := keypairs.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "keypair")
	}

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeKeypairV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	err = keypairs.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud keypair: %s", err)
	}
	d.SetId("")
	return nil
}
