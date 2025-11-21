package iam

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
var accessKeyV5NonUpdatableParams = []string{"user_id"}

func ResourceIdentityAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5AccessKeyCreate,
		ReadContext:   resourceIdentityV5AccessKeyRead,
		UpdateContext: resourceIdentityV5AccessKeyUpdate,
		DeleteContext: resourceIdentityV5AccessKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(accessKeyV5NonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityV5AccessKeyImportState,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_used_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityV5AccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	createAccessKeyBody, err := utils.FlattenResponse(createAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	accessKeyId := utils.PathSearch("access_key.access_key_id", createAccessKeyBody, "").(string)
	if accessKeyId == "" {
		return diag.Errorf("error getting IAM access key id: %s", err)
	}
	d.SetId(accessKeyId)
	mErr := multierror.Append(nil,
		d.Set("secret_access_key", utils.PathSearch("access_key.secret_access_key", createAccessKeyBody, "").(string)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM access key fields: %s", err)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "active" {
		updateAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
		updateAccessKeyPath := iamClient.Endpoint + updateAccessKeyHttpUrl
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{user_id}", userId)
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{access_key_id}", accessKeyId)
		updateAccessKeyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"status": d.Get("status").(string),
			},
		}
		_, err := iamClient.Request("PUT", updateAccessKeyPath, &updateAccessKeyOpt)
		if err != nil {
			return diag.Errorf("error updating IAM access key: %s", err)
		}
	}
	return resourceIdentityV5AccessKeyRead(ctx, d, meta)
}

func resourceIdentityV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "error get IAM user access key")
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

func resourceIdentityV5AccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	accessKeyId := d.Id()
	if d.HasChanges("status") {
		updateAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
		updateAccessKeyPath := iamClient.Endpoint + updateAccessKeyHttpUrl
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{user_id}", d.Get("user_id").(string))
		updateAccessKeyPath = strings.ReplaceAll(updateAccessKeyPath, "{access_key_id}", accessKeyId)
		updateAccessKeyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"status": d.Get("status").(string),
			},
		}
		_, err := iamClient.Request("PUT", updateAccessKeyPath, &updateAccessKeyOpt)
		if err != nil {
			return diag.Errorf("error updating IAM access key: %s", err)
		}
	}
	return resourceIdentityV5AccessKeyRead(ctx, d, meta)
}

func resourceIdentityV5AccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	deleteAccessKeyHttpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	deleteAccessKeyPath := iamClient.Endpoint + deleteAccessKeyHttpUrl
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{user_id}", d.Get("user_id").(string))
	deleteAccessKeyPath = strings.ReplaceAll(deleteAccessKeyPath, "{access_key_id}", d.Id())
	deleteAccessKeyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteAccessKeyPath, &deleteAccessKeyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM access key: %s", err)
	}
	return nil
}

func resourceIdentityV5AccessKeyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <user_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("user_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
