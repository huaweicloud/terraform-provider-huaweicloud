package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var recycleBinVolumeDeleteNonUpdatableParams = []string{
	"volume_id",
}

// @API EVS DELETE /v3/{project_id}/recycle-bin-volumes/{volume_id}
func ResourceRecycleBinVolumeDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecycleBinVolumeDeleteCreate,
		ReadContext:   resourceRecycleBinVolumeDeleteRead,
		UpdateContext: resourceRecycleBinVolumeDeleteUpdate,
		DeleteContext: resourceRecycleBinVolumeDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(recycleBinVolumeDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceRecycleBinVolumeDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v3/{project_id}/recycle-bin-volumes/{volume_id}"
		product  = "evs"
		volumeId = d.Get("volume_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{volume_id}", volumeId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 202, 204,
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting EVS recycle bin volume: %s", err)
	}

	d.SetId(volumeId)

	return resourceRecycleBinVolumeDeleteRead(ctx, d, meta)
}

func resourceRecycleBinVolumeDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceRecycleBinVolumeDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceRecycleBinVolumeDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete EVS recycle bin volume. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the tf
    state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
