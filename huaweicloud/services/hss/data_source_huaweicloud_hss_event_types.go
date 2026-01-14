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

// @API HSS GET /v5/{project_id}/event/event-type
func DataSourceEventTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attack_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"att_ck": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceEventTypesDataListSchema(),
			},
		},
	}
}

func dataSourceEventTypesDataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_type_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"event_type_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"event_num": {
							Type:     schema.TypeInt,
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

func buildEventTypesQueryParams(d *schema.ResourceData, epsId string) string {
	severityList := d.Get("severity_list").([]interface{})
	tagList := d.Get("tag_list").([]interface{})
	queryParams := fmt.Sprintf("?category=%v", d.Get("category"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_name"); ok {
		queryParams = fmt.Sprintf("%s&container_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("handle_status"); ok {
		queryParams = fmt.Sprintf("%s&handle_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if len(severityList) > 0 {
		for _, v := range severityList {
			queryParams = fmt.Sprintf("%s&severity_list=%v", queryParams, v)
		}
	}
	if v, ok := d.GetOk("attack_tag"); ok {
		queryParams = fmt.Sprintf("%s&attack_tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if len(tagList) > 0 {
		for _, v := range tagList {
			queryParams = fmt.Sprintf("%s&tag_list=%v", queryParams, v)
		}
	}
	if v, ok := d.GetOk("att_ck"); ok {
		queryParams = fmt.Sprintf("%s&att_ck=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_name"); ok {
		queryParams = fmt.Sprintf("%s&event_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/event/event-type"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEventTypesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS event types: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("data_list", flattenEventTypesDataList(
			utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventTypesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"event_type_name": utils.PathSearch("event_type_name", v, nil),
			"event_type_num":  utils.PathSearch("event_type_num", v, nil),
			"event_type_list": flattenEventTypeList(
				utils.PathSearch("event_type_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenEventTypeList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"event_type": utils.PathSearch("event_type", v, nil),
			"event_num":  utils.PathSearch("event_num", v, nil),
			"status":     utils.PathSearch("status", v, nil),
		})
	}

	return rst
}
