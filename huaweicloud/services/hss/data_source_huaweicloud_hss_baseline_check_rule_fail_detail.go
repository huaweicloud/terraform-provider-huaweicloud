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

// @API HSS GET /v5/{project_id}/baseline/check-rule/fail-detail
func DataSourceBaselineCheckRuleFailDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineCheckRuleFailDetailRead,

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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fail_detail_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceBaselineCheckRuleFailDetailListSchema(),
			},
		},
	}
}

func dataSourceBaselineCheckRuleFailDetailListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fix_fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBaselineCheckRuleFailDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?check_rule_id=%v", d.Get("check_rule_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_name"); ok {
		queryParams = fmt.Sprintf("%s&check_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("standard"); ok {
		queryParams = fmt.Sprintf("%s&standard=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBaselineCheckRuleFailDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/baseline/check-rule/fail-detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineCheckRuleFailDetailQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS baseline check rule fail detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("fail_detail_list", flattenBaselineCheckRuleFailDetailList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineCheckRuleFailDetailList(respBody interface{}) []interface{} {
	failDetailList := utils.PathSearch("fail_detail_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(failDetailList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(failDetailList))
	for _, v := range failDetailList {
		rst = append(rst, map[string]interface{}{
			"fix_fail_reason": utils.PathSearch("fix_fail_reason", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
		})
	}

	return rst
}
