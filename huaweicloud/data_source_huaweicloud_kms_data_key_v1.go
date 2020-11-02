package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"
)

func dataSourceKmsDataKeyV1() *schema.Resource {
	return &schema.Resource{
		Read: resourceKmsDataKeyV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encryption_context": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datakey_length": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plain_text": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cipher_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceKmsDataKeyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	KmsDataKeyV1Client, err := config.kmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud kms key client: %s", err)
	}

	req := &keys.DataEncryptOpts{
		KeyID:             d.Get("key_id").(string),
		EncryptionContext: d.Get("encryption_context").(string),
		DatakeyLength:     d.Get("datakey_length").(string),
	}
	log.Printf("[DEBUG] KMS get data key for key: %s", d.Get("key_id").(string))
	v, err := keys.DataEncryptGet(KmsDataKeyV1Client, req).ExtractDataKey()
	if err != nil {
		return err
	}

	//lintignore:R017
	d.SetId(time.Now().UTC().String())
	d.Set("plain_text", v.PlainText)
	d.Set("cipher_text", v.CipherText)

	return nil
}
