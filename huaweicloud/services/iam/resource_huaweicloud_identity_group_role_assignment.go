package iam

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/eps_permissions"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3GroupRoleAssignmentNonUpdatableParams = []string{
	"group_id",
	"role_id",
	"domain_id",
	"project_id",
	"enterprise_project_id",
}

// @API IAM PUT /v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}
// @API IAM PUT /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects
// @API IAM PUT /v3/projects/{project_id}/groups/{group_id}/roles/{role_id}
// @API IAM PUT /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles/{role_id}
// @API IAM GET /v3/domains/{domain_id}/groups/{group_id}/roles
// @API IAM GET /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/inherited_to_projects
// @API IAM GET /v3/projects/{project_id}/groups/{group_id}/roles
// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles
// @API IAM DELETE /v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}
// @API IAM DELETE /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects
// @API IAM DELETE /v3/projects/{project_id}/groups/{group_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles/{role_id}
func ResourceV3GroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3GroupRoleAssignmentCreate,
		ReadContext:   resourceV3GroupRoleAssignmentRead,
		UpdateContext: resourceV3GroupRoleAssignmentUpdate,
		DeleteContext: resourceV3GroupRoleAssignmentDelete,

		CustomizeDiff: config.FlexibleForceNew(v3GroupRoleAssignmentNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3GroupRoleAssignmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of user group to which the role to be authorized belongs.`,
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of role to be authorized.`,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"enterprise_project_id",
				},
				Description: `The ID of domain to assign the role in.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of project to assign the role in.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of enterprise project to assign the role in.`,
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
		},
	}
}

func assignRoleWithDomainId(client *golangsdk.ServiceClient, groupId, roleId, domainId string) error {
	opts := roles.AssignOpts{
		GroupID:  groupId,
		DomainID: domainId,
	}
	return roles.Assign(client, roleId, opts).ExtractErr()
}

func assignRoleWithProjectId(client *golangsdk.ServiceClient, groupId, roleId, domainId, projectId string) error {
	// The value of "all" means that the specified user group will be able to use all projects,
	// including existing and future projects.
	if projectId == "all" {
		if domainId == "" {
			return errors.New("the domain_id must be specified in the provider configuration")
		}
		return roles.AssignAllResources(client, domainId, groupId, roleId).ExtractErr()
	}

	// Assign role to a specified project.
	opts := roles.AssignOpts{
		GroupID:   groupId,
		ProjectID: projectId,
	}
	return roles.Assign(client, roleId, opts).ExtractErr()
}

func assignRoleWithEpsId(client *golangsdk.ServiceClient, groupId, roleId, epsId string) error {
	return eps_permissions.UserGroupPermissionsCreate(client, epsId, groupId, roleId).ExtractErr()
}

func resourceV3GroupRoleAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		resourceId    string
		err           error
		iamV3Client   *golangsdk.ServiceClient
		iamV3P0Client *golangsdk.ServiceClient

		roleId  = d.Get("role_id").(string)
		groupId = d.Get("group_id").(string)
	)

	iamV3Client, err = cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}
	iamV3P0Client, err = cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	if domainId, ok := d.GetOk("domain_id"); ok {
		err = assignRoleWithDomainId(iamV3Client, groupId, roleId, domainId.(string))
		// The value of parameter 'domain_id' is used as part of the resource ID.
		resourceId = fmt.Sprintf("%s/%s/%v", groupId, roleId, domainId)
	}
	if projectId, ok := d.GetOk("project_id"); ok {
		err = assignRoleWithProjectId(iamV3Client, groupId, roleId, cfg.DomainID, projectId.(string))
		// The value of parameter 'project_id' is used as part of the resource ID.
		resourceId = fmt.Sprintf("%s/%s/%v", groupId, roleId, projectId)
	}
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		err = assignRoleWithEpsId(iamV3P0Client, groupId, roleId, epsId.(string))
		// The value of parameter 'enterprise_project_id' is used as part of the resource ID.
		resourceId = fmt.Sprintf("%s/%s/%v", groupId, roleId, epsId)
	}
	if err != nil {
		return diag.Errorf("error assigning role (%s) to group (%s): %s", roleId, groupId, err)
	}

	d.SetId(resourceId)

	return resourceV3GroupRoleAssignmentRead(ctx, d, meta)
}

func CheckV3GroupRoleAssignmentWithDomainId(client *golangsdk.ServiceClient, groupId, roleId, domainId string) error {
	opts := roles.ListAssignmentsOpts{
		GroupID:       groupId,
		ScopeDomainID: domainId,
	}

	var isAssigned bool
	err := roles.ListAssignments(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, assignment := range assignmentList {
			if assignment.ID == roleId {
				isAssigned = true
				return false, nil
			}
		}

		return true, nil
	})

	if err == nil && !isAssigned {
		return golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/domains/{domain_id}/groups/{group_id}/roles",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("the role (%s), which assigned to the domain (%s), is not assigned to the group (%s)",
					roleId, domainId, groupId)),
			},
		}
	}
	return err
}

func CheckV3GroupRoleAssignmentWithProjectId(client *golangsdk.ServiceClient, groupId, roleId, domainId, projectId string) error {
	if projectId == "all" {
		if domainId == "" {
			return errors.New("the parameter 'domain_id' must be specified in the provider configuration")
		}
		return roles.CheckAllResourcesPermission(client, domainId, groupId, roleId).ExtractErr()
	}

	opts := roles.ListAssignmentsOpts{
		GroupID:        groupId,
		ScopeProjectID: projectId,
	}

	var isAssigned bool
	err := roles.ListAssignments(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, assignment := range assignmentList {
			if assignment.ID == roleId {
				isAssigned = true
				return false, nil
			}
		}

		return true, nil
	})

	if err == nil && !isAssigned {
		return golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/projects/{project_id}/groups/{group_id}/roles",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("the role (%s), which assigned to the project (%s), is not assigned to the group (%s)",
					roleId, projectId, groupId)),
			},
		}
	}
	return err
}

func CheckV3GroupRoleAssignmentWithEpsId(client *golangsdk.ServiceClient, groupId, roleId, epsId string) error {
	allRoles, err := eps_permissions.UserGroupPermissionsGet(client, epsId, groupId).Extract()
	if err != nil {
		return err
	}

	for _, role := range allRoles {
		if role.ID == roleId {
			return nil
		}
	}
	return golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v3.0/OS-PERMISSION/enterprise-projects/{enterpriseProjectID}/groups/{group_id}/roles",
			RequestId: "NONE",
			Body: []byte(fmt.Sprintf("the role (%s), which assigned to the enterprise project (%s), is not assigned to the group (%s)",
				roleId, epsId, groupId)),
		},
	}
}

func resourceV3GroupRoleAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		err           error
		iamV3Client   *golangsdk.ServiceClient
		iamV3P0Client *golangsdk.ServiceClient

		roleId  = d.Get("role_id").(string)
		groupId = d.Get("group_id").(string)
	)
	iamV3Client, err = cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	iamV3P0Client, err = cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	if domainId, ok := d.GetOk("domain_id"); ok {
		err = CheckV3GroupRoleAssignmentWithDomainId(iamV3Client, groupId, roleId, domainId.(string))
	}
	if projectId, ok := d.GetOk("project_id"); ok {
		err = CheckV3GroupRoleAssignmentWithProjectId(iamV3Client, groupId, roleId, cfg.DomainID, projectId.(string))
	}
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		err = CheckV3GroupRoleAssignmentWithEpsId(iamV3P0Client, groupId, roleId, epsId.(string))
	}

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving role assignment")
	}
	return nil
}

func resourceV3GroupRoleAssignmentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// All parameters are non-updatable for this resource.
	return nil
}

func unassignRoleWithDomainId(client *golangsdk.ServiceClient, groupId, roleId, domainId string) error {
	opts := roles.UnassignOpts{
		GroupID:  groupId,
		DomainID: domainId,
	}
	return roles.Unassign(client, roleId, opts).ExtractErr()
}

func unassignRoleWithProjectId(client *golangsdk.ServiceClient, groupId, roleId, domainId, projectId string) error {
	if projectId == "all" {
		if domainId == "" {
			return errors.New("the parameter 'domain_id' must be specified in the provider configuration")
		}
		return roles.UnassignAllResources(client, domainId, groupId, roleId).ExtractErr()
	}
	opts := roles.UnassignOpts{
		GroupID:   groupId,
		ProjectID: projectId,
	}
	return roles.Unassign(client, roleId, opts).ExtractErr()
}

func unassignRoleWithEpsId(client *golangsdk.ServiceClient, groupId, roleId, epsId string) error {
	return eps_permissions.UserGroupPermissionsDelete(client, epsId, groupId, roleId).ExtractErr()
}

func resourceV3GroupRoleAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		err           error
		iamV3Client   *golangsdk.ServiceClient
		iamV3P0Client *golangsdk.ServiceClient

		roleId  = d.Get("role_id").(string)
		groupId = d.Get("group_id").(string)
	)
	iamV3Client, err = cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}
	iamV3P0Client, err = cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	if domainId, ok := d.GetOk("domain_id"); ok {
		err = unassignRoleWithDomainId(iamV3Client, groupId, roleId, domainId.(string))
	}
	if projectId, ok := d.GetOk("project_id"); ok {
		err = unassignRoleWithProjectId(iamV3Client, groupId, roleId, cfg.DomainID, projectId.(string))
	}
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		err = unassignRoleWithEpsId(iamV3P0Client, groupId, roleId, epsId.(string))
	}
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error unassigning role")
	}
	return nil
}

func setV3GroupRoleAssignmentInfoWithType(d *schema.ResourceData, idInfo, assignmentType string) ([]*schema.ResourceData, error) {
	var (
		err = fmt.Errorf(`invalid format specified for import ID, want these following format:
1. <group_id>/<role_id>/<domain_id>:domain
2. <group_id>/<role_id>/all:project
3. <group_id>/<role_id>/<project_id>:project
4. <group_id>/<role_id>/<domain_id>:enterprise_project
but got '%s:%s'`, idInfo, assignmentType)
		parts = strings.Split(idInfo, "/")
	)

	if len(parts) < 3 {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(idInfo)
	d.Set("group_id", parts[0])
	d.Set("role_id", parts[1])

	switch assignmentType {
	case "domain":
		err = d.Set("domain_id", parts[2])
	case "project":
		err = d.Set("project_id", parts[2])
	case "enterprise_project":
		err = d.Set("enterprise_project_id", parts[2])
	default:
	}
	return []*schema.ResourceData{d}, err
}

func setV3GroupRoleAssignmentInfoWithoutType(d *schema.ResourceData, cfg *config.Config, iamV3Client, iamV3P0Client *golangsdk.ServiceClient,
	idInfo string) ([]*schema.ResourceData, error) {
	var (
		err = fmt.Errorf(`invalid format specified for import ID, want these following format:
1. <group_id>/<role_id>/<domain_id>
2. <group_id>/<role_id>/all
3. <group_id>/<role_id>/<project_id>
4. <group_id>/<role_id>/<enterprise_project_id>
but got '%s'`, idInfo)
		parts = strings.Split(idInfo, "/")
	)

	if len(parts) < 3 {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(idInfo)
	d.Set("group_id", parts[0])
	d.Set("role_id", parts[1])

	// Role authorization checks are performed in the order of domain, project, and enterprise_project.
	// If the check passes, the corresponding parameters are set.
	if err = CheckV3GroupRoleAssignmentWithDomainId(iamV3Client, parts[0], parts[1], parts[2]); err == nil {
		return []*schema.ResourceData{d}, d.Set("domain_id", parts[2])
	}

	if err = CheckV3GroupRoleAssignmentWithProjectId(iamV3Client, parts[0], parts[1], cfg.DomainID, parts[2]); err == nil {
		return []*schema.ResourceData{d}, d.Set("project_id", parts[2])
	}

	if err = CheckV3GroupRoleAssignmentWithEpsId(iamV3P0Client, parts[0], parts[1], parts[2]); err == nil {
		return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[2])
	}

	return []*schema.ResourceData{d}, err
}

func resourceV3GroupRoleAssignmentImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		err           error
		iamV3Client   *golangsdk.ServiceClient
		iamV3P0Client *golangsdk.ServiceClient

		importedId = d.Id()
	)
	iamV3Client, err = cfg.IdentityV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3 client: %s", err)
	}
	iamV3P0Client, err = cfg.IAMV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM v3.0 client: %s", err)
	}

	log.Printf("[Lance] The imported ID is: %s", importedId)

	parts := strings.Split(importedId, ":")
	if len(parts) >= 2 {
		return setV3GroupRoleAssignmentInfoWithType(d, parts[0], parts[1])
	}
	return setV3GroupRoleAssignmentInfoWithoutType(d, cfg, iamV3Client, iamV3P0Client, importedId)
}
