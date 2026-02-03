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

var (
	v3GroupMembershipNonUpdatableParams = []string{"group"}
	strSliceParamKeysForGroupMembership = []string{"users"}
)

// @API IAM PUT /v3/groups/{group_id}/users/{user_id}
// @API IAM DELETE /v3/groups/{group_id}/users/{user_id}
// @API IAM GET /v3/groups/{group_id}/users
// @API IAM GET /v3/groups/{group_id}
func ResourceV3GroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3GroupMembershipCreate,
		ReadContext:   resourceV3GroupMembershipRead,
		UpdateContext: resourceV3GroupMembershipUpdate,
		DeleteContext: resourceV3GroupMembershipDelete,

		CustomizeDiff: config.FlexibleForceNew(v3GroupMembershipNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the group to which the users belong.`,
			},
			"users": {
				Type:             schema.TypeList,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
				Description:      `The list of user IDs associated with the group.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal attributes.
			"users_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'users'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func addV3UserToGroup(client *golangsdk.ServiceClient, groupId, userId string) error {
	httpUrl := "v3/groups/{group_id}/users/{user_id}"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{group_id}", groupId)
	addPath = strings.ReplaceAll(addPath, "{user_id}", userId)

	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err := client.Request("PUT", addPath, &addOpt)
	return err
}

func resourceV3GroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Get("group").(string)
		userIds = d.Get("users").([]interface{})
		mErr    *multierror.Error
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	d.SetId(groupId)

	log.Printf("[DEBUG] Prepare to add users to group (%v): %v", groupId, userIds)
	for _, userId := range userIds {
		if err := addV3UserToGroup(client, groupId, userId.(string)); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error adding user (%s) to group (%s): %s", userId, groupId, err))
		}
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error adding users to group (%s): %s", groupId, err)
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForGroupMembership)
	if err != nil {
		// Don't report an error if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}
	return resourceV3GroupMembershipRead(ctx, d, meta)
}

func ListV3AssociatedUsersForGroup(client *golangsdk.ServiceClient, groupId string, usersOrigin []interface{}) ([]interface{}, error) {
	// Before listing associated users, check whether the user group exists.
	_, err := GetV3GroupById(client, groupId)
	if err != nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/groups/{group_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the group (%s) does not exist", groupId)),
			},
		}
	}

	httpUrl := "v3/groups/{group_id}/users"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{group_id}", groupId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	users := orderV3AssociatedUsersByUsersOrigin(utils.PathSearch("users[*].id", respBody, make([]interface{}, 0)).([]interface{}), usersOrigin)
	if len(users) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/groups/{group_id}/users",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the managed membership for group (%s) does not exist", groupId)),
			},
		}
	}
	return users, nil
}

func orderV3AssociatedUsersByUsersOrigin(associatedUsers, usersOrigin []interface{}) []interface{} {
	if len(usersOrigin) < 1 {
		return associatedUsers
	}

	sortedAssociatedUsers := make([]interface{}, 0, len(associatedUsers))
	associatedUsersCopy := associatedUsers
	for _, userIdOrigin := range usersOrigin {
		for index, userId := range associatedUsersCopy {
			if userId != userIdOrigin {
				continue
			}
			// Add the found user to the sorted users list.
			sortedAssociatedUsers = append(sortedAssociatedUsers, associatedUsersCopy[index])
			// Remove the processed user from the original array.
			associatedUsersCopy = append(associatedUsersCopy[:index], associatedUsersCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted users to the end of the sorted list.
	sortedAssociatedUsers = append(sortedAssociatedUsers, associatedUsersCopy...)
	return sortedAssociatedUsers
}

func resourceV3GroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		groupId     = d.Id()
		usersOrigin = d.Get("users_origin").([]interface{})
		mErr        *multierror.Error
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	associatedUsers, err := ListV3AssociatedUsersForGroup(client, groupId, usersOrigin)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "unable to query group membership")
	}

	mErr = multierror.Append(mErr,
		d.Set("group", groupId),
		d.Set("users", associatedUsers),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting group membership fields: %s", err)
	}
	return nil
}

func removeV3UserFromGroup(client *golangsdk.ServiceClient, groupId, userId string) error {
	httpUrl := "v3/groups/{group_id}/users/{user_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", groupId)
	deletePath = strings.ReplaceAll(deletePath, "{user_id}", userId)

	removeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err := client.Request("DELETE", deletePath, &removeOpt)
	return err
}

func resourceV3GroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	oldRaw, newRaw := d.GetChange("users")
	oldUserIds := oldRaw.([]interface{})
	newUserIds := newRaw.([]interface{})

	for _, userId := range oldUserIds {
		if !utils.SliceContains(newUserIds, userId) {
			log.Printf("[DEBUG] Prepare to remove user (%s) from group (%v)", userId, groupId)
			if err := removeV3UserFromGroup(client, groupId, userId.(string)); err != nil {
				return diag.Errorf("error removing user (%s) from group (%v): %s", userId, groupId, err)
			}
		}
	}

	for _, userId := range newUserIds {
		if !utils.SliceContains(oldUserIds, userId) {
			log.Printf("[DEBUG] Prepare to add user (%s) to group (%v)", userId, groupId)
			if err := addV3UserToGroup(client, groupId, userId.(string)); err != nil {
				return diag.Errorf("error adding user (%s) to group (%v): %s", userId, groupId, err)
			}
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForGroupMembership)
	if err != nil {
		// Don't report an error if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceV3GroupMembershipRead(ctx, d, meta)
}

func resourceV3GroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userIds := d.Get("users").([]interface{})
	// The value of users_origin is empty only when the resource is imported and the terraform apply command is not executed.
	// In this case, all information obtained from the remote service is used to remove user relationships from the group.
	if userIdsOrigin, ok := d.GetOk("users_origin"); ok && len(userIdsOrigin.([]interface{})) > 0 {
		log.Printf("[DEBUG] Find the custom users configuration, according to it to remove users from the group (%v)", groupId)
		userIds = userIdsOrigin.([]interface{})
	}

	log.Printf("[DEBUG] Prepare to remove all users from group (%v): %v", groupId, userIds)
	for _, userId := range userIds {
		if err := removeV3UserFromGroup(client, groupId, userId.(string)); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the user (%s) is not exist, ignore to remove it from the group", userId)
				continue
			}
			return diag.Errorf("error removing user (%s) from group (%v): %s", userId, groupId, err)
		}
	}

	return nil
}
