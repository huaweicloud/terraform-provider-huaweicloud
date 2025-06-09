package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var decryptDatakeyNonUpdatableParams = []string{
	"key_id",
	"cipher_text",
	"datakey_cipher_length",
	"sequence",
}

// @API KMS POST /v1.0/{project_id}/kms/decrypt-datakey
func ResourceKmsDecryptDatakey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsDecryptDatakeyCreate,
		ReadContext:   resourceKmsDecryptDatakeyRead,
		UpdateContext: resourceKmsDecryptDatakeyUpdate,
		DeleteContext: resourceKmsDecryptDatakeyDelete,

		CustomizeDiff: config.FlexibleForceNew(decryptDatakeyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key ID.`,
			},
			"cipher_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DEK ciphertext and metadata in hexadecimal string.`,
			},
			"datakey_cipher_length": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the byte length of the DEK ciphertext.`,
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sequence number of the request message.`,
			},
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"data_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DEK plaintext in hexadecimal string.`,
			},
			"datakey_length": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The byte length of DEK plaintext.`,
			},
			"datakey_dgst": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SHA256 value of the DEK plaintext in hexadecimal string.`,
			},
		},
	}
}

func buildDecryptDatakeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":                d.Get("key_id"),
		"cipher_text":           d.Get("cipher_text"),
		"datakey_cipher_length": d.Get("datakey_cipher_length"),
		"sequence":              utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceKmsDecryptDatakeyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/decrypt-datakey"
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDecryptDatakeyBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error decrypting KMS datakey: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("data_key", utils.PathSearch("data_key", respBody, nil)),
		d.Set("datakey_length", utils.PathSearch("datakey_length", respBody, nil)),
		d.Set("datakey_dgst", utils.PathSearch("datakey_dgst", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsDecryptDatakeyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsDecryptDatakeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsDecryptDatakeyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to decrypt a datakey.
Deleting this resource will not recover the decrypted datakey, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
