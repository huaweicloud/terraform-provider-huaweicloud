package eg

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

// @API EG GET /v1/{project_id}/subscriptions
func DataSourceEventSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the subscriptions are located.`,
			},

			// Optional parameters.
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the event channel to filter subscriptions.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The exact name of the subscription to be queried.`,
			},
			"fuzzy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the subscription to be queried for fuzzy matching.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the target connection to filter subscriptions.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting method for query results.`,
			},

			// Attributes.
			"subscriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the subscriptions that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the subscription.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the subscription.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the subscription.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the subscription.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the subscription.`,
						},
						"channel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event channel.`,
						},
						"channel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event channel.`,
						},
						"used": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The associated resources of the subscription.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the associated resource.`,
									},
									"owner": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The management tenant account to which the associated resource belongs.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The description of the associated resource.`,
									},
								},
							},
						},
						"sources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of subscription sources.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the subscription source.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the subscription source.`,
									},
									"provider_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The provider type of the subscription source.`,
									},
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter list of the subscription source, in JSON format.`,
									},
									"filter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The matching filter rules of the subscription source, in JSON format.`,
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creation time of the subscription source, in RFC3339 format.`,
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The update time of the subscription source, in RFC3339 format.`,
									},
								},
							},
						},
						"targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of subscription targets.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the subscription target.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the subscription target.`,
									},
									"provider_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The provider type of the subscription target.`,
									},
									"connection_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The target connection ID used by the subscription target.`,
									},
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter list of the subscription target, in JSON format.`,
									},
									"kafka_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The Kafka target parameter list of the subscription, in JSON format.`,
									},
									"smn_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The SMN target parameter list of the subscription, in JSON format.`,
									},
									"eg_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The EG channel target parameter list of the subscription, in JSON format.`,
									},
									"apigw_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The APIGW target parameter list of the subscription, in JSON format.`,
									},
									"retry_times": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of retry times.`,
									},
									"transform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The transform rules of the subscription target, in JSON format.`,
									},
									"dead_letter_queue": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dead letter queue parameters of the subscription, in JSON format.`,
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creation time of the subscription target, in RFC3339 format.`,
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The update time of the subscription target, in RFC3339 format.`,
									},
								},
							},
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the subscription, in RFC3339 format.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the subscription, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceEventSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	subscriptions, err := ListEventSubscriptions(client, d)
	if err != nil {
		return diag.Errorf("error querying event subscriptions: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("subscriptions", flattenEventSubscriptions(subscriptions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ListEventSubscriptions(client *golangsdk.ServiceClient, d ...*schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/subscriptions?limit={limit}"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	if len(d) > 0 {
		listPath += buildEventSubscriptionsQueryParams(d[0])
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		subscriptions := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, subscriptions...)
		if len(subscriptions) < limit {
			break
		}

		offset += len(subscriptions)
	}

	return result, nil
}

func buildEventSubscriptionsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("channel_id"); ok {
		res = fmt.Sprintf("%s&channel_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("fuzzy_name"); ok {
		res = fmt.Sprintf("%s&fuzzy_name=%v", res, v)
	}
	if v, ok := d.GetOk("connection_id"); ok {
		res = fmt.Sprintf("%s&connection_id=%v", res, v)
	}
	if v, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, v)
	}

	return res
}

func flattenSubscriptionUsed(usedList []interface{}) []map[string]interface{} {
	if len(usedList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(usedList))
	for _, used := range usedList {
		result = append(result, map[string]interface{}{
			"resource_id": utils.PathSearch("resource_id", used, nil),
			"owner":       utils.PathSearch("owner", used, nil),
			"description": utils.PathSearch("description", used, nil),
		})
	}

	return result
}

func flattenSubscriptionSources(sources []interface{}) []map[string]interface{} {
	if len(sources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sources))
	for _, source := range sources {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", source, nil),
			"name":          utils.PathSearch("name", source, nil),
			"provider_type": utils.PathSearch("provider_type", source, nil),
			"detail":        utils.JsonToString(utils.PathSearch("detail", source, nil)),
			"filter":        utils.JsonToString(utils.PathSearch("filter", source, nil)),
			"created_time":  utils.PathSearch("created_time", source, nil),
			"updated_time":  utils.PathSearch("updated_time", source, nil),
		})
	}

	return result
}

func flattenSubscriptionTargets(targets []interface{}) []map[string]interface{} {
	if len(targets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(targets))
	for _, target := range targets {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", target, nil),
			"name":              utils.PathSearch("name", target, nil),
			"provider_type":     utils.PathSearch("provider_type", target, nil),
			"connection_id":     utils.PathSearch("connection_id", target, nil),
			"detail":            utils.JsonToString(utils.PathSearch("detail", target, nil)),
			"kafka_detail":      utils.JsonToString(utils.PathSearch("kafka_detail", target, nil)),
			"smn_detail":        utils.JsonToString(utils.PathSearch("smn_detail", target, nil)),
			"eg_detail":         utils.JsonToString(utils.PathSearch("eg_detail", target, nil)),
			"apigw_detail":      utils.JsonToString(utils.PathSearch("apigw_detail", target, nil)),
			"retry_times":       utils.PathSearch("retry_times", target, nil),
			"transform":         utils.JsonToString(utils.PathSearch("transform", target, nil)),
			"dead_letter_queue": utils.JsonToString(utils.PathSearch("dead_letter_queue", target, nil)),
			"created_time":      utils.PathSearch("created_time", target, nil),
			"updated_time":      utils.PathSearch("updated_time", target, nil),
		})
	}

	return result
}

func flattenEventSubscriptions(subscriptions []interface{}) []map[string]interface{} {
	if len(subscriptions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", subscription, nil),
			"name":         utils.PathSearch("name", subscription, nil),
			"description":  utils.PathSearch("description", subscription, nil),
			"type":         utils.PathSearch("type", subscription, nil),
			"status":       utils.PathSearch("status", subscription, nil),
			"channel_id":   utils.PathSearch("channel_id", subscription, nil),
			"channel_name": utils.PathSearch("channel_name", subscription, nil),
			"used":         flattenSubscriptionUsed(utils.PathSearch("used", subscription, make([]interface{}, 0)).([]interface{})),
			"sources":      flattenSubscriptionSources(utils.PathSearch("sources", subscription, make([]interface{}, 0)).([]interface{})),
			"targets":      flattenSubscriptionTargets(utils.PathSearch("targets", subscription, make([]interface{}, 0)).([]interface{})),
			"created_time": utils.PathSearch("created_time", subscription, nil),
			"updated_time": utils.PathSearch("updated_time", subscription, nil),
		})
	}

	return result
}
