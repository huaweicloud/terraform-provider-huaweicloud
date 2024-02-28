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
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM DELETE /v3/projects/{project_id}/groups/{group_id}/roles/{role_id}
// @API IAM PUT /v3/projects/{project_id}/groups/{group_id}/roles/{role_id}
// @API IAM GET /v3/projects/{project_id}/groups/{group_id}/roles
// @API IAM DELETE /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects
// @API IAM PUT /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects
// @API IAM GET /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/inherited_to_projects
// @API IAM DELETE /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles/{role_id}
// @API IAM PUT /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles/{role_id}
// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles
// @API IAM DELETE /v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}
// @API IAM PUT /v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}
// @API IAM GET /v3/domains/{domain_id}/groups/{group_id}/roles
func ResourceIdentityGroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupRoleAssignmentCreate,
		ReadContext:   resourceIdentityGroupRoleAssignmentRead,
		DeleteContext: resourceIdentityGroupRoleAssignmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityGroupRoleAssignmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id", "enterprise_project_id",
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityGroupRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	identityClient, err := conf.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := conf.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	var resourceID string
	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	if v, ok := d.GetOk("domain_id"); ok {
		domainID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, domainID)

		opts := roles.AssignOpts{
			GroupID:  groupID,
			DomainID: domainID,
		}
		err = roles.Assign(identityClient, roleID, opts).ExtractErr()
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, projectID)

		if projectID == "all" {
			// "all" means that the specified user group will be able to use all projects,
			// including existing and future projects.
			if conf.DomainID == "" {
				return diag.Errorf("the domain_id must be specified in the provider configuration")
			}
			err = roles.AssignAllResources(identityClient, conf.DomainID, groupID, roleID).ExtractErr()
		} else {
			opts := roles.AssignOpts{
				GroupID:   groupID,
				ProjectID: projectID,
			}
			err = roles.Assign(identityClient, roleID, opts).ExtractErr()
		}
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		enterpriseProjectID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, enterpriseProjectID)

		err = eps_permissions.UserGroupPermissionsCreate(iamClient, enterpriseProjectID, groupID, roleID).ExtractErr()
	}

	if err != nil {
		return diag.Errorf("error assigning role: %s", err)
	}

	d.SetId(resourceID)
	return resourceIdentityGroupRoleAssignmentRead(ctx, d, meta)
}

func resourceIdentityGroupRoleAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	identityClient, err := conf.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := conf.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	var resourceID, respRoleID string
	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	if v, ok := d.GetOk("domain_id"); ok {
		domainID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, domainID)

		roleAssignment, err := GetGroupRoleAssignmentWithDomainID(identityClient, groupID, roleID, domainID)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error getting role assignment")
		}

		respRoleID = roleAssignment.ID
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, projectID)

		var roleAssignment roles.RoleAssignment
		if projectID == "all" {
			if conf.DomainID == "" {
				return diag.Errorf("the domain_id must be specified in the provider configuration")
			}
			roleAssignment.ID = roleID
			err = roles.CheckAllResourcesPermission(identityClient, conf.DomainID, groupID, roleID).ExtractErr()
		} else {
			roleAssignment, err = GetGroupRoleAssignmentWithProjectID(identityClient, groupID, roleID, projectID)
		}

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error getting role assignment")
		}
		respRoleID = roleAssignment.ID
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		enterpriseProjectID := v.(string)
		resourceID = fmt.Sprintf("%s/%s/%s", groupID, roleID, enterpriseProjectID)

		role, err := GetGroupRoleAssignmentWithEpsID(iamClient, groupID, roleID, enterpriseProjectID)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error getting role assignment")
		}

		respRoleID = role.ID
	}

	// update the resource ID for `huaweicloud_identity_role_assignment`
	d.SetId(resourceID)
	mErr := multierror.Append(nil,
		d.Set("role_id", respRoleID),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting role assignment fields: %s", err)
	}

	return nil
}

func GetGroupRoleAssignmentWithDomainID(identityClient *golangsdk.ServiceClient, groupID, roleID, domainID string) (roles.RoleAssignment, error) {
	opts := roles.ListAssignmentsOpts{
		GroupID:       groupID,
		ScopeDomainID: domainID,
	}

	pager := roles.ListAssignments(identityClient, opts)
	var assignment roles.RoleAssignment

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, a := range assignmentList {
			if a.ID == roleID {
				assignment = a
				return false, nil
			}
		}

		return true, nil
	})

	return assignment, err
}

func GetGroupRoleAssignmentWithProjectID(identityClient *golangsdk.ServiceClient, groupID, roleID, projectID string) (roles.RoleAssignment, error) {
	opts := roles.ListAssignmentsOpts{
		GroupID:        groupID,
		ScopeProjectID: projectID,
	}

	pager := roles.ListAssignments(identityClient, opts)
	var assignment roles.RoleAssignment

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, a := range assignmentList {
			if a.ID == roleID {
				assignment = a
				return false, nil
			}
		}

		return true, nil
	})

	return assignment, err
}

func GetGroupRoleAssignmentWithEpsID(iamClient *golangsdk.ServiceClient, groupID, roleID, enterpriseProjectID string) (eps_permissions.Role, error) {
	var assignment eps_permissions.Role

	allRole, err := eps_permissions.UserGroupPermissionsGet(iamClient, enterpriseProjectID, groupID).Extract()
	if err != nil {
		return assignment, err
	}

	for _, role := range allRole {
		if role.ID == roleID {
			assignment = role
			break
		}
	}

	return assignment, nil
}

func resourceIdentityGroupRoleAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	identityClient, err := conf.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := conf.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)

	if v, ok := d.GetOk("domain_id"); ok {
		domainID := v.(string)

		opts := roles.UnassignOpts{
			GroupID:  groupID,
			DomainID: domainID,
		}
		err = roles.Unassign(identityClient, roleID, opts).ExtractErr()
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)

		if projectID == "all" {
			if conf.DomainID == "" {
				return diag.Errorf("the domain_id must be specified in the provider configuration")
			}

			err = roles.UnassignAllResources(identityClient, conf.DomainID, groupID, roleID).ExtractErr()
		} else {
			opts := roles.UnassignOpts{
				GroupID:   groupID,
				ProjectID: projectID,
			}
			err = roles.Unassign(identityClient, roleID, opts).ExtractErr()
		}
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		enterpriseProjectID := v.(string)
		err = eps_permissions.UserGroupPermissionsDelete(iamClient, enterpriseProjectID, groupID, roleID).ExtractErr()
	}

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error unassigning role")
	}
	return nil
}

func resourceIdentityGroupRoleAssignmentImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	identityClient, err := conf.IdentityV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3 client: %s", err)
	}

	iamClient, err := conf.IAMV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id," +
			" must be <group_id>/<role_id>/<domain_id>, <group_id>/<role_id>/<project_id> or <group_id>/<role_id>/<enterprise_project_id>")
	}

	d.Set("group_id", parts[0])
	d.Set("role_id", parts[1])

	if parts[2] == "all" {
		d.Set("project_id", parts[2])
		return []*schema.ResourceData{d}, nil
	}

	if _, err = GetGroupRoleAssignmentWithDomainID(identityClient, parts[0], parts[1], parts[2]); err == nil {
		d.Set("domain_id", parts[2])
		return []*schema.ResourceData{d}, nil
	}

	if _, err = GetGroupRoleAssignmentWithProjectID(identityClient, parts[0], parts[1], parts[2]); err == nil {
		d.Set("project_id", parts[2])
		return []*schema.ResourceData{d}, nil
	}

	if _, err = GetGroupRoleAssignmentWithEpsID(iamClient, parts[0], parts[1], parts[2]); err == nil {
		d.Set("enterprise_project_id", parts[2])
		return []*schema.ResourceData{d}, nil
	}

	return nil, fmt.Errorf("error importing role assignment: %s", err)
}
