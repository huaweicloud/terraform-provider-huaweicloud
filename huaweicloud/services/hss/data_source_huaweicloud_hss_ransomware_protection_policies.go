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

// @API HSS GET /v5/{project_id}/ransomware/protection/policy
func DataSourceRansomwareProtectionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRansomwareProtectionPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bait_protection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_detection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count_associated_server": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operating_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_whitelist": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_policy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRansomwareProtectionPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=20"
	queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	if v, ok := d.GetOk("policy_id"); ok {
		queryParams = fmt.Sprintf("%s&protect_policy_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&policy_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("operating_system"); ok {
		queryParams = fmt.Sprintf("%s&operating_system=%v", queryParams, v)
	}

	return queryParams
}

func datasourceRansomwareProtectionPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
		httpUrl = "v5/{project_id}/ransomware/protection/policy"
		offset  = 0
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRansomwareProtectionPoliciesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS ransomware protection policies: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		policiesResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(policiesResp) == 0 {
			break
		}

		result = append(result, policiesResp...)
		offset += len(policiesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenPolicies(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolicies(policiesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(policiesResp))
	for _, v := range policiesResp {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("policy_id", v, nil),
			"name":                     utils.PathSearch("policy_name", v, nil),
			"protection_mode":          utils.PathSearch("protection_mode", v, nil),
			"bait_protection_status":   utils.PathSearch("bait_protection_status", v, nil),
			"deploy_mode":              utils.PathSearch("deploy_mode", v, nil),
			"protection_directory":     utils.PathSearch("protection_directory", v, nil),
			"protection_type":          utils.PathSearch("protection_type", v, nil),
			"exclude_directory":        utils.PathSearch("exclude_directory", v, nil),
			"runtime_detection_status": utils.PathSearch("runtime_detection_status", v, nil),
			"count_associated_server":  utils.PathSearch("count_associated_server", v, nil),
			"operating_system":         utils.PathSearch("operating_system", v, nil),
			"process_whitelist":        flattenProcessWhitelist(utils.PathSearch("process_whitelist", v, make([]interface{}, 0)).([]interface{})),
			"default_policy":           utils.PathSearch("default_policy", v, nil),
		})
	}

	return rst
}

func flattenProcessWhitelist(processWhitelistResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(processWhitelistResp))
	for _, v := range processWhitelistResp {
		rst = append(rst, map[string]interface{}{
			"path": utils.PathSearch("path", v, nil),
			"hash": utils.PathSearch("hash", v, nil),
		})
	}

	return rst
}
