package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/events
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/events/{name}
func DataSourceV2Events() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2EventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"event_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"involved_object": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventsInvolvedObjectSchema(),
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventsMetadataSchema(),
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reporting_component": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reporting_instance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"component": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func eventsInvolvedObjectSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func eventsMetadataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceV2EventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	events := make([]interface{}, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := GetEventDetail(client, namespace, name.(string))
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return diag.Errorf("error getting the event from the server: %s", err)
			}
		}
		events = append(events, resp)
	} else {
		resp, err := listEvents(client, namespace)
		if err != nil {
			return diag.Errorf("error finding the event list from the server: %s", err)
		}
		events = utils.PathSearch("items", resp, make([]interface{}, 0)).([]interface{})
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("events", flattenEvents(events)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetEventDetail(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getEventDetailHttpUrl := "apis/cci/v2/namespaces/{namespace}/events/{name}"
	getEventDetailPath := client.Endpoint + getEventDetailHttpUrl
	getEventDetailPath = strings.ReplaceAll(getEventDetailPath, "{namespace}", namespace)
	getEventDetailPath = strings.ReplaceAll(getEventDetailPath, "{name}", name)
	getEventDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getEventDetailResp, err := client.Request("GET", getEventDetailPath, &getEventDetailOpt)
	if err != nil {
		return getEventDetailResp, err
	}

	return utils.FlattenResponse(getEventDetailResp)
}

func listEvents(client *golangsdk.ServiceClient, namespace string) (interface{}, error) {
	listEventsHttpUrl := "apis/cci/v2/namespaces/{namespace}/events"
	listEventsPath := client.Endpoint + listEventsHttpUrl
	listEventsPath = strings.ReplaceAll(listEventsPath, "{namespace}", namespace)
	listEventsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listEventsResp, err := client.Request("GET", listEventsPath, &listEventsOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(listEventsResp)
}

func flattenEvents(events []interface{}) []interface{} {
	if len(events) == 0 {
		return nil
	}

	rst := make([]interface{}, len(events))
	for i, v := range events {
		rst[i] = map[string]interface{}{
			"action":              utils.PathSearch("action", v, nil),
			"api_version":         utils.PathSearch("apiVersion", v, nil),
			"count":               utils.PathSearch("count", v, nil),
			"event_time":          utils.PathSearch("eventTime", v, nil),
			"first_timestamp":     utils.PathSearch("firstTimestamp", v, nil),
			"kind":                utils.PathSearch("kind", v, nil),
			"last_timestamp":      utils.PathSearch("lastTimestamp", v, nil),
			"message":             utils.PathSearch("message", v, nil),
			"involved_object":     flattenInvolvedObject(utils.PathSearch("involvedObject", v, nil)),
			"metadata":            flattenMetadata(utils.PathSearch("metadata", v, nil)),
			"reason":              utils.PathSearch("reason", v, nil),
			"reporting_component": utils.PathSearch("reportingComponent", v, nil),
			"type":                utils.PathSearch("type", v, nil),
			"source":              flattenSource(utils.PathSearch("source", v, nil)),
		}
	}
	return rst
}

func flattenInvolvedObject(involvedObject interface{}) []map[string]interface{} {
	if involvedObject == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"field_path":       utils.PathSearch("fieldPath", involvedObject, nil),
			"kind":             utils.PathSearch("kind", involvedObject, nil),
			"name":             utils.PathSearch("name", involvedObject, nil),
			"namespace":        utils.PathSearch("namespace", involvedObject, nil),
			"resource_version": utils.PathSearch("resourceVersion", involvedObject, nil),
			"uid":              utils.PathSearch("uid", involvedObject, nil),
		},
	}
}

func flattenMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":               utils.PathSearch("name", metadata, nil),
			"namespace":          utils.PathSearch("namespace", metadata, nil),
			"annotations":        utils.PathSearch("annotations", metadata, nil),
			"creation_timestamp": utils.PathSearch("creationTimestamp", metadata, nil),
			"resource_version":   utils.PathSearch("resourceVersion", metadata, nil),
			"uid":                utils.PathSearch("uid", metadata, nil),
		},
	}
}

func flattenSource(source interface{}) []map[string]interface{} {
	if source == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"component": utils.PathSearch("component", source, nil),
			"host":      utils.PathSearch("host", source, nil),
		},
	}
}
