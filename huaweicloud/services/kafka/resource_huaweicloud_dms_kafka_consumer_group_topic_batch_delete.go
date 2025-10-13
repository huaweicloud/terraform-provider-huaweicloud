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

var consumerGroupTopicBatchDeleteNonUpdatableParams = []string{
	"instance_id",
	"group",
	"topics",
}

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/delete-offset
func ResourceConsumerGroupTopicBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConsumerGroupTopicBatchDeleteCreate,
		ReadContext:   resourceConsumerGroupTopicBatchDeleteRead,
		UpdateContext: resourceConsumerGroupTopicBatchDeleteUpdate,
		DeleteContext: resourceConsumerGroupTopicBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(consumerGroupTopicBatchDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the topics to be deleted in batches are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the consumer group.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of topic names to be deleted.`,
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The result of the batch delete operation.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the topic.`,
						},
						"success": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the topic was deleted successfully.`,
						},
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error code if the topic deletion failed.`,
						},
					},
				},
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

func resourceConsumerGroupTopicBatchDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/delete-offset"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{group}", d.Get("group").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateConsumerGroupTopicBatchDeleteBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting consumer group topic offsets: %s", err)
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
		d.Set("result", flattenConsumerGroupTopicBatchDeleteResult(utils.PathSearch("topics",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateConsumerGroupTopicBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topics": utils.ValueIgnoreEmpty(d.Get("topics")),
	}
}

func flattenConsumerGroupTopicBatchDeleteResult(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rest := make([]interface{}, 0, len(resp))
	for _, topic := range resp {
		rest = append(rest, map[string]interface{}{
			"name":       utils.PathSearch("name", topic, nil),
			"success":    utils.PathSearch("success", topic, nil),
			"error_code": utils.PathSearch("error_code", topic, nil),
		})
	}
	return rest
}

func resourceConsumerGroupTopicBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConsumerGroupTopicBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceConsumerGroupTopicBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch deleting subscribed topics under the consumer group.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
