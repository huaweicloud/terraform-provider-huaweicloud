package eip

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP GET /v3/{domain_id}/global-eips
func DataSourceGlobalEIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEIPsRead,

		Schema: map[string]*schema.Schema{
			"geip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"associate_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associate_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"global_eips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"geip_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frozen": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"frozen_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"polluted": {
							Type:     schema.TypeBool,
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"global_connection_bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"global_connection_bandwidth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associate_instance_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associate_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associate_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGlobalEIPsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getGEIPHttpUrl := "v3/{domain_id}/global-eips"
	getGEIPPath := client.Endpoint + getGEIPHttpUrl
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", cfg.DomainID)
	getGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGEIPPath += fmt.Sprintf("?limit=%v", pageLimit)
	if id, ok := d.GetOk("geip_id"); ok {
		getGEIPPath += fmt.Sprintf("&id=%s", id)
	}
	if name, ok := d.GetOk("name"); ok {
		getGEIPPath += fmt.Sprintf("&name=%s", name)
	}
	if bwID, ok := d.GetOk("internet_bandwidth_id"); ok {
		getGEIPPath += fmt.Sprintf("&internet_bandwidth_id=%s", bwID)
	}
	if ipAddress, ok := d.GetOk("ip_address"); ok {
		getGEIPPath += fmt.Sprintf("&ip_address=%s", ipAddress)
	}
	if status, ok := d.GetOk("status"); ok {
		getGEIPPath += fmt.Sprintf("&status=%s", status)
	}
	if epsID, ok := d.GetOk("enterprise_project_id"); ok {
		getGEIPPath += fmt.Sprintf("&enterprise_project_id=%s", epsID)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		tagsList := expandTagsMapToStringList(rawTags.(map[string]interface{}))
		for _, v := range tagsList {
			getGEIPPath += fmt.Sprintf("&tags=%s", v)
		}
	}
	if associateInstanceID, ok := d.GetOk("associate_instance_id"); ok {
		getGEIPPath += fmt.Sprintf("&associate_instance_info.instance_id=%s", associateInstanceID)
	}
	if associateInstanceType, ok := d.GetOk("associate_instance_type"); ok {
		getGEIPPath += fmt.Sprintf("&associate_instance_info.instance_type=%s", associateInstanceType)
	}

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getGEIPPath + fmt.Sprintf("&offset=%d", currentTotal)
		getGEIPResp, err := client.Request("GET", currentPath, &getGEIPOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getGEIPRespBody, err := utils.FlattenResponse(getGEIPResp)
		if err != nil {
			return diag.FromErr(err)
		}

		geips := utils.PathSearch("global_eips", getGEIPRespBody, make([]interface{}, 0)).([]interface{})
		for _, geip := range geips {
			results = append(results, map[string]interface{}{
				"id":                               utils.PathSearch("id", geip, nil),
				"access_site":                      utils.PathSearch("access_site", geip, nil),
				"geip_pool_name":                   utils.PathSearch("geip_pool_name", geip, nil),
				"internet_bandwidth_id":            utils.PathSearch("internet_bandwidth_info.id", geip, nil),
				"isp":                              utils.PathSearch("isp", geip, nil),
				"ip_version":                       int(utils.PathSearch("ip_version", geip, float64(0)).(float64)),
				"description":                      utils.PathSearch("description", geip, nil),
				"name":                             utils.PathSearch("name", geip, nil),
				"ip_address":                       utils.PathSearch("ip_address", geip, nil),
				"enterprise_project_id":            utils.PathSearch("enterprise_project_id", geip, nil),
				"frozen":                           utils.PathSearch("frozen_info", geip, false),
				"frozen_info":                      utils.PathSearch("frozen_info", geip, nil),
				"polluted":                         utils.PathSearch("polluted", geip, false),
				"status":                           utils.PathSearch("status", geip, nil),
				"created_at":                       utils.PathSearch("created_at", geip, nil),
				"updated_at":                       utils.PathSearch("updated_at", geip, nil),
				"tags":                             utils.FlattenTagsToMap(utils.PathSearch("tags", geip, nil)),
				"global_connection_bandwidth_id":   utils.PathSearch("global_connection_bandwidth_info.gcb_id", geip, nil),
				"global_connection_bandwidth_type": utils.PathSearch("global_connection_bandwidth_info.gcb_type", geip, nil),
				"associate_instance_region":        utils.PathSearch("associate_instance_info.region", geip, nil),
				"associate_instance_id":            utils.PathSearch("associate_instance_info.instance_id", geip, nil),
				"associate_instance_type":          utils.PathSearch("associate_instance_info.instance_type", geip, nil),
			})
		}

		// `current_count` means the number of `global_eips` in this page, and the limit of page is `10`.
		currentCount := utils.PathSearch("page_info.current_count", getGEIPRespBody, float64(0))
		if currentCount.(float64) < pageLimit {
			break
		}
		currentTotal += len(geips)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("global_eips", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
