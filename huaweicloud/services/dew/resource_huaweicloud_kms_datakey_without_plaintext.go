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

var datakeyWithoutPlaintextNonUpdatableParams = []string{
	"key_id",
	"key_spec",
	"datakey_length",
	"sequence",
}

// @API KMS POST /v1.0/{project_id}/kms/create-datakey-without-plaintext
func ResourceKmsDatakeyWithoutPlaintext() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsDatakeyWithoutPlaintextCreate,
		ReadContext:   resourceKmsDatakeyWithoutPlaintextRead,
		UpdateContext: resourceKmsDatakeyWithoutPlaintextUpdate,
		DeleteContext: resourceKmsDatakeyWithoutPlaintextDelete,

		CustomizeDiff: config.FlexibleForceNew(datakeyWithoutPlaintextNonUpdatableParams),

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
			"key_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the generated key bit length.`,
			},
			"datakey_length": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the bit length of the key.`,
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
				Description: `The encrypted ciphertext, Base64 encoded.`,
			},
		},
	}
}

func buildDatakeyWithoutPlaintextBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":         d.Get("key_id"),
		"key_spec":       utils.ValueIgnoreEmpty(d.Get("key_spec")),
		"datakey_length": utils.ValueIgnoreEmpty(d.Get("datakey_length")),
		"sequence":       utils.ValueIgnoreEmpty(d.Get("sequence")),
	}
	return bodyParams
}

func resourceKmsDatakeyWithoutPlaintextCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/create-datakey-without-plaintext"
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
		JSONBody:         utils.RemoveNil(buildDatakeyWithoutPlaintextBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating KMS datakey without plaintext: %s", err)
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

func resourceKmsDatakeyWithoutPlaintextRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsDatakeyWithoutPlaintextUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsDatakeyWithoutPlaintextDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to create a datakey without plaintext.
Deleting this resource will not recover the created datakey, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
