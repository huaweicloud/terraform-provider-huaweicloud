package dew

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys
// @API DEW GET /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys
// @API DEW DELETE /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/{access_key_id}
// @API DEW POST /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/enable
// @API DEW POST /v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/disable
func ResourceCpcsAppAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCpcsAppAccessKeyCreate,
		ReadContext:   resourceCpcsAppAccessKeyRead,
		UpdateContext: resourceCpcsAppAccessKeyUpdate,
		DeleteContext: resourceCpcsAppAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCpcsAppAccessKeyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"app_id",
			"key_name",
			"access_key",
			"secret_key",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application ID to which the access key belongs.`,
			},
			"key_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the access key name.`,
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `Specifies the access key AK. If omitted, the system will automatically generate it.`,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `Specifies the access key SK. If omitted, the system will automatically generate it.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the status of the access key.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application to which the access key belongs.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the access key, UNIX timestamp in milliseconds.`,
			},
			"download_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time when the access key was downloaded, UNIX timestamp in milliseconds.`,
			},
			"is_downloaded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the access key has been downloaded.`,
			},
			"is_imported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the access key is imported.`,
			},
		},
	}
}

func buildCreateAppAccessKeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_name":   d.Get("key_name"),
		"access_key": utils.ValueIgnoreEmpty(d.Get("access_key")),
		"secret_key": utils.ValueIgnoreEmpty(d.Get("secret_key")),
	}
	return bodyParams
}

func updateCpcsAppAccessKeyStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, status string) error {
	httpUrl := ""
	switch status {
	case "disable":
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/disable"
	case "enable":
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/enable"
	default:
		return fmt.Errorf("invalid status: %s", status)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", d.Get("app_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"access_key_ids": []string{d.Id()},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceCpcsAppAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", d.Get("app_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAppAccessKeyBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DEW CPCS application access key: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKeyId := utils.PathSearch("access_key_id", respBody, "").(string)
	if accessKeyId == "" {
		return diag.Errorf("unable to find the DEW CPCS application access key ID from the API response")
	}
	d.SetId(accessKeyId)

	if d.Get("status").(string) == "disable" {
		if err := updateCpcsAppAccessKeyStatus(client, d, "disable"); err != nil {
			return diag.Errorf("error disabling DEW CPCS application access key in create operation: %s", err)
		}
	}

	return resourceCpcsAppAccessKeyRead(ctx, d, meta)
}

// The value of `key_name` can be used to accurately find the target value.
func QueryCpcsAppAccessKeyByAppIdAndKeyName(client *golangsdk.ServiceClient, appId, keyName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", appId)
	requestPath += fmt.Sprintf("?key_name=%s", keyName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	accessKeyDetail := utils.PathSearch("result|[0]", respBody, nil)
	if accessKeyDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return accessKeyDetail, nil
}

func resourceCpcsAppAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "kms"
		appId   = d.Get("app_id").(string)
		keyName = d.Get("key_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	accessKeyDetail, err := QueryCpcsAppAccessKeyByAppIdAndKeyName(client, appId, keyName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DEW CPCS application access key")
	}

	// Make sure that the ID value can be written back normally when importing.
	d.SetId(utils.PathSearch("access_key_id", accessKeyDetail, "").(string))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("key_name", utils.PathSearch("key_name", accessKeyDetail, nil)),
		d.Set("access_key", utils.PathSearch("access_key", accessKeyDetail, nil)),
		d.Set("app_name", utils.PathSearch("app_name", accessKeyDetail, nil)),
		d.Set("status", utils.PathSearch("status", accessKeyDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", accessKeyDetail, nil)),
		d.Set("download_time", utils.PathSearch("download_time", accessKeyDetail, nil)),
		d.Set("is_downloaded", utils.PathSearch("is_downloaded", accessKeyDetail, nil)),
		d.Set("is_imported", utils.PathSearch("is_imported", accessKeyDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCpcsAppAccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	if d.HasChange("status") {
		if err := updateCpcsAppAccessKeyStatus(client, d, d.Get("status").(string)); err != nil {
			return diag.Errorf("error updating DEW CPCS application access key status in update operation: %s", err)
		}
	}

	return resourceCpcsAppAccessKeyRead(ctx, d, meta)
}

func resourceCpcsAppAccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}/access-keys/{access_key_id}"
		product = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	// Disable the key first and then delete it.
	if d.Get("status").(string) == "enable" {
		if err := updateCpcsAppAccessKeyStatus(client, d, "disable"); err != nil {
			return diag.Errorf("error updating DEW CPCS application access key status in delete operation: %s", err)
		}
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", d.Get("app_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{access_key_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DEW CPCS application access key: %s", err)
	}

	return nil
}

func resourceCpcsAppAccessKeyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<app_id>/<key_name>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(
		d.Set("app_id", parts[0]),
		d.Set("key_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
