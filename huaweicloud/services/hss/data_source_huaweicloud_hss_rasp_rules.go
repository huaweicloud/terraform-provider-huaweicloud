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

// @API HSS GET /v5/{project_id}/rasp/rule
func DataSourceRaspRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRaspRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
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
						"chk_feature_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"chk_feature_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chk_feature_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"feature_configure": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"optional_protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"editable": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRaspRulesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceRaspRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/rasp/rule"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRaspRulesQueryParams(d, epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving detection rules: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	hosts := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenRaspRules(hosts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRaspRules(hostsResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		rst = append(rst, map[string]interface{}{
			"chk_feature_id":             utils.PathSearch("chk_feature_id", v, nil),
			"chk_feature_name":           utils.PathSearch("chk_feature_name", v, nil),
			"chk_feature_desc":           utils.PathSearch("chk_feature_desc", v, nil),
			"os_type":                    utils.PathSearch("os_type", v, nil),
			"feature_configure":          utils.PathSearch("feature_configure", v, nil),
			"optional_protective_action": utils.PathSearch("optional_protective_action", v, nil),
			"protective_action":          utils.PathSearch("protective_action", v, nil),
			"editable":                   utils.PathSearch("editable", v, nil),
		})
	}

	return rst
}
