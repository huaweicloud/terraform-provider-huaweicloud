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

// @API HSS GET /v5/{project_id}/asset/processes/detail
func DataSourceAssetProcessDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetProcessDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
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
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_params": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"process_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"run_permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAssetProcessDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("path"); ok {
		queryParams = fmt.Sprintf("%s&path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceAssetProcessDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/asset/processes/detail"
		epsId    = cfg.GetEnterpriseProjectID(d)
		offset   = 0
		result   = make([]interface{}, 0)
		totalNum float64
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAssetProcessDetailQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving servers of a specified process: %s", err)
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

		totalNum = utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		if int(totalNum) == len(result) {
			break
		}

		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenAssetProcessDetail(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetProcessDetail(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"hash":           utils.PathSearch("hash", v, nil),
			"host_ip":        utils.PathSearch("host_ip", v, nil),
			"host_name":      utils.PathSearch("host_name", v, nil),
			"launch_params":  utils.PathSearch("launch_params", v, nil),
			"launch_time":    utils.PathSearch("launch_time", v, nil),
			"process_path":   utils.PathSearch("process_path", v, nil),
			"process_pid":    utils.PathSearch("process_pid", v, nil),
			"run_permission": utils.PathSearch("run_permission", v, nil),
			"container_id":   utils.PathSearch("container_id", v, nil),
			"container_name": utils.PathSearch("container_name", v, nil),
		})
	}

	return result
}
