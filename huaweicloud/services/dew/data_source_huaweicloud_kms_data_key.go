package dew

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DEW POST /v1.0/{project_id}/kms/create-datakey
func DataSourceKmsDataKeyV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceKmsDataKeyV1Read,

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
				Computed: true,
			},
			"cipher_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKmsDataKeyV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	kmsDataKeyV1Client, err := cfg.KmsKeyV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating kms key client: %s", err)
	}

	req := &keys.DataEncryptOpts{
		KeyID:             d.Get("key_id").(string),
		EncryptionContext: d.Get("encryption_context").(string),
		DatakeyLength:     d.Get("datakey_length").(string),
	}
	v, err := keys.DataEncryptGet(kmsDataKeyV1Client, req).ExtractDataKey()
	if err != nil {
		return diag.FromErr(err)
	}

	// lintignore:R017
	d.SetId(time.Now().UTC().String())
	d.Set("plain_text", v.PlainText)
	d.Set("cipher_text", v.CipherText)

	return nil
}
