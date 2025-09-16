package waf

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

// @API WAF GET /v1/{project_id}/waf/overviews/response-code/timeline
func DataSourceOverviewsResponseCodeTimeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewsResponseCodeTimelineRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"response_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"response_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeline": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildOverviewsResponseCodeQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	hostIds := d.Get("hosts").([]interface{})
	instanceIds := d.Get("instances").([]interface{})
	rst := fmt.Sprintf("?from=%v&to=%v", d.Get("from"), d.Get("to"))

	if len(hostIds) > 0 {
		for _, v := range hostIds {
			rst += fmt.Sprintf("&hosts=%v", v)
		}
	}
	if len(instanceIds) > 0 {
		for _, v := range instanceIds {
			rst += fmt.Sprintf("&instances=%v", v)
		}
	}
	if responseSource, ok := d.GetOk("response_source"); ok {
		rst += fmt.Sprintf("&response_source=%v", responseSource)
	}
	if groupBy, ok := d.GetOk("group_by"); ok {
		rst += fmt.Sprintf("&group_by=%v", groupBy)
	}
	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}

	return rst
}

func dataSourceOverviewsResponseCodeTimelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/overviews/response-code/timeline"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOverviewsResponseCodeQueryParams(cfg, d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving WAF overviews response code timeline: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return diag.Errorf("convert response array failed")
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("response_codes", flattenOverviewsResponseCode(respArray)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOverviewsResponseCode(rawResponseCodes []interface{}) []interface{} {
	if len(rawResponseCodes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawResponseCodes))
	for _, v := range rawResponseCodes {
		rst = append(rst, map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"timeline": flattenOverviewsResponseCodeTimeline(utils.PathSearch("timeline", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenOverviewsResponseCodeTimeline(rawTimeline []interface{}) []interface{} {
	if len(rawTimeline) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawTimeline))
	for _, v := range rawTimeline {
		rst = append(rst, map[string]interface{}{
			"time": utils.PathSearch("time", v, nil),
			"num":  utils.PathSearch("num", v, nil),
		})
	}

	return rst
}
