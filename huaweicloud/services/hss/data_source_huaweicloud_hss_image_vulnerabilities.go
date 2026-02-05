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

// @API HSS GET /v5/{project_id}/image/vulnerabilities
func DataSourceImageVulnerabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageVulnerabilitiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repair_necessity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vul_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vul_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vul_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repair_necessity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// The API documentation is `decription`, but the actual return is `description`.
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"solution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"history_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"undeal_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cve_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvss_score": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"publish_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
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

func buildImageVulnerabilitiesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("repair_necessity"); ok {
		queryParams = fmt.Sprintf("%s&repair_necessity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("vul_id"); ok {
		queryParams = fmt.Sprintf("%s&vul_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImageVulnerabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v5/{project_id}/image/vulnerabilities"
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildImageVulnerabilitiesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS image vulnerabilities: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenImageVulnerabilitiesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageVulnerabilitiesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"vul_name":         utils.PathSearch("vul_name", v, nil),
			"vul_id":           utils.PathSearch("vul_id", v, nil),
			"repair_necessity": utils.PathSearch("repair_necessity", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"solution":         utils.PathSearch("solution", v, nil),
			"url":              utils.PathSearch("url", v, nil),
			"history_number":   utils.PathSearch("history_number", v, nil),
			"undeal_number":    utils.PathSearch("undeal_number", v, nil),
			"data_list": flattenImageVulnerabilitiesSubDataList(
				utils.PathSearch("data_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenImageVulnerabilitiesSubDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"cve_id":       utils.PathSearch("cve_id", v, nil),
			"cvss_score":   utils.PathSearch("cvss_score", v, nil),
			"publish_time": utils.PathSearch("publish_time", v, nil),
			"description":  utils.PathSearch("description", v, nil),
		})
	}

	return result
}
