package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Workspace POST /v2/{project_id}/groups
// @API Workspace POST /v2/{project_id}/groups/{group_id}/actions
// @API Workspace DELETE /v2/{project_id}/groups/{group_id}
// @API Workspace GET /v2/{project_id}/groups
// @API Workspace GET /v2/{project_id}/groups/{group_id}/users
// @API Workspace PUT /v2/{project_id}/groups/{group_id}
func ResourceUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupCreate,
		ReadContext:   resourceUserGroupRead,
		UpdateContext: resourceUserGroupUpdate,
		DeleteContext: resourceUserGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the user group is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user group.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the user group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the user group.",
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of user.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of user.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The email of user.",
						},
						"phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The phone of user.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of user.",
						},
						"total_desktops": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of desktops the user has.",
						},
					},
				},
				Description: "The user information under the user group.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the user group.",
			},
		},
	}
}

func resourceUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	createOpt := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
	}
	err = groups.Create(client, createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace user group: %s", err)
	}

	d.SetId(createOpt.Name)
	err = refreshUserGroupID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	users := d.Get("users").([]interface{})
	if len(users) > 0 {
		groupID := d.Id()
		err = doActionUserGroup(client, groupID, "ADD", users)
		if err != nil {
			return diag.Errorf("error adding users to Workspace user group (%s): %s", groupID, err)
		}
	}

	return resourceUserGroupRead(ctx, d, meta)
}

func refreshUserGroupID(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	userGroups, err := groups.List(client, groups.ListOpts{})
	if err != nil {
		return fmt.Errorf("error querying Workspace user groups: %s", err)
	}

	groupID := d.Id()
	for _, group := range userGroups {
		if group.Name == groupID {
			d.SetId(group.ID)
			return nil
		}
	}

	return fmt.Errorf("the Workspace user group (%s) does not exist", groupID)
}

func resourceUserGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	userGroups, err := groups.List(client, groups.ListOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace user group")
	}

	var mErr *multierror.Error
	for _, userGroup := range userGroups {
		groupID := userGroup.ID
		if groupID != d.Id() {
			continue
		}
		mErr = multierror.Append(nil,
			d.Set("region", region),
			d.Set("name", userGroup.Name),
			d.Set("type", userGroup.Type),
			d.Set("description", userGroup.Description),
			d.Set("created_at", userGroup.CreatedAt),
		)
		userList, err := groups.ListUser(client, groupID, groups.ListUserOpts{})
		if err != nil {
			return diag.Errorf("error reading users under Workspace user group (%s): %s", groupID, err)
		}
		if len(userList) > 0 {
			users := make([]map[string]interface{}, 0, len(userList))
			for _, user := range userList {
				users = append(users, map[string]interface{}{
					"id":             user.ID,
					"name":           user.Name,
					"email":          user.Email,
					"phone":          user.Phone,
					"description":    user.Description,
					"total_desktops": user.TotalDesktops,
				})
			}
			mErr = multierror.Append(mErr, d.Set("users", users))
		}

		return diag.FromErr(mErr.ErrorOrNil())
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Workspace user group")
}

func doActionUserGroup(client *golangsdk.ServiceClient, groupID, action string, users []interface{}) error {
	userIds := make([]string, len(users))
	for i, user := range users {
		userMap := user.(map[string]interface{})
		userIds[i] = userMap["id"].(string)
	}

	actionOpts := groups.ActionOpts{
		UserIDs: userIds,
		Type:    action,
	}
	err := groups.DoAction(client, groupID, actionOpts)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	groupID := d.Id()
	updateOpts := groups.UpdateOpts{
		GroupID: groupID,
	}
	if d.HasChanges("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChanges("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	err = groups.Update(client, updateOpts)
	if err != nil {
		return diag.Errorf("error updating Workspace user group (%s): %s", groupID, err)
	}

	if d.HasChanges("users") {
		oldRaws, newRaws := d.GetChange("users")
		newUsers := newRaws.([]interface{})
		oldUsers := oldRaws.([]interface{})
		if len(oldUsers) > 0 {
			err = doActionUserGroup(client, groupID, "DELETE", oldUsers)
			if err != nil {
				return diag.Errorf("error removing users from Workspace user group (%s): %s", groupID, err)
			}
		}

		if len(newUsers) > 0 {
			err = doActionUserGroup(client, groupID, "ADD", newUsers)
			if err != nil {
				return diag.Errorf("error adding users to Workspace user group (%s): %s", groupID, err)
			}
		}
	}

	return resourceUserGroupRead(ctx, d, meta)
}

func resourceUserGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	groupID := d.Id()
	err = groups.Delete(client, groupID)
	if err != nil {
		// WKS.00170117: The tenant does not exist.
		// WKS.00170208: The user group does not exist.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", []string{"WKS.00170117", "WKS.00170208"}...),
			fmt.Sprintf("error deleting Workspace user group (%s)", groupID))
	}

	return nil
}
