package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	modelartsAlgorithmNonUpdatableParams = []string{
		"metadata.*.name",
	}
	modelartsAlgorithmNotFoundErrCodes = []string{
		"ModelArts.2755", // The resource does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/algorithms
// @API ModelArts GET /v2/{project_id}/algorithms/{algorithm_id}
// @API ModelArts PUT /v2/{project_id}/algorithms/{algorithm_id}
// @API ModelArts DELETE /v2/{project_id}/algorithms/{algorithm_id}
func ResourceAlgorithm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlgorithmCreate,
		ReadContext:   resourceAlgorithmRead,
		UpdateContext: resourceAlgorithmUpdate,
		DeleteContext: resourceAlgorithmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(modelartsAlgorithmNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the algorithm is located.`,
			},

			// Required parameters.
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the algorithm.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the algorithm.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the workspace to which the algorithm belongs.`,
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The key of the tag.`,
									},
								},
							},
						},
					},
				},
				Description: `The metadata of the algorithm.`,
			},
			"job_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        algorithmJobConfigSchemaResource(),
				Description: `The configuration of the algorithm.`,
			},

			// Optional parameters.
			"resource_requirements": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The resource constraint list of the algorithm.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The key of the resource constraint.`,
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of values corresponding to the key.`,
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The relationship between the key and values.`,
						},
					},
				},
			},
			"advanced_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldRaw, newRaw := d.GetChange("advanced_config")
					newVal := newRaw.([]interface{})
					// When advanced_config is updated to empty, the API returns '{"advanced_config": {"auto_search": {}}}',
					// this scenario needs to suppress the change.
					return len(oldRaw.([]interface{})) == 0 && (len(newVal) == 0 || newVal[0] == nil)
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_search": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reward_attrs": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The metric name.`,
												},
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The search direction.`,
												},
												"regex": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The metric regular expression.`,
												},
											},
										},
										Description: `The metric list of search.`,
									},
									"search_params": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The hyperparameter name.`,
												},
												"param_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The parameter type.`,
												},
												"lower_bound": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The lower bound of the hyperparameter.`,
												},
												"upper_bound": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The upper bound of the hyperparameter.`,
												},
												"discrete_points_num": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The number of discrete samples for a continuous hyperparameter.`,
												},
												"discrete_values": {
													Type:        schema.TypeList,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: `The list of discrete values for a discrete hyperparameter.`,
												},
											},
										},
										Description: `The parameter list of search.`,
									},
									"algo_configs": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `The name of the search algorithm.`,
												},
												"params": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: `The key of the parameter.`,
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: `The value of the parameter.`,
															},
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: `The type of the parameter.`,
															},
														},
													},
													Description: `The parameter list of the search algorithm.`,
												},
											},
										},
										Description: `The algorithm configuration of search.`,
									},
									"skip_search_params": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The hyperparameter combination to be excluded from search.`,
									},
								},
							},
							Description: `The strategy configuration of hyperparameter search.`,
						},
					},
				},
				Description: `The advanced configuration of the algorithm.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func algorithmJobConfigSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the engine specification.`,
						},
						"engine_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the engine.`,
						},
						"engine_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The version of the engine version.`,
						},
						"image_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The custom image URL of the algorithm.`,
						},
					},
				},
				Description: `The engine configuration of the algorithm.`,
			},
			"code_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The code directory of the algorithm.`,
			},
			"boot_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The boot file path under the code directory.`,
			},
			"command": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The container start command for custom image algorithm.`,
			},
			"inputs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the data input channel.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the data input channel.`,
						},
						"access_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The access method of the data input channel.`,
						},
						"remote_constraints": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The type of the data input.`,
									},
									"attributes": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
										Description:  `The attributes of the input data, in JSON format.`,
									},
								},
							},
							Description: `The constraint of the data input.`,
						},
					},
				},
				Description: `The data input list of the algorithm.`,
			},
			"outputs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the data output channel.`,
						},
						"access_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The access method of the data output channel.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the data output channel.`,
						},
					},
				},
				Description: `The data output list of the algorithm.`,
			},
			"parameters_customization": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the hyperparameter can be customized when creating a training job.`,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the parameter.`,
						},
						"constraint": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the parameter.`,
									},
									"required": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether the parameter is required.`,
									},
									"editable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether the parameter is editable.`,
									},
									"valid_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The valid type of the parameter value.`,
									},
									"valid_range": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The valid range list of the parameter value.`,
									},
								},
							},
							Description: `The constraint of the parameter.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the parameter.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the parameter.`,
						},
						"i18n_description": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"language": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The language code.`,
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The description in the specified language.`,
									},
								},
							},
						},
					},
				},
				Description: `The list of running parameters of the algorithm.`,
			},
		},
	}
}

func buildCreateAlgorithm(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"metadata":              buildAlgorithmMetadata(d.Get("metadata").([]interface{})),
		"job_config":            buildAlgorithmJobConfig(d.Get("job_config.0")),
		"resource_requirements": buildAlgorithmResourceRequirements(d.Get("resource_requirements").([]interface{})),
		"advanced_config":       buildAlgorithmAdvancedConfig(d.Get("advanced_config").([]interface{})),
	}
}

func buildAlgorithmMetadata(metadata []interface{}) map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	return map[string]interface{}{
		"name":         utils.PathSearch("name", metadata[0], nil),
		"description":  utils.ValueIgnoreEmpty(utils.PathSearch("description", metadata[0], nil)),
		"workspace_id": utils.ValueIgnoreEmpty(utils.PathSearch("workspace_id", metadata[0], nil)),
		"tags":         buildAlgorithmMetadataTags(utils.PathSearch("tags", metadata[0], make([]interface{}, 0)).([]interface{})),
	}
}

func buildAlgorithmMetadataTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return make([]map[string]interface{}, 0)
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, v := range tags {
		result = append(result, map[string]interface{}{
			"key": utils.PathSearch("key", v, nil),
		})
	}

	return result
}

func buildAlgorithmJobConfig(jobConfig interface{}) map[string]interface{} {
	return map[string]interface{}{
		"engine":    buildAlgorithmJobConfigEngine(utils.PathSearch("engine|[0]", jobConfig, nil)),
		"code_dir":  utils.ValueIgnoreEmpty(utils.PathSearch("code_dir", jobConfig, nil)),
		"boot_file": utils.ValueIgnoreEmpty(utils.PathSearch("boot_file", jobConfig, nil)),
		"command":   utils.ValueIgnoreEmpty(utils.PathSearch("command", jobConfig, nil)),
		"inputs": buildAlgorithmJobConfigInputs(utils.PathSearch("inputs",
			jobConfig, make([]interface{}, 0)).([]interface{})),
		"outputs": buildAlgorithmJobConfigOutputs(utils.PathSearch("outputs",
			jobConfig, make([]interface{}, 0)).([]interface{})),
		"parameters_customization": utils.PathSearch("parameters_customization", jobConfig, nil),
		"parameters": buildAlgorithmJobConfigParameters(utils.PathSearch("parameters",
			jobConfig, make([]interface{}, 0)).([]interface{})),
	}
}

func buildAlgorithmJobConfigEngine(engine interface{}) map[string]interface{} {
	return map[string]interface{}{
		"engine_id":      utils.ValueIgnoreEmpty(utils.PathSearch("engine_id", engine, nil)),
		"engine_name":    utils.ValueIgnoreEmpty(utils.PathSearch("engine_name", engine, nil)),
		"engine_version": utils.ValueIgnoreEmpty(utils.PathSearch("engine_version", engine, nil)),
		"image_url":      utils.ValueIgnoreEmpty(utils.PathSearch("image_url", engine, nil)),
	}
}

func buildAlgorithmJobConfigInputs(inputs []interface{}) []map[string]interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, v := range inputs {
		result = append(result, map[string]interface{}{
			"name":          utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"description":   utils.ValueIgnoreEmpty(utils.PathSearch("description", v, nil)),
			"access_method": utils.ValueIgnoreEmpty(utils.PathSearch("access_method", v, nil)),
			"remote_constraints": buildAlgorithmJobConfigInputRemoteConstraints(utils.PathSearch("remote_constraints",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func buildAlgorithmJobConfigInputRemoteConstraints(constraints []interface{}) []map[string]interface{} {
	if len(constraints) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(constraints))
	for _, v := range constraints {
		result = append(result, map[string]interface{}{
			"data_type":  utils.ValueIgnoreEmpty(utils.PathSearch("data_type", v, nil)),
			"attributes": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("attributes", v, "").(string))),
		})
	}

	return result
}

func buildAlgorithmJobConfigOutputs(outputs []interface{}) []map[string]interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(outputs))
	for _, v := range outputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"access_method": utils.ValueIgnoreEmpty(utils.PathSearch("access_method", v, nil)),
			"description":   utils.ValueIgnoreEmpty(utils.PathSearch("description", v, nil)),
		})
	}

	return result
}

func buildAlgorithmJobConfigParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, v := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"value":       utils.ValueIgnoreEmpty(utils.PathSearch("value", v, nil)),
			"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", v, nil)),
			"constraint": buildAlgorithmJobConfigParametersConstraint(utils.PathSearch("constraint[0]",
				v, map[string]interface{}{}).(map[string]interface{})),
			"i18n_description": buildAlgorithmJobConfigParametersI18nDescription(utils.PathSearch("i18n_description[0]",
				v, map[string]interface{}{}).(map[string]interface{})),
		})
	}
	return result
}

func buildAlgorithmJobConfigParametersConstraint(constraint map[string]interface{}) map[string]interface{} {
	if len(constraint) < 1 {
		return nil
	}

	return map[string]interface{}{
		"type":       utils.ValueIgnoreEmpty(utils.PathSearch("type", constraint, nil)),
		"editable":   utils.ValueIgnoreEmpty(utils.PathSearch("editable", constraint, nil)),
		"required":   utils.ValueIgnoreEmpty(utils.PathSearch("required", constraint, nil)),
		"valid_type": utils.ValueIgnoreEmpty(utils.PathSearch("valid_type", constraint, nil)),
		"valid_range": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("valid_range",
			constraint, make([]interface{}, 0)).([]interface{}))),
	}
}

func buildAlgorithmJobConfigParametersI18nDescription(i18nDescription map[string]interface{}) map[string]interface{} {
	if len(i18nDescription) < 1 {
		return nil
	}

	return map[string]interface{}{
		"language":    utils.ValueIgnoreEmpty(utils.PathSearch("language", i18nDescription, nil)),
		"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", i18nDescription, nil)),
	}
}

func buildAlgorithmResourceRequirements(requirements []interface{}) []map[string]interface{} {
	if len(requirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(requirements))
	for _, v := range requirements {
		result = append(result, map[string]interface{}{
			"key": utils.ValueIgnoreEmpty(utils.PathSearch("key", v, nil)),
			"values": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("values",
				v, make([]interface{}, 0)).([]interface{}))),
			"operator": utils.ValueIgnoreEmpty(utils.PathSearch("operator", v, nil)),
		})
	}

	return result
}

func buildAlgorithmAdvancedConfig(advancedConfigs []interface{}) map[string]interface{} {
	autoSearch := utils.PathSearch("[0].auto_search[0]", advancedConfigs, map[string]interface{}{}).(map[string]interface{})
	if len(autoSearch) < 1 {
		return nil
	}

	return map[string]interface{}{
		"auto_search": buildAlgorithmAdvancedConfigAutoSearch(autoSearch),
	}
}

func buildAlgorithmAdvancedConfigAutoSearch(autoSearch map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"reward_attrs": buildAlgorithmAdvancedConfigAutoSearchRewardAttrs(utils.PathSearch("reward_attrs",
			autoSearch, make([]interface{}, 0)).([]interface{})),
		"search_params": buildAlgorithmAdvancedConfigAutoSearchSearchParams(utils.PathSearch("search_params",
			autoSearch, make([]interface{}, 0)).([]interface{})),
		"algo_configs": buildAlgorithmAdvancedConfigAutoSearchAlgoConfigs(utils.PathSearch("algo_configs",
			autoSearch, make([]interface{}, 0)).([]interface{})),
		"skip_search_params": utils.ValueIgnoreEmpty(utils.PathSearch("skip_search_params", autoSearch, nil)),
	}
}

func buildAlgorithmAdvancedConfigAutoSearchRewardAttrs(rewardAttrs []interface{}) []map[string]interface{} {
	if len(rewardAttrs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rewardAttrs))
	for _, v := range rewardAttrs {
		result = append(result, map[string]interface{}{
			"name":  utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"mode":  utils.ValueIgnoreEmpty(utils.PathSearch("mode", v, nil)),
			"regex": utils.ValueIgnoreEmpty(utils.PathSearch("regex", v, nil)),
		})
	}
	return result
}

func buildAlgorithmAdvancedConfigAutoSearchSearchParams(searchParams []interface{}) []map[string]interface{} {
	if len(searchParams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(searchParams))
	for _, v := range searchParams {
		result = append(result, map[string]interface{}{
			"name":                utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"param_type":          utils.ValueIgnoreEmpty(utils.PathSearch("param_type", v, nil)),
			"lower_bound":         utils.ValueIgnoreEmpty(utils.PathSearch("lower_bound", v, nil)),
			"upper_bound":         utils.ValueIgnoreEmpty(utils.PathSearch("upper_bound", v, nil)),
			"discrete_points_num": utils.ValueIgnoreEmpty(utils.PathSearch("discrete_points_num", v, nil)),
			"discrete_values": utils.ValueIgnoreEmpty(utils.ExpandToStringList(
				utils.PathSearch("discrete_values", v, make([]interface{}, 0)).([]interface{}))),
		})
	}

	return result
}

func buildAlgorithmAdvancedConfigAutoSearchAlgoConfigs(algoConfigs []interface{}) []map[string]interface{} {
	if len(algoConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(algoConfigs))
	for _, v := range algoConfigs {
		result = append(result, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"params": buildAlgorithmAdvancedConfigAutoSearchAlgoConfigsParams(utils.PathSearch("params",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func buildAlgorithmAdvancedConfigAutoSearchAlgoConfigsParams(params []interface{}) []map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(params))
	for _, v := range params {
		result = append(result, map[string]interface{}{
			"key":   utils.ValueIgnoreEmpty(utils.PathSearch("key", v, nil)),
			"value": utils.ValueIgnoreEmpty(utils.PathSearch("value", v, nil)),
			"type":  utils.ValueIgnoreEmpty(utils.PathSearch("type", v, nil)),
		})
	}
	return result
}

func createAlgorithm(client *golangsdk.ServiceClient, params map[string]interface{}) (interface{}, error) {
	httpURL := "v2/{project_id}/algorithms"
	createPath := client.Endpoint + httpURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(params),
	}
	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAlgorithmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createAlgorithm(client, buildCreateAlgorithm(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts algorithm: %s", err)
	}

	algorithmId := utils.PathSearch("metadata.id", resp, "").(string)
	if algorithmId == "" {
		return diag.Errorf("unable to find the algorithm ID from the API response")
	}

	d.SetId(algorithmId)

	return resourceAlgorithmRead(ctx, d, meta)
}

func GetAlgorithmById(client *golangsdk.ServiceClient, algorithmId string) (interface{}, error) {
	httpURL := "v2/{project_id}/algorithms/{algorithm_id}"
	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{algorithm_id}", algorithmId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAlgorithmRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		algorithmId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetAlgorithmById(client, algorithmId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", modelartsAlgorithmNotFoundErrCodes...),
			fmt.Sprintf("error retrieving ModelArts algorithm (%s)", algorithmId),
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("metadata", flattenAlgorithmMetadata(utils.PathSearch("metadata", resp, nil))),
		d.Set("job_config", flattenAlgorithmJobConfig(utils.PathSearch("job_config", resp, nil))),
		d.Set("resource_requirements", flattenAlgorithmResourceRequirements(utils.PathSearch("resource_requirements",
			resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("advanced_config", flattenAlgorithmAdvancedConfig(utils.PathSearch("advanced_config", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlgorithmMetadata(metadata interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":         utils.PathSearch("name", metadata, nil),
			"description":  utils.PathSearch("description", metadata, nil),
			"workspace_id": utils.PathSearch("workspace_id", metadata, nil),
			"tags":         flattenAlgorithmMetadataTags(utils.PathSearch("tags", metadata, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenAlgorithmMetadataTags(tags []interface{}) []interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, v := range tags {
		result = append(result, map[string]interface{}{
			"key": utils.PathSearch("key", v, nil),
		})
	}

	return result
}

func flattenAlgorithmJobConfig(jobConfig interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"engine":    flattenAlgorithmJobConfigEngine(utils.PathSearch("engine", jobConfig, nil)),
			"code_dir":  utils.PathSearch("code_dir", jobConfig, nil),
			"boot_file": utils.PathSearch("boot_file", jobConfig, nil),
			"command":   utils.PathSearch("command", jobConfig, nil),
			"inputs": flattenAlgorithmJobConfigInputs(utils.PathSearch("inputs",
				jobConfig, make([]interface{}, 0)).([]interface{})),
			"outputs": flattenAlgorithmJobConfigOutputs(utils.PathSearch("outputs",
				jobConfig, make([]interface{}, 0)).([]interface{})),
			"parameters_customization": utils.PathSearch("parameters_customization", jobConfig, nil),
			"parameters": flattenAlgorithmJobConfigParameters(utils.PathSearch("parameters",
				jobConfig, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenAlgorithmJobConfigEngine(engine interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"engine_id":      utils.PathSearch("engine_id", engine, nil),
			"engine_name":    utils.PathSearch("engine_name", engine, nil),
			"engine_version": utils.PathSearch("engine_version", engine, nil),
			"image_url":      utils.PathSearch("image_url", engine, nil),
		},
	}
}

func flattenAlgorithmJobConfigInputs(inputs []interface{}) []interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(inputs))
	for _, v := range inputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"access_method": utils.PathSearch("access_method", v, nil),
			"remote_constraints": flattenAlgorithmJobConfigInputRemoteConstraints(utils.PathSearch("remote_constraints",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenAlgorithmJobConfigInputRemoteConstraints(remoteConstraints []interface{}) []interface{} {
	if len(remoteConstraints) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(remoteConstraints))
	for _, v := range remoteConstraints {
		result = append(result, map[string]interface{}{
			"data_type":  utils.PathSearch("data_type", v, nil),
			"attributes": utils.JsonToString(utils.PathSearch("attributes", v, nil)),
		})
	}

	return result
}

func flattenAlgorithmJobConfigOutputs(outputs []interface{}) []interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(outputs))
	for _, v := range outputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"access_method": utils.PathSearch("access_method", v, nil),
			"description":   utils.PathSearch("description", v, nil),
		})
	}

	return result
}

func flattenAlgorithmJobConfigParameters(parameters []interface{}) []interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(parameters))
	for _, v := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"value":       utils.PathSearch("value", v, nil),
			"constraint": flattenAlgorithmJobConfigParametersConstraint(utils.PathSearch("constraint",
				v, map[string]interface{}{}).(map[string]interface{})),
			"i18n_description": flattenAlgorithmJobConfigParametersI18nDescription(utils.PathSearch("i18n_description[0]",
				v, map[string]interface{}{}).(map[string]interface{})),
		})
	}

	return result
}

func flattenAlgorithmJobConfigParametersConstraint(constraint interface{}) []interface{} {
	if constraint == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":        utils.PathSearch("type", constraint, nil),
			"editable":    utils.PathSearch("editable", constraint, nil),
			"required":    utils.PathSearch("required", constraint, nil),
			"valid_type":  utils.PathSearch("valid_type", constraint, nil),
			"valid_range": utils.PathSearch("valid_range", constraint, nil),
		},
	}
}

func flattenAlgorithmJobConfigParametersI18nDescription(i18nDescription map[string]interface{}) []interface{} {
	if len(i18nDescription) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"language":    utils.PathSearch("language", i18nDescription, nil),
			"description": utils.PathSearch("description", i18nDescription, nil),
		},
	}
}

func flattenAlgorithmResourceRequirements(requirements []interface{}) []interface{} {
	if len(requirements) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(requirements))
	for _, v := range requirements {
		result = append(result, map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"values":   utils.PathSearch("values", v, nil),
			"operator": utils.PathSearch("operator", v, nil),
		})
	}
	return result
}

func flattenAlgorithmAdvancedConfig(advancedConfig interface{}) []interface{} {
	autoSearch := utils.PathSearch("auto_search", advancedConfig, make(map[string]interface{}, 0)).(map[string]interface{})
	if len(autoSearch) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"auto_search": flattenAlgorithmAdvancedConfigAutoSearch(autoSearch),
		},
	}
}

func flattenAlgorithmAdvancedConfigAutoSearch(autoSearch map[string]interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"skip_search_params": utils.PathSearch("skip_search_params", autoSearch, nil),
			"reward_attrs": flattenAlgorithmAdvancedConfigAutoSearchRewardAttrs(utils.PathSearch("reward_attrs",
				autoSearch, make([]interface{}, 0)).([]interface{})),
			"search_params": flattenAlgorithmAdvancedConfigAutoSearchSearchParams(utils.PathSearch("search_params",
				autoSearch, make([]interface{}, 0)).([]interface{})),
			"algo_configs": flattenAlgorithmAdvancedConfigAutoSearchAlgoConfigs(utils.PathSearch("algo_configs",
				autoSearch, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenAlgorithmAdvancedConfigAutoSearchRewardAttrs(rewardAttrs []interface{}) []interface{} {
	if len(rewardAttrs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rewardAttrs))
	for _, v := range rewardAttrs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"mode":  utils.PathSearch("mode", v, nil),
			"regex": utils.PathSearch("regex", v, nil),
		})
	}
	return result
}

func flattenAlgorithmAdvancedConfigAutoSearchSearchParams(searchParams []interface{}) []interface{} {
	if len(searchParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(searchParams))
	for _, v := range searchParams {
		result = append(result, map[string]interface{}{
			"name":                utils.PathSearch("name", v, nil),
			"param_type":          utils.PathSearch("param_type", v, nil),
			"lower_bound":         utils.PathSearch("lower_bound", v, nil),
			"upper_bound":         utils.PathSearch("upper_bound", v, nil),
			"discrete_points_num": utils.PathSearch("discrete_points_num", v, nil),
			"discrete_values":     utils.PathSearch("discrete_values", v, nil),
		})
	}
	return result
}

func flattenAlgorithmAdvancedConfigAutoSearchAlgoConfigs(algoConfigs []interface{}) []interface{} {
	if len(algoConfigs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(algoConfigs))
	for _, v := range algoConfigs {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"params": flattenAlgorithmAdvancedConfigAutoSearchAlgoConfigsParams(utils.PathSearch("params",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenAlgorithmAdvancedConfigAutoSearchAlgoConfigsParams(params []interface{}) []interface{} {
	if len(params) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(params))
	for _, v := range params {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
			"type":  utils.PathSearch("type", v, nil),
		})
	}

	return result
}

func updateAlgorithm(client *golangsdk.ServiceClient, algorithmId string, params map[string]interface{}) error {
	httpURL := "v2/{project_id}/algorithms/{algorithm_id}"
	updatePath := client.Endpoint + httpURL
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{algorithm_id}", algorithmId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: params,
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func buildUpdateAlgorithm(d *schema.ResourceData) map[string]interface{} {
	params := utils.RemoveNil(map[string]interface{}{
		"job_config":            buildUpdateAlgorithmJobConfig(d.Get("job_config.0")),
		"resource_requirements": buildAlgorithmResourceRequirements(d.Get("resource_requirements").([]interface{})),
		"advanced_config":       buildAlgorithmAdvancedConfig(d.Get("advanced_config").([]interface{})),
	})

	params["metadata"] = buildUpdateAlgorithmMetadata(d.Get("metadata").([]interface{}))

	if resourceRequirements := utils.PathSearch("resource_requirements", params, nil); resourceRequirements == nil {
		params["resource_requirements"] = make([]interface{}, 0)
	}

	if advancedConfig := utils.PathSearch("advanced_config", params, nil); advancedConfig == nil {
		params["advanced_config"] = map[string]interface{}{
			"auto_search": map[string]interface{}{
				"reward_attrs":  make([]interface{}, 0),
				"search_params": make([]interface{}, 0),
				"algo_configs":  make([]interface{}, 0),
			},
		}
	}

	return params
}

func buildUpdateAlgorithmMetadata(metadata []interface{}) map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	params := utils.RemoveNil(map[string]interface{}{
		"name":         utils.PathSearch("name", metadata[0], nil),
		"description":  utils.ValueIgnoreEmpty(utils.PathSearch("description", metadata[0], nil)),
		"workspace_id": utils.ValueIgnoreEmpty(utils.PathSearch("workspace_id", metadata[0], nil)),
	})

	// API requires: tags can only be updated to empty when it is specified as an empty list
	params["tags"] = buildAlgorithmMetadataTags(utils.PathSearch("tags", metadata[0], make([]interface{}, 0)).([]interface{}))

	return params
}

func buildUpdateAlgorithmJobConfig(jobConfig interface{}) map[string]interface{} {
	jobConfigMap := utils.RemoveNil(buildAlgorithmJobConfig(jobConfig))
	if inputs := utils.PathSearch("inputs", jobConfigMap, nil); inputs == nil {
		jobConfigMap["inputs"] = make([]interface{}, 0)
	}

	if outputs := utils.PathSearch("outputs", jobConfigMap, nil); outputs == nil {
		jobConfigMap["outputs"] = make([]interface{}, 0)
	}

	if parameters := utils.PathSearch("parameters", jobConfigMap, nil); parameters == nil {
		jobConfigMap["parameters"] = make([]interface{}, 0)
	}

	return jobConfigMap
}

func resourceAlgorithmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if err := updateAlgorithm(client, d.Id(), buildUpdateAlgorithm(d)); err != nil {
		return diag.Errorf("error updating ModelArts algorithm (%s): %s", d.Id(), err)
	}

	return resourceAlgorithmRead(ctx, d, meta)
}

func deleteAlgorithm(client *golangsdk.ServiceClient, algorithmId string) error {
	httpURL := "v2/{project_id}/algorithms/{algorithm_id}"
	deletePath := client.Endpoint + httpURL
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{algorithm_id}", algorithmId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceAlgorithmDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteAlgorithm(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", modelartsAlgorithmNotFoundErrCodes...),
			fmt.Sprintf("error deleting ModelArts algorithm (%s)", d.Id()),
		)
	}
	return nil
}
