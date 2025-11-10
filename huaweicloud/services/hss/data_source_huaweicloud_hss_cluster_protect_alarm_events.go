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

// @API HSS GET /v5/{project_id}/cluster-protect/events
func DataSourceClusterProtectAlarmEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterProtectAlarmEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"event_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enforcement_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_name": {
										Type:     schema.TypeString,
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

func buildClusterProtectAlarmEventsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceClusterProtectAlarmEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "hss"
		epsId          = cfg.GetEnterpriseProjectID(d)
		result         = make([]interface{}, 0)
		offset         = 0
		totalNum       = 0
		lastUpdateTime = 0
		httpUrl        = "v5/{project_id}/cluster-protect/events"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildClusterProtectAlarmEventsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS cluster protect all alarm events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		lastUpdateTime = int(utils.PathSearch("last_update_time", respBody, float64(0)).(float64))
		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("last_update_time", lastUpdateTime),
		d.Set("data_list", flattenClusterProtectAlarmEventsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterProtectAlarmEventsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"action":                utils.PathSearch("action", v, nil),
			"cluster_name":          utils.PathSearch("cluster_name", v, nil),
			"cluster_id":            utils.PathSearch("cluster_id", v, nil),
			"event_name":            utils.PathSearch("event_name", v, nil),
			"event_class_id":        utils.PathSearch("event_class_id", v, nil),
			"event_id":              utils.PathSearch("event_id", v, nil),
			"event_type":            utils.PathSearch("event_type", v, nil),
			"event_content":         utils.PathSearch("event_content", v, nil),
			"handle_status":         utils.PathSearch("handle_status", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
			"update_time":           utils.PathSearch("update_time", v, nil),
			"project_id":            utils.PathSearch("project_id", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"policy_name":           utils.PathSearch("policy_name", v, nil),
			"policy_id":             utils.PathSearch("policy_id", v, nil),
			"resource_info":         flattenEventsResourceInfo(utils.PathSearch("resource_info", v, nil)),
		})
	}

	return rst
}

func flattenEventsResourceInfo(resourceInfo interface{}) []interface{} {
	if resourceInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enforcement_action": utils.PathSearch("enforcementAction", resourceInfo, nil),
			"group":              utils.PathSearch("group", resourceInfo, nil),
			"message":            utils.PathSearch("message", resourceInfo, nil),
			"name":               utils.PathSearch("name", resourceInfo, nil),
			"namespace":          utils.PathSearch("namespace", resourceInfo, nil),
			"version":            utils.PathSearch("version", resourceInfo, nil),
			"kind":               utils.PathSearch("kind", resourceInfo, nil),
			"resource_name":      utils.PathSearch("resource_name", resourceInfo, nil),
		},
	}
}
