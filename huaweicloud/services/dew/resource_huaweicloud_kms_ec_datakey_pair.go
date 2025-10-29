package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ecDatakeyPairNonUpdatableParams = []string{
	"key_id",
	"key_spec",
	"with_plain_text",
	"additional_authenticated_data",
	"sequence",
}

// @API DEW POST /v1.0/{project_id}/kms/create-ec-datakey-pair
func ResourceEcDatakeyPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEcDatakeyPairCreate,
		ReadContext:   resourceEcDatakeyPairRead,
		UpdateContext: resourceEcDatakeyPairUpdate,
		DeleteContext: resourceEcDatakeyPairDelete,

		CustomizeDiff: config.FlexibleForceNew(ecDatakeyPairNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_spec": {
				Type:     schema.TypeString,
				Required: true,
			},
			"with_plain_text": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"additional_authenticated_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"public_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key_cipher_text": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key_plain_text": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"wrapped_private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"ciphertext_recipient": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func buildEcDatakeyPairBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":                        d.Get("key_id"),
		"key_spec":                      d.Get("key_spec"),
		"with_plain_text":               d.Get("with_plain_text"),
		"additional_authenticated_data": utils.ValueIgnoreEmpty(d.Get("additional_authenticated_data")),
		"sequence":                      utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceEcDatakeyPairCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/create-ec-datakey-pair"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildEcDatakeyPairBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EC datakey pair: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("public_key", utils.PathSearch("public_key", respBody, nil)),
		d.Set("private_key_cipher_text", utils.PathSearch("private_key_cipher_text", respBody, nil)),
		d.Set("private_key_plain_text", utils.PathSearch("private_key_plain_text", respBody, nil)),
		d.Set("wrapped_private_key", utils.PathSearch("wrapped_private_key", respBody, nil)),
		d.Set("ciphertext_recipient", utils.PathSearch("ciphertext_recipient", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEcDatakeyPairRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceEcDatakeyPairUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceEcDatakeyPairDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
