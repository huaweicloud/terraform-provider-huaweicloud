package hss

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

// @API HSS GET /v5/{project_id}/webtamper/static/protect-history
func DataSourceWebtamperStaticProtectHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebtamperStaticProtectHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_operation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"occur_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_operation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"process_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_cmd": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWebtamperStaticProtectHistoryQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?start_time=%v&end_time=%v&limit=200", d.Get("start_time"), d.Get("start_time"))

	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_path"); ok {
		queryParams = fmt.Sprintf("%s&file_path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_operation"); ok {
		queryParams = fmt.Sprintf("%s&file_operation=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceWebtamperStaticProtectHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/webtamper/static/protect-history"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWebtamperStaticProtectHistoryQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving static WTP events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenWebtamperStaticProtectHistory(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWebtamperStaticProtectHistory(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_name":      utils.PathSearch("host_name", v, nil),
			"occur_time":     utils.PathSearch("occur_time", v, nil),
			"file_path":      utils.PathSearch("file_path", v, nil),
			"file_operation": utils.PathSearch("file_operation", v, nil),
			"host_ip":        utils.PathSearch("host_ip", v, nil),
			"process_id":     utils.PathSearch("process_id", v, nil),
			"process_name":   utils.PathSearch("process_name", v, nil),
			"process_cmd":    utils.PathSearch("process_cmd", v, nil),
		})
	}

	return rst
}
