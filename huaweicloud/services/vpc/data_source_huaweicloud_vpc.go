package vpc

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
		},
	}
}

func dataSourceVpcV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: config.DataGetEnterpriseProjectID(d),
	}

	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Vpc := refinedVpcs[0]

	logp.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("region", config.GetRegion(d))
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
	if vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d)); err == nil {
		if resourceTags, err := tags.Get(vpcV2Client, "vpcs", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return fmtp.DiagErrorf("Error saving tags to state for VPC (%s): %s", d.Id(), err)
			}
		} else {
			logp.Printf("[WARN] Error fetching tags of VPC (%s): %s", d.Id(), err)
		}
	}

	return nil
}
