package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/encrypt-data
// @API DEW POST /v1.0/{project_id}/kms/decrypt-data
func ResourceKmsDataEncryptDecrypt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsDataEncryptDecryptCreate,
		ReadContext:   resourceKmsDataEncryptDecryptRead,
		DeleteContext: resourceKmsDataEncryptDecryptDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"plain_text": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cipher_text": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cipher_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plain_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plain_text_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDataEncryptDecryptBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":               utils.ValueIgnoreEmpty(d.Get("key_id")),
		"encryption_algorithm": utils.ValueIgnoreEmpty(d.Get("encryption_algorithm")),
		"plain_text":           utils.ValueIgnoreEmpty(d.Get("plain_text")),
		"cipher_text":          utils.ValueIgnoreEmpty(d.Get("cipher_text")),
	}
	return bodyParams
}

func resourceKmsDataEncryptDecryptCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		dataEncryptDecryptHttpUrl = "v1.0/{project_id}/kms/{action}-data"
		dataEncryptDecryptProduct = "kms"
	)

	client, err := cfg.NewServiceClient(dataEncryptDecryptProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	action := d.Get("action").(string)
	if action != "encrypt" && action != "decrypt" {
		return diag.Errorf("the 'action' value is incorrect: %s, the value can only be 'encrypt' or 'decrypt'", action)
	}
	dataEncryptDecryptPath := client.Endpoint + dataEncryptDecryptHttpUrl
	dataEncryptDecryptPath = strings.ReplaceAll(dataEncryptDecryptPath, "{project_id}", client.ProjectID)
	dataEncryptDecryptPath = strings.ReplaceAll(dataEncryptDecryptPath, "{action}", action)

	dataEncryptDecryptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	dataEncryptDecryptOpt.JSONBody = utils.RemoveNil(buildDataEncryptDecryptBodyParams(d))
	dataEncryptDecryptResp, err := client.Request("POST", dataEncryptDecryptPath, &dataEncryptDecryptOpt)
	if err != nil {
		return diag.Errorf("error running %s operation: %s", action, err)
	}

	dataEncryptDecryptRespBody, err := utils.FlattenResponse(dataEncryptDecryptResp)
	if err != nil {
		return diag.FromErr(err)
	}

	keyId, err := jmespath.Search("key_id", dataEncryptDecryptRespBody)
	if err != nil || keyId == nil {
		return diag.Errorf("error running %s operation: ID is not found in API response", action)
	}

	d.SetId(keyId.(string))

	var mErr *multierror.Error
	if action == "encrypt" {
		cipherText, err := jmespath.Search("cipher_text", dataEncryptDecryptRespBody)
		if err != nil || cipherText == nil {
			return diag.Errorf("err searching cipher text when encrypting data: %s", err)
		}

		mErr = multierror.Append(
			mErr,
			d.Set("cipher_data", cipherText),
		)
	} else {
		plainText, err := jmespath.Search("plain_text", dataEncryptDecryptRespBody)
		if err != nil || plainText == nil {
			return diag.Errorf("err searching plain text when decrypting data: %s", err)
		}

		plainTextBase64, err := jmespath.Search("plain_text_base64", dataEncryptDecryptRespBody)
		if err != nil || plainTextBase64 == nil {
			return diag.Errorf("err searching base64 plain text when encrypting data: %s", err)
		}

		mErr = multierror.Append(
			mErr,
			d.Set("plain_data", plainText),
			d.Set("plain_text_base64", plainTextBase64),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsDataEncryptDecryptRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceKmsDataEncryptDecryptDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
