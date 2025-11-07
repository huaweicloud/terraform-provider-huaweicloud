package hss

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

// @API HSS GET /v5/{project_id}/rasp/events
func DataSourceRaspEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRaspEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"app_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attack_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_src_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_stack": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_input_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_input_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_headers": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_params": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"req_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chk_probe": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chk_rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chk_rule_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exist_bug": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRaspEventsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?host_id=%v&start_time=%v&end_time=%v&limit=10", d.Get("host_id"), d.Get("start_time"), d.Get("end_time"))

	if v, ok := d.GetOk("app_type"); ok {
		queryParams = fmt.Sprintf("%s&app_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("attack_tag"); ok {
		queryParams = fmt.Sprintf("%s&attack_tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_status"); ok {
		queryParams = fmt.Sprintf("%s&protect_status=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceRaspEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/rasp/events"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRaspEventsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving application protection events: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenRaspEvents(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRaspEvents(events []interface{}) []interface{} {
	if len(events) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(events))
	for _, v := range events {
		rst = append(rst, map[string]interface{}{
			"host_name":          utils.PathSearch("host_name", v, nil),
			"private_ip":         utils.PathSearch("private_ip", v, nil),
			"event_name":         utils.PathSearch("event_name", v, nil),
			"severity":           utils.PathSearch("severity", v, nil),
			"req_src_ip":         utils.PathSearch("req_src_ip", v, nil),
			"app_stack":          utils.PathSearch("app_stack", v, nil),
			"attack_input_name":  utils.PathSearch("attack_input_name", v, nil),
			"attack_input_value": utils.PathSearch("attack_input_value", v, nil),
			"query_string":       utils.PathSearch("query_string", v, nil),
			"req_headers":        utils.PathSearch("req_headers", v, nil),
			"req_method":         utils.PathSearch("req_method", v, nil),
			"req_params":         utils.PathSearch("req_params", v, nil),
			"req_path":           utils.PathSearch("req_path", v, nil),
			"req_protocol":       utils.PathSearch("req_protocol", v, nil),
			"req_url":            utils.PathSearch("req_url", v, nil),
			"attack_tag":         utils.PathSearch("attack_tag", v, nil),
			"chk_probe":          utils.PathSearch("chk_probe", v, nil),
			"chk_rule":           utils.PathSearch("chk_rule", v, nil),
			"chk_rule_desc":      utils.PathSearch("chk_rule_desc", v, nil),
			"exist_bug":          utils.PathSearch("exist_bug", v, nil),
		})
	}

	return rst
}
