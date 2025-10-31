package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var topicMessageBatchDeleteNonUpdatableParams = []string{
	"instance_id",
	"topic",
	"partitions",
}

// @API Kafka DELETE /v2/{project_id}/kafka/instances/{instance_id}/topics/{topic}/messages
func ResourceTopicMessageBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicMessageBatchDeleteCreate,
		ReadContext:   resourceTopicMessageBatchDeleteRead,
		UpdateContext: resourceTopicMessageBatchDeleteUpdate,
		DeleteContext: resourceTopicMessageBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(topicMessageBatchDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the topic messages to be deleted are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic.`,
			},
			"partitions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The number of the partition.`,
						},
						"offset": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The offset of the message to be deleted.`,
						},
					},
				},
				Description: `The partition configuration list to which the messages to be deleted belong.`,
			},
			// Attributes.
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of the partition.`,
						},
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The operation result.`,
						},
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error code if the operation failed.`,
						},
					},
				},
				Description: `The result of the message delete operation.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateTopicMessageBatchDeleteBodyParams(partitions []interface{}) map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(partitions))
	for _, partition := range partitions {
		result = append(result, map[string]interface{}{
			"partition": utils.PathSearch("partition", partition, nil),
			"offset":    utils.PathSearch("offset", partition, nil),
		})
	}

	return map[string]interface{}{
		"partitions": result,
	}
}

func resourceTopicMessageBatchDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		httpUrl   = "v2/{project_id}/kafka/instances/{instance_id}/topics/{topic}/messages"
		topicName = d.Get("topic").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{topic}", topicName)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateTopicMessageBatchDeleteBodyParams(d.Get("partitions").([]interface{}))),
	}

	resp, err := client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting topic (%s) messages: %s", topicName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("result", flattenTopicMessageBatchDeleteResult(utils.PathSearch("partitions",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopicMessageBatchDeleteResult(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rest := make([]interface{}, 0, len(resp))
	for _, partition := range resp {
		rest = append(rest, map[string]interface{}{
			"partition":  utils.PathSearch("partition", partition, nil),
			"result":     utils.PathSearch("result", partition, nil),
			"error_code": utils.PathSearch("error_code", partition, nil),
		})
	}
	return rest
}

func resourceTopicMessageBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTopicMessageBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTopicMessageBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for deleting topic messages in batches. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
