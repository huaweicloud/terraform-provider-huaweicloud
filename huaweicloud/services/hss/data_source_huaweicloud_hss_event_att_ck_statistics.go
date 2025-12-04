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

// @API HSS GET /v5/{project_id}/event/att-ck
func DataSourceEventAttCkStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventAttCkStatisticsRead,
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
			"event_type": {
				Type:     schema.TypeInt,
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
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"att_ck": {
							Type:     schema.TypeString,
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
	}
}

func buildEventAttCkStatisticsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?category=%v", d.Get("category"))
	severityList := d.Get("severity_list").([]interface{})
	tagList := d.Get("tag_list").([]interface{})

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
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
	if v, ok := d.GetOk("event_type"); ok {
		queryParams = fmt.Sprintf("%s&event_type=%v", queryParams, v)
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
	if v, ok := d.GetOk("event_name"); ok {
		queryParams = fmt.Sprintf("%s&event_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventAttCkStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/event/att-ck"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildEventAttCkStatisticsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving statistics of ATT and CK phases: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenEventAttCkStatistics(dataResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventAttCkStatistics(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"att_ck": utils.PathSearch("att_ck", v, nil),
			"num":    utils.PathSearch("num", v, nil),
		})
	}

	return result
}
