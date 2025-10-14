package kafka

import (
	"context"
	"fmt"
	"strconv"
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

// @API Kafka PUT /v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/reset-message-offset
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks
func ResourceMessageOffsetReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMessageOffsetResetCreate,
		ReadContext:   resourceMessageOffsetResetRead,
		DeleteContext: resourceMessageOffsetResetDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the consumption progress is to be reset is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the consumer group.`,
			},
			"partition": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The partition number.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the topic.`,
			},
			"message_offset": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"timestamp"},
				Description:  `The offset to reset the consumption progress.`,
			},
			"timestamp": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The time to reset the consumption progress. The value is a Unix timestamp, in millisecond`,
			},
		},
	}
}

func resourceMessageOffsetResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		createHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/reset-message-offset"
		consumerGroup = d.Get("group").(string)
	)

	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	params, err := buildCreateTMessageOffsetResetBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{group}", consumerGroup)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resetting message offset of the consumer group (%s): %s", consumerGroup, err)
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

func buildCreateTMessageOffsetResetBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		// when topic is empty string, reset all topic
		"topic":     d.Get("topic"),
		"partition": d.Get("partition"),
	}

	rawConfig := d.GetRawConfig()
	messageOffsetRaw := utils.GetNestedObjectFromRawConfig(rawConfig, "message_offset")
	if messageOffsetRaw != nil {
		messageOffset, err := strconv.Atoi(messageOffsetRaw.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting message offset to int: %s", err)
		}

		bodyParams["message_offset"] = messageOffset
	}

	timestampRaw := utils.GetNestedObjectFromRawConfig(rawConfig, "timestamp")
	if timestampRaw != nil {
		timestamp, err := strconv.Atoi(timestampRaw.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting timestamp to int: %s", err)
		}

		bodyParams["timestamp"] = timestamp
	}
	return bodyParams, nil
}

func resourceMessageOffsetResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMessageOffsetResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for resetting the consumption progress of the consumer group.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
