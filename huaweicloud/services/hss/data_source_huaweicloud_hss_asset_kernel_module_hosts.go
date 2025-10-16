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

// @API HSS GET /v5/{project_id}/asset/host/kernel-module
func DataSourceAssetKernelModuleHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetKernelModuleHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect
			"part_match": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": {
							Type:     schema.TypeString,
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
						"kernel_module_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"srcversion": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"uid": {
										Type:     schema.TypeInt,
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
									"hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desc": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"record_time": {
										Type:     schema.TypeInt,
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

func buildAssetKernelModuleHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?name=%v&limit=200", d.Get("name"))

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_ip"); ok {
		queryParams = fmt.Sprintf("%s&host_ip=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if d.Get("part_match").(bool) {
		queryParams = fmt.Sprintf("%s&part_match=true", queryParams)
	}

	return queryParams
}

func dataSourceAssetKernelModuleHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/asset/host/kernel-module"
		epsId   = cfg.GetEnterpriseProjectID(d)
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAssetKernelModuleHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving servers for a specified kernel module: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenAssetKernelModuleHostsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetKernelModuleHostsDataList(dataResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"agent_id":           utils.PathSearch("agent_id", v, nil),
			"host_id":            utils.PathSearch("host_id", v, nil),
			"host_name":          utils.PathSearch("host_name", v, nil),
			"host_ip":            utils.PathSearch("host_ip", v, nil),
			"kernel_module_info": flattenHostKernelModuleInfo(utils.PathSearch("kernel_module_info", v, nil)),
		})
	}

	return result
}

func flattenHostKernelModuleInfo(kernelModuleInfo interface{}) []map[string]interface{} {
	if kernelModuleInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"name":        utils.PathSearch("name", kernelModuleInfo, nil),
		"file_name":   utils.PathSearch("file_name", kernelModuleInfo, nil),
		"version":     utils.PathSearch("version", kernelModuleInfo, nil),
		"srcversion":  utils.PathSearch("srcversion", kernelModuleInfo, nil),
		"path":        utils.PathSearch("path", kernelModuleInfo, nil),
		"size":        utils.PathSearch("size", kernelModuleInfo, nil),
		"mode":        utils.PathSearch("mode", kernelModuleInfo, nil),
		"uid":         utils.PathSearch("uid", kernelModuleInfo, nil),
		"ctime":       utils.PathSearch("ctime", kernelModuleInfo, nil),
		"mtime":       utils.PathSearch("mtime", kernelModuleInfo, nil),
		"hash":        utils.PathSearch("hash", kernelModuleInfo, nil),
		"desc":        utils.PathSearch("desc", kernelModuleInfo, nil),
		"record_time": utils.PathSearch("record_time", kernelModuleInfo, nil),
	}

	return []map[string]interface{}{result}
}
