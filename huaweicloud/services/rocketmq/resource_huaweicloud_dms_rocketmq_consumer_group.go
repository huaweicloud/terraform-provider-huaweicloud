package rocketmq

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RocketMQ POST /v2/{project_id}/instances/{instance_id}/groups
// @API RocketMQ DELETE /v2/{project_id}/instances/{instance_id}/groups/{group}
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/groups/{group}
// @API RocketMQ PUT /v2/{project_id}/instances/{instance_id}/groups/{group}
func ResourceDmsRocketMQConsumerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQConsumerGroupCreate,
		UpdateContext: resourceDmsRocketMQConsumerGroupUpdate,
		ReadContext:   resourceDmsRocketMQConsumerGroupRead,
		DeleteContext: resourceDmsRocketMQConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the rocketMQ instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the consumer group.`,
			},
			"retry_max_times": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the maximum number of retry times.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Specifies the consumer group is enabled or not. Default to true.`,
			},
			"brokers": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the list of associated brokers of the consumer group.`,
			},
			"consume_orderly": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to consume orderly.`,
			},
			"broadcast": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to broadcast of the consumer group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the consumer group.`,
			},
		},
	}
}

func resourceDmsRocketMQConsumerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRocketmqConsumerGroup: create DMS rocketmq consumer group
	var (
		createRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups"
		createRocketmqConsumerGroupProduct = "dmsv2"
	)
	createRocketmqConsumerGroupClient, err := cfg.NewServiceClient(createRocketmqConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createRocketmqConsumerGroupPath := createRocketmqConsumerGroupClient.Endpoint + createRocketmqConsumerGroupHttpUrl
	createRocketmqConsumerGroupPath = strings.ReplaceAll(createRocketmqConsumerGroupPath, "{project_id}",
		createRocketmqConsumerGroupClient.ProjectID)
	createRocketmqConsumerGroupPath = strings.ReplaceAll(createRocketmqConsumerGroupPath, "{instance_id}", instanceID)

	createRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createRocketmqConsumerGroupOpt.JSONBody = utils.RemoveNil(buildCreateRocketmqConsumerGroupBodyParams(d, cfg))
	createRocketmqConsumerGroupResp, err := createRocketmqConsumerGroupClient.Request("POST",
		createRocketmqConsumerGroupPath, &createRocketmqConsumerGroupOpt)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQConsumerGroup: %s", err)
	}

	createRocketmqConsumerGroupRespBody, err := utils.FlattenResponse(createRocketmqConsumerGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	name := utils.PathSearch("name", createRocketmqConsumerGroupRespBody, "").(string)
	if name == "" {
		return diag.Errorf("unable to find consumer group name from the API response")
	}

	d.SetId(instanceID + "/" + name)

	return resourceDmsRocketMQConsumerGroupRead(ctx, d, meta)
}

func buildCreateRocketmqConsumerGroupBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enabled":         utils.ValueIgnoreEmpty(d.Get("enabled")),
		"broadcast":       utils.ValueIgnoreEmpty(d.Get("broadcast")),
		"brokers":         utils.ValueIgnoreEmpty(d.Get("brokers").(*schema.Set).List()),
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"retry_max_time":  utils.ValueIgnoreEmpty(d.Get("retry_max_times")),
		"group_desc":      utils.ValueIgnoreEmpty(d.Get("description")),
		"consume_orderly": utils.ValueIgnoreEmpty(d.Get("consume_orderly")),
	}
	return bodyParams
}

func resourceDmsRocketMQConsumerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRocketmqConsumerGroupHasChanges := []string{
		"enabled",
		"broadcast",
		"retry_max_times",
		"description",
		"consume_orderly",
	}

	if d.HasChanges(updateRocketmqConsumerGroupHasChanges...) {
		// updateRocketmqConsumerGroup: update DMS rocketmq consumer group
		var (
			updateRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
			updateRocketmqConsumerGroupProduct = "dmsv2"
		)
		updateRocketmqConsumerGroupClient, err := cfg.NewServiceClient(updateRocketmqConsumerGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating DMS client: %s", err)
		}

		parts := strings.SplitN(d.Id(), "/", 2)
		if len(parts) != 2 {
			return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
		}
		instanceID := parts[0]
		name := parts[1]
		updateRocketmqConsumerGroupPath := updateRocketmqConsumerGroupClient.Endpoint + updateRocketmqConsumerGroupHttpUrl
		updateRocketmqConsumerGroupPath = strings.ReplaceAll(updateRocketmqConsumerGroupPath, "{project_id}",
			updateRocketmqConsumerGroupClient.ProjectID)
		updateRocketmqConsumerGroupPath = strings.ReplaceAll(updateRocketmqConsumerGroupPath, "{instance_id}", instanceID)
		updateRocketmqConsumerGroupPath = strings.ReplaceAll(updateRocketmqConsumerGroupPath, "{group}", name)

		updateRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateRocketmqConsumerGroupOpt.JSONBody = utils.RemoveNil(buildUpdateRocketmqConsumerGroupBodyParams(d))
		_, err = updateRocketmqConsumerGroupClient.Request("PUT", updateRocketmqConsumerGroupPath,
			&updateRocketmqConsumerGroupOpt)
		if err != nil {
			return diag.Errorf("error updating DmsRocketMQConsumerGroup: %s", err)
		}
	}

	return resourceDmsRocketMQConsumerGroupRead(ctx, d, meta)
}

func buildUpdateRocketmqConsumerGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"broadcast":       utils.ValueIgnoreEmpty(d.Get("broadcast")),
		"retry_max_time":  utils.ValueIgnoreEmpty(d.Get("retry_max_times")),
		"group_desc":      utils.ValueIgnoreEmpty(d.Get("description")),
		"consume_orderly": utils.ValueIgnoreEmpty(d.Get("consume_orderly")),
		"enabled":         utils.ValueIgnoreEmpty(d.Get("enabled")),
	}
	return bodyParams
}

func resourceDmsRocketMQConsumerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqConsumerGroup: query DMS rocketmq consumer group
	var (
		getRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getRocketmqConsumerGroupProduct = "dmsv2"
	)
	getRocketmqConsumerGroupClient, err := cfg.NewServiceClient(getRocketmqConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRocketmqConsumerGroupPath := getRocketmqConsumerGroupClient.Endpoint + getRocketmqConsumerGroupHttpUrl
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{project_id}",
		getRocketmqConsumerGroupClient.ProjectID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{instance_id}", instanceID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{group}", name)

	getRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRocketmqConsumerGroupResp, err := getRocketmqConsumerGroupClient.Request("GET", getRocketmqConsumerGroupPath,
		&getRocketmqConsumerGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQConsumerGroup")
	}

	getRocketmqConsumerGroupRespBody, err := utils.FlattenResponse(getRocketmqConsumerGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("enabled", utils.PathSearch("enabled", getRocketmqConsumerGroupRespBody, nil)),
		d.Set("broadcast", utils.PathSearch("broadcast", getRocketmqConsumerGroupRespBody, nil)),
		d.Set("brokers", utils.PathSearch("brokers", getRocketmqConsumerGroupRespBody, nil)),
		d.Set("name", name),
		d.Set("retry_max_times", utils.PathSearch("retry_max_time", getRocketmqConsumerGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("group_desc", getRocketmqConsumerGroupRespBody, nil)),
		d.Set("consume_orderly", utils.PathSearch("consume_orderly", getRocketmqConsumerGroupRespBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQConsumerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRocketmqConsumerGroup: delete DMS rocketmq consumer group
	var (
		deleteRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		deleteRocketmqConsumerGroupProduct = "dmsv2"
	)
	deleteRocketmqConsumerGroupClient, err := cfg.NewServiceClient(deleteRocketmqConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	name := parts[1]
	deleteRocketmqConsumerGroupPath := deleteRocketmqConsumerGroupClient.Endpoint + deleteRocketmqConsumerGroupHttpUrl
	deleteRocketmqConsumerGroupPath = strings.ReplaceAll(deleteRocketmqConsumerGroupPath, "{project_id}",
		deleteRocketmqConsumerGroupClient.ProjectID)
	deleteRocketmqConsumerGroupPath = strings.ReplaceAll(deleteRocketmqConsumerGroupPath, "{instance_id}", instanceID)
	deleteRocketmqConsumerGroupPath = strings.ReplaceAll(deleteRocketmqConsumerGroupPath, "{group}", name)

	deleteRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteRocketmqConsumerGroupClient.Request("DELETE", deleteRocketmqConsumerGroupPath, &deleteRocketmqConsumerGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting DmsRocketMQConsumerGroup: %s", err)
	}

	d.SetId("")

	return nil
}
