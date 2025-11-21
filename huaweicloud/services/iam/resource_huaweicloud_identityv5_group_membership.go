package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/groups/{group_id}/remove-user
// @API IAM POST /v5/groups/{group_id}/add-user
// @API IAM GET /v5/users
var groupMembershipNonUpdatableParams = []string{"group_id"}

func ResourceIdentityV5GroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5GroupMembershipCreate,
		ReadContext:   resourceIdentityV5GroupMembershipRead,
		UpdateContext: resourceIdentityV5GroupMembershipUpdate,
		DeleteContext: resourceIdentityV5GroupMembershipDelete,

		CustomizeDiff: config.FlexibleForceNew(groupMembershipNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"user_id_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_root_user": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIdentityV5GroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	userList := utils.ExpandToStringList(d.Get("user_id_list").(*schema.Set).List())

	if err := addUsersToGroupV5(iamClient, groupID, userList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)
	return resourceIdentityV5GroupMembershipRead(ctx, d, meta)
}

func resourceIdentityV5GroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getGroupHttpUrl := "v5/users"
	getGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	var allUsers []interface{}
	var marker string
	var path string
	for {
		path = iamClient.Endpoint + getGroupHttpUrl + buildQueryGroupParamPath(d.Id(), marker)
		getGroupResp, err := iamClient.Request("GET", path, &getGroupOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error getting IAM group")
		}
		getGroupRespBody, err := utils.FlattenResponse(getGroupResp)
		if err != nil {
			return diag.FromErr(err)
		}
		users := utils.PathSearch("users", getGroupRespBody, make([]interface{}, 0)).([]interface{})
		allUsers = append(allUsers, users...)

		marker = utils.PathSearch("page_info.next_marker", getGroupRespBody, "").(string)
		if marker == "" {
			break
		}
	}
	if err := d.Set("users", allUsers); err != nil {
		return diag.Errorf("error setting users fields: %s", err)
	}

	return nil
}

func resourceIdentityV5GroupMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	if d.HasChange("user_id_list") {
		oldRaw, newRaw := d.GetChange("user_id_list")
		rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))

		removeList := utils.ExpandToStringListBySet(rmSet)
		if err := removeUsersFromGroupV5(iamClient, groupID, removeList); err != nil {
			return diag.Errorf("error updating membership: %s", err)
		}

		addList := utils.ExpandToStringListBySet(addSet)
		if err := addUsersToGroupV5(iamClient, groupID, addList); err != nil {
			return diag.Errorf("error updating membership: %s", err)
		}
	}
	return resourceIdentityV5GroupMembershipRead(ctx, d, meta)
}

func resourceIdentityV5GroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	groupID := d.Get("group_id").(string)
	allUsers := utils.ExpandToStringList(d.Get("user_id_list").(*schema.Set).List())

	if err := removeUsersFromGroupV5(iamClient, groupID, allUsers); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func addUsersToGroupV5(iamClient *golangsdk.ServiceClient, groupID string, userList []string) error {
	addGroupMembershipHttpUrl := "v5/groups/{group_id}/add-user"
	addGroupMembershipPath := iamClient.Endpoint + addGroupMembershipHttpUrl
	addGroupMembershipPath = strings.ReplaceAll(addGroupMembershipPath, "{group_id}", groupID)
	for _, u := range userList {
		addGroupMembershipOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"user_id": u,
			},
		}
		addGroupMembershipResp, err := iamClient.Request("POST", addGroupMembershipPath, &addGroupMembershipOpt)
		if err != nil {
			return fmt.Errorf("error adding user %s to group %s: %v ", u, groupID, addGroupMembershipResp)
		}
	}
	return nil
}

func removeUsersFromGroupV5(iamClient *golangsdk.ServiceClient, groupID string, userList []string) error {
	removeGroupMembershipHttpUrl := "v5/groups/{group_id}/remove-user"
	removeGroupMembershipPath := iamClient.Endpoint + removeGroupMembershipHttpUrl
	removeGroupMembershipPath = strings.ReplaceAll(removeGroupMembershipPath, "{group_id}", groupID)
	for _, u := range userList {
		removeGroupMembershipOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"user_id": u,
			},
		}
		removeGroupMembershipResp, err := iamClient.Request("POST", removeGroupMembershipPath, &removeGroupMembershipOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the user %s is not exist, ignore to remove it from the group", u)
				continue
			}
			return fmt.Errorf("error removing user %s from group %s: %v ", u, groupID, removeGroupMembershipResp)
		}
	}
	return nil
}

func buildQueryGroupParamPath(groupId string, marker string) string {
	res := "?limit=100"
	if groupId != "" {
		res = fmt.Sprintf("%s&group_id=%s", res, groupId)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%s", res, marker)
	}
	return res
}
