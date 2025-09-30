package kafka

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota
func DataSourceTopicQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopicQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the topic quotas are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the topic quotas belong.`,
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The keyword of the topic quota to be queried.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All topic quotas that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the topic.`,
						},
						"producer_byte_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The producer byte rate limit. The unit is B/s.`,
						},
						"consumer_byte_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The consumer byte rate limit. The unit is B/s.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceTopicQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	quotas, err := getTopicQuotas(client, instanceId, d.Get("keyword").(string))
	if err != nil {
		return diag.Errorf("error getting topic quotas of the instance (%s): %s", instanceId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenTopicQuotas(quotas.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopicQuotas(quotas []interface{}) []interface{} {
	if len(quotas) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(quotas))
	for _, quota := range quotas {
		result = append(result, map[string]interface{}{
			"topic":              utils.PathSearch("topic", quota, nil),
			"producer_byte_rate": utils.PathSearch(`"producer-byte-rate"`, quota, nil),
			"consumer_byte_rate": utils.PathSearch(`"consumer-byte-rate"`, quota, nil),
		})
	}
	return result
}
