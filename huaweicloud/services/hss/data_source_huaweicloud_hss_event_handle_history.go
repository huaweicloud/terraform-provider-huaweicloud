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

// @API HSS GET /v5/{project_id}/event/handle-history
func DataSourceEventHandleHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventHandleHistoryRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attack_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The field not take effect.
			"event_class_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_abstract": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attack_tag": {
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
						"asset_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"occur_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handle_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventHandleHistoryQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	eventClassIds := d.Get("event_class_ids").([]interface{})

	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("attack_tag"); ok {
		queryParams = fmt.Sprintf("%s&attack_tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("asset_value"); ok {
		queryParams = fmt.Sprintf("%s&asset_value=%v", queryParams, v)
	}
	if len(eventClassIds) > 0 {
		for _, v := range eventClassIds {
			queryParams = fmt.Sprintf("%s&event_class_ids=%v", queryParams, v)
		}
	}
	if v, ok := d.GetOk("event_name"); ok {
		queryParams = fmt.Sprintf("%s&event_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_type"); ok {
		queryParams = fmt.Sprintf("%s&event_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("handle_status"); ok {
		queryParams = fmt.Sprintf("%s&handle_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		queryParams = fmt.Sprintf("%s&public_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventHandleHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/event/handle-history"
		result  = make([]interface{}, 0)
		offset  = 0
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildEventHandleHistoryQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%d", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS event handle history: %s", err)
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
		d.Set("data_list", flattenEventHandleHistory(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventHandleHistory(dataResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"event_type":     utils.PathSearch("event_type", v, nil),
			"host_name":      utils.PathSearch("host_name", v, nil),
			"event_abstract": utils.PathSearch("event_abstract", v, nil),
			"attack_tag":     utils.PathSearch("attack_tag", v, nil),
			"private_ip":     utils.PathSearch("private_ip", v, nil),
			"public_ip":      utils.PathSearch("public_ip", v, nil),
			"asset_value":    utils.PathSearch("asset_value", v, nil),
			"occur_time":     utils.PathSearch("occur_time", v, nil),
			"handle_status":  utils.PathSearch("handle_status", v, nil),
			"notes":          utils.PathSearch("notes", v, nil),
			"event_class_id": utils.PathSearch("event_class_id", v, nil),
			"event_name":     utils.PathSearch("event_name", v, nil),
			"handle_time":    utils.PathSearch("handle_time", v, nil),
			"operate_type":   utils.PathSearch("operate_type", v, nil),
			"severity":       utils.PathSearch("severity", v, nil),
			"user_name":      utils.PathSearch("user_name", v, nil),
		})
	}
	return result
}
