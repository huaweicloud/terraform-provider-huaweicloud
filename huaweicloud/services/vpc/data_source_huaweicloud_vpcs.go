package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcsRead,

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
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpcs": {
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
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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
						"secondary_cidrs": {
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

func dataSourceVpcsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	v1client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	v2Client, err := conf.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC V2 client: %s", err)
	}

	v3Client, err := conf.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: conf.DataGetEnterpriseProjectID(d),
	}

	vpcList, err := vpcs.List(v1client, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	log.Printf("[DEBUG] retrieved VPC using given filter: %+v", vpcList)

	var vpcs []map[string]interface{}
	tagFilter := d.Get("tags").(map[string]interface{})
	var ids []string
	for _, vpcResource := range vpcList {
		vpc := map[string]interface{}{
			"id":                    vpcResource.ID,
			"name":                  vpcResource.Name,
			"cidr":                  vpcResource.CIDR,
			"enterprise_project_id": vpcResource.EnterpriseProjectID,
			"status":                vpcResource.Status,
			"description":           vpcResource.Description,
		}

		if resourceTags, err := tags.Get(v2Client, "vpcs", vpcResource.ID).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)

			if !utils.HasMapContains(tagmap, tagFilter) {
				continue
			}
			vpc["tags"] = tagmap
		} else {
			// The tags api does not support eps authorization, so don't return 403 to avoid error
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				log.Printf("[WARN] error query tags of VPC (%s): %s", vpcResource.ID, err)
			} else {
				return diag.Errorf("error query tags of VPC (%s): %s", vpcResource.ID, err)
			}
		}

		// save VirtualPrivateCloudV3 extend_cidr
		res, err := obtainV3VpcResp(v3Client, vpcResource.ID)
		if err != nil {
			diag.Errorf("error retrieving VPC (%s) v3 detail: %s", vpcResource.ID, err)
		}
		vpc["secondary_cidrs"] = res.Vpc.ExtendCidrs

		vpcs = append(vpcs, vpc)
		ids = append(ids, vpcResource.ID)
	}
	log.Printf("[DEBUG] VPC List after filter, count: %d vpcs: %+v", len(vpcs), vpcs)

	mErr := d.Set("vpcs", vpcs)
	if mErr != nil {
		return diag.Errorf("set vpcs err: %s", mErr)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}
