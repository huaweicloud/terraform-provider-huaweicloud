package iam

import (
	"context"
	"fmt"
	"log"
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

var v5GroupMembershipNonUpdatableParams = []string{"group_id"}

// @API IAM POST /v5/groups/{group_id}/add-user
// @API IAM GET /v5/users
// @API IAM POST /v5/groups/{group_id}/remove-user
func ResourceV5GroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5GroupMembershipCreate,
		ReadContext:   resourceV5GroupMembershipRead,
		UpdateContext: resourceV5GroupMembershipUpdate,
		DeleteContext: resourceV5GroupMembershipDelete,

		CustomizeDiff: config.FlexibleForceNew(v5GroupMembershipNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user group.`,
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The ID of the user.`,
								utils.SchemaDescInput{Required: true},
							),
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
							Description: `Whether the user is a root user.`,
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
					},
				},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
				Description: utils.SchemaDesc(
					`The list of users associated with the group.`,
					utils.SchemaDescInput{Required: true},
				),
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			// Internal attribute(s).
			"users_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'users'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			// Deprecated parameter(s).
			"user_id_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The list of user IDs to associate with the group.`,
					utils.SchemaDescInput{
						Required:   true,
						Deprecated: true,
					},
				),
			},
		},
	}
}

func resourceV5GroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var (
		groupId = d.Get("group_id").(string)
		userIds []interface{}
	)

	// Compatible with the deprecated 'user_id_list' parameter, if the users parameter is not set, use the 'user_id_list' parameter.
	users, isConfigScriptUsers := d.GetOk("users")
	if !isConfigScriptUsers {
		userIds = d.Get("user_id_list").(*schema.Set).List()
	} else {
		userIds = utils.PathSearch("[*].user_id", users, make([]interface{}, 0)).([]interface{})
	}

	if len(userIds) == 0 {
		return diag.Errorf("At least one user is required to be added to the group")
	}

	for _, userId := range userIds {
		if err := v5AddUsersToGroup(client, groupId, userId.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(groupId)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	if isConfigScriptUsers {
		if err = refreshV5GroupMembershipUsersOrigin(d); err != nil {
			// Don't report an error if origin refresh fails
			log.Printf("[WARN] Unable to refresh the users origin values: %s", err)
		}
	}

	return resourceV5GroupMembershipRead(ctx, d, meta)
}

func refreshV5GroupMembershipUsersOrigin(d *schema.ResourceData) error {
	scriptConfigValue := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "users")
	if scriptConfigValue == nil {
		log.Printf("[WARN] Unable to get the script configuration value for the users parameter")
		return nil
	}

	return d.Set("users_origin", utils.PathSearch("[*].user_id", scriptConfigValue, make([]interface{}, 0)).([]interface{}))
}

func GetV5GroupassociateUsers(client *golangsdk.ServiceClient, groupId string, usersOrigin []interface{}) ([]interface{}, error) {
	assocatedUsers, err := listV5Users(client, fmt.Sprintf("&group_id=%s", groupId))
	if err != nil {
		return nil, err
	}

	if len(assocatedUsers) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "v5/users",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the users in group (%s) does not exist", groupId)),
			},
		}
	}

	associatedUserIds := utils.PathSearch("[*].user_id", assocatedUsers, make([]interface{}, 0)).([]interface{})
	// When the user bound in the script is manually unbound in the console, and there are users bound through other ways under the group,
	// 404 error is returned.
	if len(assocatedUsers) == 0 || (len(usersOrigin) > 0 && len(utils.FildSliceIntersection(associatedUserIds, usersOrigin)) == 0) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "v5/users",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("All locally managed users (%v) have been unbound from the group (%s)", usersOrigin, groupId)),
			},
		}
	}

	return orderV5AssociatedUsersByUsersOrigin(assocatedUsers, usersOrigin), nil
}

func orderV5AssociatedUsersByUsersOrigin(associatedUsers, usersOrigin []interface{}) []interface{} {
	if len(usersOrigin) < 1 {
		return associatedUsers
	}

	sortedAssociatedUsers := make([]interface{}, 0, len(associatedUsers))
	associatedUsersCopy := associatedUsers
	for _, userIdOrigin := range usersOrigin {
		for index, user := range associatedUsersCopy {
			if utils.PathSearch("user_id", user, "").(string) == userIdOrigin {
				// Add the found user to the sorted users list.
				sortedAssociatedUsers = append(sortedAssociatedUsers, associatedUsersCopy[index])
				// Remove the processed user from the original array.
				associatedUsersCopy = append(associatedUsersCopy[:index], associatedUsersCopy[index+1:]...)

				break
			}
		}
	}
	// Add any remaining unsorted users to the end of the sorted list.
	sortedAssociatedUsers = append(sortedAssociatedUsers, associatedUsersCopy...)
	return sortedAssociatedUsers
}

func resourceV5GroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Id()
	sortedUsers, err := GetV5GroupassociateUsers(client, groupId, d.Get("users_origin").([]interface{}))
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting users in group (%s)", groupId))
	}

	mErr := multierror.Append(
		d.Set("group_id", groupId),
		d.Set("users", flattenV5GroupAssociatedUsers(sortedUsers)),
		// Deprecated parameter(s).
		d.Set("user_id_list", utils.PathSearch("[*].user_id", sortedUsers, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting group membership fields: %s", err)
	}

	return nil
}

func flattenV5GroupAssociatedUsers(users []interface{}) []interface{} {
	if len(users) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"user_id":      utils.PathSearch("user_id", user, nil),
			"user_name":    utils.PathSearch("user_name", user, nil),
			"enabled":      utils.PathSearch("enabled", user, nil),
			"description":  utils.PathSearch("description", user, nil),
			"is_root_user": utils.PathSearch("is_root_user", user, nil),
			"created_at":   utils.PathSearch("created_at", user, nil),
			"urn":          utils.PathSearch("urn", user, nil),
		})
	}
	return result
}

func resourceV5GroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var (
		groupId           = d.Id()
		configScriptUsers = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "users")
		oldUserIds        []interface{}
		newUserIds        []interface{}
	)

	// Compatible with the deprecated user_id_list parameter, if the users parameter is not set, use the 'user_id_list' parameter.
	if v, ok := configScriptUsers.([]interface{}); ok && len(v) > 0 {
		oldRaw, newRaw := d.GetChange("users")
		oldUserIds = utils.PathSearch("[*].user_id", oldRaw, make([]interface{}, 0)).([]interface{})
		newUserIds = utils.PathSearch("[*].user_id", newRaw, make([]interface{}, 0)).([]interface{})
	} else {
		oldRaw, newRaw := d.GetChange("user_id_list")
		oldUserIds = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set)).List()
		newUserIds = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set)).List()
	}

	if len(oldUserIds) > 0 {
		for _, userId := range oldUserIds {
			if !utils.SliceContains(newUserIds, userId) {
				if err := v5RemoveUsersFromGroup(client, groupId, userId.(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	for _, userId := range newUserIds {
		if !utils.SliceContains(oldUserIds, userId) {
			if err := v5AddUsersToGroup(client, groupId, userId.(string)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	if configScriptUsers != nil {
		if err = refreshV5GroupMembershipUsersOrigin(d); err != nil {
			// Don't report an error if origin refresh fails
			log.Printf("[WARN] Unable to refresh the users origin values: %s", err)
		}
	}

	return resourceV5GroupMembershipRead(ctx, d, meta)
}

func resourceV5GroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		groupId = d.Get("group_id").(string)
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var userIds []interface{}
	configScriptUsers := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "users")
	// Compatible with the deprecated 'user_id_list' parameter, if the users parameter is not set, use the 'user_id_list' parameter.
	if v, ok := configScriptUsers.([]interface{}); ok && len(v) > 0 {
		userIds = utils.PathSearch("[*].user_id", d.Get("users").([]interface{}), make([]interface{}, 0)).([]interface{})
		// The value of users_origin is empty only when the resource is imported and the terraform apply command is not executed.
		// In this case, all information obtained from the remote service is used to remove user relationships from the group.
		if userIdsOrigin, ok := d.GetOk("users_origin"); ok && len(userIdsOrigin.([]interface{})) > 0 {
			log.Printf("[DEBUG] Find the custom users configuration, according to it to remove users from the group (%v)", groupId)
			userIds = userIdsOrigin.([]interface{})
		}
	} else {
		userIds = d.Get("user_id_list").(*schema.Set).List()
	}

	for _, userId := range userIds {
		if err := v5RemoveUsersFromGroup(client, groupId, userId.(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func v5AddUsersToGroup(client *golangsdk.ServiceClient, groupId, userId string) error {
	httpUrl := "v5/groups/{group_id}/add-user"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{group_id}", groupId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"user_id": userId,
		},
	}
	_, err := client.Request("POST", addPath, &opt)
	if err != nil {
		return fmt.Errorf("error adding user (%s) to group (%s): %s ", userId, groupId, err)
	}

	return nil
}

func v5RemoveUsersFromGroup(client *golangsdk.ServiceClient, groupId, userId string) error {
	httpUrl := "v5/groups/{group_id}/remove-user"
	removePath := client.Endpoint + httpUrl
	removePath = strings.ReplaceAll(removePath, "{group_id}", groupId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"user_id": userId,
		},
	}
	_, err := client.Request("POST", removePath, &opt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] the user (%s) is not exist, ignore to remove it from the group", userId)
		}

		return fmt.Errorf("error removing user (%s) from group (%s): %s ", userId, groupId, err)
	}

	return nil
}

func listV5Users(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		limit   = 200
		httpUrl = "v5/users"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = fmt.Sprintf("%s?limit=%v", listPath, limit)
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	listOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		users := utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}
