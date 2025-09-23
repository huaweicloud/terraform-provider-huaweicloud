package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/reset-message-offset
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks
func ResourceDmsKafkaMessageOffsetReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaMessageOffsetResetCreate,
		ReadContext:   resourceDmsKafkaMessageOffsetResetRead,
		DeleteContext: resourceDmsKafkaMessageOffsetResetDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the instance ID.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the consumer group.`,
			},
			"partition": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the partition number.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the topic name.`,
			},
			"message_offset": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"timestamp"},
				Description:  `Specifies the message offset.`,
			},
			"timestamp": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the message offset. The value is a Unix timestamp, in millisecond`,
			},
		},
	}
}

func resourceDmsKafkaMessageOffsetResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/reset-message-offset"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{group}", d.Get("group").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateKafkaMessageOffsetResetBodyParams(d)),
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resetting Kafka message offset: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// when topic is empty string, reset all topic, and a job will be created
	if _, ok := d.GetOk("topic"); !ok {
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"CREATED"},
			Target:       []string{"SUCCESS"},
			Refresh:      FilterTaskRefreshFunc(client, d.Get("instance_id").(string), "kafkaResetConsumerOffset"),
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        5 * time.Second,
			PollInterval: 5 * time.Second,
		}

		task, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			if taskID := utils.PathSearch("id", task, ""); taskID != "" {
				return diag.Errorf("error waiting for job (%s) to be done: %s", taskID, err)
			}
			return diag.Errorf("error waiting for job to be done: %s", err)
		}
	}

	return nil
}

func buildCreateKafkaMessageOffsetResetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// when topic is empty string, reset all topic
		"topic":          d.Get("topic"),
		"partition":      d.Get("partition"),
		"message_offset": utils.ValueIgnoreEmpty(d.Get("message_offset")),
		"timestamp":      utils.ValueIgnoreEmpty(d.Get("timestamp")),
	}
	return bodyParams
}

func resourceDmsKafkaMessageOffsetResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsKafkaMessageOffsetResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
