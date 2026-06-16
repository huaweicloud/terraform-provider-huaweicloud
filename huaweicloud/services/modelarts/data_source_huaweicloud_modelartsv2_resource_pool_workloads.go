package modelarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/pools/{pool_id}/workloads
func DataSourceResourcePoolWorkloads() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcePoolWorkloadsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource pool workloads are located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource pool.`,
			},

			// Optional parameters.
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the workload to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the workload to be queried.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort field of the query result.`,
			},
			"ascend": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"sort"},
				Description:  `Whether to sort the query results in ascending order.`,
			},

			// Attributes.
			"workloads": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource pool workloads that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The API version of the workload.`,
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the workload resource.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the workload.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The namespace of the workload.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the workload.`,
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the job to which the workload belongs.`,
						},
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the workload.`,
						},
						"job_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the job to which the workload belongs.`,
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource flavor of the workload.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the workload.`,
						},
						"resource_requirement": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The resource requirement of the workload.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The CPU resource of the workload.`,
									},
									"memory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The memory resource of the workload.`,
									},
									"nvidia_gpu": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The GPU resource of the workload.`,
									},
									"huawei_ascend_snt3": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The Ascend Snt3 NPU resource of the workload.`,
									},
									"huawei_ascend_snt9": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The Ascend Snt9 NPU resource of the workload.`,
									},
								},
							},
						},
						"priority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The priority of the workload.`,
						},
						"running_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The running duration of the workload, in seconds.`,
						},
						"pending_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The pending duration of the workload, in seconds.`,
						},
						"pending_position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The pending position of the workload.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the workload.`,
						},
						"gvk": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The GVK (GroupVersionKind) information of the workload.`,
						},
						"host_ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The host IPs of the workload.`,
						},
						"nodes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The node information of the workload.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The host IP of the node.`,
									},
									"npu_topology_placement": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The NPU topology placement information of the node.`,
									},
									"resource_requirement": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The resource requirement of the node.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The CPU resource of the node.`,
												},
												"memory": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The memory resource of the node.`,
												},
												"nvidia_gpu": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The GPU resource of the node.`,
												},
												"huawei_ascend_snt310": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The Ascend Snt3 NPU resource of the node.`,
												},
												"huawei_ascend_snt1980": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The Ascend Snt9 NPU resource of the node.`,
												},
											},
										},
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

func buildResourcePoolWorkloadsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, v)
	}
	if v, ok := d.GetOk("ascend"); ok {
		res = fmt.Sprintf("%s&ascend=%v", res, v)
	}

	return res
}

func listResourcePoolWorkloads(client *golangsdk.ServiceClient, poolID string, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/pools/{pool_id}/workloads?limit={limit}"
		limit   = 500
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_id}", poolID)
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

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}

		offset += len(items)
	}

	return result, nil
}

func flattenNodesResourceRequirement(req map[string]interface{}) []map[string]interface{} {
	if len(req) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"cpu":                   utils.PathSearch("cpu", req, nil),
			"memory":                utils.PathSearch("memory", req, nil),
			"nvidia_gpu":            utils.PathSearch("nvidia.com/gpu", req, nil),
			"huawei_ascend_snt310":  utils.PathSearch("huawei.com/ascend-310", req, nil),
			"huawei_ascend_snt1980": utils.PathSearch("huawei.com/ascend-1980", req, nil),
		},
	}
}

func flattenResourcePoolWorkloadsNodes(nodes []interface{}) []map[string]interface{} {
	if len(nodes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"host_ip":                utils.PathSearch("hostIp", node, nil),
			"npu_topology_placement": utils.PathSearch("npuTopologyPlacement", node, nil),
			"resource_requirement": flattenNodesResourceRequirement(
				utils.PathSearch("resourceRequirement", node,
					make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func flattenResourcePoolWorkloadsResourceRequirement(req map[string]interface{}) []map[string]interface{} {
	if len(req) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"cpu":                utils.PathSearch("cpu", req, nil),
			"memory":             utils.PathSearch("memory", req, nil),
			"nvidia_gpu":         utils.PathSearch("nvidia.com/gpu", req, nil),
			"huawei_ascend_snt3": utils.PathSearch("huawei.com/ascend-snt3", req, nil),
			"huawei_ascend_snt9": utils.PathSearch("huawei.com/ascend-snt9", req, nil),
		},
	}
}

func flattenResourcePoolWorkloads(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"api_version":      utils.PathSearch("apiVersion", item, nil),
			"kind":             utils.PathSearch("kind", item, nil),
			"type":             utils.PathSearch("type", item, nil),
			"namespace":        utils.PathSearch("namespace", item, nil),
			"name":             utils.PathSearch("name", item, nil),
			"job_name":         utils.PathSearch("jobName", item, nil),
			"uid":              utils.PathSearch("uid", item, nil),
			"job_uuid":         utils.PathSearch("jobUUID", item, nil),
			"flavor":           utils.PathSearch("flavor", item, nil),
			"status":           utils.PathSearch("status", item, nil),
			"priority":         utils.PathSearch("priority", item, nil),
			"running_duration": utils.PathSearch("runningDuration", item, nil),
			"pending_duration": utils.PathSearch("pendingDuration", item, nil),
			"pending_position": utils.PathSearch("pendingPosition", item, nil),
			"gvk":              utils.PathSearch("gvk", item, nil),
			"host_ips":         utils.PathSearch("hostIps", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("createTime", item, float64(0)).(float64))/1000, false),
			"resource_requirement": flattenResourcePoolWorkloadsResourceRequirement(
				utils.PathSearch("resourceRequirement", item,
					make(map[string]interface{})).(map[string]interface{})),
			"nodes": flattenResourcePoolWorkloadsNodes(
				utils.PathSearch("nodes", item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceResourcePoolWorkloadsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	workloads, err := listResourcePoolWorkloads(client, d.Get("pool_id").(string),
		buildResourcePoolWorkloadsQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying resource pool workloads: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workloads", flattenResourcePoolWorkloads(workloads)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
