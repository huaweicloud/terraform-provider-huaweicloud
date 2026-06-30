package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/asset/ai-component/detail
func DataSourceAiComponentDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAiComponentDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"catalogue": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ai_application": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ai_tool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"installation_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The parameters `first_scan_time` and `latest_scan_time` are of type int in the API documentation,
			// but here they are changed to type string.
			// Because it needs to meet the scenario that can be set to `0`.
			"first_scan_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latest_scan_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ai_application": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ai_tool": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"startup_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"startup_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"install_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cmdline": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ppid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listen_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listen_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listen_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"listen_status": {
										Type:     schema.TypeString,
										Computed: true,
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

func buildAiComponentDetailQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?limit=200&category=%v&catalogue=%v", d.Get("category").(string),
		d.Get("catalogue").(string))

	if v, ok := d.GetOk("server_name"); ok {
		queryParams = fmt.Sprintf("%s&server_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("server_ip"); ok {
		queryParams = fmt.Sprintf("%s&server_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("ai_application"); ok {
		queryParams = fmt.Sprintf("%s&ai_application=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("ai_tool"); ok {
		queryParams = fmt.Sprintf("%s&ai_tool=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("installation_path"); ok {
		queryParams = fmt.Sprintf("%s&installation_path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("first_scan_time"); ok {
		queryParams = fmt.Sprintf("%s&first_scan_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("latest_scan_time"); ok {
		queryParams = fmt.Sprintf("%s&latest_scan_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_name"); ok {
		queryParams = fmt.Sprintf("%s&container_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_id"); ok {
		queryParams = fmt.Sprintf("%s&container_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAiComponentDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/asset/ai-component/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAiComponentDetailQueryParams(d)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS AI component detail: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenAiComponentDetailDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAiComponentDetailDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"server_name":      utils.PathSearch("server_name", v, nil),
			"server_ip":        utils.PathSearch("server_ip", v, nil),
			"ai_application":   utils.PathSearch("ai_application", v, nil),
			"ai_tool":          utils.PathSearch("ai_tool", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"version":          utils.PathSearch("version", v, nil),
			"startup_path":     utils.PathSearch("startup_path", v, nil),
			"startup_time":     utils.PathSearch("startup_time", v, nil),
			"install_path":     utils.PathSearch("install_path", v, nil),
			"cmdline":          utils.PathSearch("cmdline", v, nil),
			"first_scan_time":  utils.PathSearch("first_scan_time", v, nil),
			"latest_scan_time": utils.PathSearch("latest_scan_time", v, nil),
			"container_name":   utils.PathSearch("container_name", v, nil),
			"container_id":     utils.PathSearch("container_id", v, nil),
			"host_id":          utils.PathSearch("host_id", v, nil),
			"pid":              utils.PathSearch("pid", v, nil),
			"ppid":             utils.PathSearch("ppid", v, nil),
			"user":             utils.PathSearch("user", v, nil),
			"net_info": flattenAiComponentDetailNetInfo(
				utils.PathSearch("net_info", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenAiComponentDetailNetInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"listen_ip":       utils.PathSearch("listen_ip", v, nil),
			"listen_protocol": utils.PathSearch("listen_protocol", v, nil),
			"listen_port":     utils.PathSearch("listen_port", v, nil),
			"listen_status":   utils.PathSearch("listen_status", v, nil),
		})
	}

	return rst
}
