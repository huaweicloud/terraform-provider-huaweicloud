package iam

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/chnsz/golangsdk/pagination"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityRoleAssignmentV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleAssignmentV3Create,
		ReadContext:   resourceIdentityRoleAssignmentV3Read,
		DeleteContext: resourceIdentityRoleAssignmentV3Delete,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"project_id", "project_name"},
				Optional:     true,
				ForceNew:     true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func GetProjectId(client *golangsdk.ServiceClient, projectInfo string) (string, error) {
	projectName, projectID := flattenProjcetInfo(projectInfo)
	if projectID != "" {
		return projectID, nil
	}

	listOpts := projects.ListOpts{
		Name: projectName,
	}
	pages, err := projects.List(client, listOpts).AllPages()
	if err != nil {
		return "", fmt.Errorf("error retrieving project list: %v", err)
	}
	projectList, err := projects.ExtractProjects(pages)
	if err != nil {
		return "", fmt.Errorf("error fetching project object: %v", err)
	}
	if len(projectList) < 1 {
		return "", fmt.Errorf("unable to find any project by name (%s), please check your input", projectName)
	}
	if len(projectList) > 1 {
		return "", fmt.Errorf("more than one projects are found within project name (%s), please use project ID "+
			"instead", projectName)
	}
	return projectList[0].ID, nil
}

func resourceIdentityRoleAssignmentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	roleID := d.Get("role_id").(string)
	groupID := d.Get("group_id").(string)
	domainID := d.Get("domain_id").(string)

	var projectID string
	if val, ok := d.GetOk("project_id"); ok {
		projectID = val.(string)
	} else {
		projectID, err = GetProjectId(identityClient, d.Get("project_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	opts := roles.AssignOpts{
		GroupID:   groupID,
		DomainID:  domainID,
		ProjectID: projectID,
	}

	err = roles.Assign(identityClient, roleID, opts).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error assigning role: %s", err)
	}

	if val, ok := d.GetOk("project_name"); ok {
		d.SetId(buildRoleAssignmentID(domainID, val.(string), groupID, roleID))
	} else {
		d.SetId(buildRoleAssignmentID(domainID, projectID, groupID, roleID))
	}

	return resourceIdentityRoleAssignmentV3Read(ctx, d, meta)
}

func flattenProjcetInfo(projectInfo string) (projectName, projectID string) {
	re, _ := regexp.Compile(`[a-f0-9]{32}$`) // project ID format
	if re.FindString(projectInfo) != "" {
		projectID = projectInfo
		return
	}
	projectName = projectInfo
	return
}

func resourceIdentityRoleAssignmentV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	roleAssignment, err := getRoleAssignment(identityClient, d)
	if err != nil {
		return fmtp.DiagErrorf("Error getting role assignment: %s", err)
	}

	domainID, projectInfo, groupID, _ := ExtractRoleAssignmentID(d.Id())
	logp.Printf("[DEBUG] Retrieved HuaweiCloud role assignment: %#v", roleAssignment)

	projectName, projcetId := flattenProjcetInfo(projectInfo)
	mErr := multierror.Append(nil,
		d.Set("role_id", roleAssignment.ID),
		d.Set("group_id", groupID),
		d.Set("domain_id", domainID),
		d.Set("project_name", projectName),
		d.Set("project_id", projcetId),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity role assignment fields: %s", err)
	}

	return nil
}

func resourceIdentityRoleAssignmentV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	domainID, projectInfo, groupID, roleID := ExtractRoleAssignmentID(d.Id())
	projectID, err := GetProjectId(identityClient, projectInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := roles.UnassignOpts{
		GroupID:   groupID,
		DomainID:  domainID,
		ProjectID: projectID,
	}
	err = roles.Unassign(identityClient, roleID, opts).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error unassigning role: %s", err)
	}

	return nil
}

func getRoleAssignment(identityClient *golangsdk.ServiceClient, d *schema.ResourceData) (*roles.RoleAssignment, error) {
	domainID, projectInfo, groupID, roleID := ExtractRoleAssignmentID(d.Id())
	projectID, err := GetProjectId(identityClient, projectInfo)
	if err != nil {
		return nil, err
	}

	opts := roles.ListAssignmentsOpts{
		GroupID:        groupID,
		ScopeDomainID:  domainID,
		ScopeProjectID: projectID,
	}
	pager := roles.ListAssignments(identityClient, opts)
	var assignment roles.RoleAssignment

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
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

	return &assignment, err
}

// Role assignments have no ID in HuaweiCloud. Build an ID out of the IDs (maybe contain project name, not ID) that make
// up the role assignment.
// Input: domain ID, project ID (or project name), group ID, role ID
func buildRoleAssignmentID(domainID, projectInfo, groupID, roleID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", domainID, projectInfo, groupID, roleID)
}

// ExtractRoleAssignmentID is a method to flatten a composite ID format into multiple IDs.
// Return: domain ID, project ID (or project name), group ID, role ID
func ExtractRoleAssignmentID(roleAssignmentID string) (string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	return split[0], split[1], split[2], split[3]
}
