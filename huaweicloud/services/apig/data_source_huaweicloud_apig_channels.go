package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels
func DataSourceChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the dedicated instance to which the channels belong.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC channel ID of the to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the channel to be queried.`,
			},
			"precise_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter name for exact matching to be queried.`,
			},
			"member_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the member group to be queried.`,
			},
			"member_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the member group to be queried.`,
			},
			"vpc_channels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All VPC channels that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC channel.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the VPC channel.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port of the backend server.`,
						},
						"balance_strategy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The distribution algorithm.`,
						},
						"member_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The member type of the VPC channel.`,
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The type of the VPC channel.`,
						},
						"member_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The parameter member groups of the VPC channels.`,
							Elem: &schema.Resource{
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
										Description: `The weight of the current member group.`,
									},
									"microservice_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The microservice version of the backend server group.`,
									},
									"microservice_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The microservice port of the backend server group.`,
									},
									"microservice_labels": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The microservice tags of the backend server group.`,
										Elem: &schema.Resource{
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
										},
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of channel, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildListChannelsParams(d *schema.ResourceData) string {
	res := ""
	if channelId, ok := d.GetOk("channel_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, channelId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name)
	}
	if preciseSearch, ok := d.GetOk("precise_search"); ok {
		res = fmt.Sprintf("%s&precise_search=%v", res, preciseSearch)
	}
	if memberGroupId, ok := d.GetOk("member_group_id"); ok {
		res = fmt.Sprintf("%s&member_group_id=%v", res, memberGroupId)
	}
	if memberGroupName, ok := d.GetOk("member_group_name"); ok {
		res = fmt.Sprintf("%s&member_group_name=%v", res, memberGroupName)
	}
	return res
}

func queryChannels(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/vpc-channels?limit=500"
		instanceId = d.Get("instance_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	queryParams := buildListChannelsParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving channels under specified "+
				"dedicated instance (%s): %s", instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		channels := utils.PathSearch("vpc_channels", respBody, make([]interface{}, 0)).([]interface{})
		if len(channels) < 1 {
			break
		}
		result = append(result, channels...)
		offset += len(channels)
	}
	return result, nil
}

func dataSourceChannelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	channels, err := queryChannels(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_channels", flattenChannels(channels)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenChannels(channels []interface{}) []interface{} {
	if len(channels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(channels))
	for _, channel := range channels {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", channel, nil),
			"name":             utils.PathSearch("name", channel, nil),
			"port":             utils.PathSearch("port", channel, nil),
			"balance_strategy": utils.PathSearch("balance_strategy", channel, nil),
			"member_type":      utils.PathSearch("member_type", channel, nil),
			"type":             utils.PathSearch("type", channel, nil),
			"member_group":     flattenChannelsMemeberGroup(utils.PathSearch("member_groups", channel, make([]interface{}, 0))),
			"created_at":       utils.PathSearch("create_time", channel, nil),
		})
	}
	return result
}

func flattenChannelsMemeberGroup(memberGroups interface{}) []map[string]interface{} {
	memberGroupsArray := memberGroups.([]interface{})
	result := make([]map[string]interface{}, len(memberGroupsArray))
	for i, memberGroup := range memberGroupsArray {
		microserviceLabels := utils.PathSearch("microservice_labels", memberGroup, make([]interface{}, 0))
		result[i] = map[string]interface{}{
			"id":                   utils.PathSearch("member_group_id", memberGroup, nil),
			"name":                 utils.PathSearch("member_group_name", memberGroup, nil),
			"description":          utils.PathSearch("member_group_remark", memberGroup, nil),
			"weight":               utils.PathSearch("member_group_weight", memberGroup, nil),
			"microservice_version": utils.PathSearch("microservice_version", memberGroup, nil),
			"microservice_port":    utils.PathSearch("microservice_port", memberGroup, nil),
			"microservice_labels":  flattenMemeberGroupMicroserviceLabels(microserviceLabels),
		}
	}

	return result
}

func flattenMemeberGroupMicroserviceLabels(microserviceLabels interface{}) []map[string]interface{} {
	microserviceLabelsArray := microserviceLabels.([]interface{})
	result := make([]map[string]interface{}, len(microserviceLabelsArray))
	for i, microserviceLabel := range microserviceLabelsArray {
		result[i] = map[string]interface{}{
			"name":  utils.PathSearch("label_name", microserviceLabel, nil),
			"value": utils.PathSearch("label_value", microserviceLabel, nil),
		}
	}

	return result
}
