package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var publicScriptExecuteNonUpdatableParams = []string{"script_uuid", "timeout", "success_rate", "execute_user",
	"script_params", "execute_batches"}

// @API COC POST /v1/job/public-scripts/{script_uuid}
func ResourcePublicScriptExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicScriptExecuteCreate,
		ReadContext:   resourcePublicScriptExecuteRead,
		UpdateContext: resourcePublicScriptExecuteUpdate,
		DeleteContext: resourcePublicScriptExecuteDelete,

		CustomizeDiff: config.FlexibleForceNew(publicScriptExecuteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"script_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"success_rate": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"execute_user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execute_batches": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_index": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"target_instances": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cloud_service_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_attributes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
						"rotation_strategy": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"script_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"param_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"param_refer": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"refer_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"param_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"param_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"execute_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPublicScriptExecuteCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"execute_param":   buildPublicScriptExecuteExecuteParamCreateOpts(d),
		"execute_batches": buildPublicScriptExecuteExecuteBatchesCreateOpts(d.Get("execute_batches")),
	}

	return bodyParams
}

func buildPublicScriptExecuteExecuteParamCreateOpts(d *schema.ResourceData) map[string]interface{} {
	param := map[string]interface{}{
		"timeout":       d.Get("timeout"),
		"success_rate":  d.Get("success_rate"),
		"execute_user":  d.Get("execute_user"),
		"script_params": buildPublicScriptExecuteScriptParamsCreateOpts(d.Get("script_params")),
	}

	return param
}

func buildPublicScriptExecuteScriptParamsCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"param_name":  raw["param_name"],
				"param_value": raw["param_value"],
				"param_refer": buildPublicScriptExecuteScriptParamsParamReferCreateOpts(raw["param_refer"]),
			}
		}
		return params
	}

	return nil
}

func buildPublicScriptExecuteScriptParamsParamReferCreateOpts(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"refer_type":    raw["refer_type"],
			"param_id":      raw["param_id"],
			"param_version": utils.ValueIgnoreEmpty(raw["param_version"]),
		}

		return param
	}

	return nil
}

func buildPublicScriptExecuteExecuteBatchesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"batch_index":       raw["batch_index"],
				"target_instances":  buildPublicScriptExecuteExecuteBatchesTargetInstancesCreateOpts(raw["target_instances"]),
				"rotation_strategy": raw["rotation_strategy"],
			}
		}
		return params
	}

	return nil
}

func buildPublicScriptExecuteExecuteBatchesTargetInstancesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"resource_id": raw["resource_id"],
				"region_id":   raw["region_id"],
				"provider":    raw["cloud_service_name"],
				"type":        raw["type"],
				"custom_attributes": buildPublicScriptExecuteExecuteBatchesTargetInstancesCustomAttributesCreateOpts(
					raw["custom_attributes"]),
			}
		}
		return params
	}

	return nil
}

func buildPublicScriptExecuteExecuteBatchesTargetInstancesCustomAttributesCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"key":   raw["key"],
				"value": raw["value"],
			}
		}
		return params
	}

	return nil
}

func resourcePublicScriptExecuteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/job/public-scripts/{script_uuid}"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	scriptUUID := d.Get("script_uuid").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{script_uuid}", scriptUUID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPublicScriptExecuteCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing the COC public script (%s): %s", scriptUUID, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening executing the COC public script response: %s", err)
	}

	executeUUID := utils.PathSearch("data", createRespBody, "").(string)
	if executeUUID == "" {
		return diag.Errorf("error executing the COC public script: can not find execute_uuid in return")
	}

	d.SetId(scriptUUID)

	mErr := multierror.Append(
		d.Set("execute_uuid", executeUUID),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePublicScriptExecuteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePublicScriptExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePublicScriptExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting public script execute resource is not supported. The public script execute resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
