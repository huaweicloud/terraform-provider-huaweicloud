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

// @API HSS GET /v5/{project_id}/asset/host/web-site
func DataSourceAssetWebsiteHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetWebsiteHostsRead,

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
			"domain": {
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
			"part_match": {
				Type:     schema.TypeBool,
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
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bind_addr": {
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
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proc_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cert_issuer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_issue_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_time": {
							Type:     schema.TypeInt,
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

func buildAssetWebsiteHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?category=%v&domain=%v&limit=200", d.Get("category"), d.Get("domain"))

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

func dataSourceAssetWebsiteHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v5/{project_id}/asset/host/web-site"
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
	getPath += buildAssetWebsiteHostsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving a specified website servers: %s", err)
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
		d.Set("data_list", flattenAssetWebsiteHostsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetWebsiteHostsDataList(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"agent_id":          utils.PathSearch("agent_id", v, nil),
			"host_id":           utils.PathSearch("host_id", v, nil),
			"host_name":         utils.PathSearch("host_name", v, nil),
			"host_ip":           utils.PathSearch("host_ip", v, nil),
			"domain":            utils.PathSearch("domain", v, nil),
			"app_name":          utils.PathSearch("app_name", v, nil),
			"path":              utils.PathSearch("path", v, nil),
			"port":              utils.PathSearch("port", v, nil),
			"bind_addr":         utils.PathSearch("bind_addr", v, nil),
			"url_path":          utils.PathSearch("url_path", v, nil),
			"uid":               utils.PathSearch("uid", v, nil),
			"gid":               utils.PathSearch("gid", v, nil),
			"mode":              utils.PathSearch("mode", v, nil),
			"pid":               utils.PathSearch("pid", v, nil),
			"proc_path":         utils.PathSearch("proc_path", v, nil),
			"is_https":          utils.PathSearch("is_https", v, nil),
			"cert_issuer":       utils.PathSearch("cert_issuer", v, nil),
			"cert_user":         utils.PathSearch("cert_user", v, nil),
			"cert_issue_time":   utils.PathSearch("cert_issue_time", v, nil),
			"cert_expired_time": utils.PathSearch("cert_expired_time", v, nil),
			"record_time":       utils.PathSearch("record_time", v, nil),
			"container_id":      utils.PathSearch("container_id", v, nil),
			"container_name":    utils.PathSearch("container_name", v, nil),
		})
	}

	return result
}
