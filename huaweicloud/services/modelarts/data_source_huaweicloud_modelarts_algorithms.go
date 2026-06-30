package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/algorithms
func DataSourceAlgorithms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlgorithmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the algorithms are located.",
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the workspace to which the algorithms belong.",
			},
			"searches": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filter condition for searching algorithms.",
			},
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The field by which the algorithms are sorted.",
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The sort order of the algorithms.",
			},

			// Attributes.
			"algorithms": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsSchema(),
				Description: "The list of algorithms that matched filter parameters.",
			},
		},
	}
}

func algorithmsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsMetadataSchema(),
				Description: "The metadata configuration of the algorithm.",
			},
			"job_config": {
				Type:        schema.TypeList,
				Elem:        algorithmsJobConfigSchema(),
				Computed:    true,
				Description: "The configuration of the algorithm.",
			},
			"resource_requirements": {
				Type:        schema.TypeList,
				Elem:        algorithmsResourceRequirementsSchema(),
				Computed:    true,
				Description: "The resource constraint list of the algorithm.",
			},
			"advanced_config": {
				Type:        schema.TypeList,
				Elem:        algorithmsAdvancedConfigSchema(),
				Computed:    true,
				Description: "The advanced configuration of the algorithm.",
			},
		},
	}
}

func algorithmsMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the algorithm.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the algorithm.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the algorithm.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the workspace to which the algorithm belongs.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name of the algorithm.",
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source type of the algorithm.",
			},
			"is_valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The availability of the algorithm.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the algorithm.",
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the tag.`,
						},
					},
				},
				Description: "The tags of the algorithm.",
			},
			"attr_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The attribute list of the algorithm.",
			},
			"version_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number of the algorithm.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the algorithm.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the algorithm, in RFC3339 format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the algorithm, in RFC3339 format.",
			},
		},
	}
}

func algorithmsJobConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The code directory of the algorithm.",
			},
			"boot_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The boot file of the algorithm.",
			},
			"command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container start command for custom image algorithm.",
			},
			"parameters_customization": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the hyperparameter can be customized.",
			},
			"parameters": {
				Type:        schema.TypeList,
				Elem:        algorithmsJobConfigParametersSchema(),
				Computed:    true,
				Description: "The running parameters of the algorithm.",
			},
			"inputs": {
				Type:        schema.TypeList,
				Elem:        algorithmsJobConfigInputsSchema(),
				Computed:    true,
				Description: "The data input list of the algorithm.",
			},
			"outputs": {
				Type:        schema.TypeList,
				Elem:        algorithmsJobConfigOutputsSchema(),
				Computed:    true,
				Description: "The data output list of the algorithm.",
			},
			"engine": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsJobConfigEngineSchema(),
				Description: "The engine of the algorithm.",
			},
		},
	}
}

func algorithmsJobConfigParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the parameter.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the parameter.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the parameter.",
			},
		},
	}
}

func algorithmsJobConfigInputsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the data input channel.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the data input channel.",
			},
		},
	}
}

func algorithmsJobConfigOutputsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the data output channel.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the data output channel.",
			},
		},
	}
}

func algorithmsJobConfigEngineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the engine.",
			},
			"engine_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the engine.",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the engine.",
			},
			"image_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The custom image URL of the algorithm.",
			},
		},
	}
}

func algorithmsResourceRequirementsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of the resource constraint.",
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The values corresponding to the key.",
			},
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relationship between the key and values.",
			},
		},
	}
}

func algorithmsAdvancedConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_search": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsAdvancedConfigAutoSearchSchema(),
				Description: "The hyperparameter search strategy.",
			},
		},
	}
}

func algorithmsAdvancedConfigAutoSearchSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"skip_search_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hyperparameter combinations to be excluded from search.",
			},
			"reward_attrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsAdvancedConfigAutoSearchRewardAttrsSchema(),
				Description: "The metric list of search.",
			},
			"search_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsAdvancedConfigAutoSearchSearchParamsSchema(),
				Description: "The search parameters.",
			},
			"algo_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsAdvancedConfigAutoSearchAlgoConfigsSchema(),
				Description: "The search algorithm configurations.",
			},
		},
	}
}

func algorithmsAdvancedConfigAutoSearchRewardAttrsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The metric name.",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The search direction.",
			},
			"regex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The metric regular expression.",
			},
		},
	}
}

func algorithmsAdvancedConfigAutoSearchSearchParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hyperparameter name.",
			},
			"param_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parameter type.",
			},
			"lower_bound": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lower bound of the hyperparameter.",
			},
			"upper_bound": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The upper bound of the hyperparameter.",
			},
			"discrete_points_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The number of discrete samples for continuous hyperparameters.",
			},
			"discrete_values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The discrete values for discrete hyperparameters.",
			},
		},
	}
}

func algorithmsAdvancedConfigAutoSearchAlgoConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The search algorithm name.",
			},
			"params": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        algorithmsAdvancedConfigAutoSearchAlgoConfigsParamsSchema(),
				Description: "The search algorithm parameters.",
			},
		},
	}
}

func algorithmsAdvancedConfigAutoSearchAlgoConfigsParamsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of the parameter.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the parameter.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the parameter.",
			},
		},
	}
}

func buildAlgorithmsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	if v, ok := d.GetOk("searches"); ok {
		res = fmt.Sprintf("%s&searches=%v", res, v)
	}

	if v, ok := d.GetOk("sort_by"); ok {
		res = fmt.Sprintf("%s&sort_by=%v", res, v)
	}

	if v, ok := d.GetOk("order"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}

	return res
}

func listAlgorithms(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/algorithms"
		// Maximum is 50.
		limit    = 50
		pageSize = 0
		result   = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?limit=%d%s", limit, buildAlgorithmsQueryParams(d))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPageSize := listPath + fmt.Sprintf("&offset=%d", pageSize)
		requestResp, err := client.Request("GET", listPathWithPageSize, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) < 1 {
			break
		}

		result = append(result, items...)
		if len(items) < limit {
			break
		}

		pageSize++
	}

	return result, nil
}

func flattenAlgorithms(algorithms []interface{}) []interface{} {
	if len(algorithms) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(algorithms))
	for _, v := range algorithms {
		rst = append(rst, map[string]interface{}{
			"metadata":   flattenAlgorithmsMetadata(utils.PathSearch("metadata", v, nil)),
			"job_config": flattenAlgorithmsJobConfig(utils.PathSearch("job_config", v, nil)),
			"resource_requirements": flattenAlgorithmsResourceRequirements(utils.PathSearch("resource_requirements",
				v, make([]interface{}, 0)).([]interface{})),
			"advanced_config": flattenAlgorithmsAdvancedConfig(utils.PathSearch("advanced_config", v, nil)),
		})
	}
	return rst
}

func flattenAlgorithmsMetadata(metadata interface{}) []interface{} {
	if metadata == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":           utils.PathSearch("id", metadata, nil),
			"name":         utils.PathSearch("name", metadata, nil),
			"description":  utils.PathSearch("description", metadata, nil),
			"workspace_id": utils.PathSearch("workspace_id", metadata, nil),
			"user_name":    utils.PathSearch("user_name", metadata, nil),
			"source":       utils.PathSearch("source", metadata, nil),
			"is_valid":     utils.PathSearch("is_valid", metadata, nil),
			"state":        utils.PathSearch("state", metadata, nil),
			"tags": flattenAlgorithmsMetadataTags(utils.PathSearch("tags",
				metadata, make([]interface{}, 0)).([]interface{})),
			"attr_list":   utils.PathSearch("attr_list", metadata, nil),
			"version_num": utils.PathSearch("version_num", metadata, nil),
			"size":        utils.PathSearch("size", metadata, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				metadata, float64(0)).(float64))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				metadata, float64(0)).(float64))/1000, false),
		},
	}
}

func flattenAlgorithmsMetadataTags(tags []interface{}) []map[string]interface{} {
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

func flattenAlgorithmsJobConfig(jobConfig interface{}) []interface{} {
	if jobConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"code_dir":                 utils.PathSearch("code_dir", jobConfig, nil),
			"boot_file":                utils.PathSearch("boot_file", jobConfig, nil),
			"command":                  utils.PathSearch("command", jobConfig, nil),
			"parameters_customization": utils.PathSearch("parameters_customization", jobConfig, nil),
			"parameters": flattenAlgorithmsJobConfigParameters(
				utils.PathSearch("parameters", jobConfig, make([]interface{}, 0)).([]interface{})),
			"inputs": flattenAlgorithmsJobConfigInputs(
				utils.PathSearch("inputs", jobConfig, make([]interface{}, 0)).([]interface{})),
			"outputs": flattenAlgorithmsJobConfigOutputs(
				utils.PathSearch("outputs", jobConfig, make([]interface{}, 0)).([]interface{})),
			"engine": flattenAlgorithmsJobConfigEngine(
				utils.PathSearch("engine", jobConfig, nil)),
		},
	}
}

func flattenAlgorithmsJobConfigParameters(parameters []interface{}) []interface{} {
	if len(parameters) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(parameters))
	for _, v := range parameters {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"value":       utils.PathSearch("value", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsJobConfigInputs(inputs []interface{}) []interface{} {
	if len(inputs) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(inputs))
	for _, v := range inputs {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsJobConfigOutputs(outputs []interface{}) []interface{} {
	if len(outputs) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(outputs))
	for _, v := range outputs {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsJobConfigEngine(engine interface{}) []interface{} {
	if engine == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"engine_id":      utils.PathSearch("engine_id", engine, nil),
			"engine_name":    utils.PathSearch("engine_name", engine, nil),
			"engine_version": utils.PathSearch("engine_version", engine, nil),
			"image_url":      utils.PathSearch("image_url", engine, nil),
		},
	}
}

func flattenAlgorithmsResourceRequirements(requirements []interface{}) []interface{} {
	if len(requirements) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(requirements))
	for _, v := range requirements {
		rst = append(rst, map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"values":   utils.PathSearch("values", v, nil),
			"operator": utils.PathSearch("operator", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsAdvancedConfig(advancedConfig interface{}) []interface{} {
	if advancedConfig == nil {
		return nil
	}

	autoSearch := utils.PathSearch("auto_search", advancedConfig, nil)
	if autoSearch == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"auto_search": flattenAlgorithmsAdvancedConfigAutoSearch(autoSearch),
		},
	}
}

func flattenAlgorithmsAdvancedConfigAutoSearch(autoSearch interface{}) []interface{} {
	if autoSearch == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"skip_search_params": utils.PathSearch("skip_search_params", autoSearch, nil),
			"reward_attrs": flattenAlgorithmsAdvancedConfigAutoSearchRewardAttrs(
				utils.PathSearch("reward_attrs", autoSearch, make([]interface{}, 0)).([]interface{})),
			"search_params": flattenAlgorithmsAdvancedConfigAutoSearchSearchParams(
				utils.PathSearch("search_params", autoSearch, make([]interface{}, 0)).([]interface{})),
			"algo_configs": flattenAlgorithmsAdvancedConfigAutoSearchAlgoConfigs(
				utils.PathSearch("algo_configs", autoSearch, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenAlgorithmsAdvancedConfigAutoSearchRewardAttrs(rewardAttrs []interface{}) []interface{} {
	if len(rewardAttrs) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(rewardAttrs))
	for _, v := range rewardAttrs {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"mode":  utils.PathSearch("mode", v, nil),
			"regex": utils.PathSearch("regex", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsAdvancedConfigAutoSearchSearchParams(searchParams []interface{}) []interface{} {
	if len(searchParams) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(searchParams))
	for _, v := range searchParams {
		rst = append(rst, map[string]interface{}{
			"name":                utils.PathSearch("name", v, nil),
			"param_type":          utils.PathSearch("param_type", v, nil),
			"lower_bound":         utils.PathSearch("lower_bound", v, nil),
			"upper_bound":         utils.PathSearch("upper_bound", v, nil),
			"discrete_points_num": utils.PathSearch("discrete_points_num", v, nil),
			"discrete_values":     utils.PathSearch("discrete_values", v, nil),
		})
	}
	return rst
}

func flattenAlgorithmsAdvancedConfigAutoSearchAlgoConfigs(algoConfigs []interface{}) []interface{} {
	if len(algoConfigs) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(algoConfigs))
	for _, v := range algoConfigs {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"params": flattenAlgorithmsAdvancedConfigAutoSearchAlgoConfigsParams(
				utils.PathSearch("params", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenAlgorithmsAdvancedConfigAutoSearchAlgoConfigsParams(params []interface{}) []interface{} {
	if len(params) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(params))
	for _, v := range params {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
			"type":  utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func dataSourceAlgorithmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	algorithms, err := listAlgorithms(client, d)
	if err != nil {
		return diag.Errorf("error querying ModelArts algorithms: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("algorithms", flattenAlgorithms(algorithms)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
