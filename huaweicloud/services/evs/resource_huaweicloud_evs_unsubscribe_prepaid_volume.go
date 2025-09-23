package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var unsubscribePrepaidVolumeNonUpdatableParams = []string{
	"volume_ids",
}

// @API EVS POST /v2/{project_id}/cloudvolumes/unsubscribe
func ResourceUnsubscribePrepaidVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnsubscribePrepaidVolumeCreate,
		ReadContext:   resourceUnsubscribePrepaidVolumeRead,
		UpdateContext: resourceUnsubscribePrepaidVolumeUpdate,
		DeleteContext: resourceUnsubscribePrepaidVolumeDelete,

		CustomizeDiff: config.FlexibleForceNew(unsubscribePrepaidVolumeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildUnsubscribePrepaidVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	volumeIdsRaw := d.Get("volume_ids").([]interface{})
	if len(volumeIdsRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"volume_ids": volumeIdsRaw,
	}

	return bodyParams
}

func resourceUnsubscribePrepaidVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/cloudvolumes/unsubscribe"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUnsubscribePrepaidVolumeBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error unsubscribing prepaid EVS volume: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceUnsubscribePrepaidVolumeRead(ctx, d, meta)
}

func resourceUnsubscribePrepaidVolumeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceUnsubscribePrepaidVolumeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceUnsubscribePrepaidVolumeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to unsubscribe prepaid volume.
Deleting this resource will not reset the unsubscribe volume, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
