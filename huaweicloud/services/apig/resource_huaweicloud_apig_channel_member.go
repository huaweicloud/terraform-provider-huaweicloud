package apig

import (
	"context"
	"fmt"
	"strconv"
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

var channelMemberNonUpdatableParams = []string{"instance_id", "vpc_channel_id"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members/{member_id}
func ResourceChannelMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelMemberCreate,
		ReadContext:   resourceChannelMemberRead,
		UpdateContext: resourceChannelMemberUpdate,
		DeleteContext: resourceChannelMemberDelete,

		CustomizeDiff: config.FlexibleForceNew(channelMemberNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceChannelMemberImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the channel member is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the channel member belongs.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC channel.`,
			},

			// Optional parameters.
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 10_000),
				Description:  `The weight value of the channel member.`,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 65_535),
				Description:  `The port number of the channel member.`,
			},
			"is_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether this member is the backup member.`,
			},
			"member_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the channel member group.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The status of the channel member.`,
			},
			"member_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IP address of channel member.`,
			},
			"ecs_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the ECS channel member.`,
			},
			"ecs_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the ECS channel member.`,
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the channel member, in RFC3339 format.`,
			},
			"member_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the member group.`,
			},
			"health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The health status of the channel member.`,
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

func buildChannelMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"members": []map[string]interface{}{
			{
				"weight":            utils.ValueIgnoreEmpty(d.Get("weight")),
				"port":              utils.ValueIgnoreEmpty(d.Get("port")),
				"is_backup":         utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_backup"),
				"member_group_name": utils.ValueIgnoreEmpty(d.Get("member_group_name")),
				"status":            utils.ValueIgnoreEmpty(d.Get("status")),
				"host":              utils.ValueIgnoreEmpty(d.Get("member_ip_address")),
				"ecs_id":            utils.ValueIgnoreEmpty(d.Get("ecs_id")),
				"ecs_name":          utils.ValueIgnoreEmpty(d.Get("ecs_name")),
			},
		},
	}
}

func createChannelMember(client *golangsdk.ServiceClient, instanceId, vpcChannelId string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{vpc_channel_id}", vpcChannelId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(utils.RemoveNil(buildChannelMemberBodyParams(d))),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceChannelMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	respBody, err := createChannelMember(client, instanceId, vpcChannelId, d)
	if err != nil {
		return diag.Errorf("error creating channel member: %s", err)
	}

	memberId := utils.PathSearch("members[0].id", respBody, "").(string)
	if memberId == "" {
		return diag.Errorf("unable to find the member ID from the API response")
	}
	d.SetId(memberId)

	return resourceChannelMemberRead(ctx, d, meta)
}

func listChannelMembers(client *golangsdk.ServiceClient, instanceId, vpcChannelId string) ([]interface{}, error) {
	var (
		result  = make([]interface{}, 0)
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members?limit={limit}"
		limit   = 100
		offset  = 0
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{instance_id}", instanceId)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{vpc_channel_id}", vpcChannelId)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPathWithLimit + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		members := utils.PathSearch("members", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, members...)
		if len(members) < limit {
			break
		}
		offset += len(members)
	}

	return result, nil
}

func GetChannelMemberById(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberId string) (interface{}, error) {
	members, err := listChannelMembers(client, instanceId, vpcChannelId)
	if err != nil {
		return nil, err
	}

	member := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", memberId), members, nil)
	if member == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return member, nil
}

func resourceChannelMemberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		vpcChannelId = d.Get("vpc_channel_id").(string)
		memberId     = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	member, err := GetChannelMemberById(client, instanceId, vpcChannelId, memberId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying channel member detail")
	}

	mErr := multierror.Append(
		d.Set("weight", utils.PathSearch("weight", member, nil)),
		d.Set("port", utils.PathSearch("port", member, nil)),
		d.Set("is_backup", utils.PathSearch("is_backup", member, nil)),
		d.Set("member_group_name", utils.PathSearch("member_group_name", member, nil)),
		d.Set("status", utils.PathSearch("status", member, nil)),
		d.Set("member_ip_address", utils.PathSearch("host", member, nil)),
		d.Set("ecs_id", utils.PathSearch("ecs_id", member, nil)),
		d.Set("ecs_name", utils.PathSearch("ecs_name", member, nil)),
		d.Set("create_time", utils.PathSearch("create_time", member, nil)),
		d.Set("member_group_id", utils.PathSearch("member_group_id", member, nil)),
		d.Set("health_status", utils.PathSearch("health_status", member, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateChannelMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"member_group_name": utils.ValueIgnoreEmpty(d.Get("member_group_name")),
		"members": []map[string]interface{}{
			{
				"port":              utils.ValueIgnoreEmpty(d.Get("port")),
				"weight":            utils.ValueIgnoreEmpty(d.Get("weight")),
				"is_backup":         utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_backup"),
				"member_group_name": utils.ValueIgnoreEmpty(d.Get("member_group_name")),
				"status":            utils.ValueIgnoreEmpty(d.Get("status")),
				"host":              utils.ValueIgnoreEmpty(d.Get("member_ip_address")),
				"ecs_id":            utils.ValueIgnoreEmpty(d.Get("ecs_id")),
				"ecs_name":          utils.ValueIgnoreEmpty(d.Get("ecs_name")),
			},
		},
	}
}

func updateChannelMember(client *golangsdk.ServiceClient, instanceId, vpcChannelId string, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{vpc_channel_id}", vpcChannelId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateChannelMemberBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceChannelMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = updateChannelMember(client, instanceId, vpcChannelId, d)
	if err != nil {
		return diag.Errorf("error updating channel member (%s): %s", d.Id(), err)
	}

	return resourceChannelMemberRead(ctx, d, meta)
}

func deleteChannelMember(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members/{member_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{vpc_channel_id}", vpcChannelId)
	deletePath = strings.ReplaceAll(deletePath, "{member_id}", memberId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceChannelMemberDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		vpcChannelId = d.Get("vpc_channel_id").(string)
		memberId     = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = deleteChannelMember(client, instanceId, vpcChannelId, memberId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting channel member (%s)",
			memberId))
	}

	return nil
}

func resourceChannelMemberImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
