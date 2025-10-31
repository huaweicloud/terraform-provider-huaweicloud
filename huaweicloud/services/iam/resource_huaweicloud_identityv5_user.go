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

// @API IAM POST /v5/users
// @API IAM GET /v5/users/{user_id}
// @API IAM PUT /v5/users/{user_id}
// @API IAM DELETE /v5/users/{user_id}
// @API IAM GET /v5/users/{user_id}/last-login
func ResourceIdentityV5User() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5UserCreate,
		ReadContext:   resourceIdentityV5UserRead,
		UpdateContext: resourceIdentityV5UserUpdate,
		DeleteContext: resourceIdentityV5UserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_root_user": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_login_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityV5UserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	createUserHttpUrl := "v5/users"
	createUserPath := iamClient.Endpoint + createUserHttpUrl
	createUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateUserBodyParams(d)),
	}
	createUserResp, err := iamClient.Request("POST", createUserPath, &createUserOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error get IAM user")
	}
	createUserBody, err := utils.FlattenResponse(createUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("user.user_id", createUserBody, "").(string)
	if id == "" {
		return common.CheckDeletedDiag(d, err, "error getting IAM user: user is not found in API response")
	}
	d.SetId(id)
	return resourceIdentityV5UserRead(ctx, d, meta)
}

func buildCreateUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name").(string),
		"enabled":     utils.ValueIgnoreEmpty(d.Get("enabled")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceIdentityV5UserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getUserHttpUrl := "v5/users/{user_id}"
	getUserPath := iamClient.Endpoint + getUserHttpUrl
	getUserPath = strings.ReplaceAll(getUserPath, "{user_id}", d.Id())
	getUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getUserResp, err := iamClient.Request("GET", getUserPath, &getUserOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error get IAM user")
	}
	getUserRespBody, err := utils.FlattenResponse(getUserResp)
	if err != nil {
		return diag.FromErr(err)
	}
	user := utils.PathSearch("user", getUserRespBody, nil)
	if user == nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM user: user is not found in API response")
	}

	tags := listUserTagInfo(user)
	getLastLoginHttpUrl := "v5/users/{user_id}/last-login"
	getLastLoginPath := iamClient.Endpoint + getLastLoginHttpUrl
	getLastLoginPath = strings.ReplaceAll(getLastLoginPath, "{user_id}", d.Id())
	getLastLoginOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLastLoginResp, err := iamClient.Request("GET", getLastLoginPath, &getLastLoginOpt)
	if err != nil {
		return diag.Errorf("error get IAM User last login time: %s", err)
	}
	userLastLogin, err := utils.FlattenResponse(getLastLoginResp)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("user_name", user, nil)),
		d.Set("description", utils.PathSearch("description", user, nil)),
		d.Set("is_root_user", utils.PathSearch("is_root_user", user, nil)),
		d.Set("created_at", utils.PathSearch("created_at", user, nil)),
		d.Set("urn", utils.PathSearch("urn", user, nil)),
		d.Set("enabled", utils.PathSearch("enabled", user, nil)),
		d.Set("last_login_at", utils.PathSearch("last_login_at", userLastLogin, nil)),
		d.Set("tags", tags),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM user fields: %s", err)
	}
	return nil
}

func listUserTagInfo(user interface{}) []interface{} {
	if user == nil {
		return nil
	}
	curJson := utils.PathSearch("tags", user, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	results := make([]interface{}, 0, len(curArray))
	for _, tag := range curArray {
		results = append(results, map[string]interface{}{
			"tag_key":   utils.PathSearch("tag_key", tag, nil),
			"tag_value": utils.PathSearch("tag_value", tag, nil),
		})
	}
	return results
}

func resourceIdentityV5UserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	updateChanges := []string{
		"name",
		"description",
		"enabled",
	}
	if d.HasChanges(updateChanges...) {
		updateUserHttpUrl := "v5/users/{user_id}"
		updateUserPath := iamClient.Endpoint + updateUserHttpUrl
		updateUserPath = strings.ReplaceAll(updateUserPath, "{user_id}", d.Id())
		updateUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateUserBodyParams(d),
		}
		_, err := iamClient.Request("PUT", updateUserPath, &updateUserOpt)
		if err != nil {
			return diag.Errorf("error updating IAM user: %s", err)
		}
	}
	return resourceIdentityV5UserRead(ctx, d, meta)
}

func buildUpdateUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_user_name":   d.Get("name").(string),
		"enabled":         d.Get("enabled"),
		"new_description": d.Get("description").(string),
	}
	return bodyParams
}

func resourceIdentityV5UserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	deleteUserHttpUrl := "v5/users/{user_id}"
	deleteUserPath := iamClient.Endpoint + deleteUserHttpUrl
	deleteUserPath = strings.ReplaceAll(deleteUserPath, "{user_id}", d.Id())
	deleteUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteUserPath, &deleteUserOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM user: %s", err)
	}
	return nil
}
