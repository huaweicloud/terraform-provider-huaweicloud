package dws

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}
func DataSourceWorkloadQueues() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceWorkloadQueuesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cluster ID to which the workload queue belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the workload queue.`,
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the cluster. Required if the cluster is a logical cluster.`,
			},
			"queues": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the workload queue.`,
						},
						"logical_cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The logical cluster name.`,
						},
						"short_query_optimize": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Short query acceleration switch.`,
						},
						"short_query_concurrency_num": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The concurrency of short queries in the workload queue.`,
						},
						"configuration": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource name.`,
									},
									"resource_value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The resource attribute value.`,
									},
									"value_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource attribute unit.`,
									},
								},
							},
							Description: `The configuration information for workload queue.`,
						},
					},
				},
				Description: `The list of the workload queues.`,
			},
		},
	}
}

func resourceWorkloadQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues"
		product = "dws"
	)

	getClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getPath := getClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))

	// If cluster is logical cluster should add parameter
	if logicalName, ok := d.GetOk("logical_cluster_name"); ok {
		getPath += "?logical_cluster_name=" + logicalName.(string)
	}
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	getResp, err := getClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	queuesListJson := utils.PathSearch("workload_queue_name_list", getRespBody, make([]interface{}, 0))
	queuesList := queuesListJson.([]interface{})

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("queues", filterQueues(queuesList, d, getClient)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterQueues(all []interface{}, d *schema.ResourceData, client *golangsdk.ServiceClient) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok {
			if !strings.Contains(v.(string), fmt.Sprint(param)) {
				continue
			}
		}

		rst = append(rst, getQueueDetail(v.(string), d, client))
	}
	return rst
}

func getQueueDetail(queueName string, d *schema.ResourceData, client *golangsdk.ServiceClient) interface{} {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getPath = strings.ReplaceAll(getPath, "{queue_name}", queueName)
	// If cluster is logical cluster should add parameter
	if logicalName, ok := d.GetOk("logical_cluster_name"); ok {
		getPath += "?logical_cluster_name=" + logicalName.(string)
	}
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		log.Printf("[WARN] failed to get the workload queue detail: %s", err)
		return nil
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		log.Printf("[WARN] failed to flatten the workload queue detail: %s", err)
		return nil
	}

	queueDetail := utils.PathSearch("workload_queue", getRespBody, nil)
	// When short query optimize is t means support short query acceleration.
	shorQueryStr := utils.PathSearch("short_query_optimize", queueDetail, "").(string)
	shorQuery := false
	if shorQueryStr == "t" {
		shorQuery = true
	}
	rst := map[string]interface{}{
		"name":                        utils.PathSearch("queue_name", queueDetail, ""),
		"logical_cluster_name":        utils.PathSearch("logical_cluster_name", queueDetail, ""),
		"short_query_optimize":        shorQuery,
		"short_query_concurrency_num": utils.PathSearch("short_query_concurrency_num", queueDetail, 0.0).(float64),
		"configuration": flattenConfiguration(utils.PathSearch("resource_item_list", queueDetail,
			make([]interface{}, 0)).([]interface{})),
	}

	return rst
}

func flattenConfiguration(queue []interface{}) []interface{} {
	if len(queue) == 0 {
		return nil
	}
	rst := make([]interface{}, len(queue))
	for i, queueDetail := range queue {
		rst[i] = map[string]interface{}{
			"resource_name":  utils.PathSearch("resource_name", queueDetail, ""),
			"resource_value": utils.PathSearch("resource_value", queueDetail, 0.0).(float64),
			"value_unit":     utils.PathSearch("value_unit", queueDetail, ""),
		}
	}
	return rst
}
