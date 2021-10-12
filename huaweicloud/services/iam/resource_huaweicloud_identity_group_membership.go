package iam

import (
	"context"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityGroupMembershipV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupMembershipV3Create,
		ReadContext:   resourceIdentityGroupMembershipV3Read,
		UpdateContext: resourceIdentityGroupMembershipV3Update,
		DeleteContext: resourceIdentityGroupMembershipV3Delete,

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

func resourceIdentityGroupMembershipV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	group := d.Get("group").(string)
	userList := utils.ExpandToStringList(d.Get("users").(*schema.Set).List())

	if err := addUsersToGroup(identityClient, group, userList); err != nil {
		return fmtp.DiagErrorf("Error adding users to identity group: %s", err)
	}

	//lintignore:R015
	d.SetId(resource.UniqueId())

	return resourceIdentityGroupMembershipV3Read(ctx, d, meta)
}

func resourceIdentityGroupMembershipV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}
	group := d.Get("group").(string)
	userList := d.Get("users").(*schema.Set)
	var ul []string

	allPages, err := users.ListInGroup(identityClient, group, users.ListOpts{}).AllPages()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Unable to query groups")
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve users: %s", err)
	}

	for _, u := range allUsers {
		if userList.Contains(u.ID) {
			ul = append(ul, u.ID)
		}
	}

	if err := d.Set("users", ul); err != nil {
		return fmtp.DiagErrorf("Error setting user list from IAM (%s), error: %s", group, err)
	}

	return nil
}

func resourceIdentityGroupMembershipV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	if d.HasChange("users") {
		group := d.Get("group").(string)

		o, n := d.GetChange("users")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := utils.ExpandToStringList(os.Difference(ns).List())
		add := utils.ExpandToStringList(ns.Difference(os).List())

		if err := removeUsersFromGroup(identityClient, group, remove); err != nil {
			return fmtp.DiagErrorf("Error update user-group-membership: %s", err)
		}

		if err := addUsersToGroup(identityClient, group, add); err != nil {
			return fmtp.DiagErrorf("Error update user-group-membership: %s", err)
		}
	}

	return resourceIdentityGroupMembershipV3Read(ctx, d, meta)
}

func resourceIdentityGroupMembershipV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	group := d.Get("group").(string)
	users := utils.ExpandToStringList(d.Get("users").(*schema.Set).List())

	if err := removeUsersFromGroup(identityClient, group, users); err != nil {
		return fmtp.DiagErrorf("Error delete user-group-membership: %s", err)
	}

	d.SetId("")
	return nil
}

func addUsersToGroup(identityClient *golangsdk.ServiceClient, group string, userList []string) error {
	for _, u := range userList {
		if r := users.AddToGroup(identityClient, group, u).ExtractErr(); r != nil {
			return fmtp.Errorf("Error add user %s to group %s: %s ", u, group, r)
		}
	}
	return nil
}

func removeUsersFromGroup(identityClient *golangsdk.ServiceClient, group string, userList []string) error {
	for _, u := range userList {
		if r := users.RemoveFromGroup(identityClient, group, u).ExtractErr(); r != nil {
			if _, ok := r.(golangsdk.ErrDefault404); ok {
				logp.Printf("[WARN] the user %s is not exist, ignore to remove it from the group", u)
				continue
			}
			return fmtp.Errorf("Error remove user %s from group %s: %s", u, group, r)
		}
	}
	return nil
}
