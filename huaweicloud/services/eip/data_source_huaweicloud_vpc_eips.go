package eip

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v2.0/{project_id}/publicips/{id}/tags
// @API EIP GET /v1/{project_id}/publicips
func DataSourceVpcEips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
				Default:      4,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"eips": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcEipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Networking client: %s", err)
	}

	clientV2, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Networking V2 client: %s", err)
	}

	listOpts := &eips.ListOpts{
		Id:                  utils.ExpandToStringList(d.Get("ids").([]interface{})),
		PublicIp:            utils.ExpandToStringList(d.Get("public_ips").([]interface{})),
		PrivateIp:           utils.ExpandToStringList(d.Get("private_ips").([]interface{})),
		PortId:              utils.ExpandToStringList(d.Get("port_ids").([]interface{})),
		IPVersion:           d.Get("ip_version").(int),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to retrieve eips: %s ", err)
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve eips: %s ", err)
	}

	log.Printf("[DEBUG] Retrieved eips using given filter: %+v", allEips)

	var eipList []map[string]interface{}
	tagFilter := d.Get("tags").(map[string]interface{})
	var ids []string
	for _, item := range allEips {
		var tagRst map[string]string

		if resourceTags, err := tags.Get(clientV2, "publicips", item.ID).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)

			if !utils.HasMapContains(tagmap, tagFilter) {
				continue
			}
			tagRst = tagmap
		} else {
			// The tags api does not support eps authorization, so don't return 403 to avoid error
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				log.Printf("[WARN] Error query tags of EIP (%s): %s", item.ID, err)
			} else {
				return diag.Errorf("error query tags of EIP (%s): %s", item.ID, err)
			}
		}

		eip := map[string]interface{}{
			"id":                    item.ID,
			"name":                  item.Alias,
			"status":                NormalizeEipStatus(item.Status),
			"type":                  item.Type,
			"private_ip":            item.PrivateAddress,
			"public_ip":             item.PublicAddress,
			"public_ipv6":           item.PublicIpv6Address,
			"port_id":               item.PortID,
			"enterprise_project_id": item.EnterpriseProjectID,
			"ip_version":            item.IpVersion,
			"bandwidth_id":          item.BandwidthID,
			"bandwidth_size":        item.BandwidthSize,
			"bandwidth_name":        item.BandwidthName,
			"bandwidth_share_type":  item.BandwidthShareType,
			"tags":                  tagRst,
			"created_at":            item.CreateTime,
		}

		eipList = append(eipList, eip)
		ids = append(ids, item.ID)
	}
	log.Printf("[DEBUG]Eips List after filter, count=%d :%+v", len(eipList), eipList)

	mErr := multierror.Append(nil,
		d.Set("eips", eipList),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("set eips err:%s", err)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}
