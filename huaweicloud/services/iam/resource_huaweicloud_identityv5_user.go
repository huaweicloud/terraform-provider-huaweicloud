package iam

import (
	"context"
	"fmt"
	"log"
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
// @API IAM GET /v5/users/{user_id}/last-login
// @API IAM PUT /v5/users/{user_id}
// @API IAM DELETE /v5/users/{user_id}
func ResourceV5User() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5UserCreate,
		ReadContext:   resourceV5UserRead,
		UpdateContext: resourceV5UserUpdate,
		DeleteContext: resourceV5UserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the user.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to enable the user.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the user.`,
			},
			"is_root_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user is root user.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the user.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the user.`,
			},
			"last_login_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last login time of the user.`,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
						"tag_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the tag.`,
						},
					},
				},
				Description: `The list of tags associated with the user.`,
			},
		},
	}
}

func resourceV5UserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createUserHttpUrl := "v5/users"
	createUserPath := client.Endpoint + createUserHttpUrl
	createUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateV5UserBodyParams(d)),
	}
	createUserResp, err := client.Request("POST", createUserPath, &createUserOpt)
	if err != nil {
		return diag.Errorf("error creating user: %s", err)
	}

	createUserBody, err := utils.FlattenResponse(createUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	userId := utils.PathSearch("user.user_id", createUserBody, "").(string)
	if userId == "" {
		return diag.Errorf("unable to find the user ID from the API response: %s", err)
	}

	d.SetId(userId)
	return resourceV5UserRead(ctx, d, meta)
}

func buildCreateV5UserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name").(string),
		"enabled":     utils.ValueIgnoreEmpty(d.Get("enabled")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func GetV5UserById(client *golangsdk.ServiceClient, userId string) (interface{}, error) {
	getHttpUrl := "v5/users/{user_id}"
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

	return utils.PathSearch("user", respBody, nil), nil
}

func resourceV5UserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	user, err := GetV5UserById(client, userId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting user (%s)", userId))
	}

	lastLoginAt, err := getV5UserLastLoginTime(client, userId)
	if err != nil {
		// To prevent the error from affecting the subsequent logic, use the log to record the error.
		log.Printf("[DEBUG] error getting user (%s) last login time: %s", userId, err)
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("user_name", user, nil)),
		d.Set("description", utils.PathSearch("description", user, nil)),
		d.Set("is_root_user", utils.PathSearch("is_root_user", user, nil)),
		d.Set("created_at", utils.PathSearch("created_at", user, nil)),
		d.Set("urn", utils.PathSearch("urn", user, nil)),
		d.Set("enabled", utils.PathSearch("enabled", user, nil)),
		d.Set("last_login_at", lastLoginAt),
		d.Set("tags", flattenV5UserTags(utils.PathSearch("tags", user, make([]interface{}, 0)).([]interface{}))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields for user (%s) resource: %s", userId, err)
	}

	return nil
}

func flattenV5UserTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	results := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		results = append(results, map[string]interface{}{
			"tag_key":   utils.PathSearch("tag_key", tag, nil),
			"tag_value": utils.PathSearch("tag_value", tag, nil),
		})
	}
	return results
}

func getV5UserLastLoginTime(client *golangsdk.ServiceClient, userId string) (string, error) {
	getHttpUrl := "v5/users/{user_id}/last-login"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{user_id}", userId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("user_last_login.last_login_at", respBody, "").(string), nil
}

func resourceV5UserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChanges("name", "description", "enabled") {
		updateUserHttpUrl := "v5/users/{user_id}"
		updateUserPath := client.Endpoint + updateUserHttpUrl
		updateUserPath = strings.ReplaceAll(updateUserPath, "{user_id}", userId)
		updateUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateV5UserBodyParams(d),
		}
		_, err := client.Request("PUT", updateUserPath, &updateUserOpt)
		if err != nil {
			return diag.Errorf("error updating user (%s): %s", userId, err)
		}
	}

	return resourceV5UserRead(ctx, d, meta)
}

func buildUpdateV5UserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_user_name":   d.Get("name").(string),
		"enabled":         d.Get("enabled"),
		"new_description": d.Get("description").(string),
	}
	return bodyParams
}

func resourceV5UserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteUserHttpUrl := "v5/users/{user_id}"
	deleteUserPath := client.Endpoint + deleteUserHttpUrl
	deleteUserPath = strings.ReplaceAll(deleteUserPath, "{user_id}", userId)
	deleteUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteUserPath, &deleteUserOpt)
	if err != nil {
		return diag.Errorf("error deleting user (%s): %s", userId, err)
	}

	return nil
}
