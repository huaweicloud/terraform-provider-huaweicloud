package swr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	iamusers "github.com/chnsz/golangsdk/openstack/identity/v3.0/users"
	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SWR GET /v2/manage/namespaces/{namespace}/access
// @API SWR POST /v2/manage/namespaces/{namespace}/access
// @API SWR DELETE /v2/manage/namespaces/{namespace}/access
// @API IAM GET /v3.0/OS-USER/users/{userID}
func ResourceSWROrganizationPermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSWROrganizationPermissionsCreate,
		ReadContext:   resourceSWROrganizationPermissionsRead,
		UpdateContext: resourceSWROrganizationPermissionsUpdate,
		DeleteContext: resourceSWROrganizationPermissionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Manage", "Write", "Read",
							}, false),
						},
					},
				},
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_permission": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSWROrganizationPermissionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create SWR client : %s", err)
	}

	nameSpace := d.Get("organization").(string)

	userRaw := d.Get("users").([]interface{})

	err = addSwrOrganizationPermissions(d, cfg, swrClient, userRaw, nameSpace)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(nameSpace)

	return resourceSWROrganizationPermissionsRead(ctx, d, meta)
}

func resourceSWROrganizationPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	access, err := namespaces.GetAccess(swrClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR")
	}

	var users []map[string]interface{}
	for _, pairObject := range access.OthersAuths {
		if pairObject.UserName == access.SelfAuth.UserName {
			continue
		}
		user := make(map[string]interface{})
		user["user_name"] = pairObject.UserName
		user["user_id"] = pairObject.UserID

		permission := resourceSWRAuthToPermission(pairObject.Auth)
		user["permission"] = permission

		users = append(users, user)
	}

	selfPermission := []map[string]interface{}{
		{
			"user_name": access.SelfAuth.UserName,
			"user_id":   access.SelfAuth.UserID,
		},
	}

	permission := resourceSWRAuthToPermission(access.SelfAuth.Auth)
	selfPermission[0]["permission"] = permission

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("organization", access.Name),
		d.Set("creator", access.CreatorName),
		d.Set("self_permission", selfPermission),
		d.Set("users", users),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceSWROrganizationPermissionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	nameSpace := d.Get("organization").(string)

	users, _ := d.GetChange("users")

	userIDs := make([]string, 0, len(d.Get("users").([]interface{})))

	for _, userRaw := range users.([]interface{}) {
		user := userRaw.(map[string]interface{})
		userIDs = append(userIDs, user["user_id"].(string))
	}

	err = namespaces.DeleteAccess(swrClient, userIDs, nameSpace).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting SWR organization permission: %s", err)
	}
	userRaw := d.Get("users").([]interface{})
	addUsers := dealAddUsers(userRaw)
	err = addSwrOrganizationPermissions(d, cfg, swrClient, addUsers, nameSpace)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSWROrganizationPermissionsRead(ctx, d, meta)
}

func addSwrOrganizationPermissions(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient,
	addUsers []interface{}, nameSpace string) error {
	users := make([]namespaces.User, len(addUsers))
	for i, raw := range addUsers {
		rawMap := raw.(map[string]interface{})
		auth := resourceSWRPermissionToAuth(rawMap["permission"].(string))
		users[i] = namespaces.User{
			UserID: rawMap["user_id"].(string),
			Auth:   auth,
		}

		if rawMap["user_name"].(string) != "" {
			users[i].UserName = rawMap["user_name"].(string)
		} else {
			iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
			if err != nil {
				return fmt.Errorf("error creating iam client: %s", err)
			}

			user, err := iamusers.Get(iamClient, rawMap["user_id"].(string)).Extract()
			if err != nil {
				return fmt.Errorf("error retrieving user(%s): %s", rawMap["user_id"].(string), err)
			}
			log.Printf("[DEBUG] retrieved user: %#v", user)

			users[i].UserName = user.Name
		}
	}
	createOpts := namespaces.CreateAccessOpts{
		Users: users,
	}

	err := namespaces.CreateAccess(client, createOpts, nameSpace).ExtractErr()

	if err != nil {
		return fmt.Errorf("error creating SWR organization permission: %s", err)
	}
	return nil
}

func resourceSWROrganizationPermissionsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	nameSpace := d.Get("organization").(string)

	userIDs := make([]string, 0, len(d.Get("users").([]interface{})))

	for _, userRaw := range d.Get("users").([]interface{}) {
		user := userRaw.(map[string]interface{})
		userIDs = append(userIDs, user["user_id"].(string))
	}

	err = namespaces.DeleteAccess(swrClient, userIDs, nameSpace).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR organization permission")
	}

	return nil
}

func resourceSWRPermissionToAuth(permission string) int {
	auth := 0
	switch permission {
	case "Manage":
		auth = 7
	case "Write":
		auth = 3
	case "Read":
		auth = 1
	}

	return auth
}

func resourceSWRAuthToPermission(auth int) string {
	permission := "Unknown"
	switch auth {
	case 7:
		permission = "Manage"
	case 3:
		permission = "Write"
	case 1:
		permission = "Read"
	}

	return permission
}
