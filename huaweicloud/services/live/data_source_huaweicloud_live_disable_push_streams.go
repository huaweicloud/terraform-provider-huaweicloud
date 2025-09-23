package live

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

// @API LIVE GET /v1/{project_id}/stream/blocks
func DataSourceDisablePushStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDisablePushStreamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ingest domain name of the disabling push stream.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name of the disabling push stream.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the stream name of the disabling push stream.`,
			},
			"blocks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the disabled push streams.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name of the disabling push stream.`,
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The stream name of the disabling push stream.`,
						},
						"resume_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time of the resuming push stream.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDisablePushStreamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	recordings, err := queryDisablePushStreams(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("blocks", flattenDisablePushStreams(recordings)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryDisablePushStreams(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/stream/blocks?size=100"
		page    = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s&domain=%v", listPath, d.Get("domain_name"))
	listPath += buildDisablePushStreamsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		// The page indicates the page number.
		// The default value is 0, which represents the first page.
		listPathWithPage := fmt.Sprintf("%s&page=%d", listPath, page)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving disabled push streams information: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		block := utils.PathSearch("blocks", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, block...)
		if len(block) == 0 {
			break
		}
		page++
	}
	return result, nil
}

func buildDisablePushStreamsQueryParams(d *schema.ResourceData) string {
	res := ""
	if appName, ok := d.GetOk("app_name"); ok {
		res = fmt.Sprintf("%s&app_name=%v", res, appName)
	}
	if streamName, ok := d.GetOk("stream_name"); ok {
		res = fmt.Sprintf("%s&stream_name=%v", res, streamName)
	}
	return res
}

func flattenDisablePushStreams(blocks []interface{}) []map[string]interface{} {
	if len(blocks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(blocks))
	for i, v := range blocks {
		result[i] = map[string]interface{}{
			"app_name":    utils.PathSearch("app_name", v, nil),
			"stream_name": utils.PathSearch("stream_name", v, nil),
			"resume_time": utils.PathSearch("resume_time", v, nil),
		}
	}
	return result
}
