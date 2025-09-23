package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var channelMemberGroupNonUpdatableParams = []string{"instance_id", "vpc_channel_id"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}
func ResourceChannelMemberGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelMemberGroupCreate,
		ReadContext:   resourceChannelMemberGroupRead,
		UpdateContext: resourceChannelMemberGroupUpdate,
		DeleteContext: resourceChannelMemberGroupDelete,

		CustomizeDiff: config.FlexibleForceNew(channelMemberGroupNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceChannelMemberGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the member group is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the member group belongs.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC channel.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the member group.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the member group.`,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Description:  `The weight value of the member group.`,
			},
			"microservice_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version of the member group.`,
			},
			"microservice_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The port number of the member group.`,
			},
			"microservice_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        channelMemberGroupMicroserviceLabelsSchema(),
				Description: `The microservice labels of the member group.`,
			},
			"reference_vpc_channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the reference load balance channel.`,
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the member group, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the member group, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func channelMemberGroupMicroserviceLabelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the microservice label.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of the microservice label.`,
			},
		},
	}
}

func buildChannelMemberGroupMicroserviceLabels(labels []interface{}) []map[string]interface{} {
	if len(labels) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"label_name":  utils.PathSearch("name", label, ""),
			"label_value": utils.PathSearch("value", label, ""),
		})
	}

	return result
}

func buildCreateChannelMemberGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"member_groups": []map[string]interface{}{
			{
				"member_group_name":        d.Get("name"),
				"member_group_remark":      utils.ValueIgnoreEmpty(d.Get("description")),
				"member_group_weight":      utils.ValueIgnoreEmpty(d.Get("weight")),
				"microservice_version":     utils.ValueIgnoreEmpty(d.Get("microservice_version")),
				"microservice_port":        utils.ValueIgnoreEmpty(d.Get("microservice_port")),
				"microservice_labels":      buildChannelMemberGroupMicroserviceLabels(d.Get("microservice_labels").([]interface{})),
				"reference_vpc_channel_id": utils.ValueIgnoreEmpty(d.Get("reference_vpc_channel_id")),
			},
		},
	}
}

func createChannelMemberGroup(client *golangsdk.ServiceClient, instanceId, vpcChannelId string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{vpc_channel_id}", vpcChannelId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateChannelMemberGroupBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceChannelMemberGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		vpcChannelId = d.Get("vpc_channel_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := createChannelMemberGroup(client, instanceId, vpcChannelId, d)
	if err != nil {
		return diag.Errorf("error creating member group: %s", err)
	}

	memberGroupId := utils.PathSearch("member_groups[0].member_group_id", respBody, "").(string)
	if memberGroupId == "" {
		return diag.Errorf("unable to find the member group ID from the API response")
	}
	d.SetId(memberGroupId)

	return resourceChannelMemberGroupRead(ctx, d, meta)
}

func flattenChannelMemberGroupMicroserviceLabels(labels []interface{}) []map[string]interface{} {
	if len(labels) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("label_name", label, ""),
			"value": utils.PathSearch("label_value", label, ""),
		})
	}

	return result
}

func GetChannelMemberGroupById(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberGroupId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{vpc_channel_id}", vpcChannelId)
	getPath = strings.ReplaceAll(getPath, "{member_group_id}", memberGroupId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceChannelMemberGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		vpcChannelId  = d.Get("vpc_channel_id").(string)
		memberGroupId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	memberGroup, err := GetChannelMemberGroupById(client, instanceId, vpcChannelId, memberGroupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying member group detail")
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("member_group_name", memberGroup, nil)),
		d.Set("description", utils.PathSearch("member_group_remark", memberGroup, nil)),
		d.Set("weight", utils.PathSearch("member_group_weight", memberGroup, nil)),
		d.Set("microservice_version", utils.PathSearch("microservice_version", memberGroup, nil)),
		d.Set("microservice_port", utils.PathSearch("microservice_port", memberGroup, nil)),
		d.Set("microservice_labels", flattenChannelMemberGroupMicroserviceLabels(
			utils.PathSearch("microservice_labels", memberGroup, make([]interface{}, 0)).([]interface{}))),
		d.Set("reference_vpc_channel_id", utils.PathSearch("reference_vpc_channel_id", memberGroup, nil)),
		d.Set("create_time", utils.PathSearch("create_time", memberGroup, nil)),
		d.Set("update_time", utils.PathSearch("update_time", memberGroup, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateChannelMemberGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"member_group_name":        d.Get("name"),
		"member_group_remark":      utils.ValueIgnoreEmpty(d.Get("description")),
		"member_group_weight":      utils.ValueIgnoreEmpty(d.Get("weight")),
		"microservice_version":     utils.ValueIgnoreEmpty(d.Get("microservice_version")),
		"microservice_port":        utils.ValueIgnoreEmpty(d.Get("microservice_port")),
		"microservice_labels":      buildChannelMemberGroupMicroserviceLabels(d.Get("microservice_labels").([]interface{})),
		"reference_vpc_channel_id": utils.ValueIgnoreEmpty(d.Get("reference_vpc_channel_id")),
	}
}

func updateChannelMemberGroup(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberGroupId string, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{vpc_channel_id}", vpcChannelId)
	updatePath = strings.ReplaceAll(updatePath, "{member_group_id}", memberGroupId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateChannelMemberGroupBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceChannelMemberGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		vpcChannelId  = d.Get("vpc_channel_id").(string)
		memberGroupId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = updateChannelMemberGroup(client, instanceId, vpcChannelId, memberGroupId, d)
	if err != nil {
		return diag.Errorf("error updating member group (%s): %s", memberGroupId, err)
	}

	return resourceChannelMemberGroupRead(ctx, d, meta)
}

func deleteChannelMemberGroup(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberGroupId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups/{member_group_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{vpc_channel_id}", vpcChannelId)
	deletePath = strings.ReplaceAll(deletePath, "{member_group_id}", memberGroupId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceChannelMemberGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		vpcChannelId  = d.Get("vpc_channel_id").(string)
		memberGroupId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = deleteChannelMemberGroup(client, instanceId, vpcChannelId, memberGroupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting member group (%s)",
			memberGroupId))
	}

	return nil
}

func resourceChannelMemberGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf(`invalid format specified for import ID, want '<instance_id>/<vpc_channel_id>/<id>',
		 but got '%s'`, importedId)
	}

	d.SetId(parts[2])
	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("vpc_channel_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
