package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPC GET /v1/{project_id}/vpcs
func DataSourceVpcIdsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcIdsV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVpcIdsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Vpc client: %s", err)
	}

	listOpts := vpcs.ListOpts{}
	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	listVpcs := make([]string, 0)

	for _, vpc := range refinedVpcs {
		listVpcs = append(listVpcs, vpc.ID)
	}
	d.SetId(listVpcs[0])
	d.Set("ids", listVpcs)

	return nil
}
