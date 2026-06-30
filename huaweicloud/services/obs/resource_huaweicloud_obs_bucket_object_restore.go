package obs

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var bucketObjectRestoreNonUpdatableParams = []string{
	"bucket",
	"key",
	"days",
	"tier",
	"version_id",
}

// @API OBS POST /{ObjectName}?restore
func ResourceBucketObjectRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBucketObjectRestoreCreate,
		ReadContext:   resourceBucketObjectRestoreRead,
		UpdateContext: resourceBucketObjectRestoreUpdate,
		DeleteContext: resourceBucketObjectRestoreDelete,

		CustomizeDiff: config.FlexibleForceNew(bucketObjectRestoreNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the archived object is located.`,
			},

			// Required parameters.
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the bucket to restore the object from.`,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the object to restore.`,
			},
			"days": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of days for which the restored object copy is valid.`,
			},
			"tier": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The restore option.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version ID of the object to restore.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func resourceBucketObjectRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf      = meta.(*config.Config)
		bucket    = d.Get("bucket").(string)
		key       = d.Get("key").(string)
		days      = d.Get("days").(int)
		versionId = d.Get("version_id").(string)
	)

	obsClient, err := conf.ObjectStorageClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating OBS client: %s", err)
	}

	input := &obs.RestoreObjectInput{
		Bucket:    bucket,
		Key:       key,
		Days:      days,
		VersionId: versionId,
	}

	tier := d.Get("tier").(string)
	switch tier {
	case "expedited":
		input.Tier = obs.RestoreTierExpedited
	case "standard":
		input.Tier = obs.RestoreTierStandard
	}

	_, err = obsClient.RestoreObject(input)
	if err != nil {
		return diag.Errorf("error restoring object %s in bucket %s: %s", key, bucket, err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceBucketObjectRestoreRead(ctx, d, meta)
}

func resourceBucketObjectRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBucketObjectRestoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBucketObjectRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for restoring an archived OBS object. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
