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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups
func DataSourceChannelMemberGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelMemberGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the list of member groups are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the VPC channel belongs.`,
			},
			"vpc_channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC channel to which the list of member groups belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the member group for fuzzy matching.`,
			},
			"precise_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The list of parameter names for exact matching.`,
			},

			// Attributes.
			"member_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        channelMemberGroupSchema(),
				Description: `The list of the member groups that matched filter parameters.`,
			},
		},
	}
}

func channelMemberGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the member group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the member group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the member group.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The weight value of the member group.`,
			},
			"microservice_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The microservice version of the member group.`,
			},
			"microservice_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The microservice port of the member group.`,
			},
			"microservice_labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        channelMemberGroupMicroserviceLabelSchema(),
				Description: `The microservice labels of the member group.`,
			},
			"reference_vpc_channel_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the referenced load channel.`,
			},
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
		},
	}
}

func channelMemberGroupMicroserviceLabelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the microservice label.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the microservice label.`,
			},
		},
	}
}

func buildChannelMemberGroupsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&member_group_name=%v", res, v)
	}
	if v, ok := d.GetOk("precise_search"); ok {
		res = fmt.Sprintf("%s&precise_search=%v", res, v)
	}

	return res
}

func listChannelMemberGroups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl      = "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}/member-groups?limit={limit}"
		instanceId   = d.Get("instance_id").(string)
		vpcChannelId = d.Get("vpc_channel_id").(string)
		limit        = 100
		offset       = 0
		result       = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{vpc_channel_id}", vpcChannelId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildChannelMemberGroupsQueryParams(d)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		memberGroups := utils.PathSearch("member_groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, memberGroups...)
		if len(memberGroups) < limit {
			break
		}
		offset += len(memberGroups)
	}

	return result, nil
}

func dataSourceChannelMemberGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	resp, err := listChannelMemberGroups(client, d)
	if err != nil {
		return diag.Errorf("error querying member groups: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("member_groups", flattenMemberGroups(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMemberGroups(memberGroups []interface{}) []interface{} {
	if len(memberGroups) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(memberGroups))
	for _, memberGroup := range memberGroups {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("member_group_id", memberGroup, nil),
			"name":                 utils.PathSearch("member_group_name", memberGroup, nil),
			"description":          utils.PathSearch("member_group_remark", memberGroup, nil),
			"weight":               utils.PathSearch("member_group_weight", memberGroup, nil),
			"microservice_version": utils.PathSearch("microservice_version", memberGroup, nil),
			"microservice_port":    utils.PathSearch("microservice_port", memberGroup, nil),
			"microservice_labels": flattenMemberGroupsMicroserviceLabels(
				utils.PathSearch("microservice_labels", memberGroup, make([]interface{}, 0)).([]interface{})),
			"reference_vpc_channel_id": utils.PathSearch("reference_vpc_channel_id", memberGroup, nil),
			"create_time":              utils.PathSearch("create_time", memberGroup, nil),
			"update_time":              utils.PathSearch("update_time", memberGroup, nil),
		})
	}

	return result
}

func flattenMemberGroupsMicroserviceLabels(microserviceLabels []interface{}) []interface{} {
	if len(microserviceLabels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(microserviceLabels))
	for _, microserviceLabel := range microserviceLabels {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("label_name", microserviceLabel, nil),
			"value": utils.PathSearch("label_value", microserviceLabel, nil),
		})
	}

	return result
}
