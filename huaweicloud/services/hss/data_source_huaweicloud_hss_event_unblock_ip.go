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

// @API HSS GET /v5/{project_id}/event/blocked-ip
func DataSourceEventUnblockIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventUnblockIpRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"intercept_status": {
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
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intercept_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"intercept_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"latest_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventUnblockIpQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=100"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("last_days"); ok {
		queryParams = fmt.Sprintf("%s&last_days=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("src_ip"); ok {
		queryParams = fmt.Sprintf("%s&src_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("intercept_status"); ok {
		queryParams = fmt.Sprintf("%s&intercept_status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventUnblockIpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		httpUrl = "v5/{project_id}/event/blocked-ip"
		epsId   = cfg.GetEnterpriseProjectID(d)
		offset  = 0
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildEventUnblockIpQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS unblock IP: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
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

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataList(dataResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"host_id":          utils.PathSearch("host_id", v, nil),
			"host_name":        utils.PathSearch("host_name", v, nil),
			"src_ip":           utils.PathSearch("src_ip", v, nil),
			"login_type":       utils.PathSearch("login_type", v, nil),
			"intercept_num":    utils.PathSearch("intercept_num", v, nil),
			"intercept_status": utils.PathSearch("intercept_status", v, nil),
			"block_time":       utils.PathSearch("block_time", v, nil),
			"latest_time":      utils.PathSearch("latest_time", v, nil),
		})
	}

	return result
}
