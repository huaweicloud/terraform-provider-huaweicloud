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

// @API HSS GET /v5/{project_id}/image/whitelists
func DataSourceImageWhitelists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageWhitelistsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"global_image_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vul_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vul_type": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vul_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vul_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vul_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cves": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cve_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvss": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildImageWhitelistsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?global_image_type=%v&type=%v&limit=200", d.Get("global_image_type"), d.Get("type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if v, ok := d.GetOk("vul_name"); ok {
		queryParams = fmt.Sprintf("%s&vul_name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("vul_type"); ok {
		queryParams = fmt.Sprintf("%s&vul_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImageWhitelistsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/image/whitelists"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildImageWhitelistsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving image whitelists: %s", err)
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
		d.Set("data_list", flattenImageWhitelists(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageWhitelists(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"vul_id":      utils.PathSearch("vul_id", v, nil),
			"vul_name":    utils.PathSearch("vul_name", v, nil),
			"vul_type":    utils.PathSearch("vul_type", v, nil),
			"cves":        flattenImageWhitelistsCves(utils.PathSearch("cves", v, make([]interface{}, 0)).([]interface{})),
			"rule_type":   utils.PathSearch("rule_type", v, nil),
			"query_info":  flattenImageWhitelistsQueryInfo(utils.PathSearch("query_info", v, nil)),
			"image_info":  flattenImageWhitelistsImageInfo(utils.PathSearch("image_info", v, make([]interface{}, 0)).([]interface{})),
			"description": utils.PathSearch("description", v, nil),
		})
	}

	return result
}

func flattenImageWhitelistsCves(rawCveList []interface{}) []interface{} {
	if len(rawCveList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawCveList))
	for _, v := range rawCveList {
		result = append(result, map[string]interface{}{
			"cve_id": utils.PathSearch("cve_id", v, nil),
			"cvss":   utils.PathSearch("cvss", v, nil),
		})
	}

	return result
}

func flattenImageWhitelistsQueryInfo(rawInfo interface{}) []map[string]interface{} {
	if rawInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"image_type": utils.PathSearch("image_type", rawInfo, nil),
	}

	return []map[string]interface{}{result}
}

func flattenImageWhitelistsImageInfo(imageInfo []interface{}) []interface{} {
	if len(imageInfo) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(imageInfo))
	for _, v := range imageInfo {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"image_id":   utils.PathSearch("image_id", v, nil),
			"image_name": utils.PathSearch("image_name", v, nil),
		})
	}

	return result
}
