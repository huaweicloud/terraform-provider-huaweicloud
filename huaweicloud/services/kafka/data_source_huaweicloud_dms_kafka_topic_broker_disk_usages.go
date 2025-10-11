package kafka

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

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/topics/diskusage
func DataSourceTopicBrokerDiskUsages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopicBrokerDiskUsagesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the topic broker disk usages are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the kafka instance.`,
			},
			"min_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The minimum disk size threshold to be queried.`,
			},
			"top": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of top topics to be queried.`,
			},
			"percentage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The percentage threshold to be queried.`,
			},
			"disk_usages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the broker.`,
						},
						"data_disk_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The total disk capacity.`,
						},
						"data_disk_use": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The used disk capacity.`,
						},
						"data_disk_free": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The free disk capacity.`,
						},
						"data_disk_use_percentage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The disk usage percentage.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the broker.`,
						},
						"topics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of topic disk usage.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The size of the disk usage.`,
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the topic.`,
									},
									"topic_partition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The partition of the topic.`,
									},
									"percentage": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The percentage of the disk usage.`,
									},
								},
							},
						},
					},
				},
				Description: `The disk usage of topics on brokers that match the filter parameters.`,
			},
		},
	}
}

func buildTopicBrokerDiskUsagesQueryParams(d *schema.ResourceData) string {
	res := ""

	if minSize, ok := d.GetOk("min_size"); ok {
		res += fmt.Sprintf("&minSize=%s", minSize)
	}
	if top, ok := d.GetOk("top"); ok {
		res += fmt.Sprintf("&top=%d", top)
	}
	if percentage, ok := d.GetOk("percentage"); ok {
		res += fmt.Sprintf("&percentage=%s", percentage)
	}

	if res != "" {
		res = "?" + res
	}
	return res
}

func getTopicBrokerDiskUsages(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/topics/diskusage"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildTopicBrokerDiskUsagesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceTopicBrokerDiskUsagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := getTopicBrokerDiskUsages(client, d)
	if err != nil {
		return diag.Errorf("error getting topic broker disk usages: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("disk_usages", flattenTopicBrokerDiskUsages(utils.PathSearch("broker_list",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopicBrokerDiskUsages(brokers []interface{}) []interface{} {
	if len(brokers) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(brokers))
	for _, v := range brokers {
		rst = append(rst, map[string]interface{}{
			"broker_name":              utils.PathSearch("broker_name", v, nil),
			"data_disk_size":           utils.PathSearch("data_disk_size", v, nil),
			"data_disk_use":            utils.PathSearch("data_disk_use", v, nil),
			"data_disk_free":           utils.PathSearch("data_disk_free", v, nil),
			"data_disk_use_percentage": utils.PathSearch("data_disk_use_percentage", v, nil),
			"status":                   utils.PathSearch("status", v, nil),
			"topics": flattenTopicBrokerDiskUsagesBrokerTopics(utils.PathSearch("topic_list",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenTopicBrokerDiskUsagesBrokerTopics(topics []interface{}) []interface{} {
	if len(topics) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(topics))
	for _, v := range topics {
		rst = append(rst, map[string]interface{}{
			"size":            utils.PathSearch("size", v, nil),
			"topic_name":      utils.PathSearch("topic_name", v, nil),
			"topic_partition": utils.PathSearch("topic_partition", v, nil),
			"percentage":      utils.PathSearch("percentage", v, nil),
		})
	}

	return rst
}
