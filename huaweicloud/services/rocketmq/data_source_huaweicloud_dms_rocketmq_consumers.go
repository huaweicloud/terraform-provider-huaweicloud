package rocketmq

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

const pageLimit = 10

// @API RocketMQ GET /v2/rocketmq/{project_id}/instances/{instance_id}/groups/{group}/clients
func DataSourceDmsRocketmqConsumers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketmqConsumersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the consumer group name.`,
			},
			"is_detail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to query the consumer details.`,
			},
			"clients": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of consumer subscription details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the subscription list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the name of the subscribed topic.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the subscription type.`,
									},
									"expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the subscription tag.`,
									},
								},
							},
						},
						"language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the client language.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the client version.`,
						},
						"client_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the client ID.`,
						},
						"client_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the client address.`,
						},
					},
				},
			},
			"online": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the consumer group is online.`,
			},
			"subscription_consistency": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether subscriptions are consistent.`,
			},
		},
	}
}

func dataSourceDmsRocketmqConsumersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listHttpUrl := "v2/rocketmq/{project_id}/instances/{instance_id}/groups/{group}/clients"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{group}", d.Get("group").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	listPath += fmt.Sprintf("?limit=%v", pageLimit)
	listPath += fmt.Sprintf("&is_detail=%v", d.Get("is_detail"))

	var offset int
	var online bool
	var subscriptionConsistency bool
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving consumers: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		consumers := utils.PathSearch("clients", listRespBody, make([]interface{}, 0)).([]interface{})
		for _, consumer := range consumers {
			results = append(results, map[string]interface{}{
				"language":       utils.PathSearch("language", consumer, nil),
				"version":        utils.PathSearch("version", consumer, nil),
				"client_id":      utils.PathSearch("client_id", consumer, nil),
				"client_address": utils.PathSearch("client_addr", consumer, nil),
				"subscriptions": flattenRocketMQConsumersClientSubscriptions(
					utils.PathSearch("subscriptions", consumer, make([]interface{}, 0)).([]interface{})),
			})
		}

		// `-1` means to the end
		nextOffset := utils.PathSearch("next_offset", listRespBody, float64(-1)).(float64)
		if int(nextOffset) == -1 {
			online = utils.PathSearch("online", listRespBody, false).(bool)
			subscriptionConsistency = utils.PathSearch("subscription_consistency", listRespBody, false).(bool)

			break
		}

		offset = int(nextOffset)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("clients", results),
		d.Set("online", online),
		d.Set("subscription_consistency", subscriptionConsistency),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRocketMQConsumersClientSubscriptions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"topic":      utils.PathSearch("topic", params, nil),
			"type":       utils.PathSearch("type", params, nil),
			"expression": utils.PathSearch("expression", params, nil),
		})
	}
	return rst
}
