package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/users
func DataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the users are located.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name to be queried.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user description for fuzzy matching.`,
			},
			"active_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The activation type of the user.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user group name for exact matching.`,
			},
			"is_query_total_desktops": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query the number of desktops bound to the user.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userSchema(),
				Description: `The list of users that matched filter parameters.`,
			},
		},
	}
}

func userSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user.`,
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SID of the user.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of user.`,
			},
			"user_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address of the user.`,
			},
			"total_desktops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of desktops bound to the user.`,
			},
			"user_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The phone number of the user.`,
			},
			"active_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The activation type of the user.`,
			},
			"is_pre_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user is a pre-created user.`,
			},
			"account_expires": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account expired time.`,
			},
			"password_never_expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the password never expires.`,
			},
			"account_expired": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the account has expired.`,
			},
			"enable_change_password": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user is allowed to change password.`,
			},
			"next_login_change_password": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the password needs to be reset on next login.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the user.`,
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the account is locked.`,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the account is disabled.`,
			},
			"share_space_subscription": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user has subscribed to collaboration.`,
			},
			"share_space_desktops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of collaboration desktops bound to the user.`,
			},
			"group_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of group name that the user has joined.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the enterprise project.`,
			},
			"user_info_map": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user information mapping, including user service level, operation mode and type.`,
			},
		},
	}
}

func buildListUsersQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("user_name"); ok {
		res = fmt.Sprintf("%s&user_name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("active_type"); ok {
		res = fmt.Sprintf("%s&active_type=%v", res, v)
	}
	if v, ok := d.GetOk("group_name"); ok {
		res = fmt.Sprintf("%s&group_name=%v", res, v)
	}
	v := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_query_total_desktops")
	if v != nil {
		res = fmt.Sprintf("%s&is_query_total_desktops=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	return res
}

func listUsers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/users?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildListUsersQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%s", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		users := utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		offset += len(users)
	}
	return result, nil
}

func flattenUsers(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		user := map[string]interface{}{
			"id":                         utils.PathSearch("id", item, nil),
			"sid":                        utils.PathSearch("sid", item, nil),
			"user_name":                  utils.PathSearch("user_name", item, nil),
			"user_email":                 utils.PathSearch("user_email", item, nil),
			"total_desktops":             utils.PathSearch("total_desktops", item, nil),
			"user_phone":                 utils.PathSearch("user_phone", item, nil),
			"active_type":                utils.PathSearch("active_type", item, nil),
			"is_pre_user":                utils.PathSearch("is_pre_user", item, nil),
			"account_expires":            utils.PathSearch("account_expires", item, nil),
			"password_never_expired":     utils.PathSearch("password_never_expired", item, nil),
			"account_expired":            utils.PathSearch("account_expired", item, nil),
			"enable_change_password":     utils.PathSearch("enable_change_password", item, nil),
			"next_login_change_password": utils.PathSearch("next_login_change_password", item, nil),
			"description":                utils.PathSearch("description", item, nil),
			"locked":                     utils.PathSearch("locked", item, nil),
			"disabled":                   utils.PathSearch("disabled", item, nil),
			"share_space_subscription":   utils.PathSearch("share_space_subscription", item, nil),
			"share_space_desktops":       utils.PathSearch("share_space_desktops", item, nil),
			"group_names":                utils.PathSearch("group_names", item, make([]interface{}, 0)),
			"enterprise_project_id":      utils.PathSearch("enterprise_project_id", item, nil),
			"user_info_map":              utils.PathSearch("user_info_map", item, nil),
		}
		result = append(result, user)
	}

	return result
}

func dataSourceUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := listUsers(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace users: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenUsers(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
