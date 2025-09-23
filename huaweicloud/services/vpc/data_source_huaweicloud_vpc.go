package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v2.0/{project_id}/vpcs/{id}/tags
// @API VPC GET /v1/{project_id}/vpcs
func DataSourceVpcV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routes": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "use huaweicloud_vpc_route_table data source to get all routes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secondary_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVpcV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	refinedVpcs, err := vpcs.List(v1Client, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Vpc := refinedVpcs[0]

	log.Printf("[INFO] retrieved VPC using given filter (%s): %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("region", region)
	d.Set("name", Vpc.Name)
	d.Set("cidr", Vpc.CIDR)
	d.Set("enterprise_project_id", Vpc.EnterpriseProjectID)
	d.Set("status", Vpc.Status)
	d.Set("description", Vpc.Description)

	var s []map[string]interface{}
	for _, route := range Vpc.Routes {
		mapping := map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
		s = append(s, mapping)
	}
	d.Set("routes", s)

	// save VirtualPrivateCloudV2 tags
	if v2Client, err := cfg.NetworkingV2Client(region); err == nil {
		if resourceTags, err := tags.Get(v2Client, "vpcs", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return diag.Errorf("error saving tags to state for VPC (%s): %s", d.Id(), err)
			}
		} else {
			log.Printf("[WARN] Error fetching tags of VPC (%s): %s", d.Id(), err)
		}
	}

	// save VirtualPrivateCloudV3 extend_cidr
	v3Client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	res, err := obtainV3VpcResp(v3Client, d.Id())
	if err != nil {
		diag.Errorf("error retrieving VPC (%s) v3 detail: %s", d.Id(), err)
	}
	d.Set("secondary_cidrs", utils.PathSearch("vpc.extend_cidrs", res, nil))

	return nil
}
