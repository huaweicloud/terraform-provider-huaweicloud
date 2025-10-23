package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW GET /v1/{project_id}/dew/cpcs/resource-info
func DataSourceResourceInfos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceInfosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_service": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_distribution": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"ccsp_service": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_distribution": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"encrypt_decrypt": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sign_verify": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"kms": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timestamp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"colla_sign": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"otp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"db_encrypt": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"file_encrypt": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"digit_seal": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ssl_vpn": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"vsm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpcs_cluster_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpcs_instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"app": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"kms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"result": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aes_256": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sm4": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rsa_2048": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rsa_3072": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rsa_4096": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ec_p256": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ec_p384": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sm2": {
										Type:     schema.TypeInt,
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

func dataSourceResourceInfosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/resource-info"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving resource distribution information: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
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
		d.Set("cloud_service", flattenCloudServiceInfo(utils.PathSearch("cloud_service", respBody, nil))),
		d.Set("ccsp_service", flattenCcspServiceInfo(utils.PathSearch("ccsp_service", respBody, nil))),
		d.Set("vsm", flattenVsmResourceInfo(utils.PathSearch("vsm", respBody, nil))),
		d.Set("app", flattenApplicationInfo(utils.PathSearch("app", respBody, nil))),
		d.Set("kms", flattenKmsResourceInfo(utils.PathSearch("kms", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCloudServiceInfo(cloudServiceInfo interface{}) []interface{} {
	if cloudServiceInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"service_num":           utils.PathSearch("service_num", cloudServiceInfo, nil),
		"resource_num":          utils.PathSearch("resource_num", cloudServiceInfo, nil),
		"resource_distribution": flattenResourceDistribution(utils.PathSearch("resource_distribution", cloudServiceInfo, nil)),
	}

	return []interface{}{result}
}

func flattenResourceDistribution(resourceDistribution interface{}) []interface{} {
	if resourceDistribution == nil {
		return nil
	}

	result := map[string]interface{}{
		"kms": utils.PathSearch("kms", resourceDistribution, nil),
	}

	return []interface{}{result}
}

func flattenCcspServiceInfo(ccspServiceInfo interface{}) []interface{} {
	if ccspServiceInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"cluster_num":           utils.PathSearch("cluster_num", ccspServiceInfo, nil),
		"instance_num":          utils.PathSearch("instance_num", ccspServiceInfo, nil),
		"instance_quota":        utils.PathSearch("instance_quota", ccspServiceInfo, nil),
		"instance_distribution": flattenInstanceDistribution(utils.PathSearch("instance_distribution", ccspServiceInfo, nil)),
	}

	return []interface{}{result}
}

func flattenInstanceDistribution(instanceDistribution interface{}) []interface{} {
	if instanceDistribution == nil {
		return nil
	}

	result := map[string]interface{}{
		"encrypt_decrypt": utils.PathSearch("encrypt_decrypt", instanceDistribution, nil),
		"sign_verify":     utils.PathSearch("sign_verify", instanceDistribution, nil),
		"kms":             utils.PathSearch("kms", instanceDistribution, nil),
		"timestamp":       utils.PathSearch("timestamp", instanceDistribution, nil),
		"colla_sign":      utils.PathSearch("colla_sign", instanceDistribution, nil),
		"otp":             utils.PathSearch("otp", instanceDistribution, nil),
		"db_encrypt":      utils.PathSearch("db_encrypt", instanceDistribution, nil),
		"file_encrypt":    utils.PathSearch("file_encrypt", instanceDistribution, nil),
		"digit_seal":      utils.PathSearch("digit_seal", instanceDistribution, nil),
		"ssl_vpn":         utils.PathSearch("ssl_vpn", instanceDistribution, nil),
	}

	return []interface{}{result}
}

func flattenVsmResourceInfo(vsmResourceInfo interface{}) []interface{} {
	if vsmResourceInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"cluster_num":       utils.PathSearch("cluster_num", vsmResourceInfo, nil),
		"cpcs_cluster_num":  utils.PathSearch("cpcs_cluster_num", vsmResourceInfo, nil),
		"instance_num":      utils.PathSearch("instance_num", vsmResourceInfo, nil),
		"cpcs_instance_num": utils.PathSearch("cpcs_instance_num", vsmResourceInfo, nil),
		"instance_quota":    utils.PathSearch("instance_quota", vsmResourceInfo, nil),
	}

	return []interface{}{result}
}

func flattenApplicationInfo(appInfo interface{}) []interface{} {
	if appInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"app_num": utils.PathSearch("app_num", appInfo, nil),
	}

	return []interface{}{result}
}

func flattenKmsResourceInfo(kmsResourceInfo interface{}) []interface{} {
	if kmsResourceInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"total_num": utils.PathSearch("total_num", kmsResourceInfo, nil),
		"result":    flattenKmsInfo(utils.PathSearch("result", kmsResourceInfo, nil)),
	}

	return []interface{}{result}
}

func flattenKmsInfo(resultInfo interface{}) []interface{} {
	if resultInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"aes_256":  utils.PathSearch("AES_256", resultInfo, nil),
		"sm4":      utils.PathSearch("sm4", resultInfo, nil),
		"rsa_2048": utils.PathSearch("rsa_2048", resultInfo, nil),
		"rsa_3072": utils.PathSearch("rsa_3072", resultInfo, nil),
		"rsa_4096": utils.PathSearch("rsa_4096", resultInfo, nil),
		"ec_p256":  utils.PathSearch("ec_p256", resultInfo, nil),
		"ec_p384":  utils.PathSearch("ec_p384", resultInfo, nil),
		"sm2":      utils.PathSearch("sm2", resultInfo, nil),
	}

	return []interface{}{result}
}
