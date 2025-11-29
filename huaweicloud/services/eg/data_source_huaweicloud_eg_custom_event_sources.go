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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
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
			"fuzzy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the channels to be queried for fuzzy matching.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting method for query results.`,
			},

			// Attributes
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
						"error_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The error information of the custom event source.`,
							Elem:        errorInfoSchema(),
						},
					},
				},
			},
		},
	}
}

func errorInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error code of current source.`,
			},
			"error_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error detail of current source.`,
			},
			"error_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error message of current source.`,
			},
		},
	}
}

func buildEventSourcesQueryParams(d *schema.ResourceData, providerTypeInput ...string) string {
	res := ""
	if sourceName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, sourceName)
	}
	if channelId, ok := d.GetOk("channel_id"); ok {
		res = fmt.Sprintf("%s&channel_id=%v", res, channelId)
	}
	if fuzzyName, ok := d.GetOk("fuzzy_name"); ok {
		res = fmt.Sprintf("%s&fuzzy_name=%v", res, fuzzyName)
	}
	if sort, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, sort)
	}

	if len(providerTypeInput) > 0 {
		res = fmt.Sprintf("%s&provider_type=%v", res, providerTypeInput[0])
	} else if typeVal, ok := d.GetOk("provider_type"); ok {
		res = fmt.Sprintf("%s&provider_type=%v", res, typeVal)
	}
	return res
}

func queryEventSources(client *golangsdk.ServiceClient, d *schema.ResourceData, providerType ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/sources?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildEventSourcesQueryParams(d, providerType...)

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
			// If the offset is greater than or equal to the total number, an error will be returned, not empey page (record list).
			parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", "APIGW.0106")
			if _, ok := parsedErr.(golangsdk.ErrDefault404); ok {
				break
			}
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		eventSources := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(eventSources) < 1 {
			break
		}
		result = append(result, eventSources...)
		offset += len(eventSources)
	}

	return result, nil
}

func flattenDataCustomEventSources(eventSources []interface{}) []interface{} {
	result := make([]interface{}, 0, len(eventSources))

	for _, eventSource := range eventSources {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", eventSource, nil),
			"channel_id":   utils.PathSearch("channel_id", eventSource, nil),
			"channel_name": utils.PathSearch("channel_name", eventSource, nil),
			"name":         utils.PathSearch("name", eventSource, nil),
			"type":         utils.PathSearch("type", eventSource, nil),
			"description":  utils.PathSearch("description", eventSource, nil),
			"status":       utils.PathSearch("status", eventSource, nil),
			"created_at":   utils.PathSearch("created_at", eventSource, nil),
			"updated_at":   utils.PathSearch("updated_at", eventSource, nil),
			"error_info": flattenConnectionErrorInfo(utils.PathSearch(
				"error_info", eventSource, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func filterDataCustomEventSources(d *schema.ResourceData, eventSources []interface{}) []interface{} {
	// Copy slice contents without having to worry about underlying reuse issues.
	result := eventSources
	if sourceId, ok := d.GetOk("source_id"); ok {
		result = utils.PathSearch(fmt.Sprintf("[?id=='%s']", sourceId), result, make([]interface{}, 0)).([]interface{})
	}
	return result
}

func dataSourceCustomEventSourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	eventSources, err := queryEventSources(client, d, "CUSTOM")
	if err != nil {
		return diag.Errorf("error querying custom event sources: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sources", flattenDataCustomEventSources(filterDataCustomEventSources(d, eventSources))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of EG custom event sources: %s", err)
	}
	return nil
}
