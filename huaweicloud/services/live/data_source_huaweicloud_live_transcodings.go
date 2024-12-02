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

// @API LIVE GET /v1/{project_id}/template/transcodings
func DataSourceLiveTranscodings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLiveTranscodingsRead,

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
				Description: `Specifies the ingest domain name to which the transcoding templates blong.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name of the transcoding template.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the transcoding templates.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name of the transcoding template.`,
						},
						"quality_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The video quality information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The transcoding template name.`,
									},
									"quality": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The video quality.`,
									},
									"video_encoding": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The video encoding format.`,
									},
									"width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The video long edge (width of horizontal screen, height of vertical screen), in pixels.`,
									},
									"height": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The video short edge (horizontal screen height, vertical screen width), in pixels.`,
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The bitrate of the transcoding video, in Kbps.`,
									},
									"frame_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The frame rate of transcoding video, in fps.`,
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The protocol type of transcoding output.`,
									},
									"low_bitrate_hd": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Whether to enable high-definition and low bitrate.`,
									},
									"gop": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The I frame interval, in seconds.`,
									},
									"bitrate_adaptive": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The adaptive bitrate.`,
									},
									"i_frame_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The maximum I frame interval, in frame.`,
									},
									"i_frame_policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The encoding output I frame policy.`,
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

func dataSourceLiveTranscodingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	transcodings, err := queryLiveTranscodings(d, client)
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
		d.Set("templates", flattenLiveTranscodings(transcodings)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryLiveTranscodings(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/template/transcodings?size=100"
		page    = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s&domain=%v", listPath, d.Get("domain_name"))
	if v, ok := d.GetOk("app_name"); ok {
		listPath = fmt.Sprintf("%s&app_name=%v", listPath, v)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		// The page indicates the page number.
		// The default value is 0, which represents the first page.
		listPathWithPage := fmt.Sprintf("%s&page=%d", listPath, page)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving transcodings: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		transcoding := utils.PathSearch("templates", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, transcoding...)
		if len(transcoding) == 0 {
			break
		}
		page++
	}
	return result, nil
}

func flattenLiveTranscodings(transcodings []interface{}) []map[string]interface{} {
	if len(transcodings) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(transcodings))
	for i, v := range transcodings {
		result[i] = map[string]interface{}{
			"app_name":     utils.PathSearch("app_name", v, nil),
			"quality_info": flattenQualityInfos(utils.PathSearch("quality_info", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return result
}

func flattenQualityInfos(qualityInfos []interface{}) []map[string]interface{} {
	if len(qualityInfos) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(qualityInfos))
	for i, v := range qualityInfos {
		result[i] = map[string]interface{}{
			"name":             utils.PathSearch("templateName", v, nil),
			"quality":          utils.PathSearch("quality", v, nil),
			"video_encoding":   utils.PathSearch("codec", v, nil),
			"width":            utils.PathSearch("width", v, nil),
			"height":           utils.PathSearch("height", v, nil),
			"bitrate":          utils.PathSearch("bitrate", v, nil),
			"frame_rate":       utils.PathSearch("video_frame_rate", v, nil),
			"protocol":         utils.PathSearch("protocol", v, nil),
			"low_bitrate_hd":   utils.PathSearch("hdlb", v, nil),
			"gop":              utils.PathSearch("gop", v, nil),
			"bitrate_adaptive": utils.PathSearch("bitrate_adaptive", v, nil),
			"i_frame_interval": utils.PathSearch("iFrameInterval", v, nil),
			"i_frame_policy":   utils.PathSearch("i_frame_policy", v, nil),
		}
	}
	return result
}
