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

// @API HSS GET /v5/{project_id}/baseline/check-rule/hosts
func DataSourceBaselineCheckRuleHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineCheckRuleHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"check_rule_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_name": {
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
				Elem:     dataSourceBaselineCheckRuleHostsSchema(),
			},
		},
	}
}

func dataSourceBaselineCheckRuleHostsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"baseline_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"passed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"diff_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_fix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_verify": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_click": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cancel_ignore_enable_click": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"result_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fix_failed_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBaselineCheckRuleHostsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?check_rule_id=%s&standard=%s&limit=200", d.Get("check_rule_id"),
		d.Get("standard"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("check_name"); ok {
		queryParams = fmt.Sprintf("%s&check_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_type"); ok {
		queryParams = fmt.Sprintf("%s&check_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("result_type"); ok {
		queryParams = fmt.Sprintf("%s&result_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBaselineCheckRuleHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum float64
		httpUrl  = "v5/{project_id}/baseline/check-rule/hosts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineCheckRuleHostsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS baseline check rule hosts: %s", err)
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
		d.Set("data_list", flattenBaselineCheckRuleHostsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineCheckRuleHostsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"host_id":                    utils.PathSearch("host_id", v, nil),
			"host_name":                  utils.PathSearch("host_name", v, nil),
			"check_name":                 utils.PathSearch("check_name", v, nil),
			"baseline_name":              utils.PathSearch("baseline_name", v, nil),
			"host_public_ip":             utils.PathSearch("host_public_ip", v, nil),
			"host_private_ip":            utils.PathSearch("host_private_ip", v, nil),
			"scan_time":                  utils.PathSearch("scan_time", v, nil),
			"failed_num":                 utils.PathSearch("failed_num", v, nil),
			"passed_num":                 utils.PathSearch("passed_num", v, nil),
			"diff_description":           utils.PathSearch("diff_description", v, nil),
			"description":                utils.PathSearch("description", v, nil),
			"host_type":                  utils.PathSearch("host_type", v, nil),
			"enable_fix":                 utils.PathSearch("enable_fix", v, nil),
			"enable_verify":              utils.PathSearch("enable_verify", v, nil),
			"enable_click":               utils.PathSearch("enable_click", v, nil),
			"cancel_ignore_enable_click": utils.PathSearch("cancel_ignore_enable_click", v, nil),
			"result_type":                utils.PathSearch("result_type", v, nil),
			"fix_failed_reason":          utils.PathSearch("fix_failed_reason", v, nil),
			"cluster_id":                 utils.PathSearch("cluster_id", v, nil),
		})
	}

	return rst
}
