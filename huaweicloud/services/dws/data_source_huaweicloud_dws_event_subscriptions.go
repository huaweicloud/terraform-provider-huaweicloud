package dws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/event-subs
func DataSourceEventSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceEventSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of event subscription.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of source event.`,
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The category of source event.`,
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The severity of source event.`,
			},
			"notification_target_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of notification target.`,
			},
			"enable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether the event subscription is enabled.`,
			},
			"event_subscriptions": {
				Type:        schema.TypeList,
				Elem:        eventSubSchema(),
				Computed:    true,
				Description: `The list of event subscriptions.`,
			},
		},
	}
}

func eventSubSchema() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of event subscription.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the event subscription.`,
			},
			"enable": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the event subscription is enabled.`,
			},
			"notification_target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The notification target.`,
			},
			"notification_target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of notification target.`,
			},
			"notification_target_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of notification target. Currently only **SMN** is supported.`,
			},
			"source_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of source event.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of source event.`,
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category of source event.`,
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The severity of source event.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time zone of the event subscription.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID of the event subscription.`,
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language of the event subscription.`,
			},
			"name_space": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name space of the event subscription.`,
			},
		},
	}

	return &nodeResource
}

func resourceEventSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getDwsEventSubs: Query the DWS event subscription.
	var (
		getDwsEventSubsHttpUrl = "v2/{project_id}/event-subs"
		getDwsEventSubsProduct = "dws"
	)
	getDwsEventSubsClient, err := cfg.NewServiceClient(getDwsEventSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getDwsEventSubsPath := getDwsEventSubsClient.Endpoint + getDwsEventSubsHttpUrl
	getDwsEventSubsPath = strings.ReplaceAll(getDwsEventSubsPath, "{project_id}", getDwsEventSubsClient.ProjectID)

	getDwsEventSubsResp, err := pagination.ListAllItems(
		getDwsEventSubsClient,
		"offset",
		getDwsEventSubsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("retrieving DWS event subscriptions: %s", err)
	}

	getDwsEventSubsRespJson, err := json.Marshal(getDwsEventSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDwsEventSubsRespBody interface{}
	err = json.Unmarshal(getDwsEventSubsRespJson, &getDwsEventSubsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	disasterList := utils.PathSearch("event_subscriptions", getDwsEventSubsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("event_subscriptions", filterEventSubs(flattenEventSubs(disasterList), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventSubs(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
			"source_id":                utils.PathSearch("source_id", v, nil),
			"source_type":              utils.PathSearch("source_type", v, nil),
			"category":                 utils.PathSearch("category", v, nil),
			"severity":                 utils.PathSearch("severity", v, nil),
			"enable":                   fmt.Sprint(utils.PathSearch("enable", v, nil)),
			"notification_target":      utils.PathSearch("notification_target", v, nil),
			"notification_target_name": utils.PathSearch("notification_target_name", v, nil),
			"notification_target_type": utils.PathSearch("notification_target_type", v, nil),
			"time_zone":                utils.PathSearch("time_zone", v, nil),
			"project_id":               utils.PathSearch("project_id", v, nil),
			"language":                 utils.PathSearch("language", v, nil),
			"name_space":               utils.PathSearch("name_space", v, nil),
		})
	}
	return rst
}

func filterEventSubs(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("notification_target_name"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("notification_target_name", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("enable"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("enable", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("source_type"); ok {
			if !isAllContain(fmt.Sprint(param), fmt.Sprint(utils.PathSearch("source_type", v, nil))) {
				continue
			}
		}
		if param, ok := d.GetOk("severity"); ok {
			if !isAllContain(fmt.Sprint(param), fmt.Sprint(utils.PathSearch("severity", v, nil))) {
				continue
			}
		}
		if param, ok := d.GetOk("category"); ok {
			if !isAllContain(fmt.Sprint(param), fmt.Sprint(utils.PathSearch("category", v, nil))) {
				continue
			}
		}
		rst = append(rst, v)
	}
	return rst
}

func isAllContain(filter string, target string) bool {
	filterList := strings.Split(filter, ",")
	for _, v := range filterList {
		if strings.Contains(target, v) {
			continue
		}
		return false
	}
	return true
}
