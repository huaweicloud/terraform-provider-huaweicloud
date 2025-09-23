package eg

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/eg/v1/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG POST /v1/{project_id}/subscriptions
// @API EG GET /v1/{project_id}/subscriptions/{subscription_id}
// @API EG PUT /v1/{project_id}/subscriptions/{subscription_id}
// @API EG DELETE /v1/{project_id}/subscriptions/{subscription_id}
func ResourceEventSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventSubscriptionCreate,
		ReadContext:   resourceEventSubscriptionRead,
		UpdateContext: resourceEventSubscriptionUpdate,
		DeleteContext: resourceEventSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the event subscription is located.",
			},
			"channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The channel ID to which the event subscription belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the event subscription.",
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The provider type of the event source.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the event source.",
						},
						"filter_rule": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The filter rule of the event source",
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								"The ID of the event source.",
								utils.SchemaDescInput{
									Required: true,
								}),
						},
						"detail_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name (key) of the source configuration detail.",
						},
						"detail": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The configuration source of the event target, in JSON format.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the event source.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the event source.",
						},
					},
				},
				Description: "The list of the event sources.",
			},
			"targets": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The provider type of the event target.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the event target.",
						},
						"detail_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name (key) of the target configuration detail.",
						},
						"detail": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The configuration detail of the event target, in JSON format.",
						},
						"transform": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The transform configuration of the event target, in JSON format.",
						},
						"connection_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The connection ID of the EG event target.",
						},
						"dead_letter_queue": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  "The specified queue to which failure events sent, in JSON format.",
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								"The ID of the event target.",
								utils.SchemaDescInput{
									Required: true,
								}),
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the event target.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the event target.",
						},
					},
				},
				Description: "The list of the event targets.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the event subscription.",
			},
			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the event subscription.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the event subscription.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the event subscription.",
			},
		},
	}
}

func unmarshalEventSubscriptionParamsters(paramName, paramVal string) map[string]interface{} {
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(paramVal), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the %s, not json format", paramName)
	}
	return parseResult
}

func buildEventSourcesOpts(newSources *schema.Set) []interface{} {
	result := make([]interface{}, 0, newSources.Len())
	for _, val := range newSources.List() {
		newSource, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		element := map[string]interface{}{
			"provider_type":                   newSource["provider_type"],
			"name":                            newSource["name"],
			"filter":                          unmarshalEventSubscriptionParamsters("filter rule of event source", newSource["filter_rule"].(string)),
			newSource["detail_name"].(string): unmarshalEventSubscriptionParamsters("event source detail", newSource["detail"].(string)),
		}
		if sourceId, ok := newSource["id"].(string); ok && sourceId != "" {
			// The ID can be omitted, a new source will be created in this scenario.
			element["id"] = sourceId
		}
		result = append(result, element)
	}
	return result
}

func buildEventTargetsOpts(newTargets *schema.Set) []interface{} {
	result := make([]interface{}, 0, newTargets.Len())
	for _, val := range newTargets.List() {
		newTarget, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		element := map[string]interface{}{
			"provider_type":                   newTarget["provider_type"],
			"name":                            newTarget["name"],
			newTarget["detail_name"].(string): unmarshalEventSubscriptionParamsters("event target detail", newTarget["detail"].(string)),
			"transform":                       unmarshalEventSubscriptionParamsters("transform of event target", newTarget["transform"].(string)),
			"connection_id":                   utils.ValueIgnoreEmpty(utils.PathSearch("connection_id", newTarget, nil)),
		}
		if queueRaw := newTarget["dead_letter_queue"].(string); queueRaw != "" {
			element["dead_letter_queue"] = unmarshalEventSubscriptionParamsters("dead letter queue of event target", queueRaw)
		}
		if targetId, ok := newTarget["id"].(string); ok && targetId != "" {
			// The ID can be omitted, a new source will be created in this scenario.
			element["id"] = targetId
		}
		result = append(result, element)
	}
	return result
}

func resourceEventSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		opts   = subscriptions.CreateOpts{
			ChannelId:   d.Get("channel_id").(string),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Sources:     buildEventSourcesOpts(d.Get("sources").(*schema.Set)),
			Targets:     buildEventTargetsOpts(d.Get("targets").(*schema.Set)),
		}
	)

	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := subscriptions.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating EG subscription: %s", err)
	}
	d.SetId(resp.ID)

	return resourceEventSubscriptionRead(ctx, d, meta)
}

func flattenEventSources(sourcesResp []map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(sourcesResp))
	for _, source := range sourcesResp {
		element := map[string]interface{}{
			"provider_type": source["provider_type"],
			"name":          source["name"],
			"id":            source["id"],
			"created_at":    source["created_time"],
			"updated_at":    source["updated_time"],
		}
		// find the key name of the target details
		for key, value := range source {
			if strings.Contains(key, "detail") {
				jsonDetail, err := json.Marshal(value)
				if err != nil {
					log.Printf("[ERROR] unable to convert the detail of the event source, not json format")
				} else {
					element["detail_name"] = key
					element["detail"] = string(jsonDetail)
				}
				break
			}
		}

		jsonFilter, err := json.Marshal(source["filter"])
		if err != nil {
			log.Printf("[ERROR] unable to convert the event source filter rule, not json format")
		} else {
			element["filter_rule"] = string(jsonFilter)
		}

		result = append(result, element)
	}
	return result
}

func flattenEventTargets(targetsResp []map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(targetsResp))
	for _, target := range targetsResp {
		element := map[string]interface{}{
			"provider_type": target["provider_type"],
			"name":          target["name"],
			"connection_id": target["connection_id"],
			"id":            target["id"],
			"created_at":    target["created_time"],
			"updated_at":    target["updated_time"],
		}
		// find the key name of the target details
		for key, value := range target {
			if strings.Contains(key, "detail") {
				jsonDetail, err := json.Marshal(value)
				if err != nil {
					log.Printf("[ERROR] unable to convert the detail of the event target, not json format")
				} else {
					element["detail_name"] = key
					element["detail"] = string(jsonDetail)
				}
				break
			}
		}

		if deadLetterQueue, ok := target["dead_letter_queue"]; ok {
			jsonQueue, err := json.Marshal(deadLetterQueue)
			if err != nil {
				log.Printf("[ERROR] unable to convert the dead letter queue of the event target, not json format")
			} else {
				element["dead_letter_queue"] = string(jsonQueue)
			}
		}

		if transform, ok := target["transform"]; ok {
			jsonTransform, err := json.Marshal(transform)
			if err != nil {
				log.Printf("[ERROR] unable to convert the transform configuration of the event target, not json format")
			} else {
				element["transform"] = string(jsonTransform)
			}
		}

		result = append(result, element)
	}
	return result
}

func resourceEventSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		subscriptionId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := subscriptions.Get(client, subscriptionId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EG subscription")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("channel_id", resp.ChannelId),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("sources", flattenEventSources(resp.Sources)),
		d.Set("targets", flattenEventTargets(resp.Targets)),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedTime),
		d.Set("updated_at", resp.UpdatedTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving EG event subscription (%s) fields: %s", subscriptionId, err)
	}
	return nil
}

func resourceEventSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		subscriptionId = d.Id()
		opts           = subscriptions.UpdateOpts{
			SubscriptionId: subscriptionId,
			Description:    utils.String(d.Get("description").(string)),
			Sources:        buildEventSourcesOpts(d.Get("sources").(*schema.Set)),
			Targets:        buildEventTargetsOpts(d.Get("targets").(*schema.Set)),
		}
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	_, err = subscriptions.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating EG subscription (%s): %s", subscriptionId, err)
	}
	return resourceEventSubscriptionRead(ctx, d, meta)
}

func DeleteEventSubscription(client *golangsdk.ServiceClient, subscriptionId string) error {
	return subscriptions.Delete(client, subscriptionId)
}

func resourceEventSubscriptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		subscriptionId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	err = DeleteEventSubscription(client, subscriptionId)
	if err != nil {
		return diag.Errorf("error deleting EG subscription (%s): %s", subscriptionId, err)
	}
	return nil
}
