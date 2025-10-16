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

// @API HSS GET /v5/{project_id}/app/policy
func DataSourceAppWhitelistPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppWhitelistPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"learning_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"intercept": {
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
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"learning_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"learning_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"specified_dir": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dir_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"file_extension_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"intercept": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_detect": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"not_effect_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"effect_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trust_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"suspicious_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"malicious_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unknown_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"abnormal_info_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"abnormal_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"abnormal_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"auto_confirm": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_policy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"host_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildAppWhitelistPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("policy_name"); ok {
		queryParams = fmt.Sprintf("%s&policy_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_type"); ok {
		queryParams = fmt.Sprintf("%s&policy_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("learning_status"); ok {
		queryParams = fmt.Sprintf("%s&learning_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("intercept"); ok {
		queryParams = fmt.Sprintf("%s&intercept=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceAppWhitelistPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v5/{project_id}/app/policy"
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
	getPath += buildAppWhitelistPoliciesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving process whitelist policies: %s", err)
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
		d.Set("data_list", flattenAppWhitelistPolicies(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppWhitelistPolicies(appPolicies []interface{}) []interface{} {
	if len(appPolicies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(appPolicies))
	for _, v := range appPolicies {
		rst = append(rst, map[string]interface{}{
			"policy_id":           utils.PathSearch("policy_id", v, nil),
			"policy_name":         utils.PathSearch("policy_name", v, nil),
			"policy_type":         utils.PathSearch("policy_type", v, nil),
			"learning_status":     utils.PathSearch("learning_status", v, nil),
			"learning_days":       utils.PathSearch("learning_days", v, nil),
			"specified_dir":       utils.PathSearch("specified_dir", v, nil),
			"dir_list":            utils.PathSearch("dir_list", v, nil),
			"file_extension_list": utils.PathSearch("file_extension_list", v, nil),
			"intercept":           utils.PathSearch("intercept", v, nil),
			"auto_detect":         utils.PathSearch("auto_detect", v, nil),
			"not_effect_host_num": utils.PathSearch("not_effect_host_num", v, nil),
			"effect_host_num":     utils.PathSearch("effect_host_num", v, nil),
			"trust_num":           utils.PathSearch("trust_num", v, nil),
			"suspicious_num":      utils.PathSearch("suspicious_num", v, nil),
			"malicious_num":       utils.PathSearch("malicious_num", v, nil),
			"unknown_num":         utils.PathSearch("unknown_num", v, nil),
			"abnormal_info_list": flattenAbnormalInfos(
				utils.PathSearch("abnormal_info_list", v, make([]interface{}, 0)).([]interface{})),
			"auto_confirm":   utils.PathSearch("auto_confirm", v, nil),
			"default_policy": utils.PathSearch("default_policy", v, nil),
			"host_id_list":   utils.PathSearch("host_id_list", v, nil),
		})
	}

	return rst
}

func flattenAbnormalInfos(abnormalInfos []interface{}) []interface{} {
	if len(abnormalInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(abnormalInfos))
	for _, v := range abnormalInfos {
		rst = append(rst, map[string]interface{}{
			"abnormal_type":        utils.PathSearch("abnormal_type", v, nil),
			"abnormal_description": utils.PathSearch("abnormal_description", v, nil),
		})
	}

	return rst
}
