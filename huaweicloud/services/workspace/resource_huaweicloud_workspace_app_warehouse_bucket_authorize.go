package workspace

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

var appWarehouseBucketAuthorizeNonUpdatableParams = []string{"bucket_name"}

// @API Workspace POST /v1/{project_id}/app-warehouse/bucket
func ResourceAppWarehouseBucketAuthorize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppWarehouseBucketAuthorizeCreate,
		ReadContext:   resourceAppWarehouseBucketAuthorizeRead,
		UpdateContext: resourceAppWarehouseBucketAuthorizeUpdate,
		DeleteContext: resourceAppWarehouseBucketAuthorizeDelete,

		CustomizeDiff: config.FlexibleForceNew(appWarehouseBucketAuthorizeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the app repository bucket is located.`,
			},

			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the bucket to be assigned.`,
			},
		},
	}
}

func buildAppWarehouseBucketAuthorizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"bucket_name": utils.ValueIgnoreEmpty(d.Get("bucket_name")),
	}
}

func resourceAppWarehouseBucketAuthorizeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-warehouse/bucket"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAppWarehouseBucketAuthorizeBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error assigning app repository bucket: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAppWarehouseBucketAuthorizeRead(ctx, d, meta)
}

func resourceAppWarehouseBucketAuthorizeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppWarehouseBucketAuthorizeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppWarehouseBucketAuthorizeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for assigning an app repository bucket.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
