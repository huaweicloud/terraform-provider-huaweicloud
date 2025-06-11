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

var encryptDatakeyNonUpdatableParams = []string{
	"key_id",
	"plain_text",
	"datakey_plain_length",
	"sequence",
}

// @API DEW POST /v1.0/{project_id}/kms/encrypt-datakey
func ResourceKmsEncryptDatakey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsEncryptDatakeyCreate,
		ReadContext:   resourceKmsEncryptDatakeyRead,
		UpdateContext: resourceKmsEncryptDatakeyUpdate,
		DeleteContext: resourceKmsEncryptDatakeyDelete,

		CustomizeDiff: config.FlexibleForceNew(encryptDatakeyNonUpdatableParams),

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
			"plain_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the plaintext of data encryption key.`,
			},
			"datakey_plain_length": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the byte length of the DEK plaintext.`,
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
			"cipher_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DEK ciphertext in hexadecimal.`,
			},
		},
	}
}

func buildEncryptDatakeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":               d.Get("key_id"),
		"plain_text":           d.Get("plain_text"),
		"datakey_plain_length": d.Get("datakey_plain_length"),
		"sequence":             utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceKmsEncryptDatakeyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/encrypt-datakey"
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
		JSONBody:         utils.RemoveNil(buildEncryptDatakeyBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error encrypting KMS datakey: %s", err)
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
		d.Set("cipher_text", utils.PathSearch("cipher_text", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsEncryptDatakeyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsEncryptDatakeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsEncryptDatakeyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to encrypt a datakey.
Deleting this resource will not recover the encrypted datakey, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
