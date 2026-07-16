package gaussdb

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

var geminiDBDownloadBackupsAuthorizeParams = []string{
	"backup_id",
}

// @API GaussDB POST /v3/{project_id}/backups/{backup_id}/download/authorization
func ResourceGaussDBDownloadBackupsAuthorize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBDownloadBackupsAuthorizeCreate,
		ReadContext:   resourceGaussDBDownloadBackupsAuthorizeRead,
		UpdateContext: resourceGaussDBDownloadBackupsAuthorizeUpdate,
		DeleteContext: resourceGaussDBDownloadBackupsAuthorizeDelete,

		CustomizeDiff: config.FlexibleForceNew(geminiDBDownloadBackupsAuthorizeParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_paths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGaussDBDownloadBackupsAuthorizeCreate(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/{backup_id}/download/authorization"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	backupID := d.Get("backup_id").(string)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{backup_id}", backupID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error authorizing to download GaussDB backup (%s): %s", backupID, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(backupID)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("bucket", utils.PathSearch("bucket", createRespBody, nil)),
		d.Set("file_paths", utils.PathSearch("file_paths", createRespBody, make([]interface{}, 0))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBDownloadBackupsAuthorizeUpdate(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBDownloadBackupsAuthorizeRead(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBDownloadBackupsAuthorizeDelete(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	errorMsg := "This resource is a one-time action resource for authorizing users to download backups. " +
		"Deleting this resource will not revoke the authorization, but will only remove the resource " +
		"information from the tfstate file."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource deletion is not supported",
			Detail:   errorMsg,
		},
	}
}
