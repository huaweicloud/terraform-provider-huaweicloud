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

// @API HSS GET /v5/{project_id}/cluster-protect/policy
func DataSourceClusterProtectPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterProtectPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"general_policy_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"malicious_image_policy_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_policy_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"images": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"labels": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"labels_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"white_images": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_version": {
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

func buildClusterProtectPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceClusterProtectPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		region                  = cfg.GetRegion(d)
		epsId                   = cfg.GetEnterpriseProjectID(d)
		product                 = "hss"
		httpUrl                 = "v5/{project_id}/cluster-protect/policy"
		dataListResult          = make([]interface{}, 0)
		offset                  = 0
		generalPolicyNum        float64
		maliciousImagePolicyNum float64
		securityPolicyNum       float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildClusterProtectPoliciesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		getResp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS cluster protect policies: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		generalPolicyNum = utils.PathSearch("general_policy_num", getRespBody, float64(0)).(float64)
		maliciousImagePolicyNum = utils.PathSearch("malicious_image_policy_num", getRespBody, float64(0)).(float64)
		securityPolicyNum = utils.PathSearch("security_policy_num", getRespBody, float64(0)).(float64)

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		dataListResult = append(dataListResult, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("general_policy_num", generalPolicyNum),
		d.Set("malicious_image_policy_num", maliciousImagePolicyNum),
		d.Set("security_policy_num", securityPolicyNum),
		d.Set("data_list", flattenClusterProtectPoliciesDataList(dataListResult)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterProtectPoliciesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"content":        utils.PathSearch("content", v, nil),
			"deploy_content": utils.PathSearch("deploy_content", v, nil),
			"parameters":     utils.PathSearch("parameters", v, nil),
			"policy_name":    utils.PathSearch("policy_name", v, nil),
			"policy_id":      utils.PathSearch("policy_id", v, nil),
			"resources":      flattenClusterProtectPoliciesResources(utils.PathSearch("resources", v, make([]interface{}, 0)).([]interface{})),
			"template_id":    utils.PathSearch("template_id", v, nil),
			"template_name":  utils.PathSearch("template_name", v, nil),
			"template_type":  utils.PathSearch("template_type", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"image_num":      utils.PathSearch("image_num", v, nil),
			"labels_num":     utils.PathSearch("labels_num", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"white_images":   flattenClusterProtectPoliciesWhiteImages(utils.PathSearch("white_images", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenClusterProtectPoliciesResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"cluster_id":   utils.PathSearch("cluster_id", v, nil),
			"cluster_name": utils.PathSearch("cluster_name", v, nil),
			"images":       utils.PathSearch("images", v, nil),
			"labels":       utils.PathSearch("labels", v, nil),
			"namespace":    utils.PathSearch("namespace", v, nil),
		})
	}

	return rst
}

func flattenClusterProtectPoliciesWhiteImages(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"cluster_id":    utils.PathSearch("cluster_id", v, nil),
			"image_name":    utils.PathSearch("image_name", v, nil),
			"image_version": utils.PathSearch("image_version", v, nil),
		})
	}

	return rst
}
