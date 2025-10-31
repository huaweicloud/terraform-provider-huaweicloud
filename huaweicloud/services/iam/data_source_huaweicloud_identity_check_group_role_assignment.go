package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityCheckGroupRoleAssignment
// @API IAM HEAD /v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}
// @API IAM HEAD /v3/projects/{project_id}/groups/{group_id}/roles/{role_id}
// @API IAM HEAD /v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects
func DataSourceIdentityCheckGroupRoleAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCheckGroupRoleAssignmentRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"project_id"},
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityCheckGroupRoleAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	roleId := d.Get("role_id").(string)
	domainId := d.Get("domain_id").(string)
	projectId := d.Get("project_id").(string)
	var checkGroupRoleAssignmentPath string
	if domainId != "" {
		checkGroupRoleAssignmentPath = iamClient.Endpoint + "v3/domains/{domain_id}/groups/{group_id}/roles/{role_id}"
		checkGroupRoleAssignmentPath = strings.ReplaceAll(checkGroupRoleAssignmentPath, "{domain_id}", domainId)
	} else {
		if projectId == "all" {
			// "all" means to check whether the user group has the specified permissions
			// for all projects, including existing and future projects.
			checkGroupRoleAssignmentPath = iamClient.Endpoint +
				"v3/OS-INHERIT/domains/{domain_id}/groups/{group_id}/roles/{role_id}/inherited_to_projects"
			checkGroupRoleAssignmentPath = strings.ReplaceAll(checkGroupRoleAssignmentPath, "{domain_id}", cfg.DomainID)
		} else {
			checkGroupRoleAssignmentPath = iamClient.Endpoint + "v3/projects/{project_id}/groups/{group_id}/roles/{role_id}"
			checkGroupRoleAssignmentPath = strings.ReplaceAll(checkGroupRoleAssignmentPath, "{project_id}", projectId)
		}
	}
	checkGroupRoleAssignmentPath = strings.ReplaceAll(checkGroupRoleAssignmentPath, "{group_id}", groupId)
	checkGroupRoleAssignmentPath = strings.ReplaceAll(checkGroupRoleAssignmentPath, "{role_id}", roleId)
	options := golangsdk.RequestOpts{
		OkCodes: []int{204, 404},
	}
	response, err := iamClient.Request("HEAD", checkGroupRoleAssignmentPath, &options)
	if err != nil {
		return diag.Errorf("error checkGroupRoleAssignment: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate UUID: %s", err)
	}
	d.SetId(id)
	if response.StatusCode == 204 {
		err = d.Set("result", true)
	} else if response.StatusCode == 404 {
		err = d.Set("result", false)
	}
	if err != nil {
		return diag.Errorf("error set result filed: %s", err)
	}
	return nil
}
