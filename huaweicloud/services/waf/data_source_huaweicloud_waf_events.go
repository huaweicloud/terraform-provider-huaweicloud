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

// @API WAF GET /v1/{project_id}/waf/event
func DataSourceWafAttackEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafAttackEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"recent": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"attacks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payload_location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"headers": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cookie": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"response_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"response_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAttackEventsQueryParams(d *schema.ResourceData, epsId string) string {
	res := "?pagesize=100"
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsId)
	}
	if v, ok := d.GetOk("recent"); ok {
		res = fmt.Sprintf("%s&recent=%v", res, v)
	}
	if v, ok := d.GetOk("from"); ok {
		res = fmt.Sprintf("%s&from=%v", res, v)
	}
	if v, ok := d.GetOk("to"); ok {
		res = fmt.Sprintf("%s&to=%v", res, v)
	}
	if v, ok := d.GetOk("attacks"); ok {
		attacks, ok := v.([]interface{})
		if ok {
			for _, attack := range attacks {
				res = fmt.Sprintf("%s&attacks=%v", res, attack)
			}
		}
	}
	if v, ok := d.GetOk("hosts"); ok {
		hosts, ok := v.([]interface{})
		if ok {
			for _, host := range hosts {
				res = fmt.Sprintf("%s&hosts=%v", res, host)
			}
		}
	}
	if v, ok := d.GetOk("sips"); ok {
		sips, ok := v.([]interface{})
		if ok {
			for _, sip := range sips {
				res = fmt.Sprintf("%s&sips=%v", res, sip)
			}
		}
	}

	return res
}

func dataSourceWafAttackEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/waf/event"
		epsId       = cfg.GetEnterpriseProjectID(d)
		result      = make([]interface{}, 0)
		currentPage = 1
	)
	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAttackEventsQueryParams(d, epsId)
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
			return diag.Errorf("error retrieving WAF events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		event := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(event) == 0 {
			break
		}

		result = append(result, event...)
		currentPage++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenAttackEvents(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttackEvents(rawEvents []interface{}) []interface{} {
	if len(rawEvents) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawEvents))
	for _, v := range rawEvents {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"time":             utils.PathSearch("time", v, nil),
			"policyid":         utils.PathSearch("policyid", v, nil),
			"sip":              utils.PathSearch("sip", v, nil),
			"host":             utils.PathSearch("host", v, nil),
			"url":              utils.PathSearch("url", v, nil),
			"attack":           utils.PathSearch("attack", v, nil),
			"rule":             utils.PathSearch("rule", v, nil),
			"payload":          utils.PathSearch("payload", v, nil),
			"payload_location": utils.PathSearch("payload_location", v, nil),
			"action":           utils.PathSearch("action", v, nil),
			"request_line":     utils.PathSearch("request_line", v, nil),
			"headers":          utils.PathSearch("headers", v, nil),
			"cookie":           utils.PathSearch("cookie", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"process_time":     utils.PathSearch("process_time", v, nil),
			"region":           utils.PathSearch("region", v, nil),
			"host_id":          utils.PathSearch("host_id", v, nil),
			"response_time":    utils.PathSearch("response_time", v, nil),
			"response_size":    utils.PathSearch("response_size", v, nil),
			"response_body":    utils.PathSearch("response_body", v, nil),
			"request_body":     utils.PathSearch("request_body", v, nil),
		})
	}

	return rst
}
