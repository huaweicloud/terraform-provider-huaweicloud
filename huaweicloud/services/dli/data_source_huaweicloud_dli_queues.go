package dli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI GET /v1.0/{project_id}/queues
// @API DLI GET /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}/queues
func DataSourceQueues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQueuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the queues.`,
			},
			"elastic_resource_pool_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"queue_type", "with_privilege", "with_charge_info", "tags"},
				Description:   `The name of the elastic resource pool.`,
			},
			"queue_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"elastic_resource_pool_name"},
				Description:  `The name of the queue to be queried.`,
			},
			"queue_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the queue to be queried.`,
			},
			"with_privilege": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to return permission information.`,
			},
			"with_charge_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to return charge information.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tags to filter queues.`,
			},

			// Attributes
			"queues": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of queues that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the queue.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the queue.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the queue.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the queue.`,
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The engine type of the queue.`,
						},
						"scaling_policies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The scaling policies of the queue.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The priority of the scaling policy.`,
									},
									"impact_start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The start time of the scaling policy.`,
									},
									"impact_stop_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The stop time of the scaling policy.`,
									},
									"min_cu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The minimum CU of the scaling policy.`,
									},
									"max_cu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The maximum CU of the scaling policy.`,
									},
									"inherit_elastic_resource_pool_max_cu": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether to inherit the maximum CU of the elastic resource pool.`,
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

func buildElasticResourcePoolQueuesQueryParams(d *schema.ResourceData) string {
	res := ""

	if queueName, ok := d.GetOk("queue_name"); ok {
		res = fmt.Sprintf("%s&queue_name=%v", res, queueName)
	}

	return res
}

func buildAllQueuesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("with_privilege"); ok {
		res = fmt.Sprintf("%s&with-priv=%v", res, v)
	}
	if v, ok := d.GetOk("with_charge_info"); ok {
		res = fmt.Sprintf("%s&with-charge-info=%v", res, v)
	}
	if v, ok := d.GetOk("queue_type"); ok {
		res = fmt.Sprintf("%s&queue_type=%v", res, v)
	}
	if v, ok := d.GetOk("tags"); ok {
		res = fmt.Sprintf("%s&tags=%v", res, v)
	}

	return res
}

func listElasticResourcePoolQueues(client *golangsdk.ServiceClient, elasticResourcePoolName string, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}/queues?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{elastic_resource_pool_name}", elasticResourcePoolName)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if queryParams != "" {
		listPath += queryParams
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		queues := utils.PathSearch("queues", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, queues...)
		if len(queues) < limit {
			break
		}

		offset += len(queues)
	}

	return result, nil
}

func listAllQueues(client *golangsdk.ServiceClient, queryParams string) ([]interface{}, error) {
	var (
		httpUrl     = "v1.0/{project_id}/queues?page-size={page-size}"
		pageSize    = 100
		currentPage = 1
		result      = make([]interface{}, 0)
		respBody    interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{page-size}", strconv.Itoa(pageSize))
	if queryParams != "" {
		listPath += queryParams
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPage := listPath + fmt.Sprintf("&current-page=%d", currentPage)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err = utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		queues := utils.PathSearch("queues", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, queues...)
		if len(queues) < pageSize {
			break
		}

		currentPage++
	}

	return result, nil
}

func flattenQueueScalingPolicies(policies []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"priority":                             utils.PathSearch("priority", policy, nil),
			"impact_start_time":                    utils.PathSearch("impact_start_time", policy, nil),
			"impact_stop_time":                     utils.PathSearch("impact_stop_time", policy, nil),
			"min_cu":                               utils.PathSearch("min_cu", policy, nil),
			"max_cu":                               utils.PathSearch("max_cu", policy, nil),
			"inherit_elastic_resource_pool_max_cu": utils.PathSearch("inherit_elastic_resource_pool_max_cu", policy, nil),
		})
	}

	return result
}

func flattenElasticResourcePoolQueues(queues []interface{}) []map[string]interface{} {
	if len(queues) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(queues))
	for _, queue := range queues {
		result = append(result, map[string]interface{}{
			"name":                  utils.PathSearch("queue_name", queue, nil),
			"type":                  utils.PathSearch("queue_type", queue, nil),
			"owner":                 utils.PathSearch("owner", queue, nil),
			"engine":                utils.PathSearch("engine", queue, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", queue, nil),
			"scaling_policies": flattenQueueScalingPolicies(
				utils.PathSearch("queue_scaling_policies", queue, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", queue, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	var queues []interface{}
	if elasticResourcePoolName, ok := d.GetOk("elastic_resource_pool_name"); ok {
		queues, err = listElasticResourcePoolQueues(client, elasticResourcePoolName.(string),
			buildElasticResourcePoolQueuesQueryParams(d))
		if err != nil {
			return diag.Errorf("error retrieving queues from elastic resource pool: %s", err)
		}
	} else {
		queues, err = listAllQueues(client, buildAllQueuesQueryParams(d))
		if err != nil {
			return diag.Errorf("error retrieving all queues: %s", err)
		}
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("queues", flattenElasticResourcePoolQueues(queues)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
