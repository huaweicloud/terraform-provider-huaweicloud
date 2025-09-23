package kafka

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

// @API Kafka POST /v2/{project_id}/kafka/instances/{instance_id}/group
// @API Kafka PUT /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/groups/{group}
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/groups/batch-delete
func ResourceDmsKafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaConsumerGroupCreate,
		UpdateContext: resourceDmsKafkaConsumerGroupUpdate,
		ReadContext:   resourceDmsKafkaConsumerGroupRead,
		DeleteContext: resourceDmsKafkaConsumerGroupDelete,
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
				Description: `Specifies the ID of the Kafka instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the consumer group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the consumer group.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the state of the consumer group.`,
			},
			"coordinator_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the coordinator id of the consumer group.`,
			},
			"lag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the lag number of the consumer group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the consumer group.`,
			},
		},
	}
}

func resourceDmsKafkaConsumerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createKafkaConsumerGroup: create DMS Kafka consumer group
	var (
		createKafkaConsumerGroupHttpUrl = "v2/{project_id}/kafka/instances/{instance_id}/group"
		createKafkaConsumerGroupProduct = "dms"
	)
	createKafkaConsumerGroupClient, err := cfg.NewServiceClient(createKafkaConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createKafkaConsumerGroupPath := createKafkaConsumerGroupClient.Endpoint + createKafkaConsumerGroupHttpUrl
	createKafkaConsumerGroupPath = strings.ReplaceAll(createKafkaConsumerGroupPath, "{project_id}",
		createKafkaConsumerGroupClient.ProjectID)
	createKafkaConsumerGroupPath = strings.ReplaceAll(createKafkaConsumerGroupPath, "{instance_id}", instanceID)

	createKafkaConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createKafkaConsumerGroupOpt.JSONBody = utils.RemoveNil(buildCreateKafkaConsumerGroupBodyParams(d, cfg))
	_, createErr := createKafkaConsumerGroupClient.Request("POST",
		createKafkaConsumerGroupPath, &createKafkaConsumerGroupOpt)

	if createErr != nil {
		return diag.Errorf("error creating DMS Kafka consumer group: %s", createErr)
	}

	d.SetId(instanceID + "/" + d.Get("name").(string))
	return resourceDmsKafkaConsumerGroupRead(ctx, d, meta)
}

func buildCreateKafkaConsumerGroupBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name": d.Get("name"),
		"group_desc": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceDmsKafkaConsumerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateKafkaConsumerGroupHasChanges := []string{
		"description",
	}

	if d.HasChanges(updateKafkaConsumerGroupHasChanges...) {
		// updateKafkaConsumerGroup: update DMS Kafka consumer group
		var (
			updateKafkaConsumerGroupHttpUrl = "v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}"
			updateKafkaConsumerGroupProduct = "dms"
		)
		updateKafkaConsumerGroupClient, err := cfg.NewServiceClient(updateKafkaConsumerGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating DMS Client: %s", err)
		}

		parts := strings.Split(d.Id(), "/")
		if len(parts) != 2 {
			return diag.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
		}
		instanceID := parts[0]
		name := parts[1]
		updateKafkaConsumerGroupPath := updateKafkaConsumerGroupClient.Endpoint + updateKafkaConsumerGroupHttpUrl
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{engine}", "kafka")
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{project_id}",
			updateKafkaConsumerGroupClient.ProjectID)
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{instance_id}", instanceID)
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{group}", name)

		updateKafkaConsumerGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateKafkaConsumerGroupOpt.JSONBody = utils.RemoveNil(buildUpdateKafkaConsumerGroupBodyParams(d))
		_, err = updateKafkaConsumerGroupClient.Request("PUT", updateKafkaConsumerGroupPath,
			&updateKafkaConsumerGroupOpt)
		if err != nil {
			return diag.Errorf("error updating DMS Kafka consumer group: %s", err)
		}
	}

	return resourceDmsKafkaConsumerGroupRead(ctx, d, meta)
}

func buildUpdateKafkaConsumerGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name": d.Get("name"),
		"group_desc": d.Get("description"),
	}
	return bodyParams
}

func resourceDmsKafkaConsumerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getKafkaConsumerGroup: query DMS Kafka consumer group
	var (
		getKafkaConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getKafkaConsumerGroupProduct = "dms"
	)
	getKafkaConsumerGroupClient, err := cfg.NewServiceClient(getKafkaConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	name := parts[1]
	getKafkaConsumerGroupPath := getKafkaConsumerGroupClient.Endpoint + getKafkaConsumerGroupHttpUrl
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{project_id}",
		getKafkaConsumerGroupClient.ProjectID)
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{instance_id}", instanceID)
	getKafkaConsumerGroupPath = strings.ReplaceAll(getKafkaConsumerGroupPath, "{group}", name)

	getKafkaConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getKafkaConsumerGroupResp, err := getKafkaConsumerGroupClient.Request("GET", getKafkaConsumerGroupPath,
		&getKafkaConsumerGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS Kafka consumer group")
	}

	getKafkaConsumerGroupRespBody, err := utils.FlattenResponse(getKafkaConsumerGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupJson := utils.PathSearch("group", getKafkaConsumerGroupRespBody, nil)
	state := utils.PathSearch("state", groupJson, nil)
	if state == "DEAD" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", name),
		d.Set("description", utils.PathSearch("group_desc", groupJson, nil)),
		d.Set("state", utils.PathSearch("state", groupJson, nil)),
		d.Set("coordinator_id", utils.PathSearch("coordinator_id", groupJson, 0)),
		d.Set("lag", utils.PathSearch("lag", groupJson, 0)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			(int64(utils.PathSearch("createdAt", groupJson, float64(0)).(float64)))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsKafkaConsumerGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteKafkaConsumerGroup: delete DMS Kafka consumer group
	var (
		deleteKafkaConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/batch-delete"
		deleteKafkaConsumerGroupProduct = "dms"
	)
	deleteKafkaConsumerGroupClient, err := cfg.NewServiceClient(deleteKafkaConsumerGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	deleteKafkaConsumerGroupPath := deleteKafkaConsumerGroupClient.Endpoint + deleteKafkaConsumerGroupHttpUrl
	deleteKafkaConsumerGroupPath = strings.ReplaceAll(deleteKafkaConsumerGroupPath, "{project_id}",
		deleteKafkaConsumerGroupClient.ProjectID)
	deleteKafkaConsumerGroupPath = strings.ReplaceAll(deleteKafkaConsumerGroupPath, "{instance_id}", instanceID)

	deleteKafkaConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	deleteKafkaConsumerGroupOpt.JSONBody = utils.RemoveNil(buildBatchDeleteKafkaConsumerGroupBodyParams(d))
	_, err = deleteKafkaConsumerGroupClient.Request("POST", deleteKafkaConsumerGroupPath, &deleteKafkaConsumerGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting DMS Kafka consumer group: %s", err)
	}
	return resourceDmsKafkaConsumerGroupRead(ctx, d, meta)
}

func buildBatchDeleteKafkaConsumerGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_ids": []string{d.Get("name").(string)},
	}
	return bodyParams
}
