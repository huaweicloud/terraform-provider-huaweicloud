package eg

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/eg/v1/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceEventSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventSubscriptionCreate,
		ReadContext:   resourceEventSubscriptionRead,
		UpdateContext: resourceEventSubscriptionUpdate,
		DeleteContext: resourceEventSubscriptionDelete,

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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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
					},
				},
				Description: "The list of the event sources.",
			},
			"targets": {
				Type:     schema.TypeList,
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

func buildEventSourcesOpts(sources []interface{}) []interface{} {
	result := make([]interface{}, 0, len(sources))
	for _, val := range sources {
		source, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"provider_type": source["provider_type"],
			"name":          source["name"],
			"filter":        unmarshalEventSubscriptionParamsters("filter rule of event source", source["filter_rule"].(string)),
		})
	}
	return result
}

func buildEventTargetsOpts(targets []interface{}) []interface{} {
	result := make([]interface{}, 0, len(targets))
	for _, val := range targets {
		target, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		element := map[string]interface{}{
			"provider_type":                target["provider_type"],
			"name":                         target["name"],
			"connection_id":                target["connection_id"],
			target["detail_name"].(string): unmarshalEventSubscriptionParamsters("event target detail", target["detail"].(string)),
			"transform":                    unmarshalEventSubscriptionParamsters("transform of event target", target["transform"].(string)),
		}
		if queueRaw := target["dead_letter_queue"].(string); queueRaw != "" {
			element["dead_letter_queue"] = unmarshalEventSubscriptionParamsters("dead letter queue of event target", queueRaw)
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
			Sources:     buildEventSourcesOpts(d.Get("sources").([]interface{})),
			Targets:     buildEventTargetsOpts(d.Get("targets").([]interface{})),
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
			Sources:        buildEventSourcesOpts(d.Get("sources").([]interface{})),
			Targets:        buildEventTargetsOpts(d.Get("targets").([]interface{})),
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

	err = subscriptions.Delete(client, subscriptionId)
	if err != nil {
		return diag.Errorf("error deleting EG subscription (%s): %s", subscriptionId, err)
	}
	return nil
}
