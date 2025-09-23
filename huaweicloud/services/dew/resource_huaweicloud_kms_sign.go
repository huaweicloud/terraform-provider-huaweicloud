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

var signNonUpdatableParams = []string{
	"key_id",
	"message",
	"signing_algorithm",
	"message_type",
	"sequence",
}

// @API DEW POST /v1.0/{project_id}/kms/sign
func ResourceKmsSign() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsSignCreate,
		ReadContext:   resourceKmsSignRead,
		UpdateContext: resourceKmsSignUpdate,
		DeleteContext: resourceKmsSignDelete,

		CustomizeDiff: config.FlexibleForceNew(signNonUpdatableParams),

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
			"signature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The signature value.`,
			},
		},
	}
}

func buildSignBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":            d.Get("key_id"),
		"message":           d.Get("message"),
		"signing_algorithm": d.Get("signing_algorithm"),
		"message_type":      utils.ValueIgnoreEmpty(d.Get("message_type")),
		"sequence":          utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceKmsSignCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/sign"
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
		JSONBody:         utils.RemoveNil(buildSignBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error running sign : %s", err)
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
		d.Set("signature", utils.PathSearch("signature", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsSignRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsSignUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsSignDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to sign message.
Deleting this resource will not recover sign, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
