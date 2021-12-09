package vpc

import (
	"context"
	"strings"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
					},
				},
			},
		},
	}
}

func dataSourceVpcsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	vpcV2Client, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC V2 client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	vpcList, err := vpcs.List(client, listOpts)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve vpcs: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved Vpc using given filter: %+v", vpcList)

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

		if resourceTags, err := tags.Get(vpcV2Client, "vpcs", vpcResource.ID).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			tags, isMatch := filterByTags(tagmap, tagFilter)
			if isMatch {
				vpc["tags"] = tags
			} else {
				continue
			}
		} else {
			return fmtp.DiagErrorf("Error query tags of VPC (%s): %s", d.Id(), err)
		}

		vpcs = append(vpcs, vpc)
		ids = append(ids, vpcResource.ID)
	}
	logp.Printf("[DEBUG]Vpc List after filter, count=%d :%+v", len(vpcs), vpcs)

	mErr := d.Set("vpcs", vpcs)
	if mErr != nil {
		return fmtp.DiagErrorf("set vpcs err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}

// if filterTags = {"foo":"a,b"} , the rawTags must hava a key is "foo" and value is "a" OR "b"
func filterByTags(rawTags map[string]string, filterTags map[string]interface{}) (map[string]string, bool) {
	if len(filterTags) < 1 {
		return rawTags, true
	}
	if len(filterTags) > 0 && len(rawTags) < 1 {
		return nil, false
	}

	isMatch := true
	for key, value := range filterTags {
		isMatch = isMatch && isTagMatch(rawTags, key, value.(string))
	}

	if isMatch {
		return rawTags, true
	}

	return nil, false
}

func isTagMatch(rawTags map[string]string, filterKey, filterValue string) bool {
	if rawTag, ok := rawTags[filterKey]; ok {
		if filterValue != "" {
			filterTagValues := strings.Split(filterValue, ",")
			return utils.StrSliceContains(filterTagValues, rawTag)

		} else {
			return true
		}

	} else {
		return false
	}

}
