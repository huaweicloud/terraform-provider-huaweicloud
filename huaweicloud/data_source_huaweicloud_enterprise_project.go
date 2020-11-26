package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/eps/v1/enterpriseprojects"
)

func DataSourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnterpriseProjectRead,
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

func dataSourceEnterpriseProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	epsClient, err := config.EnterpriseProjectClient(region)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud eps client %s", err)
	}

	listOpts := enterpriseprojects.ListOpts{
		Name:   d.Get("name").(string),
		ID:     d.Get("id").(string),
		Status: d.Get("status").(int),
	}
	projects, err := enterpriseprojects.List(epsClient, listOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error retriving enterprise projects %s", err)
	}

	if len(projects) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(projects) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	project := projects[0]

	d.SetId(project.ID)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("status", project.Status)
	d.Set("created_at", project.CreatedAt)
	d.Set("updated_at", project.UpdatedAt)

	return nil
}
