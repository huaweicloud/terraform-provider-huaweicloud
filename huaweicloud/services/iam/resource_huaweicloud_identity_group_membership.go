package iam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM PUT /v3/groups/{group_id}/users/{user_id}
// @API IAM DELETE /v3/groups/{group_id}/users/{user_id}
// @API IAM GET /v3/groups/{group_id}/users
// @API IAM GET /v3/groups/{group_id}
func ResourceIdentityGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupMembershipCreate,
		ReadContext:   resourceIdentityGroupMembershipRead,
		UpdateContext: resourceIdentityGroupMembershipUpdate,
		DeleteContext: resourceIdentityGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIdentityGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group").(string)
	userList := utils.ExpandToStringList(d.Get("users").(*schema.Set).List())

	if err := addUsersToGroup(identityClient, groupID, userList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)
	return resourceIdentityGroupMembershipRead(ctx, d, meta)
}

func resourceIdentityGroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group").(string)
	userList := d.Get("users").(*schema.Set)

	// check whether the user group exists
	_, err = groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "unable to query group")
	}

	allPages, err := users.ListInGroup(identityClient, groupID, nil).AllPages()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "unable to query group membership")
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return diag.Errorf("unable to retrieve users: %s", err)
	}

	var ul []string
	for _, u := range allUsers {
		if userList.Contains(u.ID) {
			ul = append(ul, u.ID)
		}
	}

	return diag.FromErr(d.Set("users", ul))
}

func resourceIdentityGroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group").(string)
	if d.HasChange("users") {
		oldRaw, newRaw := d.GetChange("users")
		rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))

		removeList := utils.ExpandToStringListBySet(rmSet)
		if err := removeUsersFromGroup(identityClient, groupID, removeList); err != nil {
			return diag.Errorf("error updating membership: %s", err)
		}

		addList := utils.ExpandToStringListBySet(addSet)
		if err := addUsersToGroup(identityClient, groupID, addList); err != nil {
			return diag.Errorf("error updating membership: %s", err)
		}
	}

	return resourceIdentityGroupMembershipRead(ctx, d, meta)
}

func resourceIdentityGroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group").(string)
	allUsers := utils.ExpandToStringList(d.Get("users").(*schema.Set).List())

	if err := removeUsersFromGroup(identityClient, groupID, allUsers); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func addUsersToGroup(identityClient *golangsdk.ServiceClient, groupID string, userList []string) error {
	for _, u := range userList {
		if r := users.AddToGroup(identityClient, groupID, u).ExtractErr(); r != nil {
			return fmt.Errorf("error adding user %s to group %s: %s ", u, groupID, r)
		}
	}
	return nil
}

func removeUsersFromGroup(identityClient *golangsdk.ServiceClient, groupID string, userList []string) error {
	for _, u := range userList {
		if r := users.RemoveFromGroup(identityClient, groupID, u).ExtractErr(); r != nil {
			if _, ok := r.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the user %s is not exist, ignore to remove it from the group", u)
				continue
			}
			return fmt.Errorf("error removing user %s from group %s: %s", u, groupID, r)
		}
	}
	return nil
}
