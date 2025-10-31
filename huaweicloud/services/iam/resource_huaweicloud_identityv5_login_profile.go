package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/users/{user_id}/login-profile
// @API IAM GET /v5/users/{user_id}/login-profile
// @API IAM PUT /v5/users/{user_id}/login-profile
// @API IAM DELETE /v5/users/{user_id}/login-profile
var loginProfileNonUpdatableParams = []string{"user_id"}

func ResourceIdentityV5LoginProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5LoginProfileCreate,
		ReadContext:   resourceIdentityV5LoginProfileRead,
		UpdateContext: resourceIdentityV5LoginProfileUpdate,
		DeleteContext: resourceIdentityV5LoginProfileDelete,

		CustomizeDiff: config.FlexibleForceNew(loginProfileNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password_reset_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"password_expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityV5LoginProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	userId := d.Get("user_id").(string)
	createLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	createLoginProfilePath := iamClient.Endpoint + createLoginProfileHttpUrl
	createLoginProfilePath = strings.ReplaceAll(createLoginProfilePath, "{user_id}", userId)
	createLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateLoginProfileBodyParams(d),
	}
	createLoginProfileResp, err := iamClient.Request("POST", createLoginProfilePath, &createLoginProfileOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error get IAM login profile")
	}
	createLoginProfileBody, err := utils.FlattenResponse(createLoginProfileResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("login_profile.user_id", createLoginProfileBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating IAM login profile: user_id is not found in API response: %s", err)
	}
	d.SetId(id)
	return resourceIdentityV5LoginProfileRead(ctx, d, meta)
}

func buildCreateLoginProfileBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password":                d.Get("password").(string),
		"password_reset_required": d.Get("password_reset_required"),
	}
	return bodyParams
}

func resourceIdentityV5LoginProfileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	getLoginProfilePath := iamClient.Endpoint + getLoginProfileHttpUrl
	getLoginProfilePath = strings.ReplaceAll(getLoginProfilePath, "{user_id}", d.Id())
	getLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLoginProfileResp, err := iamClient.Request("GET", getLoginProfilePath, &getLoginProfileOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM login profile")
	}
	getLoginProfileRespBody, err := utils.FlattenResponse(getLoginProfileResp)
	if err != nil {
		return diag.FromErr(err)
	}
	loginProfile := utils.PathSearch("login_profile", getLoginProfileRespBody, nil)
	if loginProfile == nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM login profile : profile is not found in API response")
	}
	mErr := multierror.Append(nil,
		d.Set("user_id", utils.PathSearch("user_id", loginProfile, nil)),
		d.Set("password_reset_required", utils.PathSearch("password_reset_required", loginProfile, nil)),
		d.Set("created_at", utils.PathSearch("created_at", loginProfile, nil)),
		d.Set("password_expires_at", utils.PathSearch("password_expires_at", loginProfile, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting login profile fields: %s", err)
	}
	return nil
}

func resourceIdentityV5LoginProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	updateChanges := []string{
		"password",
		"password_reset_required",
	}
	if d.HasChanges(updateChanges...) {
		updateLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
		updateLoginProfilePath := iamClient.Endpoint + updateLoginProfileHttpUrl
		updateLoginProfilePath = strings.ReplaceAll(updateLoginProfilePath, "{user_id}", d.Id())
		baseBody := map[string]interface{}{}
		if d.HasChange("password") {
			baseBody["password"] = d.Get("password").(string)
		}
		if d.HasChange("password_reset_required") {
			baseBody["password_reset_required"] = d.Get("password_reset_required")
		}
		updateLoginProfileOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         baseBody,
		}
		_, err := iamClient.Request("PUT", updateLoginProfilePath, &updateLoginProfileOpt)
		if err != nil {
			return diag.Errorf("error updating IAM login profile: %s", err)
		}
	}
	return resourceIdentityV5LoginProfileRead(ctx, d, meta)
}

func resourceIdentityV5LoginProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	deleteLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	deleteLoginProfilePath := iamClient.Endpoint + deleteLoginProfileHttpUrl
	deleteLoginProfilePath = strings.ReplaceAll(deleteLoginProfilePath, "{user_id}", d.Id())
	deleteLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteLoginProfilePath, &deleteLoginProfileOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM login profile: %s", err)
	}
	return nil
}
