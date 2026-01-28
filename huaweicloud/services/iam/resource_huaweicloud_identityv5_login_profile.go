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

var v5LoginProfileNonUpdatableParams = []string{"user_id"}

// @API IAM POST /v5/users/{user_id}/login-profile
// @API IAM GET /v5/users/{user_id}/login-profile
// @API IAM PUT /v5/users/{user_id}/login-profile
// @API IAM DELETE /v5/users/{user_id}/login-profile
func ResourceV5LoginProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5LoginProfileCreate,
		ReadContext:   resourceV5LoginProfileRead,
		UpdateContext: resourceV5LoginProfileUpdate,
		DeleteContext: resourceV5LoginProfileDelete,

		CustomizeDiff: config.FlexibleForceNew(v5LoginProfileNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user.`,
			},
			"password_reset_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the user needs to reset the password at the next login.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: `The password of the user login.`,
			},
			// Attributes.
			"password_expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The password expiration time of the user.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the login profile.`,
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

func resourceV5LoginProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	createLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	createLoginProfilePath := client.Endpoint + createLoginProfileHttpUrl
	createLoginProfilePath = strings.ReplaceAll(createLoginProfilePath, "{user_id}", userId)
	createLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateV5LoginProfileBodyParams(d),
	}
	createLoginProfileResp, err := client.Request("POST", createLoginProfilePath, &createLoginProfileOpt)
	if err != nil {
		return diag.Errorf("error creating login profile for user (%s): %s", userId, err)
	}

	createLoginProfileBody, err := utils.FlattenResponse(createLoginProfileResp)
	if err != nil {
		return diag.FromErr(err)
	}

	loginProfileUserId := utils.PathSearch("login_profile.user_id", createLoginProfileBody, "").(string)
	if loginProfileUserId == "" {
		return diag.Errorf("unable to find user ID from API response: %s", err)
	}

	d.SetId(loginProfileUserId)
	return resourceV5LoginProfileRead(ctx, d, meta)
}

func buildCreateV5LoginProfileBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": d.Get("password").(string),
		// For the API, this parameter is required.
		"password_reset_required": d.Get("password_reset_required"),
	}
	return bodyParams
}

func getV5LoginProfileConfig(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	getHttpUrl := "v5/users/{user_id}/login-profile"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{user_id}", userId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("login_profile", respBody, nil), nil
}

func GetV5LoginProfile(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	loginProfile, err := getV5LoginProfileConfig(client, userId)
	if err != nil {
		return nil, err
	}

	if loginProfile == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/users/{user_id}/login-profile",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the login profile of the user (%s) does not exist", userId)),
			},
		}
	}

	return loginProfile, nil
}

func resourceV5LoginProfileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	loginProfile, err := GetV5LoginProfile(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting login profile")
	}

	if loginProfile == nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find login profile of the user (%s) in API response", d.Id()))
	}

	mErr := multierror.Append(
		d.Set("user_id", utils.PathSearch("user_id", loginProfile, nil)),
		d.Set("password_reset_required", utils.PathSearch("password_reset_required", loginProfile, nil)),
		d.Set("created_at", utils.PathSearch("created_at", loginProfile, nil)),
		d.Set("password_expires_at", utils.PathSearch("password_expires_at", loginProfile, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5LoginProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	updateChanges := []string{
		"password",
		"password_reset_required",
	}
	if d.HasChanges(updateChanges...) {
		updateLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
		updateLoginProfilePath := client.Endpoint + updateLoginProfileHttpUrl
		updateLoginProfilePath = strings.ReplaceAll(updateLoginProfilePath, "{user_id}", userId)
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
		_, err := client.Request("PUT", updateLoginProfilePath, &updateLoginProfileOpt)
		if err != nil {
			return diag.Errorf("error updating login profile of the user (%s): %s", userId, err)
		}
	}

	return resourceV5LoginProfileRead(ctx, d, meta)
}

func resourceV5LoginProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteLoginProfileHttpUrl := "v5/users/{user_id}/login-profile"
	deleteLoginProfilePath := client.Endpoint + deleteLoginProfileHttpUrl
	deleteLoginProfilePath = strings.ReplaceAll(deleteLoginProfilePath, "{user_id}", userId)
	deleteLoginProfileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteLoginProfilePath, &deleteLoginProfileOpt)
	if err != nil {
		return diag.Errorf("error deleting login profile of the user (%s): %s", userId, err)
	}

	return nil
}
