package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"
	"time"
)

func dataSourceKmsDataKeyV1() *schema.Resource {
	return &schema.Resource{
		Read: resourceKmsDataKeyV1Read,

		Schema: map[string]*schema.Schema{
			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"encryption_context": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"datakey_length": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plain_text": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cipher_text": &schema.Schema{
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

	d.SetId(time.Now().UTC().String())
	d.Set("plain_text", v.PlainText)
	d.Set("cipher_text", v.CipherText)

	return nil
}
