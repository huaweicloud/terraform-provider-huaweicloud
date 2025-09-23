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

// @API WAF GET /v1/{projectid}/waf/rule/whiteblackip
func DataSourceAllWhiteblackipRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAllWhiteblackipRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policyids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
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
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"white": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ip_group": {
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
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"time_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"terminal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAllWhiteblackipRulesQueryParams(d *schema.ResourceData, epsId string) string {
	res := "?pagesize=1000"
	if v, ok := d.GetOk("policyids"); ok {
		res = fmt.Sprintf("%s&policyids=%v", res, v)
	}
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	return res
}

func dataSourceAllWhiteblackipRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		mErr        *multierror.Error
		httpUrl     = "v1/{projectid}/waf/rule/whiteblackip"
		allRules    []interface{}
		currentPage = 1
		epsId       = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{projectid}", client.ProjectID)
	requestPath += buildAllWhiteblackipRulesQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithPage := fmt.Sprintf("%s&page=%d", requestPath, currentPage)
		resp, err := client.Request("GET", requestPathWithPage, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving WAF all whiteblackip rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		rulesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(rulesResp) == 0 {
			break
		}

		allRules = append(allRules, rulesResp...)
		currentPage++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("items", flattenAllWhiteblackipRules(allRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAllWhiteblackipRules(rules []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"policyid":    utils.PathSearch("policyid", v, nil),
			"timestamp":   utils.PathSearch("timestamp", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"addr":        utils.PathSearch("addr", v, nil),
			"white":       utils.PathSearch("white", v, nil),
			"ip_group":    flattenWhiteblackipRulesIpGroup(utils.PathSearch("ip_group", v, nil)),
			"time_mode":   utils.PathSearch("time_mode", v, nil),
			"start":       utils.PathSearch("start", v, nil),
			"terminal":    utils.PathSearch("terminal", v, nil),
		})
	}
	return rst
}

func flattenWhiteblackipRulesIpGroup(ipGroup interface{}) []interface{} {
	if ipGroup == nil {
		return nil
	}

	rst := map[string]interface{}{
		"id":   utils.PathSearch("id", ipGroup, nil),
		"name": utils.PathSearch("name", ipGroup, nil),
		"size": utils.PathSearch("size", ipGroup, nil),
	}

	return []interface{}{rst}
}
