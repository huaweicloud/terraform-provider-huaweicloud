// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v3/{project_id}/templates
// @API DLI GET /v3/{project_id}/templates/{id}
// @API DLI PUT /v3/{project_id}/templates/{id}
// @API DLI POST /v1.0/{project_id}/sqls-deletion
func ResourceSparkTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSparkTemplateCreate,
		UpdateContext: resourceSparkTemplateUpdate,
		ReadContext:   resourceSparkTemplateRead,
		DeleteContext: resourceSparkTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the spark template.`,
			},
			"body": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        SparkTemplateBodySchema(),
				Required:    true,
				Description: `The content of the spark template.`,
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The group of the spark template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the spark template.`,
			},
		},
	}
}

func SparkTemplateBodySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The DLI queue name.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The spark job name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Name of the package that is of the JAR or pyFile type.`,
			},
			"main_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Java/Spark main class of the template.`,
			},
			"app_parameters": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Input parameters of the main class, that is application parameters.`,
			},
			"specification": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Compute resource type. Currently, resource types A, B, and C are available.`,
			},
			"jars": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Name of the package that is of the JAR type and has been uploaded to the DLI resource management system.`,
			},
			"python_files": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Name of the package that is of the PyFile type and has been uploaded to the DLI resource management system.`,
			},
			"files": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Name of the package that is of the file type and has been uploaded to the DLI resource management system.`,
			},
			"modules": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Name of the dependent system resource module.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        SparkTemplateResourceSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The list of resource objects.`,
			},
			"dependent_packages": {
				Type:        schema.TypeList,
				Elem:        SparkTemplateGroupSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The list of package resource objects.`,
			},
			"configurations": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The configuration items of the DLI spark.`,
			},
			"driver_memory": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Driver memory of the Spark application, for example, 2 GB and 2048 MB.`,
			},
			"driver_cores": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of CPU cores of the Spark application driver.`,
			},
			"executor_memory": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Executor memory of the Spark application, for example, 2 GB and 2048 MB.`,
			},
			"executor_cores": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of CPU cores of each Executor in the Spark application.`,
			},
			"num_executors": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of Executors in a Spark application.`,
			},
			"obs_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `OBS bucket for storing the Spark jobs.`,
			},
			"auto_recovery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the retry function.`,
			},
			"max_retry_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Maximum retry times.`,
			},
		},
	}
	return &sc
}

func SparkTemplateGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `User group name.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        SparkTemplateResourceSchema(),
				Optional:    true,
				Computed:    true,
				Description: `User group resource.`,
			},
		},
	}
	return &sc
}

func SparkTemplateResourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Resource name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Resource type.`,
			},
		},
	}
	return &sc
}

func resourceSparkTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSparkTemplate: create a Spark template.
	var (
		createSparkTemplateHttpUrl = "v3/{project_id}/templates"
		createSparkTemplateProduct = "dli"
	)
	createSparkTemplateClient, err := cfg.NewServiceClient(createSparkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createSparkTemplatePath := createSparkTemplateClient.Endpoint + createSparkTemplateHttpUrl
	createSparkTemplatePath = strings.ReplaceAll(createSparkTemplatePath, "{project_id}", createSparkTemplateClient.ProjectID)

	createSparkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createSparkTemplateOpt.JSONBody = utils.RemoveNil(buildSparkTemplateRequestBodyParams(d))
	createSparkTemplateResp, err := createSparkTemplateClient.Request("POST", createSparkTemplatePath, &createSparkTemplateOpt)
	if err != nil {
		return diag.Errorf("error creating DLI spark template: %s", err)
	}

	createSparkTemplateRespBody, err := utils.FlattenResponse(createSparkTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("id", createSparkTemplateRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("unable to find the DLI spark template ID from the API response")
	}
	d.SetId(templateId)

	return resourceSparkTemplateRead(ctx, d, meta)
}

func buildSparkTemplateRequestBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        "SPARK",
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"group":       utils.ValueIgnoreEmpty(d.Get("group")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"body":        buildSparkTemplateContent(d.Get("body")),
	}
	return bodyParams
}

func buildSparkTemplateContent(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"queue":           utils.ValueIgnoreEmpty(raw["queue_name"]),
			"name":            utils.ValueIgnoreEmpty(raw["name"]),
			"file":            utils.ValueIgnoreEmpty(raw["app_name"]),
			"className":       utils.ValueIgnoreEmpty(raw["main_class"]),
			"args":            utils.ValueIgnoreEmpty(raw["app_parameters"]),
			"sc_type":         utils.ValueIgnoreEmpty(raw["specification"]),
			"jars":            utils.ValueIgnoreEmpty(raw["jars"]),
			"pyFiles":         utils.ValueIgnoreEmpty(raw["python_files"]),
			"files":           utils.ValueIgnoreEmpty(raw["files"]),
			"modules":         utils.ValueIgnoreEmpty(raw["modules"]),
			"resources":       buildSparkTemplateResource(raw["resources"]),
			"groups":          buildSparkTemplateGroup(raw["dependent_packages"]),
			"conf":            utils.ValueIgnoreEmpty(raw["configurations"]),
			"driverMemory":    utils.ValueIgnoreEmpty(raw["driver_memory"]),
			"driverCores":     utils.ValueIgnoreEmpty(raw["driver_cores"]),
			"executorMemory":  utils.ValueIgnoreEmpty(raw["executor_memory"]),
			"executorCores":   utils.ValueIgnoreEmpty(raw["executor_cores"]),
			"numExecutors":    utils.ValueIgnoreEmpty(raw["num_executors"]),
			"obs_bucket":      utils.ValueIgnoreEmpty(raw["obs_bucket"]),
			"auto_recovery":   utils.ValueIgnoreEmpty(raw["auto_recovery"]),
			"max_retry_times": utils.ValueIgnoreEmpty(raw["max_retry_times"]),
		}
		return params
	}
	return nil
}

func buildSparkTemplateResource(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIgnoreEmpty(raw["name"]),
				"type": utils.ValueIgnoreEmpty(raw["type"]),
			}
		}
		return rst
	}
	return nil
}

func buildSparkTemplateGroup(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":      utils.ValueIgnoreEmpty(raw["name"]),
				"resources": buildSparkTemplateResource(raw["resources"]),
			}
		}
		return rst
	}
	return nil
}

func resourceSparkTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSparkTemplate: Query the Spark template.
	var (
		getSparkTemplateHttpUrl = "v3/{project_id}/templates/{id}"
		getSparkTemplateProduct = "dli"
	)
	getSparkTemplateClient, err := cfg.NewServiceClient(getSparkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getSparkTemplatePath := getSparkTemplateClient.Endpoint + getSparkTemplateHttpUrl
	getSparkTemplatePath = strings.ReplaceAll(getSparkTemplatePath, "{project_id}", getSparkTemplateClient.ProjectID)
	getSparkTemplatePath = strings.ReplaceAll(getSparkTemplatePath, "{id}", d.Id())

	getSparkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSparkTemplateResp, err := getSparkTemplateClient.Request("GET", getSparkTemplatePath, &getSparkTemplateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DLI spark template")
	}

	getSparkTemplateRespBody, err := utils.FlattenResponse(getSparkTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getSparkTemplateRespBody, nil)),
		d.Set("group", utils.PathSearch("group", getSparkTemplateRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getSparkTemplateRespBody, nil)),
		d.Set("body", flattenGetSparkTemplateResponseBody(getSparkTemplateRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSparkTemplateResponseBody(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("body", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"queue_name":         utils.PathSearch("queue", curJson, nil),
			"name":               utils.PathSearch("name", curJson, nil),
			"app_name":           utils.PathSearch("file", curJson, nil),
			"main_class":         utils.PathSearch("className", curJson, nil),
			"app_parameters":     utils.PathSearch("args", curJson, nil),
			"specification":      utils.PathSearch("sc_type", curJson, nil),
			"jars":               utils.PathSearch("jars", curJson, nil),
			"python_files":       utils.PathSearch("pyFiles", curJson, nil),
			"files":              utils.PathSearch("files", curJson, nil),
			"modules":            utils.PathSearch("modules", curJson, nil),
			"resources":          flattenSparkTemplateResources(curJson),
			"dependent_packages": flattenSparkTemplateGroups(curJson),
			"configurations":     utils.PathSearch("conf", curJson, nil),
			"driver_memory":      utils.PathSearch("driverMemory", curJson, nil),
			"driver_cores":       utils.PathSearch("driverCores", curJson, nil),
			"executor_memory":    utils.PathSearch("executorMemory", curJson, nil),
			"executor_cores":     utils.PathSearch("executorCores", curJson, nil),
			"num_executors":      utils.PathSearch("numExecutors", curJson, nil),
			"obs_bucket":         utils.PathSearch("obs_bucket", curJson, nil),
			"auto_recovery":      utils.PathSearch("auto_recovery", curJson, nil),
			"max_retry_times":    utils.PathSearch("max_retry_times", curJson, nil),
		},
	}
	return rst
}

func flattenSparkTemplateResources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"type": utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func flattenSparkTemplateGroups(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":      utils.PathSearch("name", v, nil),
			"resources": flattenSparkTemplateResources(v),
		})
	}
	return rst
}

func resourceSparkTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSparkTemplateChanges := []string{
		"type",
		"name",
		"group",
		"description",
		"body",
	}

	if d.HasChanges(updateSparkTemplateChanges...) {
		// updateSparkTemplate: update Spark template
		var (
			updateSparkTemplateHttpUrl = "v3/{project_id}/templates/{id}"
			updateSparkTemplateProduct = "dli"
		)
		updateSparkTemplateClient, err := cfg.NewServiceClient(updateSparkTemplateProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateSparkTemplatePath := updateSparkTemplateClient.Endpoint + updateSparkTemplateHttpUrl
		updateSparkTemplatePath = strings.ReplaceAll(updateSparkTemplatePath, "{project_id}", updateSparkTemplateClient.ProjectID)
		updateSparkTemplatePath = strings.ReplaceAll(updateSparkTemplatePath, "{id}", d.Id())

		updateSparkTemplateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateSparkTemplateOpt.JSONBody = utils.RemoveNil(buildSparkTemplateRequestBodyParams(d))
		_, err = updateSparkTemplateClient.Request("PUT", updateSparkTemplatePath, &updateSparkTemplateOpt)
		if err != nil {
			return diag.Errorf("error updating DLI spark template: %s", err)
		}
	}
	return resourceSparkTemplateRead(ctx, d, meta)
}

func resourceSparkTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSparkTemplate: delete Spark template
	var (
		deleteSparkTemplateHttpUrl = "v1.0/{project_id}/sqls-deletion"
		deleteSparkTemplateProduct = "dli"
	)
	deleteSparkTemplateClient, err := cfg.NewServiceClient(deleteSparkTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteSparkTemplatePath := deleteSparkTemplateClient.Endpoint + deleteSparkTemplateHttpUrl
	deleteSparkTemplatePath = strings.ReplaceAll(deleteSparkTemplatePath, "{project_id}", deleteSparkTemplateClient.ProjectID)

	deleteSparkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteSparkTemplateOpt.JSONBody = utils.RemoveNil(buildDeleteSparkTemplateBodyParams(d))
	_, err = deleteSparkTemplateClient.Request("POST", deleteSparkTemplatePath, &deleteSparkTemplateOpt)
	if err != nil {
		return diag.Errorf("error deleting DLI spark template: %s", err)
	}

	return nil
}

func buildDeleteSparkTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_ids": []string{d.Id()},
	}
	return bodyParams
}
