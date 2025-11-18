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

// @API HSS GET /v5/{project_id}/image/baseline/risk-configs
func DataSourceImageBaselineRiskConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageBaselineRiskConfigsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
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
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_rule_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"failed_rule_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"check_type_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildImageBaselineRiskConfigsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	queryParams = fmt.Sprintf("%s&image_type=%v", queryParams, d.Get("image_type"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_version"); ok {
		queryParams = fmt.Sprintf("%s&image_version=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_id"); ok {
		queryParams = fmt.Sprintf("%s&image_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_name"); ok {
		queryParams = fmt.Sprintf("%s&check_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("standard"); ok {
		queryParams = fmt.Sprintf("%s&standard=%v", queryParams, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImageBaselineRiskConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/image/baseline/risk-configs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildImageBaselineRiskConfigsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS image baseline risk configs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenImageBaselineRiskConfigsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageBaselineRiskConfigsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"severity":        utils.PathSearch("severity", v, nil),
			"check_name":      utils.PathSearch("check_name", v, nil),
			"check_type":      utils.PathSearch("check_type", v, nil),
			"standard":        utils.PathSearch("standard", v, nil),
			"check_rule_num":  utils.PathSearch("check_rule_num", v, nil),
			"failed_rule_num": utils.PathSearch("failed_rule_num", v, nil),
			"check_type_desc": utils.PathSearch("check_type_desc", v, nil),
		})
	}

	return rst
}
