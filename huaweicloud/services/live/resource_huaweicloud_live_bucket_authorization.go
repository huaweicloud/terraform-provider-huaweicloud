// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Live
// ---------------------------------------------------------------

package live

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live PUT /v1/{project_id}/obs/authority
func ResourceLiveBucketAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiveBucketAuthorizationCreate,
		ReadContext:   resourceLiveBucketAuthorizationRead,
		DeleteContext: resourceLiveBucketAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the bucket name of the OBS.`,
			},
		},
	}
}

func resourceLiveBucketAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createLiveBucketAuthorization: create Live bucket authorization
	var (
		createLiveBucketAuthorizationHttpUrl = "v1/{project_id}/obs/authority"
		createLiveBucketAuthorizationProduct = "live"
	)
	createLiveBucketAuthorizationClient, err := cfg.NewServiceClient(createLiveBucketAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	createLiveBucketAuthorizationPath := createLiveBucketAuthorizationClient.Endpoint + createLiveBucketAuthorizationHttpUrl
	createLiveBucketAuthorizationPath = strings.ReplaceAll(createLiveBucketAuthorizationPath, "{project_id}",
		createLiveBucketAuthorizationClient.ProjectID)

	createLiveBucketAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createLiveBucketAuthorizationOpt.JSONBody = utils.RemoveNil(buildLiveBucketAuthorizationBodyParams(d, 1))
	_, err = createLiveBucketAuthorizationClient.Request("PUT", createLiveBucketAuthorizationPath,
		&createLiveBucketAuthorizationOpt)
	if err != nil {
		return diag.Errorf("error creating Live bucket Authorization: %s", err)
	}

	d.SetId(d.Get("bucket").(string))

	return resourceLiveBucketAuthorizationRead(ctx, d, meta)
}

func buildLiveBucketAuthorizationBodyParams(d *schema.ResourceData, operation int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bucket":    utils.ValueIgnoreEmpty(d.Get("bucket")),
		"operation": operation,
	}
	return bodyParams
}

func resourceLiveBucketAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("bucket", d.Id()),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLiveBucketAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteLiveBucketAuthorization: Delete Live bucket Authorization
	var (
		deleteLiveBucketAuthorizationHttpUrl = "v1/{project_id}/obs/authority"
		deleteLiveBucketAuthorizationProduct = "live"
	)
	deleteLiveBucketAuthorizationClient, err := cfg.NewServiceClient(deleteLiveBucketAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Live Client: %s", err)
	}

	deleteLiveBucketAuthorizationPath := deleteLiveBucketAuthorizationClient.Endpoint + deleteLiveBucketAuthorizationHttpUrl
	deleteLiveBucketAuthorizationPath = strings.ReplaceAll(deleteLiveBucketAuthorizationPath, "{project_id}",
		deleteLiveBucketAuthorizationClient.ProjectID)

	deleteLiveBucketAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteLiveBucketAuthorizationOpt.JSONBody = utils.RemoveNil(buildLiveBucketAuthorizationBodyParams(d, 0))
	_, err = deleteLiveBucketAuthorizationClient.Request("PUT", deleteLiveBucketAuthorizationPath,
		&deleteLiveBucketAuthorizationOpt)
	if err != nil {
		return diag.Errorf("error deleting Live bucket Authorization: %s", err)
	}

	return nil
}
