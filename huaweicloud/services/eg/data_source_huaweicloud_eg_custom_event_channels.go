package eg

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
				Description: `The region where the custom event channels are located.`,
			},
			"channel_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The channel ID used to query specified custom event channel.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The channel name used to query specified custom event channel.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the custom event channels belong.`,
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the custom event channel.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom event channel.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the custom event channel.`,
						},
						"provider_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the custom event channel.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the custom event channel belongs.`,
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
							Description: `The creation time of the custom event channel.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the custom event channel.`,
						},
					},
				},
			},
		},
	}
}

func buildEventChannelsQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, apiName)
	}
	if channelId, ok := d.GetOk("channel_id"); ok {
		res = fmt.Sprintf("%s&channel_id=%v", res, channelId)
	}
	return res
}

func queryEventChannels(client *golangsdk.ServiceClient, d *schema.ResourceData, providerType string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/channels?provider_type={provider_type}&limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{provider_type}", providerType)
	listPath += buildEventChannelsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		channels := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(channels) < 1 {
			break
		}
		result = append(result, channels...)
		offset += len(channels)
	}

	return result, nil
}

func flattenEventChannels(channels []interface{}) []interface{} {
	result := make([]interface{}, 0, len(channels))

	for _, channel := range channels {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", channel, nil),
			"name":                  utils.PathSearch("name", channel, nil),
			"provider_type":         utils.PathSearch("provider_type", channel, nil),
			"enterprise_project_id": utils.PathSearch("eps_id", channel, nil),
			"cross_account_ids":     utils.PathSearch("policy.Principal.IAM", channel, make([]interface{}, 0)),
			"created_at":            utils.PathSearch("created_time", channel, nil),
			"updated_at":            utils.PathSearch("updated_time", channel, nil),
		})
	}

	return result
}

func filterEventChannels(cfg *config.Config, d *schema.ResourceData, channels []interface{}) []interface{} {
	// Copy slice contents without having to worry about underlying reuse issues.
	result := channels
	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		result = utils.PathSearch(fmt.Sprintf("[?eps_id=='%s']", epsId), result, make([]interface{}, 0)).([]interface{})
	}
	return result
}

func dataSourceCustomEventChannelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	channaels, err := queryEventChannels(client, d, "CUSTOM")
	if err != nil {
		return diag.Errorf("error querying custom event channels: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("channels", flattenEventChannels(filterEventChannels(cfg, d, channaels))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG custom event channels: %s", err)
	}
	return nil
}
