// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v1/{project_id}/models
// @API ModelArts DELETE /v1/{project_id}/models/{id}
// @API ModelArts GET /v1/{project_id}/models/{id}
func ResourceModelartsModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsModelCreate,
		ReadContext:   resourceModelartsModelRead,
		DeleteContext: resourceModelartsModelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
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
				ForceNew:    true,
				Description: `Model name, which consists of 1 to 64 characters.`,
			},
			"model_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Model type.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Model version in the format of Digit.Digit.Digit.`,
			},
			"source_location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `OBS path where the model is located or the SWR image location.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Model description that consists of 1 to 100 characters.`,
			},
			"model_docs": {
				Type:        schema.TypeList,
				Elem:        modelartsModelModelDocsSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `List of model description documents. A maximum of three documents are supported.`,
			},
			"template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     modelartsModelTemplateSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_copy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to enable image replication.`,
			},
			"execution_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `OBS path for storing the execution code.`,
			},
			"source_job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `ID of the source training job.`,
			},
			"source_job_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Version of the source training job.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Model source type`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"install_type": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Deployment type. Only lowercase letters are supported.`,
			},
			"initial_config": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return utils.JSONStringsEqual(old, new)
				},
				Description: `The model configuration file.`,
			},
			"model_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Model algorithm.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Model runtime environment.`,
			},
			"metrics": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Model precision.`,
			},
			"dependencies": {
				Type:        schema.TypeList,
				Elem:        modelartsModelDependencySchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Package required for inference code and model.`,
			},
			"schema_doc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Download address of the model schema file.`,
			},
			"image_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Image path generated after model packaging.`,
			},
			"model_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Model size, in bytes.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model status.`,
			},
			"model_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model source.`,
			},
			"tunable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a model can be tuned.`,
			},
			"market_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a model is subscribed from AI Gallery.`,
			},
			"publishable_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a model is subscribed from AI Gallery.`,
			},
		},
	}
}

func modelartsModelModelDocsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"doc_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `HTTP(S) link of the document.`,
			},
			"doc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Document name, which must start with a letter. Enter 1 to 48 characters.`,
			},
		},
	}
	return &sc
}

func modelartsModelTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"infer_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `ID of the input and output mode.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the used template.`,
			},
			"template_inputs": {
				Type:        schema.TypeList,
				Elem:        modelartsModelTemplateTemplateInputSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Template input configuration, specifying the source path for configuring a model.`,
			},
		},
	}
	return &sc
}

func modelartsModelTemplateTemplateInputSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"input": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Template input path, which can be a path to an OBS file or directory.`,
			},
			"input_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Input item ID, which is obtained from template details.`,
			},
		},
	}
	return &sc
}

func modelartsModelDependencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"installer": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Installation mode. Only **pip** is supported.`,
			},
			"packages": {
				Type:        schema.TypeList,
				Elem:        modelartsModelDependencyPackageSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Collection of dependency packages.`,
			},
		},
	}
	return &sc
}

func modelartsModelDependencyPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"package_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Version of a dependency package.`,
			},
			"package_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Name of a dependency package.`,
			},
			"restraint": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Version restriction, which can be **EXACT**, **ATLEAST**, or **ATMOST**.`,
			},
		},
	}
	return &sc
}

func resourceModelartsModelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createModel: create a Modelarts model.
	var (
		createModelHttpUrl = "v1/{project_id}/models"
		createModelProduct = "modelarts"
	)
	createModelClient, err := cfg.NewServiceClient(createModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	createModelPath := createModelClient.Endpoint + createModelHttpUrl
	createModelPath = strings.ReplaceAll(createModelPath, "{project_id}", createModelClient.ProjectID)

	createModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	createModelOpt.JSONBody = utils.RemoveNil(buildCreateModelBodyParams(d))
	createModelResp, err := createModelClient.Request("POST", createModelPath, &createModelOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts model: %s", err)
	}

	createModelRespBody, err := utils.FlattenResponse(createModelResp)
	if err != nil {
		return diag.FromErr(err)
	}

	modelId := utils.PathSearch("model_id", createModelRespBody, "").(string)
	if modelId == "" {
		return diag.Errorf("unable to find the ModelArts model ID from the API response")
	}
	d.SetId(modelId)

	err = createModelWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts model (%s) creation to complete: %s", d.Id(), err)
	}

	return resourceModelartsModelRead(ctx, d, meta)
}

func buildCreateModelBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"model_docs":         buildCreateModelReqBodyModelDocs(d.Get("model_docs")),
		"template":           buildCreateModelReqBodyTemplate(d.Get("template")),
		"model_version":      utils.ValueIgnoreEmpty(d.Get("version")),
		"source_job_version": utils.ValueIgnoreEmpty(d.Get("source_job_version")),
		"source_location":    utils.ValueIgnoreEmpty(d.Get("source_location")),
		"source_copy":        utils.ValueIgnoreEmpty(d.Get("source_copy")),
		"initial_config":     utils.ValueIgnoreEmpty(d.Get("initial_config")),
		"execution_code":     utils.ValueIgnoreEmpty(d.Get("execution_code")),
		"source_job_id":      utils.ValueIgnoreEmpty(d.Get("source_job_id")),
		"model_type":         utils.ValueIgnoreEmpty(d.Get("model_type")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"runtime":            utils.ValueIgnoreEmpty(d.Get("runtime")),
		"model_metrics":      utils.ValueIgnoreEmpty(d.Get("metrics")),
		"source_type":        utils.ValueIgnoreEmpty(d.Get("source_type")),
		"dependencies":       buildCreateModelReqBodyDependency(d.Get("dependencies")),
		"workspace_id":       utils.ValueIgnoreEmpty(d.Get("workspace_id")),
		"model_algorithm":    utils.ValueIgnoreEmpty(d.Get("model_algorithm")),
		"model_name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"install_type":       utils.ValueIgnoreEmpty(d.Get("install_type")),
	}
	return bodyParams
}

func buildCreateModelReqBodyModelDocs(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"doc_url":  utils.ValueIgnoreEmpty(raw["doc_url"]),
				"doc_name": utils.ValueIgnoreEmpty(raw["doc_name"]),
			}
		}
		return rst
	}
	return nil
}

func buildCreateModelReqBodyTemplate(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"infer_format":    utils.ValueIgnoreEmpty(raw["infer_format"]),
			"template_id":     utils.ValueIgnoreEmpty(raw["template_id"]),
			"template_inputs": buildTemplateTemplateInput(raw["template_inputs"]),
		}
		return params
	}
	return nil
}

func buildTemplateTemplateInput(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"input":    utils.ValueIgnoreEmpty(raw["input"]),
				"input_id": utils.ValueIgnoreEmpty(raw["input_id"]),
			}
		}
		return rst
	}
	return nil
}

func buildCreateModelReqBodyDependency(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"installer": utils.ValueIgnoreEmpty(raw["installer"]),
				"packages":  buildDependencypackage(raw["packages"]),
			}
		}
		return rst
	}
	return nil
}

func buildDependencypackage(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"package_version": utils.ValueIgnoreEmpty(raw["package_version"]),
				"package_name":    utils.ValueIgnoreEmpty(raw["package_name"]),
				"restraint":       utils.ValueIgnoreEmpty(raw["restraint"]),
			}
		}
		return rst
	}
	return nil
}

func createModelWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createModelWaitingHttpUrl = "v1/{project_id}/models/{id}"
				createModelWaitingProduct = "modelarts"
			)
			createModelWaitingClient, err := cfg.NewServiceClient(createModelWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating ModelArts Client: %s", err)
			}

			createModelWaitingPath := createModelWaitingClient.Endpoint + createModelWaitingHttpUrl
			createModelWaitingPath = strings.ReplaceAll(createModelWaitingPath, "{project_id}", createModelWaitingClient.ProjectID)
			createModelWaitingPath = strings.ReplaceAll(createModelWaitingPath, "{id}", d.Id())

			createModelWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}

			createModelWaitingResp, err := createModelWaitingClient.Request("GET", createModelWaitingPath, &createModelWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createModelWaitingRespBody, err := utils.FlattenResponse(createModelWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`model_status`, createModelWaitingRespBody, "").(string)

			targetStatus := []string{
				"published",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createModelWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"publishing",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createModelWaitingRespBody, "PENDING", nil
			}

			return createModelWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsModelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getModel: Query the Modelarts model.
	var (
		getModelHttpUrl = "v1/{project_id}/models/{id}"
		getModelProduct = "modelarts"
	)
	getModelClient, err := cfg.NewServiceClient(getModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	getModelPath := getModelClient.Endpoint + getModelHttpUrl
	getModelPath = strings.ReplaceAll(getModelPath, "{project_id}", getModelClient.ProjectID)
	getModelPath = strings.ReplaceAll(getModelPath, "{id}", d.Id())

	getModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getModelResp, err := getModelClient.Request("GET", getModelPath, &getModelOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts model")
	}

	getModelRespBody, err := utils.FlattenResponse(getModelResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("model_docs", flattenGetModelResponseBodyModelDocs(getModelRespBody)),
		d.Set("version", utils.PathSearch("model_version", getModelRespBody, nil)),
		d.Set("source_job_version", utils.PathSearch("source_job_version", getModelRespBody, nil)),
		d.Set("source_location", utils.PathSearch("source_location", getModelRespBody, nil)),
		d.Set("source_copy", fmt.Sprintf("%v", utils.PathSearch("source_copy", getModelRespBody, nil))),
		d.Set("execution_code", utils.PathSearch("execution_code", getModelRespBody, nil)),
		d.Set("source_job_id", utils.PathSearch("source_job_id", getModelRespBody, nil)),
		d.Set("model_type", utils.PathSearch("model_type", getModelRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getModelRespBody, nil)),
		d.Set("runtime", utils.PathSearch("runtime", getModelRespBody, nil)),
		d.Set("metrics", utils.PathSearch("model_metrics", getModelRespBody, nil)),
		d.Set("source_type", utils.PathSearch("source_type", getModelRespBody, nil)),
		d.Set("dependencies", flattenGetModelResponseBodyDependency(getModelRespBody)),
		d.Set("workspace_id", utils.PathSearch("workspace_id", getModelRespBody, nil)),
		d.Set("model_algorithm", utils.PathSearch("model_algorithm", getModelRespBody, nil)),
		d.Set("name", utils.PathSearch("model_name", getModelRespBody, nil)),
		d.Set("install_type", utils.PathSearch("install_type", getModelRespBody, nil)),
		d.Set("schema_doc", utils.PathSearch("schema_doc", getModelRespBody, nil)),
		d.Set("image_address", utils.PathSearch("image_address", getModelRespBody, nil)),
		d.Set("model_size", utils.PathSearch("model_size", getModelRespBody, nil)),
		d.Set("status", utils.PathSearch("model_status", getModelRespBody, nil)),
		d.Set("model_source", utils.PathSearch("model_source", getModelRespBody, nil)),
		d.Set("tunable", utils.PathSearch("tunable", getModelRespBody, nil)),
		d.Set("market_flag", utils.PathSearch("market_flag", getModelRespBody, nil)),
		d.Set("publishable_flag", utils.PathSearch("publishable_flag", getModelRespBody, nil)),
		d.Set("initial_config", utils.PathSearch("config", getModelRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetModelResponseBodyModelDocs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("model_docs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"doc_url":  utils.PathSearch("doc_url", v, nil),
			"doc_name": utils.PathSearch("doc_name", v, nil),
		})
	}
	return rst
}

func flattenGetModelResponseBodyDependency(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("dependencies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"installer": utils.PathSearch("installer", v, nil),
			"packages":  flattenDependencyPackages(v),
		})
	}
	return rst
}

func flattenDependencyPackages(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("packages", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"package_version": utils.PathSearch("package_version", v, nil),
			"package_name":    utils.PathSearch("package_name", v, nil),
			"restraint":       utils.PathSearch("restraint", v, nil),
		})
	}
	return rst
}

func resourceModelartsModelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteModel: delete Modelarts model
	var (
		deleteModelHttpUrl = "v1/{project_id}/models/{id}"
		deleteModelProduct = "modelarts"
	)
	deleteModelClient, err := cfg.NewServiceClient(deleteModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	deleteModelPath := deleteModelClient.Endpoint + deleteModelHttpUrl
	deleteModelPath = strings.ReplaceAll(deleteModelPath, "{project_id}", deleteModelClient.ProjectID)
	deleteModelPath = strings.ReplaceAll(deleteModelPath, "{id}", d.Id())

	deleteModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	_, err = deleteModelClient.Request("DELETE", deleteModelPath, &deleteModelOpt)
	if err != nil {
		return diag.Errorf("error deleting Modelarts model: %s", err)
	}

	return nil
}
