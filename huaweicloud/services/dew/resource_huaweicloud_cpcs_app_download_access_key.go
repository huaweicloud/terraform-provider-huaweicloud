package dew

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

var appAccessKeyNonUpdatableParams = []string{"app_id", "access_key_id"}

// @API DEW GET /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/{access_key_id}
func ResourceCpcsAppDownloadAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCpcsAppDownloadAccessKeyCreate,
		ReadContext:   resourceCpcsAppDownloadAccessKeyRead,
		UpdateContext: resourceCpcsAppDownloadAccessKeyUpdate,
		DeleteContext: resourceCpcsAppDownloadAccessKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(appAccessKeyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_imported": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceCpcsAppDownloadAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/{access_key_id}"
		product     = "kms"
		appId       = d.Get("app_id").(string)
		accessKeyId = d.Get("access_key_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", appId)
	requestPath = strings.ReplaceAll(requestPath, "{access_key_id}", accessKeyId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error downloading DEW CPCS application access key: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(accessKeyId)

	mErr := multierror.Append(
		d.Set("access_key", utils.PathSearch("access_key", respBody, nil)),
		d.Set("secret_key", utils.PathSearch("secret_key", respBody, nil)),
		d.Set("key_name", utils.PathSearch("key_name", respBody, nil)),
		d.Set("is_imported", utils.PathSearch("is_imported", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DEW CPCS application access key attributes: %s", err)
	}

	return resourceCpcsAppDownloadAccessKeyRead(ctx, d, meta)
}

func resourceCpcsAppDownloadAccessKeyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCpcsAppDownloadAccessKeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceCpcsAppDownloadAccessKeyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to download CPC application access key.
Deleting this resource will not recover the downloaded access key, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
