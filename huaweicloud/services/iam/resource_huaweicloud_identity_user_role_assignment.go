package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/eps_permissions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM PUT /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/users/{user_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/users/{user_id}/roles/{role_id}
// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/users/{user_id}/roles
func ResourceIdentityUserRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserRoleAssignmentCreate,
		ReadContext:   resourceIdentityUserRoleAssignmentRead,
		DeleteContext: resourceIdentityUserRoleAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityUserRoleAssignmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityUserRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	userID := d.Get("user_id").(string)
	roleID := d.Get("role_id").(string)
	enterpriseProjectID := d.Get("enterprise_project_id").(string)

	err = eps_permissions.UserPermissionsCreate(client, enterpriseProjectID, userID, roleID).ExtractErr()
	if err != nil {
		return diag.Errorf("error assigning role (%s) to enterprise project (%s) and user (%s): %s",
			roleID, enterpriseProjectID, userID, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", enterpriseProjectID, roleID, userID))

	return resourceIdentityUserRoleAssignmentRead(ctx, d, meta)
}

func resourceIdentityUserRoleAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	userID := d.Get("user_id").(string)
	roleID := d.Get("role_id").(string)
	enterpriseProjectID := d.Get("enterprise_project_id").(string)
	role, err := GetUserRoleAssignmentWithEpsID(client, userID, roleID, enterpriseProjectID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting role assignment")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", userID, roleID, enterpriseProjectID))

	mErr := multierror.Append(nil,
		d.Set("role_id", role.ID),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting role assignment fields: %s", err)
	}

	return nil
}

func GetUserRoleAssignmentWithEpsID(client *golangsdk.ServiceClient, userID, roleID, enterpriseProjectID string) (eps_permissions.Role, error) {
	var assignment eps_permissions.Role

	roles, err := eps_permissions.UserPermissionsGet(client, enterpriseProjectID, userID).Extract()
	if err != nil {
		return assignment, err
	}

	for _, role := range roles {
		if role.ID == roleID {
			assignment = role
			break
		}
	}

	return assignment, nil
}

func resourceIdentityUserRoleAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	userID := d.Get("user_id").(string)
	roleID := d.Get("role_id").(string)
	enterpriseProjectID := d.Get("enterprise_project_id").(string)

	err = eps_permissions.UserPermissionsDelete(client, enterpriseProjectID, userID, roleID).ExtractErr()
	if err != nil {
		errMessage := fmt.Sprintf("error unassigning role (%s) from enterprise project (%s) and user (%s)",
			roleID, enterpriseProjectID, userID)

		return common.CheckDeletedDiag(d, err, errMessage)
	}

	return nil
}

func resourceIdentityUserRoleAssignmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id," +
			" must be <user_id>/<role_id>/<enterprise_project_id>")
	}

	d.Set("user_id", parts[0])
	d.Set("role_id", parts[1])
	d.Set("enterprise_project_id", parts[2])

	return []*schema.ResourceData{d}, nil
}
