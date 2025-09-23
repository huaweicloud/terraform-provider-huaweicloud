package eg

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EG GET /v1/{project_id}/channels
func DataSourceEventChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventChannelsRead,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the event channels are located.`,
			},
			"provider_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the event channels to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The channel name used to query specified event channel.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the event channels belong.`,
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event channel.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event channel.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the event channel.`,
						},
						"provider_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the event channel.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the event channel belongs.`,
						},
						"cross_account_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of domain IDs (other tenants) for the cross-account policy.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the event channel.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the event channel.`,
						},
					},
				},
				Description: `All event channels that match the filter parameters.`,
			},
		},
	}
}

func dataSourceEventChannelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	channels, err := queryEventChannels(client, d)
	if err != nil {
		return diag.Errorf("error querying event channels: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("channels", flattenDataEventChannels(filterDataEventChannels(cfg, d, channels))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG event channels: %s", err)
	}
	return nil
}
