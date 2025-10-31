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

// DataSourceIdentityCheckAgencyRoleAssignment
// @API IAM HEAD /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM HEAD /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM HEAD /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects
func DataSourceIdentityCheckAgencyRoleAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCheckAgencyRoleAssignmentRead,

		Schema: map[string]*schema.Schema{
			"agency_id": {
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

func dataSourceIdentityCheckAgencyRoleAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	agencyId := d.Get("agency_id").(string)
	roleId := d.Get("role_id").(string)
	domainId := d.Get("domain_id").(string)
	projectId := d.Get("project_id").(string)
	var checkAgencyRoleAssignmentPath string
	if domainId != "" {
		checkAgencyRoleAssignmentPath = iamClient.Endpoint + "v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}"
		checkAgencyRoleAssignmentPath = strings.ReplaceAll(checkAgencyRoleAssignmentPath, "{domain_id}", domainId)
	} else {
		if projectId == "all" {
			// "all" means to check whether the agency has the specified permissions
			// for all projects, including existing and future projects.
			checkAgencyRoleAssignmentPath = iamClient.Endpoint +
				"v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects"
			checkAgencyRoleAssignmentPath = strings.ReplaceAll(checkAgencyRoleAssignmentPath, "{domain_id}", cfg.DomainID)
		} else {
			checkAgencyRoleAssignmentPath = iamClient.Endpoint + "v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}"
			checkAgencyRoleAssignmentPath = strings.ReplaceAll(checkAgencyRoleAssignmentPath, "{project_id}", projectId)
		}
	}
	checkAgencyRoleAssignmentPath = strings.ReplaceAll(checkAgencyRoleAssignmentPath, "{agency_id}", agencyId)
	checkAgencyRoleAssignmentPath = strings.ReplaceAll(checkAgencyRoleAssignmentPath, "{role_id}", roleId)
	options := golangsdk.RequestOpts{
		OkCodes: []int{204, 404},
	}
	response, err := iamClient.Request("HEAD", checkAgencyRoleAssignmentPath, &options)
	if err != nil {
		return diag.Errorf("error checkAgencyRoleAssignment: %s", err)
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
