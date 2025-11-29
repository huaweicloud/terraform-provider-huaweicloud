package eg

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

// @API EG GET /v1/{project_id}/connections
func DataSourceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the connections are located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The exact name of the connection to be queried.`,
			},
			"fuzzy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the connection to be queried for fuzzy matching.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting method for query results.`,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the connection.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the connection.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the connection.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC to which the connection belongs.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the subnet to which the connection belongs.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the connection.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the connection.`,
						},
						"kafka_detail": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the Kafka instance.`,
									},
									"connect_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The connection address of the Kafka instance.`,
									},
									"security_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The security protocol used for the connection.`,
									},
									"enable_sasl_ssl": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether SASL_SSL is enabled for the Kafka instance.`,
									},
									"user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The username of the Kafka instance.`,
									},
									"acks": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `The number of confirmation signals the producer
											needs to receive to consider the message sent successfully.`,
									},
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The connection address of Kafka instance.`,
									},
								},
							},
							Description: `The Kafka detail information for the connection.`,
						},
						"agency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user delegation name used by the private network connection.`,
						},
						"flavor": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the flavor.`,
									},
									"concurrency_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The concurrency type of the flavor.`,
									},
									"concurrency": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The concurrency value of the flavor.`,
									},
									"bandwidth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The bandwidth type of the flavor.`,
									},
								},
							},
							Description: `The flavor information of the connection.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the connection, in UTC format.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the connection, in UTC format.`,
						},
						"error_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error code.`,
									},
									"error_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The detailed error information.`,
									},
									"error_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error message.`,
									},
								},
							},
							Description: `The error information of the connection.`,
						},
					},
				},
				Description: `All connections that match the filter parameters.`,
			},
		},
	}
}

func dataSourceConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	connections, err := listConnections(client, d)
	if err != nil {
		return diag.Errorf("error getting EG connections: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connections", flattenConnections(connections)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listConnections(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/connections"
		limit   = 500
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildConnectionsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		connections := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, connections...)
		if len(connections) < limit {
			break
		}

		offset += len(connections)
	}

	return result, nil
}

func buildConnectionsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("fuzzy_name"); ok {
		res = fmt.Sprintf("%s&fuzzy_name=%v", res, v)
	}

	if v, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, v)
	}

	return res
}

func flattenConnectionKafkaDetail(kafkaDetail interface{}) []interface{} {
	if kafkaDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"instance_id":       utils.PathSearch("instance_id", kafkaDetail, nil),
			"connect_address":   utils.PathSearch("addr", kafkaDetail, nil),
			"security_protocol": utils.PathSearch("security_protocol", kafkaDetail, nil),
			"enable_sasl_ssl":   utils.PathSearch("enable_sasl_ssl", kafkaDetail, nil),
			"user_name":         utils.PathSearch("user_name", kafkaDetail, nil),
			"acks":              utils.PathSearch("acks", kafkaDetail, nil),
			"address":           utils.PathSearch("addr", kafkaDetail, nil),
		},
	}
}

func flattenConnectionFlavor(flavor interface{}) []interface{} {
	if flavor == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"name":             utils.PathSearch("name", flavor, nil),
			"concurrency_type": utils.PathSearch("concurrency_type", flavor, nil),
			"concurrency":      utils.PathSearch("concurrency", flavor, nil),
			"bandwidth_type":   utils.PathSearch("bandwidth_type", flavor, nil),
		},
	}
}

func flattenConnectionErrorInfo(errorInfo interface{}) []interface{} {
	if errorInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"error_code":   utils.PathSearch("error_code", errorInfo, nil),
			"error_detail": utils.PathSearch("error_detail", errorInfo, nil),
			"error_msg":    utils.PathSearch("error_msg", errorInfo, nil),
		},
	}
}

func flattenConnections(connections []interface{}) []map[string]interface{} {
	if len(connections) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(connections))
	for _, connection := range connections {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", connection, nil),
			"name":         utils.PathSearch("name", connection, nil),
			"description":  utils.PathSearch("description", connection, nil),
			"vpc_id":       utils.PathSearch("vpc_id", connection, nil),
			"subnet_id":    utils.PathSearch("subnet_id", connection, nil),
			"type":         utils.PathSearch("type", connection, nil),
			"status":       utils.PathSearch("status", connection, nil),
			"kafka_detail": flattenConnectionKafkaDetail(utils.PathSearch("kafka_detail", connection, nil)),
			"agency":       utils.PathSearch("agency", connection, nil),
			"flavor":       flattenConnectionFlavor(utils.PathSearch("flavor", connection, nil)),
			"created_time": utils.PathSearch("created_time", connection, nil),
			"updated_time": utils.PathSearch("updated_time", connection, nil),
			"error_info":   flattenConnectionErrorInfo(utils.PathSearch("error_info", connection, nil)),
		})
	}

	return result
}
