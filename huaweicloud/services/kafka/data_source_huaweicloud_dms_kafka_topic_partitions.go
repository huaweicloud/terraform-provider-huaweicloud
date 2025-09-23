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

// @API Kafka GET /v2/{project_id}/kafka/instances/{instance_id}/topics/{topic}/partitions
func DataSourceDmsKafkaTopicPartitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaTopicPartitionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic name.`,
			},
			"partitions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the partitions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the partition ID.`,
						},
						"start_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the start offset.`,
						},
						"last_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the last offset.`,
						},
						"last_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
						"message_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the message count.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDmsKafkaTopicPartitionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	getTopicPartitionsHttpUrl := "v2/{project_id}/kafka/instances/{instance_id}/topics/{topic}/partitions"
	getTopicPartitionsPath := client.Endpoint + getTopicPartitionsHttpUrl
	getTopicPartitionsPath = strings.ReplaceAll(getTopicPartitionsPath, "{project_id}", client.ProjectID)
	getTopicPartitionsPath = strings.ReplaceAll(getTopicPartitionsPath, "{instance_id}", d.Get("instance_id").(string))
	getTopicPartitionsPath = strings.ReplaceAll(getTopicPartitionsPath, "{topic}", d.Get("topic").(string))

	getTopicPartitionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pagelimit is `10`
	getTopicPartitionsPath += fmt.Sprintf("?limit=%v", pageLimit)
	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getTopicPartitionsPath + fmt.Sprintf("&offset=%d", currentTotal)
		getTopicPartitionsResp, err := client.Request("GET", currentPath, &getTopicPartitionsOpt)
		if err != nil {
			return diag.Errorf("error retrieving partitions: %s", err)
		}
		getTopicPartitionsRespBody, err := utils.FlattenResponse(getTopicPartitionsResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		partitions := utils.PathSearch("partitions", getTopicPartitionsRespBody, make([]interface{}, 0)).([]interface{})
		for _, partition := range partitions {
			results = append(results, map[string]interface{}{
				"partition":    utils.PathSearch("partition", partition, nil),
				"start_offset": utils.PathSearch("start_offset", partition, nil),
				"last_offset":  utils.PathSearch("last_offset", partition, nil),
				"last_update_time": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("last_update_time", partition, float64(0)).(float64))/1000, false),
				"message_count": utils.PathSearch("message_count", partition, nil),
			})
		}

		currentTotal += len(partitions)
		total := utils.PathSearch("total", getTopicPartitionsRespBody, float64(0)).(float64)
		if currentTotal == int(total) {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("partitions", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
