package dli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI GET /v2.0/{project_id}/datasource/enhanced-connections
func DataSourceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     connectionsSchema(),
			},
		},
	}
}

func connectionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"queues": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_privis": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"elastic_resource_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func convertTagsToString(tagRaw map[string]interface{}) string {
	if len(tagRaw) < 1 {
		return ""
	}

	tagList := make([]string, 0, len(tagRaw))
	for k, v := range tagRaw {
		tagList = append(tagList, k+"="+v.(string))
	}

	result := strings.Join(tagList, ",")
	return result
}

func buildListDataSourceConnectionsQueryParams(d *schema.ResourceData) string {
	queryParam := "status=ACTIVE"
	if v, ok := d.GetOk("name"); ok {
		queryParam += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagRaw := convertTagsToString(v.(map[string]interface{}))
		if tagRaw != "" {
			queryParam += fmt.Sprintf("&tags=%v", tagRaw)
		}
	}

	if queryParam != "" {
		queryParam = "?" + queryParam
	}

	return queryParam
}

func dataSourceConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2.0/{project_id}/datasource/enhanced-connections"
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI v2.0 client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	listDataSourceConnectionsQueryParams := buildListDataSourceConnectionsQueryParams(d)
	getPath += listDataSourceConnectionsQueryParams

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving datasource connections")
	}

	listRespJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("connections", flattenListConnections(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getFormatTime(key string, raw interface{}) string {
	return utils.FormatTimeStampRFC3339(int64(utils.PathSearch(key, raw, float64(0)).(float64))/1000, false)
}

func flattenListConnections(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("connections", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
			"queues":                 flattenQueues(utils.PathSearch("available_queue_info", v, make([]interface{}, 0))),
			"vpc_id":                 utils.PathSearch("dest_vpc_id", v, nil),
			"subnet_id":              utils.PathSearch("dest_network_id", v, nil),
			"is_privis":              utils.PathSearch("isPrivis", v, nil).(bool),
			"created_at":             getFormatTime("create_time", v),
			"hosts":                  flattenHosts(utils.PathSearch("hosts", v, make([]interface{}, 0))),
			"elastic_resource_pools": flattenResourcePools(utils.PathSearch("elastic_resource_pools", v, make([]interface{}, 0))),
			"routes":                 flattenRoutes(utils.PathSearch("routes", v, make([]interface{}, 0))),
		}
	}
	return rst
}

func getRaw(raw interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":         utils.PathSearch("peer_id", raw, nil),
		"status":     utils.PathSearch("status", raw, nil),
		"name":       utils.PathSearch("name", raw, nil),
		"error_msg":  utils.PathSearch("err_msg", raw, nil),
		"updated_at": getFormatTime("update_time", raw),
	}
}
func flattenQueues(queues interface{}) []map[string]interface{} {
	curArray := queues.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, queue := range curArray {
		result[i] = getRaw(queue)
	}
	return result
}

func flattenHosts(hosts interface{}) []map[string]interface{} {
	curArray := hosts.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, host := range curArray {
		result[i] = map[string]interface{}{
			"name": utils.PathSearch("name", host, nil),
			"ip":   utils.PathSearch("ip", host, nil),
		}
	}
	return result
}

func flattenResourcePools(resourcePools interface{}) []map[string]interface{} {
	curArray := resourcePools.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, resourcePool := range curArray {
		result[i] = getRaw(resourcePool)
	}
	return result
}

func flattenRoutes(routes interface{}) []map[string]interface{} {
	curArray := routes.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(curArray))
	for i, route := range curArray {
		result[i] = map[string]interface{}{
			"name":       utils.PathSearch("name", route, nil),
			"cidr":       utils.PathSearch("cidr", route, nil),
			"created_at": getFormatTime("create_time", route),
		}
	}
	return result
}
