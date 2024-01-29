package eps

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EPS GET /v1.0/enterprise-projects
func DataSourceEnterpriseProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenEnterpriseProjectResponseBody(project enterpriseprojects.Project) map[string]interface{} {
	result := map[string]interface{}{
		"id":          project.ID,
		"name":        project.Name,
		"type":        project.Type,
		"status":      project.Status,
		"description": project.Description,
		"created_at":  project.CreatedAt,
		"updated_at":  project.UpdatedAt,
	}

	return result
}

func dataSourceEnterpriseProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.EnterpriseProjectClient(region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	opt := enterpriseprojects.ListOpts{
		Name:    d.Get("name").(string),
		ID:      d.Get("enterprise_project_id").(string),
		Status:  d.Get("status").(int),
		Type:    d.Get("type").(string),
		SortKey: "name",
		SortDir: "asc",
	}
	projects, err := enterpriseprojects.List(client, opt).Extract()

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving enterprise projects")
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := make([]interface{}, 0, len(projects))
	for _, project := range projects {
		result = append(result, flattenEnterpriseProjectResponseBody(project))
	}

	mErr := multierror.Append(nil,
		d.Set("enterprise_projects", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
