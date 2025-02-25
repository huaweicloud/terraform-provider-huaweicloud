package rabbitmq

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

const pageLimit = 10

// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/tasks
func DataSourceDmsRabbitMQBackgroundTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRabbitMQBackgroundTasksRead,

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
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time of task where the query starts. The format is YYYYMMDDHHmmss.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time of task where the query ends. The format is YYYYMMDDHHmmss.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the task list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task name.`,
						},
						"params": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task parameters.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task status.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the username.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDmsRabbitMQBackgroundTasksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return kafka.DataSourceDmsKafkaBackgroundTasksRead(ctx, d, meta)
}
