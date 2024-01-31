package ecs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS GET /v1/{project_id}/cloudservers/os-server-groups
func DataSourceComputeServerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeServerGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"servergroups": {
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
						"policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"members": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceComputeServerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	ecsClient, err := conf.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	pages, err := servergroups.List(ecsClient).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list server groups: %s", err)
	}

	allServerGroups, err := servergroups.ExtractServerGroups(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve server groups: %s", err)
	}

	filter := map[string]interface{}{
		"Name": d.Get("name").(string),
	}
	filteredServerGroups, err := utils.FilterSliceWithField(allServerGroups, filter)

	if err != nil {
		return diag.Errorf("filter server groups failed: %s", err)
	}

	var serverGroupsToSet []map[string]interface{}
	var serverGroupsIds []string

	for _, item := range filteredServerGroups {
		serverGroupInAll := item.(servergroups.ServerGroup)
		serverGroupID := serverGroupInAll.ID
		serverGroupsIds = append(serverGroupsIds, serverGroupID)

		serverGroupToSet := map[string]interface{}{
			"id":       serverGroupID,
			"name":     serverGroupInAll.Name,
			"members":  serverGroupInAll.Members,
			"policies": serverGroupInAll.Policies,
		}
		serverGroupsToSet = append(serverGroupsToSet, serverGroupToSet)
	}

	d.SetId(hashcode.Strings(serverGroupsIds))
	mErr := multierror.Append(nil,
		d.Set("servergroups", serverGroupsToSet),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
