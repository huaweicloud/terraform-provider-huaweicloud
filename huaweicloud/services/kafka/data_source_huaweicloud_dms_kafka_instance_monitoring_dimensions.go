package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/ces-hierarchy
func DataSourceInstanceMonitoringDimensions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceMonitoringDimensionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the Kafka instance monitoring dimensions are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"dimensions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of monitoring dimensions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the monitoring dimension.`,
						},
						"metrics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of monitoring metric names.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"key_name": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of keys used for monitoring queries.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"dim_router": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of monitoring dimension routes.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"children": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of child dimensions.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the child dimension.`,
									},
									"metrics": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of monitoring metric names.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"key_name": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of keys used for monitoring queries.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"dim_router": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of monitoring dimension routes.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of instance IDs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
					},
				},
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the node.`,
						},
					},
				},
			},
			"queues": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of queues.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the topic.`,
						},
						"partitions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of partitions.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the partition.`,
									},
								},
							},
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of consumer groups.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the consumer group.`,
						},
						"queues": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of queues.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the topic.`,
									},
									"partitions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of partitions.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The name of the partition.`,
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

func dataSourceInstanceMonitoringDimensionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/ces-hierarchy"
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving Kafka instance (%s) monitoring dimensions: %s", instanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("dimensions", flattenInstanceMonitoringDimensions(utils.PathSearch("dimensions",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("instance_ids", flattenMonitoringInstanceIds(utils.PathSearch("instance_ids",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("nodes", flattenInstanceMonitoringNodes(utils.PathSearch("nodes",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("queues", flattenInstanceMonitoringQueues(utils.PathSearch("queues",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("groups", flattenInstanceMonitoringGroups(utils.PathSearch("groups",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceMonitoringDimensions(dimensions []interface{}) []interface{} {
	if len(dimensions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dimensions))
	for _, dimension := range dimensions {
		result = append(result, map[string]interface{}{
			"name":       utils.PathSearch("name", dimension, nil),
			"metrics":    utils.PathSearch("metrics", dimension, nil),
			"key_name":   utils.PathSearch("key_name", dimension, nil),
			"dim_router": utils.PathSearch("dim_router", dimension, nil),
			"children": flattenInstanceMonitoringDimensionChildren(utils.PathSearch("children",
				dimension, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenInstanceMonitoringDimensionChildren(children []interface{}) []interface{} {
	if len(children) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(children))
	for _, child := range children {
		result = append(result, map[string]interface{}{
			"name":       utils.PathSearch("name", child, nil),
			"metrics":    utils.PathSearch("metrics", child, nil),
			"key_name":   utils.PathSearch("key_name", child, nil),
			"dim_router": utils.PathSearch("dim_router", child, nil),
		})
	}
	return result
}

func flattenMonitoringInstanceIds(instanceIds []interface{}) []interface{} {
	if len(instanceIds) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(instanceIds))
	for _, instanceId := range instanceIds {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", instanceId, nil),
		})
	}
	return result
}

func flattenInstanceMonitoringNodes(nodes []interface{}) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", node, nil),
		})
	}
	return result
}

func flattenInstanceMonitoringQueues(queues []interface{}) []interface{} {
	if len(queues) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(queues))
	for _, queue := range queues {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", queue, nil),
			"partitions": flattenInstanceMonitoringQueuePartitions(utils.PathSearch("partitions",
				queue, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenInstanceMonitoringQueuePartitions(partitions []interface{}) []interface{} {
	if len(partitions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(partitions))
	for _, partition := range partitions {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", partition, nil),
		})
	}
	return result
}

func flattenInstanceMonitoringGroups(groups []interface{}) []interface{} {
	if len(groups) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", group, nil),
			"queues": flattenInstanceMonitoringQueues(utils.PathSearch("queues",
				group, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}
