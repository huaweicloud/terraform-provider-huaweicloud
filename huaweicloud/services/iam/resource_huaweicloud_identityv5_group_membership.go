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
			"user_id_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of user IDs to associate with the group.`,
			},
			// Attributes.
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
				Description: `The list of users associated with the group.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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

	groupId := d.Get("group_id").(string)
	users := utils.ExpandToStringList(d.Get("user_id_list").(*schema.Set).List())
	if err := v5AddUsersToGroup(client, groupId, users); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupId)
	return resourceV5GroupMembershipRead(ctx, d, meta)
}

func GetV5GroupassociateUsers(client *golangsdk.ServiceClient, groupId string) ([]interface{}, error) {
	users, err := listV5Users(client, fmt.Sprintf("&group_id=%s", groupId))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "v5/users",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the users in group (%s) does not exist", groupId)),
			},
		}
	}

	return users, nil
}

func resourceV5GroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	users, err := GetV5GroupassociateUsers(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting users in group (%s)", d.Id()))
	}

	mErr := multierror.Append(
		d.Set("group_id", d.Id()),
		d.Set("user_id_list", utils.PathSearch("[*].user_id", users, nil)),
		d.Set("users", flattenV5GroupAssociatedUsers(users)),
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

	groupId := d.Get("group_id").(string)
	if d.HasChange("user_id_list") {
		oldRaw, newRaw := d.GetChange("user_id_list")
		rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))

		removeList := utils.ExpandToStringListBySet(rmSet)
		if err := v5RemoveUsersFromGroup(client, groupId, removeList); err != nil {
			return diag.FromErr(err)
		}

		addList := utils.ExpandToStringListBySet(addSet)
		if err := v5AddUsersToGroup(client, groupId, addList); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV5GroupMembershipRead(ctx, d, meta)
}

func resourceV5GroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	allUsers := utils.ExpandToStringList(d.Get("user_id_list").(*schema.Set).List())
	if err := v5RemoveUsersFromGroup(client, groupId, allUsers); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func v5AddUsersToGroup(client *golangsdk.ServiceClient, groupId string, users []string) error {
	httpUrl := "v5/groups/{group_id}/add-user"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{group_id}", groupId)
	for _, user := range users {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"user_id": user,
			},
		}
		_, err := client.Request("POST", addPath, &opt)
		if err != nil {
			return fmt.Errorf("error adding user (%s) to group (%s): %s ", user, groupId, err)
		}
	}

	return nil
}

func v5RemoveUsersFromGroup(client *golangsdk.ServiceClient, groupId string, users []string) error {
	httpUrl := "v5/groups/{group_id}/remove-user"
	removePath := client.Endpoint + httpUrl
	removePath = strings.ReplaceAll(removePath, "{group_id}", groupId)
	for _, user := range users {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"user_id": user,
			},
		}
		_, err := client.Request("POST", removePath, &opt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the user (%s) is not exist, ignore to remove it from the group", user)
				continue
			}

			return fmt.Errorf("error removing user (%s) from group (%s): %s ", user, groupId, err)
		}
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
