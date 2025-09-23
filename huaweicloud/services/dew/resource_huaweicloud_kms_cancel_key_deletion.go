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

var cancelKeyDeletionNonUpdatableParams = []string{
	"key_id",
	"sequence",
}

// @API DEW POST /v1.0/{project_id}/kms/cancel-key-deletion
func ResourceKmsCancelKeyDeletion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsCancelKeyDeletionCreate,
		ReadContext:   resourceKmsCancelKeyDeletionRead,
		UpdateContext: resourceKmsCancelKeyDeletionUpdate,
		DeleteContext: resourceKmsCancelKeyDeletionDelete,

		CustomizeDiff: config.FlexibleForceNew(cancelKeyDeletionNonUpdatableParams),

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
			"key_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the KMS key.`,
			},
		},
	}
}

func buildCancelKeyDeletionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":   d.Get("key_id"),
		"sequence": utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return bodyParams
}

func resourceKmsCancelKeyDeletionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/cancel-key-deletion"
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
		JSONBody:         utils.RemoveNil(buildCancelKeyDeletionBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error canceling the deletion of key: %s", err)
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
		d.Set("key_state", utils.PathSearch("key_state", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsCancelKeyDeletionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsCancelKeyDeletionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsCancelKeyDeletionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to cancel key deletion.
Deleting this resource will not recover the cancellation of the key deletion, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
