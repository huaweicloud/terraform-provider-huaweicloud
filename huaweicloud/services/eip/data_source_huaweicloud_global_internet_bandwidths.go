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

// @API EIP GET /v3/{domain_id}/geip/internet-bandwidths
func DataSourceGlobalInternetBandwidths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalInternetBandwidthsRead,

		Schema: map[string]*schema.Schema{
			"bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_site": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
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
			"internet_bandwidths": {
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
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ingress_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ratio_95peak": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"frozen_info": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGlobalInternetBandwidthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getInternetBandwidthsHttpUrl := "v3/{domain_id}/geip/internet-bandwidths"
	getInternetBandwidthsPath := client.Endpoint + getInternetBandwidthsHttpUrl
	getInternetBandwidthsPath = strings.ReplaceAll(getInternetBandwidthsPath, "{domain_id}", cfg.DomainID)
	getInternetBandwidthsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getInternetBandwidthsPath += fmt.Sprintf("?limit=%v", pageLimit)
	if id, ok := d.GetOk("bandwidth_id"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&id=%s", id)
	}
	if size, ok := d.GetOk("size"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&size=%s", size)
	}
	if name, ok := d.GetOk("name"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&name=%s", name)
	}
	if accessSite, ok := d.GetOk("access_site"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&access_site=%s", accessSite)
	}
	if bwType, ok := d.GetOk("type"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&type=%s", bwType)
	}
	if status, ok := d.GetOk("status"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&status=%s", status)
	}
	if epsID, ok := d.GetOk("enterprise_project_id"); ok {
		getInternetBandwidthsPath += fmt.Sprintf("&enterprise_project_id=%s", epsID)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		tagsList := expandTagsMapToStringList(rawTags.(map[string]interface{}))
		for _, v := range tagsList {
			getInternetBandwidthsPath += fmt.Sprintf("&tags=%s", v)
		}
	}

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getInternetBandwidthsPath + fmt.Sprintf("&offset=%d", currentTotal)
		getInternetBandwidthsResp, err := client.Request("GET", currentPath, &getInternetBandwidthsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getInternetBandwidthsRespBody, err := utils.FlattenResponse(getInternetBandwidthsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		bws := utils.PathSearch("internet_bandwidths", getInternetBandwidthsRespBody, make([]interface{}, 0)).([]interface{})
		for _, bw := range bws {
			results = append(results, map[string]interface{}{
				"id":                    utils.PathSearch("id", bw, nil),
				"access_site":           utils.PathSearch("access_site", bw, nil),
				"isp":                   utils.PathSearch("isp", bw, nil),
				"charge_mode":           utils.PathSearch("charge_mode", bw, nil),
				"size":                  utils.PathSearch("size", bw, 0),
				"ingress_size":          utils.PathSearch("ingress_size", bw, 0),
				"description":           utils.PathSearch("description", bw, nil),
				"name":                  utils.PathSearch("name", bw, nil),
				"type":                  utils.PathSearch("type", bw, nil),
				"enterprise_project_id": utils.PathSearch("enterprise_project_id", bw, nil),
				"ratio_95peak":          utils.PathSearch("ratio_95peak", bw, 0),
				"frozen_info":           utils.PathSearch("frozen_info", bw, nil),
				"status":                utils.PathSearch("status", bw, nil),
				"created_at":            utils.PathSearch("created_at", bw, nil),
				"updated_at":            utils.PathSearch("updated_at", bw, nil),
				"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", bw, nil)),
			})
		}

		// `current_count` means the number of `internet_banwidths` in this page, and the limit of page is `10`.
		currentCount := utils.PathSearch("page_info.current_count", getInternetBandwidthsRespBody, float64(0))
		if currentCount.(float64) < pageLimit {
			break
		}
		currentTotal += len(bws)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("internet_bandwidths", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func expandTagsMapToStringList(tagMap map[string]interface{}) []string {
	if len(tagMap) < 1 {
		return nil
	}

	taglist := make([]string, 0)
	for k, v := range tagMap {
		tag := k + "|" + v.(string)
		taglist = append(taglist, tag)
	}

	return taglist
}
