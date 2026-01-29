package iam

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

// @API IAM POST /v5/users/{user_id}/access-keys
// @API IAM DELETE /v5/users/{user_id}/access-keys/{access_key_id}
// @API IAM GET /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys/{access_key_id}/last-used
// @API IAM PUT /v5/users/{user_id}/access-keys/{access_key_id}
var v5AccessKeyNonUpdatableParams = []string{"user_id"}

func ResourceV5AccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5AccessKeyCreate,
		ReadContext:   resourceV5AccessKeyRead,
		UpdateContext: resourceV5AccessKeyUpdate,
		DeleteContext: resourceV5AccessKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(v5AccessKeyNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5AccessKeyImportState,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The status of the access key.`,
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the generated access key.`,
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The generated secret access key.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the access key.`,
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the access key was last used.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5AccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	createAccessKeyHttpUrl := "v5/users/{user_id}/access-keys"
	createAccessKeyPath := iamClient.Endpoint + createAccessKeyHttpUrl
	createAccessKeyPath = strings.ReplaceAll(createAccessKeyPath, "{user_id}", userId)
	createAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccessKeyResp, err := iamClient.Request("POST", createAccessKeyPath, &createAccessKeyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM access key for user (%s): %s", userId, err)
	}

	createAccessKeyBody, err := utils.FlattenResponse(createAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKeyId := utils.PathSearch("access_key.access_key_id", createAccessKeyBody, "").(string)
	if accessKeyId == "" {
		return diag.Errorf("unable to find the IAM access key ID from the API response: %s", err)
	}

	d.SetId(accessKeyId)
	mErr := multierror.Append(
		d.Set("secret_access_key", utils.PathSearch("access_key.secret_access_key", createAccessKeyBody, "").(string)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving secret access key field to state: %s", err)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "active" {
		err = updateV5AccessKeyStatus(iamClient, userId, accessKeyId, v.(string))
		if err != nil {
			return diag.Errorf("error updating status of access key for user (%s): %s", userId, err)
		}
	}

	return resourceV5AccessKeyRead(ctx, d, meta)
}

func updateV5AccessKeyStatus(client *golangsdk.ServiceClient, userId string, accessKeyId string, status string) error {
	httpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{user_id}", userId)
	updatePath = strings.ReplaceAll(updatePath, "{access_key_id}", accessKeyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"status": status,
		},
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	getAccessKeyHttpUrl := "v5/users/{user_id}/access-keys"
	getAccessKeyPath := iamClient.Endpoint + getAccessKeyHttpUrl
	getAccessKeyPath = strings.ReplaceAll(getAccessKeyPath, "{user_id}", userId)
	getAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccessKeyResp, err := iamClient.Request("GET", getAccessKeyPath, &getAccessKeyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving IAM access key of the user (%s)", userId))
	}

	getAccessKeyRespBody, err := utils.FlattenResponse(getAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	accessKey := utils.PathSearch(fmt.Sprintf("access_keys[?access_key_id=='%s']|[0]", d.Id()), getAccessKeyRespBody, nil)
	if accessKey == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr := multierror.Append(nil,
		d.Set("access_key_id", utils.PathSearch("access_key_id", accessKey, nil)),
		d.Set("created_at", utils.PathSearch("created_at", accessKey, "").(string)),
		d.Set("status", utils.PathSearch("status", accessKey, "").(string)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM access key fields: %s", err)
	}

	accessKeyId := utils.PathSearch("access_key_id", accessKey, "").(string)
	if accessKeyId != "" {
		getLastUsedHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}/last-used"
		getLastUsedPath := iamClient.Endpoint + getLastUsedHttpUrl
		getLastUsedPath = strings.ReplaceAll(getLastUsedPath, "{user_id}", userId)
		getLastUsedPath = strings.ReplaceAll(getLastUsedPath, "{access_key_id}", accessKeyId)
		getLastUsedOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getLastUsedResp, err := iamClient.Request("GET", getLastUsedPath, &getLastUsedOpt)
		if err != nil {
			return diag.Errorf("error get IAM User last login time: %s", err)
		}
		lastUsed, err := utils.FlattenResponse(getLastUsedResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if e := d.Set("last_used_at", utils.PathSearch("access_key_last_used.last_used_at", lastUsed, nil)); e != nil {
			return diag.FromErr(e)
		}
	}
	return nil
}

func resourceV5AccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		userId      = d.Get("user_id").(string)
		accessKeyId = d.Id()
	)

	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = updateV5AccessKeyStatus(client, userId, accessKeyId, d.Get("status").(string))
	if err != nil {
		return diag.Errorf("error updating status of access key for user (%s): %s", userId, err)
	}

	return resourceV5AccessKeyRead(ctx, d, meta)
}

func resourceV5AccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Get("user_id").(string)
	)

	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	deleteAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	deleteAccessKeyPath := iamClient.Endpoint + deleteAccessKeyHttpUrl
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{user_id}", userId)
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{access_key_id}", d.Id())
	deleteAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteAccessKeyPath, &deleteAccessKeyOpt)
	if err != nil {
		return diag.Errorf("error deleting access key for user (%s): %s", userId, err)
	}
	return nil
}

func resourceV5AccessKeyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<user_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("user_id", parts[0])
}
