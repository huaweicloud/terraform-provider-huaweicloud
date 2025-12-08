package kafka

import (
	"context"
	"fmt"
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
// @API Kafka GET /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}
// @API Kafka PUT /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}
// @API Kafka DELETE /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}
func ResourceDmsKafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaConsumerGroupCreate,
		UpdateContext: resourceDmsKafkaConsumerGroupUpdate,
		ReadContext:   resourceDmsKafkaConsumerGroupRead,
		DeleteContext: resourceDmsKafkaConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceConsumerGroupImportState,
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
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the consumer group.`,
			},
			// Deprecated attribute(s).
			"lag": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: utils.SchemaDesc(
					`The lag number of the consumer group.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
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
		return diag.Errorf("error creating DMS client: %s", err)
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
			return diag.Errorf("error creating DMS client: %s", err)
		}

		updateKafkaConsumerGroupPath := updateKafkaConsumerGroupClient.Endpoint + updateKafkaConsumerGroupHttpUrl
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{engine}", "kafka")
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{project_id}",
			updateKafkaConsumerGroupClient.ProjectID)
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{instance_id}", d.Get("instance_id").(string))
		updateKafkaConsumerGroupPath = strings.ReplaceAll(updateKafkaConsumerGroupPath, "{group}", d.Get("name").(string))

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

func GetConsumerGroupByName(client *golangsdk.ServiceClient, instanceId, groupName string) (interface{}, error) {
	httpUrl := "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{group}", groupName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	consumerGroup := utils.PathSearch("group", respBody, nil)
	state := utils.PathSearch("state", consumerGroup, nil)
	// DEAD means the consumer group has been deleted.
	if state == "DEAD" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       getPath,
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the consumer group with name '%s' has been removed", groupName)),
			},
		}
	}

	return consumerGroup, nil
}

func resourceDmsKafkaConsumerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		groupName  = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	consumerGroup, err := GetConsumerGroupByName(client, instanceId, groupName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving consumer group (%s) under instance (%s)",
			groupName, instanceId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", groupName),
		d.Set("description", utils.PathSearch("group_desc", consumerGroup, nil)),
		d.Set("state", utils.PathSearch("state", consumerGroup, nil)),
		d.Set("coordinator_id", utils.PathSearch("coordinator_id", consumerGroup, 0)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			(int64(utils.PathSearch("createdAt", consumerGroup, float64(0)).(float64)))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsKafkaConsumerGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		httpUrl           = "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}"
		instanceId        = d.Get("instance_id").(string)
		consumerGroupName = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{group}", consumerGroupName)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// When the instance does not exist, the status code is 404.
	// DMS.111400864: The consumer group does not exist.
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DMS.111400864"),
			fmt.Sprintf("error deleting consumer group (%s) under Kafka instance (%s)", consumerGroupName, instanceId),
		)
	}

	return nil
}

func resourceConsumerGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<name>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
