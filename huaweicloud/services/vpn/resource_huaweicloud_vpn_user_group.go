package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var userGroupNonUpdatableParams = []string{"vpn_server_id"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}
// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/add-users
// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/remove-users
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/users
func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		UpdateContext: resourceUserGroupUpdate,
		ReadContext:   resourceUserGroupRead,
		DeleteContext: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserGroupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(userGroupNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The VPN server ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the user group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the user group.`,
			},
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        userSchema(),
				Description: `The user list.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user group type.`,
			},
			"user_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func userSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The username.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user description.`,
			},
		},
	}
	return &sc
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createUserGroupHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups"
		createUserGroupProduct = "vpn"
	)
	createUserGroupClient, err := conf.NewServiceClient(createUserGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	vpnServerId := d.Get("vpn_server_id").(string)
	createUserGroupPath := createUserGroupClient.Endpoint + createUserGroupHttpUrl
	createUserGroupPath = strings.ReplaceAll(createUserGroupPath, "{project_id}", createUserGroupClient.ProjectID)
	createUserGroupPath = strings.ReplaceAll(createUserGroupPath, "{vpn_server_id}", vpnServerId)

	createUserGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createUserGroupOpt.JSONBody = utils.RemoveNil(buildCreateUserGroupBodyParams(d))
	createUserGroupResp, err := createUserGroupClient.Request("POST", createUserGroupPath, &createUserGroupOpt)
	if err != nil {
		return diag.Errorf("error creating VPN user group: %s", err)
	}

	createUserGroupRespBody, err := utils.FlattenResponse(createUserGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("user_group.id", createUserGroupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN user group: ID is not found in API response")
	}
	d.SetId(id)

	users := d.Get("users").(*schema.Set).List()
	if len(users) > 0 {
		err = addUserToUserGroup(createUserGroupClient, vpnServerId, id, users)
		if err != nil {
			return diag.Errorf("error adding users to VPN user group: %s", err)
		}
	}

	// The creation interface is asynchronous.
	// If the user group information disappears, then the creation fails.
	// Wait for a while to check if the creation is successful.
	// lintignore:R018
	time.Sleep(30 * time.Second)

	return resourceUserGroupRead(ctx, d, meta)
}

func buildCreateUserGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_group": map[string]interface{}{
			"name":        d.Get("name"),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceUserGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	getUserGroupProduct := "vpn"
	getUserGroupClient, err := conf.NewServiceClient(getUserGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	getUserGroupRespBody, err := GetUserGroup(getUserGroupClient, serverId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN user group")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("user_group.name", getUserGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("user_group.description", getUserGroupRespBody, nil)),
		d.Set("type", utils.PathSearch("user_group.type", getUserGroupRespBody, nil)),
		d.Set("user_number", utils.PathSearch("user_group.user_number", getUserGroupRespBody, nil)),
		d.Set("created_at", utils.PathSearch("user_group.created_at", getUserGroupRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("user_group.updated_at", getUserGroupRespBody, nil)),
	)

	users, err := getUsersFromUserGroup(getUserGroupClient, serverId, id)
	if err != nil {
		return diag.Errorf("error retrieving users in the VPN user group: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("users", flattenUsers(users)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetUserGroup(client *golangsdk.ServiceClient, serverId, id string) (interface{}, error) {
	getUserGroupHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}"
	getUserGroupPath := buildUserGroupURL(client, getUserGroupHttpUrl, serverId, id)

	getUserGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getUserGroupResp, err := client.Request("GET", getUserGroupPath, &getUserGroupOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getUserGroupResp)
}

func getUsersFromUserGroup(client *golangsdk.ServiceClient, serverId, id string) (interface{}, error) {
	getUsersUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/users"
	baseUsersPath := buildUserGroupURL(client, getUsersUrl, serverId, id)
	baseUsersPath += "?limit=50"
	getUsersPath := baseUsersPath
	getUsersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	rst := make([]interface{}, 0)

	for {
		getUsersResp, err := client.Request("GET", getUsersPath, &getUsersOpt)
		if err != nil {
			return nil, err
		}

		getUsersRespBody, err := utils.FlattenResponse(getUsersResp)
		if err != nil {
			return nil, err
		}

		users := utils.PathSearch("users", getUsersRespBody, make([]interface{}, 0)).([]interface{})
		if len(users) > 0 {
			rst = append(rst, users...)
		}

		marker := utils.PathSearch("page_info.next_marker", getUsersRespBody, nil)
		if marker == nil {
			break
		}
		getUsersPath = fmt.Sprintf("%s&marker=%v", baseUsersPath, marker.(string))
	}

	return rst, nil
}

func flattenUsers(resp interface{}) []interface{} {
	if users, ok := resp.([]interface{}); ok {
		rst := make([]interface{}, len(users))

		for i, user := range users {
			userMap := user.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"id":          utils.PathSearch("id", userMap, nil),
				"name":        utils.PathSearch("name", userMap, nil),
				"description": utils.PathSearch("description", userMap, nil),
			}
		}
		return rst
	}
	return nil
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	updateUserGroupClient, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	updateUserGrouphasChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateUserGrouphasChanges...) {
		updateUserGroupHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}"
		updateUserGroupPath := buildUserGroupURL(updateUserGroupClient, updateUserGroupHttpUrl, serverId, id)

		updateUserGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateUserGroupOpt.JSONBody = buildUpdateUserGroupBodyParams(d)
		_, err = updateUserGroupClient.Request("PUT", updateUserGroupPath, &updateUserGroupOpt)
		if err != nil {
			return diag.Errorf("error updating VPN user group: %s", err)
		}
	}

	if d.HasChange("users") {
		oldUsers, newUsers := d.GetChange("users")
		rmUsers := oldUsers.(*schema.Set).Difference(newUsers.(*schema.Set))
		addUsers := newUsers.(*schema.Set).Difference(oldUsers.(*schema.Set))

		if rmUsers.Len() > 0 {
			err = removeUserFromUserGroup(updateUserGroupClient, serverId, id, rmUsers.List())
			if err != nil {
				return diag.Errorf("error removing users from the VPN user group: %s", err)
			}
		}

		if addUsers.Len() > 0 {
			err = addUserToUserGroup(updateUserGroupClient, serverId, id, addUsers.List())
			if err != nil {
				return diag.Errorf("error adding users to the VPN user group: %s", err)
			}
		}
	}

	return resourceUserGroupRead(ctx, d, meta)
}

func buildUserGroupURL(client *golangsdk.ServiceClient, urlTemplate, serverId, id string) string {
	url := client.Endpoint + urlTemplate
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{vpn_server_id}", serverId)
	url = strings.ReplaceAll(url, "{group_id}", id)
	return url
}

func buildUpdateUserGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_group": map[string]interface{}{
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}
	return bodyParams
}

func removeUserFromUserGroup(client *golangsdk.ServiceClient, serverId, id string, users []interface{}) error {
	removeUserHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/remove-users"
	removeUserPath := buildUserGroupURL(client, removeUserHttpUrl, serverId, id)

	removeUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	removeUserOpt.JSONBody = buildUserBodyParams(users)
	_, err := client.Request("POST", removeUserPath, &removeUserOpt)
	return err
}

func addUserToUserGroup(client *golangsdk.ServiceClient, serverId, id string, users []interface{}) error {
	addUserHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}/add-users"
	addUserPath := buildUserGroupURL(client, addUserHttpUrl, serverId, id)

	addUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	addUserOpt.JSONBody = buildUserBodyParams(users)
	_, err := client.Request("POST", addUserPath, &addUserOpt)
	return err
}

func buildUserBodyParams(users []interface{}) map[string]interface{} {
	rst := make([]interface{}, len(users))

	for i, user := range users {
		userMap := user.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"id": userMap["id"],
		}
	}
	return map[string]interface{}{
		"users": rst,
	}
}

func resourceUserGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		deleteUserGroupHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/groups/{group_id}"
		deleteUserGroupProduct = "vpn"
	)
	deleteUserGroupClient, err := conf.NewServiceClient(deleteUserGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	deleteUserGroupPath := buildUserGroupURL(deleteUserGroupClient, deleteUserGroupHttpUrl, serverId, id)

	deleteUserGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteUserGroupClient.Request("DELETE", deleteUserGroupPath, &deleteUserGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN user group")
	}

	return nil
}

func resourceUserGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <vpn_server_id>/<id>")
	}

	d.Set("vpn_server_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
