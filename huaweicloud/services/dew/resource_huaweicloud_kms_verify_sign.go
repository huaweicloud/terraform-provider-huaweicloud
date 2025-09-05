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

var verifySignNonUpdatableParams = []string{
	"key_id",
	"message",
	"signature",
	"signing_algorithm",
	"message_type",
	"sequence",
}

// @API DEW POST /v1.0/{project_id}/kms/verify
func ResourceKmsVerifySign() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsVerifySignCreate,
		ReadContext:   resourceKmsVerifySignRead,
		UpdateContext: resourceKmsVerifySignUpdate,
		DeleteContext: resourceKmsVerifySignDelete,

		CustomizeDiff: config.FlexibleForceNew(verifySignNonUpdatableParams),

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
			"message": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the message digest or message.`,
			},
			"signature": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signature value to be verified.`,
			},
			"signing_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signature algorithm.`,
			},
			"message_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the message.`,
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
			"signature_valid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The validity of signature verification.`,
			},
		},
	}
}

func buildVerifySignBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":            d.Get("key_id"),
		"message":           d.Get("message"),
		"signature":         d.Get("signature"),
		"signing_algorithm": d.Get("signing_algorithm"),
		"message_type":      utils.ValueIgnoreEmpty(d.Get("message_type")),
		"sequence":          utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceKmsVerifySignCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/verify"
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
		JSONBody:         utils.RemoveNil(buildVerifySignBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error running verifying signature : %s", err)
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
		d.Set("signature_valid", utils.PathSearch("signature_valid", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsVerifySignRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsVerifySignUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsVerifySignDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to verify signature.
Deleting this resource will not recover verify, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
