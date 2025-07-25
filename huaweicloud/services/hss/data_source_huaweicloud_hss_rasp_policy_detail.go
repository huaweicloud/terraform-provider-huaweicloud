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

// @API HSS GET /v5/{project_id}/rasp/policy/detail
func DataSourceRaspPolicyDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRaspPolicyDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_list": {
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
						"feature_configure": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"optional_protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
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

func buildRaspPolicyDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?policy_id=%v", d.Get("policy_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceRaspPolicyDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/rasp/policy/detail"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRaspPolicyDetailQueryParams(d, epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the protection policy details: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyDetails := utils.PathSearch("rule_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_name", utils.PathSearch("policy_name", getRespBody, nil)),
		d.Set("os_type", utils.PathSearch("os_type", getRespBody, nil)),
		d.Set("rule_list", flattenRaspPolicyDetail(policyDetails)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRaspPolicyDetail(hostsResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(hostsResp))
	for _, v := range hostsResp {
		rst = append(rst, map[string]interface{}{
			"chk_feature_id":             utils.PathSearch("chk_feature_id", v, nil),
			"chk_feature_name":           utils.PathSearch("chk_feature_name", v, nil),
			"chk_feature_desc":           utils.PathSearch("chk_feature_desc", v, nil),
			"feature_configure":          utils.PathSearch("feature_configure", v, nil),
			"protective_action":          utils.PathSearch("protective_action", v, nil),
			"optional_protective_action": utils.PathSearch("optional_protective_action", v, nil),
			"enabled":                    utils.PathSearch("enabled", v, nil),
			"editable":                   utils.PathSearch("editable", v, nil),
		})
	}

	return rst
}
