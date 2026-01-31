package lakeformation

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

var instanceRecoverNonUpdatableParams = []string{
	"instance_id",
}

// @API LakeFormation POST /v1/{project_id}/instances/{instance_id}/recover
func ResourceInstanceRecover() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceRecoverCreate,
		ReadContext:   resourceInstanceRecoverRead,
		UpdateContext: resourceInstanceRecoverUpdate,
		DeleteContext: resourceInstanceRecoverDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceRecoverNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instance needs to be recovered is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to be recovered from the recycle bin.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceInstanceRecoverCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v1/{project_id}/instances/{instance_id}/recover"
	)

	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating LakeFormation client: %s", err)
	}

	recoverPath := client.Endpoint + httpUrl
	recoverPath = strings.ReplaceAll(recoverPath, "{project_id}", client.ProjectID)
	recoverPath = strings.ReplaceAll(recoverPath, "{instance_id}", instanceId)

	recoverOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("POST", recoverPath, &recoverOpt)
	if err != nil {
		return diag.Errorf("error recovering instance (%s): %s", instanceId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceInstanceRecoverRead(ctx, d, meta)
}

func resourceInstanceRecoverRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceRecoverUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceRecoverDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for recovering the LakeFormation instance from the
recycle bin. Deleting this resource will not clear the corresponding request record, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
