package eps

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
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
	}
}

func dataSourceEnterpriseProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	epsClient, err := config.EnterpriseProjectClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud eps client %s", err)
	}

	listOpts := enterpriseprojects.ListOpts{
		Name:   d.Get("name").(string),
		ID:     d.Get("id").(string),
		Status: d.Get("status").(int),
	}
	projects, err := enterpriseprojects.List(epsClient, listOpts).Extract()

	if err != nil {
		return fmtp.DiagErrorf("Error retriving enterprise projects %s", err)
	}

	if len(projects) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(projects) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	project := projects[0]

	d.SetId(project.ID)
	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("status", project.Status),
		d.Set("created_at", project.CreatedAt),
		d.Set("updated_at", project.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting HuaweiCloud enterprise project fields: %w", err)
	}

	return nil
}
