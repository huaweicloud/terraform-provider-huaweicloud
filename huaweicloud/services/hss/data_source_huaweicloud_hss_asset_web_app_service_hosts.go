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

// @API HSS GET /v5/{project_id}/asset/web-app-and-services
func DataSourceAssetWebAppServiceHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetWebAppServiceHostsRead,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"install_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"part_match": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"catalogue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"install_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ctime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mtime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"atime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"proc_path": {
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
						"record_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAssetWebAppServiceHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?category=%v&catalogue=%v&name=%v&limit=200", d.Get("category"), d.Get("catalogue"), d.Get("name"))

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("version"); ok {
		queryParams = fmt.Sprintf("%s&version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("install_dir"); ok {
		queryParams = fmt.Sprintf("%s&install_dir=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if d.Get("part_match").(bool) {
		queryParams = fmt.Sprintf("%s&part_match=true", queryParams)
	}

	return queryParams
}

func dataSourceAssetWebAppServiceHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/asset/web-app-and-services"
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
	getPath += buildAssetWebAppServiceHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the hosts: %s", err)
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
		d.Set("data_list", flattenAssetWebAppServiceHosts(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetWebAppServiceHosts(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"catalogue":      utils.PathSearch("catalogue", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"version":        utils.PathSearch("version", v, nil),
			"agent_id":       utils.PathSearch("agent_id", v, nil),
			"install_path":   utils.PathSearch("install_path", v, nil),
			"config_path":    utils.PathSearch("config_path", v, nil),
			"uid":            utils.PathSearch("uid", v, nil),
			"gid":            utils.PathSearch("gid", v, nil),
			"mode":           utils.PathSearch("mode", v, nil),
			"ctime":          utils.PathSearch("ctime", v, nil),
			"mtime":          utils.PathSearch("mtime", v, nil),
			"atime":          utils.PathSearch("atime", v, nil),
			"pid":            utils.PathSearch("pid", v, nil),
			"proc_path":      utils.PathSearch("proc_path", v, nil),
			"container_id":   utils.PathSearch("container_id", v, nil),
			"container_name": utils.PathSearch("container_name", v, nil),
			"record_time":    utils.PathSearch("record_time", v, nil),
			"host_id":        utils.PathSearch("host_id", v, nil),
			"host_name":      utils.PathSearch("host_name", v, nil),
			"host_ip":        utils.PathSearch("host_ip", v, nil),
		})
	}

	return result
}
