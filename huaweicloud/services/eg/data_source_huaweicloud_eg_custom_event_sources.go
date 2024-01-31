package eg

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eg/v1/source/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG GET /v1/{project_id}/sources
func DataSourceCustomEventSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomEventSourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the custom event sources are located.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the custom event channel to which the custom event sources belong.`,
			},
			"source_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The event source ID used to query specified custom event source.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The event source name used to query specified custom event source.`,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the custom event source.`,
						},
						"channel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the custom event channel to which the custom event source belong.`,
						},
						"channel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom event channel to which the custom event source belong.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom event source.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the custom event source.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the custom event source.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the custom event source.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the custom event source.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the custom event source.`,
						},
					},
				},
			},
		},
	}
}

func filterEventSources(d *schema.ResourceData, eventSources []custom.Source) ([]interface{}, error) {
	filter := map[string]interface{}{
		"ID": d.Get("source_id"),
	}

	filterResult, err := utils.FilterSliceWithField(eventSources, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting list of custom event sources: %s", err)
	}
	return filterResult, nil
}

func flattenCustomEventSources(eventSources []interface{}) []map[string]interface{} {
	if len(eventSources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(eventSources))
	for i, val := range eventSources {
		eventSource := val.(custom.Source)
		result[i] = map[string]interface{}{
			"id":           eventSource.ID,
			"channel_id":   eventSource.ChannelId,
			"channel_name": eventSource.ChannelName,
			"name":         eventSource.Name,
			"type":         eventSource.Type,
			"description":  eventSource.Description,
			"status":       eventSource.Status,
			"created_at":   eventSource.CreatedTime,
			"updated_at":   eventSource.UpdatedTime,
		}
	}
	return result
}

func dataSourceCustomEventSourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		opts   = custom.ListOpts{
			ChannelId:    d.Get("channel_id").(string),
			ProviderType: "CUSTOM",
			Name:         d.Get("name").(string),
		}
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := custom.List(client, opts)
	if err != nil {
		return diag.Errorf("error querying custom event sources: %s", err)
	}
	filterResult, err := filterEventSources(d, resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sources", flattenCustomEventSources(filterResult)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG custom event sources: %s", err)
	}
	return nil
}
