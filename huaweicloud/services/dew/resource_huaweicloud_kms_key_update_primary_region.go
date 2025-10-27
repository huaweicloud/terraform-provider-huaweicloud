package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var primaryRegionNonUpdatableParams = []string{"key_id", "primary_region"}

// @API DEW PUT /v2/{project_id}/kms/keys/{key_id}/update-primary-region
func ResourceKeyUpdatePrimaryRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyUpdatePrimaryRegionCreate,
		ReadContext:   resourceKeyUpdatePrimaryRegionRead,
		UpdateContext: resourceKeyUpdatePrimaryRegionUpdate,
		DeleteContext: resourceKeyUpdatePrimaryRegionDelete,

		CustomizeDiff: config.FlexibleForceNew(primaryRegionNonUpdatableParams),

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
			"primary_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildKeyUpdatePrimaryRegionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"primary_region": d.Get("primary_region"),
	}

	return bodyParams
}

func resourceKeyUpdatePrimaryRegionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/kms/keys/{key_id}/update-primary-region"
		keyId   = d.Get("key_id").(string)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{key_id}", keyId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildKeyUpdatePrimaryRegionBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating KMS key primary region: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	return nil
}

func resourceKeyUpdatePrimaryRegionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKeyUpdatePrimaryRegionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKeyUpdatePrimaryRegionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
