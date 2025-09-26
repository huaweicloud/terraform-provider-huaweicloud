package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/members
func DataSourceChannelMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelMembersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the channel members are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the channel members belong.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC channel to which the members belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the channel member to be queried for fuzzy matching.`,
			},
			"member_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the channel member group to be queried for fuzzy matching.`,
			},
			"member_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the channel member group to be queried.`,
			},
			"precise_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parameter name for exact matching to be queried.`,
			},

			// Attributes.
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        channelMemberSchema(),
				Description: `The list of the channel members that matched filter parameters.`,
			},
		},
	}
}

func channelMemberSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the channel member.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC channel.`,
			},
			"member_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the channel member group.`,
			},
			"member_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the channel member group.`,
			},
			"member_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of the channel member.`,
			},
			"ecs_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the ECS instance.`,
			},
			"ecs_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the ECS instance.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The port of the channel member.`,
			},
			"is_backup": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the channel member is a backup node.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the channel member.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The weight value of the channel member.`,
			},
			"health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The health status of the channel member.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the channel member was added to the VPC channel, in RFC3339 format.`,
			},
		},
	}
}

func buildChannelMembersQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("member_group_name"); ok {
		res = fmt.Sprintf("%s&member_group_name=%v", res, v)
	}
	if v, ok := d.GetOk("member_group_id"); ok {
		res = fmt.Sprintf("%s&member_group_id=%v", res, v)
	}
	if v, ok := d.GetOk("precise_search"); ok {
		res = fmt.Sprintf("%s&precise_search=%v", res, v)
	}

	return res
}

func listChannelMembers(client *golangsdk.ServiceClient, instanceId, vpcChannelId, memberGroupName string,
	d ...*schema.ResourceData) ([]interface{}, error) {
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
	if len(d) > 0 {
		listPathWithLimit += buildChannelMembersQueryParams(d[0])
	} else {
		listPathWithLimit = fmt.Sprintf(`%s&member_group_name=%s&precise_search=member_group_name`, listPathWithLimit, memberGroupName)
	}

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

func flattenMembers(members []interface{}) []interface{} {
	if len(members) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(members))
	for _, item := range members {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", item, nil),
			"vpc_channel_id":    utils.PathSearch("vpc_channel_id", item, nil),
			"member_group_name": utils.PathSearch("member_group_name", item, nil),
			"member_group_id":   utils.PathSearch("member_group_id", item, nil),
			"member_ip_address": utils.PathSearch("host", item, nil),
			"ecs_id":            utils.PathSearch("ecs_id", item, nil),
			"ecs_name":          utils.PathSearch("ecs_name", item, nil),
			"port":              utils.PathSearch("port", item, nil),
			"is_backup":         utils.PathSearch("is_backup", item, nil),
			"status":            utils.PathSearch("status", item, nil),
			"weight":            utils.PathSearch("weight", item, nil),
			"health_status":     utils.PathSearch("health_status", item, nil),
			"create_time":       utils.PathSearch("create_time", item, nil),
		})
	}

	return result
}

func dataSourceChannelMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		instanceId      = d.Get("instance_id").(string)
		vpcChannelId    = d.Get("vpc_channel_id").(string)
		memberGroupName = d.Get("member_group_name").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	members, err := listChannelMembers(client, instanceId, vpcChannelId, memberGroupName, d)
	if err != nil {
		return diag.Errorf("error querying channel members: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("members", flattenMembers(members)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
