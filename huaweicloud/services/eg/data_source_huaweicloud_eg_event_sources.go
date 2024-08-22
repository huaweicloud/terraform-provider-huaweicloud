package eg

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG GET /v1/{project_id}/sources
func DataSourceEventSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventSourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the event sources are located.`,
			},
			"provider_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the event sources to be queried.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the event channel to which the event sources belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the event source to be queried.`,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event source.`,
						},
						"channel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event channel to which the event source belong.`,
						},
						"channel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event channel to which the event source belong.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event source.`,
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name of the official event source.`,
						},
						"provider_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The provider type of the event source.`,
						},
						"event_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the event type.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The description of the event type.`,
									},
								},
							},
							Description: `The event types that official event source provided.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the event source.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the event source.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the event source.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the event source.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the event source.`,
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The message instance link information encapsulated in json format.`,
						},
					},
				},
				Description: `All event sources that match the filter parameters.`,
			},
		},
	}
}

func flattenSourceEventTypes(eventTypes []interface{}) []interface{} {
	result := make([]interface{}, 0, len(eventTypes))

	for _, eventType := range eventTypes {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", eventType, nil),
			"description": utils.PathSearch("description", eventType, nil),
		})
	}

	return result
}

func flattenDataEventSources(eventSources []interface{}) []interface{} {
	result := make([]interface{}, 0, len(eventSources))

	for _, eventSource := range eventSources {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", eventSource, nil),
			"channel_id":    utils.PathSearch("channel_id", eventSource, nil),
			"channel_name":  utils.PathSearch("channel_name", eventSource, nil),
			"name":          utils.PathSearch("name", eventSource, nil),
			"label":         utils.PathSearch("label", eventSource, nil),
			"provider_type": utils.PathSearch("provider_type", eventSource, nil),
			"event_types":   flattenSourceEventTypes(utils.PathSearch("event_types", eventSource, make([]interface{}, 0)).([]interface{})),
			"type":          utils.PathSearch("type", eventSource, nil),
			"description":   utils.PathSearch("description", eventSource, nil),
			"status":        utils.PathSearch("status", eventSource, nil),
			"created_at":    utils.PathSearch("created_time", eventSource, nil),
			"updated_at":    utils.PathSearch("updated_time", eventSource, nil),
			"detail":        parseCustomEventSourceDetail(utils.PathSearch("detail", eventSource, nil)),
		})
	}

	return result
}

func dataSourceEventSourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	eventSources, err := queryEventSources(client, d)
	if err != nil {
		return diag.Errorf("error querying event sources: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sources", flattenDataEventSources(eventSources)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG event sources: %s", err)
	}
	return nil
}
