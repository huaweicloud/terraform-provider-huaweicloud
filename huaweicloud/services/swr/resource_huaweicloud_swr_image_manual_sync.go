package swr

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

var swrImageManualSyncNonUpdatableParams = []string{
	"organization", "repository", "image_tag", "target_region", "target_organization", "override",
}

// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/sync_images
func ResourceSwrImageManualSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrImageManualSyncCreate,
		ReadContext:   resourceSwrImageManualSyncRead,
		UpdateContext: resourceSwrImageManualSyncUpdate,
		DeleteContext: resourceSwrImageManualSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(swrImageManualSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the organization.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the repository.`,
			},
			"image_tag": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the iamge tags.`,
			},
			"target_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the target region name.`,
			},
			"target_organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the target organization name.`,
			},
			"override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to overwrite. Default to **false**.`,
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

func resourceSwrImageManualSyncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)

	createHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/sync_images"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{namespace}", organization)
	createPath = strings.ReplaceAll(createPath, "{repository}", repository)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrImageManualSyncBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR image manual sync: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceSwrImageManualSyncRead(ctx, d, meta)
}

func buildCreateSwrImageManualSyncBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"remoteRegionId":  d.Get("target_region"),
		"remoteNamespace": d.Get("target_organization"),
		"imageTag":        d.Get("image_tag"),
		"override":        utils.ValueIgnoreEmpty(d.Get("override")),
	}
	return bodyParams
}

func resourceSwrImageManualSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrImageManualSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrImageManualSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR image manual sync resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
