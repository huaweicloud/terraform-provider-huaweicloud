package iam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users
// @API IAM GET /v5/users/{user_id}
// @API IAM GET /v5/users/{user_id}/last-login
// @API IAM GET /v5/users/{user_id}/login-profile
func DataSourceV5Users() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5UserRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the user group to which the users belong.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the user.`,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the user is enabled.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
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
						"last_login_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last login time of the user.`,
						},
						"password_reset_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the password needs to be reset when the user logs in next time.`,
						},
						"password_expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The password expiration time of the user.`,
						},
					},
				},
				Description: `The list of users that matched filter parameters.`,
			},
		},
	}
}

func dataSourceV5UserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	users, err := listV5Users(client, buildListV5UsersQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying users: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	if userId, ok := d.GetOk("user_id"); ok {
		users = utils.PathSearch(fmt.Sprintf("[?user_id=='%s']", userId.(string)), users, make([]interface{}, 0)).([]interface{})
	}

	mErr := multierror.Append(
		d.Set("users", flattenV5Users(client, users)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting users fields: %s", err)
	}

	return nil
}

func buildListV5UsersQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("group_id"); ok {
		return fmt.Sprintf("?group_id=%v", v)
	}

	return ""
}

func flattenV5UsersTags(tags []interface{}) []interface{} {
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

func flattenV5Users(client *golangsdk.ServiceClient, users []interface{}) interface{} {
	if len(users) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		userId := utils.PathSearch("user_id", user, "").(string)
		userInfo, err := GetV5UserById(client, userId)
		if err != nil {
			return fmt.Errorf("error retrieving user (%s) information: %s", userId, err)
		}

		userLastLoginAt, err := getV5UserLastLoginTime(client, userId)
		if err != nil {
			// To avoid the error of this interface, causing the data source to fail to use normally, use log to record the error.
			log.Printf("[ERROR] error retrieving user (%s) last login time: %s", userId, err)
		}

		loginProfile, err := getV5LoginProfileConfig(client, userId)
		if err != nil {
			// For new created users, if the login password is not set, this interface will report an error,
			// so use log to record the error.
			// The error message as follows:
			// { "error_code": "IAM.0004", "error_msg": "Could not find login-profile: {user_id}" }
			log.Printf("[ERROR] error retrieving user (%s) login profile: %s", userId, err)
		}

		result = append(result, map[string]interface{}{
			"enabled":                 utils.PathSearch("enabled", user, nil),
			"user_name":               utils.PathSearch("user_name", user, nil),
			"description":             utils.PathSearch("description", user, nil),
			"is_root_user":            utils.PathSearch("is_root_user", user, nil),
			"created_at":              utils.PathSearch("created_at", user, nil),
			"urn":                     utils.PathSearch("urn", user, nil),
			"user_id":                 utils.PathSearch("user_id", user, nil),
			"tags":                    flattenV5UsersTags(utils.PathSearch("tags", userInfo, make([]interface{}, 0)).([]interface{})),
			"last_login_at":           userLastLoginAt,
			"password_reset_required": utils.PathSearch("password_reset_required", loginProfile, nil),
			"password_expires_at":     utils.PathSearch("password_expires_at", loginProfile, nil),
		})
	}

	return result
}
