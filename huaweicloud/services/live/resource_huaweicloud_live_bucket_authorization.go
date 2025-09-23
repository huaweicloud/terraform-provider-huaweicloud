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
)

// @API Live PUT /v1/{project_id}/obs/authority
func ResourceBucketAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBucketAuthorizationCreate,
		ReadContext:   resourceBucketAuthorizationRead,
		DeleteContext: resourceBucketAuthorizationDelete,
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

func resourceBucketAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/obs/authority"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBucketAuthorizationBodyParams(d, 1),
	}
	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating Live bucket authorization: %s", err)
	}

	d.SetId(d.Get("bucket").(string))

	return resourceBucketAuthorizationRead(ctx, d, meta)
}

func buildBucketAuthorizationBodyParams(d *schema.ResourceData, operation int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bucket":    d.Get("bucket"),
		"operation": operation,
	}
	return bodyParams
}

func resourceBucketAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceBucketAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/obs/authority"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBucketAuthorizationBodyParams(d, 0),
	}
	_, err = client.Request("PUT", requestPath,
		&requestOpt)
	if err != nil {
		return diag.Errorf("error deleting Live bucket authorization: %s", err)
	}

	return nil
}
