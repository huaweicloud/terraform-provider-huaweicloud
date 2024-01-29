package eg

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eg/v1/channel/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG GET /v1/{project_id}/channels
func DataSourceCustomEventChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomEventChannelsRead,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the custom event channels are located.",
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The channel ID used to query specified custom event channel.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The channel name used to query specified custom event channel.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the custom event channels belong.",
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the custom event channel.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the custom event channel.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the custom event channel.",
						},
						"provider_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the custom event channel.",
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the enterprise project to which the custom event channel belongs.",
						},
						"cross_account_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The list of domain IDs (other tenants) for the cross-account policy.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the custom event channel.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the custom event channel.",
						},
					},
				},
			},
		},
	}
}

func filterCustomEventChannels(channels []custom.Channel, epsId string) ([]interface{}, error) {
	filter := map[string]interface{}{
		"EnterpriseProjectId": epsId,
	}

	filterResult, err := utils.FilterSliceWithField(channels, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting list of custom event channels: %s", err)
	}
	return filterResult, nil
}

func flattenCustomEventChannels(channels []interface{}) []map[string]interface{} {
	if len(channels) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(channels))
	for i, val := range channels {
		channel := val.(custom.Channel)
		result[i] = map[string]interface{}{
			"id":                    channel.ID,
			"name":                  channel.Name,
			"provider_type":         channel.ProviderType,
			"enterprise_project_id": channel.EnterpriseProjectId,
			"cross_account_ids":     channel.Policy.Principal.IAM,
			"created_at":            channel.CreatedTime,
			"updated_at":            channel.UpdatedTime,
		}
	}
	return result
}

func dataSourceCustomEventChannelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error querying custom event channels: %s", err)
	}
	filterResult, err := filterCustomEventChannels(resp, cfg.GetEnterpriseProjectID(d))
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
		d.Set("channels", flattenCustomEventChannels(filterResult)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG custom event channels: %s", err)
	}
	return nil
}
